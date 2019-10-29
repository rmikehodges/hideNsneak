package cmd

import (
	"bufio"
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var cfgFile string

var configTemplate = `
{
    "aws_access_id": "",
    "aws_secret_key": "",
    "aws_bucket_name": "",
    "digitalocean_token":"",
    "azure_tenant_id": "",
    "azure_client_id": "",
    "azure_client_secret": "",
    "azure_location": "",
    "azure_subscription_id": "",
    "google_credentials_path": "",
    "google_project_id": "",
    "public_key":"",
    "private_key":"",
    "do_user": "root",
    "ec2_user": "ubuntu",
    "google_user": "ubuntu"
}
`

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "hidensneak",
	Short: "hideNsneak is an application that will help you automate cloud management.",
	Long: `
	__     __     __         _______                              __    
	|  |--.|__|.--|  |.-----.|    |  |.-----..-----..-----..---.-.|  |--.
	|     ||  ||  _  ||  -__||       ||__ --||     ||  -__||  _  ||    < 
	|__|__||__||_____||_____||__|____||_____||__|__||_____||___._||__|__|
																		 

hideNsneak is a CLI that empowers red teamers during penetration testing.
This application is a tool that automates deployment, management, and destruction
of cloud infrastructure.`,
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)

	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.

	// config file is a parameter but its set to default here
	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "$GOPATH/src/github.com/rmikehodges/hideNsneak/config/config.json", "config file")
}

// initConfig expands the filepath for the user
func initConfig() {
	goPath := os.Getenv("GOPATH")
	cfgFile = goPath + "/src/github.com/rmikehodges/hideNsneak/config/config.json"

	if _, err := os.Stat(cfgFile); os.IsNotExist(err) {
		f, err := os.Create(cfgFile)
		defer f.Close()
		w := bufio.NewWriter(f)
		gg, err := w.WriteString(configTemplate)
	}

}
