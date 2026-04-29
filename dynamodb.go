package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"
)

func dynamodbPrintUsage() {
	fmt.Print(`Usage: mocklib-cli dynamodb <action> [flags]

Actions:
  create-table   Create a DynamoDB table
  put-item       Put an item into a table
  get-item       Get an item from a table
  update-item    Update an item in a table
  delete-item    Delete an item from a table
  scan           Scan all items in a table
  query          Query a table with a condition

Flags:
  create-table:
    --name   Table name (required)
    --key    Partition key name (required)

  put-item:
    --table  Table name (required)
    --item   Item as JSON, e.g. '{"pk":{"S":"val"}}' (required)

  get-item:
    --table  Table name (required)
    --key    Key as JSON, e.g. '{"pk":{"S":"val"}}' (required)

  update-item:
    --table    Table name (required)
    --key      Key as JSON (required)
    --updates  UpdateExpression string (required)

  delete-item:
    --table  Table name (required)
    --key    Key as JSON (required)

  scan:
    --table  Table name (required)

  query:
    --table      Table name (required)
    --condition  KeyConditionExpression string (required)
`)
}

// dynamodbRequest sends a DynamoDB request using JSON body and X-Amz-Target header.
func dynamodbRequest(target string, body map[string]interface{}) (map[string]interface{}, error) {
	return makeJSONRequest("POST", "/aws/dynamodb", body, map[string]string{
		"X-Amz-Target": "DynamoDB_20120810." + target,
	})
}

func runDynamoDB(args []string) {
	if len(args) == 0 || args[0] == "--help" || args[0] == "-h" {
		dynamodbPrintUsage()
		os.Exit(0)
	}

	action := args[0]
	rest := args[1:]

	switch action {
	case "create-table":
		fs := flag.NewFlagSet("create-table", flag.ExitOnError)
		name := fs.String("name", "", "Table name")
		key := fs.String("key", "", "Partition key name")
		fs.Parse(rest)
		requireArg("name", *name)
		requireArg("key", *key)

		resp, err := dynamodbRequest("CreateTable", map[string]interface{}{
			"TableName": *name,
			"AttributeDefinitions": []map[string]string{
				{"AttributeName": *key, "AttributeType": "S"},
			},
			"KeySchema": []map[string]string{
				{"AttributeName": *key, "KeyType": "HASH"},
			},
			"BillingMode": "PAY_PER_REQUEST",
		})
		if err != nil {
			fatal("%v", err)
		}
		printJSON(resp)

	case "put-item":
		fs := flag.NewFlagSet("put-item", flag.ExitOnError)
		table := fs.String("table", "", "Table name")
		item := fs.String("item", "", "Item as JSON")
		fs.Parse(rest)
		requireArg("table", *table)
		requireArg("item", *item)

		var itemVal interface{}
		if err := json.Unmarshal([]byte(*item), &itemVal); err != nil {
			fatal("invalid item JSON: %v", err)
		}
		resp, err := dynamodbRequest("PutItem", map[string]interface{}{
			"TableName": *table,
			"Item":      itemVal,
		})
		if err != nil {
			fatal("%v", err)
		}
		printJSON(resp)

	case "get-item":
		fs := flag.NewFlagSet("get-item", flag.ExitOnError)
		table := fs.String("table", "", "Table name")
		key := fs.String("key", "", "Key as JSON")
		fs.Parse(rest)
		requireArg("table", *table)
		requireArg("key", *key)

		var keyVal interface{}
		if err := json.Unmarshal([]byte(*key), &keyVal); err != nil {
			fatal("invalid key JSON: %v", err)
		}
		resp, err := dynamodbRequest("GetItem", map[string]interface{}{
			"TableName": *table,
			"Key":       keyVal,
		})
		if err != nil {
			fatal("%v", err)
		}
		printJSON(resp)

	case "update-item":
		fs := flag.NewFlagSet("update-item", flag.ExitOnError)
		table := fs.String("table", "", "Table name")
		key := fs.String("key", "", "Key as JSON")
		updates := fs.String("updates", "", "UpdateExpression string")
		fs.Parse(rest)
		requireArg("table", *table)
		requireArg("key", *key)
		requireArg("updates", *updates)

		var keyVal interface{}
		if err := json.Unmarshal([]byte(*key), &keyVal); err != nil {
			fatal("invalid key JSON: %v", err)
		}
		resp, err := dynamodbRequest("UpdateItem", map[string]interface{}{
			"TableName":        *table,
			"Key":              keyVal,
			"UpdateExpression": *updates,
		})
		if err != nil {
			fatal("%v", err)
		}
		printJSON(resp)

	case "delete-item":
		fs := flag.NewFlagSet("delete-item", flag.ExitOnError)
		table := fs.String("table", "", "Table name")
		key := fs.String("key", "", "Key as JSON")
		fs.Parse(rest)
		requireArg("table", *table)
		requireArg("key", *key)

		var keyVal interface{}
		if err := json.Unmarshal([]byte(*key), &keyVal); err != nil {
			fatal("invalid key JSON: %v", err)
		}
		resp, err := dynamodbRequest("DeleteItem", map[string]interface{}{
			"TableName": *table,
			"Key":       keyVal,
		})
		if err != nil {
			fatal("%v", err)
		}
		printJSON(resp)

	case "scan":
		fs := flag.NewFlagSet("scan", flag.ExitOnError)
		table := fs.String("table", "", "Table name")
		fs.Parse(rest)
		requireArg("table", *table)

		resp, err := dynamodbRequest("Scan", map[string]interface{}{
			"TableName": *table,
		})
		if err != nil {
			fatal("%v", err)
		}
		printJSON(resp)

	case "query":
		fs := flag.NewFlagSet("query", flag.ExitOnError)
		table := fs.String("table", "", "Table name")
		condition := fs.String("condition", "", "KeyConditionExpression string")
		fs.Parse(rest)
		requireArg("table", *table)
		requireArg("condition", *condition)

		resp, err := dynamodbRequest("Query", map[string]interface{}{
			"TableName":              *table,
			"KeyConditionExpression": *condition,
		})
		if err != nil {
			fatal("%v", err)
		}
		printJSON(resp)

	default:
		fmt.Fprintf(os.Stderr, "error: unknown dynamodb action %q\n\n", action)
		dynamodbPrintUsage()
		os.Exit(1)
	}
}
