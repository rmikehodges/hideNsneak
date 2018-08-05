package deployer

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"os/exec"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/cloudfront"
	"github.com/aws/aws-sdk-go/service/ec2"
)

//checkEc2KeyExistence queries the Amazon EC2 API for the keypairs with the specified keyname
//Returns true if the resulting array is > 0, false otherwise
func checkEC2KeyExistence(secret string, accessID string, region string, keyName string) bool {
	// keyFingerprint := genEC2KeyFingerprint(privateKey)

	svc := ec2.New(session.New(&aws.Config{
		Region:      aws.String(region),
		Credentials: credentials.NewStaticCredentials(accessID, secret, ""),
	}))
	keyPairOutput, _ := svc.DescribeKeyPairs(&ec2.DescribeKeyPairsInput{
		KeyNames: aws.StringSlice([]string{keyName}),
	})
	if len(keyPairOutput.KeyPairs) == 0 {
		return false
	}
	return true
}

func deleteCloudFront(id string, ETag string, secret string, accessID string) error {
	svc := cloudfront.New(session.New(&aws.Config{
		Credentials: credentials.NewStaticCredentials(accessID, secret, ""),
	}))
	_, err := svc.DeleteDistribution(&cloudfront.DeleteDistributionInput{
		Id:      aws.String(id),
		IfMatch: aws.String(ETag),
	})
	if err != nil {
		return fmt.Errorf("Error deleting instance, instance is now disabled: %s", err)
	}
	return nil
}

func importEC2Key(secret string, accessID string, region string, pubKey []byte, keyName string) error {
	svc := ec2.New(session.New(&aws.Config{
		Region:      aws.String(region),
		Credentials: credentials.NewStaticCredentials(accessID, secret, ""),
	}))

	_, err := svc.ImportKeyPair(&ec2.ImportKeyPairInput{
		KeyName:           aws.String(keyName),
		PublicKeyMaterial: pubKey,
	})

	return err
}

func genEC2KeyFingerprint(privateKey string) (keyFingerprint string) {
	args1 := []string{"pkey", "-in", privateKey, "-pubout", "-outform", "DER"}
	args2 := []string{"md5", "-c"}

	pipeReader, pipeWriter := io.Pipe()

	cmd1 := exec.Command("openssl", args1...)
	cmd2 := exec.Command("openssl", args2...)

	cmd1.Stdout = pipeWriter
	cmd2.Stdin = pipeReader

	var cmd2Out bytes.Buffer

	cmd2.Stdout = &cmd2Out

	cmd1.Start()
	cmd2.Start()
	cmd1.Wait()
	pipeWriter.Close()
	cmd2.Wait()

	keyFingerprint = strings.Split(strings.TrimSpace(cmd2Out.String()), " ")[1]

	return
}

func describeInstanceSecurityGroup(region string, ip string, secret string, accessID string,
	securityGroupName string) *ec2.DescribeSecurityGroupsOutput {
	svc := ec2.New(session.New(&aws.Config{
		Region:      aws.String(region),
		Credentials: credentials.NewStaticCredentials(accessID, secret, ""),
	}))
	securityGroups, err := svc.DescribeSecurityGroups(&ec2.DescribeSecurityGroupsInput{})
	if err != nil {
		log.Println("Error describing security group for AWS Instance")
	}
	return securityGroups
}

func createDefaultSecurityGroup(securityGroup string, region string, secret string, accessID string) error {
	svc := ec2.New(session.New(&aws.Config{
		Region:      aws.String(region),
		Credentials: credentials.NewStaticCredentials(accessID, secret, ""),
	}))
	_, err := svc.CreateSecurityGroup(&ec2.CreateSecurityGroupInput{
		GroupName:   aws.String(securityGroup),
		Description: aws.String("hidensneak default security group"),
	})

	if err != nil {
		if aerr, ok := err.(awserr.Error); ok {
			switch aerr.Code() {
			case "InvalidGroup.Duplicate":
				return nil
			}
		}
		return err
	}
	ec2Permissions := []*ec2.IpPermission{
		(&ec2.IpPermission{}).
			SetIpProtocol("tcp").
			SetFromPort(int64(0)).
			SetToPort(int64(65535)).
			SetIpRanges([]*ec2.IpRange{
				{CidrIp: aws.String("0.0.0.0/0")},
			}),
		(&ec2.IpPermission{}).
			SetIpProtocol("udp").
			SetFromPort(int64(0)).
			SetToPort(int64(65535)).
			SetIpRanges([]*ec2.IpRange{
				{CidrIp: aws.String("0.0.0.0/0")},
			}),
	}
	// fmt.Printf("Created security group %s with VPC %s.\n", aws.StringValue(createRes.GroupId), instance.Cloud.ID)
	_, err = svc.AuthorizeSecurityGroupIngress(&ec2.AuthorizeSecurityGroupIngressInput{
		GroupName:     aws.String(securityGroup),
		IpPermissions: ec2Permissions,
	})
	if err != nil {
		return err
	}
	return nil
}

//checkEc2KeyExistence queries the Amazon EC2 API for the security groups
//with the specified name
//Returns true if the resulting array is > 0, false otherwise
func checkEC2SecurityGroupExistence(secret string, accessID string, region string, securityGroupName string) (bool, string) {
	svc := ec2.New(session.New(&aws.Config{
		Region:      aws.String(region),
		Credentials: credentials.NewStaticCredentials(accessID, secret, ""),
	}))
	securityGroupOutput, _ := svc.DescribeSecurityGroups(&ec2.DescribeSecurityGroupsInput{
		GroupNames: aws.StringSlice([]string{securityGroupName}),
	})

	if len(securityGroupOutput.SecurityGroups) == 0 {
		return false, ""
	}

	return true, *securityGroupOutput.SecurityGroups[0].GroupId
}

func compareEC2Config(instanceOne EC2ConfigWrapper, instanceTwo EC2ConfigWrapper) bool {
	if instanceOne.DefaultUser == instanceTwo.DefaultUser &&
		instanceOne.InstanceType == instanceTwo.InstanceType &&
		instanceOne.KeyPairName == instanceTwo.KeyPairName {
		return true
	}
	return false
}
