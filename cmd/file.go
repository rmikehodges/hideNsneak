package cmd

import (
	"fmt"

	"github.com/rmikehodges/hideNsneak/deployer"

	"github.com/spf13/cobra"
)

var localFilePath string
var remoteFilePath string
var instanceFileIndex []int

// helloCmd represents the hello command
var file = &cobra.Command{
	Use:   "file",
	Short: "file",
	Long:  `file`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Run 'file --help' for usage.")
	},
}

var filePush = &cobra.Command{
	Use:   "push",
	Short: "send a file or directory",
	Long:  `send a file or directory from your local host to a remote server via absolute filepath`,
	Run: func(cmd *cobra.Command, args []string) {
		apps := []string{"sync-push"}
		playbook := deployer.GeneratePlaybookFile(apps)

		marshalledState := deployer.TerraformStateMarshaller()

		list := deployer.ListInstances(marshalledState, cfgFile)

		var instances []deployer.ListStruct

		for _, num := range instanceFileIndex {
			instances = append(instances, list[num])
		}

		hostFile := deployer.GenerateHostFile(instances, domain, burpFile, localFilePath, remoteFilePath,
			execCommand, socatPort, socatIP, nmapOutput, nmapCommands,
			cobaltStrikeLicense, cobaltStrikePassword, cobaltStrikeC2Path, cobaltStrikeFile, cobaltStrikeKillDate,
			ufwAction, ufwTCPPorts, ufwUDPPorts)

		deployer.WriteToFile("ansible/hosts.yml", hostFile)
		deployer.WriteToFile("ansible/main.yml", playbook)

		deployer.ExecAnsible("hosts.yml", "main.yml")
	},
}

var filePull = &cobra.Command{
	Use:   "pull",
	Short: "get a file or directory",
	Long:  `get a file or directory from your remote server to your local host via absolute filepath`,
	Run: func(cmd *cobra.Command, args []string) {
		apps := []string{"sync-pull"}

		playbook := deployer.GeneratePlaybookFile(apps)

		marshalledState := deployer.TerraformStateMarshaller()

		list := deployer.ListInstances(marshalledState, cfgFile)

		var instances []deployer.ListStruct

		for _, num := range instanceFileIndex {
			instances = append(instances, list[num])
		}

		hostFile := deployer.GenerateHostFile(instances, domain, burpFile, localFilePath, remoteFilePath,
			execCommand, socatPort, socatIP, nmapOutput, nmapCommands,
			cobaltStrikeLicense, cobaltStrikePassword, cobaltStrikeC2Path, cobaltStrikeFile, cobaltStrikeKillDate,
			ufwAction, ufwTCPPorts, ufwUDPPorts)

		deployer.WriteToFile("ansible/hosts.yml", hostFile)
		deployer.WriteToFile("ansible/main.yml", playbook)

		deployer.ExecAnsible("hosts.yml", "main.yml")
	},
}

func init() {
	rootCmd.AddCommand(file)
	file.AddCommand(filePush, filePull)

	filePush.PersistentFlags().IntSliceVarP(&instanceFileIndex, "id", "i", []int{}, "[Required] the id(s) for the remote server i.e. 1 or 1,2,3")
	filePush.MarkFlagRequired("id")
	filePush.PersistentFlags().StringVarP(&localFilePath, "local", "l", "", "[Required] the local file or directory absolute path")
	filePush.MarkPersistentFlagRequired("host")
	filePush.PersistentFlags().StringVarP(&remoteFilePath, "remote", "r", "", "[Required] the remote directory path to write to")
	filePush.MarkPersistentFlagRequired("remote")

	filePull.PersistentFlags().IntSliceVarP(&instanceFileIndex, "id", "i", []int{}, "[Required] the id(s) for the remote server i.e. 1 or 1,2,3")
	filePull.MarkFlagRequired("id")
	filePull.PersistentFlags().StringVarP(&localFilePath, "local", "l", "", "[Required] the local directory path to write to")
	filePull.MarkPersistentFlagRequired("host")
	filePull.PersistentFlags().StringVarP(&remoteFilePath, "remote", "r", "", "[Required] the remote file or directory absolute path")
	filePull.MarkPersistentFlagRequired("remote")
}
