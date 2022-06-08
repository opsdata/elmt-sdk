package v1

import (
	"context"

	metav1 "github.com/opsdata/common-base/pkg/meta/v1"
	v1 "github.com/opsdata/elmt-api/apiserver/v1"
	rest "github.com/opsdata/elmt-sdk/rest"
)

// ZbxCmdGetter
// - method to return a ZbxCmdInterface
// - a group's client should implement this interface
//
type ZbxCmdGetter interface {
	ZbxCmd() ZbxCmdInterface
}

type ZbxCmdInterface interface {
	GetZbxItem(ctx context.Context, item_name string, opts metav1.GetOptions) (*v1.Indicator, error)
	GetZbxHost(ctx context.Context, host_name string, opts metav1.GetOptions) (*v1.ZbxHost, error)
}

type zbxcmd struct {
	client rest.Interface
}

func newZbxCmd(c *APIV1Client) *zbxcmd {
	return &zbxcmd{
		client: c.RESTClient(),
	}
}

func (z *zbxcmd) GetZbxItem(ctx context.Context, item_name string, options metav1.GetOptions) (result *v1.Indicator, err error) {
	result = &v1.Indicator{}

	err = z.client.Get().
		Resource("zbxitems").
		Name(item_name).
		VersionedParams(options).
		Do(ctx).
		Into(result)

	return
}

func (z *zbxcmd) GetZbxHost(ctx context.Context, host_name string, options metav1.GetOptions) (result *v1.ZbxHost, err error) {
	result = &v1.ZbxHost{}

	err = z.client.Get().
		Resource("zbxhosts").
		Name(host_name).
		VersionedParams(options).
		Do(ctx).
		Into(result)

	return
}
