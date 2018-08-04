package deployer

import (
	"context"
	"fmt"
	"io/ioutil"
	"os"
	"strings"

	"github.com/digitalocean/godo"
	"golang.org/x/crypto/ssh"
	"golang.org/x/oauth2"
)

func compareDOConfig(instanceOne DOConfigWrapper, instanceTwo DOConfigWrapper) bool {
	if instanceOne.Image == instanceTwo.Image &&
		instanceOne.Fingerprint == instanceTwo.Fingerprint &&
		instanceOne.PrivateKey == instanceTwo.PrivateKey &&
		instanceOne.Size == instanceTwo.Size &&
		instanceOne.DefaultUser == instanceTwo.DefaultUser {
		return true
	}
	return false
}

func genDOKeyFingerprint(publicKey string) (keyFingerprint string) {
	key, err := ioutil.ReadFile(publicKey)

	if err != nil {
		fmt.Println("Specified DO public key does not exist")
		os.Exit(1)
	}
	pubKey, _, _, _, err := ssh.ParseAuthorizedKey(key)

	if err != nil {
		fmt.Println("Specified DO public key is not formatted correctly")
		os.Exit(1)
	}

	return strings.TrimSpace(ssh.FingerprintLegacyMD5(pubKey))
}

//GetDoRegions returns the list of available regions for digital ocean
func GetDoRegions(configFile string) (regions []string) {
	config := createConfig(configFile)
	client := newDOClient(config.DigitaloceanToken)
	regions, err := doRegions(client)
	if err != nil {
		fmt.Println("Error retrieving region list for validation")
	}
	return
}

func doRegions(client *godo.Client) ([]string, error) {
	var slugs []string
	regions, _, err := client.Regions.List(context.TODO(), &godo.ListOptions{})
	if err != nil {
		return slugs, err
	}
	for _, r := range regions {
		slugs = append(slugs, r.Slug)
	}
	return slugs, nil
}

//Token provides a function to retrieve a new digitalocean token for the service to make API calls
func (t *Token) Token() (*oauth2.Token, error) {
	token := &oauth2.Token{
		AccessToken: t.AccessToken,
	}
	return token, nil
}

func newDOClient(token string) *godo.Client {
	t := &Token{AccessToken: token}
	oa := oauth2.NewClient(oauth2.NoContext, t)
	return godo.NewClient(oa)
}
