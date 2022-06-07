package v1

import (
	"context"
	"time"

	metav1 "github.com/opsdata/common-base/pkg/meta/v1"
	v1 "github.com/opsdata/elmt-api/apiserver/v1"
	rest "github.com/opsdata/elmt-sdk/rest"
)

// UsersGetter
// - method to return a UserInterface.
// - a group's client should implement this interface
//
type UsersGetter interface {
	Users() UserInterface
}

type UserInterface interface {
	Create(ctx context.Context, user *v1.User, opts metav1.CreateOptions) (*v1.User, error)
	Update(ctx context.Context, user *v1.User, opts metav1.UpdateOptions) (*v1.User, error)
	Delete(ctx context.Context, name string, opts metav1.DeleteOptions) error
	DeleteCollection(ctx context.Context, opts metav1.DeleteOptions, listOpts metav1.ListOptions) error
	Get(ctx context.Context, name string, opts metav1.GetOptions) (*v1.User, error)
	List(ctx context.Context, opts metav1.ListOptions) (*v1.UserList, error)
}

type users struct {
	client rest.Interface
}

func newUsers(c *APIV1Client) *users {
	return &users{
		client: c.RESTClient(),
	}
}

func (c *users) Get(ctx context.Context, name string, options metav1.GetOptions) (result *v1.User, err error) {
	result = &v1.User{}

	err = c.client.Get().
		Resource("users").
		Name(name).
		VersionedParams(options).
		Do(ctx).
		Into(result)

	return
}

func (c *users) List(ctx context.Context, opts metav1.ListOptions) (result *v1.UserList, err error) {
	var timeout time.Duration
	if opts.TimeoutSeconds != nil {
		timeout = time.Duration(*opts.TimeoutSeconds) * time.Second
	}

	result = &v1.UserList{}

	err = c.client.Get().
		Resource("users").
		VersionedParams(opts).
		Timeout(timeout).
		Do(ctx).
		Into(result)

	return
}

func (c *users) Create(ctx context.Context, user *v1.User, opts metav1.CreateOptions) (result *v1.User, err error) {
	result = &v1.User{}

	err = c.client.Post().
		Resource("users").
		VersionedParams(opts).
		Body(user).
		Do(ctx).
		Into(result)

	return
}

func (c *users) Update(ctx context.Context, user *v1.User, opts metav1.UpdateOptions) (result *v1.User, err error) {
	result = &v1.User{}

	err = c.client.Put().
		Resource("users").
		Name(user.Name).
		VersionedParams(opts).
		Body(user).
		Do(ctx).
		Into(result)

	return
}

func (c *users) Delete(ctx context.Context, name string, opts metav1.DeleteOptions) error {
	return c.client.Delete().
		Resource("users").
		Name(name).
		Body(&opts).
		Do(ctx).
		Error()
}

func (c *users) DeleteCollection(ctx context.Context, opts metav1.DeleteOptions, listOpts metav1.ListOptions) error {
	var timeout time.Duration
	if listOpts.TimeoutSeconds != nil {
		timeout = time.Duration(*listOpts.TimeoutSeconds) * time.Second
	}

	return c.client.Delete().
		Resource("users").
		VersionedParams(listOpts).
		Timeout(timeout).
		Body(&opts).
		Do(ctx).
		Error()
}
