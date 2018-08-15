package deployer

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"regexp"
	"strconv"
	"strings"
	"text/template"

	"gopkg.in/yaml.v2"
)

// var config configStruct

////////////////////////
//Miscellaneous Functions
////////////////////////

func createConfig(configFilePath string) (config configStruct) {
	var configContents, _ = ioutil.ReadFile(configFilePath)

	json.Unmarshal(configContents, &config)

	return
}

func AskForConfirmation() bool {
	var response string
	_, err := fmt.Scanln(&response)
	if err != nil {
		log.Fatal(err)
	}
	okayResponses := []string{"y", "Y", "yes", "Yes", "YES"}
	nokayResponses := []string{"n", "N", "no", "No", "NO"}
	if ContainsString(okayResponses, response) {
		return true
	} else if ContainsString(nokayResponses, response) {
		return false
	} else {
		fmt.Println("Please type yes or no and then press enter:")
		return AskForConfirmation()
	}
}

// You might want to put the following two functions in a separate utility package.

// posString returns the first index of element in slice.
// If slice does not contain element, returns -1.
func PosString(slice []string, element string) int {
	for index, elem := range slice {
		if elem == element {
			return index
		}
	}
	return -1
}

func checkErr(err error) {
	if err != nil {
		fmt.Println(err)
		log.Fatal()
	}
}
func templateCounter() func() int {
	i := -1
	return func() int {
		i++
		return i
	}
}

func removeSpaces(input string) (newString string) {
	newString = strings.ToLower(input)
	newString = strings.Replace(newString, " ", "_", -1)

	return
}

//ContainsString checks to see if the array contains the target string
func ContainsString(s []string, e string) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false
}

//ContainsInt checks to see if the array contains the target int
func ContainsInt(s []int, e int) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false
}

func execCmd(binary string, args []string, filepath string) string {
	var stdout, stderr bytes.Buffer

	cmd := exec.Command(binary, args...)
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	cmd.Dir = filepath

	err := cmd.Run()
	if err != nil {
		fmt.Println(stderr.String())
		log.Fatalf("cmd.Run() failed with %s\n", err)
	}

	return stdout.String()
}

//IsValidNumberInput takes in a string and checks if the numbers are valid
func IsValidNumberInput(input string) error {
	sliceToParse := strings.Split(input, ",")

	for _, num := range sliceToParse {
		_, err := strconv.Atoi(num)
		if err != nil {
			dashSlice := strings.Split(num, "-")
			if len(dashSlice) != 2 {
				return err
			} else {
				_, err := strconv.Atoi(dashSlice[0])
				if err != nil {
					return err
				}
				_, err = strconv.Atoi(dashSlice[1])
				if err != nil {
					return err
				}
			}
			continue
		}
	}
	return nil
}

//ExpandNumberInput expands input string and returns a list of ints
func ExpandNumberInput(input string) []int {
	var result []int
	sliceToParse := strings.Split(input, ",")

	for _, num := range sliceToParse {
		getInt, err := strconv.Atoi(num)
		if err != nil {
			sliceToSplit := strings.Split(num, "-")
			firstNum, err := strconv.Atoi(sliceToSplit[0])
			if err != nil {
				continue
			}
			secondNum, err := strconv.Atoi(sliceToSplit[1])
			if err != nil {
				continue
			}
			for i := firstNum; i <= secondNum; i++ {
				result = append(result, i)
			}
		} else {
			result = append(result, getInt)
		}
	}
	return result
}

func stringInSlice(a string, list []string) bool {
	for _, b := range list {
		if b == a {
			return true
		}
	}
	return false
}

func GetEC2DataToDestroy(instanceNames []string) (newInstanceNames []string) {
	var tempList []string
	newInstanceNames = instanceNames
	for _, name := range instanceNames {
		moduleNameList := strings.Split(name, ".")
		moduleNameList = moduleNameList[:4]
		moduleName := strings.Join(moduleNameList, ".")
		match, _ := regexp.MatchString(`module\.ec2Deploy[1-9]+\.module\.aws\-[a-zA-Z0-9-]+`, moduleName)
		if match {
			if !ContainsString(tempList, moduleName) {
				tempList = append(tempList, moduleName)
				dataElementList := []string{moduleName + ".aws_ami.ubuntu",
					moduleName + ".aws_subnet_ids.all", moduleName + ".aws_vpc.default"}
				newInstanceNames = append(newInstanceNames, dataElementList...)
			}
		}

	}
	return
}

//WriteToFile opens, clears and writes to file
func WriteToFile(path string, content string) {
	file, err := os.Create(path)
	checkErr(err)

	_, err = file.Write([]byte(content))
	checkErr(err)
	defer file.Close()
}

//ValidateNumberOfInstances makes sure that the number input is actually available in our list of active instances
func ValidateNumberOfInstances(numberInput []int, listType string, configFile string) error {
	marshalledState := TerraformStateMarshaller()

	switch listType {
	case "instance":
		list := ListInstances(marshalledState, configFile)
		largestInstanceNum := FindLargestNumber(numberInput)

		//make sure the largestInstanceNumToInstall is not bigger than totalInstancesAvailable
		if len(list) < largestInstanceNum {
			return errors.New("the number you entered is too big; try running `list` to see the number of instances you have")
		}
	case "api":
		list := ListAPIs(marshalledState)
		largestInstanceNum := FindLargestNumber(numberInput)

		//make sure the largestInstanceNumToInstall is not bigger than totalInstancesAvailable
		if len(list) < largestInstanceNum {
			return errors.New("the number you entered is too big; try running `list` to see the number of instances you have")
		}
	case "domainfront":
		list := ListDomainFronts(marshalledState)
		largestInstanceNum := FindLargestNumber(numberInput)

		//make sure the largestInstanceNumToInstall is not bigger than totalInstancesAvailable
		if len(list) < largestInstanceNum {
			return errors.New("the number you entered is too big; try running `list` to see the number of instances you have")
		}
	default:
		return fmt.Errorf("Unknown list type specified")
	}

	return nil
}

//InstanceDiff takes the old list of instances and the new list of instances and proceeds to
//check each instance in the new list against the old list. If its not in the old list, it
//appends it to output.
func InstanceDiff(instancesOld []ListStruct, instancesNew []ListStruct) (instancesOut []ListStruct) {
	if len(instancesOld) == 0 {
		instancesOut = instancesNew
	} else {
		for _, instance := range instancesNew {
			for index, check := range instancesOld {
				if check.IP == instance.IP {
					break
				}
				if index == len(instancesOld)-1 {
					instancesOut = append(instancesOut, instance)
					break
				}
			}
		}
	}

	return
}

/////////////////////
//Ansible Functions
/////////////////////

// getAnsibleDirectory expands the filepath for the user
func getAnsibleDirectory() (ansDirectory string) {
	goPath := os.Getenv("GOPATH")
	ansDirectory = goPath + "/src/github.com/rmikehodges/hideNsneak/ansible"

	return
}

//GeneratePlaybookFile generates an ansible playbook
func GeneratePlaybookFile(apps []string) string {
	var playbookStruct ansiblePlaybook

	playbookStruct.GenerateDefault()

	for _, app := range apps {
		playbookStruct.Roles = append(playbookStruct.Roles, app)
	}

	playbookList := []ansiblePlaybook{playbookStruct}

	playbook, err := yaml.Marshal(playbookList)

	if err != nil {
		fmt.Println("Error marshalling playbook")
	}

	return string(playbook)
}

//GenerateHostsFile generates an ansible host file
func GenerateHostFile(instances []ListStruct, domain string, burpFile string,
	hostFilePath string, remoteFilePath string, execCommand string, socatPort string, socatIP string, nmapOutput string, nmapCommands map[int][]string,
	cobaltStrikeLicense string, cobaltStrikePassword string, cobaltStrikeC2Path string, cobaltStrikeFile string, cobaltStrikeKillDate string,
	ufwAction string, ufwTcpPort []string, ufwUdpPort []string) string {
	var inventory ansibleInventory

	inventory.All.Hosts = make(map[string]ansibleHost)
	for index, instance := range instances {
		inventory.All.Hosts[instance.IP] = ansibleHost{
			AnsibleHost:           instance.IP,
			AnsiblePrivateKey:     instance.PrivateKey,
			AnsibleUser:           instance.Username,
			AnsibleAdditionalOpts: "-o StrictHostKeyChecking=no",
			AnsibleDomain:         domain,
			BurpFile:              burpFile,
			HostAbsPath:           hostFilePath,
			RemoteAbsPath:         remoteFilePath,
			ExecCommand:           execCommand,
			NmapCommands:          nmapCommands[index],
			NmapOutput:            nmapOutput,
			SocatPort:             socatPort,
			SocatIP:               socatIP,
			CobaltStrikeFile:      cobaltStrikeFile,
			CobaltStrikeLicense:   cobaltStrikeLicense,
			CobaltStrikeC2Path:    cobaltStrikeC2Path,
			CobaltStrikePassword:  cobaltStrikePassword,
			CobaltStrikeKillDate:  cobaltStrikeKillDate,
			UfwAction:             ufwAction,
			UfwTCPPort:            ufwTcpPort,
			UfwUDPPort:            ufwUdpPort,
		}
	}

	hostFile, err := yaml.Marshal(inventory)

	if err != nil {
		fmt.Println("problem marshalling inventory file")
	}

	return string(hostFile)
}

func ExecAnsible(hostsFile string, playbook string) {
	filepath := getAnsibleDirectory()

	// var stdout, stderr bytes.Buffer
	binary, err := exec.LookPath("ansible-playbook")

	checkErr(err)

	args := []string{"-i", hostsFile, playbook}
	cmd := exec.Command(binary, args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Dir = filepath

	err = cmd.Run()
	if err != nil {
		fmt.Println(err)
	}

	return
}

/////////////////////
//Terraform Functions
/////////////////////

// getTerraformDirectory expands the filepath for the user
func getTerraformDirectory() (tfDirectory string) {
	goPath := os.Getenv("GOPATH")
	tfDirectory = goPath + "/src/github.com/rmikehodges/hideNsneak/terraform"

	return
}

//Hack for backend config
func execBashTerraform(args string, filepath string) string {
	var stdout, stderr bytes.Buffer

	binary, err := exec.LookPath("terraform")

	args = binary + " " + args

	checkErr(err)

	cmd := exec.Command("/bin/bash", "-c", args)

	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	cmd.Dir = filepath

	err = cmd.Run()
	if err != nil {
		fmt.Println(stderr.String())
	}

	return stdout.String()
}

func execTerraform(args []string, filepath string) string {
	var stdout, stderr bytes.Buffer

	binary, err := exec.LookPath("terraform")

	checkErr(err)

	cmd := exec.Command(binary, args...)
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	cmd.Dir = filepath

	err = cmd.Run()
	if err != nil {
		fmt.Println(stderr.String())
	}

	return stdout.String()
}

//InitializeTerraformFiles Creates the base templates for
//the terraform infrastructure
func InitializeTerraformFiles(configFile string) {

	config := createConfig(configFile)

	secrets, err := template.New("secrets").Parse(templateSecrets)

	if err != nil {
		fmt.Println(err)
	}

	secretBuff := new(bytes.Buffer)

	err = secrets.Execute(secretBuff, &config)

	if err != nil {
		fmt.Println(err)
	}

	backendBucket, err := template.New("backend").Parse(backend)

	if err != nil {
		fmt.Println(err)
	}

	backendBuff := new(bytes.Buffer)

	err = backendBucket.Execute(backendBuff, &config)

	mainFile, err := os.Create(tfMainFile)
	checkErr(err)
	defer mainFile.Close()

	varFile, err := os.Create(tfVariablesFile)
	checkErr(err)
	defer varFile.Close()

	tfvarsFile, err := os.Create(tfVarsFile)
	checkErr(err)
	defer tfvarsFile.Close()

	fmt.Println("Creating Terraform Files...")

	mainFile.Write(backendBuff.Bytes())
	varFile.Write([]byte(variables))
	tfvarsFile.Write(secretBuff.Bytes())
}

//TerraformApply runs the init, plan, and apply commands for our
//generated terraform templates
func TerraformApply(configFile string) {
	filepath := getTerraformDirectory()
	config := createConfig(configFile)

	// Initializing Terraform
	args := "init -backend-config=\"access_key=" + config.AwsAccessID + "\" -backend-config=\"secret_key=" + config.AwsSecretKey + "\""

	execBashTerraform(args, filepath)

	//Applying Changes Identified in tfplan
	fmt.Println("Applying Terraform Changes...")
	argsSlice := []string{"apply", "-input=false", "-auto-approve"}

	execTerraform(argsSlice, filepath)

}

func TerraformDestroy(nameList []string, configFile string) {
	filepath := getTerraformDirectory()
	config := createConfig(configFile)

	//Initializing Terraform
	args := "init -backend-config=\"access_key=" + config.AwsAccessID + "\" -backend-config=\"secret_key=" + config.AwsSecretKey + "\""

	execBashTerraform(args, filepath)

	argsSlice := []string{"destroy", "-auto-approve"}

	for _, name := range nameList {
		argsSlice = append(argsSlice, "-target", name)
	}
	fmt.Println("Destroying Terraform Targets...")

	execTerraform(argsSlice, filepath)
}

//TerraforrmOutputMarshaller runs the terraform output command
//and marshalls the resulting JSON into a TerraformOutput struct
func TerraformStateMarshaller() (outputStruct State) {
	filepath := getTerraformDirectory()
	args := []string{"state", "pull"}
	output := execTerraform(args, filepath)

	json.Unmarshal([]byte(output), &outputStruct)

	return
}

//CreateTerraformMain takes in a string containing all the necessary calls
//for the main.tf file
func CreateTerraformMain(masterString string, configFile string) {

	InitializeTerraformFiles(configFile)

	//Opening Main.tf to append parsed template
	filepath := getTerraformDirectory()
	newFilePath := filepath + "/main.tf"
	mainFile, err := os.OpenFile(newFilePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	checkErr(err)

	//Writing the masterString to file. masterString was instantiated in master_list.go
	_, err = mainFile.WriteString(masterString)
	checkErr(err)

	err = mainFile.Close()
	checkErr(err)
}

func writeGoogleFrontFiles(googleFront GooglefrontConfigWrapper) (indexFilePath string, packageFilePath string) {
	indexFilePath = "/tmp/index.js"
	packageFilePath = "/tmp/package.json"

	indexFile, err := os.Create(indexFilePath)
	checkErr(err)
	packageFile, err := os.Create(packageFilePath)
	checkErr(err)

	t, err := template.New("index").Parse(googleDomainFrontCode)
	t.Execute(indexFile, googleFront)

	packageFile.WriteString(googlefrontPackage)

	return
}

//ProviderCheck takes in a user-defined array of
//providers and validates they are supported
func ProviderCheck(providerArray []string) bool {
	for _, p := range providerArray {
		if strings.ToUpper(p) != "AWS" &&
			strings.ToUpper(p) != "DO" &&
			strings.ToUpper(p) != "GOOGLE" &&
			strings.ToUpper(p) != "AZURE" {
			return false
		}
	}
	return true
}

func mergeMap(map1 map[string]string, map2 map[string]string) map[string]string {
	for key, value := range map2 {
		map1[key] = value
	}
	return map1
}

func FindLargestNumber(nums []int) int {
	var n, largest int
	for _, i := range nums {
		if i > n {
			n = i
			largest = n
		}
	}
	return largest
}

func PrintProxyChains(socksList string) (proxies string) {
	socksArray := strings.Split(socksList, "\n")
	for _, command := range socksArray {
		args := strings.Split(command, " ")
		for index, arg := range args {
			var port string
			if arg == "-D" {
				port = args[index+1]
				proxies = proxies + fmt.Sprintf("socks5 127.0.0.1 %s\n", port)
			}
		}
	}
	return
}

func PrintSocksd(socksList string) (proxies string) {
	socksArray := strings.Split(socksList, "\n")
	proxies = "{\"proxies\":[\n\t\"upstreams\":["
	for index, command := range socksArray {
		var port string
		var ip string
		args := strings.Split(command, " ")
		for index, arg := range args {

			if arg == "-D" {
				port = args[index+1]
			}
			if strings.Contains(arg, "@") {
				ip = strings.Split(arg, "@")[1]
			}
		}
		upstream := fmt.Sprintf("\t\t{\"type\": \"socks5\", \"address\": \"127.0.0.1:%s\", \"target\": \"%s\"}", port, ip)
		if index != len(socksArray)-1 {
			upstream = upstream + ","
		}
		proxies = proxies + "\n" + upstream

	}
	proxies = proxies + "\n\t\t]\n\t}\n\t]\n}"
	return
}

func DestroySOCKS(ip string) {
	args := []string{"-f", "ssh.*-D [0-9]{4,}.*@" + ip}
	cmd := exec.Command("pkill", args...)

	cmd.Run()
}

//createSingleSOCKS initiates a SOCKS Proxy on the local host with the specifed ipv4 address
func CreateSingleSOCKS(privateKey string, username string, ipv4 string, port int) (err error) {
	portString := strconv.Itoa(port)
	args := []string{"-D", portString, "-o", "StrictHostKeyChecking=no", "-N", "-f", "-i", privateKey, username + "@" + ipv4}
	cmd := exec.Command("ssh", args...)
	err = cmd.Start()
	if err != nil {
		fmt.Println(err)
		return
	}
	return

}

//listStructSOrt takes a list of listStructs and goes ahead and
//sorts them based on position. This is important because it will
//ensure that the order of the elements remain the same on each list call
func listStructSort(listStructs []ListStruct) (finalList []ListStruct) {
	for index := range listStructs {
		for _, list := range listStructs {
			if list.Place == index {
				finalList = append(finalList, list)
				break
			}
		}
	}
	return
}

func ListProxies(instances []ListStruct) (output string) {
	for _, instance := range instances {
		ip := instance.IP
		args := []string{"-f", "ssh.*-D.*" + ip}
		cmd := exec.Command("pgrep", args...)
		out, err := cmd.Output()

		if len(out) > 0 {
			pid := strings.TrimSpace(string(out))

			args = []string{"-o", "command", pid}
			cmd = exec.Command("ps", args...)
			out, err = cmd.Output()
			if err != nil {
				fmt.Println(err)
			}

			commandOutput := strings.Split(string(out), "COMMAND")[1]

			output = output + "PID: " + pid + " Command: " + strings.TrimSpace(commandOutput) + "\n"
		}

	}
	return
}

func ListDomainFronts(state State) (domainFronts []DomainFrontOutput) {
	for _, module := range state.Modules {
		var domainFrontOutput DomainFrontOutput
		if len(module.Path) > 1 && len(module.Resources) != 0 {
			for _, resource := range module.Resources {
				if strings.Contains(module.Path[1], "cloudfrontDeploy") {
					domainFrontOutput.Provider = "AWS"
					domainFrontOutput.ID = resource.Primary.Attributes["id"].(string)
					domainFrontOutput.Etag = resource.Primary.Attributes["etag"].(string)
					domainFrontOutput.Status = resource.Primary.Attributes["status"].(string)
					domainFrontOutput.Name = "module." + module.Path[1]
					for key, value := range resource.Primary.Attributes {
						if strings.Contains(key, "domain_name") {
							if strings.Contains(key, "origin") {
								domainFrontOutput.Origin = value.(string)
							} else {
								domainFrontOutput.Invoke = value.(string)
							}
						}
					}
					domainFronts = append(domainFronts, domainFrontOutput)
				} else if strings.Contains(module.Path[1], "azurefrontDeploy") {
					domainFrontOutput.Provider = "AZURE"
					// domainFronts = append(domainFronts, domainFrontOutput)
				} else if strings.Contains(module.Path[1], "googlefrontDeploy") {
					if resource.Type == "google_cloudfunctions_function" {

						domainFrontOutput.Provider = "GOOGLE"
						domainFrontOutput.Origin = resource.Primary.Attributes["labels.target"].(string)
						domainFrontOutput.Invoke = resource.Primary.Attributes["https_trigger_url"].(string)

						result, _ := strconv.ParseBool(resource.Primary.Attributes["trigger_http"].(string))
						if result {
							domainFrontOutput.Status = "Enabled"
						} else {
							domainFrontOutput.Status = "Disabled"
						}
						domainFrontOutput.Name = "module." + module.Path[1]
						domainFrontOutput.FunctionName = resource.Primary.Attributes["name"].(string)
						domainFrontOutput.RestrictUA = resource.Primary.Attributes["description"].(string)

						domainFronts = append(domainFronts, domainFrontOutput)
					}

				}
			}

		}
	}
	return
}

func ListAPIs(state State) (apiOutputs []APIOutput) {
	for _, module := range state.Modules {
		var apiOutput APIOutput
		if len(module.Path) > 1 && len(module.Resources) != 0 && strings.Contains(module.Path[1], "awsAPIDeploy") {
			apiOutput.Provider = "AWS"
			apiOutput.Name = "module." + module.Path[1]
			for _, resource := range module.Resources {
				switch resource.Type {
				case "aws_api_gateway_deployment":
					apiOutput.InvokeURI = resource.Primary.Attributes["invoke_url"].(string)
				case "aws_api_gateway_integration":
					apiOutput.TargetURI = resource.Primary.Attributes["uri"].(string)
				default:
					continue
				}
			}
			apiOutputs = append(apiOutputs, apiOutput)
		}

	}
	return
}

func ListInstances(state State, configFile string) (hostOutput []ListStruct) {
	config := createConfig(configFile)

	privateKey := config.PrivateKey

	for _, module := range state.Modules {
		if len(module.Resources) != 0 {
			var tempOutput []ListStruct
			if len(module.Path) > 1 {

				for name, resource := range module.Resources {
					fullName := "module." + strings.Join(module.Path[1:], ".module.") + "." + name
					nameSlice := strings.Split(name, ".")
					finalString := nameSlice[len(nameSlice)-1]
					count, err := strconv.Atoi(finalString)
					if err == nil {

						index := "[" + finalString + "]"

						newName := strings.Join(nameSlice[:len(nameSlice)-1], ".")

						fullName = "module." + strings.Join(module.Path[1:], ".module.") + "." + newName + index
					}
					switch resource.Type {
					case "digitalocean_droplet":
						tempOutput = append(tempOutput, ListStruct{
							IP:         resource.Primary.Attributes["ipv4_address"].(string),
							Provider:   "DigitalOcean",
							Region:     resource.Primary.Attributes["region"].(string),
							Name:       fullName,
							Place:      count,
							PrivateKey: privateKey,
							Username:   config.DOUser,
						})
					case "aws_instance":
						tempOutput = append(tempOutput, ListStruct{
							IP:         resource.Primary.Attributes["public_ip"].(string),
							Provider:   "AWS",
							Region:     resource.Primary.Attributes["availability_zone"].(string),
							Name:       fullName,
							Place:      count,
							PrivateKey: privateKey,
							Username:   config.EC2User,
						})
					default:
						continue
					}
				}
				hostOutput = append(hostOutput, listStructSort(tempOutput)...)
			}
		}
	}
	return
}

//InstanceDeploy takes input from the user interface in order to divide and deploy appropriate regions
//it takes in a TerraformOutput struct, makes the appropriate edits, and returns that same struct
func InstanceDeploy(providers []string, awsRegions []string, doRegions []string, azureRegions []string,
	googleRegions []string, count int, keyName string, wrappers ConfigWrappers, configFile string) ConfigWrappers {

	config := createConfig(configFile)
	doModuleCount := wrappers.DropletModuleCount
	awsModuleCount := wrappers.EC2ModuleCount

	//Gather the count per provider and the remainder
	countPerProvider := count / len(providers)

	remainderForProviders := count % len(providers)

	for _, provider := range providers {
		switch strings.ToUpper(provider) {
		case "AWS":
			countPerAWSregion := countPerProvider / len(awsRegions)

			remainderForAWSRegion := countPerProvider % len(awsRegions)

			if remainderForProviders > 0 {
				remainderForAWSRegion = remainderForAWSRegion + 1
				remainderForProviders = remainderForProviders - 1
			}

			for _, region := range awsRegions {

				regionCount := countPerAWSregion

				if remainderForAWSRegion > 0 {
					regionCount = regionCount + 1
					remainderForAWSRegion = remainderForAWSRegion - 1
				}
				//TODO: Add custom input
				if regionCount > 0 {
					result := checkEC2KeyExistence(config.AwsSecretKey, config.AwsAccessID, region, keyName)

					if !result {
						publicKeyBytes, err := ioutil.ReadFile(config.PublicKey)
						if err != nil {
							fmt.Printf("Error reading public key: %s", err)
						}

						err = importEC2Key(config.AwsSecretKey, config.AwsAccessID, region, publicKeyBytes, keyName)
						if err != nil {
							fmt.Printf("There was an errror importing your key to EC2: %s", err)
						} else {
							fmt.Println("Success for importing AWS key for region: " + region)
						}
					}

					//Attempting to add default security group, if it exists function will return nil
					err := createDefaultSecurityGroup("hidensneak", region, config.AwsSecretKey, config.AwsAccessID)

					if err != nil {
						fmt.Printf("Default security group creation failed: %s \n You may have to manually allow traffic to your instance", err)
					}

					newEC2RegionConfig := EC2ConfigWrapper{
						InstanceType: "t2.micro",
						PrivateKey:   config.PrivateKey,
						PublicKey:    config.PublicKey,
						KeyPairName:  keyName,
						DefaultUser:  "ubuntu",
						RegionMap:    make(map[string]int),
					}
					newEC2RegionConfig.RegionMap[region] = regionCount

					if len(wrappers.EC2) == 0 {
						awsModuleCount = 1
						newEC2RegionConfig.ModuleName = "ec2Deploy" + strconv.Itoa(awsModuleCount)
						awsModuleCount = awsModuleCount + 1
						wrappers.EC2 = append(wrappers.EC2, newEC2RegionConfig)
						continue
					}
					for index, config := range wrappers.EC2 {
						if compareEC2Config(config, newEC2RegionConfig) {
							if config.RegionMap[region] > 0 {
								config.RegionMap[region] = config.RegionMap[region] + regionCount
							} else {
								config.RegionMap[region] = regionCount
							}
							break
						} else if index == len(wrappers.DO)-1 {
							awsModuleCount = awsModuleCount + 1
							newEC2RegionConfig.ModuleName = "ec2Deploy" + strconv.Itoa(awsModuleCount)
							wrappers.EC2 = append(wrappers.EC2, newEC2RegionConfig)
							awsModuleCount = awsModuleCount + 1
						}
					}
				}
			}
		case "DO":
			countPerDOregion := countPerProvider / len(doRegions)

			remainderForDORegion := countPerProvider % len(doRegions)

			if remainderForProviders > 0 {
				remainderForDORegion = remainderForDORegion + 1
				remainderForProviders = remainderForProviders - 1
			}

			for _, region := range doRegions {
				regionCount := countPerDOregion

				if remainderForDORegion > 0 {
					regionCount = regionCount + 1
					remainderForDORegion = remainderForDORegion - 1
				}
				//TODO: Add custom input
				if regionCount > 0 {
					newDORegionConfig := DOConfigWrapper{
						Image:       "ubuntu-16-04-x64",
						PrivateKey:  config.PrivateKey,
						Fingerprint: genDOKeyFingerprint(config.PublicKey),
						Size:        "512mb",
						DefaultUser: "root",
						RegionMap:   make(map[string]int),
					}
					newDORegionConfig.RegionMap[region] = regionCount

					if len(wrappers.DO) == 0 {
						doModuleCount = 1
						newDORegionConfig.ModuleName = "doDropletDeploy" + strconv.Itoa(doModuleCount)
						doModuleCount = doModuleCount + 1
						wrappers.DO = append(wrappers.DO, newDORegionConfig)
						continue
					}
					for index, config := range wrappers.DO {
						if compareDOConfig(config, newDORegionConfig) {
							if config.RegionMap[region] > 0 {
								config.RegionMap[region] = config.RegionMap[region] + regionCount
							} else {
								config.RegionMap[region] = regionCount
							}
							break
						} else if index == len(wrappers.DO)-1 {
							doModuleCount = doModuleCount + 1
							newDORegionConfig.ModuleName = "doDropletDeploy" + strconv.Itoa(doModuleCount)
							wrappers.DO = append(wrappers.DO, newDORegionConfig)
							doModuleCount = doModuleCount + 1
						}
					}
				}
			}
		default:
			continue
		}
	}
	return wrappers
}

//APIDeploy takes argruments to deploy an API Gateway
func APIDeploy(provider string, targetURI string, wrappers ConfigWrappers) ConfigWrappers {
	moduleCount := wrappers.AWSAPIModuleCount

	if strings.ToUpper(provider) == "AWS" {
		if len(wrappers.AWSAPI) > 0 {
			for _, wrapper := range wrappers.AWSAPI {
				if targetURI == wrapper.TargetURI {
					continue
				}
				wrappers.AWSAPI = append(wrappers.AWSAPI, AWSApiConfigWrapper{
					ModuleName: "awsAPIDeploy" + strconv.Itoa(moduleCount+1),
					TargetURI:  targetURI,
				})
				moduleCount = moduleCount + 1
			}
		} else {
			wrappers.AWSAPI = append(wrappers.AWSAPI, AWSApiConfigWrapper{
				ModuleName: "awsAPIDeploy" + strconv.Itoa(moduleCount+1),
				TargetURI:  targetURI,
			})
			moduleCount = moduleCount + 1
		}
	} else if strings.ToUpper(provider) == "ALIBABA" {
	}
	wrappers.AWSAPIModuleCount = moduleCount
	return wrappers
}

func DomainFrontDeploy(provider string, origin string, restrictUA string,
	functionName string, frontedDomain string, wrappers ConfigWrappers) ConfigWrappers {
	cloudfrontmMduleCount := wrappers.CloudfrontModuleCount

	googlefrontModuleCount := wrappers.GooglefrontModuleCount
	switch strings.ToUpper(provider) {
	case "AWS":
		if len(wrappers.Cloudfront) > 0 {
			for _, wrapper := range wrappers.Cloudfront {
				if origin == wrapper.Origin {
					continue
				}
				wrappers.Cloudfront = append(wrappers.Cloudfront, CloudfrontConfigWrapper{
					ModuleName: "cloudfrontDeploy" + strconv.Itoa(cloudfrontmMduleCount+1),
					Origin:     origin,
					Enabled:    "true",
				})
			}
		} else {
			wrappers.Cloudfront = append(wrappers.Cloudfront, CloudfrontConfigWrapper{
				ModuleName: "cloudfrontDeploy" + strconv.Itoa(cloudfrontmMduleCount+1),
				Origin:     origin,
				Enabled:    "true",
			})
		}
	case "GOOGLE":
		if len(wrappers.Googlefront) > 0 {
			for _, wrapper := range wrappers.Googlefront {
				if origin == wrapper.HostURL {
					continue
				}
				tempConfig := GooglefrontConfigWrapper{
					ModuleName:    "googlefrontDeploy" + strconv.Itoa(googlefrontModuleCount+1),
					FrontedDomain: frontedDomain,
					HostURL:       "https://" + origin,
					Host:          origin,
					Enabled:       true,
					RestrictUA:    restrictUA,
					FunctionName:  functionName,
				}
				indexFile, packageFile := writeGoogleFrontFiles(tempConfig)

				tempConfig.SourceFile = indexFile
				tempConfig.PackageFile = packageFile

				wrappers.Googlefront = append(wrappers.Googlefront, tempConfig)
			}
		} else {
			tempConfig := GooglefrontConfigWrapper{
				ModuleName:    "googlefrontDeploy" + strconv.Itoa(googlefrontModuleCount+1),
				FrontedDomain: frontedDomain,
				HostURL:       "https://" + origin,
				Host:          origin,
				Enabled:       true,
				RestrictUA:    restrictUA,
				FunctionName:  functionName,
			}
			indexFile, packageFile := writeGoogleFrontFiles(tempConfig)

			tempConfig.SourceFile = indexFile
			tempConfig.PackageFile = packageFile

			wrappers.Googlefront = append(wrappers.Googlefront, tempConfig)
			wrappers.GooglefrontModuleCount = wrappers.GooglefrontModuleCount + 1
		}
	case "AZURE":
		fmt.Println("Current unsupported provider")
		return wrappers
	default:
		return wrappers
	}

	return wrappers
}

//AWSCloufFrontDestroy uses the deleteCloudFront function to delete
//the specified cloudfront due to the problems with terraforms destruction process
func AWSCloudFrontDestroy(output DomainFrontOutput, configFile string) error {
	config := createConfig(configFile)
	filepath := getTerraformDirectory()

	err := deleteCloudFront(output.ID, output.Etag, config.AwsAccessID, config.AwsSecretKey)
	if err != nil {
		return err
	}

	args := []string{"state", "rm", output.Name}
	execTerraform(args, filepath)
	return nil
}
