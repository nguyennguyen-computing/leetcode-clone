package execution

import (
	"testing"

	"leetcode-clone-backend/pkg/models"
)

func TestExecutionService_ValidateCode(t *testing.T) {
	es := NewExecutionService()

	tests := []struct {
		name     string
		code     string
		language string
		wantErr  bool
	}{
		{
			name:     "Valid JavaScript code",
			code:     "function solution(input) { return input; }",
			language: models.LanguageJavaScript,
			wantErr:  false,
		},
		{
			name:     "Valid Python code",
			code:     "def solution(input):\n    return input",
			language: models.LanguagePython,
			wantErr:  false,
		},
		{
			name:     "Valid Java code",
			code:     "public String solution(String input) { return input; }",
			language: models.LanguageJava,
			wantErr:  false,
		},
		{
			name:     "Dangerous JavaScript code - fs access",
			code:     "const fs = require('fs'); function solution(input) { return input; }",
			language: models.LanguageJavaScript,
			wantErr:  true,
		},
		{
			name:     "Dangerous Python code - os import",
			code:     "import os\ndef solution(input):\n    return input",
			language: models.LanguagePython,
			wantErr:  true,
		},
		{
			name:     "Dangerous Java code - Runtime access",
			code:     "Runtime.getRuntime().exec(\"ls\"); public String solution(String input) { return input; }",
			language: models.LanguageJava,
			wantErr:  true,
		},
		{
			name:     "Code too long",
			code:     string(make([]byte, 60000)), // Exceeds 50KB limit
			language: models.LanguageJavaScript,
			wantErr:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := es.validateCode(tt.code, tt.language)
			if (err != nil) != tt.wantErr {
				t.Errorf("ExecutionService.validateCode() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestExecutionService_IsLanguageSupported(t *testing.T) {
	es := NewExecutionService()

	tests := []struct {
		language string
		expected bool
	}{
		{models.LanguageJavaScript, true},
		{models.LanguagePython, true},
		{models.LanguageJava, true},
		{"cpp", false},
		{"go", false},
		{"", false},
	}

	for _, tt := range tests {
		t.Run(tt.language, func(t *testing.T) {
			result := es.isLanguageSupported(tt.language)
			if result != tt.expected {
				t.Errorf("ExecutionService.isLanguageSupported(%s) = %v, want %v", tt.language, result, tt.expected)
			}
		})
	}
}
func TestExecutionService_WrapCode(t *testing.T) {
	es := NewExecutionService()

	tests := []struct {
		name     string
		code     string
		language string
	}{
		{
			name:     "JavaScript code wrapping",
			code:     "function solution(input) { return input.trim(); }",
			language: models.LanguageJavaScript,
		},
		{
			name:     "Python code wrapping",
			code:     "def solution(input_data):\n    return input_data.strip()",
			language: models.LanguagePython,
		},
		{
			name:     "Java code wrapping",
			code:     "public String solution(String input) { return input.trim(); }",
			language: models.LanguageJava,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var wrapped string
			switch tt.language {
			case models.LanguageJavaScript:
				wrapped = es.wrapJavaScriptCode(tt.code)
			case models.LanguagePython:
				wrapped = es.wrapPythonCode(tt.code)
			case models.LanguageJava:
				wrapped = es.wrapJavaCode(tt.code)
			}

			// Check that the wrapped code contains the original code
			if !contains(wrapped, tt.code) {
				t.Errorf("Wrapped code does not contain original code")
			}

			// Check that the wrapped code contains input reading logic
			if !contains(wrapped, "input") {
				t.Errorf("Wrapped code does not contain input reading logic")
			}
		})
	}
}

func TestExecutionService_BuildDockerCommand(t *testing.T) {
	es := NewExecutionService()
	execDir := "/tmp/test"
	codeFile := "/tmp/test/solution.js"

	tests := []struct {
		name     string
		language string
		expected []string
	}{
		{
			name:     "JavaScript Docker command",
			language: models.LanguageJavaScript,
			expected: []string{"run", "--rm", "--network=none", "--read-only"},
		},
		{
			name:     "Python Docker command",
			language: models.LanguagePython,
			expected: []string{"run", "--rm", "--network=none", "--read-only"},
		},
		{
			name:     "Java Docker command",
			language: models.LanguageJava,
			expected: []string{"run", "--rm", "--network=none", "--read-only"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cmd := es.buildDockerCommand(execDir, codeFile, tt.language)

			// Check that all expected security flags are present
			for _, expected := range tt.expected {
				if !containsString(cmd, expected) {
					t.Errorf("Docker command missing expected flag: %s", expected)
				}
			}

			// Check that memory and CPU limits are set
			if !containsString(cmd, "--memory=") {
				t.Errorf("Docker command missing memory limit")
			}
			if !containsString(cmd, "--cpus=") {
				t.Errorf("Docker command missing CPU limit")
			}
		})
	}
}

// Helper functions for tests
func contains(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || len(s) > len(substr) && (s[:len(substr)] == substr || s[len(s)-len(substr):] == substr || containsSubstring(s, substr)))
}

func containsSubstring(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}

func containsString(slice []string, item string) bool {
	for _, s := range slice {
		if s == item || containsSubstring(s, item) {
			return true
		}
	}
	return false
}
