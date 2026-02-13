#!/bin/bash
# MockLib CLI - Microservices Architecture Setup
# Creates a complete microservices infrastructure

set -e

echo "Microservices Infrastructure Setup"
echo "===================================="
echo ""

# Configuration
ENV="${ENV:-staging}"
CIDR_BLOCK="10.${ENV:0:1}.0.0/16"

echo "Environment: $ENV"
echo "CIDR Block: $CIDR_BLOCK"
echo ""

# Create VPC
echo "1. Creating VPC..."
VPC_ID=$(./mocklib mocklib_vpc_create "$CIDR_BLOCK")
echo "   ✓ VPC: $VPC_ID"
echo ""

# Create Lambda functions for each microservice
echo "2. Creating microservices..."

API_GATEWAY=$(./mocklib mocklib_lambda_create "${ENV}-api-gateway" "nodejs18.x" "512")
echo "   ✓ API Gateway: $API_GATEWAY"

AUTH_SERVICE=$(./mocklib mocklib_lambda_create "${ENV}-auth-service" "python3.9" "256")
echo "   ✓ Auth Service: $AUTH_SERVICE"

USER_SERVICE=$(./mocklib mocklib_lambda_create "${ENV}-user-service" "python3.9" "256")
echo "   ✓ User Service: $USER_SERVICE"

ORDER_SERVICE=$(./mocklib mocklib_lambda_create "${ENV}-order-service" "nodejs18.x" "512")
echo "   ✓ Order Service: $ORDER_SERVICE"

PAYMENT_SERVICE=$(./mocklib mocklib_lambda_create "${ENV}-payment-service" "python3.9" "512")
echo "   ✓ Payment Service: $PAYMENT_SERVICE"

echo ""

# Create DynamoDB tables
echo "3. Creating databases..."

USERS_TABLE=$(./mocklib mocklib_dynamodb_create_table "${ENV}_users" "user_id" "S")
echo "   ✓ Users Table: $USERS_TABLE"

ORDERS_TABLE=$(./mocklib mocklib_dynamodb_create_table "${ENV}_orders" "order_id" "S")
echo "   ✓ Orders Table: $ORDERS_TABLE"

SESSIONS_TABLE=$(./mocklib mocklib_dynamodb_create_table "${ENV}_sessions" "session_id" "S")
echo "   ✓ Sessions Table: $SESSIONS_TABLE"

echo ""

# Create SQS queues
echo "4. Creating message queues..."

EMAIL_QUEUE=$(./mocklib mocklib_sqs_create_queue "${ENV}-email-queue" "120")
echo "   ✓ Email Queue: $EMAIL_QUEUE"

EVENTS_QUEUE=$(./mocklib mocklib_sqs_create_queue "${ENV}-events-queue" "60")
echo "   ✓ Events Queue: $EVENTS_QUEUE"

DLQ=$(./mocklib mocklib_sqs_create_queue "${ENV}-dead-letter-queue" "300")
echo "   ✓ Dead Letter Queue: $DLQ"

echo ""

# Generate infrastructure config
echo "5. Generating configuration..."

cat > "${ENV}-infrastructure.env" <<EOF
# MockFactory Microservices Infrastructure - $ENV
# Generated: $(date)

VPC_ID=$VPC_ID
CIDR_BLOCK=$CIDR_BLOCK

# Lambda Functions
API_GATEWAY_FUNCTION=$API_GATEWAY
AUTH_SERVICE_FUNCTION=$AUTH_SERVICE
USER_SERVICE_FUNCTION=$USER_SERVICE
ORDER_SERVICE_FUNCTION=$ORDER_SERVICE
PAYMENT_SERVICE_FUNCTION=$PAYMENT_SERVICE

# DynamoDB Tables
USERS_TABLE=$USERS_TABLE
ORDERS_TABLE=$ORDERS_TABLE
SESSIONS_TABLE=$SESSIONS_TABLE

# SQS Queues
EMAIL_QUEUE_URL=$EMAIL_QUEUE
EVENTS_QUEUE_URL=$EVENTS_QUEUE
DLQ_URL=$DLQ

# Metadata
ENVIRONMENT=$ENV
CREATED_AT=$(date -u +"%Y-%m-%dT%H:%M:%SZ")
EOF

echo "   ✓ Configuration saved to: ${ENV}-infrastructure.env"
echo ""

echo "✅ Infrastructure setup complete!"
echo ""
echo "Load configuration:"
echo "  source ${ENV}-infrastructure.env"
echo ""
echo "Cost estimate: ~\$0.15/hour for this infrastructure"
