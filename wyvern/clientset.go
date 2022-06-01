package wyvern

import (
	"github.com/opsdata/elmt-sdk/rest"
	"github.com/opsdata/elmt-sdk/wyvern/service/elmt"
)

// Interface defines method used to return client interface:
// - 项目级别接口
type Interface interface {
	Elmt() elmt.ElmtInterface
}

// Clientset contains the clients for groups. Each group has exactly one
// version included in a Clientset.
type Clientset struct {
	elmt *elmt.ElmtClient
}

var _ Interface = &Clientset{}

// Elmt retrieves the ElmtClient.
func (c *Clientset) Elmt() elmt.ElmtInterface {
	return c.elmt
}

// NewForConfig creates a new Clientset for the given config.
// If config's RateLimiter is not set and QPS and Burst are acceptable,
// NewForConfig will generate a rate-limiter in configShallowCopy.
func NewForConfig(c *rest.Config) (*Clientset, error) {
	configShallowCopy := *c

	var (
		cs  Clientset
		err error
	)

	cs.elmt, err = elmt.NewForConfig(&configShallowCopy)
	if err != nil {
		return nil, err
	}
	return &cs, nil
}

// NewForConfigOrDie creates a new Clientset for the given config and
// panics if there is an error in the config.
func NewForConfigOrDie(c *rest.Config) *Clientset {
	var cs Clientset
	cs.elmt = elmt.NewForConfigOrDie(c)
	return &cs
}

// New creates a new Clientset for the given RESTClient.
func New(c rest.Interface) *Clientset {
	var cs Clientset
	cs.elmt = elmt.New(c)
	return &cs
}
