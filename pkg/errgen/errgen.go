package errgen

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"text/template"

	"gopkg.in/yaml.v3"
)

type Error struct {
	Code    int    `yaml:"code"`
	Status  int    `yaml:"status"`
	Message string `yaml:"message"`
}

type Config struct {
	PackageName string
	Errors      map[string]Error
}

const templateFileName = "errors.tmpl"

func LoadConfig(configFile string, packageName string) (*Config, error) {
	cfg := Config{
		PackageName: packageName,
	}

	data, err := ioutil.ReadFile(configFile)
	if err != nil {
		return nil, err
	}
	if err := yaml.Unmarshal(data, &cfg.Errors); err != nil {
		return nil, fmt.Errorf("unmarshal YAML data: %v", err)
	}

	return &cfg, nil
}

func Generate(cfg *Config, outputFile string) error {
	errorsTemplate := template.Must(template.ParseFiles(templateFileName))
	dir := filepath.Dir(outputFile)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("failed to create output directory: %v", err)
	}

	file, err := os.Create(outputFile)
	if err != nil {
		return fmt.Errorf("failed to create errors file %s: %v", outputFile, err)
	}
	defer file.Close()

	if err := errorsTemplate.Execute(file, cfg); err != nil {
		return fmt.Errorf("failed to execute errorsTemplate: %v", err)
	}

	return nil
}
