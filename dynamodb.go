package main

import (
	"encoding/json"
	"fmt"
)

// DynamoDB command implementations

func dynamodbCreateTable(tableName, partitionKey string, partitionKeyType ...string) error {
	keyType := "S" // Default to String
	if len(partitionKeyType) > 0 {
		keyType = partitionKeyType[0]
	}

	reqBody := map[string]interface{}{
		"Action":           "CreateTable",
		"TableName":        tableName,
		"PartitionKey":     partitionKey,
		"PartitionKeyType": keyType,
	}

	resp, err := makeRequest("POST", "/aws/dynamodb", reqBody)
	if err != nil {
		return err
	}

	// Print table name for easy capture
	fmt.Println(resp["TableName"])
	return nil
}

func dynamodbPutItem(tableName, itemJSON string) error {
	var item map[string]interface{}
	if err := json.Unmarshal([]byte(itemJSON), &item); err != nil {
		return fmt.Errorf("invalid JSON: %w", err)
	}

	reqBody := map[string]interface{}{
		"Action":    "PutItem",
		"TableName": tableName,
		"Item":      item,
	}

	resp, err := makeRequest("POST", "/aws/dynamodb", reqBody)
	if err != nil {
		return err
	}

	printJSON(resp)
	return nil
}

func dynamodbGetItem(tableName, keyJSON string) error {
	var key map[string]interface{}
	if err := json.Unmarshal([]byte(keyJSON), &key); err != nil {
		return fmt.Errorf("invalid JSON: %w", err)
	}

	reqBody := map[string]interface{}{
		"Action":    "GetItem",
		"TableName": tableName,
		"Key":       key,
	}

	resp, err := makeRequest("POST", "/aws/dynamodb", reqBody)
	if err != nil {
		return err
	}

	printJSON(resp)
	return nil
}

func dynamodbDeleteTable(tableName string) error {
	reqBody := map[string]interface{}{
		"Action":    "DeleteTable",
		"TableName": tableName,
	}

	_, err := makeRequest("POST", "/aws/dynamodb", reqBody)
	if err != nil {
		return err
	}

	fmt.Printf("Deleted DynamoDB table: %s\n", tableName)
	return nil
}
