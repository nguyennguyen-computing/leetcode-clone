package execution

import (
	"context"
	"fmt"
	"leetcode-clone-backend/pkg/models"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"
)

// ExecutionResult represents the result of code execution
type ExecutionResult struct {
	Status          string       `json:"status"`
	Output          string       `json:"output"`
	ErrorMessage    string       `json:"error_message,omitempty"`
	RuntimeMs       int          `json:"runtime_ms"`
	MemoryKb        int          `json:"memory_kb"`
	TestCasesPassed int          `json:"test_cases_passed"`
	TotalTestCases  int          `json:"total_test_cases"`
	TestResults     []TestResult `json:"test_results,omitempty"`
}

// TestResult represents the result of a single test case
type TestResult struct {
	Input          string `json:"input"`
	ExpectedOutput string `json:"expected_output"`
	ActualOutput   string `json:"actual_output"`
	Passed         bool   `json:"passed"`
	RuntimeMs      int    `json:"runtime_ms"`
	MemoryKb       int    `json:"memory_kb"`
}

// ExecutionServiceInterface defines the interface for code execution
type ExecutionServiceInterface interface {
	ExecuteCode(code, language string, testCases []models.TestCase) (*ExecutionResult, error)
	ValidateCode(code, language string) error
}

// ExecutionService handles code execution in sandboxed environments
type ExecutionService struct {
	dockerImage    string
	timeoutSeconds int
	memoryLimitMB  int
	tempDir        string
}

// NewExecutionService creates a new execution service
func NewExecutionService() *ExecutionService {
	return &ExecutionService{
		dockerImage:    "ubuntu:22.04",
		timeoutSeconds: 10,  // 10 seconds timeout
		memoryLimitMB:  128, // 128MB memory limit
		tempDir:        "/tmp/leetcode-execution",
	}
}

// ExecuteCode runs the provided code against test cases in a sandboxed environment
func (es *ExecutionService) ExecuteCode(code, language string, testCases []models.TestCase) (*ExecutionResult, error) {
	// Validate language support
	if !es.isLanguageSupported(language) {
		return &ExecutionResult{
			Status:       models.StatusInternalError,
			ErrorMessage: fmt.Sprintf("Unsupported language: %s", language),
		}, nil
	}

	// Create temporary directory for this execution
	execDir, err := es.createTempDir()
	if err != nil {
		return &ExecutionResult{
			Status:       models.StatusInternalError,
			ErrorMessage: "Failed to create temporary directory",
		}, err
	}
	defer os.RemoveAll(execDir)

	// Prepare code file
	codeFile, err := es.prepareCodeFile(execDir, code, language)
	if err != nil {
		return &ExecutionResult{
			Status:       models.StatusInternalError,
			ErrorMessage: "Failed to prepare code file",
		}, err
	}

	// Execute against test cases
	result := &ExecutionResult{
		TotalTestCases: len(testCases),
		TestResults:    make([]TestResult, 0, len(testCases)),
	}

	totalRuntime := 0
	maxMemory := 0

	for _, testCase := range testCases {
		testResult, err := es.executeTestCase(execDir, codeFile, language, testCase)
		if err != nil {
			result.Status = models.StatusInternalError
			result.ErrorMessage = err.Error()
			return result, nil
		}

		result.TestResults = append(result.TestResults, *testResult)
		totalRuntime += testResult.RuntimeMs
		if testResult.MemoryKb > maxMemory {
			maxMemory = testResult.MemoryKb
		}

		if testResult.Passed {
			result.TestCasesPassed++
		} else {
			// If any test case fails, determine the failure reason
			if strings.Contains(testResult.ActualOutput, "timeout") {
				result.Status = models.StatusTimeLimitExceeded
			} else if strings.Contains(testResult.ActualOutput, "memory") {
				result.Status = models.StatusMemoryLimitExceeded
			} else if strings.Contains(testResult.ActualOutput, "error") || strings.Contains(testResult.ActualOutput, "Error") {
				result.Status = models.StatusRuntimeError
			} else {
				result.Status = models.StatusWrongAnswer
			}
			result.ErrorMessage = fmt.Sprintf("Test case failed: expected %s, got %s", testResult.ExpectedOutput, testResult.ActualOutput)
			break
		}
	}

	// If all test cases passed
	if result.TestCasesPassed == result.TotalTestCases {
		result.Status = models.StatusAccepted
	}

	result.RuntimeMs = totalRuntime / len(testCases) // Average runtime
	result.MemoryKb = maxMemory

	return result, nil
}

// executeTestCase runs a single test case
func (es *ExecutionService) executeTestCase(execDir, codeFile, language string, testCase models.TestCase) (*TestResult, error) {
	start := time.Now()

	// Create input file
	inputFile := filepath.Join(execDir, "input.txt")
	if err := os.WriteFile(inputFile, []byte(testCase.Input), 0644); err != nil {
		return nil, fmt.Errorf("failed to create input file: %v", err)
	}

	// Prepare Docker command
	dockerCmd := es.buildDockerCommand(execDir, codeFile, language)

	// Create context with timeout
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(es.timeoutSeconds)*time.Second)
	defer cancel()

	// Execute in Docker container
	cmd := exec.CommandContext(ctx, "docker", dockerCmd...)
	output, err := cmd.CombinedOutput()

	runtime := int(time.Since(start).Milliseconds())

	result := &TestResult{
		Input:          testCase.Input,
		ExpectedOutput: strings.TrimSpace(testCase.ExpectedOutput),
		ActualOutput:   strings.TrimSpace(string(output)),
		RuntimeMs:      runtime,
		MemoryKb:       es.estimateMemoryUsage(string(output)), // Simple estimation
	}

	// Check for timeout
	if ctx.Err() == context.DeadlineExceeded {
		result.ActualOutput = "timeout"
		result.Passed = false
		return result, nil
	}

	// Check for execution errors
	if err != nil {
		result.ActualOutput = fmt.Sprintf("error: %v", err)
		result.Passed = false
		return result, nil
	}

	// Compare outputs
	result.Passed = result.ExpectedOutput == result.ActualOutput

	return result, nil
}

// buildDockerCommand constructs the Docker command for code execution
func (es *ExecutionService) buildDockerCommand(execDir, codeFile, language string) []string {
	baseCmd := []string{
		"run",
		"--rm",
		"--network=none",                            // No network access
		"--read-only",                               // Read-only filesystem
		"--tmpfs", "/tmp:rw,noexec,nosuid,size=10m", // Limited temp space
		fmt.Sprintf("--memory=%dm", es.memoryLimitMB), // Memory limit
		fmt.Sprintf("--cpus=0.5"),                     // CPU limit
		"--user", "nobody",                            // Run as nobody user
		"-v", fmt.Sprintf("%s:/workspace:ro", execDir), // Mount code directory as read-only
		"-w", "/workspace",
	}

	switch language {
	case models.LanguageJavaScript:
		return append(baseCmd, "node:18-alpine", "timeout", fmt.Sprintf("%ds", es.timeoutSeconds), "node", filepath.Base(codeFile))
	case models.LanguagePython:
		return append(baseCmd, "python:3.11-alpine", "timeout", fmt.Sprintf("%ds", es.timeoutSeconds), "python3", filepath.Base(codeFile))
	case models.LanguageJava:
		className := strings.TrimSuffix(filepath.Base(codeFile), ".java")
		return append(baseCmd, "openjdk:17-alpine", "sh", "-c",
			fmt.Sprintf("javac %s && timeout %ds java %s", filepath.Base(codeFile), es.timeoutSeconds, className))
	default:
		return append(baseCmd, "alpine:latest", "echo", "Unsupported language")
	}
}

// prepareCodeFile creates the code file with proper extension and security measures
func (es *ExecutionService) prepareCodeFile(execDir, code, language string) (string, error) {
	// Security: Validate and sanitize code
	if err := es.validateCode(code, language); err != nil {
		return "", err
	}

	var filename string
	var finalCode string

	switch language {
	case models.LanguageJavaScript:
		filename = "solution.js"
		finalCode = es.wrapJavaScriptCode(code)
	case models.LanguagePython:
		filename = "solution.py"
		finalCode = es.wrapPythonCode(code)
	case models.LanguageJava:
		filename = "Solution.java"
		finalCode = es.wrapJavaCode(code)
	default:
		return "", fmt.Errorf("unsupported language: %s", language)
	}

	codeFile := filepath.Join(execDir, filename)
	if err := os.WriteFile(codeFile, []byte(finalCode), 0644); err != nil {
		return "", fmt.Errorf("failed to write code file: %v", err)
	}

	return codeFile, nil
}

// ValidateCode performs security validation on the submitted code (exported for handlers)
func (es *ExecutionService) ValidateCode(code, language string) error {
	return es.validateCode(code, language)
}

// validateCode performs security validation on the submitted code
func (es *ExecutionService) validateCode(code, language string) error {
	// Check for dangerous patterns
	dangerousPatterns := []string{
		"import os", "import sys", "import subprocess", "import socket",
		"exec(", "eval(", "__import__", "open(", "file(",
		"System.exit", "Runtime.getRuntime", "ProcessBuilder",
		"require('fs')", "require('os')", "require('child_process')",
		"process.exit", "require('net')", "require('http')",
	}

	codeLower := strings.ToLower(code)
	for _, pattern := range dangerousPatterns {
		if strings.Contains(codeLower, strings.ToLower(pattern)) {
			return fmt.Errorf("code contains potentially dangerous pattern: %s", pattern)
		}
	}

	// Check code length (prevent extremely large submissions)
	if len(code) > 50000 { // 50KB limit
		return fmt.Errorf("code exceeds maximum length limit")
	}

	return nil
}

// wrapJavaScriptCode wraps user code with input/output handling
func (es *ExecutionService) wrapJavaScriptCode(userCode string) string {
	return fmt.Sprintf(`
const fs = require('fs');

// Read input
const input = fs.readFileSync('/workspace/input.txt', 'utf8').trim();

// User's solution code
%s

// Execute and output result
try {
    const result = solution(input);
    console.log(result);
} catch (error) {
    console.error('Runtime Error:', error.message);
}
`, userCode)
}

// wrapPythonCode wraps user code with input/output handling
func (es *ExecutionService) wrapPythonCode(userCode string) string {
	return fmt.Sprintf(`
import sys

# Read input
with open('/workspace/input.txt', 'r') as f:
    input_data = f.read().strip()

# User's solution code
%s

# Execute and output result
try:
    result = solution(input_data)
    print(result)
except Exception as error:
    print(f'Runtime Error: {error}', file=sys.stderr)
`, userCode)
}

// wrapJavaCode wraps user code with input/output handling
func (es *ExecutionService) wrapJavaCode(userCode string) string {
	return fmt.Sprintf(`
import java.io.*;
import java.util.*;

public class Solution {
    %s
    
    public static void main(String[] args) {
        try {
            Scanner scanner = new Scanner(new File("/workspace/input.txt"));
            String input = scanner.nextLine();
            scanner.close();
            
            Solution sol = new Solution();
            String result = sol.solution(input);
            System.out.println(result);
        } catch (Exception error) {
            System.err.println("Runtime Error: " + error.getMessage());
        }
    }
}
`, userCode)
}

// Helper methods
func (es *ExecutionService) isLanguageSupported(language string) bool {
	supportedLanguages := []string{
		models.LanguageJavaScript,
		models.LanguagePython,
		models.LanguageJava,
	}

	for _, supported := range supportedLanguages {
		if language == supported {
			return true
		}
	}
	return false
}

func (es *ExecutionService) createTempDir() (string, error) {
	// Ensure base temp directory exists
	if err := os.MkdirAll(es.tempDir, 0755); err != nil {
		return "", err
	}

	// Create unique directory for this execution
	return os.MkdirTemp(es.tempDir, "exec-*")
}

func (es *ExecutionService) estimateMemoryUsage(output string) int {
	// Simple memory estimation based on output length
	// In a real implementation, you'd use Docker stats API
	baseMemory := 1024                // 1MB base
	outputMemory := len(output) / 100 // Rough estimation
	return baseMemory + outputMemory
}
