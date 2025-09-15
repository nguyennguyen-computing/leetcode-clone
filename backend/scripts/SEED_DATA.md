# Seed Data Documentation

This document describes the sample data that gets created when running the database seeding script.

## Overview

The seed script populates the database with realistic test data to help with development and testing. It creates users, problems, test cases, submissions, and user progress data.

## Running the Seed Script

### Option 1: Using the shell script
```bash
cd backend/scripts
./seed.sh
```

### Option 2: Using Make (from project root)
```bash
make seed-db
```

### Option 3: Direct Go execution
```bash
cd backend/scripts
go run seed_data.go
```

## Environment Variables

The seed script uses the following environment variables (with defaults):

- `DB_HOST` (default: localhost)
- `DB_PORT` (default: 5432)
- `DB_USER` (default: postgres)
- `DB_PASSWORD` (default: password)
- `DB_NAME` (default: leetcode_clone)

## Sample Data Created

### Users (5 total)

| Username | Email | Role | Password | Description |
|----------|-------|------|----------|-------------|
| admin | admin@leetcodeclone.com | Admin | admin123 | System administrator |
| john_doe | john.doe@example.com | User | john_doe123 | Active user with multiple submissions |
| jane_smith | jane.smith@example.com | User | jane_smith123 | Experienced user with good solutions |
| alice_johnson | alice.johnson@example.com | User | alice_johnson123 | User focused on dynamic programming |
| bob_wilson | bob.wilson@example.com | User | bob_wilson123 | Beginner user with some failed attempts |

### Problems (8 total)

| Title | Slug | Difficulty | Tags | Description |
|-------|------|------------|------|-------------|
| Two Sum | two-sum | Easy | Array, Hash Table | Find two numbers that add up to target |
| Add Two Numbers | add-two-numbers | Medium | Linked List, Math, Recursion | Add two numbers represented as linked lists |
| Longest Substring Without Repeating Characters | longest-substring-without-repeating-characters | Medium | Hash Table, String, Sliding Window | Find longest substring without repeating chars |
| Median of Two Sorted Arrays | median-of-two-sorted-arrays | Hard | Array, Binary Search, Divide and Conquer | Find median of two sorted arrays |
| Valid Parentheses | valid-parentheses | Easy | String, Stack | Check if parentheses are valid |
| Climbing Stairs | climbing-stairs | Easy | Math, Dynamic Programming, Memoization | Count ways to climb stairs |
| Maximum Subarray | maximum-subarray | Medium | Array, Divide and Conquer, Dynamic Programming | Find maximum sum subarray |
| Coin Change | coin-change | Medium | Array, Dynamic Programming, BFS | Find minimum coins for amount |

### Test Cases (40+ total)

Each problem has 3-5 test cases:
- **Public test cases**: Visible to users (match examples in problem description)
- **Hidden test cases**: Used for final evaluation (edge cases, corner cases)

Examples:
- Two Sum: Basic cases, negative numbers, duplicate values
- Valid Parentheses: Simple cases, nested brackets, invalid combinations
- Climbing Stairs: Small values, edge cases

### Submissions (9 total)

Sample submissions showing different scenarios:

#### Successful Submissions
- **john_doe**: Two Sum (JavaScript) - Accepted in 68ms
- **john_doe**: Valid Parentheses (Python) - Accepted in 32ms
- **jane_smith**: Two Sum (Python) - Accepted in 52ms
- **jane_smith**: Longest Substring (Python) - Accepted in 76ms
- **alice_johnson**: Climbing Stairs (Python) - Accepted in 28ms
- **alice_johnson**: Coin Change (JavaScript) - Accepted in 84ms

#### Failed Submissions
- **john_doe**: Climbing Stairs (JavaScript) - Wrong Answer (inefficient recursive solution)
- **jane_smith**: Maximum Subarray (Java) - Time Limit Exceeded (O(n³) solution)
- **bob_wilson**: Valid Parentheses (Java) - Compile Error (missing return statement)

### User Progress (12 records)

Tracks each user's progress on problems they've attempted:

#### john_doe
- Two Sum: ✅ Solved (1 attempt)
- Valid Parentheses: ✅ Solved (1 attempt)
- Climbing Stairs: ❌ Not solved (3 attempts)
- Add Two Numbers: ❌ Not solved (2 attempts)

#### jane_smith
- Two Sum: ✅ Solved (1 attempt)
- Longest Substring: ✅ Solved (2 attempts)
- Maximum Subarray: ❌ Not solved (4 attempts)
- Median of Two Sorted Arrays: ❌ Not solved (1 attempt)

#### alice_johnson
- Climbing Stairs: ✅ Solved (1 attempt)
- Coin Change: ✅ Solved (2 attempts)
- Two Sum: ❌ Not solved (1 attempt)

#### bob_wilson
- Valid Parentheses: ❌ Not solved (5 attempts)
- Two Sum: ❌ Not solved (2 attempts)

## Data Relationships

The seed data maintains proper referential integrity:

1. **Submissions** reference valid users and problems
2. **Test Cases** belong to specific problems
3. **User Progress** tracks attempts and links to best submissions
4. **Best Submission IDs** in user progress point to actual accepted submissions

## Use Cases

This seed data supports testing of:

### User Authentication & Authorization
- Login with different user types (admin vs regular users)
- Role-based access control

### Problem Browsing & Filtering
- Filter by difficulty (Easy/Medium/Hard)
- Filter by tags (Array, String, Dynamic Programming, etc.)
- Search by title or description

### Code Submission & Evaluation
- Submit code in different languages (JavaScript, Python, Java)
- Handle various submission statuses (Accepted, Wrong Answer, TLE, etc.)
- Test case execution and validation

### Progress Tracking
- View solved/unsolved problems
- Track submission history
- Monitor improvement over time

### Performance Metrics
- Runtime and memory usage tracking
- Comparison between different solutions
- Success rate analysis

## Customization

To modify the seed data:

1. Edit the data structures in `seed_data.go`
2. Add new problems, users, or test cases
3. Adjust the relationships between entities
4. Run the seed script to apply changes

## Cleanup

To clear all seeded data:

```bash
make clean-db
```

Or manually truncate tables:
```sql
TRUNCATE TABLE user_progress, submissions, test_cases, problems, users RESTART IDENTITY CASCADE;
```