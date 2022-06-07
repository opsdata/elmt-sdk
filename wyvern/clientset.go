package wyvern

//
// wyvern: 虚拟组织名称
//

import (
	"github.com/opsdata/elmt-sdk/rest"
	"github.com/opsdata/elmt-sdk/wyvern/service/elmt"
)

type Interface interface {
	Elmt() elmt.ElmtInterface
}

type Clientset struct {
	elmt *elmt.ElmtClient
}

var _ Interface = &Clientset{}

func (c *Clientset) Elmt() elmt.ElmtInterface {
	return c.elmt
}

func NewForConfig(c *rest.Config) (*Clientset, error) {
	var (
		cs  Clientset
		err error
	)

	cs.elmt, err = elmt.NewForConfig(c)
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
