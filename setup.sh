#!/bin/bash
terraformPathDir="/usr/local/bin/"

uname=$(uname | awk '{print tolower($0)}')
terraformLink="https://releases.hashicorp.com/terraform/0.11.7/terraform_0.11.7_${uname}_amd64.zip"

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
    curl -sLo /tmp/terraform_${uname}.zip $terraformLink; sudo unzip /tmp/terraform_${uname}.zip -d /usr/local/bin/
fi

echo "Installing Ansible...."

if ! exists "ansible"
then
    # #If on Mac and experiencing errors, use the following command
    # sudo CFLAGS=-Qunused-arguments CPPFLAGS=-Qunused-arguments pip install ansible
    sudo pip install ansible
fi

echo "Ensuring Go dependencies are met"
go get github.com/rmikehodges/hideNsneak

echo "Instantiating Backend DynamoDB Table"

cd terraform/backend
terraform init -input=true
terraform apply
cd ../../

echo "If this the table already exists, you are good to go"

echo "All requirements met!"
