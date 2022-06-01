package clientcmd

import (
	"io/ioutil"
	"path"

	yaml "gopkg.in/yaml.v3"

	"github.com/opsdata/common-base/pkg/util/homedir"
)

// Defines some useful for const variables.
const (
	RecommendedConfigPathFlag   = "elmtconfig"
	RecommendedConfigPathEnvVar = "ELMTCONFIG"
	RecommendedHomeDir          = ".elmt"
	RecommendedFileName         = "config"
	RecommendedSchemaName       = "scheme"
)

// Defines some useful variables.
var (
	RecommendedConfigDir  = path.Join(homedir.HomeDir(), RecommendedHomeDir)
	RecommendedHomeFile   = path.Join(RecommendedConfigDir, RecommendedFileName)
	RecommendedSchemaFile = path.Join(RecommendedConfigDir, RecommendedSchemaName)
)

// LoadFromFile load config from file.
func LoadFromFile(filename string) (*Config, error) {
	iamconfigBytes, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	config, err := Load(iamconfigBytes)
	if err != nil {
		return nil, err
	}

	// set LocationOfOrigin on every Cluster, User, and Context
	config.AuthInfo.LocationOfOrigin = filename
	config.Server.LocationOfOrigin = filename

	if config.AuthInfo == nil {
		config.AuthInfo = &AuthInfo{}
	}

	if config.Server == nil {
		config.Server = &Server{}
	}

	return config, nil
}

// Load takes a byte slice and deserializes the contents into Config object.
// Encapsulates deserialization without assuming the source is a file.
func Load(data []byte) (*Config, error) {
	config := NewConfig()
	// if there's no data in a file, return the default object instead of failing (DecodeInto reject empty input)
	if len(data) == 0 {
		return config, nil
	}

	if err := yaml.Unmarshal(data, config); err != nil {
		return nil, err
	}

	return config, nil
}
