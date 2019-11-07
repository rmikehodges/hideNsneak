Welcome to hideNsneak.
===============================
![Alt text](assets/logo.png "hideNsneak")
This application assists in managing attack infrastructure for penetration testers by providing an interface to rapidly deploy, manage, and take down various cloud services. These include VMs, domain fronting, Cobalt Strike servers, API gateways, and firewalls.

Black Hat Arsenal Video Demo Video - https://youtu.be/8YTYScLn7pA

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
hideNsneak provides a simple interface that allows penetration testers to build and manage infrastructure -- one that requires minimal overhead.

hideNsneak can: 

* *`deploy`, `destroy`, and `list`*
	1. Cloud instances via EC2 and Digital Ocean (Google Cloud, Azure, and Alibaba Cloud coming soon)
	2. API Gateway (AWS)
	3. Domain fronts via AWS Cloudfront and Google Cloud Functions (Azure CDN coming soon)

* *Proxy through infrastructure*
* *Deploy C2 redirectors*
* *Send and receive files*
* *Distributed Port Scanning*
* *Remote installations of Burp Collaborator, Cobalt Strike, Socat, LetsEncrypt, GoPhish, and SQLMAP*
* *Share and manage infrastructure across various teams.

Running locally
---------------
*A few disclosures for V 1.1:*
* At this time, all hosts are assumed `Ubuntu 16.04 Linux`.
* SSH Agents aren't supported
* All config changes can be done via `hideNsneak setup`
* `setup.sh` script now just checks and installs dependencies

1. Create a new AWS S3 bucket in `us-east-1`
	- Ensure this is not public as it will hold your terraform state
2. `./setup.sh`
3. `go build -o hideNsneak main.go`
4. `hideNsneak setup aws -s <secret> -a <access> -b <bucket>`
  - This step is required for first run
5. `hideNsneak setup --help`
6. From here you can setup your other keys for other services or run solely on AWS


Commands
---------
* `hidensneak help` --> run this anytime to get available commands 

* `hidensneak setup aws`
* `hidensneak setup azure`
* `hidensneak setup do`
* `hidensneak setup ssh`

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



Error: Error locking state: Error acquiring the state lock: ConditionalCheckFailedException: The conditional request failed
	status code: 400, request id: P7BUM7NA56LQEJQC20A3SE2SOVVV4KQNSO5AEMVJF66Q9ASUAAJG
Lock Info:
  ID:        4919d588-6b29-4aa7-d917-2bcb67c14ab4
  
If this does not go away after another user has finished deploying then it is usually due to to Terraform not automatically unlocking your state in the face of errors. This can be fixed by running the following:
* `terraform force-unlock <ID> $GOPATH/src/github.com/rmikehodges/hideNsneak/terraform`

Note that this will unlock the state so it may have an adverse affect on any other writes happening in the state so make sure your other users are not actively deploying/destroying anything when you run this.

If you encounter an error along the following:

* `Error: module.googlefrontDeploy2.google_storage_bucket.bucket: configuration for module.googlefrontDeploy2.provider.google is not present; a provider configuration block is required for all operations`

This often means that there are items in the state you are not accounting for. This can be remediated by performing the following:

* `cd terraform`
* `terraform state list` - this will provide you the list of resources in the state
* `terraform state rm <offending module or resource from above>` - this will remove the offending resource from the state and you should be good to go. If the resource is still active then ensure you manually delete it.



Contributions
-------------
We would love to have you contribute to hideNsneak. Feel free to pull the repo and start contributing, we will review pull requests as we receive them. If you feel like some things need improvement or some features need adding, feel free to open up an issue and hopefully -- someone will pick it up. 


License 
-------
[MIT](LICENSE)
