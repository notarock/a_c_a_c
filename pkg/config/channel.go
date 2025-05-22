package config

import (
	"fmt"
	"os"
	"strconv"

	"gopkg.in/yaml.v3"
)

type ChannelConfig struct {
	Bots     []string  `yaml:"bots"`
	Channels []Channel `yaml:"channels"`
}

type Channel struct {
	Name      string   `yaml:"name"`
	Frequency int      `yaml:"frequency"`
	AllowBits bool     `yaml:"allow_bits"`
	ExtraBots []string `yaml:"extra_bots"`
}

func LoadChannelConfig(path string) (cfg ChannelConfig, err error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return cfg, fmt.Errorf("failed to read config file %s: %v", path, err)
	}

	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return cfg, fmt.Errorf("failed to parse config file %s: %v", path, err)
	}

	defaultFrequency, err := GetDefaultFrequency()
	if err != nil {
		return cfg, fmt.Errorf("failed to get default chat frequency: %v", err)
	}

	// Apply defaults
	for i := range cfg.Channels {
		ch := &cfg.Channels[i]

		if ch.Name == "" {
			return cfg, fmt.Errorf("channel name is required")
		}

		if ch.Frequency == 0 {
			ch.Frequency = defaultFrequency
		}

		if ch.ExtraBots == nil {
			ch.ExtraBots = []string{}
		}
	}

	fmt.Printf("Parsed config:\n%+v\n", cfg)

	return cfg, nil
}

func GetDefaultFrequency() (int, error) {
	COUNTDOWN := os.Getenv("COUNTDOWN")
	countdownInterval, err := strconv.Atoi(COUNTDOWN)
	return countdownInterval, err
}
