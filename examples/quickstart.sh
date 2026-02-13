#!/bin/bash
# MockLib CLI - Quick Start Example

set -e  # Exit on error

# Ensure API key is set
if [ -z "$MOCKFACTORY_API_KEY" ]; then
    echo "Error: MOCKFACTORY_API_KEY environment variable required"
    echo "Set it with: export MOCKFACTORY_API_KEY='mf_...'"
    exit 1
fi

echo "MockLib CLI - Quick Start"
echo "=========================="
echo ""

# Assuming mocklib binary is in PATH or current directory
MOCKLIB="${MOCKLIB:-./mocklib}"

# Create VPC
echo "Creating VPC..."
VPC_ID=$($MOCKLIB mocklib_vpc_create "10.0.0.0/16")
echo "✓ Created VPC: $VPC_ID"
echo ""

# Create Lambda function
echo "Creating Lambda function..."
LAMBDA_NAME=$($MOCKLIB mocklib_lambda_create "demo-function" "python3.9" "256")
echo "✓ Created Lambda: $LAMBDA_NAME"
echo ""

# Create DynamoDB table
echo "Creating DynamoDB table..."
TABLE_NAME=$($MOCKLIB mocklib_dynamodb_create_table "users" "user_id" "S")
echo "✓ Created DynamoDB table: $TABLE_NAME"
echo ""

# Create SQS queue
echo "Creating SQS queue..."
QUEUE_URL=$($MOCKLIB mocklib_sqs_create_queue "background-jobs" "30")
echo "✓ Created SQS queue: $QUEUE_URL"
echo ""

# Put item in DynamoDB
echo "Adding user to DynamoDB..."
$MOCKLIB mocklib_dynamodb_put_item "$TABLE_NAME" '{"user_id":"123","name":"John Doe","email":"john@example.com"}'
echo "✓ Added user"
echo ""

# Get item from DynamoDB
echo "Retrieving user from DynamoDB..."
$MOCKLIB mocklib_dynamodb_get_item "$TABLE_NAME" '{"user_id":"123"}'
echo ""

# Send SQS message
echo "Sending message to SQS..."
MSG_ID=$($MOCKLIB mocklib_sqs_send_message "$QUEUE_URL" "Hello from MockLib CLI!")
echo "✓ Sent message: $MSG_ID"
echo ""

# Receive SQS messages
echo "Receiving messages from SQS..."
$MOCKLIB mocklib_sqs_receive_messages "$QUEUE_URL" "10"
echo ""

# List all VPCs
echo "Listing all VPCs..."
$MOCKLIB mocklib_vpc_list
echo ""

echo "✅ Demo complete!"
echo ""
echo "Cleanup commands:"
echo "  $MOCKLIB mocklib_vpc_delete '$VPC_ID'"
echo "  $MOCKLIB mocklib_lambda_delete '$LAMBDA_NAME'"
echo "  $MOCKLIB mocklib_dynamodb_delete_table '$TABLE_NAME'"
echo "  $MOCKLIB mocklib_sqs_delete_queue '$QUEUE_URL'"
