package cmd

import (
	"fmt"
	"strings"
	"time"

	"github.com/rmikehodges/hideNsneak/deployer"

	"github.com/schollz/progressbar"
	"github.com/spf13/cobra"
)

var instanceProviders []string
var instancePrivateKey string
var instancePublicKey string
var instanceCount int
var regionAws []string
var regionDo []string
var regionAzure []string
var regionGoogle []string
var instanceDestroyIndices string

var instance = &cobra.Command{
	Use:   "instance",
	Short: "instance parent command",
	Long:  `parent command for managing instances`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Run 'instance --help' for usage.")
	},
}

var instanceDeploy = &cobra.Command{
	Use:   "deploy",
	Short: "deploys instances",
	Long:  `deploys instances for AWS, Azure, Digital Ocean, or Google Cloud`,
	Args: func(cmd *cobra.Command, args []string) error {
		if !deployer.ProviderCheck(instanceProviders) {
			return fmt.Errorf("invalid providers specified: %v", instanceProviders)
		}
		if deployer.ContainsString(instanceProviders, "DO") {
			availableDORegions := deployer.GetDoRegions(cfgFile)
			var unavailableRegions []string
			for _, region := range regionDo {
				if !deployer.ContainsString(availableDORegions, strings.ToLower(region)) {
					unavailableRegions = append(unavailableRegions, region)
				}
			}
			if len(unavailableRegions) != 0 {
				return fmt.Errorf("digitalocean region(s) not available: %s", strings.Join(unavailableRegions, ","))
			}
		}

		return nil

	},
	Run: func(cmd *cobra.Command, args []string) {

		marshalledState := deployer.TerraformStateMarshaller()
		wrappers := deployer.CreateWrappersFromState(marshalledState)

		oldList := deployer.ListInstances(marshalledState)

		wrappers = deployer.InstanceDeploy(instanceProviders, regionAws, regionDo, regionAzure, regionGoogle, instanceCount, instancePrivateKey, instancePublicKey, "hidensneak", wrappers, cfgFile)

		mainFile := deployer.CreateMasterFile(wrappers)

		deployer.CreateTerraformMain(mainFile, cfgFile)

		deployer.TerraformApply(cfgFile)

		fmt.Println("Waiting for instances to initialize...")

		bar := progressbar.New(120)
		for i := 0; i < 120; i++ {
			bar.Add(1)
			time.Sleep(1 * time.Second)
		}
		fmt.Println("")
		fmt.Println("Restricting Ports to only port 22...")

		marshalledState = deployer.TerraformStateMarshaller()
		newList := deployer.ListInstances(marshalledState)
		firewallList := deployer.InstanceDiff(oldList, newList)

		apps := []string{"firewall"}
		playbook := deployer.GeneratePlaybookFile(apps)

		ufwTCPPorts = []string{"22"}
		ufwAction = "add"

		hostFile := deployer.GenerateHostFile(firewallList, fqdn, domain, burpFile, localFilePath, remoteFilePath,
			execCommand, socatPort, socatIP, nmapOutput, nmapCommands,
			cobaltStrikeLicense, cobaltStrikePassword, cobaltStrikeC2Path, cobaltStrikeFile, cobaltStrikeKillDate,
			ufwAction, ufwTCPPorts, ufwUDPPorts)

		deployer.WriteToFile("ansible/hosts.yml", hostFile)
		deployer.WriteToFile("ansible/main.yml", playbook)

		deployer.ExecAnsible("hosts.yml", "main.yml")

	},
}

var instanceDestroy = &cobra.Command{
	Use:   "destroy",
	Short: "destroys instances",
	Long:  `destroys instances by choosing an index`,
	Args: func(cmd *cobra.Command, args []string) error {
		err := deployer.IsValidNumberInput(instanceDestroyIndices)

		if err != nil {
			return err
		}

		expandedNumIndex := deployer.ExpandNumberInput(instanceDestroyIndices)

		err = deployer.ValidateNumberOfInstances(expandedNumIndex, "instance")

		if err != nil {
			return err
		}

		return err
	},
	Run: func(cmd *cobra.Command, args []string) {
		marshalledState := deployer.TerraformStateMarshaller()

		list := deployer.ListInstances(marshalledState)

		var namesToDelete []string

		expandedNumIndex := deployer.ExpandNumberInput(instanceDestroyIndices)

		for _, numIndex := range expandedNumIndex {
			namesToDelete = append(namesToDelete, list[numIndex].Name)
		}

		emptyEC2Modules := deployer.CheckForEmptyEC2Module(namesToDelete, marshalledState)

		namesToDelete = append(namesToDelete, emptyEC2Modules...)

		deployer.TerraformDestroy(namesToDelete, cfgFile)
		return
	},
}

var instanceList = &cobra.Command{
	Use:   "list",
	Short: "detailed list of instances",
	Long:  `list instances and shows their index, IP, provider, region, and name`,
	Run: func(cmd *cobra.Command, args []string) {
		marshalledState := deployer.TerraformStateMarshaller()

		list := deployer.ListInstances(marshalledState)

		for index, item := range list {
			fmt.Print(index)
			fmt.Println(" : " + item.String())
		}
	},
}

func init() {
	rootCmd.AddCommand(instance)
	instance.AddCommand(instanceDeploy, instanceDestroy, instanceList)

	instanceDeploy.PersistentFlags().StringSliceVarP(&instanceProviders, "providers", "p", nil, "[Required] comma seperated list conatinaing any of the following available providers: AWS,DO")
	instanceDeploy.MarkPersistentFlagRequired("providers")

	instanceDeploy.PersistentFlags().IntVarP(&instanceCount, "count", "c", 0, "[Required] number of instances to deploy")
	instanceDeploy.MarkPersistentFlagRequired("count")

	instanceDeploy.PersistentFlags().StringVarP(&instancePrivateKey, "privatekey", "v", "", "[Required] local file path to private ssh key to connect to instance")
	instanceDeploy.MarkPersistentFlagRequired("privatekey")

	instanceDeploy.PersistentFlags().StringVarP(&instancePublicKey, "publickey", "b", "", "[Required] local file path to the public key corresponding to the specified private key")
	instanceDeploy.MarkPersistentFlagRequired("publickey")

	instanceDestroy.PersistentFlags().StringVarP(&instanceDestroyIndices, "input", "i", "", "[Required] indices of instances to destroy")
	instanceDestroy.MarkPersistentFlagRequired("input")

	//TODO: default all regions
	instanceDeploy.PersistentFlags().StringSliceVar(&regionAws, "region-aws",
		[]string{"us-east-1", "us-east-2", "us-west-1", "us-west-2", "ca-central-1", "eu-central-1", "eu-west-1", "eu-west-2", "eu-west-3", "ap-northeast-1", "ap-northeast-2", "ap-southeast-1", "ap-southeast-2", "ap-south-1", "sa-east-1"},
		"[Optional] comma seperated list of regions for aws")
	instanceDeploy.PersistentFlags().StringSliceVar(&regionDo, "region-do",
		[]string{"nyc1", "sgp1", "lon1", "nyc3", "ams3", "fra1", "tor1", "sfo2", "blr1"},
		"[Optional] comma sperated list of digital ocean regions")
	instanceDeploy.PersistentFlags().StringSliceVar(&regionAzure, "region-azure", []string{"westus", "centralus"}, "[Optional] comma seperated list of regions for azure")
	instanceDeploy.PersistentFlags().StringSliceVar(&regionGoogle, "region-google", []string{"us-west1", "us-east1"}, "[Optional] comma seperated list of regions for google cloud")

}
