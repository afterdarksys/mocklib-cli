package main

import (
	"encoding/base64"
	"encoding/json"
	"flag"
	"fmt"
	"os"
)

func lambdaPrintUsage() {
	fmt.Print(`Usage: mocklib-cli lambda <action> [flags]

Actions:
  create    Create a Lambda function
  list      List all Lambda functions
  get       Get a Lambda function
  invoke    Invoke a Lambda function
  delete    Delete a Lambda function

Flags:
  create:
    --name      Function name (required)
    --runtime   Runtime, e.g. python3.11, nodejs18.x (required)
    --handler   Handler, e.g. index.handler (required)
    --zip       Path to deployment ZIP file (required)

  get / delete:
    --name      Function name (required)

  invoke:
    --name      Function name (required)
    --payload   JSON payload string (default: {})
`)
}

func runLambda(args []string) {
	if len(args) == 0 || args[0] == "--help" || args[0] == "-h" {
		lambdaPrintUsage()
		os.Exit(0)
	}

	action := args[0]
	rest := args[1:]

	switch action {
	case "create":
		fs := flag.NewFlagSet("create", flag.ExitOnError)
		name := fs.String("name", "", "Function name")
		runtime := fs.String("runtime", "", "Runtime (e.g. python3.11, nodejs18.x)")
		handler := fs.String("handler", "", "Handler (e.g. index.handler)")
		zipPath := fs.String("zip", "", "Path to deployment ZIP file")
		fs.Parse(rest)
		requireArg("name", *name)
		requireArg("runtime", *runtime)
		requireArg("handler", *handler)
		requireArg("zip", *zipPath)

		zipData, err := os.ReadFile(*zipPath)
		if err != nil {
			fatal("read zip file: %v", err)
		}
		b64Zip := base64.StdEncoding.EncodeToString(zipData)

		body := map[string]interface{}{
			"FunctionName": *name,
			"Runtime":      *runtime,
			"Handler":      *handler,
			"Role":         "arn:aws:iam::123456789012:role/mock-role",
			"Code": map[string]string{
				"ZipFile": b64Zip,
			},
		}
		resp, err := makeJSONRequest("POST", "/lambda/2015-03-31/functions", body, nil)
		if err != nil {
			fatal("%v", err)
		}
		printJSON(resp)

	case "list":
		resp, err := makeJSONRequest("GET", "/lambda/2015-03-31/functions", nil, nil)
		if err != nil {
			fatal("%v", err)
		}
		printJSON(resp)

	case "get":
		fs := flag.NewFlagSet("get", flag.ExitOnError)
		name := fs.String("name", "", "Function name")
		fs.Parse(rest)
		requireArg("name", *name)

		resp, err := makeJSONRequest("GET", "/lambda/2015-03-31/functions/"+*name, nil, nil)
		if err != nil {
			fatal("%v", err)
		}
		printJSON(resp)

	case "invoke":
		fs := flag.NewFlagSet("invoke", flag.ExitOnError)
		name := fs.String("name", "", "Function name")
		payload := fs.String("payload", "{}", "JSON payload string")
		fs.Parse(rest)
		requireArg("name", *name)

		var payloadVal interface{}
		if err := json.Unmarshal([]byte(*payload), &payloadVal); err != nil {
			fatal("invalid payload JSON: %v", err)
		}

		resp, err := makeJSONRequest(
			"POST",
			"/lambda/2015-03-31/functions/"+*name+"/invocations",
			payloadVal,
			nil,
		)
		if err != nil {
			fatal("%v", err)
		}
		printJSON(resp)

	case "delete":
		fs := flag.NewFlagSet("delete", flag.ExitOnError)
		name := fs.String("name", "", "Function name")
		fs.Parse(rest)
		requireArg("name", *name)

		resp, err := makeJSONRequest("DELETE", "/lambda/2015-03-31/functions/"+*name, nil, nil)
		if err != nil {
			fatal("%v", err)
		}
		if resp != nil {
			printJSON(resp)
		} else {
			fmt.Printf("Deleted Lambda function: %s\n", *name)
		}

	default:
		fmt.Fprintf(os.Stderr, "error: unknown lambda action %q\n\n", action)
		lambdaPrintUsage()
		os.Exit(1)
	}
}
