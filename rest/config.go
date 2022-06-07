package rest

import (
	"crypto/tls"
	"crypto/x509"
	"encoding/base64"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	gruntime "runtime"
	"strings"
	"time"

	"github.com/opsdata/common-base/pkg/runtime"
	"github.com/opsdata/common-base/pkg/scheme"
	"github.com/opsdata/common-base/pkg/version"

	"github.com/opsdata/elmt-sdk/third_party/forked/gorequest"
)

// Config holds the common attributes that can be passed to a IAM client on
// initialization.
type Config struct {
	Host    string
	APIPath string
	ContentConfig

	// Server requires Basic authentication
	Username  string
	Password  string
	SecretID  string
	SecretKey string

	// Server requires Bearer authentication
	BearerToken string

	// Path to a file containing a BearerToken.
	BearerTokenFile string

	// TLSClientConfig contains settings to enable transport layer security
	TLSClientConfig

	// UserAgent is an optional field that specifies the caller of this request
	UserAgent string

	// The maximum length of time to wait before giving up on a server request, and a value of zero
	// means no timeout
	Timeout time.Duration

	MaxRetries    int
	RetryInterval time.Duration

	// JSON-RPC API information for Zabbix
	ZabbixApiUrl  string
	ZabbixApiUser string
	ZabbixApiPass string
}

// ContentConfig defines config for content.
type ContentConfig struct {
	ServiceName        string
	AcceptContentTypes string
	ContentType        string
	GroupVersion       *scheme.GroupVersion
	Negotiator         runtime.ClientNegotiator
}

type sanitizedConfig *Config

// GoString implements fmt.GoStringer and sanitizes(清洁/消毒) sensitive fields of Config
// to prevent accidental leaking via logs.
func (c *Config) GoString() string {
	return c.String()
}

// String implements fmt.Stringer and sanitizes sensitive fields of Config to
// prevent accidental leaking via logs.
func (c *Config) String() string {
	if c == nil {
		return "<nil>"
	}

	// Explicitly mark non-empty credential fields as redacted
	cc := sanitizedConfig(CopyConfig(c))

	if cc.Password != "" {
		cc.Password = "--- REDACTED ---"
	}

	if cc.BearerToken != "" {
		cc.BearerToken = "--- REDACTED ---"
	}

	if cc.SecretKey != "" {
		cc.SecretKey = "--- REDACTED ---"
	}

	return fmt.Sprintf("%#v", cc)
}

// RESTClientFor
// - 创建RESTClient客户端
// - return a RESTClient that satisfies the requested attributes on a client Config object
// - a RESTClient created by this method is generic: it expects to operate on an API that follows
// the ELMT conventions
func RESTClientFor(config *Config) (*RESTClient, error) {
	if config.GroupVersion == nil {
		return nil, fmt.Errorf("GroupVersion is required when initializing a RESTClient")
	}

	if config.Negotiator == nil {
		return nil, fmt.Errorf("NegotiatedSerializer is required when initializing a RESTClient")
	}

	// 生成基本的HTTP请求路径
	// - baseURL=http://127.0.0.1:8080
	// - versionedAPIPath=/v1
	baseURL, versionedAPIPath, err := defaultServerURLFor(config)
	if err != nil {
		return nil, err
	}

	tlsConfig, err := TLSConfigFor(config)
	if err != nil {
		return nil, err
	}

	// Only retry when get a server side error
	client := gorequest.New().TLSClientConfig(tlsConfig).Timeout(config.Timeout).
		Retry(config.MaxRetries, config.RetryInterval, http.StatusInternalServerError)

	// NOTICE: must set DoNotClearSuperAgent to true, or the client will clean header befor http.Do
	client.DoNotClearSuperAgent = true

	var gv scheme.GroupVersion

	if config.GroupVersion != nil {
		gv = *config.GroupVersion
	}

	clientContent := ClientContentConfig{
		Username:           config.Username,
		Password:           config.Password,
		SecretID:           config.SecretID,
		SecretKey:          config.SecretKey,
		BearerToken:        config.BearerToken,
		BearerTokenFile:    config.BearerTokenFile,
		TLSClientConfig:    config.TLSClientConfig,
		AcceptContentTypes: config.AcceptContentTypes,
		ContentType:        config.ContentType,
		GroupVersion:       gv,
		Negotiator:         config.Negotiator,
	}

	return NewRESTClient(baseURL, versionedAPIPath, clientContent, client)
}

// TLSConfigFor
// - return a tls.Config that will provide the transport level security defined
// by the provided Config
// - it will return nil if no transport level security is requested
func TLSConfigFor(c *Config) (*tls.Config, error) {
	if !(c.HasCA() || c.HasCertAuth() || c.Insecure || len(c.ServerName) > 0) {
		return nil, nil
	}

	if c.HasCA() && c.Insecure {
		return nil, fmt.Errorf("specifying a root certificates file with the insecure flag is not allowed")
	}

	if err := LoadTLSFiles(c); err != nil {
		return nil, err
	}

	tlsConfig := &tls.Config{
		// Can't use SSLv3 because of POODLE and BEAST
		// Can't use TLSv1.0 because of POODLE and BEAST using CBC cipher
		// Can't use TLSv1.1 because of RC4 cipher usage
		MinVersion: tls.VersionTLS12,
		//nolint: gosec
		InsecureSkipVerify: c.Insecure,
		ServerName:         c.ServerName,
		NextProtos:         c.NextProtos,
	}

	if c.HasCA() {
		tlsConfig.RootCAs = rootCertPool(c.CAData)
	}

	var staticCert *tls.Certificate

	// Treat cert as static if either key or cert was data, not a file
	if c.HasCertAuth() {
		// If key/cert were provided, verify them before setting up
		// tlsConfig.GetClientCertificate.
		cert, err := tls.X509KeyPair(c.CertData, c.KeyData)
		if err != nil {
			return nil, err
		}

		staticCert = &cert
	}

	if c.HasCertAuth() {
		tlsConfig.GetClientCertificate = func(*tls.CertificateRequestInfo) (*tls.Certificate, error) {
			// Note: static key/cert data always take precedence over cert
			// callback.
			if staticCert != nil {
				return staticCert, nil
			}

			// Both c.TLS.CertData/KeyData were unset and GetCert didn't return
			// anything. Return an empty tls.Certificate, no client cert will
			// be sent to the server.
			return &tls.Certificate{}, nil
		}
	}

	return tlsConfig, nil
}

// rootCertPool
// - return nil if caData is empty: when passed along, this will mean "use system CAs"
// - when caData is not empty, it will be the ONLY information used in the CertPool
func rootCertPool(caData []byte) *x509.CertPool {
	if len(caData) == 0 {
		return nil
	}

	// If we have caData, use it
	certPool := x509.NewCertPool()
	certPool.AppendCertsFromPEM(caData)

	return certPool
}

// LoadTLSFiles
// - copy the data from the CertFile, KeyFile, and CAFile fields into the CertData,
// KeyData, and CAFile fields, or returns an error
// - if no error is returned, all three fields are either populated or were empty to start
func LoadTLSFiles(c *Config) error {
	var err error

	c.CAData, err = dataFromSliceOrFile(c.CAData, c.CAFile)
	if err != nil {
		return err
	}

	c.CertData, err = dataFromSliceOrFile(c.CertData, c.CertFile)
	if err != nil {
		return err
	}

	c.KeyData, err = dataFromSliceOrFile(c.KeyData, c.KeyFile)
	if err != nil {
		return err
	}

	return nil
}

// dataFromSliceOrFile
// - return data from the slice (if non-empty), or from the file
// - or an error if an error occurred reading the file.
func dataFromSliceOrFile(data []byte, file string) ([]byte, error) {
	if len(data) > 0 {
		return base64.StdEncoding.DecodeString(string(data))
	}

	if len(file) > 0 {
		fileData, err := ioutil.ReadFile(file)
		if err != nil {
			return []byte{}, err
		}

		return fileData, nil
	}

	return nil, nil
}

// SetELMTDefaults
// - set default values on the provided client config for accessing the
// ELMT API
// - or returns an error if any of the defaults are impossible or invalid
func SetELMTDefaults(config *Config) error {
	if len(config.UserAgent) == 0 {
		config.UserAgent = DefaultUserAgent()
	}

	return nil
}

// DefaultUserAgent returns a User-Agent string built from static global vars.
func DefaultUserAgent() string {
	return buildUserAgent(
		adjustCommand(os.Args[0]),
		adjustVersion(version.Get().GitVersion),
		gruntime.GOOS,
		gruntime.GOARCH,
		adjustCommit(version.Get().GitCommit))
}

// buildUserAgent builds a User-Agent string from given args.
func buildUserAgent(command, version, os, arch, commit string) string {
	return fmt.Sprintf(
		"%s/%s (%s/%s) elmt/%s", command, version, os, arch, commit)
}

// adjustCommand returns the last component of the OS-specific command path for use in User-Agent.
func adjustCommand(p string) string {
	// Unlikely, but better than returning "".
	if len(p) == 0 {
		return "unknown"
	}

	return filepath.Base(p)
}

// adjustVersion strips "alpha", "beta", etc. from version in form
// major.minor.patch-[alpha|beta|etc].
func adjustVersion(v string) string {
	if len(v) == 0 {
		return "unknown"
	}

	seg := strings.SplitN(v, "-", 2)

	return seg[0]
}

// adjustCommit returns sufficient significant figures of the commit's git hash.
func adjustCommit(c string) string {
	if len(c) == 0 {
		return "unknown"
	}

	if len(c) > 7 {
		return c[:7]
	}

	return c
}

// AddUserAgent adds a http User-Agent header.
func AddUserAgent(config *Config, userAgent string) *Config {
	fullUserAgent := DefaultUserAgent() + "/" + userAgent
	config.UserAgent = fullUserAgent

	return config
}

// CopyConfig returns a copy of the given config.
func CopyConfig(config *Config) *Config {
	return &Config{
		Host:            config.Host,
		APIPath:         config.APIPath,
		ContentConfig:   config.ContentConfig,
		Username:        config.Username,
		Password:        config.Password,
		SecretID:        config.SecretID,
		SecretKey:       config.SecretKey,
		BearerToken:     config.BearerToken,
		BearerTokenFile: config.BearerTokenFile,
		UserAgent:       config.UserAgent,
		Timeout:         config.Timeout,

		TLSClientConfig: TLSClientConfig{
			Insecure:   config.TLSClientConfig.Insecure,
			ServerName: config.TLSClientConfig.ServerName,
			CertFile:   config.TLSClientConfig.CertFile,
			KeyFile:    config.TLSClientConfig.KeyFile,
			CAFile:     config.TLSClientConfig.CAFile,
			CertData:   config.TLSClientConfig.CertData,
			KeyData:    config.TLSClientConfig.KeyData,
			CAData:     config.TLSClientConfig.CAData,
			NextProtos: config.TLSClientConfig.NextProtos,
		},
	}
}
