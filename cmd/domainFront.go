package cmd

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/rmikehodges/hideNsneak/deployer"

	"github.com/spf13/cobra"
)

var domainFrontProvider string
var domainFrontIndex int
var domainFrontOrigin string
var restrictUA string
var functionName string
var frontedDomain string

// helloCmd represents the hello command
var domainFront = &cobra.Command{
	Use:   "domainfront",
	Short: "domainfront",
	Long:  `domainfront parent command`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Run 'domainfront --help' for usage.")
	},
}

var domainFrontDeploy = &cobra.Command{
	Use:   "deploy",
	Short: "deploys a domain front",
	Long:  `deploys a domain front with a specified cloud provider`,
	Args: func(cmd *cobra.Command, args []string) error {
		switch strings.ToUpper(domainFrontProvider) {
		case "AWS":
		case "GOOGLE":
			match, _ := regexp.MatchString("^[a-zA-Z][a-zA-Z0-9]+", functionName)

			if functionName == "" {
				return fmt.Errorf("Google Domain Fronts must have a function name (-n)")
			} else if !match {
				return fmt.Errorf("Google Domain Front function must begin with a letter and can only contain letters and numbers")
			}
		case "AZURE":
			return fmt.Errorf("Currently unsupported provider")
		default:
			return fmt.Errorf("Unknown provider")
		}
		return nil
	},
	Run: func(cmd *cobra.Command, args []string) {
		domainFrontProvider = strings.ToUpper(domainFrontProvider)
		marshalledState := deployer.TerraformStateMarshaller()
		wrappers := deployer.CreateWrappersFromState(marshalledState, cfgFile)
		wrappers = deployer.DomainFrontDeploy(domainFrontProvider, domainFrontOrigin,
			restrictUA, functionName, frontedDomain, wrappers)

		mainFile := deployer.CreateMasterFile(wrappers)

		deployer.CreateTerraformMain(mainFile, cfgFile)

		deployer.TerraformApply(cfgFile)
	},
}

var domainFrontDestroy = &cobra.Command{
	Use:   "destroy",
	Short: "destroy domain front",
	Long:  `destroys a domain front specified by an index`,
	Args: func(cmd *cobra.Command, args []string) error {
		marshalledState := deployer.TerraformStateMarshaller()
		list := deployer.ListDomainFronts(marshalledState)

		if domainFrontIndex > len(list)-1 {
			return fmt.Errorf("domainfront index not in bounds")
		}
		if list[domainFrontIndex].Provider == "AWS" {
			if list[domainFrontIndex].Status == "Enabled" {
				return fmt.Errorf("domainfront must be disabled to destroy")
			}
		}
		return nil
	},
	Run: func(cmd *cobra.Command, args []string) {
		marshalledState := deployer.TerraformStateMarshaller()
		list := deployer.ListDomainFronts(marshalledState)

		currentDomainfront := list[domainFrontIndex]

		if list[domainFrontIndex].Provider == "AWS" {
			fmt.Println(deployer.AWSCloudFrontDestroy(currentDomainfront, cfgFile))
		} else {
			deployer.TerraformDestroy([]string{list[domainFrontIndex].Name}, cfgFile)
		}
	},
}

var domainFrontDisable = &cobra.Command{
	Use:   "disable",
	Short: "disable domainfront",
	Long:  `disables an enabled Cloudfront domainfront`,
	Args: func(cmd *cobra.Command, args []string) error {
		marshalledState := deployer.TerraformStateMarshaller()
		list := deployer.ListDomainFronts(marshalledState)
		if domainFrontIndex > len(list)-1 {
			return fmt.Errorf("domainfront index not in bounds")
		}

		if list[domainFrontIndex].Provider == "AWS" {
			if list[domainFrontIndex].Status == "Disabled" {
				return fmt.Errorf("domainfront is already disabled")
			}
		}

		return nil
	},
	Run: func(cmd *cobra.Command, args []string) {
		marshalledState := deployer.TerraformStateMarshaller()
		list := deployer.ListDomainFronts(marshalledState)
		wrappers := deployer.CreateWrappersFromState(marshalledState, cfgFile)
		switch list[domainFrontIndex].Provider {
		case "AWS":
			for index, front := range wrappers.Cloudfront {
				if list[domainFrontIndex].ID == front.ID {
					wrappers.Cloudfront[index].Enabled = "false"
				}
			}
		case "AZURE":
			//TODO: Implement Azure CDN domain fronting
		case "GOOGLE":
			fmt.Println("Disabling Google Domain Fronts is not currently supported")
			fmt.Println("Exiting....")
			// for index, front := range wrappers.Googlefront {
			// 	if list[domainFrontIndex].Invoke == front.InvokeURI {
			// 		wrappers.Googlefront[index].Enabled = false
			// 	}
			// }
			return
		default:
		}
		mainFile := deployer.CreateMasterFile(wrappers)

		deployer.CreateTerraformMain(mainFile, cfgFile)

		deployer.TerraformApply(cfgFile)

	},
}

var domainFrontEnable = &cobra.Command{
	Use:   "enable",
	Short: "enable",
	Long:  `enables a disabled Cloudfront domain front`,
	Args: func(cmd *cobra.Command, args []string) error {
		marshalledState := deployer.TerraformStateMarshaller()
		list := deployer.ListDomainFronts(marshalledState)
		if domainFrontIndex > len(list)-1 {
			return fmt.Errorf("domainfront index not in bounds")
		}

		if list[domainFrontIndex].Provider == "AWS" {
			if list[domainFrontIndex].Status == "Enabled" {
				return fmt.Errorf("domainfront is already enabled")
			}
		}

		return nil
	},
	Run: func(cmd *cobra.Command, args []string) {
		marshalledState := deployer.TerraformStateMarshaller()
		list := deployer.ListDomainFronts(marshalledState)
		wrappers := deployer.CreateWrappersFromState(marshalledState, cfgFile)

		switch list[domainFrontIndex].Provider {
		case "AWS":
			for index, front := range wrappers.Cloudfront {
				if list[domainFrontIndex].ID == front.ID {
					wrappers.Cloudfront[index].Enabled = "true"
				}
			}
		case "AZURE":
		case "GOOGLE":
			fmt.Println("Enabling Google Domain Fronts is not currently supported")
			fmt.Println("Exiting....")
			// for index, front := range wrappers.Googlefront {
			// 	if list[domainFrontIndex].Invoke == front.InvokeURI {
			// 		wrappers.Googlefront[index].Enabled = false
			// 	}
			// }
			return
		default:
		}

		mainFile := deployer.CreateMasterFile(wrappers)

		deployer.CreateTerraformMain(mainFile, cfgFile)

		deployer.TerraformApply(cfgFile)
	},
}

var domainFrontList = &cobra.Command{
	Use:   "list",
	Short: "list domain fronts",
	Long:  `list all domain fronts and their index, origin domain, invoke url, provider, and terraform name`,
	Run: func(cmd *cobra.Command, args []string) {
		marshalledState := deployer.TerraformStateMarshaller()
		list := deployer.ListDomainFronts(marshalledState)
		for index, front := range list {
			fmt.Print(index)
			fmt.Println(front)
		}

		return
	},
}

func init() {
	rootCmd.AddCommand(domainFront)
	domainFront.AddCommand(domainFrontDeploy, domainFrontDestroy, domainFrontList, domainFrontEnable, domainFrontDisable)

	domainFrontDeploy.PersistentFlags().StringVarP(&domainFrontProvider, "provider", "p", "", "[Required] specify the provider to use. only GOOGLE and AWS are supported")
	domainFrontDeploy.MarkPersistentFlagRequired("provider")

	domainFrontDeploy.PersistentFlags().StringVarP(&domainFrontOrigin, "target", "t", "", "[Required] the target domain or IP address of your C2 server. For Google your C2 server must support HTTPS.")
	domainFrontDeploy.MarkPersistentFlagRequired("target")

	domainFrontDeploy.PersistentFlags().StringVarP(&frontedDomain, "frontedDomain", "d", "", "[Required for Google] the domain to front as and use in defensive measures i.e. inbox.google.com")
	domainFrontDeploy.PersistentFlags().StringVarP(&functionName, "name", "n", "", "[Required for Google] the function name of the Google cloud function i.e actionjacksonb")
	domainFrontDeploy.PersistentFlags().StringVar(&restrictUA, "restrictua", "", "[Optional for Google] the User Agent header to check for valid requests. only matched User Agents will be forwarded to the C2")

	domainFrontEnable.PersistentFlags().IntVarP(&domainFrontIndex, "id", "i", 0, "[Required] the indices of the domain front(s) i.e. 1-3,5")
	domainFrontEnable.MarkPersistentFlagRequired("id")

	domainFrontDisable.PersistentFlags().IntVarP(&domainFrontIndex, "id", "i", 0, "[Required] the indices of the domain front(s) i.e. 1-3,5")
	domainFrontDisable.MarkPersistentFlagRequired("id")

	domainFrontDestroy.PersistentFlags().IntVarP(&domainFrontIndex, "id", "i", 0, "[Required] the indices of the domain front(s) i.e. 1-3,5")
	domainFrontDestroy.MarkPersistentFlagRequired("id")

}
