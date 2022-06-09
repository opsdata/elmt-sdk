package v1

import (
	v1 "github.com/opsdata/elmt-api/apiserver/v1"

	"github.com/opsdata/common-base/pkg/runtime"
	"github.com/opsdata/elmt-sdk/rest"
)

// APIV1Interface
// - methods to work with ELMT resources
//
type APIV1Interface interface {
	RESTClient() rest.Interface
	UsersGetter
	SecretsGetter
	PoliciesGetter
	ZbxCmdGetter
}

/*
 * APIV1Client:
 * - be used to interact with features provided by the group (apiserver)
 * - implement the APIV1Interface interface
 */

type APIV1Client struct {
	restClient rest.Interface
}

func (c *APIV1Client) RESTClient() rest.Interface {
	if c == nil {
		return nil
	}
	return c.restClient
}

func (c *APIV1Client) Users() UserInterface {
	return newUsers(c)
}

func (c *APIV1Client) Secrets() SecretInterface {
	return newSecrets(c)
}

func (c *APIV1Client) Policies() PolicyInterface {
	return newPolicies(c)
}

func (c *APIV1Client) ZbxCmd() ZbxCmdInterface {
	return newZbxCmd(c)
}

/*
 * Methods to initiate a new APIV1Client:
 * - New
 * - NewForConfig, NewForConfigOrDie
 */

func New(c rest.Interface) *APIV1Client {
	return &APIV1Client{c}
}

func NewForConfig(c *rest.Config) (*APIV1Client, error) {
	config := *c
	setConfigDefaults(&config)

	client, err := rest.RESTClientFor(&config)
	if err != nil {
		return nil, err
	}

	return &APIV1Client{client}, nil
}

func NewForConfigOrDie(c *rest.Config) *APIV1Client {
	client, err := NewForConfig(c)
	if err != nil {
		panic(err)
	}
	return client
}

func setConfigDefaults(config *rest.Config) {
	gv := v1.SchemeGroupVersion
	config.GroupVersion = &gv
	config.APIPath = ""
	config.Negotiator = runtime.NewSimpleClientNegotiator()

	if config.UserAgent == "" {
		config.UserAgent = rest.DefaultUserAgent()
	}
}
