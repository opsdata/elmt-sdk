package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"path/filepath"

	"github.com/ory/ladon"

	metav1 "github.com/opsdata/common-base/pkg/meta/v1"
	authzclientv1 "github.com/opsdata/elmt-sdk/wyvern/service/elmt/authz/v1"

	"github.com/opsdata/elmt-sdk/tools/clientcmd"
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
	client, err := authzclientv1.NewForConfig(config)
	if err != nil {
		panic(err.Error())
	}

	request := &ladon.Request{
		Resource: "resources:articles:ladon-introduction",
		Action:   "delete",
		Subject:  "users:peter",
		Context: ladon.Context{
			"remoteIP": "192.168.0.5",
		},
	}

	// Authorize the request
	fmt.Println("Authorize request...")
	ret, err := client.Authz().Authorize(context.TODO(), request, metav1.AuthorizeOptions{})
	if err != nil {
		panic(err.Error())
	}

	fmt.Printf("Authorize response: %s.\n", ret.ToString())
}
