package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"strings"
)

const (
	defaultAPIURL = "https://api.mockfactory.io/v1"
	cliVersion    = "1.0.0"
)

var (
	apiKey string
	apiURL string
)

// ─── HTTP helpers ────────────────────────────────────────────────────────────

// makeFormRequest sends a POST with form-encoded body.
func makeFormRequest(path string, fields map[string]string) (map[string]interface{}, error) {
	form := url.Values{}
	for k, v := range fields {
		form.Set(k, v)
	}
	req, err := http.NewRequest("POST", apiURL+path, strings.NewReader(form.Encode()))
	if err != nil {
		return nil, fmt.Errorf("create request: %w", err)
	}
	req.Header.Set("Authorization", "Bearer "+apiKey)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("User-Agent", "mocklib-cli/"+cliVersion)
	return doRequest(req)
}

// makeJSONRequest sends a request with a JSON body and optional extra headers.
func makeJSONRequest(method, path string, body interface{}, extraHeaders map[string]string) (map[string]interface{}, error) {
	var r io.Reader
	if body != nil {
		data, err := json.Marshal(body)
		if err != nil {
			return nil, fmt.Errorf("marshal request: %w", err)
		}
		r = bytes.NewBuffer(data)
	}
	req, err := http.NewRequest(method, apiURL+path, r)
	if err != nil {
		return nil, fmt.Errorf("create request: %w", err)
	}
	req.Header.Set("Authorization", "Bearer "+apiKey)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("User-Agent", "mocklib-cli/"+cliVersion)
	for k, v := range extraHeaders {
		req.Header.Set(k, v)
	}
	return doRequest(req)
}

func doRequest(req *http.Request) (map[string]interface{}, error) {
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("send request: %w", err)
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("read response: %w", err)
	}
	if resp.StatusCode >= 400 {
		return nil, fmt.Errorf("API error (%d): %s", resp.StatusCode, strings.TrimSpace(string(body)))
	}
	var result map[string]interface{}
	if len(body) == 0 {
		return result, nil
	}
	if err := json.Unmarshal(body, &result); err != nil {
		// Response may be a plain string or non-JSON; wrap it.
		return map[string]interface{}{"response": string(body)}, nil
	}
	return result, nil
}

func printJSON(v interface{}) {
	data, _ := json.MarshalIndent(v, "", "  ")
	fmt.Println(string(data))
}

func fatal(format string, a ...interface{}) {
	fmt.Fprintf(os.Stderr, "error: "+format+"\n", a...)
	os.Exit(1)
}

func requireArg(name, val string) {
	if val == "" {
		fatal("flag --%s is required", name)
	}
}

// ─── Usage ───────────────────────────────────────────────────────────────────

func printUsage() {
	fmt.Print(`mocklib-cli - MockFactory.io command-line interface

Usage:
  mocklib-cli <service> <action> [flags]

Environment:
  MOCKFACTORY_API_KEY   API key (required)
  MOCKFACTORY_API_URL   Base URL (default: https://api.mockfactory.io/v1)

Services:
  sts         AWS Security Token Service
  ec2         AWS Elastic Compute Cloud
  route53     AWS Route 53 DNS
  iam         AWS Identity and Access Management
  lambda      AWS Lambda
  sns         AWS Simple Notification Service
  sqs         AWS Simple Queue Service
  dynamodb    AWS DynamoDB
  vpc         AWS Virtual Private Cloud
  storage     Object storage (S3-compatible)
  oci         Oracle Cloud Infrastructure
  gcp         Google Cloud Platform
  azure       Microsoft Azure

Run 'mocklib-cli <service> --help' for service-specific actions and flags.

Version: ` + cliVersion + "\n")
}

// ─── Main ────────────────────────────────────────────────────────────────────

func main() {
	apiKey = os.Getenv("MOCKFACTORY_API_KEY")
	apiURL = os.Getenv("MOCKFACTORY_API_URL")
	if apiURL == "" {
		apiURL = defaultAPIURL
	}

	if len(os.Args) < 2 || os.Args[1] == "--help" || os.Args[1] == "-h" {
		printUsage()
		os.Exit(0)
	}

	service := os.Args[1]
	args := os.Args[2:]

	// Allow --help / -h on a service without requiring the API key.
	isHelp := len(args) == 0 ||
		args[0] == "--help" || args[0] == "-h" ||
		(len(args) >= 2 && (args[1] == "--help" || args[1] == "-h"))

	if !isHelp && apiKey == "" {
		fmt.Fprintln(os.Stderr, "error: MOCKFACTORY_API_KEY environment variable is required")
		fmt.Fprintln(os.Stderr, "  export MOCKFACTORY_API_KEY='mf_...'")
		os.Exit(1)
	}

	switch service {
	case "sts":
		runSTS(args)
	case "ec2":
		runEC2(args)
	case "route53":
		runRoute53(args)
	case "iam":
		runIAM(args)
	case "lambda":
		runLambda(args)
	case "sns":
		runSNS(args)
	case "sqs":
		runSQS(args)
	case "dynamodb":
		runDynamoDB(args)
	case "vpc":
		runVPC(args)
	case "storage":
		runStorage(args)
	case "oci":
		runOCI(args)
	case "gcp":
		runGCP(args)
	case "azure":
		runAzure(args)
	default:
		fmt.Fprintf(os.Stderr, "error: unknown service %q\n\n", service)
		printUsage()
		os.Exit(1)
	}
}
