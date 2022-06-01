package v1

import (
	"github.com/opsdata/common-base/pkg/runtime"
	v1 "github.com/opsdata/elmt-api/apiserver/v1"
	"github.com/opsdata/elmt-sdk/rest"
)

// APIV1Interface has methods to work with elmt resources:
// - 服务级别接口
type APIV1Interface interface {
	RESTClient() rest.Interface
	SecretsGetter
	UsersGetter
	ZbxCmdGetter
}

// APIV1Client is used to interact with features provided by the group.
type APIV1Client struct {
	restClient rest.Interface
}

// Users create and return user rest client.
func (c *APIV1Client) Users() UserInterface {
	return newUsers(c)
}

// Secrets create and return secret rest client.
func (c *APIV1Client) Secrets() SecretInterface {
	return newSecrets(c)
}

// Secrets create and return secret rest client.
func (c *APIV1Client) ZbxCmd() ZbxCmdInterface {
	return newZbxCmd(c)
}

// NewForConfig creates a new APIV1Client for the given config.
func NewForConfig(c *rest.Config) (*APIV1Client, error) {
	config := *c
	setConfigDefaults(&config)

	client, err := rest.RESTClientFor(&config)
	if err != nil {
		return nil, err
	}
	return &APIV1Client{client}, nil
}

// NewForConfigOrDie creates a new APIV1Client for the given config and
// panics if there is an error in the config.
func NewForConfigOrDie(c *rest.Config) *APIV1Client {
	client, err := NewForConfig(c)
	if err != nil {
		panic(err)
	}
	return client
}

// New creates a new APIV1Client for the given RESTClient.
func New(c rest.Interface) *APIV1Client {
	return &APIV1Client{c}
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

// RESTClient returns a RESTClient that is used to communicate
// with API server by this client implementation.
func (c *APIV1Client) RESTClient() rest.Interface {
	if c == nil {
		return nil
	}
	return c.restClient
}