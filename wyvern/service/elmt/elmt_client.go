package elmt

import (
	"github.com/opsdata/elmt-sdk/rest"

	apiv1 "github.com/opsdata/elmt-sdk/wyvern/service/elmt/apiserver/v1"
	authzv1 "github.com/opsdata/elmt-sdk/wyvern/service/elmt/authz/v1"
)

// ElmtInterface
// - hold the methods that elmt server-supported API services, versions and resources
type ElmtInterface interface {
	APIV1() apiv1.APIV1Interface
	AuthzV1() authzv1.AuthzV1Interface
}

// ElmtClient contains the clients for elmt service. Each elmt service has exactly one
// version included in an ElmtClient.
type ElmtClient struct {
	apiV1   *apiv1.APIV1Client
	authzV1 *authzv1.AuthzV1Client
}

// APIV1 retrieves the APIV1Client.
func (c *ElmtClient) APIV1() apiv1.APIV1Interface {
	return c.apiV1
}

// AuthzV1 retrieves the AuthzV1Client.
func (c *ElmtClient) AuthzV1() authzv1.AuthzV1Interface {
	return c.authzV1
}

// NewForConfig creates a new ElmtV1Client for the given config.
func NewForConfig(c *rest.Config) (*ElmtClient, error) {
	configShallowCopy := *c

	var (
		ec  ElmtClient
		err error
	)

	ec.apiV1, err = apiv1.NewForConfig(&configShallowCopy)
	if err != nil {
		return nil, err
	}

	ec.authzV1, err = authzv1.NewForConfig(&configShallowCopy)
	if err != nil {
		return nil, err
	}

	return &ec, nil
}

// NewForConfigOrDie creates a new ElmtClient for the given config and
// panics if there is an error in the config.
func NewForConfigOrDie(c *rest.Config) *ElmtClient {
	var ec ElmtClient
	ec.apiV1 = apiv1.NewForConfigOrDie(c)
	return &ec
}

// New creates a new ElmtClient for the given RESTClient.
func New(c rest.Interface) *ElmtClient {
	var ec ElmtClient
	ec.apiV1 = apiv1.New(c)
	return &ec
}
