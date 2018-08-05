package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/rmikehodges/hideNsneak/deployer"

	"github.com/spf13/cobra"
)

var installArgs string
var burpCmd string
var installIndex string
var fqdn string
var domain string
var burpFile string

var install = &cobra.Command{
	Use:   "install",
	Short: "install",
	Long:  `install parent command`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Run 'install --help' for usage.")
	},
}

var collaboratorInstall = &cobra.Command{
	Use:   "collaborator",
	Short: "installs burp collaborator server",
	Long:  `installs and starts a burp collaborator with the domain on the remote server`,
	Args: func(cmd *cobra.Command, args []string) error {
		err := deployer.IsValidNumberInput(installIndex)

		if err != nil {
			return err
		}

		expandedInstallIndex := deployer.ExpandNumberInput(installIndex)

		err = deployer.ValidateNumberOfInstances(expandedInstallIndex, "instance", cfgFile)

		if err != nil {
			return err
		}

		return err
	},
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("WARNING: Its best to obtain your wildcard letsencrypt certificate prior to installation")
		fmt.Println("Do you still wish to continue?")
		if !deployer.AskForConfirmation() {
			return
		}

		apps := []string{"collaborator"}

		playbook := deployer.GeneratePlaybookFile(apps)

		marshalledState := deployer.TerraformStateMarshaller()

		list := deployer.ListInstances(marshalledState, cfgFile)

		var instances []deployer.ListStruct

		expandedNumIndex := deployer.ExpandNumberInput(installIndex)

		for _, num := range expandedNumIndex {
			instances = append(instances, list[num])
		}

		fqdn = domain

		hostFile := deployer.GenerateHostFile(instances, fqdn, domain, burpFile, localFilePath, remoteFilePath,
			execCommand, socatPort, socatIP, nmapOutput, nmapCommands,
			cobaltStrikeLicense, cobaltStrikePassword, cobaltStrikeC2Path, cobaltStrikeFile, cobaltStrikeKillDate,
			ufwAction, ufwTCPPorts, ufwUDPPorts)

		deployer.WriteToFile("ansible/hosts.yml", hostFile)
		deployer.WriteToFile("ansible/main.yml", playbook)

		deployer.ExecAnsible("hosts.yml", "main.yml")

		fmt.Println("Next Steps:")

		fmt.Println("1. Set this IP address to be both the primary and secondary nameserver for your domain")

		fmt.Println("Note: In order to have valid HTTPS on the collaborator server you must obtain a wildcard certificate from letsencrypt")
	},
}

var cobaltStrikeInstall = &cobra.Command{
	Use:   "cobaltstrike",
	Short: "installs Cobalt Strike",
	Long:  `installs, starts, and optionally licenses Cobaltstrike on the remote server with the specified malleable C2 profile and password`,
	Args: func(cmd *cobra.Command, args []string) error {
		if _, err := os.Stat(cobaltStrikeFile); os.IsNotExist(err) {
			return fmt.Errorf("cobaltstrike file does not exist")
		}

		if len(strings.Split(filepath.Base(cobaltStrikeFile), ".")) != 2 {
			return fmt.Errorf("cobaltstrike file must be in tgz format as only linux teamservers are supported")
		}

		if strings.Split(filepath.Base(cobaltStrikeFile), ".")[1] != "tgz" {
			return fmt.Errorf("cobaltstrike file must be in tgz format as only linux teamservers are supported")
		}

		err := deployer.IsValidNumberInput(installIndex)

		if err != nil {
			return err
		}

		expandedNumIndex := deployer.ExpandNumberInput(installIndex)

		err = deployer.ValidateNumberOfInstances(expandedNumIndex, "instance", cfgFile)

		if err != nil {
			return err
		}

		return err
	},
	Run: func(cmd *cobra.Command, args []string) {
		apps := []string{"cobaltstrike"}

		playbook := deployer.GeneratePlaybookFile(apps)

		marshalledState := deployer.TerraformStateMarshaller()

		list := deployer.ListInstances(marshalledState, cfgFile)

		var instances []deployer.ListStruct

		expandedNumIndex := deployer.ExpandNumberInput(installIndex)

		for _, num := range expandedNumIndex {
			instances = append(instances, list[num])
		}

		hostFile := deployer.GenerateHostFile(instances, fqdn, domain, burpFile, localFilePath, remoteFilePath,
			execCommand, socatPort, socatIP, nmapOutput, nmapCommands,
			cobaltStrikeLicense, cobaltStrikePassword, cobaltStrikeC2Path, cobaltStrikeFile, cobaltStrikeKillDate,
			ufwAction, ufwTCPPorts, ufwUDPPorts)

		deployer.WriteToFile("ansible/hosts.yml", hostFile)
		deployer.WriteToFile("ansible/main.yml", playbook)

		deployer.ExecAnsible("hosts.yml", "main.yml")
	},
}

var goPhishInstall = &cobra.Command{
	Use:   "gophish",
	Short: "installs Gophish",
	Long:  `installs gophish on the remote server`,
	Args: func(cmd *cobra.Command, args []string) error {
		err := deployer.IsValidNumberInput(installIndex)

		if err != nil {
			return err
		}

		expandedNumIndex := deployer.ExpandNumberInput(installIndex)

		err = deployer.ValidateNumberOfInstances(expandedNumIndex, "instance", cfgFile)

		if err != nil {
			return err
		}

		return err
	},
	Run: func(cmd *cobra.Command, args []string) {
		apps := []string{"gophish"}

		playbook := deployer.GeneratePlaybookFile(apps)

		marshalledState := deployer.TerraformStateMarshaller()

		list := deployer.ListInstances(marshalledState, cfgFile)

		var instances []deployer.ListStruct

		expandedNumIndex := deployer.ExpandNumberInput(installIndex)

		for _, num := range expandedNumIndex {
			instances = append(instances, list[num])
		}

		hostFile := deployer.GenerateHostFile(instances, fqdn, domain, burpFile, localFilePath, remoteFilePath,
			execCommand, socatPort, socatIP, nmapOutput, nmapCommands,
			cobaltStrikeLicense, cobaltStrikePassword, cobaltStrikeC2Path, cobaltStrikeFile, cobaltStrikeKillDate,
			ufwAction, ufwTCPPorts, ufwUDPPorts)

		deployer.WriteToFile("ansible/hosts.yml", hostFile)
		deployer.WriteToFile("ansible/main.yml", playbook)

		deployer.ExecAnsible("hosts.yml", "main.yml")
	},
}

var letsEncryptInstall = &cobra.Command{
	Use:   "letsencrypt",
	Short: "Installs Letsencrypt",
	Long:  `Installs Letsencrypt with the specified domain on the specified server`,
	Args: func(cmd *cobra.Command, args []string) error {
		err := deployer.IsValidNumberInput(installIndex)

		if err != nil {
			return err
		}

		expandedNumIndex := deployer.ExpandNumberInput(installIndex)

		err = deployer.ValidateNumberOfInstances(expandedNumIndex, "instance", cfgFile)

		if err != nil {
			return err
		}

		return err
	},
	Run: func(cmd *cobra.Command, args []string) {
		apps := []string{"letsencrypt"}

		playbook := deployer.GeneratePlaybookFile(apps)

		marshalledState := deployer.TerraformStateMarshaller()

		list := deployer.ListInstances(marshalledState, cfgFile)

		var instances []deployer.ListStruct

		expandedNumIndex := deployer.ExpandNumberInput(installIndex)

		for _, num := range expandedNumIndex {
			instances = append(instances, list[num])
		}

		hostFile := deployer.GenerateHostFile(instances, fqdn, domain, burpFile, localFilePath, remoteFilePath,
			execCommand, socatPort, socatIP, nmapOutput, nmapCommands,
			cobaltStrikeLicense, cobaltStrikePassword, cobaltStrikeC2Path, cobaltStrikeFile, cobaltStrikeKillDate,
			ufwAction, ufwTCPPorts, ufwUDPPorts)

		deployer.WriteToFile("ansible/hosts.yml", hostFile)
		deployer.WriteToFile("ansible/main.yml", playbook)

		deployer.ExecAnsible("hosts.yml", "main.yml")
	},
}

var nmapInstall = &cobra.Command{
	Use:   "nmap",
	Short: "installs nmap",
	Long:  `installs nmap on the remote server`,
	Args: func(cmd *cobra.Command, args []string) error {
		err := deployer.IsValidNumberInput(installIndex)

		if err != nil {
			return err
		}

		expandedNumIndex := deployer.ExpandNumberInput(installIndex)

		err = deployer.ValidateNumberOfInstances(expandedNumIndex, "instance", cfgFile)

		if err != nil {
			return err
		}

		return err
	},
	Run: func(cmd *cobra.Command, args []string) {
		apps := []string{"nmap"}

		playbook := deployer.GeneratePlaybookFile(apps)

		marshalledState := deployer.TerraformStateMarshaller()

		list := deployer.ListInstances(marshalledState, cfgFile)

		var instances []deployer.ListStruct

		expandedNumIndex := deployer.ExpandNumberInput(installIndex)

		for _, num := range expandedNumIndex {
			instances = append(instances, list[num])
		}

		hostFile := deployer.GenerateHostFile(instances, fqdn, domain, burpFile, localFilePath, remoteFilePath,
			execCommand, socatPort, socatIP, nmapOutput, nmapCommands,
			cobaltStrikeLicense, cobaltStrikePassword, cobaltStrikeC2Path, cobaltStrikeFile, cobaltStrikeKillDate,
			ufwAction, ufwTCPPorts, ufwUDPPorts)

		deployer.WriteToFile("ansible/hosts.yml", hostFile)
		deployer.WriteToFile("ansible/main.yml", playbook)

		deployer.ExecAnsible("hosts.yml", "main.yml")
	},
}

var socatInstall = &cobra.Command{
	Use:   "socat",
	Short: "Installs Socat",
	Long:  `Installs Socat to remote server`,
	Args: func(cmd *cobra.Command, args []string) error {
		err := deployer.IsValidNumberInput(installIndex)

		if err != nil {
			return err
		}

		expandedNumIndex := deployer.ExpandNumberInput(installIndex)

		err = deployer.ValidateNumberOfInstances(expandedNumIndex, "instance", cfgFile)

		if err != nil {
			return err
		}

		return err
	},
	Run: func(cmd *cobra.Command, args []string) {
		apps := []string{"socat"}

		playbook := deployer.GeneratePlaybookFile(apps)

		marshalledState := deployer.TerraformStateMarshaller()

		list := deployer.ListInstances(marshalledState, cfgFile)

		var instances []deployer.ListStruct

		expandedNumIndex := deployer.ExpandNumberInput(installIndex)

		for _, num := range expandedNumIndex {
			instances = append(instances, list[num])
		}

		hostFile := deployer.GenerateHostFile(instances, fqdn, domain, burpFile, localFilePath, remoteFilePath,
			execCommand, socatPort, socatIP, nmapOutput, nmapCommands,
			cobaltStrikeLicense, cobaltStrikePassword, cobaltStrikeC2Path, cobaltStrikeFile, cobaltStrikeKillDate,
			ufwAction, ufwTCPPorts, ufwUDPPorts)

		deployer.WriteToFile("ansible/hosts.yml", hostFile)
		deployer.WriteToFile("ansible/main.yml", playbook)

		deployer.ExecAnsible("hosts.yml", "main.yml")
	},
}

var sqlMapInstall = &cobra.Command{
	Use:   "sqlmap",
	Short: "Installs SQLmap",
	Long:  `Installs SQLmap to remote server`,
	Args: func(cmd *cobra.Command, args []string) error {
		err := deployer.IsValidNumberInput(installIndex)

		if err != nil {
			return err
		}

		expandedNumIndex := deployer.ExpandNumberInput(installIndex)

		err = deployer.ValidateNumberOfInstances(expandedNumIndex, "instance", cfgFile)

		if err != nil {
			return err
		}

		return err
	},
	Run: func(cmd *cobra.Command, args []string) {
		apps := []string{"sqlmap"}

		playbook := deployer.GeneratePlaybookFile(apps)

		marshalledState := deployer.TerraformStateMarshaller()

		list := deployer.ListInstances(marshalledState, cfgFile)

		var instances []deployer.ListStruct

		expandedNumIndex := deployer.ExpandNumberInput(installIndex)

		for _, num := range expandedNumIndex {
			instances = append(instances, list[num])
		}

		hostFile := deployer.GenerateHostFile(instances, fqdn, domain, burpFile, localFilePath, remoteFilePath,
			execCommand, socatPort, socatIP, nmapOutput, nmapCommands,
			cobaltStrikeLicense, cobaltStrikePassword, cobaltStrikeC2Path, cobaltStrikeFile, cobaltStrikeKillDate,
			ufwAction, ufwTCPPorts, ufwUDPPorts)

		deployer.WriteToFile("ansible/hosts.yml", hostFile)
		deployer.WriteToFile("ansible/main.yml", playbook)

		deployer.ExecAnsible("hosts.yml", "main.yml")
	},
}

// var empireInstall = &cobra.Command{
// 	Use:   "empire",
// 	Short: "Installs Powershell Empire",
// 	Long:  `Installs Powershell Empire to remote server`,
// 	Args: func(cmd *cobra.Command, args []string) {
// err := deployer.IsValidNumberInput(installIndex)

// if err != nil {
// 	return err
// }

// expandedNumIndex := deployer.ExpandNumberInput(installIndex)

// err = deployer.ValidateNumberOfInstances(expandedNumIndex)

// if err != nil {
// 	return err
// }

// return err
// 	},
// 	Run: func(cmd *cobra.Command, args []string) {
// 		apps := []string{"empire"}

// 		playbook := deployer.GeneratePlaybookFile(apps)

// 		masrshalledState := deployer.TerraformStateMarshaller()

// 		list := deployer.ListInstances(marshalledState, cfgFile)
// 		var instances []deployer.ListStruct

// expandedNumIndex := deployer.ExpandNumberInput(installIndex)

// for _, num := range expandedNumIndex {
// 	instances = append(instances, list[num])
// }

// 		hostFile := deployer.GenerateHostFile(instances, fqdn, domain, burpFile, localFilePath, remoteFilePath,
// 			execCommand, socatPort, socatIP, nmapOutput, nmapCommands,
// 			cobaltStrikeLicense, cobaltStrikePassword, cobaltStrikeC2Path, cobaltStrikeFile, cobaltStrikeKillDate,
// 			ufwAction, ufwTCPPorts, ufwUDPPorts)

// 		deployer.WriteToFile("ansible/hosts.yml", hostFile)
// 		deployer.WriteToFile("ansible/main.yml", playbook)

// 		deployer.ExecAnsible("hosts.yml", "main.yml")
// 	},
// }

func init() {
	rootCmd.AddCommand(install)
	install.AddCommand(collaboratorInstall, cobaltStrikeInstall, goPhishInstall, letsEncryptInstall, nmapInstall, socatInstall, sqlMapInstall /*, empireInstall*/)

	collaboratorInstall.PersistentFlags().StringVarP(&installIndex, "id", "i", "", "[Required] the id for the install")
	collaboratorInstall.MarkPersistentFlagRequired("id")
	collaboratorInstall.PersistentFlags().StringVarP(&domain, "domain", "d", "", "[Required] domain the collaborator instance will use")
	collaboratorInstall.MarkPersistentFlagRequired("domain")
	collaboratorInstall.PersistentFlags().StringVarP(&burpFile, "burpFile", "b", "", "[Required] the local filepath to the burp pro jar file")
	collaboratorInstall.MarkPersistentFlagRequired("burpFile")

	cobaltStrikeInstall.PersistentFlags().StringVarP(&installIndex, "id", "i", "", "[Required] the id for the install")
	cobaltStrikeInstall.MarkFlagRequired("id")
	cobaltStrikeInstall.PersistentFlags().StringVarP(&domain, "domain", "d", "", "[Not-In-Use] the domain for the teamserver. this functionality is not currently in use")
	cobaltStrikeInstall.PersistentFlags().StringVarP(&cobaltStrikeFile, "file", "f", "", "[Required] local filepath of the cobaltstrike tgz file")
	cobaltStrikeInstall.MarkPersistentFlagRequired("file")

	goPhishInstall.PersistentFlags().StringVarP(&installIndex, "id", "i", "", "[Required] the id for the install")
	goPhishInstall.MarkFlagRequired("id")
	goPhishInstall.PersistentFlags().StringVarP(&domain, "domain", "d", "", "[Optional] the domain for the gophish server")

	letsEncryptInstall.PersistentFlags().StringVarP(&installIndex, "id", "i", "", "[Required] the id for the install")
	letsEncryptInstall.MarkFlagRequired("id")
	//TODO Check this and how it applies to letsencrypt
	letsEncryptInstall.PersistentFlags().StringVarP(&fqdn, "fqdn", "f", "", "[Required] the fqdn of the server to generate a certificate for")
	letsEncryptInstall.MarkPersistentFlagRequired("fqdn")
	letsEncryptInstall.PersistentFlags().StringVarP(&domain, "domain", "d", "", "[Required] the domain of the server to generate a certificate for")
	letsEncryptInstall.MarkPersistentFlagRequired("domain")

	nmapInstall.PersistentFlags().StringVarP(&installIndex, "id", "i", "", "[Required] the id for the install")
	nmapInstall.MarkFlagRequired("id")

	socatInstall.PersistentFlags().StringVarP(&installIndex, "id", "i", "", "[Required] the id for the install")
	socatInstall.MarkFlagRequired("id")

	sqlMapInstall.PersistentFlags().StringVarP(&installIndex, "id", "i", "", "[Required] the id for the install")
	sqlMapInstall.MarkFlagRequired("id")

	// empireInstall.PersistentFlags().IntSliceVarP(&installIndex, "id", "i", []int{}, "Specify the id for the install")
	// empireInstall.MarkFlagRequired("id")
}
