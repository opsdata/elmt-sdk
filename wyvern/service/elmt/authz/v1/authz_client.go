package v1

import (
	v1 "github.com/opsdata/elmt-api/authz/v1"

	"github.com/opsdata/common-base/pkg/runtime"
	"github.com/opsdata/elmt-sdk/rest"
)

// AuthzV1Interface interface:
// - methods to work with ELMT resources
//
type AuthzV1Interface interface {
	RESTClient() rest.Interface
	Authz() AuthzInterface
}

/*
 * AuthzV1Client:
 * - be used to interact with features provided by the group (authz)
 * - implement the AuthzV1Interface interface
 */

type AuthzV1Client struct {
	restClient rest.Interface
}

func (c *AuthzV1Client) Authz() AuthzInterface {
	return newAuthz(c)
}

func (c *AuthzV1Client) RESTClient() rest.Interface {
	if c == nil {
		return nil
	}

	return c.restClient
}

/*
 * Methods to initiate a new AuthzV1Client:
 * - New
 * - NewForConfig, NewForConfigOrDie
 */

func New(c rest.Interface) *AuthzV1Client {
	return &AuthzV1Client{c}
}

func NewForConfig(c *rest.Config) (*AuthzV1Client, error) {
	config := *c
	setConfigDefaults(&config)

	client, err := rest.RESTClientFor(&config)
	if err != nil {
		return nil, err
	}

	return &AuthzV1Client{client}, nil
}

func NewForConfigOrDie(c *rest.Config) *AuthzV1Client {
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
