package v1

import (
	"context"
	"time"

	metav1 "github.com/opsdata/common-base/pkg/meta/v1"
	v1 "github.com/opsdata/elmt-api/apiserver/v1"
	rest "github.com/opsdata/elmt-sdk/rest"
)

// PoliciesGetter
// - method to return a PolicyInterface.
// - a group's client should implement this interface
//
type PoliciesGetter interface {
	Policies() PolicyInterface
}

type PolicyInterface interface {
	Create(ctx context.Context, policy *v1.Policy, opts metav1.CreateOptions) (*v1.Policy, error)
	Update(ctx context.Context, policy *v1.Policy, opts metav1.UpdateOptions) (*v1.Policy, error)
	Delete(ctx context.Context, name string, opts metav1.DeleteOptions) error
	DeleteCollection(ctx context.Context, opts metav1.DeleteOptions, listOpts metav1.ListOptions) error
	Get(ctx context.Context, name string, opts metav1.GetOptions) (*v1.Policy, error)
	List(ctx context.Context, opts metav1.ListOptions) (*v1.PolicyList, error)
}

type policies struct {
	client rest.Interface
}

func newPolicies(c *APIV1Client) *policies {
	return &policies{
		client: c.RESTClient(),
	}
}

func (c *policies) Get(ctx context.Context, name string, options metav1.GetOptions) (result *v1.Policy, err error) {
	result = &v1.Policy{}
	err = c.client.Get().
		Resource("policies").
		Name(name).
		VersionedParams(options).
		Do(ctx).
		Into(result)

	return
}

func (c *policies) List(ctx context.Context, opts metav1.ListOptions) (result *v1.PolicyList, err error) {
	var timeout time.Duration
	if opts.TimeoutSeconds != nil {
		timeout = time.Duration(*opts.TimeoutSeconds) * time.Second
	}

	result = &v1.PolicyList{}
	err = c.client.Get().
		Resource("policies").
		VersionedParams(opts).
		Timeout(timeout).
		Do(ctx).
		Into(result)

	return
}

func (c *policies) Create(ctx context.Context, policy *v1.Policy,
	opts metav1.CreateOptions) (result *v1.Policy, err error) {
	result = &v1.Policy{}
	err = c.client.Post().
		Resource("policies").
		VersionedParams(opts).
		Body(policy).
		Do(ctx).
		Into(result)

	return
}

func (c *policies) Update(ctx context.Context, policy *v1.Policy,
	opts metav1.UpdateOptions) (result *v1.Policy, err error) {
	result = &v1.Policy{}
	err = c.client.Put().
		Resource("policies").
		Name(policy.Name).
		VersionedParams(opts).
		Body(policy).
		Do(ctx).
		Into(result)

	return
}

func (c *policies) Delete(ctx context.Context, name string, opts metav1.DeleteOptions) error {
	return c.client.Delete().
		Resource("policies").
		Name(name).
		Body(&opts).
		Do(ctx).
		Error()
}

func (c *policies) DeleteCollection(ctx context.Context, opts metav1.DeleteOptions, listOpts metav1.ListOptions) error {
	var timeout time.Duration
	if listOpts.TimeoutSeconds != nil {
		timeout = time.Duration(*listOpts.TimeoutSeconds) * time.Second
	}

	return c.client.Delete().
		Resource("policies").
		VersionedParams(listOpts).
		Timeout(timeout).
		Body(&opts).
		Do(ctx).
		Error()
}
