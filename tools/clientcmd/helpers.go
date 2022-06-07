package clientcmd

import (
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"strconv"
	"time"

	yaml "gopkg.in/yaml.v3"
)

// ParseTimeout
// - return a parsed duration from a string
// - a duration string value must be a positive integer, optionally followed by a corresponding time unit (s|m|h).
func ParseTimeout(duration string) (time.Duration, error) {
	if i, err := strconv.ParseInt(duration, 10, 64); err == nil && i >= 0 {
		return (time.Duration(i) * time.Second), nil
	}

	if requestTimeout, err := time.ParseDuration(duration); err == nil {
		return requestTimeout, nil
	}

	return 0, fmt.Errorf(
		"invalid timeout value. Timeout must be a single integer in seconds, or an integer followed by a corresponding time unit (e.g. 1s | 2m | 3h)",
	)
}

const (
	RecommendedConfigPathFlag   = "elmtconfig"
	RecommendedConfigPathEnvVar = "ELMTCONFIG"
	RecommendedHomeDir          = ".elmt"
	RecommendedFileName         = "config"
	RecommendedSchemaName       = "scheme"
)

var (
	RecommendedConfigDir  = path.Join(os.Getenv("HOME"), RecommendedHomeDir)
	RecommendedHomeFile   = path.Join(RecommendedConfigDir, RecommendedFileName)
	RecommendedSchemaFile = path.Join(RecommendedConfigDir, RecommendedSchemaName)
)

// Load
// - take a byte slice and deserialize the contents into Config object
// - encapsulate deserialization without assuming the source is a file
func Load(data []byte) (*Config, error) {
	config := NewConfig()

	// If there is no data in a file, return the default object instead of failing
	if len(data) == 0 {
		return config, nil
	}

	if err := yaml.Unmarshal(data, config); err != nil {
		return nil, err
	}

	return config, nil
}

// LoadFromFile
// - load config from file
func LoadFromFile(filename string) (*Config, error) {
	elmtconfigBytes, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	config, err := Load(elmtconfigBytes)
	if err != nil {
		return nil, err
	}

	// Set the LocationOfOrigin
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
