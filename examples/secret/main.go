package main

import (
	"bufio"
	"context"
	"flag"
	"fmt"
	"os"
	"path/filepath"

	metav1 "github.com/opsdata/common-base/pkg/meta/v1"
	v1 "github.com/opsdata/elmt-api/apiserver/v1"

	"github.com/opsdata/elmt-sdk/tools/clientcmd"
	"github.com/opsdata/elmt-sdk/wyvern/service/elmt"
)

func main() {
	var elmtconfig *string

	if home := os.Getenv("HOME"); home != "" {
		elmtconfig = flag.String("elmtconfig", filepath.Join(home, ".elmt", "config"), "absolute path to the elmtconfig file")
	} else {
		elmtconfig = flag.String("elmtconfig", "", "absolute path to the elmtconfig file")
	}

	flag.Parse()

	// Use the current context in elmtconfig
	config, err := clientcmd.BuildConfigFromFlags("", *elmtconfig)
	if err != nil {
		panic(err.Error())
	}

	// Create the elmtclient
	elmtclient, err := elmt.NewForConfig(config)
	if err != nil {
		panic(err.Error())
	}

	secretsClient := elmtclient.APIV1().Secrets()

	secret := &v1.Secret{
		ObjectMeta: metav1.ObjectMeta{
			Name: "sdk",
		},
		Expires:     3724075800,
		Description: "test secret for sdk",
	}

	// Create secret
	fmt.Println("Creating secret...")
	ret, err := secretsClient.Create(context.TODO(), secret, metav1.CreateOptions{})
	if err != nil {
		panic(err.Error())
	}
	fmt.Printf("Created secret %q.\n", ret.Name)

	// Delete secret
	defer func() {
		fmt.Println("Deleting secret...")
		if err := secretsClient.Delete(context.TODO(), "sdk", metav1.DeleteOptions{}); err != nil {
			fmt.Printf("Delete secret failed: %s\n", err.Error())
			return
		}
		fmt.Println("Deleted secret.")
	}()

	// Get secret
	prompt()
	fmt.Println("Geting secret...")
	ret, err = secretsClient.Get(context.TODO(), "sdk", metav1.GetOptions{})
	if err != nil {
		panic(err.Error())
	}

	fmt.Printf("Get secret %q.\n", ret.Name)

	// Update secret
	prompt()
	fmt.Println("Updating secret...")
	secret = &v1.Secret{
		ObjectMeta: metav1.ObjectMeta{
			Name: "sdk",
		},
		Expires:     4071231000,
		Description: "test secret for sdk_update",
	}
	ret, err = secretsClient.Update(context.TODO(), secret, metav1.UpdateOptions{})
	if err != nil {
		panic(err.Error())
	}
	fmt.Printf("Updated secret..., new expires: %d\n", ret.Expires)

	// List secrets
	prompt()
	fmt.Println("Listing secrets...")
	list, err := secretsClient.List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		panic(err)
	}
	for _, d := range list.Items {
		fmt.Printf(" * %s (secretID: %s, secretKey: %s, expires: %d)\n",
			d.Name, d.SecretID, d.SecretKey, d.Expires)
	}
}

func prompt() {
	fmt.Printf("-> Press Return key to continue.")
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		break
	}
	if err := scanner.Err(); err != nil {
		panic(err)
	}
	fmt.Println()
}
