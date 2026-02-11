package lib

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

// Config represents our configuration structure
type Config struct {
	Project struct {
		Name    string `yaml:"name"`
		Version string `yaml:"version"`
	} `yaml:"project"`

	Git struct {
		RepositoryPath string `yaml:"repository_path"`
		DefaultBranch  string `yaml:"default_branch"`
	} `yaml:"git"`

	Output struct {
		Format   string `yaml:"format"`
		Filename string `yaml:"filename"`
	} `yaml:"output"`

	AI struct {
		Enabled  bool   `yaml:"enabled"`
		Provider string `yaml:"provider"`
		Model    string `yaml:"model"`
	} `yaml:"ai"`

	Categories []string `yaml:"categories"`
}

func LoadConfig(filename string) (*Config, error) {
	//Read the file
	data, err := os.ReadFile(filename)
	if err != nil {
		return nil, fmt.Errorf("failed to read config file: %w", err)
	}

	//Create a Config struct instance to hold the data
	var config Config

	//Unmarshal YAML data into the Config struct
	err = yaml.Unmarshal(data, &config)
	if err != nil {
		return nil, fmt.Errorf("failed to parse config: %w", err)
	}

	return &config, nil
}

func PrintConfig(config *Config) {
	fmt.Println("Current Configuration:")
	fmt.Println()
	fmt.Printf("  Project: %s (v%s)\n", config.Project.Name, config.Project.Version)
	fmt.Printf("  Repository: %s\n", config.Git.RepositoryPath)
	fmt.Printf("  Output: %s (%s)\n", config.Output.Filename, config.Output.Format)
	fmt.Printf("  AI: %v (%s)\n", config.AI.Enabled, config.AI.Provider)
	fmt.Printf("  Categories: %v\n", config.Categories)
}
