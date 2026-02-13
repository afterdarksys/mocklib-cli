package main

import (
	"fmt"
	"strconv"
)

// Lambda command implementations

func lambdaCreate(functionName, runtime string, memoryMB ...string) error {
	memory := 256
	if len(memoryMB) > 0 {
		var err error
		memory, err = strconv.Atoi(memoryMB[0])
		if err != nil {
			return fmt.Errorf("invalid memory_mb: %w", err)
		}
	}

	reqBody := map[string]interface{}{
		"Action":       "CreateFunction",
		"FunctionName": functionName,
		"Runtime":      runtime,
		"Handler":      "index.handler",
		"MemoryMB":     memory,
		"Timeout":      30,
	}

	resp, err := makeRequest("POST", "/aws/lambda", reqBody)
	if err != nil {
		return err
	}

	// Print function name for easy capture
	fmt.Println(resp["FunctionName"])
	return nil
}

func lambdaInvoke(functionName, payload string) error {
	reqBody := map[string]interface{}{
		"Action":       "Invoke",
		"FunctionName": functionName,
		"Payload":      payload,
	}

	resp, err := makeRequest("POST", "/aws/lambda", reqBody)
	if err != nil {
		return err
	}

	printJSON(resp)
	return nil
}

func lambdaDelete(functionName string) error {
	reqBody := map[string]interface{}{
		"Action":       "DeleteFunction",
		"FunctionName": functionName,
	}

	_, err := makeRequest("POST", "/aws/lambda", reqBody)
	if err != nil {
		return err
	}

	fmt.Printf("Deleted Lambda function: %s\n", functionName)
	return nil
}

func lambdaList() error {
	reqBody := map[string]interface{}{
		"Action": "ListFunctions",
	}

	resp, err := makeRequest("POST", "/aws/lambda", reqBody)
	if err != nil {
		return err
	}

	printJSON(resp)
	return nil
}
