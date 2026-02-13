package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

// VPC command implementations

func vpcCreate(cidrBlock string, tags ...string) error {
	reqBody := map[string]interface{}{
		"Action":             "CreateVpc",
		"CidrBlock":          cidrBlock,
		"EnableDnsHostnames": true,
		"EnableDnsSupport":   true,
	}

	// Parse optional tags (format: key=value key2=value2)
	if len(tags) > 0 {
		tagsMap := make(map[string]string)
		for _, tag := range tags {
			// Simple parsing - could be enhanced
			tagsMap["Name"] = tag
		}
		reqBody["Tags"] = tagsMap
	}

	resp, err := makeRequest("POST", "/aws/vpc", reqBody)
	if err != nil {
		return err
	}

	// Print VPC ID for easy capture
	fmt.Println(resp["VpcId"])
	return nil
}

func vpcDelete(vpcID string) error {
	reqBody := map[string]interface{}{
		"Action": "DeleteVpc",
		"VpcId":  vpcID,
	}

	_, err := makeRequest("POST", "/aws/vpc", reqBody)
	if err != nil {
		return err
	}

	fmt.Printf("Deleted VPC: %s\n", vpcID)
	return nil
}

func vpcList() error {
	reqBody := map[string]interface{}{
		"Action": "DescribeVpcs",
	}

	resp, err := makeRequest("POST", "/aws/vpc", reqBody)
	if err != nil {
		return err
	}

	printJSON(resp)
	return nil
}

// makeRequest is a helper to make authenticated API requests
func makeRequest(method, path string, body interface{}) (map[string]interface{}, error) {
	url := apiURL + path

	var reqBody io.Reader
	if body != nil {
		jsonData, err := json.Marshal(body)
		if err != nil {
			return nil, fmt.Errorf("marshal request: %w", err)
		}
		reqBody = bytes.NewBuffer(jsonData)
	}

	req, err := http.NewRequest(method, url, reqBody)
	if err != nil {
		return nil, fmt.Errorf("create request: %w", err)
	}

	req.Header.Set("Authorization", "Bearer "+apiKey)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("User-Agent", "mocklib-cli/0.1.0")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("send request: %w", err)
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("read response: %w", err)
	}

	if resp.StatusCode >= 400 {
		return nil, fmt.Errorf("API error (%d): %s", resp.StatusCode, string(respBody))
	}

	var result map[string]interface{}
	if err := json.Unmarshal(respBody, &result); err != nil {
		return nil, fmt.Errorf("unmarshal response: %w", err)
	}

	return result, nil
}
