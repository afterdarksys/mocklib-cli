#!/bin/bash
# MockLib CLI - CI/CD Pipeline Example
# Demonstrates ephemeral test infrastructure setup and teardown

set -e

echo "CI Pipeline - Setting up test infrastructure"
echo "============================================="

# Trap to ensure cleanup on exit
cleanup() {
    echo ""
    echo "Cleaning up test infrastructure..."

    if [ -n "$VPC_ID" ]; then
        ./mocklib mocklib_vpc_delete "$VPC_ID" || true
    fi

    if [ -n "$LAMBDA_NAME" ]; then
        ./mocklib mocklib_lambda_delete "$LAMBDA_NAME" || true
    fi

    if [ -n "$TABLE_NAME" ]; then
        ./mocklib mocklib_dynamodb_delete_table "$TABLE_NAME" || true
    fi

    if [ -n "$QUEUE_URL" ]; then
        ./mocklib mocklib_sqs_delete_queue "$QUEUE_URL" || true
    fi

    echo "✓ Cleanup complete"
}

trap cleanup EXIT

# Setup phase
echo ""
echo "Phase 1: Infrastructure Setup"
echo "------------------------------"

VPC_ID=$(./mocklib mocklib_vpc_create "172.31.0.0/16")
echo "✓ VPC: $VPC_ID"

LAMBDA_NAME=$(./mocklib mocklib_lambda_create "test-api" "nodejs18.x" "512")
echo "✓ Lambda: $LAMBDA_NAME"

TABLE_NAME=$(./mocklib mocklib_dynamodb_create_table "test_users" "id" "S")
echo "✓ DynamoDB: $TABLE_NAME"

QUEUE_URL=$(./mocklib mocklib_sqs_create_queue "test-queue" "60")
echo "✓ SQS: $QUEUE_URL"

# Export for test suite
export TEST_VPC_ID="$VPC_ID"
export TEST_LAMBDA_NAME="$LAMBDA_NAME"
export TEST_TABLE_NAME="$TABLE_NAME"
export TEST_QUEUE_URL="$QUEUE_URL"

echo ""
echo "Phase 2: Running Tests"
echo "----------------------"

# Run your actual test suite here
# Example: npm test, pytest, go test, etc.
echo "Running integration tests..."

# Simulate test execution
sleep 2

echo "✓ All tests passed"

echo ""
echo "Phase 3: Cleanup"
echo "----------------"
# Cleanup happens automatically via trap

exit 0
