package config

import (
	"fmt"
	"os"

	"github.com/hashhavoc/teller/pkg/api/alex"
	"github.com/hashhavoc/teller/pkg/api/gobob"
	"github.com/hashhavoc/teller/pkg/api/hiro"
	"github.com/hashhavoc/teller/pkg/api/stxtools"

	"gopkg.in/yaml.v2"
)

type Config struct {
	Path      string          `yaml:"-"`
	Endpoints ConfigEndpoints `yaml:"endpoints"`
	Wallets   []string        `yaml:"wallets"`
	OpenAI    OpenAIConfig    `yaml:"openai"`
}

type ConfigEndpoints struct {
	Hiro     string `yaml:"hiro"`
	Ord      string `yaml:"ord"`
	Alex     string `yaml:"alex"`
	StxTools string `yaml:"stxtools"`
	Bob      string `yaml:"bob"`
}

type OpenAIConfig struct {
	APIKey  string `yaml:"api_key"`
	Model   string `yaml:"model"`
	BaseURL string `yaml:"base_url"`
}

func NewConfig(path string) *Config {
	config := &Config{
		Path: path,
		Endpoints: ConfigEndpoints{
			Hiro:     hiro.DefaultApiBase,
			Alex:     alex.DefaultApiBase,
			StxTools: stxtools.DefaultApiBase,
			Bob:      gobob.DefaultApiBase,
		},
		OpenAI: OpenAIConfig{
			APIKey:  "",                          // Will be set from environment or config
			Model:   "gpt-4",                     // Default model
			BaseURL: "https://api.openai.com/v1", // Default OpenAI base URL
		},
	}
	return config
}

func (c *Config) ReadConfig() error {
	bytes, err := os.ReadFile(c.Path)
	if err != nil {
		return err
	}
	if err := yaml.Unmarshal(bytes, c); err != nil {
		return err
	}

	// Ensure wallets are unique
	walletMap := make(map[string]bool)
	for _, wallet := range c.Wallets {
		if _, exists := walletMap[wallet]; exists {
			return fmt.Errorf("duplicate wallet found: %s", wallet)
		}
		walletMap[wallet] = true
	}

	return nil
}

func (c *Config) AddWallet(wallet string) error {
	// Check for duplicate before adding
	for _, w := range c.Wallets {
		if w == wallet {
			return fmt.Errorf("duplicate wallet: %s", wallet)
		}
	}
	// If not found, append the wallet
	c.Wallets = append(c.Wallets, wallet)
	return nil
}

func (c *Config) RemoveWallet(wallet string) {
	for i, w := range c.Wallets {
		if w == wallet {
			c.Wallets = append(c.Wallets[:i], c.Wallets[i+1:]...)
			break
		}
	}
}

func (c *Config) WriteConfig() error {
	bytes, err := yaml.Marshal(c)
	if err != nil {
		return err
	}
	return os.WriteFile(c.Path, bytes, 0644)
}
