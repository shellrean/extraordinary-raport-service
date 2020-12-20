package config

import (
    "fmt"
    "os"
    "flag"

    "gopkg.in/yaml.v2"
)

func NewConfig(configPath string) (*Config, error) {
    // Create config structure
    config := &Config{}

    // Open config file
    file, err := os.Open(configPath)
    if err != nil {
        return nil, err
    }

    defer file.Close()

    // Init new YAML Decode
    d := yaml.NewDecoder(file)

    // Start YAML decoding from file
    if err := d.Decode(&config); err != nil {
        return nil, err
    }

    return config, nil
}

// ValdiateConfigPath just makes sure, that the path provided is a file
// that can be read
func ValdiateConfigPath(path string) error {
    s, err := os.Stat(path)
    if err != nil {
        return err
    }

    if s.IsDir() {
        return fmt.Errorf("'%s' is a directory, not a normal file", path)
    }

    return nil
}

// ParseFlags will create and parse the CLI flags
// and return the path to be used elsewhere
func ParseFlags() (string, error) {
    // String the containts the configurated configuration path
    var configPath string

    // Set up a CLI path called "-config" to allow users 
    // to supply the configuration file
    flag.StringVar(&configPath, "config", "./config.yml", "path to config file")

    // Actually parse the flags
    flag.Parse()

    // Validate the path first
    if err := ValdiateConfigPath(configPath); err != nil {
        return "", err
    }

    // Return the configuration path
    return configPath, nil
}