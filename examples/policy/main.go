package main

import (
	"bufio"
	"context"
	"flag"
	"fmt"
	"os"
	"path/filepath"

	"github.com/ory/ladon"

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

	policiesClient := elmtclient.APIV1().Policies()

	var policyConditions = ladon.Conditions{
		"owner": &ladon.EqualsSubjectCondition{},
	}

	policy := &v1.Policy{
		ObjectMeta: metav1.ObjectMeta{
			Name: "sdk",
		},
		Policy: v1.AuthzPolicy{
			DefaultPolicy: ladon.DefaultPolicy{
				Description: "description",
				Subjects:    []string{"user"},
				Effect:      ladon.AllowAccess,
				Resources:   []string{"articles:<[0-9]+>"},
				Actions:     []string{"create", "update"},
				Conditions:  policyConditions,
			},
		},
	}

	// Create policy
	fmt.Println("Creating policy...")
	ret, err := policiesClient.Create(context.TODO(), policy, metav1.CreateOptions{})
	if err != nil {
		panic(err.Error())
	}
	fmt.Printf("Created policy %q.\n", ret.Name)

	// Delete policy
	defer func() {
		fmt.Println("Deleting policy...")
		if err := policiesClient.Delete(context.TODO(), "sdk", metav1.DeleteOptions{Unscoped: true}); err != nil {
			fmt.Printf("Delete policy failed: %s\n", err.Error())
			return
		}
		fmt.Println("Deleted policy.")
	}()

	// Get policy
	prompt()
	fmt.Println("Geting policy...")
	ret, err = policiesClient.Get(context.TODO(), "sdk", metav1.GetOptions{})
	if err != nil {
		panic(err.Error())
	}
	fmt.Printf("Get policy %q.\n", ret.Name)

	// Update policy
	prompt()
	fmt.Println("Updating policy...")
	policy = &v1.Policy{
		ObjectMeta: metav1.ObjectMeta{
			Name: "sdk",
		},
		Policy: v1.AuthzPolicy{
			DefaultPolicy: ladon.DefaultPolicy{
				Description: "description - update",
				Subjects:    []string{"user"},
				Effect:      ladon.AllowAccess,
				Resources:   []string{"articles:<[0-9]+>"},
				Actions:     []string{"create", "update"},
				Conditions:  policyConditions,
			},
		},
	}
	ret, err = policiesClient.Update(context.TODO(), policy, metav1.UpdateOptions{})
	if err != nil {
		panic(err.Error())
	}
	fmt.Printf("Updated policy..., new policy: `%s`\n", ret.Policy.Description)

	// List policys
	prompt()
	fmt.Println("Listing policies...")
	list, err := policiesClient.List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		panic(err)
	}
	for _, d := range list.Items {
		fmt.Printf(" * %s (policy: `%s`)\n", d.Name, d.Policy.Description)
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
