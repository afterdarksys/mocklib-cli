package main

import (
	"flag"
	"fmt"
	"os"
)

func storagePrintUsage() {
	fmt.Print(`Usage: mocklib-cli storage <action> [flags]

Actions:
  create-bucket   Create a storage bucket
  upload          Upload an object to a bucket
  download        Download an object from a bucket
  delete-object   Delete an object from a bucket
  delete-bucket   Delete a bucket

Flags:
  create-bucket:
    --name    Bucket name (required)

  upload:
    --bucket  Bucket name (required)
    --key     Object key / path (required)
    --file    Local file path to upload (required)

  download:
    --bucket  Bucket name (required)
    --key     Object key / path (required)

  delete-object:
    --bucket  Bucket name (required)
    --key     Object key / path (required)

  delete-bucket:
    --name    Bucket name (required)
`)
}

func runStorage(args []string) {
	if len(args) == 0 || args[0] == "--help" || args[0] == "-h" {
		storagePrintUsage()
		os.Exit(0)
	}

	action := args[0]
	rest := args[1:]

	switch action {
	case "create-bucket":
		fs := flag.NewFlagSet("create-bucket", flag.ExitOnError)
		name := fs.String("name", "", "Bucket name")
		fs.Parse(rest)
		requireArg("name", *name)

		resp, err := makeFormRequest("/storage/bucket", map[string]string{
			"Action":     "CreateBucket",
			"BucketName": *name,
		})
		if err != nil {
			fatal("%v", err)
		}
		printJSON(resp)

	case "upload":
		fs := flag.NewFlagSet("upload", flag.ExitOnError)
		bucket := fs.String("bucket", "", "Bucket name")
		key := fs.String("key", "", "Object key")
		filePath := fs.String("file", "", "Local file path")
		fs.Parse(rest)
		requireArg("bucket", *bucket)
		requireArg("key", *key)
		requireArg("file", *filePath)

		data, err := os.ReadFile(*filePath)
		if err != nil {
			fatal("read file: %v", err)
		}

		resp, err := makeJSONRequest("PUT", "/storage/bucket/"+*bucket+"/"+*key,
			map[string]interface{}{
				"Body": string(data),
			}, nil)
		if err != nil {
			fatal("%v", err)
		}
		printJSON(resp)

	case "download":
		fs := flag.NewFlagSet("download", flag.ExitOnError)
		bucket := fs.String("bucket", "", "Bucket name")
		key := fs.String("key", "", "Object key")
		fs.Parse(rest)
		requireArg("bucket", *bucket)
		requireArg("key", *key)

		resp, err := makeJSONRequest("GET", "/storage/bucket/"+*bucket+"/"+*key, nil, nil)
		if err != nil {
			fatal("%v", err)
		}
		printJSON(resp)

	case "delete-object":
		fs := flag.NewFlagSet("delete-object", flag.ExitOnError)
		bucket := fs.String("bucket", "", "Bucket name")
		key := fs.String("key", "", "Object key")
		fs.Parse(rest)
		requireArg("bucket", *bucket)
		requireArg("key", *key)

		resp, err := makeJSONRequest("DELETE", "/storage/bucket/"+*bucket+"/"+*key, nil, nil)
		if err != nil {
			fatal("%v", err)
		}
		if resp != nil {
			printJSON(resp)
		} else {
			fmt.Printf("Deleted object %s from bucket %s\n", *key, *bucket)
		}

	case "delete-bucket":
		fs := flag.NewFlagSet("delete-bucket", flag.ExitOnError)
		name := fs.String("name", "", "Bucket name")
		fs.Parse(rest)
		requireArg("name", *name)

		resp, err := makeFormRequest("/storage/bucket", map[string]string{
			"Action":     "DeleteBucket",
			"BucketName": *name,
		})
		if err != nil {
			fatal("%v", err)
		}
		if resp != nil {
			printJSON(resp)
		} else {
			fmt.Printf("Deleted bucket: %s\n", *name)
		}

	default:
		fmt.Fprintf(os.Stderr, "error: unknown storage action %q\n\n", action)
		storagePrintUsage()
		os.Exit(1)
	}
}
