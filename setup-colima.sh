#!/bin/bash

# Setup script for Colima Docker container management

echo "Setting up Colima for Docker container management..."

# Check if Colima is installed
if ! command -v colima &> /dev/null; then
    echo "Colima is not installed. Installing via Homebrew..."
    if command -v brew &> /dev/null; then
        brew install colima
    else
        echo "Homebrew not found. Please install Colima manually:"
        echo "https://github.com/abiosoft/colima#installation"
        exit 1
    fi
fi

# Check if Docker CLI is installed
if ! command -v docker &> /dev/null; then
    echo "Docker CLI is not installed. Installing via Homebrew..."
    if command -v brew &> /dev/null; then
        brew install docker
    else
        echo "Homebrew not found. Please install Docker CLI manually"
        exit 1
    fi
fi

# Start Colima with appropriate settings for development
echo "Starting Colima with development settings..."
colima start --cpu 4 --memory 8 --disk 50

# Verify Docker is working
echo "Verifying Docker setup..."
docker --version
docker info

echo "Colima setup complete!"
echo ""
echo "To start the development environment:"
echo "  docker-compose -f docker-compose.dev.yml up"
echo ""
echo "To start the production environment:"
echo "  docker-compose up"
echo ""
echo "To stop Colima:"
echo "  colima stop"