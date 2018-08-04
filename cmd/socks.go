package cmd

import (
	"fmt"
	"strings"

	"github.com/rmikehodges/hideNsneak/deployer"

	"github.com/spf13/cobra"
)

var socksPort int
var socksInstanceInput string

// helloCmd represents the hello command
var socks = &cobra.Command{
	Use:   "socks",
	Short: "socks",
	Long:  `socks`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Run 'socks --help' for usage.")
	},
}

var socksDeploy = &cobra.Command{
	Use:   "deploy",
	Short: "Deploy SOCKS Proxy",
	Long:  `Deploy SOCKS Proxy`,
	Args: func(cmd *cobra.Command, args []string) error {
		err := deployer.IsValidNumberInput(socksInstanceInput)

		if err != nil {
			return err
		}

		expandedNumIndex := deployer.ExpandNumberInput(socksInstanceInput)

		err = deployer.ValidateNumberOfInstances(expandedNumIndex, "instance")

		if err != nil {
			return err
		}

		return err
	},
	Run: func(cmd *cobra.Command, args []string) {
		marshalledState := deployer.TerraformStateMarshaller()

		list := deployer.ListInstances(marshalledState)

		expandedNumIndex := deployer.ExpandNumberInput(socksInstanceInput)

		for _, num := range expandedNumIndex {
			err := deployer.CreateSingleSOCKS(list[num].PrivateKey, list[num].Username, list[num].IP, socksPort)
			if err != nil {
				fmt.Println("SOCKS creation failed for " + list[num].IP)
			}
			socksPort = socksPort + 1
		}

	},
}

var socksDestroy = &cobra.Command{
	Use:   "destroy",
	Short: "Destroy a SOCKS Proxy",
	Long:  `Destroy a SOCKS Proxy`,
	Args: func(cmd *cobra.Command, args []string) error {
		err := deployer.IsValidNumberInput(socksInstanceInput)

		if err != nil {
			return err
		}

		expandedNumIndex := deployer.ExpandNumberInput(socksInstanceInput)

		err = deployer.ValidateNumberOfInstances(expandedNumIndex, "instance")

		if err != nil {
			return err
		}

		return err
	},
	Run: func(cmd *cobra.Command, args []string) {
		marshalledState := deployer.TerraformStateMarshaller()

		list := deployer.ListInstances(marshalledState)

		expandedNumIndex := deployer.ExpandNumberInput(socksInstanceInput)

		for _, num := range expandedNumIndex {
			deployer.DestroySOCKS(list[num].IP)
		}

	},
}

var socksList = &cobra.Command{
	Use:   "list",
	Short: "List available SOCKS Proxies",
	Long:  `List available SOCKS Proxies`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Pulling Terraform State...")
		marshalledState := deployer.TerraformStateMarshaller()

		list := deployer.ListInstances(marshalledState)

		output := deployer.ListProxies(list)

		fmt.Println(output)
	},
}

var proxychains = &cobra.Command{
	Use:   "proxychains",
	Short: "Proxychains Config",
	Long:  `Prints out the proper proxychains configuration`,
	Run: func(cmd *cobra.Command, args []string) {
		marshalledState := deployer.TerraformStateMarshaller()

		list := deployer.ListInstances(marshalledState)

		output := deployer.ListProxies(list)

		fmt.Println(deployer.PrintProxyChains(output))
	},
}

var socksd = &cobra.Command{
	Use:   "socksd",
	Short: "SOCKSd config",
	Long:  `Prints out the proper socksd config`,
	Run: func(cmd *cobra.Command, args []string) {
		marshalledState := deployer.TerraformStateMarshaller()

		list := deployer.ListInstances(marshalledState)

		output := deployer.ListProxies(list)

		output = strings.TrimSpace(output)

		fmt.Println(deployer.PrintSocksd(output))
	},
}

func init() {
	rootCmd.AddCommand(socks)
	socks.AddCommand(socksDeploy, socksDestroy, socksList, proxychains, socksd)

	socksDeploy.PersistentFlags().IntVarP(&socksPort, "port", "p", 8081, "[Optional] port to start incrementing from for socks proxies")

	socksDeploy.PersistentFlags().StringVarP(&socksInstanceInput, "index", "i", "", "[Required] indices of the instances to deploy a socks proxy to")
	socksDeploy.MarkPersistentFlagRequired("index")

	socksDestroy.PersistentFlags().StringVarP(&socksInstanceInput, "index", "i", "", "[Required] indices of the instances to destroy the socks proxy for")
	socksDestroy.MarkPersistentFlagRequired("index")

}
