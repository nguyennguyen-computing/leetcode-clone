package models

import (
	"encoding/json"
	"testing"
	"time"
)

func TestExamples_MarshalUnmarshal(t *testing.T) {
	examples := Examples{
		{
			Input:       "[2,7,11,15], 9",
			Output:      "[0,1]",
			Explanation: "Because nums[0] + nums[1] == 9, we return [0, 1].",
		},
		{
			Input:  "[3,2,4], 6",
			Output: "[1,2]",
		},
	}

	// Test marshaling
	data, err := json.Marshal(examples)
	if err != nil {
		t.Fatalf("Failed to marshal examples: %v", err)
	}

	// Test unmarshaling
	var unmarshaled Examples
	err = json.Unmarshal(data, &unmarshaled)
	if err != nil {
		t.Fatalf("Failed to unmarshal examples: %v", err)
	}

	if len(unmarshaled) != len(examples) {
		t.Errorf("Expected %d examples, got %d", len(examples), len(unmarshaled))
	}

	if unmarshaled[0].Input != examples[0].Input {
		t.Errorf("Expected input %s, got %s", examples[0].Input, unmarshaled[0].Input)
	}
}

func TestTemplateCode_MarshalUnmarshal(t *testing.T) {
	templateCode := TemplateCode{
		"javascript": "function twoSum(nums, target) {\n    // Your code here\n}",
		"python":     "def two_sum(nums, target):\n    # Your code here\n    pass",
		"java":       "public int[] twoSum(int[] nums, int target) {\n    // Your code here\n}",
	}

	// Test marshaling
	data, err := json.Marshal(templateCode)
	if err != nil {
		t.Fatalf("Failed to marshal template code: %v", err)
	}

	// Test unmarshaling
	var unmarshaled TemplateCode
	err = json.Unmarshal(data, &unmarshaled)
	if err != nil {
		t.Fatalf("Failed to unmarshal template code: %v", err)
	}

	if len(unmarshaled) != len(templateCode) {
		t.Errorf("Expected %d templates, got %d", len(templateCode), len(unmarshaled))
	}

	if unmarshaled["javascript"] != templateCode["javascript"] {
		t.Errorf("Expected JavaScript template to match")
	}
}

func TestUser_JSONSerialization(t *testing.T) {
	user := &User{
		ID:           1,
		Username:     "testuser",
		Email:        "test@example.com",
		PasswordHash: "hashedpassword",
		IsAdmin:      false,
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}

	// Test marshaling
	data, err := json.Marshal(user)
	if err != nil {
		t.Fatalf("Failed to marshal user: %v", err)
	}

	// Verify password hash is not included in JSON
	var jsonMap map[string]interface{}
	err = json.Unmarshal(data, &jsonMap)
	if err != nil {
		t.Fatalf("Failed to unmarshal JSON: %v", err)
	}

	if _, exists := jsonMap["password_hash"]; exists {
		t.Error("Password hash should not be included in JSON output")
	}

	if jsonMap["username"] != user.Username {
		t.Errorf("Expected username %s, got %v", user.Username, jsonMap["username"])
	}
}

func TestProblem_JSONSerialization(t *testing.T) {
	problem := &Problem{
		ID:          1,
		Title:       "Two Sum",
		Slug:        "two-sum",
		Description: "Given an array of integers nums and an integer target...",
		Difficulty:  DifficultyEasy,
		Tags:        []string{"Array", "Hash Table"},
		Examples: Examples{
			{
				Input:       "[2,7,11,15], 9",
				Output:      "[0,1]",
				Explanation: "Because nums[0] + nums[1] == 9, we return [0, 1].",
			},
		},
		Constraints: "2 <= nums.length <= 10^4",
		TemplateCode: TemplateCode{
			"javascript": "function twoSum(nums, target) {\n    // Your code here\n}",
		},
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	// Test marshaling
	data, err := json.Marshal(problem)
	if err != nil {
		t.Fatalf("Failed to marshal problem: %v", err)
	}

	// Test unmarshaling
	var unmarshaled Problem
	err = json.Unmarshal(data, &unmarshaled)
	if err != nil {
		t.Fatalf("Failed to unmarshal problem: %v", err)
	}

	if unmarshaled.Title != problem.Title {
		t.Errorf("Expected title %s, got %s", problem.Title, unmarshaled.Title)
	}

	if len(unmarshaled.Examples) != len(problem.Examples) {
		t.Errorf("Expected %d examples, got %d", len(problem.Examples), len(unmarshaled.Examples))
	}

	if len(unmarshaled.TemplateCode) != len(problem.TemplateCode) {
		t.Errorf("Expected %d template codes, got %d", len(problem.TemplateCode), len(unmarshaled.TemplateCode))
	}
}

func TestSubmission_StatusConstants(t *testing.T) {
	validStatuses := []string{
		StatusAccepted,
		StatusWrongAnswer,
		StatusTimeLimitExceeded,
		StatusMemoryLimitExceeded,
		StatusRuntimeError,
		StatusCompileError,
		StatusInternalError,
	}

	for _, status := range validStatuses {
		if status == "" {
			t.Error("Status constant should not be empty")
		}
	}
}

func TestProblem_DifficultyConstants(t *testing.T) {
	validDifficulties := []string{
		DifficultyEasy,
		DifficultyMedium,
		DifficultyHard,
	}

	for _, difficulty := range validDifficulties {
		if difficulty == "" {
			t.Error("Difficulty constant should not be empty")
		}
	}
}

func TestProblem_LanguageConstants(t *testing.T) {
	validLanguages := []string{
		LanguageJavaScript,
		LanguagePython,
		LanguageJava,
	}

	for _, language := range validLanguages {
		if language == "" {
			t.Error("Language constant should not be empty")
		}
	}
}