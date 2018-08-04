#!/bin/bash
terraformPathDir="/usr/local/bin/"

uname=$(uname | awk '{print tolower($0)}')
terraformLink="https://releases.hashicorp.com/terraform/0.11.7/terraform_0.11.7_${uname}_amd64.zip"
ansibleProviderLink="https://github.com/nbering/terraform-provider-ansible/releases/download/v0.0.4/terraform-provider-ansible-${uname}_amd64.zip"
ansibleProviderFile="terraform-provider-ansible_v0.0.4"

REQUIRED_COMMANDS=('unzip' 'python' 'pip' 'go')
MISSING_COMMANDS=()

exists() {
    cmd=$(command -v "$1")
    [[ -n "$cmd" && -x "$cmd" ]]
}

echo "Checking initial requirements"

# Once we support Windows, remove this check
if [[ "$uname" != "darwin" && "$uname" != "linux" ]]
then
    echo "System is not Linux or macOS, program cannot be executed"
    exit 1
fi

for c in "${REQUIRED_COMMANDS[@]}"
do
    if ! exists $c
    then
        MISSING_COMMANDS=("${MISSING_COMMANDS[@]}" $c)
    fi
done

if [ "${#MISSING_COMMANDS[@]}" -ne 0 ]
then
    echo "error, missing commands: ${MISSING_COMMANDS[@]}"
    exit 1
fi

echo "Checking Terraform Installation...."
if ! exists "terraform"
then
    curl -sLo /tmp/terraform_${uname}.zip $terraformLink; unzip /tmp/terraform_${uname}.zip -d /usr/local/bin/
fi

echo "Installing Ansible Provider...."
if ! [ -f $HOME/.terraform.d/plugins/$ansibleProviderFile ]
then
    if ! [ -d $HOME/.terraform.d/plugins ]
    then
        mkdir -p $HOME/.terraform.d/plugins
    fi

    curl -sLo "$HOME/.terraform.d/plugins/$ansibleProviderFile.zip" $ansibleProviderLink
    unzip "$HOME/.terraform.d/plugins/$ansibleProviderFile.zip" -d $HOME/.terraform.d/plugins/
    rm "$HOME/.terraform.d/plugins/$ansibleProviderFile.zip"
    mv $HOME/.terraform.d/plugins/*/$ansibleProviderFile $HOME/.terraform.d/plugins/
    rm -r "$HOME/.terraform.d/plugins/${uname}_amd64"
fi

echo "Cleaning up provider files..."

echo "Installing Ansible...."

if ! exists "ansible"
then
    # #If on Mac and experiencing errors, use the following command
    # sudo CFLAGS=-Qunused-arguments CPPFLAGS=-Qunused-arguments pip install ansible
    sudo pip install ansible
fi

echo "Grabbing Go dependencies"
go get github.com/rmikehodges/hideNsneak/deployer
go get github.com/rmikehodges/hideNsneak/cmd

echo "Instantiating Backend DynamoDB Table"

cd terraform/backend
terraform init -input=true
terraform apply
cd ../../

echo "If this the table already exists, you are good to go"

echo "All requirements met!"