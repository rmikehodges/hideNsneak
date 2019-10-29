package cmd

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/rmikehodges/hideNsneak/deployer"
)

var aws_secret_key string
var aws_access_key string
var aws_bucket_name string
var digitalocean_token string
var azure_tenant_id string
var azure_client_id string
var azure_client_secret string
var azure_location string
var azure_subscription_id string
var google_credentials_path string
var google_project_id string
var public_key string
var private_key string
var do_user string
var ec2_user string
var google_user string

goPath := os.Getenv("GOPATH")
configPath = goPath + "/src/github.com/rmikehodges/hideNsneak/config/config.json"

var setup = &cobra.Command{
	Use:   "setup",
	Short: "instance parent command",
	Long:  `parent command for managing instances`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Run 'instance --help' for usage.")
	},
}

func init() {
	rootCmd.AddCommand(setup)
}

var aws = &cobra.Command{
	Use:   "aws",
	Short: "setup aws variables",
	Long:  `Setup your aws secrets as well as dynamo table`,
	Run: func(cmd *cobra.Command, args []string) {
		config = deployer.RetrieveConfig(configPath)

		config.aws_access_key = aws_access_key
		config.aws_secret_key = aws_secret_key

		deployer.UpdateConfig(configPath)

		deployer.InitializeBackendDDB(aws_access_key, aws_secret_key)
	},
}

var config = &cobra.Command{
	Use:   "config",
	Short: "dump the config file",
	Long:  `Shows the contents of the config file`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println(deployer.RetrieveConfig(configPath))
	},
}

var do = &cobra.Command{
	Use:   "do",
	Short: "setup do variables",
	Long:  `Setup your do secrets`,
	Run: func(cmd *cobra.Command, args []string) {
		config = deployer.RetrieveConfig(configPath)

		config.digitalocean_token = digitalocean_token

		deployer.UpdateConfig(configPath)
	},
}

var azure = &cobra.Command{
	Use:   "azure",
	Short: "setup azure variables",
	Long:  `Setup your azure secrets`,
	Run: func(cmd *cobra.Command, args []string) {
		config = deployer.RetrieveConfig(configPath)

		config.azure_tenant_id = azure_tenant_id
		config.azure_client_id = azure_client_id
		config.azure_client_secret = azure_client_secret
		config.azure_location = azure_location
		config.azure_subscription_id = azure_subscription_id

		deployer.UpdateConfig(configPath)
	},
}

var ssh = &cobra.Command{
	Use:   "ssh",
	Short: "setup ssh key",
	Long:  `Setup your ssh key to control your infrastructure`,
	Run: func(cmd *cobra.Command, args []string) {
		config = deployer.RetrieveConfig()

		config.public_key = public_key
		config.private_key = private_key

		deployer.UpdateConfig(configPath)
	},
}

func init() {
	rootCmd.AddCommand(setup)
	setup.AddCommand(aws, azure, do)

	aws.PersistentFlags().StringSliceVarP(&aws_secret_key, "secret", "s", nil, "[Required] AWS secret key")
	aws.MarkPersistentFlagRequired("secret")
	aws.PersistentFlags().StringSliceVarP(&aws_secret_key, "access", "a", nil, "[Required] AWS access key")
	aws.MarkPersistentFlagRequired("access")

	azure.PersistentFlags().StringSliceVarP(&azure_tenant_id, "tenant", "t", nil, "[Required] azure tenant id")
	azure.PersistentFlags().StringSliceVarP(&azure_client_id, "client", "c", nil, "[Required] azure client id")
	azure.PersistentFlags().StringSliceVarP(&azure_client_secret, "secret", "s", nil, "[Required] azure client secret")
	azure.PersistentFlags().StringSliceVarP(&azure_location, "location", "l", nil, "[Required] azure location")
	azure.PersistentFlags().StringSliceVarP(&azure_subscription_id, "sub", "i", nil, "[Required] azure subscription id")

	ssh.PersistentFlags().StringSliceVarP(&private_key, "private", "i", nil, "[Required] ssh private key")
	ssh.MarkPersistentFlagRequired("private")
	ssh.PersistentFlags().StringSliceVarP(&public_key, "public", "p", nil, "[Required] ssh public key")
	ssh.MarkPersistentFlagRequired("public")

	do.PersistentFlags().StringSliceVarP(&private_key, "token", "t", nil, "[Required] do secret token")
	do.MarkPersistentFlagRequired("token")
}
