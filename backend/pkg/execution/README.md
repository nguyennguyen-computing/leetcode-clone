# Code Execution Service

The execution service provides secure, sandboxed code execution for the LeetCode clone platform. It supports multiple programming languages and implements comprehensive security measures to prevent malicious code execution.

## Features

### Supported Languages
- **JavaScript** (Node.js 18)
- **Python** (Python 3.11)
- **Java** (OpenJDK 17)

### Security Measures
- **Docker Sandboxing**: Each code execution runs in an isolated Docker container
- **Network Isolation**: Containers have no network access (`--network=none`)
- **Read-only Filesystem**: Containers use read-only filesystems with limited temp space
- **Resource Limits**: CPU and memory constraints prevent resource exhaustion
- **Code Validation**: Static analysis prevents dangerous imports and system calls
- **Timeout Protection**: Execution time limits prevent infinite loops
- **User Isolation**: Code runs as `nobody` user with minimal privileges

### Resource Limits
- **Memory**: 128MB per execution
- **CPU**: 0.5 CPU cores
- **Timeout**: 10 seconds per test case
- **Temp Space**: 10MB read-write temporary filesystem
- **Code Size**: 50KB maximum code length

## API Endpoints

### POST /api/v1/execute/run
Executes code against public test cases (for development/testing).

**Request:**
```json
{
  "code": "function solution(input) { return input.trim(); }",
  "language": "javascript",
  "problem_id": 1
}
```

**Response:**
```json
{
  "status": "Accepted",
  "runtime_ms": 45,
  "memory_kb": 1024,
  "test_cases_passed": 2,
  "total_test_cases": 2,
  "test_results": [
    {
      "input": "hello",
      "expected_output": "hello",
      "actual_output": "hello",
      "passed": true,
      "runtime_ms": 23,
      "memory_kb": 512
    }
  ]
}
```

### POST /api/v1/execute/submit
Executes code against all test cases (including hidden ones) for submission.

### POST /api/v1/execute/validate
Validates code without executing it (syntax and security checks).

**Request:**
```json
{
  "code": "function solution(input) { return input; }",
  "language": "javascript"
}
```

**Response:**
```json
{
  "valid": true
}
```

### GET /api/v1/execute/languages
Returns supported programming languages with templates.

## Code Wrapping

The service automatically wraps user code with input/output handling:

### JavaScript
```javascript
const fs = require('fs');
const input = fs.readFileSync('/workspace/input.txt', 'utf8').trim();

// User's solution code here

try {
    const result = solution(input);
    console.log(result);
} catch (error) {
    console.error('Runtime Error:', error.message);
}
```

### Python
```python
import sys

with open('/workspace/input.txt', 'r') as f:
    input_data = f.read().strip()

# User's solution code here

try:
    result = solution(input_data)
    print(result)
except Exception as error:
    print(f'Runtime Error: {error}', file=sys.stderr)
```

### Java
```java
import java.io.*;
import java.util.*;

public class Solution {
    // User's solution code here
    
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
```

## Security Validation

The service performs static analysis to detect dangerous patterns:

### Blocked Patterns
- File system access: `import os`, `require('fs')`, `open(`, `file(`
- Process execution: `exec(`, `eval(`, `Runtime.getRuntime`, `ProcessBuilder`
- Network access: `require('net')`, `require('http')`, `import socket`
- System calls: `System.exit`, `process.exit`, `__import__`

### Code Length Limits
- Maximum 50KB per submission
- Prevents extremely large code submissions

## Docker Configuration

The service uses the following Docker images:
- **JavaScript**: `node:18-alpine`
- **Python**: `python:3.11-alpine`
- **Java**: `openjdk:17-alpine`

### Container Security
```bash
docker run \
  --rm \
  --network=none \
  --read-only \
  --tmpfs /tmp:rw,noexec,nosuid,size=10m \
  --memory=128m \
  --cpus=0.5 \
  --user nobody \
  -v /path/to/code:/workspace:ro \
  -w /workspace \
  node:18-alpine timeout 10s node solution.js
```

## Testing

### Unit Tests
```bash
go test ./pkg/execution/...
```

### Integration Tests (requires Docker)
```bash
go test -tags=integration ./pkg/execution/...
```

## Error Handling

The service handles various error conditions:

- **Compile Error**: Syntax errors in user code
- **Runtime Error**: Exceptions during execution
- **Time Limit Exceeded**: Code execution timeout
- **Memory Limit Exceeded**: Memory usage exceeds limits
- **Wrong Answer**: Output doesn't match expected result
- **Internal Error**: System or Docker errors

## Performance Considerations

- **Parallel Execution**: Multiple test cases can be executed concurrently
- **Container Reuse**: Docker containers are created per execution (not reused for security)
- **Cleanup**: Temporary files and containers are automatically cleaned up
- **Resource Monitoring**: Memory and CPU usage are tracked and limited

## Future Enhancements

- Support for additional languages (C++, Go, Rust)
- More sophisticated memory tracking using Docker stats API
- Code complexity analysis and optimization suggestions
- Execution result caching for identical submissions
- Advanced security scanning with static analysis tools