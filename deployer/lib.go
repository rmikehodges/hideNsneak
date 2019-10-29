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

	"gopkg.in/yaml.v2"
)

// var config configStruct

////////////////////////
//Miscellaneous Functions
////////////////////////

func RetrieveConfig(configFilePath string) (config configStruct) {
	config := createConfig(configFilePath)
	return
}

func UpdateConfig(configFilePath string) {
	config := createConfig(configFilePath)

	file, _ := json.MarshalIndent(data, "", " ")

	err = ioutil.WriteFile(configFilePath, file, 0500)
	if err != nil {
		fmt.Println("Unable to write to config")
	}

}

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
