package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var cfgFile string

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
	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "$GOPATH/src/github.com/rmikehodges/hideNsneak/config/config.json", "config file")
}

// initConfig expands the filepath for the user
func initConfig() {
	if cfgFile == "$GOPATH/src/github.com/rmikehodges/hideNsneak/config/config.json" {
		goPath := os.Getenv("GOPATH")
		cfgFile = goPath + "/src/github.com/rmikehodges/hideNsneak/config/config.json"
	}
}
