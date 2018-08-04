package cmd

import (
	"fmt"

	"github.com/rmikehodges/hideNsneak/deployer"

	"github.com/spf13/cobra"
)

var ufwAction string
var ufwTCPPorts []string
var ufwUDPPorts []string
var ufwIndices string

// helloCmd represents the hello command
var firewall = &cobra.Command{
	Use:   "firewall",
	Short: "firewall",
	Long:  `firewall`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Run 'firewall --help' for usage.")
	},
}

var firewallAdd = &cobra.Command{
	Use:   "add",
	Short: "add a ufw firewall rule",
	Long:  `adds a ufw firewall rules to target host containing the tcp and udp port specifications set out by the user`,
	Args: func(cmd *cobra.Command, args []string) error {
		_, err := deployer.ValidatePorts(ufwTCPPorts)
		if err != nil {
			return err
		}
		_, err = deployer.ValidatePorts(ufwUDPPorts)
		if err != nil {
			return err
		}

		err = deployer.IsValidNumberInput(ufwIndices)
		if err != nil {
			return err
		}

		expandedNumIndex := deployer.ExpandNumberInput(ufwIndices)

		err = deployer.ValidateNumberOfInstances(expandedNumIndex, "instance")
		if err != nil {
			return err
		}

		return err
	},
	Run: func(cmd *cobra.Command, args []string) {
		ufwTCPPorts, _ := deployer.ValidatePorts(ufwTCPPorts)

		ufwUDPPorts, _ := deployer.ValidatePorts(ufwUDPPorts)

		apps := []string{"firewall"}
		playbook := deployer.GeneratePlaybookFile(apps)

		marshalledState := deployer.TerraformStateMarshaller()

		list := deployer.ListInstances(marshalledState)

		var instances []deployer.ListStruct

		expandedNumIndex := deployer.ExpandNumberInput(ufwIndices)

		for _, num := range expandedNumIndex {
			instances = append(instances, list[num])
		}

		ufwAction = "add"

		hostFile := deployer.GenerateHostFile(instances, fqdn, domain, burpFile, localFilePath, remoteFilePath,
			execCommand, socatPort, socatIP, nmapOutput, nmapCommands,
			cobaltStrikeLicense, cobaltStrikePassword, cobaltStrikeC2Path, cobaltStrikeFile, cobaltStrikeKillDate,
			ufwAction, ufwTCPPorts, ufwUDPPorts)

		deployer.WriteToFile("ansible/hosts.yml", hostFile)
		deployer.WriteToFile("ansible/main.yml", playbook)

		deployer.ExecAnsible("hosts.yml", "main.yml", "ansible")
	},
}

var firewallDelete = &cobra.Command{
	Use:   "delete",
	Short: "delete a ufw firewall rule",
	Long:  `adds a ufw firewall rules to target host containing the tcp and udp port specifications set out by the user`,
	Args: func(cmd *cobra.Command, args []string) error {
		_, err := deployer.ValidatePorts(ufwTCPPorts)
		if err != nil {
			return err
		}
		_, err = deployer.ValidatePorts(ufwUDPPorts)
		if err != nil {
			return err
		}

		err = deployer.IsValidNumberInput(ufwIndices)
		if err != nil {
			return err
		}

		expandedNumIndex := deployer.ExpandNumberInput(ufwIndices)

		err = deployer.ValidateNumberOfInstances(expandedNumIndex, "instance")
		if err != nil {
			return err
		}

		return err
	},
	Run: func(cmd *cobra.Command, args []string) {
		ufwTCPPorts, _ = deployer.ValidatePorts(ufwTCPPorts)

		ufwUDPPorts, _ := deployer.ValidatePorts(ufwUDPPorts)

		apps := []string{"firewall"}

		playbook := deployer.GeneratePlaybookFile(apps)

		marshalledState := deployer.TerraformStateMarshaller()

		list := deployer.ListInstances(marshalledState)

		var instances []deployer.ListStruct

		expandedNumIndex := deployer.ExpandNumberInput(ufwIndices)

		for _, num := range expandedNumIndex {
			instances = append(instances, list[num])
		}

		ufwAction = "delete"

		hostFile := deployer.GenerateHostFile(instances, fqdn, domain, burpFile, localFilePath, remoteFilePath,
			execCommand, socatPort, socatIP, nmapOutput, nmapCommands,
			cobaltStrikeLicense, cobaltStrikePassword, cobaltStrikeC2Path, cobaltStrikeFile, cobaltStrikeKillDate,
			ufwAction, ufwTCPPorts, ufwUDPPorts)

		deployer.WriteToFile("ansible/hosts.yml", hostFile)
		deployer.WriteToFile("ansible/main.yml", playbook)

		deployer.ExecAnsible("hosts.yml", "main.yml", "ansible")
	},
}

var firewallList = &cobra.Command{
	Use:   "list",
	Short: "list ufw firewall rules",
	Long:  `lists all of the ufw firewall rules on the specifiec host`,
	Args: func(cmd *cobra.Command, args []string) error {
		err := deployer.IsValidNumberInput(ufwIndices)
		if err != nil {
			return err
		}

		expandedNumIndex := deployer.ExpandNumberInput(ufwIndices)

		err = deployer.ValidateNumberOfInstances(expandedNumIndex, "instance")
		if err != nil {
			return err
		}

		return err
	},
	Run: func(cmd *cobra.Command, args []string) {
		apps := []string{"firewall"}

		playbook := deployer.GeneratePlaybookFile(apps)

		marshalledState := deployer.TerraformStateMarshaller()

		list := deployer.ListInstances(marshalledState)

		var instances []deployer.ListStruct

		expandedNumIndex := deployer.ExpandNumberInput(ufwIndices)

		for _, num := range expandedNumIndex {
			instances = append(instances, list[num])
		}

		ufwAction = "list"

		hostFile := deployer.GenerateHostFile(instances, fqdn, domain, burpFile, localFilePath, remoteFilePath,
			execCommand, socatPort, socatIP, nmapOutput, nmapCommands,
			cobaltStrikeLicense, cobaltStrikePassword, cobaltStrikeC2Path, cobaltStrikeFile, cobaltStrikeKillDate,
			ufwAction, ufwTCPPorts, ufwUDPPorts)

		deployer.WriteToFile("ansible/hosts.yml", hostFile)
		deployer.WriteToFile("ansible/main.yml", playbook)

		deployer.ExecAnsible("hosts.yml", "main.yml", "ansible")
	},
}

func init() {
	rootCmd.AddCommand(firewall)
	firewall.AddCommand(firewallAdd, firewallDelete, firewallList)

	firewallAdd.PersistentFlags().StringVarP(&ufwIndices, "id", "i", "", "[Required] the id for the remote server")
	firewallAdd.MarkFlagRequired("id")

	firewallAdd.PersistentFlags().StringSliceVarP(&ufwTCPPorts, "tcp", "t", []string{}, "[Optional] the tcp ports to add i.e. 22,23")
	firewallAdd.PersistentFlags().StringSliceVarP(&ufwUDPPorts, "udp", "u", []string{}, "[Optional] the udp ports to add i.e. 500,53")

	firewallDelete.PersistentFlags().StringVarP(&ufwIndices, "id", "i", "", "[Required] the id for the remote server")
	firewallDelete.MarkFlagRequired("id")

	firewallDelete.PersistentFlags().StringSliceVarP(&ufwTCPPorts, "tcp", "t", []string{}, "[Optional] the tcp ports to delete i.e. 22,23")
	firewallDelete.PersistentFlags().StringSliceVarP(&ufwUDPPorts, "udp", "u", []string{}, "[Optional] the udp ports to delete i.e. 500,53")

	firewallList.PersistentFlags().StringVarP(&ufwIndices, "id", "i", "", "[Required] the id for the remote server")
	firewallList.MarkFlagRequired("id")

}
