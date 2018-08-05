package deployer

import "strings"

type State struct {
	// Version is the state file protocol version.
	Version int `json:"version"`

	// TFVersion is the version of Terraform that wrote this state.
	TFVersion string `json:"terraform_version,omitempty"`

	// Serial is incremented on any operation that modifies
	// the State file. It is used to detect potentially conflicting
	// updates.
	Serial int64 `json:"serial"`

	// Modules contains all the modules in a breadth-first order
	Modules []ModuleState `json:"modules"`
}

type ModuleState struct {
	// Path is the import path from the root module. Modules imports are
	// always disjoint, so the path represents amodule tree
	Path []string `json:"path"`

	// Locals are kept only transiently in-memory, because we can always
	// re-compute them.
	Locals map[string]interface{} `json:"-"`

	// Outputs declared by the module and maintained for each module
	// even though only the root module technically needs to be kept.
	// This allows operators to inspect values at the boundaries.
	Outputs map[string]*OutputState `json:"outputs"`

	// Resources is a mapping of the logically named resource to
	// the state of the resource. Each resource may actually have
	// N instances underneath, although a user only needs to think
	// about the 1:1 case.
	Resources map[string]ResourceState `json:"resources"`

	// Dependencies are a list of things that this module relies on
	// existing to remain intact. For example: an module may depend
	// on a VPC ID given by an aws_vpc resource.
	//
	// Terraform uses this information to build valid destruction
	// orders and to warn the user if they're destroying a module that
	// another resource depends on.
	//
	// Things can be put into this list that may not be managed by
	// Terraform. If Terraform doesn't find a matching ID in the
	// overall state, then it assumes it isn't managed and doesn't
	// worry about it.
	Dependencies []string `json:"depends_on"`
}

type OutputState struct {
	// Sensitive describes whether the output is considered sensitive,
	// which may lead to masking the value on screen in some cases.
	Sensitive bool `json:"sensitive"`
	// Type describes the structure of Value. Valid values are "string",
	// "map" and "list"
	Type string `json:"type"`
	// Value contains the value of the output, in the structure described
	// by the Type field.
	Value interface{} `json:"value"`
}

type ResourceState struct {
	// This is filled in and managed by Terraform, and is the resource
	// type itself such as "mycloud_instance". If a resource provider sets
	// this value, it won't be persisted.
	Type string `json:"type"`

	// Dependencies are a list of things that this resource relies on
	// existing to remain intact. For example: an AWS instance might
	// depend on a subnet (which itself might depend on a VPC, and so
	// on).
	//
	// Terraform uses this information to build valid destruction
	// orders and to warn the user if they're destroying a resource that
	// another resource depends on.
	//
	// Things can be put into this list that may not be managed by
	// Terraform. If Terraform doesn't find a matching ID in the
	// overall state, then it assumes it isn't managed and doesn't
	// worry about it.
	Dependencies []string `json:"depends_on"`

	// Primary is the current active instance for this resource.
	// It can be replaced but only after a successful creation.
	// This is the instances on which providers will act.
	Primary InstanceState `json:"primary"`

	// Provider is used when a resource is connected to a provider with an alias.
	// If this string is empty, the resource is connected to the default provider,
	// e.g. "aws_instance" goes with the "aws" provider.
	// If the resource block contained a "provider" key, that value will be set here.
	Provider string `json:"provider"`
}

type InstanceState struct {
	// A unique ID for this resource. This is opaque to Terraform
	// and is only meant as a lookup mechanism for the providers.
	ID string `json:"id"`

	// Attributes are basic information about the resource. Any keys here
	// are accessible in variable format within Terraform configurations:
	// ${resourcetype.name.attribute}.
	Attributes map[string]interface{} `json:"attributes"`

	// Meta is a simple K/V map that is persisted to the State but otherwise
	// ignored by Terraform core. It's meant to be used for accounting by
	// external client code. The value here must only contain Go primitives
	// and collections.
	Meta map[string]interface{} `json:"meta"`

	// Tainted is used to mark a resource for recreation.
	Tainted bool `json:"tainted"`
}

type ConfigWrappers struct {
	EC2                    []EC2ConfigWrapper
	EC2ModuleCount         int
	DO                     []DOConfigWrapper
	DropletModuleCount     int
	AWSAPI                 []AWSApiConfigWrapper
	AWSAPIModuleCount      int
	Cloudfront             []CloudfrontConfigWrapper
	CloudfrontModuleCount  int
	Googlefront            []GooglefrontConfigWrapper
	GooglefrontModuleCount int
}

type ListStruct struct {
	IP         string
	Provider   string
	Region     string
	Name       string
	Place      int
	Username   string
	PrivateKey string
}

func (listStruct *ListStruct) String() string {
	return ("IP: " + listStruct.IP + " - Provider: " + listStruct.Provider + " - Region: " + listStruct.Region + " - Name: " + listStruct.Name)
}

type APIOutput struct {
	TargetURI string
	InvokeURI string
	Provider  string
	Name      string
}

func (output APIOutput) String() string {
	return " - Target: " + output.TargetURI + " - Invoke URI: " + output.InvokeURI + " - Provider: " + output.Provider + " - Name: " + output.Name
}

type DomainFrontOutput struct {
	Name         string
	Origin       string
	ID           string
	Invoke       string
	Provider     string
	Etag         string
	Status       string
	FunctionName string
	RestrictUA   string
}

func (output DomainFrontOutput) String() string {
	if output.Provider == "GOOGLE" {
		origin := strings.Replace(output.Origin, "_", ".", -1)
		restrictedUA := output.RestrictUA
		if output.RestrictUA == "" {
			restrictedUA = "N/A"
		}
		return " - Origin: " + origin + " - Invoke: " + output.Invoke + " - Status: " + output.Status + " - RestrictedUA: " + restrictedUA + " - Provider: " + output.Provider
	}
	return " - Origin: " + output.Origin + " - Invoke: " + output.Invoke + " - Status: " + output.Status + " - Provider: " + output.Provider
}
