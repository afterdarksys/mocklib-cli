package main

import (
	"flag"
	"fmt"
	"os"
)

func ociPrintUsage() {
	fmt.Print(`Usage: mocklib-cli oci <action> [flags]

Actions:
  get-namespace      Get the Object Storage namespace
  list-buckets       List buckets in a namespace
  create-bucket      Create a bucket in a namespace
  list-instances     List compute instances
  create-instance    Create a compute instance
  stop-instance      Stop a compute instance
  start-instance     Start a compute instance
  delete-instance    Delete a compute instance
  create-vcn         Create a Virtual Cloud Network
  list-vcns          List all VCNs

Flags:
  list-buckets:
    --ns      Object Storage namespace (required)

  create-bucket:
    --ns      Object Storage namespace (required)
    --name    Bucket name (required)

  create-instance:
    --shape   Instance shape, e.g. VM.Standard2.1 (required)
    --image   Image OCID (required)
    --subnet  Subnet OCID (required)

  stop-instance / start-instance / delete-instance:
    --id      Instance OCID (required)

  create-vcn:
    --cidr    CIDR block, e.g. 10.0.0.0/16 (required)
`)
}

func runOCI(args []string) {
	if len(args) == 0 || args[0] == "--help" || args[0] == "-h" {
		ociPrintUsage()
		os.Exit(0)
	}

	action := args[0]
	rest := args[1:]

	switch action {
	case "get-namespace":
		resp, err := makeJSONRequest("GET", "/n/", nil, nil)
		if err != nil {
			fatal("%v", err)
		}
		printJSON(resp)

	case "list-buckets":
		fs := flag.NewFlagSet("list-buckets", flag.ExitOnError)
		ns := fs.String("ns", "", "Object Storage namespace")
		fs.Parse(rest)
		requireArg("ns", *ns)

		resp, err := makeJSONRequest("GET", "/n/"+*ns+"/b/", nil, nil)
		if err != nil {
			fatal("%v", err)
		}
		printJSON(resp)

	case "create-bucket":
		fs := flag.NewFlagSet("create-bucket", flag.ExitOnError)
		ns := fs.String("ns", "", "Object Storage namespace")
		name := fs.String("name", "", "Bucket name")
		fs.Parse(rest)
		requireArg("ns", *ns)
		requireArg("name", *name)

		resp, err := makeJSONRequest("POST", "/n/"+*ns+"/b/",
			map[string]interface{}{"name": *name}, nil)
		if err != nil {
			fatal("%v", err)
		}
		printJSON(resp)

	case "list-instances":
		resp, err := makeJSONRequest("GET", "/20160918/instances", nil, nil)
		if err != nil {
			fatal("%v", err)
		}
		printJSON(resp)

	case "create-instance":
		fs := flag.NewFlagSet("create-instance", flag.ExitOnError)
		shape := fs.String("shape", "", "Instance shape")
		image := fs.String("image", "", "Image OCID")
		subnet := fs.String("subnet", "", "Subnet OCID")
		fs.Parse(rest)
		requireArg("shape", *shape)
		requireArg("image", *image)
		requireArg("subnet", *subnet)

		resp, err := makeJSONRequest("POST", "/20160918/instances",
			map[string]interface{}{
				"shape":              *shape,
				"imageId":            *image,
				"subnetId":           *subnet,
				"availabilityDomain": "AD-1",
			}, nil)
		if err != nil {
			fatal("%v", err)
		}
		printJSON(resp)

	case "stop-instance":
		fs := flag.NewFlagSet("stop-instance", flag.ExitOnError)
		id := fs.String("id", "", "Instance OCID")
		fs.Parse(rest)
		requireArg("id", *id)

		resp, err := makeJSONRequest("POST", "/20160918/instances/"+*id+"/actions/stop", nil, nil)
		if err != nil {
			fatal("%v", err)
		}
		printJSON(resp)

	case "start-instance":
		fs := flag.NewFlagSet("start-instance", flag.ExitOnError)
		id := fs.String("id", "", "Instance OCID")
		fs.Parse(rest)
		requireArg("id", *id)

		resp, err := makeJSONRequest("POST", "/20160918/instances/"+*id+"/actions/start", nil, nil)
		if err != nil {
			fatal("%v", err)
		}
		printJSON(resp)

	case "delete-instance":
		fs := flag.NewFlagSet("delete-instance", flag.ExitOnError)
		id := fs.String("id", "", "Instance OCID")
		fs.Parse(rest)
		requireArg("id", *id)

		resp, err := makeJSONRequest("DELETE", "/20160918/instances/"+*id, nil, nil)
		if err != nil {
			fatal("%v", err)
		}
		if resp != nil {
			printJSON(resp)
		} else {
			fmt.Printf("Deleted OCI instance: %s\n", *id)
		}

	case "create-vcn":
		fs := flag.NewFlagSet("create-vcn", flag.ExitOnError)
		cidr := fs.String("cidr", "", "CIDR block")
		fs.Parse(rest)
		requireArg("cidr", *cidr)

		resp, err := makeJSONRequest("POST", "/20160918/vcns",
			map[string]interface{}{"cidrBlock": *cidr}, nil)
		if err != nil {
			fatal("%v", err)
		}
		printJSON(resp)

	case "list-vcns":
		resp, err := makeJSONRequest("GET", "/20160918/vcns", nil, nil)
		if err != nil {
			fatal("%v", err)
		}
		printJSON(resp)

	default:
		fmt.Fprintf(os.Stderr, "error: unknown oci action %q\n\n", action)
		ociPrintUsage()
		os.Exit(1)
	}
}
