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

	usersClient := elmtclient.APIV1().Users()

	user := &v1.User{
		ObjectMeta: metav1.ObjectMeta{
			Name: "sdk",
		},
		Password: "sdk@2022",
	}

	// Create user
	fmt.Println("Creating user...")
	ret, err := usersClient.Create(context.TODO(), user, metav1.CreateOptions{})
	if err != nil {
		panic(err.Error())
	}
	fmt.Printf("Created user %q.\n", ret.Name)

	// Delete secret
	defer func() {
		fmt.Println("Deleting user...")
		if err := usersClient.Delete(context.TODO(), "sdk", metav1.DeleteOptions{}); err != nil {
			fmt.Printf("Delete user failed: %s\n", err.Error())
			return
		}
		fmt.Println("Deleted user.")
	}()

	// Get user
	prompt()
	fmt.Println("Geting user...")
	ret, err = usersClient.Get(context.TODO(), user.Name, metav1.GetOptions{})
	if err != nil {
		panic(err.Error())
	}

	fmt.Printf("Get user %q.\n", ret.Name)

	// Update user
	prompt()
	fmt.Println("Updating user...")
	user = &v1.User{
		ObjectMeta: metav1.ObjectMeta{
			Name: "sdk",
		},
	}
	ret, err = usersClient.Update(context.TODO(), user, metav1.UpdateOptions{})
	if err != nil {
		panic(err.Error())
	}
	fmt.Println("Updated user...")
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
