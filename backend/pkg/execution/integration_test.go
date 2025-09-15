//go:build integration
// +build integration

package execution

import (
	"os/exec"
	"testing"

	"leetcode-clone-backend/pkg/models"
)

// TestExecutionService_Integration tests the execution service with actual Docker
// Run with: go test -tags=integration ./pkg/execution/...
func TestExecutionService_Integration(t *testing.T) {
	// Check if Docker is available
	if _, err := exec.LookPath("docker"); err != nil {
		t.Skip("Docker not available, skipping integration test")
	}

	es := NewExecutionService()

	testCases := []models.TestCase{
		{
			Input:          "hello",
			ExpectedOutput: "hello",
		},
	}

	tests := []struct {
		name     string
		code     string
		language string
		wantPass bool
	}{
		{
			name:     "JavaScript simple echo",
			code:     "function solution(input) { return input.trim(); }",
			language: models.LanguageJavaScript,
			wantPass: true,
		},
		{
			name:     "Python simple echo",
			code:     "def solution(input_data):\n    return input_data.strip()",
			language: models.LanguagePython,
			wantPass: true,
		},
		{
			name:     "Java simple echo",
			code:     "public String solution(String input) { return input.trim(); }",
			language: models.LanguageJava,
			wantPass: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := es.ExecuteCode(tt.code, tt.language, testCases)
			if err != nil {
				t.Fatalf("ExecuteCode() error = %v", err)
			}

			if tt.wantPass && result.Status != models.StatusAccepted {
				t.Errorf("Expected status %s, got %s. Error: %s", models.StatusAccepted, result.Status, result.ErrorMessage)
			}

			if tt.wantPass && result.TestCasesPassed != len(testCases) {
				t.Errorf("Expected %d test cases to pass, got %d", len(testCases), result.TestCasesPassed)
			}
		})
	}
}

// TestExecutionService_SecurityValidation tests security measures
func TestExecutionService_SecurityValidation(t *testing.T) {
	// Check if Docker is available
	if _, err := exec.LookPath("docker"); err != nil {
		t.Skip("Docker not available, skipping integration test")
	}

	es := NewExecutionService()

	testCases := []models.TestCase{
		{
			Input:          "test",
			ExpectedOutput: "test",
		},
	}

	maliciousCodes := []struct {
		name     string
		code     string
		language string
	}{
		{
			name:     "JavaScript file system access",
			code:     "const fs = require('fs'); function solution(input) { fs.readFileSync('/etc/passwd'); return input; }",
			language: models.LanguageJavaScript,
		},
		{
			name:     "Python os commands",
			code:     "import os\ndef solution(input_data):\n    os.system('ls')\n    return input_data",
			language: models.LanguagePython,
		},
	}

	for _, tt := range maliciousCodes {
		t.Run(tt.name, func(t *testing.T) {
			// Should fail at validation stage
			err := es.ValidateCode(tt.code, tt.language)
			if err == nil {
				t.Error("Expected validation to fail for malicious code")
			}
		})
	}
}
