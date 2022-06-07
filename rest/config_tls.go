package rest

import "fmt"

// TLSClientConfig contains settings to enable transport layer security.
type TLSClientConfig struct {
	// Server should be accessed without verifying the TLS certificate. For testing only.
	Insecure bool

	// ServerName is passed to the server for SNI and is used in the client to check server
	// ceritificates against. If ServerName is empty, the hostname used to contact the
	// server is used.
	ServerName string

	CertFile string // Server requires TLS client certificate authentication
	KeyFile  string // Server requires TLS client certificate authentication
	CAFile   string // Trusted root certificates for server

	// CertData holds PEM-encoded bytes (typically read from a client certificate file).
	// CertData takes precedence over CertFile
	CertData []byte
	// KeyData holds PEM-encoded bytes (typically read from a client certificate key file).
	// KeyData takes precedence over KeyFile
	KeyData []byte
	// CAData holds PEM-encoded bytes (typically read from a root certificates bundle).
	// CAData takes precedence over CAFile
	CAData []byte

	// NextProtos is a list of supported application level protocols, in order of preference.
	// Used to populate tls.Config.NextProtos.
	// To indicate to the server http/1.1 is preferred over http/2, set to ["http/1.1", "h2"] (though the server is free
	// to ignore that preference).
	// To use only http/1.1, set to ["http/1.1"].
	NextProtos []string
}

var (
	_ fmt.Stringer   = TLSClientConfig{}
	_ fmt.GoStringer = TLSClientConfig{}
)

type sanitizedTLSClientConfig TLSClientConfig

// GoString implements fmt.GoStringer and sanitizes sensitive fields of
// TLSClientConfig to prevent accidental leaking via logs.
func (c TLSClientConfig) GoString() string {
	return c.String()
}

// String implements fmt.Stringer and sanitizes sensitive fields of
// TLSClientConfig to prevent accidental leaking via logs.
func (c TLSClientConfig) String() string {
	// Copy the config and explicitly mark non-empty credential fields as redacted
	cc := sanitizedTLSClientConfig{
		Insecure:   c.Insecure,
		ServerName: c.ServerName,
		CertFile:   c.CertFile,
		KeyFile:    c.KeyFile,
		CAFile:     c.CAFile,
		CertData:   c.CertData,
		KeyData:    c.KeyData,
		CAData:     c.CAData,
		NextProtos: c.NextProtos,
	}

	if len(cc.CertData) != 0 {
		cc.CertData = []byte("--- TRUNCATED ---")
	}

	if len(cc.KeyData) != 0 {
		cc.KeyData = []byte("--- REDACTED ---")
	}

	return fmt.Sprintf("%#v", cc)
}

// HasCA returns whether the configuration has a certificate authority or not.
func (c TLSClientConfig) HasCA() bool {
	return len(c.CAData) > 0 || len(c.CAFile) > 0
}

// HasCertAuth returns whether the configuration has certificate authentication or not.
func (c TLSClientConfig) HasCertAuth() bool {
	return (len(c.CertData) != 0 || len(c.CertFile) != 0) && (len(c.KeyData) != 0 || len(c.KeyFile) != 0)
}
