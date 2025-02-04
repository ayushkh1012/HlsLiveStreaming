package utils

import (
    "fmt"
    "gopkg.in/yaml.v2"
    "io/ioutil"
)

type Config struct {
    Server struct {
        Host       string `yaml:"host"`
        Port       int    `yaml:"port"`
        WindowSize int    `yaml:"window_size"`
    } `yaml:"server"`
    Paths struct {
        Media     string `yaml:"media"`
        Ads       string `yaml:"ads"`
        Manifests string `yaml:"manifests"`
    } `yaml:"paths"`
    Ads []struct {
        ID       string `yaml:"id"`
        Duration int    `yaml:"duration"`
    } `yaml:"ads"`
}

func LoadConfig(configPath string) (*Config, error) {
    data, err := ioutil.ReadFile(configPath)
    if err != nil {
        return nil, fmt.Errorf("error reading config file: %v", err)
    }

    var config Config
    if err := yaml.Unmarshal(data, &config); err != nil {
        return nil, fmt.Errorf("error parsing config file: %v", err)
    }

    // Validate the config
    if err := config.Validate(); err != nil {
        return nil, fmt.Errorf("invalid configuration: %v", err)
    }

    return &config, nil
}

func (c *Config) Validate() error {
    if c.Server.Port <= 0 || c.Server.Port > 65535 {
        return fmt.Errorf("invalid port number: %d", c.Server.Port)
    }
    if c.Server.WindowSize <= 0 {
        return fmt.Errorf("invalid window size: %d", c.Server.WindowSize)
    }
    return nil
}
