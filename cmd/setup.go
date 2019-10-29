package cmd

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/rmikehodges/hideNsneak/deployer"
)

//TODO: Go says to remove underscores, would need to be updated throughout the file
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
		configStruct = deployer.RetrieveConfig(cfgFile)

		configStruct.aws_access_key = aws_access_key
		configStruct.aws_secret_key = aws_secret_key
		configStruct.aws_bucket_name = aws_bucket_name

		deployer.UpdateConfig(configStruct)

		deployer.InitializeBackendDDB(aws_access_key, aws_secret_key)
	},
}

var configStruct = &cobra.Command{
	Use:   "config",
	Short: "dump the config file",
	Long:  `Shows the contents of the config file`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println(deployer.RetrieveConfig(cfgFile))
	},
}

var do = &cobra.Command{
	Use:   "do",
	Short: "setup do variables",
	Long:  `Setup your do secrets`,
	Run: func(cmd *cobra.Command, args []string) {
		configStruct = deployer.RetrieveConfig(cfgFile)

		configStruct.digitalocean_token = digitalocean_token

		deployer.UpdateConfig(configStruct)
	},
}

var azure = &cobra.Command{
	Use:   "azure",
	Short: "setup azure variables",
	Long:  `Setup your azure secrets`,
	Run: func(cmd *cobra.Command, args []string) {
		configStruct = deployer.RetrieveConfig(cfgFile)

		configStruct.azure_tenant_id = azure_tenant_id
		configStruct.azure_client_id = azure_client_id
		configStruct.azure_client_secret = azure_client_secret
		configStruct.azure_location = azure_location
		configStruct.azure_subscription_id = azure_subscription_id

		deployer.UpdateConfig(configStruct)
	},
}

var ssh = &cobra.Command{
	Use:   "ssh",
	Short: "setup ssh key",
	Long:  `Setup your ssh key to control your infrastructure`,
	Run: func(cmd *cobra.Command, args []string) {
		configStruct = deployer.RetrieveConfig()

		configStruct.public_key = public_key
		configStruct.private_key = private_key

		deployer.UpdateConfig(configStruct)
	},
}

func init() {
	rootCmd.AddCommand(setup)
	setup.AddCommand(aws, azure, do)

	aws.PersistentFlags().StringVarP(&aws_secret_key, "secret", "s", "", "[Required] AWS secret key")
	aws.MarkPersistentFlagRequired("secret")
	aws.PersistentFlags().StringVarP(&aws_access_key, "access", "a", "", "[Required] AWS access key")
	aws.MarkPersistentFlagRequired("access")
	aws.PersistentFlags().StringVarP(&aws_bucket_name, "bucket", "b", "", "[Required] AWS bucket for state management")
	aws.MarkPersistentFlagRequired("bucket")

	azure.PersistentFlags().StringVarP(&azure_tenant_id, "tenant", "t", "", "[Required] azure tenant id")
	azure.PersistentFlags().StringVarP(&azure_client_id, "client", "c", "", "[Required] azure client id")
	azure.PersistentFlags().StringVarP(&azure_client_secret, "secret", "s", "", "[Required] azure client secret")
	azure.PersistentFlags().StringVarP(&azure_location, "location", "l", "", "[Required] azure location")
	azure.PersistentFlags().StringVarP(&azure_subscription_id, "sub", "i", "", "[Required] azure subscription id")

	ssh.PersistentFlags().StringVarP(&private_key, "private", "i", "", "[Required] ssh private key")
	ssh.MarkPersistentFlagRequired("private")
	ssh.PersistentFlags().StringVarP(&public_key, "public", "p", "", "[Required] ssh public key")
	ssh.MarkPersistentFlagRequired("public")

	do.PersistentFlags().StringVarP(&private_key, "token", "t", "", "[Required] do secret token")
	do.MarkPersistentFlagRequired("token")
}
