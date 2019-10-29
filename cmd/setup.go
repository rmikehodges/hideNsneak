package cmd

import (
	"fmt"
	"io/ioutil"

	"github.com/spf13/cobra"

	"github.com/rmikehodges/hideNsneak/deployer"
)

var awsSecretKey string
var awsAccessID string
var awsS3BucketName string
var digitaloceanToken string
var azureTenantID string
var azureClientID string
var azureClientSecret string
var azureLocation string
var azureSubscriptionID string
var googleCredentialsPath string
var googleProject string
var publicKey string
var privateKey string
var doUser string
var ec2User string
var googleUser string

var setup = &cobra.Command{
	Use:   "setup",
	Short: "setup parent command",
	Long:  `parent command for managing setup`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Run 'setup --help' for usage.")
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
		configStruct := deployer.RetrieveConfig(cfgFile)

		configStruct.AwsAccessID = awsAccessID
		configStruct.AwsSecretKey = awsSecretKey
		configStruct.AwsS3BucketName = awsS3BucketName

		deployer.UpdateConfig(cfgFile, configStruct)

		deployer.InitializeBackendDDB(awsAccessID, awsSecretKey)
	},
}

var config = &cobra.Command{
	Use:   "config",
	Short: "dump the config file",
	Long:  `Shows the contents of the config file`,
	Run: func(cmd *cobra.Command, args []string) {
		data, err := ioutil.ReadFile(cfgFile)
		if err != nil {
			fmt.Println(err)
		}
		fmt.Println(string(data))
	},
}

var do = &cobra.Command{
	Use:   "do",
	Short: "setup do variables",
	Long:  `Setup your do secrets`,
	Run: func(cmd *cobra.Command, args []string) {
		configStruct := deployer.RetrieveConfig(cfgFile)

		configStruct.DigitaloceanToken = digitaloceanToken

		deployer.UpdateConfig(cfgFile, configStruct)
	},
}

var azure = &cobra.Command{
	Use:   "azure",
	Short: "setup azure variables",
	Long:  `Setup your azure secrets`,
	Run: func(cmd *cobra.Command, args []string) {
		configStruct := deployer.RetrieveConfig(cfgFile)

		configStruct.AzureTenantID = azureTenantID
		configStruct.AzureClientID = azureClientID
		configStruct.AzureClientSecret = azureClientSecret
		configStruct.AzureLocation = azureLocation
		configStruct.AzureSubscriptionID = azureSubscriptionID

		deployer.UpdateConfig(cfgFile, configStruct)
	},
}

var ssh = &cobra.Command{
	Use:   "ssh",
	Short: "setup ssh key",
	Long:  `Setup your ssh key to control your infrastructure`,
	Run: func(cmd *cobra.Command, args []string) {
		configStruct := deployer.RetrieveConfig(cfgFile)

		configStruct.PublicKey = publicKey
		configStruct.PrivateKey = privateKey

		deployer.UpdateConfig(cfgFile, configStruct)
	},
}

func init() {
	rootCmd.AddCommand(setup)
	setup.AddCommand(aws, azure, do, config, ssh)

	aws.PersistentFlags().StringVarP(&awsSecretKey, "secret", "s", "", "[Required] AWS secret key")
	aws.MarkPersistentFlagRequired("secret")
	aws.PersistentFlags().StringVarP(&awsAccessID, "access", "a", "", "[Required] AWS access key")
	aws.MarkPersistentFlagRequired("access")
	aws.PersistentFlags().StringVarP(&awsS3BucketName, "bucket", "b", "", "[Required] AWS bucket for state management")
	aws.MarkPersistentFlagRequired("bucket")

	azure.PersistentFlags().StringVarP(&azureTenantID, "tenant", "t", "", "[Required] azure tenant id")
	azure.PersistentFlags().StringVarP(&azureClientID, "client", "c", "", "[Required] azure client id")
	azure.PersistentFlags().StringVarP(&azureClientSecret, "secret", "s", "", "[Required] azure client secret")
	azure.PersistentFlags().StringVarP(&azureLocation, "location", "l", "", "[Required] azure location")
	azure.PersistentFlags().StringVarP(&azureSubscriptionID, "sub", "i", "", "[Required] azure subscription id")

	ssh.PersistentFlags().StringVarP(&privateKey, "private", "i", "", "[Required] ssh private key")
	ssh.MarkPersistentFlagRequired("private")
	ssh.PersistentFlags().StringVarP(&publicKey, "public", "p", "", "[Required] ssh public key")
	ssh.MarkPersistentFlagRequired("public")

	do.PersistentFlags().StringVarP(&digitaloceanToken, "token", "t", "", "[Required] do secret token")
	do.MarkPersistentFlagRequired("token")
}
