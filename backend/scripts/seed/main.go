package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/lib/pq"
	_ "github.com/lib/pq"
	"golang.org/x/crypto/bcrypt"
)

// Database connection configuration
type Config struct {
	Host     string
	Port     string
	User     string
	Password string
	DBName   string
}

// Models for seeding data
type User struct {
	ID           int       `db:"id"`
	Username     string    `db:"username"`
	Email        string    `db:"email"`
	PasswordHash string    `db:"password_hash"`
	IsAdmin      bool      `db:"is_admin"`
	CreatedAt    time.Time `db:"created_at"`
	UpdatedAt    time.Time `db:"updated_at"`
}

type Example struct {
	Input       string `json:"input"`
	Output      string `json:"output"`
	Explanation string `json:"explanation,omitempty"`
}

type Problem struct {
	ID           int               `db:"id"`
	Title        string            `db:"title"`
	Slug         string            `db:"slug"`
	Description  string            `db:"description"`
	Difficulty   string            `db:"difficulty"`
	Tags         pq.StringArray    `db:"tags"`
	Examples     []Example         `db:"examples"`
	Constraints  string            `db:"constraints"`
	TemplateCode map[string]string `db:"template_code"`
	CreatedAt    time.Time         `db:"created_at"`
	UpdatedAt    time.Time         `db:"updated_at"`
}

type TestCase struct {
	ID             int       `db:"id"`
	ProblemID      int       `db:"problem_id"`
	Input          string    `db:"input"`
	ExpectedOutput string    `db:"expected_output"`
	IsHidden       bool      `db:"is_hidden"`
	CreatedAt      time.Time `db:"created_at"`
}

type Submission struct {
	ID              int       `db:"id"`
	UserID          int       `db:"user_id"`
	ProblemID       int       `db:"problem_id"`
	Language        string    `db:"language"`
	Code            string    `db:"code"`
	Status          string    `db:"status"`
	RuntimeMs       *int      `db:"runtime_ms"`
	MemoryKb        *int      `db:"memory_kb"`
	TestCasesPassed int       `db:"test_cases_passed"`
	TotalTestCases  int       `db:"total_test_cases"`
	ErrorMessage    *string   `db:"error_message"`
	SubmittedAt     time.Time `db:"submitted_at"`
}

type UserProgress struct {
	UserID           int        `db:"user_id"`
	ProblemID        int        `db:"problem_id"`
	IsSolved         bool       `db:"is_solved"`
	BestSubmissionID *int       `db:"best_submission_id"`
	Attempts         int        `db:"attempts"`
	FirstSolvedAt    *time.Time `db:"first_solved_at"`
}

func main() {
	config := Config{
		Host:     getEnv("DB_HOST", "localhost"),
		Port:     getEnv("DB_PORT", "5432"),
		User:     getEnv("DB_USER", "postgres"),
		Password: getEnv("DB_PASSWORD", "password"),
		DBName:   getEnv("DB_NAME", "leetcode_clone"),
	}

	db, err := connectDB(config)
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}
	defer db.Close()

	log.Println("Starting database seeding...")

	// Clear existing data (optional - comment out if you want to preserve existing data)
	if err := clearExistingData(db); err != nil {
		log.Fatal("Failed to clear existing data:", err)
	}

	// Seed users
	users, err := seedUsers(db)
	if err != nil {
		log.Fatal("Failed to seed users:", err)
	}
	log.Printf("Seeded %d users", len(users))

	// Seed problems
	problems, err := seedProblems(db)
	if err != nil {
		log.Fatal("Failed to seed problems:", err)
	}
	log.Printf("Seeded %d problems", len(problems))

	// Seed test cases
	testCases, err := seedTestCases(db, problems)
	if err != nil {
		log.Fatal("Failed to seed test cases:", err)
	}
	log.Printf("Seeded %d test cases", len(testCases))

	// Seed submissions
	submissions, err := seedSubmissions(db, users, problems)
	if err != nil {
		log.Fatal("Failed to seed submissions:", err)
	}
	log.Printf("Seeded %d submissions", len(submissions))

	// Seed user progress
	progress, err := seedUserProgress(db, users, problems, submissions)
	if err != nil {
		log.Fatal("Failed to seed user progress:", err)
	}
	log.Printf("Seeded %d user progress records", len(progress))

	log.Println("Database seeding completed successfully!")
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func connectDB(config Config) (*sql.DB, error) {
	connStr := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		config.Host, config.Port, config.User, config.Password, config.DBName)
	
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, err
	}

	if err := db.Ping(); err != nil {
		return nil, err
	}

	return db, nil
}

func clearExistingData(db *sql.DB) error {
	tables := []string{"user_progress", "submissions", "test_cases", "problems", "users"}
	
	for _, table := range tables {
		query := fmt.Sprintf("TRUNCATE TABLE %s RESTART IDENTITY CASCADE", table)
		if _, err := db.Exec(query); err != nil {
			return fmt.Errorf("failed to truncate %s: %w", table, err)
		}
	}
	
	return nil
}

func hashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(bytes), err
}

func seedUsers(db *sql.DB) ([]User, error) {
	users := []User{
		{
			Username: "admin",
			Email:    "admin@leetcodeclone.com",
			IsAdmin:  true,
		},
		{
			Username: "john_doe",
			Email:    "john.doe@example.com",
			IsAdmin:  false,
		},
		{
			Username: "jane_smith",
			Email:    "jane.smith@example.com",
			IsAdmin:  false,
		},
		{
			Username: "alice_johnson",
			Email:    "alice.johnson@example.com",
			IsAdmin:  false,
		},
		{
			Username: "bob_wilson",
			Email:    "bob.wilson@example.com",
			IsAdmin:  false,
		},
	}

	var createdUsers []User

	for _, user := range users {
		// Hash password (using username as password for demo purposes)
		hashedPassword, err := hashPassword(user.Username + "123")
		if err != nil {
			return nil, err
		}

		query := `
			INSERT INTO users (username, email, password_hash, is_admin, created_at, updated_at)
			VALUES ($1, $2, $3, $4, $5, $6)
			RETURNING id, username, email, password_hash, is_admin, created_at, updated_at`

		now := time.Now()
		var createdUser User
		err = db.QueryRow(query, user.Username, user.Email, hashedPassword, user.IsAdmin, now, now).Scan(
			&createdUser.ID, &createdUser.Username, &createdUser.Email, &createdUser.PasswordHash,
			&createdUser.IsAdmin, &createdUser.CreatedAt, &createdUser.UpdatedAt)
		if err != nil {
			return nil, err
		}

		createdUsers = append(createdUsers, createdUser)
	}

	return createdUsers, nil
}

func seedProblems(db *sql.DB) ([]Problem, error) {
	problems := []Problem{
		{
			Title:       "Two Sum",
			Slug:        "two-sum",
			Description: "Given an array of integers nums and an integer target, return indices of the two numbers such that they add up to target.\n\nYou may assume that each input would have exactly one solution, and you may not use the same element twice.\n\nYou can return the answer in any order.",
			Difficulty:  "Easy",
			Tags:        pq.StringArray{"Array", "Hash Table"},
			Examples: []Example{
				{
					Input:       "nums = [2,7,11,15], target = 9",
					Output:      "[0,1]",
					Explanation: "Because nums[0] + nums[1] == 9, we return [0, 1].",
				},
				{
					Input:       "nums = [3,2,4], target = 6",
					Output:      "[1,2]",
					Explanation: "Because nums[1] + nums[2] == 6, we return [1, 2].",
				},
			},
			Constraints: "2 <= nums.length <= 10^4\n-10^9 <= nums[i] <= 10^9\n-10^9 <= target <= 10^9\nOnly one valid answer exists.",
			TemplateCode: map[string]string{
				"javascript": "/**\n * @param {number[]} nums\n * @param {number} target\n * @return {number[]}\n */\nvar twoSum = function(nums, target) {\n    \n};",
				"python":     "class Solution:\n    def twoSum(self, nums: List[int], target: int) -> List[int]:\n        ",
				"java":       "class Solution {\n    public int[] twoSum(int[] nums, int target) {\n        \n    }\n}",
			},
		},
		{
			Title:       "Add Two Numbers",
			Slug:        "add-two-numbers",
			Description: "You are given two non-empty linked lists representing two non-negative integers. The digits are stored in reverse order, and each of their nodes contains a single digit. Add the two numbers and return the sum as a linked list.\n\nYou may assume the two numbers do not contain any leading zero, except the number 0 itself.",
			Difficulty:  "Medium",
			Tags:        pq.StringArray{"Linked List", "Math", "Recursion"},
			Examples: []Example{
				{
					Input:       "l1 = [2,4,3], l2 = [5,6,4]",
					Output:      "[7,0,8]",
					Explanation: "342 + 465 = 807.",
				},
				{
					Input:       "l1 = [0], l2 = [0]",
					Output:      "[0]",
					Explanation: "0 + 0 = 0.",
				},
			},
			Constraints: "The number of nodes in each linked list is in the range [1, 100].\n0 <= Node.val <= 9\nIt is guaranteed that the list represents a number that does not have leading zeros.",
			TemplateCode: map[string]string{
				"javascript": "/**\n * Definition for singly-linked list.\n * function ListNode(val, next) {\n *     this.val = (val===undefined ? 0 : val)\n *     this.next = (next===undefined ? null : next)\n * }\n */\n/**\n * @param {ListNode} l1\n * @param {ListNode} l2\n * @return {ListNode}\n */\nvar addTwoNumbers = function(l1, l2) {\n    \n};",
				"python":     "# Definition for singly-linked list.\n# class ListNode:\n#     def __init__(self, val=0, next=None):\n#         self.val = val\n#         self.next = next\nclass Solution:\n    def addTwoNumbers(self, l1: Optional[ListNode], l2: Optional[ListNode]) -> Optional[ListNode]:\n        ",
				"java":       "/**\n * Definition for singly-linked list.\n * public class ListNode {\n *     int val;\n *     ListNode next;\n *     ListNode() {}\n *     ListNode(int val) { this.val = val; }\n *     ListNode(int val, ListNode next) { this.val = val; this.next = next; }\n * }\n */\nclass Solution {\n    public ListNode addTwoNumbers(ListNode l1, ListNode l2) {\n        \n    }\n}",
			},
		},
		{
			Title:       "Longest Substring Without Repeating Characters",
			Slug:        "longest-substring-without-repeating-characters",
			Description: "Given a string s, find the length of the longest substring without repeating characters.",
			Difficulty:  "Medium",
			Tags:        pq.StringArray{"Hash Table", "String", "Sliding Window"},
			Examples: []Example{
				{
					Input:       "s = \"abcabcbb\"",
					Output:      "3",
					Explanation: "The answer is \"abc\", with the length of 3.",
				},
				{
					Input:       "s = \"bbbbb\"",
					Output:      "1",
					Explanation: "The answer is \"b\", with the length of 1.",
				},
				{
					Input:       "s = \"pwwkew\"",
					Output:      "3",
					Explanation: "The answer is \"wke\", with the length of 3.",
				},
			},
			Constraints: "0 <= s.length <= 5 * 10^4\ns consists of English letters, digits, symbols and spaces.",
			TemplateCode: map[string]string{
				"javascript": "/**\n * @param {string} s\n * @return {number}\n */\nvar lengthOfLongestSubstring = function(s) {\n    \n};",
				"python":     "class Solution:\n    def lengthOfLongestSubstring(self, s: str) -> int:\n        ",
				"java":       "class Solution {\n    public int lengthOfLongestSubstring(String s) {\n        \n    }\n}",
			},
		},
		{
			Title:       "Median of Two Sorted Arrays",
			Slug:        "median-of-two-sorted-arrays",
			Description: "Given two sorted arrays nums1 and nums2 of size m and n respectively, return the median of the two sorted arrays.\n\nThe overall run time complexity should be O(log (m+n)).",
			Difficulty:  "Hard",
			Tags:        pq.StringArray{"Array", "Binary Search", "Divide and Conquer"},
			Examples: []Example{
				{
					Input:       "nums1 = [1,3], nums2 = [2]",
					Output:      "2.00000",
					Explanation: "merged array = [1,2,3] and median is 2.",
				},
				{
					Input:       "nums1 = [1,2], nums2 = [3,4]",
					Output:      "2.50000",
					Explanation: "merged array = [1,2,3,4] and median is (2 + 3) / 2 = 2.5.",
				},
			},
			Constraints: "nums1.length == m\nnums2.length == n\n0 <= m <= 1000\n0 <= n <= 1000\n1 <= m + n <= 2000\n-10^6 <= nums1[i], nums2[i] <= 10^6",
			TemplateCode: map[string]string{
				"javascript": "/**\n * @param {number[]} nums1\n * @param {number[]} nums2\n * @return {number}\n */\nvar findMedianSortedArrays = function(nums1, nums2) {\n    \n};",
				"python":     "class Solution:\n    def findMedianSortedArrays(self, nums1: List[int], nums2: List[int]) -> float:\n        ",
				"java":       "class Solution {\n    public double findMedianSortedArrays(int[] nums1, int[] nums2) {\n        \n    }\n}",
			},
		},
		{
			Title:       "Valid Parentheses",
			Slug:        "valid-parentheses",
			Description: "Given a string s containing just the characters '(', ')', '{', '}', '[' and ']', determine if the input string is valid.\n\nAn input string is valid if:\n1. Open brackets must be closed by the same type of brackets.\n2. Open brackets must be closed in the correct order.\n3. Every close bracket has a corresponding open bracket of the same type.",
			Difficulty:  "Easy",
			Tags:        pq.StringArray{"String", "Stack"},
			Examples: []Example{
				{
					Input:  "s = \"()\"",
					Output: "true",
				},
				{
					Input:  "s = \"()[]{}\"",
					Output: "true",
				},
				{
					Input:  "s = \"(]\"",
					Output: "false",
				},
			},
			Constraints: "1 <= s.length <= 10^4\ns consists of parentheses only '()[]{}'.",
			TemplateCode: map[string]string{
				"javascript": "/**\n * @param {string} s\n * @return {boolean}\n */\nvar isValid = function(s) {\n    \n};",
				"python":     "class Solution:\n    def isValid(self, s: str) -> bool:\n        ",
				"java":       "class Solution {\n    public boolean isValid(String s) {\n        \n    }\n}",
			},
		},
		{
			Title:       "Climbing Stairs",
			Slug:        "climbing-stairs",
			Description: "You are climbing a staircase. It takes n steps to reach the top.\n\nEach time you can either climb 1 or 2 steps. In how many distinct ways can you climb to the top?",
			Difficulty:  "Easy",
			Tags:        pq.StringArray{"Math", "Dynamic Programming", "Memoization"},
			Examples: []Example{
				{
					Input:       "n = 2",
					Output:      "2",
					Explanation: "There are two ways to climb to the top.\n1. 1 step + 1 step\n2. 2 steps",
				},
				{
					Input:       "n = 3",
					Output:      "3",
					Explanation: "There are three ways to climb to the top.\n1. 1 step + 1 step + 1 step\n2. 1 step + 2 steps\n3. 2 steps + 1 step",
				},
			},
			Constraints: "1 <= n <= 45",
			TemplateCode: map[string]string{
				"javascript": "/**\n * @param {number} n\n * @return {number}\n */\nvar climbStairs = function(n) {\n    \n};",
				"python":     "class Solution:\n    def climbStairs(self, n: int) -> int:\n        ",
				"java":       "class Solution {\n    public int climbStairs(int n) {\n        \n    }\n}",
			},
		},
		{
			Title:       "Maximum Subarray",
			Slug:        "maximum-subarray",
			Description: "Given an integer array nums, find the contiguous subarray (containing at least one number) which has the largest sum and return its sum.\n\nA subarray is a contiguous part of an array.",
			Difficulty:  "Medium",
			Tags:        pq.StringArray{"Array", "Divide and Conquer", "Dynamic Programming"},
			Examples: []Example{
				{
					Input:       "nums = [-2,1,-3,4,-1,2,1,-5,4]",
					Output:      "6",
					Explanation: "[4,-1,2,1] has the largest sum = 6.",
				},
				{
					Input:       "nums = [1]",
					Output:      "1",
					Explanation: "The subarray [1] has the largest sum = 1.",
				},
				{
					Input:       "nums = [5,4,-1,7,8]",
					Output:      "23",
					Explanation: "[5,4,-1,7,8] has the largest sum = 23.",
				},
			},
			Constraints: "1 <= nums.length <= 10^5\n-10^4 <= nums[i] <= 10^4",
			TemplateCode: map[string]string{
				"javascript": "/**\n * @param {number[]} nums\n * @return {number}\n */\nvar maxSubArray = function(nums) {\n    \n};",
				"python":     "class Solution:\n    def maxSubArray(self, nums: List[int]) -> int:\n        ",
				"java":       "class Solution {\n    public int maxSubArray(int[] nums) {\n        \n    }\n}",
			},
		},
		{
			Title:       "Coin Change",
			Slug:        "coin-change",
			Description: "You are given an integer array coins representing coins of different denominations and an integer amount representing a total amount of money.\n\nReturn the fewest number of coins that you need to make up that amount. If that amount of money cannot be made up by any combination of the coins, return -1.\n\nYou may assume that you have an infinite number of each kind of coin.",
			Difficulty:  "Medium",
			Tags:        pq.StringArray{"Array", "Dynamic Programming", "Breadth-First Search"},
			Examples: []Example{
				{
					Input:       "coins = [1,3,4], amount = 6",
					Output:      "2",
					Explanation: "The minimum number of coins is 2 (3 + 3 = 6).",
				},
				{
					Input:       "coins = [2], amount = 3",
					Output:      "-1",
					Explanation: "The amount cannot be made up with the given coins.",
				},
			},
			Constraints: "1 <= coins.length <= 12\n1 <= coins[i] <= 2^31 - 1\n0 <= amount <= 10^4",
			TemplateCode: map[string]string{
				"javascript": "/**\n * @param {number[]} coins\n * @param {number} amount\n * @return {number}\n */\nvar coinChange = function(coins, amount) {\n    \n};",
				"python":     "class Solution:\n    def coinChange(self, coins: List[int], amount: int) -> int:\n        ",
				"java":       "class Solution {\n    public int coinChange(int[] coins, int amount) {\n        \n    }\n}",
			},
		},
	}

	var createdProblems []Problem

	for _, problem := range problems {
		// Convert examples to JSON
		examplesJSON, err := json.Marshal(problem.Examples)
		if err != nil {
			return nil, err
		}

		// Convert template code to JSON
		templateCodeJSON, err := json.Marshal(problem.TemplateCode)
		if err != nil {
			return nil, err
		}

		query := `
			INSERT INTO problems (title, slug, description, difficulty, tags, examples, constraints, template_code, created_at, updated_at)
			VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
			RETURNING id, title, slug, description, difficulty, tags, examples, constraints, template_code, created_at, updated_at`

		now := time.Now()
		var createdProblem Problem
		var examplesStr, templateCodeStr string

		err = db.QueryRow(query, problem.Title, problem.Slug, problem.Description, problem.Difficulty,
			problem.Tags, examplesJSON, problem.Constraints, templateCodeJSON, now, now).Scan(
			&createdProblem.ID, &createdProblem.Title, &createdProblem.Slug, &createdProblem.Description,
			&createdProblem.Difficulty, &createdProblem.Tags, &examplesStr, &createdProblem.Constraints,
			&templateCodeStr, &createdProblem.CreatedAt, &createdProblem.UpdatedAt)
		if err != nil {
			return nil, err
		}

		// Parse back the JSON fields
		if err := json.Unmarshal([]byte(examplesStr), &createdProblem.Examples); err != nil {
			return nil, err
		}
		if err := json.Unmarshal([]byte(templateCodeStr), &createdProblem.TemplateCode); err != nil {
			return nil, err
		}

		createdProblems = append(createdProblems, createdProblem)
	}

	return createdProblems, nil
}

func seedTestCases(db *sql.DB, problems []Problem) ([]TestCase, error) {
	// Map problem slugs to IDs for easier lookup
	problemMap := make(map[string]int)
	for _, p := range problems {
		problemMap[p.Slug] = p.ID
	}

	testCasesData := map[string][]TestCase{
		"two-sum": {
			{Input: "[2,7,11,15]\n9", ExpectedOutput: "[0,1]", IsHidden: false},
			{Input: "[3,2,4]\n6", ExpectedOutput: "[1,2]", IsHidden: false},
			{Input: "[3,3]\n6", ExpectedOutput: "[0,1]", IsHidden: false},
			{Input: "[1,2,3,4,5]\n9", ExpectedOutput: "[3,4]", IsHidden: true},
			{Input: "[-1,-2,-3,-4,-5]\n-8", ExpectedOutput: "[2,4]", IsHidden: true},
		},
		"add-two-numbers": {
			{Input: "[2,4,3]\n[5,6,4]", ExpectedOutput: "[7,0,8]", IsHidden: false},
			{Input: "[0]\n[0]", ExpectedOutput: "[0]", IsHidden: false},
			{Input: "[9,9,9,9,9,9,9]\n[9,9,9,9]", ExpectedOutput: "[8,9,9,9,0,0,0,1]", IsHidden: false},
			{Input: "[1,8]\n[0]", ExpectedOutput: "[1,8]", IsHidden: true},
			{Input: "[5]\n[5]", ExpectedOutput: "[0,1]", IsHidden: true},
		},
		"longest-substring-without-repeating-characters": {
			{Input: "abcabcbb", ExpectedOutput: "3", IsHidden: false},
			{Input: "bbbbb", ExpectedOutput: "1", IsHidden: false},
			{Input: "pwwkew", ExpectedOutput: "3", IsHidden: false},
			{Input: "", ExpectedOutput: "0", IsHidden: true},
			{Input: "dvdf", ExpectedOutput: "3", IsHidden: true},
		},
		"median-of-two-sorted-arrays": {
			{Input: "[1,3]\n[2]", ExpectedOutput: "2.00000", IsHidden: false},
			{Input: "[1,2]\n[3,4]", ExpectedOutput: "2.50000", IsHidden: false},
			{Input: "[0,0]\n[0,0]", ExpectedOutput: "0.00000", IsHidden: true},
			{Input: "[]\n[1]", ExpectedOutput: "1.00000", IsHidden: true},
			{Input: "[2]\n[]", ExpectedOutput: "2.00000", IsHidden: true},
		},
		"valid-parentheses": {
			{Input: "()", ExpectedOutput: "true", IsHidden: false},
			{Input: "()[]{}", ExpectedOutput: "true", IsHidden: false},
			{Input: "(]", ExpectedOutput: "false", IsHidden: false},
			{Input: "([)]", ExpectedOutput: "false", IsHidden: true},
			{Input: "{[]}", ExpectedOutput: "true", IsHidden: true},
		},
		"climbing-stairs": {
			{Input: "2", ExpectedOutput: "2", IsHidden: false},
			{Input: "3", ExpectedOutput: "3", IsHidden: false},
			{Input: "1", ExpectedOutput: "1", IsHidden: true},
			{Input: "4", ExpectedOutput: "5", IsHidden: true},
			{Input: "5", ExpectedOutput: "8", IsHidden: true},
		},
		"maximum-subarray": {
			{Input: "[-2,1,-3,4,-1,2,1,-5,4]", ExpectedOutput: "6", IsHidden: false},
			{Input: "[1]", ExpectedOutput: "1", IsHidden: false},
			{Input: "[5,4,-1,7,8]", ExpectedOutput: "23", IsHidden: false},
			{Input: "[-1]", ExpectedOutput: "-1", IsHidden: true},
			{Input: "[-2,-1]", ExpectedOutput: "-1", IsHidden: true},
		},
		"coin-change": {
			{Input: "[1,3,4]\n6", ExpectedOutput: "2", IsHidden: false},
			{Input: "[2]\n3", ExpectedOutput: "-1", IsHidden: false},
			{Input: "[1]\n0", ExpectedOutput: "0", IsHidden: false},
			{Input: "[1,2,5]\n11", ExpectedOutput: "3", IsHidden: true},
			{Input: "[2]\n1", ExpectedOutput: "-1", IsHidden: true},
		},
	}

	var createdTestCases []TestCase

	for slug, testCases := range testCasesData {
		problemID, exists := problemMap[slug]
		if !exists {
			continue
		}

		for _, testCase := range testCases {
			testCase.ProblemID = problemID

			query := `
				INSERT INTO test_cases (problem_id, input, expected_output, is_hidden, created_at)
				VALUES ($1, $2, $3, $4, $5)
				RETURNING id, problem_id, input, expected_output, is_hidden, created_at`

			var createdTestCase TestCase
			err := db.QueryRow(query, testCase.ProblemID, testCase.Input, testCase.ExpectedOutput,
				testCase.IsHidden, time.Now()).Scan(
				&createdTestCase.ID, &createdTestCase.ProblemID, &createdTestCase.Input,
				&createdTestCase.ExpectedOutput, &createdTestCase.IsHidden, &createdTestCase.CreatedAt)
			if err != nil {
				return nil, err
			}

			createdTestCases = append(createdTestCases, createdTestCase)
		}
	}

	return createdTestCases, nil
}

func seedSubmissions(db *sql.DB, users []User, problems []Problem) ([]Submission, error) {
	// Create sample submissions for different users and problems
	submissionsData := []struct {
		UserIndex    int
		ProblemIndex int
		Language     string
		Code         string
		Status       string
		RuntimeMs    *int
		MemoryKb     *int
		TestsPassed  int
		TotalTests   int
		ErrorMsg     *string
	}{
		// John Doe submissions
		{
			UserIndex: 1, ProblemIndex: 0, Language: "javascript", Status: "Accepted",
			RuntimeMs: intPtr(68), MemoryKb: intPtr(44200), TestsPassed: 5, TotalTests: 5,
			Code: `var twoSum = function(nums, target) {
    const map = new Map();
    for (let i = 0; i < nums.length; i++) {
        const complement = target - nums[i];
        if (map.has(complement)) {
            return [map.get(complement), i];
        }
        map.set(nums[i], i);
    }
    return [];
};`,
		},
		{
			UserIndex: 1, ProblemIndex: 4, Language: "python", Status: "Accepted",
			RuntimeMs: intPtr(32), MemoryKb: intPtr(16400), TestsPassed: 5, TotalTests: 5,
			Code: `class Solution:
    def isValid(self, s: str) -> bool:
        stack = []
        mapping = {")": "(", "}": "{", "]": "["}
        
        for char in s:
            if char in mapping:
                if not stack or stack.pop() != mapping[char]:
                    return False
            else:
                stack.append(char)
        
        return not stack`,
		},
		{
			UserIndex: 1, ProblemIndex: 5, Language: "javascript", Status: "Wrong Answer",
			RuntimeMs: nil, MemoryKb: nil, TestsPassed: 2, TotalTests: 5,
			ErrorMsg: strPtr("Wrong Answer on test case 3"),
			Code: `var climbStairs = function(n) {
    if (n <= 2) return n;
    return climbStairs(n-1) + climbStairs(n-2);
};`,
		},
		// Jane Smith submissions
		{
			UserIndex: 2, ProblemIndex: 0, Language: "python", Status: "Accepted",
			RuntimeMs: intPtr(52), MemoryKb: intPtr(15100), TestsPassed: 5, TotalTests: 5,
			Code: `class Solution:
    def twoSum(self, nums: List[int], target: int) -> List[int]:
        num_map = {}
        for i, num in enumerate(nums):
            complement = target - num
            if complement in num_map:
                return [num_map[complement], i]
            num_map[num] = i
        return []`,
		},
		{
			UserIndex: 2, ProblemIndex: 2, Language: "python", Status: "Accepted",
			RuntimeMs: intPtr(76), MemoryKb: intPtr(14800), TestsPassed: 5, TotalTests: 5,
			Code: `class Solution:
    def lengthOfLongestSubstring(self, s: str) -> int:
        char_map = {}
        left = 0
        max_length = 0
        
        for right in range(len(s)):
            if s[right] in char_map:
                left = max(left, char_map[s[right]] + 1)
            char_map[s[right]] = right
            max_length = max(max_length, right - left + 1)
        
        return max_length`,
		},
		{
			UserIndex: 2, ProblemIndex: 6, Language: "java", Status: "Time Limit Exceeded",
			RuntimeMs: nil, MemoryKb: nil, TestsPassed: 3, TotalTests: 5,
			ErrorMsg: strPtr("Time Limit Exceeded"),
			Code: `class Solution {
    public int maxSubArray(int[] nums) {
        int maxSum = Integer.MIN_VALUE;
        for (int i = 0; i < nums.length; i++) {
            for (int j = i; j < nums.length; j++) {
                int sum = 0;
                for (int k = i; k <= j; k++) {
                    sum += nums[k];
                }
                maxSum = Math.max(maxSum, sum);
            }
        }
        return maxSum;
    }
}`,
		},
		// Alice Johnson submissions
		{
			UserIndex: 3, ProblemIndex: 5, Language: "python", Status: "Accepted",
			RuntimeMs: intPtr(28), MemoryKb: intPtr(16200), TestsPassed: 5, TotalTests: 5,
			Code: `class Solution:
    def climbStairs(self, n: int) -> int:
        if n <= 2:
            return n
        
        prev2, prev1 = 1, 2
        for i in range(3, n + 1):
            current = prev1 + prev2
            prev2, prev1 = prev1, current
        
        return prev1`,
		},
		{
			UserIndex: 3, ProblemIndex: 7, Language: "javascript", Status: "Accepted",
			RuntimeMs: intPtr(84), MemoryKb: intPtr(42300), TestsPassed: 5, TotalTests: 5,
			Code: `var coinChange = function(coins, amount) {
    const dp = new Array(amount + 1).fill(amount + 1);
    dp[0] = 0;
    
    for (let i = 1; i <= amount; i++) {
        for (let coin of coins) {
            if (coin <= i) {
                dp[i] = Math.min(dp[i], dp[i - coin] + 1);
            }
        }
    }
    
    return dp[amount] > amount ? -1 : dp[amount];
};`,
		},
		// Bob Wilson submissions
		{
			UserIndex: 4, ProblemIndex: 4, Language: "java", Status: "Compile Error",
			RuntimeMs: nil, MemoryKb: nil, TestsPassed: 0, TotalTests: 5,
			ErrorMsg: strPtr("Compile Error: missing return statement"),
			Code: `class Solution {
    public boolean isValid(String s) {
        Stack<Character> stack = new Stack<>();
        for (char c : s.toCharArray()) {
            if (c == '(' || c == '[' || c == '{') {
                stack.push(c);
            } else {
                if (stack.isEmpty()) return false;
                char top = stack.pop();
                if ((c == ')' && top != '(') || 
                    (c == ']' && top != '[') || 
                    (c == '}' && top != '{')) {
                    return false;
                }
            }
        }
        // Missing return statement
    }
}`,
		},
	}

	var createdSubmissions []Submission

	for _, subData := range submissionsData {
		if subData.UserIndex >= len(users) || subData.ProblemIndex >= len(problems) {
			continue
		}

		submission := Submission{
			UserID:          users[subData.UserIndex].ID,
			ProblemID:       problems[subData.ProblemIndex].ID,
			Language:        subData.Language,
			Code:            subData.Code,
			Status:          subData.Status,
			RuntimeMs:       subData.RuntimeMs,
			MemoryKb:        subData.MemoryKb,
			TestCasesPassed: subData.TestsPassed,
			TotalTestCases:  subData.TotalTests,
			ErrorMessage:    subData.ErrorMsg,
		}

		query := `
			INSERT INTO submissions (user_id, problem_id, language, code, status, runtime_ms, memory_kb, 
				test_cases_passed, total_test_cases, error_message, submitted_at)
			VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)
			RETURNING id, user_id, problem_id, language, code, status, runtime_ms, memory_kb, 
				test_cases_passed, total_test_cases, error_message, submitted_at`

		var createdSubmission Submission
		err := db.QueryRow(query, submission.UserID, submission.ProblemID, submission.Language,
			submission.Code, submission.Status, submission.RuntimeMs, submission.MemoryKb,
			submission.TestCasesPassed, submission.TotalTestCases, submission.ErrorMessage,
			time.Now()).Scan(
			&createdSubmission.ID, &createdSubmission.UserID, &createdSubmission.ProblemID,
			&createdSubmission.Language, &createdSubmission.Code, &createdSubmission.Status,
			&createdSubmission.RuntimeMs, &createdSubmission.MemoryKb, &createdSubmission.TestCasesPassed,
			&createdSubmission.TotalTestCases, &createdSubmission.ErrorMessage, &createdSubmission.SubmittedAt)
		if err != nil {
			return nil, err
		}

		createdSubmissions = append(createdSubmissions, createdSubmission)
	}

	return createdSubmissions, nil
}

func seedUserProgress(db *sql.DB, users []User, problems []Problem, submissions []Submission) ([]UserProgress, error) {
	// Create a map of successful submissions by user and problem
	successfulSubmissions := make(map[string]int) // key: "userID-problemID", value: submissionID
	
	for _, sub := range submissions {
		if sub.Status == "Accepted" {
			key := fmt.Sprintf("%d-%d", sub.UserID, sub.ProblemID)
			successfulSubmissions[key] = sub.ID
		}
	}

	progressData := []struct {
		UserIndex    int
		ProblemIndex int
		IsSolved     bool
		Attempts     int
	}{
		// John Doe progress
		{UserIndex: 1, ProblemIndex: 0, IsSolved: true, Attempts: 1},   // Two Sum - solved
		{UserIndex: 1, ProblemIndex: 4, IsSolved: true, Attempts: 1},   // Valid Parentheses - solved
		{UserIndex: 1, ProblemIndex: 5, IsSolved: false, Attempts: 3},  // Climbing Stairs - attempted but not solved
		{UserIndex: 1, ProblemIndex: 1, IsSolved: false, Attempts: 2},  // Add Two Numbers - attempted but not solved
		
		// Jane Smith progress
		{UserIndex: 2, ProblemIndex: 0, IsSolved: true, Attempts: 1},   // Two Sum - solved
		{UserIndex: 2, ProblemIndex: 2, IsSolved: true, Attempts: 2},   // Longest Substring - solved
		{UserIndex: 2, ProblemIndex: 6, IsSolved: false, Attempts: 4},  // Maximum Subarray - attempted but not solved
		{UserIndex: 2, ProblemIndex: 3, IsSolved: false, Attempts: 1},  // Median of Two Sorted Arrays - attempted
		
		// Alice Johnson progress
		{UserIndex: 3, ProblemIndex: 5, IsSolved: true, Attempts: 1},   // Climbing Stairs - solved
		{UserIndex: 3, ProblemIndex: 7, IsSolved: true, Attempts: 2},   // Coin Change - solved
		{UserIndex: 3, ProblemIndex: 0, IsSolved: false, Attempts: 1},  // Two Sum - attempted
		
		// Bob Wilson progress
		{UserIndex: 4, ProblemIndex: 4, IsSolved: false, Attempts: 5},  // Valid Parentheses - multiple failed attempts
		{UserIndex: 4, ProblemIndex: 0, IsSolved: false, Attempts: 2},  // Two Sum - attempted
	}

	var createdProgress []UserProgress

	for _, progData := range progressData {
		if progData.UserIndex >= len(users) || progData.ProblemIndex >= len(problems) {
			continue
		}

		userID := users[progData.UserIndex].ID
		problemID := problems[progData.ProblemIndex].ID
		
		progress := UserProgress{
			UserID:    userID,
			ProblemID: problemID,
			IsSolved:  progData.IsSolved,
			Attempts:  progData.Attempts,
		}

		// Set best submission ID and first solved time if problem is solved
		if progData.IsSolved {
			key := fmt.Sprintf("%d-%d", userID, problemID)
			if submissionID, exists := successfulSubmissions[key]; exists {
				progress.BestSubmissionID = &submissionID
				solvedTime := time.Now().Add(-time.Duration(progData.ProblemIndex*24) * time.Hour) // Vary solved times
				progress.FirstSolvedAt = &solvedTime
			}
		}

		query := `
			INSERT INTO user_progress (user_id, problem_id, is_solved, best_submission_id, attempts, first_solved_at)
			VALUES ($1, $2, $3, $4, $5, $6)
			RETURNING user_id, problem_id, is_solved, best_submission_id, attempts, first_solved_at`

		var createdUserProgress UserProgress
		err := db.QueryRow(query, progress.UserID, progress.ProblemID, progress.IsSolved,
			progress.BestSubmissionID, progress.Attempts, progress.FirstSolvedAt).Scan(
			&createdUserProgress.UserID, &createdUserProgress.ProblemID, &createdUserProgress.IsSolved,
			&createdUserProgress.BestSubmissionID, &createdUserProgress.Attempts, &createdUserProgress.FirstSolvedAt)
		if err != nil {
			return nil, err
		}

		createdProgress = append(createdProgress, createdUserProgress)
	}

	return createdProgress, nil
}

// Helper functions
func intPtr(i int) *int {
	return &i
}

func strPtr(s string) *string {
	return &s
}