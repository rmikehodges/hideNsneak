package cmd

import (
	"fmt"
	"regexp"

	"github.com/rmikehodges/hideNsneak/deployer"

	"github.com/spf13/cobra"
)

var nmapPorts []string
var nmapHostFile string
var nmapCommand string
var nmapOutput string
var nmapEvasive bool
var nmapCommands map[int][]string
var execCommand string
var socatPort string
var socatIP string
var cobaltStrikeLicense string
var cobaltStrikeFile string
var cobaltStrikePassword string
var cobaltStrikeC2Path string
var cobaltStrikeKillDate string

var commandIndices []int

var exec = &cobra.Command{
	Use:   "exec",
	Short: "execute",
	Long:  `exec parent command`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Run 'exec --help' for usage.")
	},
}

var command = &cobra.Command{
	Use:   "command",
	Short: "execute shell command",
	Long:  `executes the specified shell command on the specifed remote server`,
	Run: func(cmd *cobra.Command, args []string) {
		apps := []string{"exec"}

		playbook := deployer.GeneratePlaybookFile(apps)

		marshalledState := deployer.TerraformStateMarshaller()

		list := deployer.ListInstances(marshalledState)

		var instances []deployer.ListStruct

		for _, num := range commandIndices {
			instances = append(instances, list[num])
		}

		hostFile := deployer.GenerateHostFile(instances, fqdn, domain, burpFile, localFilePath, remoteFilePath,
			execCommand, socatPort, socatIP, nmapOutput, nmapCommands,
			cobaltStrikeLicense, cobaltStrikePassword, cobaltStrikeC2Path, cobaltStrikeFile, cobaltStrikeKillDate,
			ufwAction, ufwTCPPorts, ufwUDPPorts)

		deployer.WriteToFile("ansible/hosts.yml", hostFile)
		deployer.WriteToFile("ansible/main.yml", playbook)

		deployer.ExecAnsible("hosts.yml", "main.yml", "ansible")
	},
}

var nmap = &cobra.Command{
	Use:   "nmap",
	Short: "execute nmap",
	Long:  `splits and distributes the nmap job between the specified hosts and optionally performs randomization of targets for evasion`,
	Args: func(cmd *cobra.Command, args []string) error {
		_, err := deployer.ValidatePorts(nmapPorts)
		if err != nil {
			return err
		}
		_, err = deployer.ParseIPFile(nmapHostFile)
		if err != nil {
			return err
		}

		err = deployer.ValidateNumberOfInstances(commandIndices, "instance")
		if err != nil {
			return err
		}
		return nil
	},
	Run: func(cmd *cobra.Command, args []string) {
		apps := []string{"nmap", "nmap-exec"}

		playbook := deployer.GeneratePlaybookFile(apps)

		marshalledState := deployer.TerraformStateMarshaller()

		list := deployer.ListInstances(marshalledState)

		var instances []deployer.ListStruct

		for _, num := range commandIndices {
			instances = append(instances, list[num])
		}

		nmapCommands := deployer.SplitNmapCommandsIntoHosts(nmapPorts, nmapHostFile, nmapCommand, len(instances), nmapEvasive)

		hostFile := deployer.GenerateHostFile(instances, fqdn, domain, burpFile, localFilePath, remoteFilePath,
			execCommand, socatPort, socatIP, nmapOutput, nmapCommands,
			cobaltStrikeLicense, cobaltStrikePassword, cobaltStrikeC2Path, cobaltStrikeFile, cobaltStrikeKillDate,
			ufwAction, ufwTCPPorts, ufwUDPPorts)

		deployer.WriteToFile("ansible/hosts.yml", hostFile)
		deployer.WriteToFile("ansible/main.yml", playbook)

		deployer.ExecAnsible("hosts.yml", "main.yml", "ansible")
	},
}

var socatRedirect = &cobra.Command{
	Use:   "socat-redirect",
	Short: "socat redirector",
	Long:  "sets up a socat redirector on the specified port with the specifed target",
	Args: func(cmd *cobra.Command, args []string) error {
		return deployer.ValidateNumberOfInstances(commandIndices, "instance")
	},
	Run: func(cmd *cobra.Command, args []string) {
		apps := []string{"socat", "socat-exec"}

		playbook := deployer.GeneratePlaybookFile(apps)

		marshalledState := deployer.TerraformStateMarshaller()

		list := deployer.ListInstances(marshalledState)

		var instances []deployer.ListStruct

		for _, num := range commandIndices {
			instances = append(instances, list[num])
		}

		hostFile := deployer.GenerateHostFile(instances, fqdn, domain, burpFile, localFilePath, remoteFilePath,
			execCommand, socatPort, socatIP, nmapOutput, nmapCommands,
			cobaltStrikeLicense, cobaltStrikePassword, cobaltStrikeC2Path, cobaltStrikeFile, cobaltStrikeKillDate,
			ufwAction, ufwTCPPorts, ufwUDPPorts)

		deployer.WriteToFile("ansible/hosts.yml", hostFile)
		deployer.WriteToFile("ansible/main.yml", playbook)

		deployer.ExecAnsible("hosts.yml", "main.yml", "ansible")
	},
}

// var empireRun = &cobra.Command{
// 	Use:   "empire-run",
// 	Short: "runs powershell empire",
// 	Long:  `starts powershell empire in a screen session`,
// 	Args: func(cmd *cobra.Command, args []string) error {
// 		return deployer.ValidateNumberOfInstances(commandIndices)
// 	},
// 	Run: func(cmd *cobra.command, args []string) {
// 		apps := []string{"empire", "empire-exec"}

// 		playbook := deployer.GeneratePlaybookFile(apps)

// 		marshalledState := deployer.TerraformStateMarshaller()

// 		list := deployer.ListInstances(marshalledState)

// 		var instances []deployer.ListStruct

// 		for _, num := range commandIndices {
// 			instances = append(instances, list[num])
// 		}

// 		hostFile := deployer.GenerateHostFile(instances, fqdn, domain, burpFile, localFilePath, remoteFilePath,
// 			execCommand, socatPort, socatIP, nmapOutput, nmapCommands,
// 			cobaltStrikeLicense, cobaltStrikePassword, cobaltStrikeC2Path, cobaltStrikeFile, cobaltStrikeKillDate,
// 			ufwAction, ufwTCPPorts, ufwUDPPorts)

// 		deployer.WriteToFile("ansible/hosts.yml", hostFile)
// 		deployer.WriteToFile("ansible/main.yml", playbook)

// 		deployer.ExecAnsible("hosts.yml", "main.yml", "ansible")
// 	},
// }

var cobaltStrikeRun = &cobra.Command{
	Use:   "cobaltstrike-run",
	Short: "updates and runs cobalt strike teamserver",
	Long:  "updates the cobalt strike teamserver with the licensse and starts the teamserver with specified profile and password. if a cobalstrike tgz file is specified this command will also install cobalstrike before running it",
	Args: func(cmd *cobra.Command, args []string) error {
		match, _ := regexp.MatchString(`20[0-9]{2}\-[0=1]{1}[0-9]{1}\-[0-3]{1}[0-9]{1}`, cobaltStrikeKillDate)
		if !match {
			return fmt.Errorf("invalid kill date format, need YYYY-MM-DD")
		}
		err := deployer.ValidateNumberOfInstances(commandIndices, "instance")
		return err
	},
	Run: func(cmd *cobra.Command, args []string) {
		var apps []string
		if cobaltStrikeFile != "" {
			apps = []string{"cobaltstrike", "cobaltstrike-exec"}
		} else {
			apps = []string{"cobaltstrike-exec"}
		}

		playbook := deployer.GeneratePlaybookFile(apps)

		marshalledState := deployer.TerraformStateMarshaller()

		list := deployer.ListInstances(marshalledState)

		var instances []deployer.ListStruct

		for _, num := range commandIndices {
			instances = append(instances, list[num])
		}

		remoteFilePath = "/opt/cobaltstrike"

		hostFile := deployer.GenerateHostFile(instances, fqdn, domain, burpFile, localFilePath, remoteFilePath,
			execCommand, socatPort, socatIP, nmapOutput, nmapCommands,
			cobaltStrikeLicense, cobaltStrikePassword, cobaltStrikeC2Path, cobaltStrikeFile, cobaltStrikeKillDate,
			ufwAction, ufwTCPPorts, ufwUDPPorts)

		deployer.WriteToFile("ansible/hosts.yml", hostFile)
		deployer.WriteToFile("ansible/main.yml", playbook)

		deployer.ExecAnsible("hosts.yml", "main.yml", "ansible")
	},
}

var collaboratorRun = &cobra.Command{
	Use:   "collaborator-run",
	Short: "Starts burp collaborator server",
	Long:  "Checks for burp collaborator installation, installs if it does not exist, and starts it",
	Args: func(cmd *cobra.Command, args []string) error {
		err := deployer.ValidateNumberOfInstances(commandIndices, "instance")
		return err
	},
	Run: func(cmd *cobra.Command, args []string) {
		var apps []string
		if domain != "" && burpFile != "" {
			apps = []string{"collaborator", "collaborator-exec"}
			fmt.Println("WARNING: Its best to obtain your wildcard letsencrypt certificate prior to installation")
			fmt.Println("Do you still wish to continue?")
			if !deployer.AskForConfirmation() {
				return
			}
		} else {
			apps = []string{"collaborator-exec"}
		}

		playbook := deployer.GeneratePlaybookFile(apps)

		marshalledState := deployer.TerraformStateMarshaller()

		list := deployer.ListInstances(marshalledState)

		var instances []deployer.ListStruct

		for _, num := range commandIndices {
			instances = append(instances, list[num])
		}

		fqdn = domain

		hostFile := deployer.GenerateHostFile(instances, fqdn, domain, burpFile, localFilePath, remoteFilePath,
			execCommand, socatPort, socatIP, nmapOutput, nmapCommands,
			cobaltStrikeLicense, cobaltStrikePassword, cobaltStrikeC2Path, cobaltStrikeFile, cobaltStrikeKillDate,
			ufwAction, ufwTCPPorts, ufwUDPPorts)

		deployer.WriteToFile("ansible/hosts.yml", hostFile)
		deployer.WriteToFile("ansible/main.yml", playbook)

		deployer.ExecAnsible("hosts.yml", "main.yml", "ansible")

		if domain != "" && burpFile != "" {

			fmt.Println("Next Steps:")

			fmt.Println("1. Set this IP address to be both the primary and secondary nameserver for your domain")

			fmt.Println("Note: In order to have valid HTTPS on the collaborator server you must obtain a wildcard certificate from letsencrypt")
		}
	},
}

// ---
// # Starts up Cobalt Strike
// cobaltstrike_license
// public_ip
// password
// path_to_malleable_c2

// ---
// # Synchronize two directories on one remote host.
// remote_absolute_path
// host_absolute_path

func init() {
	rootCmd.AddCommand(exec)
	exec.AddCommand(command, nmap, socatRedirect, cobaltStrikeRun, collaboratorRun /*, empireRun*/)

	command.PersistentFlags().IntSliceVarP(&commandIndices, "id", "i", []int{}, "[Required] the id(s) for the remote server")
	command.MarkFlagRequired("id")
	command.PersistentFlags().StringVarP(&execCommand, "command", "c", "", "[Required] the command you want to execute on the remote server. must be encapsulated in quotes")
	command.MarkPersistentFlagRequired("command")

	nmap.PersistentFlags().IntSliceVarP(&commandIndices, "id", "i", []int{}, "[Required] the id(s) for the scanning servers")
	nmap.MarkPersistentFlagRequired("id")
	nmap.PersistentFlags().StringVarP(&nmapHostFile, "hostFile", "f", "", "[Required] the local filepath to the line seperated file containing the subnets and/or IP addresses")
	nmap.MarkPersistentFlagRequired("hostFile")
	nmap.PersistentFlags().StringSliceVarP(&nmapPorts, "ports", "p", []string{}, "[Required] the port range to be passed to nmap i.e 21-23,443-445,8080-8081,8443")
	nmap.MarkPersistentFlagRequired("ports")
	nmap.PersistentFlags().StringVarP(&nmapCommand, "nmapCommand", "n", "", "[Required] the base of the nmap command to be run. this will execlude the -iL,-p, and -oA options as they are aded elsewhere i.e. nmap -sV -sT --max-rate=250")
	nmap.MarkPersistentFlagRequired("nmapCommand")
	nmap.PersistentFlags().StringVarP(&nmapOutput, "nmapOutput", "o", "", "[Required] the local directory for output from the remote servers to be saved")
	nmap.MarkPersistentFlagRequired("nmapOutput")
	nmap.PersistentFlags().BoolVarP(&nmapEvasive, "nmapEvasion", "e", false, "[Optional] scan evasively. this will split up the job more randomly and take considerably longer. defaults to false")

	socatRedirect.PersistentFlags().IntSliceVarP(&commandIndices, "id", "i", []int{}, "[Required] the id(s) for the remote server")
	socatRedirect.MarkFlagRequired("id")
	socatRedirect.PersistentFlags().StringVarP(&socatPort, "port", "p", "", "[Required] the port you want to forward between your C2 server")
	socatRedirect.MarkPersistentFlagRequired("port")
	socatRedirect.PersistentFlags().StringVarP(&socatIP, "target", "t", "", "[Required] the ip address of your C2 server")
	socatRedirect.MarkPersistentFlagRequired("target")

	cobaltStrikeRun.PersistentFlags().IntSliceVarP(&commandIndices, "id", "i", []int{}, "[Required] the id for the remote server")
	cobaltStrikeRun.MarkFlagRequired("id")
	cobaltStrikeRun.PersistentFlags().StringVarP(&cobaltStrikeLicense, "license", "l", "", "[Required] cobaltstrike license")
	cobaltStrikeRun.MarkPersistentFlagRequired("license")
	cobaltStrikeRun.PersistentFlags().StringVarP(&cobaltStrikePassword, "password", "p", "", "[Required] password for the cobaltstrike teamserver")
	cobaltStrikeRun.MarkPersistentFlagRequired("password")
	cobaltStrikeRun.PersistentFlags().StringVarP(&cobaltStrikeC2Path, "c2", "c", "", "[Required] local path to the malleable C2 profile to be used with the teamserver")
	cobaltStrikeRun.MarkPersistentFlagRequired("c2")
	cobaltStrikeRun.PersistentFlags().StringVarP(&cobaltStrikeFile, "file", "f", "", "[Optional] local filepath of the cobaltstrike tgz file. use this option if you want this command to install cobaltstrike as well")
	cobaltStrikeRun.PersistentFlags().StringVarP(&cobaltStrikeKillDate, "kill", "k", "", "[Required] kill date for cobaltstrike beacons YYYY-MM-DD i.e. 2018-08-08")
	cobaltStrikeRun.MarkPersistentFlagRequired("kill")

	collaboratorRun.PersistentFlags().IntSliceVarP(&commandIndices, "id", "i", []int{}, "[Required] the id for the install (Required)")
	collaboratorRun.MarkPersistentFlagRequired("id")
	collaboratorRun.PersistentFlags().StringVarP(&domain, "domain", "d", "", "[Optional] domain for the burp collaborator install. use this option if you wish to install burp collaborator as well")
	collaboratorRun.PersistentFlags().StringVarP(&burpFile, "burpFile", "b", "", "[Optional] local file path to the burp pro jar file. use this option if you wish to install burp collaborator as well")

	// empireRun.PersistentFlags().IntSliceVarP(&commandIndices, "id", "i", []int{}, "Specify the id for the install (Required)")
	// empireRun.MarkPersistentFlagRequired("id")
}
