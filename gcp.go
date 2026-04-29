package main

import (
	"flag"
	"fmt"
	"os"
	"strings"
)

func gcpPrintUsage() {
	fmt.Print(`Usage: mocklib-cli gcp <action> [flags]

Actions:
  list-zones        List compute zones in a project
  create-instance   Create a compute instance
  list-instances    List compute instances in a zone
  delete-instance   Delete a compute instance
  create-network    Create a VPC network
  list-networks     List VPC networks
  create-firewall   Create a firewall rule

Flags:
  list-zones:
    --project   GCP project ID (required)

  create-instance:
    --project      GCP project ID (required)
    --zone         Compute zone, e.g. us-central1-a (required)
    --name         Instance name (required)
    --machine-type Machine type, e.g. n1-standard-1 (required)

  list-instances / delete-instance:
    --project   GCP project ID (required)
    --zone      Compute zone (required)
    --name      Instance name (delete-instance only, required)

  create-network:
    --project   GCP project ID (required)
    --name      Network name (required)

  list-networks:
    --project   GCP project ID (required)

  create-firewall:
    --project   GCP project ID (required)
    --name      Firewall rule name (required)
    --network   Network name or URL (required)
    --protocol  Protocol, e.g. tcp, udp, icmp (required)
    --ports     Comma-separated ports, e.g. 80,443 (required)
`)
}

func gcpBase(project string) string {
	return "/gcp/compute/v1/projects/" + project
}

func runGCP(args []string) {
	if len(args) == 0 || args[0] == "--help" || args[0] == "-h" {
		gcpPrintUsage()
		os.Exit(0)
	}

	action := args[0]
	rest := args[1:]

	switch action {
	case "list-zones":
		fs := flag.NewFlagSet("list-zones", flag.ExitOnError)
		project := fs.String("project", "", "GCP project ID")
		fs.Parse(rest)
		requireArg("project", *project)

		resp, err := makeJSONRequest("GET", gcpBase(*project)+"/zones", nil, nil)
		if err != nil {
			fatal("%v", err)
		}
		printJSON(resp)

	case "create-instance":
		fs := flag.NewFlagSet("create-instance", flag.ExitOnError)
		project := fs.String("project", "", "GCP project ID")
		zone := fs.String("zone", "", "Compute zone")
		name := fs.String("name", "", "Instance name")
		machineType := fs.String("machine-type", "", "Machine type")
		fs.Parse(rest)
		requireArg("project", *project)
		requireArg("zone", *zone)
		requireArg("name", *name)
		requireArg("machine-type", *machineType)

		resp, err := makeJSONRequest(
			"POST",
			gcpBase(*project)+"/zones/"+*zone+"/instances",
			map[string]interface{}{
				"name":        *name,
				"machineType": "zones/" + *zone + "/machineTypes/" + *machineType,
			},
			nil,
		)
		if err != nil {
			fatal("%v", err)
		}
		printJSON(resp)

	case "list-instances":
		fs := flag.NewFlagSet("list-instances", flag.ExitOnError)
		project := fs.String("project", "", "GCP project ID")
		zone := fs.String("zone", "", "Compute zone")
		fs.Parse(rest)
		requireArg("project", *project)
		requireArg("zone", *zone)

		resp, err := makeJSONRequest("GET", gcpBase(*project)+"/zones/"+*zone+"/instances", nil, nil)
		if err != nil {
			fatal("%v", err)
		}
		printJSON(resp)

	case "delete-instance":
		fs := flag.NewFlagSet("delete-instance", flag.ExitOnError)
		project := fs.String("project", "", "GCP project ID")
		zone := fs.String("zone", "", "Compute zone")
		name := fs.String("name", "", "Instance name")
		fs.Parse(rest)
		requireArg("project", *project)
		requireArg("zone", *zone)
		requireArg("name", *name)

		resp, err := makeJSONRequest(
			"DELETE",
			gcpBase(*project)+"/zones/"+*zone+"/instances/"+*name,
			nil, nil,
		)
		if err != nil {
			fatal("%v", err)
		}
		if resp != nil {
			printJSON(resp)
		} else {
			fmt.Printf("Deleted GCP instance: %s\n", *name)
		}

	case "create-network":
		fs := flag.NewFlagSet("create-network", flag.ExitOnError)
		project := fs.String("project", "", "GCP project ID")
		name := fs.String("name", "", "Network name")
		fs.Parse(rest)
		requireArg("project", *project)
		requireArg("name", *name)

		resp, err := makeJSONRequest(
			"POST",
			gcpBase(*project)+"/global/networks",
			map[string]interface{}{
				"name":                  *name,
				"autoCreateSubnetworks": false,
			},
			nil,
		)
		if err != nil {
			fatal("%v", err)
		}
		printJSON(resp)

	case "list-networks":
		fs := flag.NewFlagSet("list-networks", flag.ExitOnError)
		project := fs.String("project", "", "GCP project ID")
		fs.Parse(rest)
		requireArg("project", *project)

		resp, err := makeJSONRequest("GET", gcpBase(*project)+"/global/networks", nil, nil)
		if err != nil {
			fatal("%v", err)
		}
		printJSON(resp)

	case "create-firewall":
		fs := flag.NewFlagSet("create-firewall", flag.ExitOnError)
		project := fs.String("project", "", "GCP project ID")
		name := fs.String("name", "", "Firewall rule name")
		network := fs.String("network", "", "Network name or URL")
		protocol := fs.String("protocol", "", "Protocol (tcp, udp, icmp)")
		ports := fs.String("ports", "", "Comma-separated ports, e.g. 80,443")
		fs.Parse(rest)
		requireArg("project", *project)
		requireArg("name", *name)
		requireArg("network", *network)
		requireArg("protocol", *protocol)
		requireArg("ports", *ports)

		resp, err := makeJSONRequest(
			"POST",
			gcpBase(*project)+"/global/firewalls",
			map[string]interface{}{
				"name":    *name,
				"network": *network,
				"allowed": []map[string]interface{}{
					{
						"IPProtocol": *protocol,
						"ports":      splitCSV(*ports),
					},
				},
			},
			nil,
		)
		if err != nil {
			fatal("%v", err)
		}
		printJSON(resp)

	default:
		fmt.Fprintf(os.Stderr, "error: unknown gcp action %q\n\n", action)
		gcpPrintUsage()
		os.Exit(1)
	}
}

// splitCSV splits a comma-separated string into a slice, trimming spaces.
func splitCSV(s string) []string {
	var parts []string
	for _, p := range strings.Split(s, ",") {
		if t := strings.TrimSpace(p); t != "" {
			parts = append(parts, t)
		}
	}
	return parts
}
