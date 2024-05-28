package app

import (
	"encoding/base64"
	"fmt"
	"os"
	"path"

	"gopkg.in/yaml.v3"
)

const (
	configFileName = "config.yaml"
)

type config struct {
	GhAccessToken string `yaml:"github-access-token,omitempty"`
	BearerToken   string `yaml:"bearer-token,omitempty"`

	decoded bool
}

// Encode encodes appropriate fields in the config to base64. If the config is already encoded, this is a no-op.
func (c *config) Encode(encoding *base64.Encoding) {
	// todo: can result in partially encoded config

	defer func() {
		c.decoded = false
	}()

	if !c.decoded {
		return
	}

	c.GhAccessToken = encoding.EncodeToString([]byte(c.GhAccessToken))
	c.BearerToken = encoding.EncodeToString([]byte(c.BearerToken))
}

func (c *config) Decode(encoding *base64.Encoding) error {
	// todo: can result in partially decoded config

	defer func() {
		c.decoded = true
	}()

	if c.decoded {
		return nil
	}

	if decoded, err := encoding.DecodeString(c.GhAccessToken); err != nil {
		return fmt.Errorf("failed to decode github access token: %w", err)
	} else {
		c.GhAccessToken = string(decoded)
	}

	if decoded, err := encoding.DecodeString(c.BearerToken); err != nil {
		return fmt.Errorf("failed to decode github access token: %w", err)
	} else {
		c.BearerToken = string(decoded)
	}

	return nil
}

type configManager struct {
	Root     string
	Encoding *base64.Encoding
	Config   *config
}

func newConfigManager(configRootDir string) (*configManager, error) {
	if configRootDir == "" {
		userConfigDir, err := os.UserConfigDir()
		if err != nil {
			return nil, fmt.Errorf("could not determine default config dir: %w", err)
		}

		configRootDir = path.Join(userConfigDir, "marina")
	}

	cm := &configManager{
		Root:     configRootDir,
		Encoding: base64.StdEncoding,
		Config:   &config{},
	}

	configPath := path.Join(configRootDir, configFileName)
	data, err := os.ReadFile(configPath)
	if err != nil && !os.IsNotExist(err) {
		return nil, fmt.Errorf("failed to open file: %w", err)
	} else if err == nil {
		if err := yaml.Unmarshal(data, cm.Config); err != nil {
			return nil, fmt.Errorf("failed to unmarshal config: %w", err)
		}

		if err := cm.Config.Decode(cm.Encoding); err != nil {
			return nil, fmt.Errorf("failed to decode config: %w", err)
		}
	}

	err = cm.Config.Decode(cm.Encoding)
	if err != nil {
		return nil, fmt.Errorf("failed to decode config: %w", err)
	}

	return cm, nil
}

func (cm *configManager) Close() error {
	cm.Config.Encode(cm.Encoding)

	if err := os.MkdirAll(cm.Root, 0777); err != nil {
		return fmt.Errorf("could not create config dir: %w", err)
	}

	configPath := path.Join(cm.Root, configFileName)

	configFile, err := os.Create(configPath)
	if err != nil {
		return fmt.Errorf("failed to open file: %w", err)
	}
	defer configFile.Close()

	data, err := yaml.Marshal(cm.Config)
	if err != nil {
		return fmt.Errorf("failed to marshal config: %w", err)
	}

	if _, err := configFile.Write(data); err != nil {
		return fmt.Errorf("failed to write config: %w", err)
	}

	return nil
}
