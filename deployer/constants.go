package deployer

var filepath = getTerraformDirectory()
var tfMainFile = filepath + "/main.tf"
var tfVariablesFile = filepath + "/variables.tf"
var tfVarsFile = filepath + "/terraform.tfvars"

const backend = `terraform {
	backend "s3" {
		bucket           = "{{.AwsS3BucketName}}"
		dynamodb_table   = "terraform-state-lock-dynamo"
		key              = "filename.tfstate"
		region           = "us-east-1"
		encrypt 	     = true 
	}
}
`

const templateSecrets = `
do_token = "{{.DigitaloceanToken}}"

aws_access_key = "{{.AwsAccessID}}"

aws_secret_key = "{{.AwsSecretKey}}"

azure_tenant_id = "{{.AzureTenantID}}"

azure_client_id = "{{.AzureClientID}}"

azure_client_secret = "{{.AzureClientSecret}}"

azure_subscription_id = "{{.AzureSubscriptionID}}"

google_credentials_path = "{{.GoogleCredentialsPath}}"

google_project = "{{.GoogleProject}}"
`

const variables = `
variable "do_token" {}

variable "aws_access_key" {}

variable "aws_secret_key" {}

variable "azure_tenant_id" {}

variable "azure_client_id" {}

variable "azure_client_secret" {}

variable "azure_subscription_id" {}

variable "google_credentials_path" {}

variable "google_project" {}
`

const outputs = `output "providers" {
	value = "${map(
	  "AWS", map(
		"instances", concat({{.ModuleNames}}),
		"security_group", list(map()), 
		"api", list(map()),
		"domain_front", list(map())),
	  "DO", map(
		"instances", list(map()),
		"firewalls", list(map())),
	  "GOOGLE", map(
		"instances", list(map())),
	  "AZURE", map(
		"instances", list(map())))
	  }"
  }`

///////////////////// MODULES /////////////////////

const mainEc2Module = `
	module "{{.ModuleName}}" {
	source          = "modules/ec2-deployment"

	region_count         = "${map({{$c := counter}}{{range $key, $value := .RegionMap}}{{if call $c}}, {{end}}"{{$key}}",{{$value}}{{end}})}"
	aws_instance_type    = "{{.InstanceType}}"
	ec2_default_user     = "{{.DefaultUser}}"
	aws_access_key       = "${var.aws_access_key}"
	aws_secret_key       = "${var.aws_secret_key}"
	aws_keypair_name     = "{{.KeyPairName}}"
	aws_private_key_file = "{{.PrivateKey}}"
	aws_public_key_file  = "{{.PublicKey}}"
  }
`

const mainAWSAPIModule = `module "{{.ModuleName}}" {
	source = "modules/aws-api-gateway"
  
	aws_access_key = "${var.aws_access_key}"
	aws_secret_key = "${var.aws_secret_key}"
  
	aws_api_target_uri = "{{.TargetURI}}"
  }
  
`

const mainCloudfrontModule = `module "{{.ModuleName}}" {
	source = "modules/cloudfront-deployment"
  
	aws_access_key = "${var.aws_access_key}"
	aws_secret_key = "${var.aws_secret_key}"
  
	cloudfront_origin = "{{.Origin}}"
	cloudfront_enabled = {{.Enabled}}
  }`

const mainDropletModule = `
  module "{{.ModuleName}}" {
	  source              = "modules/droplet-deployment"
	  do_region_count     = "${map({{$c := counter}}{{range $key, $value := .RegionMap}}{{if call $c}}, {{end}}"{{$key}}",{{$value}}{{end}})}"
	  do_token            = "${var.do_token}"
	  do_image            = "{{.Image}}"
	  do_private_key      = "{{.PrivateKey}}"
	  do_ssh_fingerprint  = "{{.Fingerprint}}"
	  do_size             = "{{.Size}}"
	  do_default_user     = "{{.DefaultUser}}"
  }
`

const azureCdnModule = `
	module "azure-cdn-{{.Endpoint}}" {
		source                  = "modules/azure-cdn-deployment"
		azure_subscription_id   = "${var.azure_subscription_id}"
		azure_tenant_id         = "${var.azure_tenant_id}"
		azure_client_id         = "${var.azure_client_id}"
		azure_client_secret     = "${var.azure_client_secret}"
		azure_cdn_hostname      = "{{.HostName}}"
		azure_cdn_profile_name  = "{{.ProfileName}}"
		azure_cdn_endpoint_name = "{{.EndpointName}}"
		azure_location          = "{{.Location}}"
	}
`

//TODO: need to run removeSpaces() on region when this is set
const azureModule = `
	module "azure-{{.Location}}" {
		source                = "modules/azure-deployment"
		azure_subscription_id = "${var.azure_subscription_id}"
		azure_tenant_id       = "${var.azure_tenant_id}"
		azure_client_id       = "${var.azure_client_id}"
		azure_client_secret   = "${var.azure_client_secret}"
		azure_location        = "{{.Location}}"
		azure_instance_count  = {{.Count}}
		azure_vm_size 		  = "{{.VMSize}}"
		azure_environment 	  = "{{.Environment}}"
		azure_public_key_file = "{{.PublicKey}}
		azure_private_key_file = "{{.PrivateKey}}
		ansible_groups       = "[]"
	}
`

const cloudfrontModule = `
	module "cloudfront-{{.Region}}" {
		source            = "modules/cloudfront-deployment"
		cloudfront_origin = "{{.Origin}}"
		aws_access_key    = "${var.aws_access_key}"
		aws_secret_key    = "${var.aws_secret_key}"
		aws_region 		  = "{{.Region}}"
	}
`

const googleCloudModule = `
	module "google-cloud-{{.Region}}" {
		source               	 = "modules/gcp-deployment"
		gcp_region          	 = "{{.Region}}"
		gcp_project          	 = "{{.Project}}"
		gcp_instance_count   	 = {{.Count}}
		gcp_ssh_user         	 = "{{.SSHUser}}"
		gcp_ssh_pub_key_file 	 = "{{.SSHPubKeyFile}}"
		gcp_ssh_private_key_file = "{{.SSHPrivateKeyFile}}"
		gcp_machine_type	 	 = "{{.MachineType}}"
		gcp_image			 	 = "{{.Image}}"
		ansible_groups       = "[]"
	}
`

const apiGatewayModule = `
	module "apigateway-{{.TargetUri}}" {
		source 				 = "modules/api-gateway"
		aws_access_key    	 = "${var.aws_access_key}"
		aws_secret_key    	 = "${var.aws_secret_key}"
		aws_api_target_uri 	 = "{{.TargetURI}"
		aws_api_stage_name	 = "{{.StageName}}"
  	}
`

const googlefrontModule = `
module "{{.ModuleName}}" {
  source = "modules/gcf-deployment"

  package_file = "{{.PackageFile}}"

  redirector_file = "{{.SourceFile}}"

  function_name = "{{.FunctionName}}"

  region = "us-central1"

  gcp_project = "${var.google_project}"

  enabled = {{.Enabled}}

  target = "{{.Host}}"

  restrictua = "{{.RestrictUA}}"

  google_credentials_path = "${var.google_credentials_path}"
}`

const googleDomainFrontCode = `
let httpProxy = require('http-proxy'),
    ip = require('ip');

let proxy = httpProxy.createProxyServer({secure: true});

let host = "{{.Host}}"
let target = "{{.HostURL}}"


let frontedDomain = "https://{{.FrontedDomain}}"
let restrictUA = "{{.RestrictUA}}"
let restrictSubnet = "{{.RestrictSubnet}}"
let restrictHeader = "{{.RestrictHeader}}"
let restrictValue = "{{.RestrictHeaderValue}}"


exports.redirector = (req, res) => {
  	let requestIP = req.ip
 

    if (req.method == "GET" || req.method == "POST") {
        if (restrictUA != "" && restrictUA != req.get('User-Agent')) {
            res.redirect(frontedDomain)
            return
        }
        req.host = host
        proxy.web(req, res, { target: target });
    } else {
        res.redirect(frontedDomain)
    }
};`

const googlefrontPackage = `{
	"name": "sample-http",
	"version": "0.0.1",
	"dependencies": {
		 "http-proxy": "1.17.0", 
	  "ip": "1.1.5"
	}
  }`
