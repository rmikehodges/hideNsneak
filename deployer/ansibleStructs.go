package deployer

type ansiblePretask struct {
	Name        string `yaml:"name"`
	Raw         string `yaml:"raw"`
	Register    string `yaml:"register"`
	ChangedWhen string `yaml:"changed_when"`
}

type ansiblePlaybook struct {
	Name        string           `yaml:"name"`
	Hosts       string           `yaml:"hosts"`
	Become      bool             `yaml:"become"`
	GatherFacts bool             `yaml:"gather_facts"`
	PreTasks    []ansiblePretask `yaml:"pre_tasks"`
	Roles       []string         `yaml:"roles"`
}

func (playbook *ansiblePlaybook) GenerateDefault() {
	playbook.Name = "install all packages"
	playbook.Hosts = "all"
	playbook.Become = true
	playbook.GatherFacts = false
	playbook.PreTasks = append(playbook.PreTasks, ansiblePretask{
		Name:        "initialization steps",
		Raw:         "test -e /usr/bin/python || (apt -y update && apt install -y python-minimal)",
		Register:    "output",
		ChangedWhen: `output.stdout != ""`,
	})
	playbook.Roles = []string{"common"}
}

type ansibleInventory struct {
	All struct {
		Hosts map[string]ansibleHost `yaml:"hosts"`
	} `yaml:"all"`
}

type ansibleHost struct {
	AnsibleHost           string   `yaml:"ansible_host"`
	AnsibleUser           string   `yaml:"ansible_user"`
	AnsiblePrivateKey     string   `yaml:"ansible_ssh_private_key_file"`
	AnsibleAdditionalOpts string   `yaml:"ansible_ssh_common_args"`
	AnsibleDomain         string   `yaml:"ansible_domain_name"`
	BurpFile              string   `yaml:"burp_file"`
	HostAbsPath           string   `yaml:"host_absolute_path"`
	RemoteAbsPath         string   `yaml:"remote_absolute_path"`
	ExecCommand           string   `yaml:"remote_command"`
	SocatPort             string   `yaml:"socat_port"`
	SocatIP               string   `yaml:"socat_ip"`
	NmapCommands          []string `yaml:"nmap_commands"`
	NmapOutput            string   `yaml:"nmap_output"`
	CobaltStrikeLicense   string   `yaml:"cobaltstrike_license"`
	CobaltStrikePassword  string   `yaml:"password"`
	CobaltStrikeC2Path    string   `yaml:"path_to_malleable_c2"`
	CobaltStrikeFile      string   `yaml:"cobaltstrike_file"`
	CobaltStrikeKillDate  string   `yaml:"cobaltstrike_kill_date"`
	UfwAction             string   `yaml:"ufw_action"`
	UfwTCPPort            []string `yaml:"ufw_tcp_port"`
	UfwUDPPort            []string `yaml:"ufw_udp_port"`
	KeyPath               string   `yaml:"key_path"`
	CertPath              string   `yaml:"cert_path"`
}
