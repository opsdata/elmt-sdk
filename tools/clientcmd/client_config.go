package clientcmd

import (
	"net/url"
	"time"

	restclient "github.com/opsdata/elmt-sdk/rest"
)

// Server contains information about how to communicate with the elmt api server.
type Server struct {
	LocationOfOrigin string
	Timeout          time.Duration `yaml:"timeout,omitempty"        mapstructure:"timeout,omitempty"`
	MaxRetries       int           `yaml:"max-retries,omitempty"    mapstructure:"max-retries,omitempty"`
	RetryInterval    time.Duration `yaml:"retry-interval,omitempty" mapstructure:"retry-interval,omitempty"`
	Address          string        `yaml:"address,omitempty"        mapstructure:"address,omitempty"`

	// TLSServerName is used to check server certificate.
	// If TLSServerName is empty, the hostname used to contact the server is used.
	// +optional
	TLSServerName string `yaml:"tls-server-name,omitempty" mapstructure:"tls-server-name,omitempty"`

	// InsecureSkipTLSVerify skips the validity check for the server's certificate.
	// This will make your HTTPS connections insecure.
	// +optional
	InsecureSkipTLSVerify bool `yaml:"insecure-skip-tls-verify,omitempty" mapstructure:"insecure-skip-tls-verify,omitempty"`

	// CertificateAuthority is the path to a cert file for the certificate authority.
	// +optional
	CertificateAuthority string `yaml:"certificate-authority,omitempty" mapstructure:"certificate-authority,omitempty"`

	// CertificateAuthorityData contains PEM-encoded certificate authority certificates.
	// Overrides CertificateAuthority
	// +optional
	CertificateAuthorityData string `yaml:"certificate-authority-data,omitempty" mapstructure:"certificate-authority-data,omitempty"`
}

// AuthInfo contains information that describes identity information.
type AuthInfo struct {
	LocationOfOrigin string

	Username  string `yaml:"username,omitempty" mapstructure:"username,omitempty"`
	Password  string `yaml:"password,omitempty" mapstructure:"password,omitempty"`
	SecretID  string `yaml:"secret-id,omitempty"  mapstructure:"secret-id,omitempty"`
	SecretKey string `yaml:"secret-key,omitempty" mapstructure:"secret-key,omitempty"`

	// Token is the bearer token for authentication to the elmt cluster.
	// +optional
	Token string `yaml:"token,omitempty" mapstructure:"token,omitempty"`

	// ClientCertificate is the path to a client certificate file for TLS.
	// +optional
	ClientCertificate string `yaml:"client-certificate,omitempty" mapstructure:"client-certificate,omitempty"`

	// ClientCertificateData contains PEM-encoded data from a client cert file for TLS. Overrides ClientCertificate.
	// +optional
	ClientCertificateData string `yaml:"client-certificate-data,omitempty" mapstructure:"client-certificate-data,omitempty"`

	// ClientKey is the path to a client key file for TLS.
	// +optional
	ClientKey string `yaml:"client-key,omitempty" mapstructure:"client-key,omitempty"`

	// ClientKeyData contains PEM-encoded data from a client key file for TLS. It overrides ClientKey.
	// +optional
	ClientKeyData string `yaml:"client-key-data,omitempty" mapstructure:"client-key-data,omitempty"`
}

// ZabbixInfo contains information that describes Zabbix JSON-RPC API information.
type ZabbixInfo struct {
	ApiUrl  string `yaml:"api-url,omitempty" mapstructure:"api-url,omitempty"`
	ApiUser string `yaml:"api-user,omitempty" mapstructure:"api-user,omitempty"`
	ApiPass string `yaml:"api-pass,omitempty" mapstructure:"api-pass,omitempty"`
}

// Config defines a config struct used by sdk.
type Config struct {
	APIVersion string      `yaml:"apiVersion,omitempty" mapstructure:"apiVersion,omitempty"`
	Server     *Server     `yaml:"server,omitempty"     mapstructure:"server,omitempty"`
	AuthInfo   *AuthInfo   `yaml:"user,omitempty"       mapstructure:"user,omitempty"`
	ZabbixInfo *ZabbixInfo `yaml:"zabbix,omitempty"     mapstructure:"zabbix,omitempty"`
}

// NewConfig is a convenience function that returns a new Config object with non-nil maps.
func NewConfig() *Config {
	return &Config{
		Server:     &Server{},
		AuthInfo:   &AuthInfo{},
		ZabbixInfo: &ZabbixInfo{},
	}
}

// ClientConfig interface
// - be used to make it easy to get an API server client
// - it returns a complete client config
type ClientConfig interface {
	ClientConfig() (*restclient.Config, error)
}

type DirectClientConfig struct {
	config Config
}

// getServer returns the clientcmdapi.Server, or an error if a required server is not found.
func (config *DirectClientConfig) getServer() Server {
	return *config.config.Server
}

// getAuthInfo returns the clientcmdapi.AuthInfo, or an error if a required auth info is not found.
func (config *DirectClientConfig) getAuthInfo() AuthInfo {
	return *config.config.AuthInfo
}

// getZabbixInfo returns the clientcmdapi.ZabbixInfo, or an error if a required zabbix info is not found.
func (config *DirectClientConfig) getZabbixInfo() ZabbixInfo {
	return *config.config.ZabbixInfo
}

// ConfirmUsable
// - look a particular context and determine if that particular part of the config is useable
// - there might still be errors in the config, but no errors in the sections requested or referenced
// - it does not return early so that it can find as many errors as possible.
func (config *DirectClientConfig) ConfirmUsable() error {
	validationErrors := make([]error, 0)

	authInfo := config.getAuthInfo()
	validationErrors = append(validationErrors, validateAuthInfo(authInfo)...)

	server := config.getServer()
	validationErrors = append(validationErrors, validateServerInfo(server)...)

	// when direct client config is specified, and the only error is that no server is defined, we should
	// return a standard "no config" error
	if len(validationErrors) == 1 && validationErrors[0] == ErrEmptyServer {
		return newErrConfigurationInvalid([]error{ErrEmptyConfig})
	}

	return newErrConfigurationInvalid(validationErrors)
}

// ClientConfig implements ClientConfig interface.
func (config *DirectClientConfig) ClientConfig() (*restclient.Config, error) {
	user := config.getAuthInfo()
	server := config.getServer()
	zabbix := config.getZabbixInfo()

	if err := config.ConfirmUsable(); err != nil {
		return nil, err
	}

	clientConfig := &restclient.Config{
		BearerToken:   user.Token,
		Username:      user.Username,
		Password:      user.Password,
		SecretID:      user.SecretID,
		SecretKey:     user.SecretKey,
		Host:          server.Address,
		Timeout:       server.Timeout,
		MaxRetries:    server.MaxRetries,
		RetryInterval: server.RetryInterval,

		// TLS
		TLSClientConfig: restclient.TLSClientConfig{
			Insecure:   server.InsecureSkipTLSVerify,
			ServerName: server.TLSServerName,
			CertFile:   user.ClientCertificate,
			KeyFile:    user.ClientKey,
			CertData:   []byte(user.ClientCertificateData),
			KeyData:    []byte(user.ClientKeyData),
			CAFile:     server.CertificateAuthority,
			CAData:     []byte(server.CertificateAuthorityData),
		},

		// Zabbix JSON-RPC
		ZabbixApiUrl:  zabbix.ApiUrl,
		ZabbixApiUser: zabbix.ApiUser,
		ZabbixApiPass: zabbix.ApiPass,
	}

	if u, err := url.ParseRequestURI(clientConfig.Host); err == nil && u.Opaque == "" && len(u.Path) > 1 {
		u.RawQuery = ""
		u.Fragment = ""
		clientConfig.Host = u.String()
	}

	return clientConfig, nil
}

func NewClientConfigFromConfig(config *Config) ClientConfig {
	return &DirectClientConfig{*config}
}

func NewClientConfigFromBytes(configBytes []byte) (ClientConfig, error) {
	config, err := Load(configBytes)
	if err != nil {
		return nil, err
	}
	return &DirectClientConfig{*config}, nil
}

func RESTConfigFromELMTConfig(configBytes []byte) (*restclient.Config, error) {
	clientConfig, err := NewClientConfigFromBytes(configBytes)
	if err != nil {
		return nil, err
	}
	return clientConfig.ClientConfig()
}

// BuildConfigFromFlags
// - a helper function that builds configs from a server url and .elmtconfig filepath
//
// - 示例: config, err := clientcmd.BuildConfigFromFlags("", "/root/.elmt/config.yaml")
//
func BuildConfigFromFlags(serverURL, elmtconfigPath string) (*restclient.Config, error) {
	config, err := LoadFromFile(elmtconfigPath)
	if err != nil {
		return nil, err
	}

	if len(serverURL) > 0 {
		config.Server.Address = serverURL
	}

	directClientConfig := &DirectClientConfig{*config}
	return directClientConfig.ClientConfig()
}
