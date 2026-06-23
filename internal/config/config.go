package config

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
)

type Config struct {
	Commands []Command `json:"commands"`
}

type Command struct {
	Name string   `json:"name"`
	Args []string `json:"args"`
}

const Starter = `{
  "commands": []
}
`

func DefaultPath() (string, error) {
	dir, err := os.UserConfigDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(dir, "lazycont", "config.json"), nil
}

func LoadDefault() (Config, string, error) {
	path, err := DefaultPath()
	if err != nil {
		return Config{}, "", err
	}
	cfg, err := Load(path)
	return cfg, path, err
}

func Ensure(path string) error {
	if strings.TrimSpace(path) == "" {
		return errors.New("config path is required")
	}
	if _, err := os.Stat(path); err == nil {
		return nil
	} else if !errors.Is(err, os.ErrNotExist) {
		return err
	}
	if err := os.MkdirAll(filepath.Dir(path), 0o700); err != nil {
		return err
	}
	return os.WriteFile(path, []byte(Starter), 0o600)
}

func Load(path string) (Config, error) {
	file, err := os.Open(path)
	if errors.Is(err, os.ErrNotExist) {
		return Config{}, nil
	}
	if err != nil {
		return Config{}, err
	}
	defer file.Close()

	var cfg Config
	decoder := json.NewDecoder(file)
	decoder.DisallowUnknownFields()
	if err := decoder.Decode(&cfg); err != nil {
		return Config{}, err
	}
	if err := decoder.Decode(&struct{}{}); !errors.Is(err, io.EOF) {
		return Config{}, errors.New("config contains trailing JSON data")
	}
	if err := cfg.normalize(); err != nil {
		return Config{}, err
	}
	return cfg, nil
}

func (c *Config) normalize() error {
	for idx := range c.Commands {
		command := &c.Commands[idx]
		command.Name = strings.TrimSpace(command.Name)
		if command.Name == "" {
			return fmt.Errorf("commands[%d].name is required", idx)
		}
		if len(command.Args) == 0 || strings.TrimSpace(command.Args[0]) == "" {
			return fmt.Errorf("commands[%d].args must start with a container subcommand", idx)
		}
		for argIndex := range command.Args {
			command.Args[argIndex] = strings.TrimSpace(command.Args[argIndex])
		}
	}
	return nil
}
