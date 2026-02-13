# MockLib CLI

Command-line interface for MockFactory cloud emulation, powered by [gosh](https://github.com/mumoshu/gosh).

## Features

- **Standalone CLI**: Run commands directly: `mocklib vpc create 10.0.0.0/16`
- **Shell-sourceable**: Load functions into your shell environment
- **Hot Reloading**: Automatic rebuild during development with `go run`
- **Single Binary**: No dependencies, just download and run
- **Type-Safe**: Written in Go, leverages our Go SDK internally

## Installation

### Option 1: Download Binary (Recommended)

```bash
# Download latest release
curl -LO https://github.com/mockfactory/mocklib-cli/releases/latest/download/mocklib
chmod +x mocklib
sudo mv mocklib /usr/local/bin/
```

### Option 2: Install from Source

```bash
go install github.com/mockfactory/mocklib-cli@latest
```

### Option 3: Build Locally

```bash
git clone https://github.com/mockfactory/mocklib-cli
cd mocklib-cli
go build -o mocklib
```

## Configuration

Set your API key as an environment variable:

```bash
export MOCKFACTORY_API_KEY="mf_..."
```

Optional: Set custom API URL:

```bash
export MOCKFACTORY_API_URL="https://api.mockfactory.io/v1"
```

## Usage Modes

### Mode 1: Direct CLI Commands

Run the binary with function names and arguments:

```bash
# Create VPC
mocklib mocklib_vpc_create "10.0.0.0/16"

# Create Lambda function
mocklib mocklib_lambda_create "my-function" "python3.9" "256"

# Create DynamoDB table
mocklib mocklib_dynamodb_create_table "users" "user_id"

# Create SQS queue
mocklib mocklib_sqs_create_queue "background-jobs"
```

### Mode 2: Shell-Sourced Functions (gosh style)

Source the CLI to load functions into your current shell:

```bash
# Start interactive shell (with hot reload)
./mocklib

# Now use functions directly in your shell:
VPC_ID=$(mocklib_vpc_create "10.0.0.0/16")
echo "Created VPC: $VPC_ID"

LAMBDA=$(mocklib_lambda_create "api-handler" "nodejs18.x" "512")
echo "Created Lambda: $LAMBDA"

# List all VPCs
mocklib_vpc_list
```

### Mode 3: Shell Scripts

Use in bash scripts:

```bash
#!/bin/bash
source <(mocklib shell)  # Future: load functions into script

VPC_ID=$(mocklib_vpc_create "10.0.0.0/16")
SUBNET_ID=$(mocklib_subnet_create "$VPC_ID" "10.0.1.0/24")

echo "Infrastructure ready!"
```

## Available Commands

### VPC Operations

```bash
# Create VPC (prints VPC ID)
mocklib_vpc_create "10.0.0.0/16" ["tag-name"]

# Delete VPC
mocklib_vpc_delete "vpc-abc123"

# List VPCs (JSON output)
mocklib_vpc_list
```

### Lambda Operations

```bash
# Create function (prints function name)
mocklib_lambda_create "function-name" "runtime" [memory_mb]

# Invoke function (prints response)
mocklib_lambda_invoke "function-name" '{"key": "value"}'

# Delete function
mocklib_lambda_delete "function-name"

# List functions
mocklib_lambda_list
```

### DynamoDB Operations

```bash
# Create table (prints table name)
mocklib_dynamodb_create_table "table-name" "partition-key" [key-type]

# Put item
mocklib_dynamodb_put_item "table-name" '{"user_id": "123", "name": "John"}'

# Get item
mocklib_dynamodb_get_item "table-name" '{"user_id": "123"}'

# Delete table
mocklib_dynamodb_delete_table "table-name"
```

### SQS Operations

```bash
# Create queue (prints queue URL)
mocklib_sqs_create_queue "queue-name" [visibility_timeout]

# Send message (prints message ID)
mocklib_sqs_send_message "queue-url" "message body"

# Receive messages
mocklib_sqs_receive_messages "queue-url" [max_messages]

# Delete queue
mocklib_sqs_delete_queue "queue-url"
```

### Storage Operations

```bash
# Create bucket (prints bucket name)
mocklib_storage_create_bucket "bucket-name" ["s3"|"gcs"|"azure"]

# Delete bucket
mocklib_storage_delete_bucket "bucket-name"
```

## Examples

### Quick Infrastructure Setup

```bash
export MOCKFACTORY_API_KEY="mf_..."

# Create infrastructure
VPC=$(mocklib mocklib_vpc_create "10.0.0.0/16")
LAMBDA=$(mocklib mocklib_lambda_create "api" "python3.9" "256")
TABLE=$(mocklib mocklib_dynamodb_create_table "users" "id")
QUEUE=$(mocklib mocklib_sqs_create_queue "jobs")

echo "Infrastructure created!"
echo "VPC: $VPC"
echo "Lambda: $LAMBDA"
echo "Table: $TABLE"
echo "Queue: $QUEUE"
```

### CI/CD Pipeline

```bash
# In GitHub Actions, CircleCI, etc.
- name: Setup test infrastructure
  run: |
    export MOCKFACTORY_API_KEY=${{ secrets.MOCKFACTORY_API_KEY }}

    # Download CLI
    curl -LO https://releases.mockfactory.io/mocklib
    chmod +x mocklib

    # Create test environment
    VPC_ID=$(./mocklib mocklib_vpc_create "172.16.0.0/16")
    echo "TEST_VPC_ID=$VPC_ID" >> $GITHUB_ENV

- name: Run integration tests
  run: npm test
  env:
    TEST_VPC_ID: ${{ env.TEST_VPC_ID }}
```

### Interactive Development Shell

```bash
# Start hot-reloading development shell
go run .

# Your functions are now available and auto-reload on code changes!
VPC_ID=$(mocklib_vpc_create "10.0.0.0/16")
mocklib_lambda_create "test-fn" "python3.9"
```

## Output Formats

- **Resource IDs**: Commands that create resources print just the ID/name for easy capture
- **JSON**: List and get commands return formatted JSON
- **Status**: Delete commands print confirmation messages

```bash
# Capture output
VPC_ID=$(mocklib mocklib_vpc_create "10.0.0.0/16")
# Prints: vpc-abc123

# JSON output
mocklib mocklib_vpc_list
# Prints:
# {
#   "Vpcs": [
#     {
#       "VpcId": "vpc-abc123",
#       "CidrBlock": "10.0.0.0/16",
#       "State": "available"
#     }
#   ]
# }
```

## Error Handling

```bash
# Check exit codes
if ! VPC_ID=$(mocklib mocklib_vpc_create "10.0.0.0/16"); then
    echo "Failed to create VPC"
    exit 1
fi

# API errors are printed to stderr
mocklib mocklib_vpc_delete "invalid-vpc-id"
# Error: API error (404): VPC not found
```

## Development

### Hot Reloading

The gosh framework enables hot reloading during development:

```bash
# Run in dev mode
go run .

# Edit any .go file - functions auto-reload!
# No need to restart the shell
```

### Adding New Commands

1. Add function to appropriate file (e.g., `vpc.go`)
2. Export it in `main.go`:
   ```go
   sh.Export("mocklib_new_command", newCommand)
   ```
3. That's it! Hot reload will pick it up

### Testing

```bash
# Run tests
go test ./...

# Test specific function
export MOCKFACTORY_API_KEY="mf_test_..."
go run . mocklib_vpc_create "10.0.0.0/16"
```

## Comparison with Other Tools

| Tool | Type | Dependencies | Hot Reload | Type Safety |
|------|------|--------------|------------|-------------|
| **mocklib-cli** | Go binary + gosh | None | ✅ | ✅ |
| mocklib.sh | Bash script | curl, jq | ❌ | ❌ |
| Python SDK | Library | Python, requests | ❌ | ⚠️ |
| Go SDK | Library | Go runtime | ❌ | ✅ |

**Use mocklib-cli when:**
- You need a standalone tool for CI/CD
- You want shell scriptability without bash complexity
- You're developing/debugging and want hot reload
- You need fast execution (compiled binary)

**Use other SDKs when:**
- You're writing application code (Python, Node.js, Go, PHP)
- You need programmatic control
- You're building integrations

## Architecture

```
mocklib-cli (gosh framework)
    ├── Standalone mode: ./mocklib command args
    ├── Interactive mode: ./mocklib (starts shell)
    └── Uses MockFactory Go SDK internally
        └── API calls to https://api.mockfactory.io/v1
```

The CLI is a thin gosh wrapper around MockFactory's API, providing:
- Argument parsing via gosh's reflection
- Shell environment integration
- Hot reloading during development
- Clean CLI interface

## License

MIT

## Links

- **MockFactory**: https://mockfactory.io
- **Documentation**: https://docs.mockfactory.io
- **Go SDK**: https://github.com/mockfactory/mocklib-go
- **Gosh Framework**: https://github.com/mumoshu/gosh
