package rest

import (
	"net/url"
	"strings"

	"github.com/opsdata/common-base/pkg/runtime"
	"github.com/opsdata/common-base/pkg/scheme"

	"github.com/opsdata/elmt-sdk/third_party/forked/gorequest"
)

// Interface captures the set of operations for generically interacting with ELMT REST apis.
type Interface interface {
	Verb(verb string) *Request
	Post() *Request
	Put() *Request
	Get() *Request
	Delete() *Request
	APIVersion() scheme.GroupVersion
}

// TLSConfig holds the information needed to set up a TLS transport.
type TLSConfig struct {
	CAFile         string // Path of the PEM-encoded server trusted root certificates.
	CertFile       string // Path of the PEM-encoded client certificate.
	KeyFile        string // Path of the PEM-encoded client key.
	ReloadTLSFiles bool   // Set to indicate that the original config provided files, and that they should be reloaded
	Insecure       bool   // Server should be accessed without verifying the certificate. For testing only.
	ServerName     string // Override for the server name passed to the server for SNI and used to verify certificates.
	CAData         []byte // Bytes of the PEM-encoded server trusted root certificates. Supercedes CAFile.
	CertData       []byte // Bytes of the PEM-encoded client certificate. Supercedes CertFile.
	KeyData        []byte // Bytes of the PEM-encoded client key. Supercedes KeyFile.
}

// ClientContentConfig controls how RESTClient communicates with the server.
type ClientContentConfig struct {
	Username     string
	Password     string
	SecretID     string
	SecretKey    string
	GroupVersion scheme.GroupVersion
	Negotiator   runtime.ClientNegotiator

	// Server requires Bearer authentication.
	BearerToken string

	// Path to a file containing a BearerToken.
	// If set, the contents are periodically read.
	// The last successfully read value takes precedence over BearerToken.
	BearerTokenFile string

	// AcceptContentTypes specifies the types the client will accept and is optional.
	// If not set, ContentType will be used to define the Accept header.
	AcceptContentTypes string

	// ContentType specifies the wire format used to communicate with the server.
	// This value will be set as the Accept header on requests made to the server if
	// AcceptContentTypes is not set, and as the default content type on any object
	// sent to the server. If not set, "application/json" is used.
	ContentType string

	TLSClientConfig
}

// HasBasicAuth returns whether the configuration has basic authentication or not.
func (c *ClientContentConfig) HasBasicAuth() bool {
	return len(c.Username) != 0
}

// HasTokenAuth returns whether the configuration has token authentication or not.
func (c *ClientContentConfig) HasTokenAuth() bool {
	return len(c.BearerToken) != 0 || len(c.BearerTokenFile) != 0
}

// HasKeyAuth returns whether the configuration has secretId/secretKey authentication or not.
func (c *ClientContentConfig) HasKeyAuth() bool {
	return len(c.SecretID) != 0 && len(c.SecretKey) != 0
}

// RESTClient
// - impose common ELMT API conventions on a set of resource paths
type RESTClient struct {
	// the root URL for all invocations of the client, such as http://elmt.api.opsdata.cn:8080, etc.
	base *url.URL

	// a path segment connecting the base URL to the resource root, such as /v1, /v2, etc.
	versionedAPIPath string

	// stand for the client group, eg: elmt.api, elmt.authz
	group string

	// describe how a RESTClient encodes and decodes responses
	content ClientContentConfig

	Client *gorequest.SuperAgent
}

// NewRESTClient
// - create a new RESTClient
// - the client performs generic REST functions such as Get, Put, Post, and Delete on specified paths.
func NewRESTClient(baseURL *url.URL, versionedAPIPath string,
	config ClientContentConfig, client *gorequest.SuperAgent) (*RESTClient, error) {
	if len(config.ContentType) == 0 {
		config.ContentType = "application/json"
	}

	base := *baseURL
	if !strings.HasSuffix(base.Path, "/") {
		base.Path += "/"
	}

	base.RawQuery = ""
	base.Fragment = ""

	return &RESTClient{
		base:             &base,
		group:            config.GroupVersion.Group,
		versionedAPIPath: versionedAPIPath,
		content:          config,
		Client:           client,
	}, nil
}

// Verb begins a Verb request.
func (c *RESTClient) Verb(verb string) *Request {
	return NewRequest(c).Verb(verb)
}

// Post begins a POST request. Short for c.Verb("POST").
func (c *RESTClient) Post() *Request {
	return c.Verb("POST")
}

// Put begins a PUT request. Short for c.Verb("PUT").
func (c *RESTClient) Put() *Request {
	return c.Verb("PUT")
}

// Get begins a GET request. Short for c.Verb("GET").
func (c *RESTClient) Get() *Request {
	return c.Verb("GET")
}

// Delete begins a DELETE request. Short for c.Verb("DELETE").
func (c *RESTClient) Delete() *Request {
	return c.Verb("DELETE")
}

// APIVersion returns the APIVersion this RESTClient is expected to use.
func (c *RESTClient) APIVersion() scheme.GroupVersion {
	return c.content.GroupVersion
}
