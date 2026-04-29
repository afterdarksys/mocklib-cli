package main

import (
	"flag"
	"fmt"
	"os"
)

func azurePrintUsage() {
	fmt.Print(`Usage: mocklib-cli azure <action> [flags]

Actions:
  list-rgs      List resource groups
  create-rg     Create a resource group
  create-vnet   Create a virtual network
  list-vnets    List virtual networks in a resource group
  create-nsg    Create a network security group
  create-vm     Create a virtual machine
  list-vms      List virtual machines in a resource group
  get-vm        Get a virtual machine
  start-vm      Start a virtual machine
  stop-vm       Stop (deallocate) a virtual machine
  delete-vm     Delete a virtual machine

Flags (all actions):
  --sub    Azure subscription ID (required)

  list-rgs:        (only --sub)

  create-rg:
    --name       Resource group name (required)
    --location   Azure region, e.g. eastus (required)

  create-vnet:
    --rg     Resource group name (required)
    --name   VNet name (required)
    --cidr   Address space CIDR (required)

  list-vnets:
    --rg     Resource group name (required)

  create-nsg:
    --rg     Resource group name (required)
    --name   NSG name (required)

  create-vm:
    --rg     Resource group name (required)
    --name   VM name (required)
    --size   VM size, e.g. Standard_B1s (required)

  list-vms / get-vm / start-vm / stop-vm / delete-vm:
    --rg     Resource group name (required)
    --name   VM name (get/start/stop/delete only, required)
`)
}

func azureBase(sub string) string {
	return "/azure/subscriptions/" + sub
}

func runAzure(args []string) {
	if len(args) == 0 || args[0] == "--help" || args[0] == "-h" {
		azurePrintUsage()
		os.Exit(0)
	}

	action := args[0]
	rest := args[1:]

	switch action {
	case "list-rgs":
		fs := flag.NewFlagSet("list-rgs", flag.ExitOnError)
		sub := fs.String("sub", "", "Azure subscription ID")
		fs.Parse(rest)
		requireArg("sub", *sub)

		resp, err := makeJSONRequest("GET", azureBase(*sub)+"/resourceGroups", nil, nil)
		if err != nil {
			fatal("%v", err)
		}
		printJSON(resp)

	case "create-rg":
		fs := flag.NewFlagSet("create-rg", flag.ExitOnError)
		sub := fs.String("sub", "", "Azure subscription ID")
		name := fs.String("name", "", "Resource group name")
		location := fs.String("location", "", "Azure region")
		fs.Parse(rest)
		requireArg("sub", *sub)
		requireArg("name", *name)
		requireArg("location", *location)

		resp, err := makeJSONRequest(
			"PUT",
			azureBase(*sub)+"/resourceGroups/"+*name,
			map[string]interface{}{"location": *location},
			nil,
		)
		if err != nil {
			fatal("%v", err)
		}
		printJSON(resp)

	case "create-vnet":
		fs := flag.NewFlagSet("create-vnet", flag.ExitOnError)
		sub := fs.String("sub", "", "Azure subscription ID")
		rg := fs.String("rg", "", "Resource group name")
		name := fs.String("name", "", "VNet name")
		cidr := fs.String("cidr", "", "Address space CIDR")
		fs.Parse(rest)
		requireArg("sub", *sub)
		requireArg("rg", *rg)
		requireArg("name", *name)
		requireArg("cidr", *cidr)

		resp, err := makeJSONRequest(
			"PUT",
			azureBase(*sub)+"/resourceGroups/"+*rg+"/providers/Microsoft.Network/virtualNetworks/"+*name,
			map[string]interface{}{
				"location": "eastus",
				"properties": map[string]interface{}{
					"addressSpace": map[string]interface{}{
						"addressPrefixes": []string{*cidr},
					},
				},
			},
			nil,
		)
		if err != nil {
			fatal("%v", err)
		}
		printJSON(resp)

	case "list-vnets":
		fs := flag.NewFlagSet("list-vnets", flag.ExitOnError)
		sub := fs.String("sub", "", "Azure subscription ID")
		rg := fs.String("rg", "", "Resource group name")
		fs.Parse(rest)
		requireArg("sub", *sub)
		requireArg("rg", *rg)

		resp, err := makeJSONRequest(
			"GET",
			azureBase(*sub)+"/resourceGroups/"+*rg+"/providers/Microsoft.Network/virtualNetworks",
			nil, nil,
		)
		if err != nil {
			fatal("%v", err)
		}
		printJSON(resp)

	case "create-nsg":
		fs := flag.NewFlagSet("create-nsg", flag.ExitOnError)
		sub := fs.String("sub", "", "Azure subscription ID")
		rg := fs.String("rg", "", "Resource group name")
		name := fs.String("name", "", "NSG name")
		fs.Parse(rest)
		requireArg("sub", *sub)
		requireArg("rg", *rg)
		requireArg("name", *name)

		resp, err := makeJSONRequest(
			"PUT",
			azureBase(*sub)+"/resourceGroups/"+*rg+"/providers/Microsoft.Network/networkSecurityGroups/"+*name,
			map[string]interface{}{"location": "eastus", "properties": map[string]interface{}{}},
			nil,
		)
		if err != nil {
			fatal("%v", err)
		}
		printJSON(resp)

	case "create-vm":
		fs := flag.NewFlagSet("create-vm", flag.ExitOnError)
		sub := fs.String("sub", "", "Azure subscription ID")
		rg := fs.String("rg", "", "Resource group name")
		name := fs.String("name", "", "VM name")
		size := fs.String("size", "", "VM size, e.g. Standard_B1s")
		fs.Parse(rest)
		requireArg("sub", *sub)
		requireArg("rg", *rg)
		requireArg("name", *name)
		requireArg("size", *size)

		resp, err := makeJSONRequest(
			"PUT",
			azureBase(*sub)+"/resourceGroups/"+*rg+"/providers/Microsoft.Compute/virtualMachines/"+*name,
			map[string]interface{}{
				"location": "eastus",
				"properties": map[string]interface{}{
					"hardwareProfile": map[string]string{"vmSize": *size},
				},
			},
			nil,
		)
		if err != nil {
			fatal("%v", err)
		}
		printJSON(resp)

	case "list-vms":
		fs := flag.NewFlagSet("list-vms", flag.ExitOnError)
		sub := fs.String("sub", "", "Azure subscription ID")
		rg := fs.String("rg", "", "Resource group name")
		fs.Parse(rest)
		requireArg("sub", *sub)
		requireArg("rg", *rg)

		resp, err := makeJSONRequest(
			"GET",
			azureBase(*sub)+"/resourceGroups/"+*rg+"/providers/Microsoft.Compute/virtualMachines",
			nil, nil,
		)
		if err != nil {
			fatal("%v", err)
		}
		printJSON(resp)

	case "get-vm":
		fs := flag.NewFlagSet("get-vm", flag.ExitOnError)
		sub := fs.String("sub", "", "Azure subscription ID")
		rg := fs.String("rg", "", "Resource group name")
		name := fs.String("name", "", "VM name")
		fs.Parse(rest)
		requireArg("sub", *sub)
		requireArg("rg", *rg)
		requireArg("name", *name)

		resp, err := makeJSONRequest(
			"GET",
			azureBase(*sub)+"/resourceGroups/"+*rg+"/providers/Microsoft.Compute/virtualMachines/"+*name,
			nil, nil,
		)
		if err != nil {
			fatal("%v", err)
		}
		printJSON(resp)

	case "start-vm":
		fs := flag.NewFlagSet("start-vm", flag.ExitOnError)
		sub := fs.String("sub", "", "Azure subscription ID")
		rg := fs.String("rg", "", "Resource group name")
		name := fs.String("name", "", "VM name")
		fs.Parse(rest)
		requireArg("sub", *sub)
		requireArg("rg", *rg)
		requireArg("name", *name)

		resp, err := makeJSONRequest(
			"POST",
			azureBase(*sub)+"/resourceGroups/"+*rg+"/providers/Microsoft.Compute/virtualMachines/"+*name+"/start",
			nil, nil,
		)
		if err != nil {
			fatal("%v", err)
		}
		if resp != nil {
			printJSON(resp)
		} else {
			fmt.Printf("Started VM: %s\n", *name)
		}

	case "stop-vm":
		fs := flag.NewFlagSet("stop-vm", flag.ExitOnError)
		sub := fs.String("sub", "", "Azure subscription ID")
		rg := fs.String("rg", "", "Resource group name")
		name := fs.String("name", "", "VM name")
		fs.Parse(rest)
		requireArg("sub", *sub)
		requireArg("rg", *rg)
		requireArg("name", *name)

		resp, err := makeJSONRequest(
			"POST",
			azureBase(*sub)+"/resourceGroups/"+*rg+"/providers/Microsoft.Compute/virtualMachines/"+*name+"/deallocate",
			nil, nil,
		)
		if err != nil {
			fatal("%v", err)
		}
		if resp != nil {
			printJSON(resp)
		} else {
			fmt.Printf("Stopped VM: %s\n", *name)
		}

	case "delete-vm":
		fs := flag.NewFlagSet("delete-vm", flag.ExitOnError)
		sub := fs.String("sub", "", "Azure subscription ID")
		rg := fs.String("rg", "", "Resource group name")
		name := fs.String("name", "", "VM name")
		fs.Parse(rest)
		requireArg("sub", *sub)
		requireArg("rg", *rg)
		requireArg("name", *name)

		resp, err := makeJSONRequest(
			"DELETE",
			azureBase(*sub)+"/resourceGroups/"+*rg+"/providers/Microsoft.Compute/virtualMachines/"+*name,
			nil, nil,
		)
		if err != nil {
			fatal("%v", err)
		}
		if resp != nil {
			printJSON(resp)
		} else {
			fmt.Printf("Deleted VM: %s\n", *name)
		}

	default:
		fmt.Fprintf(os.Stderr, "error: unknown azure action %q\n\n", action)
		azurePrintUsage()
		os.Exit(1)
	}
}
