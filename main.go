package main

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/mumoshu/gosh"
)

const (
	defaultAPIURL = "https://api.mockfactory.io/v1"
)

var apiKey string
var apiURL string

func main() {
	// Get configuration from environment
	apiKey = os.Getenv("MOCKFACTORY_API_KEY")
	apiURL = os.Getenv("MOCKFACTORY_API_URL")
	if apiURL == "" {
		apiURL = defaultAPIURL
	}

	if apiKey == "" {
		fmt.Fprintln(os.Stderr, "Error: MOCKFACTORY_API_KEY environment variable required")
		fmt.Fprintln(os.Stderr, "Set it with: export MOCKFACTORY_API_KEY='mf_...'")
		os.Exit(1)
	}

	// Create gosh shell
	sh := gosh.NewShell()

	// Export VPC commands
	sh.Export("mocklib_vpc_create", vpcCreate)
	sh.Export("mocklib_vpc_delete", vpcDelete)
	sh.Export("mocklib_vpc_list", vpcList)

	// Export Lambda commands
	sh.Export("mocklib_lambda_create", lambdaCreate)
	sh.Export("mocklib_lambda_invoke", lambdaInvoke)
	sh.Export("mocklib_lambda_delete", lambdaDelete)
	sh.Export("mocklib_lambda_list", lambdaList)

	// Export DynamoDB commands
	sh.Export("mocklib_dynamodb_create_table", dynamodbCreateTable)
	sh.Export("mocklib_dynamodb_put_item", dynamodbPutItem)
	sh.Export("mocklib_dynamodb_get_item", dynamodbGetItem)
	sh.Export("mocklib_dynamodb_delete_table", dynamodbDeleteTable)

	// Export SQS commands
	sh.Export("mocklib_sqs_create_queue", sqsCreateQueue)
	sh.Export("mocklib_sqs_send_message", sqsSendMessage)
	sh.Export("mocklib_sqs_receive_messages", sqsReceiveMessages)
	sh.Export("mocklib_sqs_delete_queue", sqsDeleteQueue)

	// Export Storage commands
	sh.Export("mocklib_storage_create_bucket", storageCreateBucket)
	sh.Export("mocklib_storage_delete_bucket", storageDeleteBucket)

	// Run the shell
	if err := sh.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}

// printJSON prints a value as formatted JSON
func printJSON(v interface{}) {
	data, _ := json.MarshalIndent(v, "", "  ")
	fmt.Println(string(data))
}
