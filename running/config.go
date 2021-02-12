package running

import (
	"fmt"
	"github.com/blang/vfs"
	"gopkg.in/yaml.v2"
	"strings"
)

const ConfigFileName = "runnr.yml"

type Config struct {
	Outpath    string `yml:"outpath"`
	Executable string `yml:"executable"`
}

func ReadConfig(fs vfs.Filesystem) (*Config, error) {
	configBytes, err := vfs.ReadFile(fs, ConfigFileName)
	if err != nil {
		return nil, err
	}
	config := &Config{}
	err = yaml.Unmarshal(configBytes, config)
	if err != nil {
		return nil, err
	}
	return config, validateConfig(config)
}

func validateConfig(config *Config) error {
	if config.Executable == "" {
		return fmt.Errorf("must specify 'execute' variable in runnr config")
	}
	if config.Outpath == "" {
		return fmt.Errorf("must specify 'output' variable in runnr config")
	}
	if strings.HasSuffix(config.Outpath, ".go") {
		return fmt.Errorf("outpath variable should not end with .go, you might have mixed up execute and outpath variables")
	}
	return nil
}
