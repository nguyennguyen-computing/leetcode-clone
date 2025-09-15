# Problem Service API Endpoints

This document describes the API endpoints implemented for the problem management system.

## Public Endpoints (No Authentication Required)

### List Problems
- **GET** `/api/v1/problems`
- **Description**: Retrieve a list of problems with optional filtering
- **Query Parameters**:
  - `difficulty`: Comma-separated list of difficulties (Easy, Medium, Hard)
  - `tags`: Comma-separated list of tags
  - `limit`: Maximum number of results (default: 50, max: 100)
  - `offset`: Number of results to skip (default: 0)
  - `sort_by`: Sort field (title, difficulty, created_at)
  - `sort_order`: Sort order (asc, desc)
- **Response**: 
  ```json
  {
    "problems": [...],
    "count": 10
  }
  ```

### Search Problems
- **GET** `/api/v1/problems/search`
- **Description**: Search problems by title or description
- **Query Parameters**:
  - `q`: Search query (required)
  - Same filtering parameters as List Problems
- **Response**:
  ```json
  {
    "problems": [...],
    "count": 5,
    "query": "two sum"
  }
  ```

### Get Problem by ID
- **GET** `/api/v1/problems/:id`
- **Description**: Retrieve a specific problem by ID
- **Response**: Problem object

### Get Problem by Slug
- **GET** `/api/v1/problems/slug/:slug`
- **Description**: Retrieve a specific problem by slug
- **Response**: Problem object

### Get Test Cases
- **GET** `/api/v1/problems/:id/testcases`
- **Description**: Retrieve test cases for a problem
- **Query Parameters**:
  - `public`: Set to "true" to get only public test cases
- **Response**:
  ```json
  {
    "test_cases": [...],
    "count": 3
  }
  ```

## Admin Endpoints (Authentication + Admin Role Required)

### Create Problem
- **POST** `/api/v1/admin/problems`
- **Description**: Create a new problem
- **Request Body**: Problem object
- **Response**: Created problem object

### Update Problem
- **PUT** `/api/v1/admin/problems/:id`
- **Description**: Update an existing problem
- **Request Body**: Problem object
- **Response**: Updated problem object

### Delete Problem
- **DELETE** `/api/v1/admin/problems/:id`
- **Description**: Delete a problem and all its test cases
- **Response**: 204 No Content

### Create Test Case
- **POST** `/api/v1/admin/problems/:id/testcases`
- **Description**: Create a new test case for a problem
- **Request Body**: Test case object
- **Response**: Created test case object

### Update Test Case
- **PUT** `/api/v1/admin/testcases/:id`
- **Description**: Update an existing test case
- **Request Body**: Test case object
- **Response**: Updated test case object

### Delete Test Case
- **DELETE** `/api/v1/admin/testcases/:id`
- **Description**: Delete a test case
- **Response**: 204 No Content

## Data Models

### Problem Object
```json
{
  "id": 1,
  "title": "Two Sum",
  "slug": "two-sum",
  "description": "Given an array of integers...",
  "difficulty": "Easy",
  "tags": ["Array", "Hash Table"],
  "examples": [
    {
      "input": "[2,7,11,15], 9",
      "output": "[0,1]",
      "explanation": "Because nums[0] + nums[1] == 9, we return [0, 1]."
    }
  ],
  "constraints": "2 <= nums.length <= 10^4",
  "template_code": {
    "javascript": "function twoSum(nums, target) {\n    // Your code here\n}",
    "python": "def two_sum(nums, target):\n    # Your code here\n    pass",
    "java": "public int[] twoSum(int[] nums, int target) {\n    // Your code here\n}"
  },
  "created_at": "2023-01-01T00:00:00Z",
  "updated_at": "2023-01-01T00:00:00Z"
}
```

### Test Case Object
```json
{
  "id": 1,
  "problem_id": 1,
  "input": "[2,7,11,15]\n9",
  "expected_output": "[0,1]",
  "is_hidden": false,
  "created_at": "2023-01-01T00:00:00Z"
}
```

## Error Responses

All endpoints return consistent error responses:

```json
{
  "error": "Error message",
  "details": "Detailed error information"
}
```

Common HTTP status codes:
- `400 Bad Request`: Invalid input or validation errors
- `401 Unauthorized`: Authentication required
- `403 Forbidden`: Admin access required
- `404 Not Found`: Resource not found
- `409 Conflict`: Duplicate resource (e.g., slug already exists)
- `500 Internal Server Error`: Server error

## Authentication

Admin endpoints require:
1. Valid JWT token in Authorization header: `Authorization: Bearer <token>`
2. User must have `is_admin: true` in their user record

## Validation Rules

### Problem Validation
- Title: Required, max 200 characters
- Description: Required
- Difficulty: Must be "Easy", "Medium", or "Hard"
- Examples: At least one example required
- Template Code: At least one language template required
- Supported languages: "javascript", "python", "java"

### Test Case Validation
- Problem ID: Required, must reference existing problem
- Input: Required
- Expected Output: Required

### Filter Validation
- Difficulty: Must be valid difficulty values
- Limit: Max 100, default 50
- Offset: Must be >= 0
- Sort By: Must be "title", "difficulty", or "created_at"
- Sort Order: Must be "asc" or "desc"