Welcome to hideNsneak.
===============================
![Alt text](assets/logo.png "hideNsneak")
This application assists in managing attack infrastructure for penetration testers by providing an interface to rapidly deploy, manage, and take down various cloud services. These include VMs, domain fronting, Cobalt Strike servers, API gateways, and firewalls.


Table of contents 
------------------
  * [Overview](#overview)
  * [Running locally](#running-locally)
  * [Commands](#commands)
  * [Organization](#organization)
  * [Contributions](#contributions)
  * [License](#license)


Overview
---------
hideNsneak provides a simple interface that allows penetration testers to build ephemeral infrastructure -- one that requires minimal overhead. 
hideNsneak can: 

* *`deploy`, `destroy`, and `list`*
	1. Cloud instances via EC2, Google Cloud, Digital Ocean, Azure, and Alibaba Cloud
	2. API Gateway (AWS)
	3. Domain fronts via CloudFront and Azure Cloudfront

* *Proxy into said infrastructure*
* *Send and receive files*
* *Port scanning via NMAP*
* *Remote installations of Burp Collab, Cobalt Strike, Socat, LetsEncrypt, GoPhish, and SQLMAP*


Running locally
---------------
*A few disclosures for V 1.0:*
* At this time, all hosts are assumed `Ubuntu 16.04 Linux`. In the future, we're hoping to add on a docker container to decrease initial setup time
* The only cloud providers currently setup are AWS and Digital Ocean
* *You need to make sure that go is installed.* Instructions can be found [here](https://golang.org/dl/)

1. Create a new AWS S3 bucket in `us-east-1`
	- Ensure this is not public as it will hold your terraform state
2. `go get github.com/rmikehodges/hideNsneak`
3. `cd $GOPATH/src/github.com/rmikehodges/hideNsneak`
4. `./setup.sh`
5. `cp config/example-config.json config/config.json` 
	- fill in the values
	- aws_access_id, aws_secret_key, aws_bucket_name, public_key, private_key, ec2_user, and do_user are required at minimum
  - all operators working on the same state must have config values filled in all the same fields
6. now you can use the program by running `./hidensneak [command]`

Commands
---------
* `hidensneak help` --> run this anytime to get available commands 

* `hidensneak instance deploy`
* `hidensneak instance destroy`
* `hidensneak instance list`

* `hidensneak api deploy`
* `hidensneak api destroy`
* `hidensneak api list`

* `hidensneak domainfront enable`
* `hidensneak domainfront disable`
* `hidensneak domainfront deploy`
* `hidensneak domainfront destroy`
* `hidensneak domainfront list`

* `hidensneak firewall add`
* `hidensneak firewall list`
* `hidensneak firewall delete`

* `hidensneak exec command -c`
* `hidensneak exec nmap`
* `hidensneak exec socat-redirect`
* `hidensneak exec cobaltstrike-run`
* `hidensneak exec collaborator-run`

* `hidensneak socks deploy`
* `hidensneak socks list`
* `hidensneak socks destroy`
* `hidensneak socks proxychains`
* `hidensneak socks socksd`

* `hidensneak install burp`
* `hidensneak install cobaltstrike`
* `hidensneak install socat`
* `hidensneak install letsencrypt`
* `hidensneak install gophish`
* `hidensneak install nmap`
* `hidensneak install sqlmap`

* `hidensneak file push`
* `hidensneak file pull`

For all commands, you can run `--help` after any of them to get guidance on what flags to use.


Organization
------------
* `_terraform` --> terraform modules 
* `_ansible` --> ansible roles and playbooks 
* `_assets` --> random assets for the beauty of this project
* `_cmd` --> frontend interface package 
* `_deployer` --> backend commands and structs
* `main.go` --> where the magic happens 

IAM Permissions
-------------
Google Domain Fronting
* App Engine API Enabled
* Cloud Functions API Enabled
* Project editor or higher permissions


Miscellaneous
-------------
A default security group `hideNsneak` is made in all AWS regions that is full-open. All instances are configured with `iptables` to *only allow port 22/tcp* upon provisioning. 

If your program starts throwing terraform errors indicating a resource is not found, then you may need to remove the problematic terraform resources. You can do this by running the following:

`cd $GOPATH/src/github.com/rmikehodges/hideNsneak/terraform`

`terraform state rm <name of problem resource>`

This resource will need to be cleaned up manually if it still exists.

Troubleshooting
---------------

Error: configuration for `module name here` is not present; a provider configuration block is required for all operations

This is usually due to artifacts being left in the state from old deployments. Below are instructions on how to remove those artifacts from your state. If they are live resources, they will need to be manually destroyed via the cloud provider's administration panel.
* `cd $GOPATH/src/github.com/rmikehodges/hideNsneak/terraform`
* `terraform state rm <module or resource name>`


Contributions
-------------
We would love to have you contribute to hideNsneak. Feel free to pull the repo and start contributing, we will review pull requests as we receive them. If you feel like some things need improvement or some features need adding, feel free to open up an issue and hopefully -- someone will pick it up. 


License 
-------
MIT
