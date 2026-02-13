package main

import (
	"fmt"
	"strconv"
)

// SQS command implementations

func sqsCreateQueue(queueName string, visibilityTimeout ...string) error {
	timeout := 30
	if len(visibilityTimeout) > 0 {
		var err error
		timeout, err = strconv.Atoi(visibilityTimeout[0])
		if err != nil {
			return fmt.Errorf("invalid visibility_timeout: %w", err)
		}
	}

	reqBody := map[string]interface{}{
		"Action":            "CreateQueue",
		"QueueName":         queueName,
		"VisibilityTimeout": timeout,
	}

	resp, err := makeRequest("POST", "/aws/sqs", reqBody)
	if err != nil {
		return err
	}

	// Print queue URL for easy capture
	fmt.Println(resp["QueueUrl"])
	return nil
}

func sqsSendMessage(queueURL, messageBody string) error {
	reqBody := map[string]interface{}{
		"Action":      "SendMessage",
		"QueueUrl":    queueURL,
		"MessageBody": messageBody,
	}

	resp, err := makeRequest("POST", "/aws/sqs", reqBody)
	if err != nil {
		return err
	}

	// Print message ID
	fmt.Println(resp["MessageId"])
	return nil
}

func sqsReceiveMessages(queueURL string, maxMessages ...string) error {
	max := 1
	if len(maxMessages) > 0 {
		var err error
		max, err = strconv.Atoi(maxMessages[0])
		if err != nil {
			return fmt.Errorf("invalid max_messages: %w", err)
		}
	}

	reqBody := map[string]interface{}{
		"Action":              "ReceiveMessage",
		"QueueUrl":            queueURL,
		"MaxNumberOfMessages": max,
	}

	resp, err := makeRequest("POST", "/aws/sqs", reqBody)
	if err != nil {
		return err
	}

	printJSON(resp)
	return nil
}

func sqsDeleteQueue(queueURL string) error {
	reqBody := map[string]interface{}{
		"Action":   "DeleteQueue",
		"QueueUrl": queueURL,
	}

	_, err := makeRequest("POST", "/aws/sqs", reqBody)
	if err != nil {
		return err
	}

	fmt.Printf("Deleted SQS queue: %s\n", queueURL)
	return nil
}
