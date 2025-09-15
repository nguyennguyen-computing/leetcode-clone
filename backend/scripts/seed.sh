#!/bin/bash

# Seed database script for LeetCode Clone
# This script runs the Go seed program to populate the database with sample data

set -e

echo "🌱 Starting database seeding..."

# Check if Go is installed
if ! command -v go &> /dev/null; then
    echo "❌ Go is not installed. Please install Go first."
    exit 1
fi

# Set default environment variables if not provided
export DB_HOST=${DB_HOST:-localhost}
export DB_PORT=${DB_PORT:-5432}
export DB_USER=${DB_USER:-postgres}
export DB_PASSWORD=${DB_PASSWORD:-password}
export DB_NAME=${DB_NAME:-leetcode_clone}

echo "📊 Database configuration:"
echo "  Host: $DB_HOST"
echo "  Port: $DB_PORT"
echo "  User: $DB_USER"
echo "  Database: $DB_NAME"
echo ""

# Navigate to the scripts directory
cd "$(dirname "$0")"

# Run the seed program
echo "🚀 Running seed program..."
go run seed_data.go

echo ""
echo "✅ Database seeding completed successfully!"
echo ""
echo "📋 Sample data created:"
echo "  • 5 users (admin, john_doe, jane_smith, alice_johnson, bob_wilson)"
echo "  • 8 coding problems (Easy, Medium, Hard difficulties)"
echo "  • 40+ test cases for problems"
echo "  • 9 sample submissions with various statuses"
echo "  • User progress tracking data"
echo ""
echo "🔐 Default passwords: [username]123 (e.g., admin123, john_doe123)"
echo ""
echo "🎯 You can now test the application with realistic data!"