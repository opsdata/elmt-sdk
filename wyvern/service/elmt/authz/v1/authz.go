package v1

import (
	"context"

	"github.com/ory/ladon"

	metav1 "github.com/opsdata/common-base/pkg/meta/v1"
	authzv1 "github.com/opsdata/elmt-api/authz/v1"
	rest "github.com/opsdata/elmt-sdk/rest"
)

// AuthzInterface interface
// - methods to work with Authz resources
//
type AuthzInterface interface {
	Authorize(ctx context.Context, request *ladon.Request, opts metav1.AuthorizeOptions) (*authzv1.Response, error)
}

type authz struct {
	client rest.Interface
}

func newAuthz(c *AuthzV1Client) *authz {
	return &authz{
		client: c.RESTClient(),
	}
}

func (c *authz) Authorize(ctx context.Context, request *ladon.Request, opts metav1.AuthorizeOptions) (result *authzv1.Response, err error) {
	result = &authzv1.Response{}

	err = c.client.Post().
		Resource("authz").
		VersionedParams(opts).
		Body(request).
		Do(ctx).
		Into(result)

	return
}
