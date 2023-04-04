package errgen

import (
	"bytes"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"

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

type Cell struct {
	Value string
	Width int
}

func LoadConfig(configFile string, packageName string) (*Config, error) {
	cfg := Config{
		PackageName: packageName,
	}
	data, err := os.ReadFile(configFile)
	if err != nil {
		return nil, err
	}
	if err := yaml.Unmarshal(data, &cfg.Errors); err != nil {
		return nil, fmt.Errorf("unmarshal YAML data: %v", err)
	}

	return &cfg, nil
}

func Generate(cfg *Config, outputFile string, fileType string) error {
	if err := os.MkdirAll(filepath.Dir(outputFile), 0755); err != nil {
		return fmt.Errorf("create output directory: %v", err)
	}

	file, err := os.Create(outputFile)
	if err != nil {
		return fmt.Errorf("create errors file %s: %v", outputFile, err)
	}
	defer file.Close()

	switch fileType {
	case "md":
		w := GenerateMarkdown(cfg)
		_, err := fmt.Fprint(file, w)
		if err != nil {
			return err
		}
	case "go":
		if err := errorsTemplate.Execute(file, cfg); err != nil {
			return fmt.Errorf("execute errorsTemplate: %v", err)
		}
	default:
		return errors.New(fmt.Sprintf("not support file type: %s", fileType))
	}

	return nil
}

func GenerateMarkdown(cfg *Config) *bytes.Buffer {
	lengths := maxLength(cfg)

	b := &bytes.Buffer{}
	writeRow(b, []Cell{
		{
			Value: "Status",
			Width: lengths["Status"],
		},
		{
			Value: "Code",
			Width: lengths["Code"],
		},
		{
			Value: "Message",
			Width: lengths["Message"],
		},
	})
	b.WriteString(fmt.Sprintf("|%s|%s|%s|", strings.Repeat("-", lengths["Status"]), strings.Repeat("-", lengths["Code"]), strings.Repeat("-", lengths["Message"])))

	b.WriteString("\n")
	for _, v := range cfg.Errors {
		writeRow(b, []Cell{
			{
				Value: strconv.Itoa(v.Status),
				Width: lengths["Status"],
			},
			{
				Value: strconv.Itoa(v.Code),
				Width: lengths["Code"],
			},
			{
				Value: v.Message,
				Width: lengths["Message"],
			},
		})
	}

	return b

}

func writeRow(b *bytes.Buffer, row []Cell) {
	for i, s := range row {
		cell := alignCenter(s.Value, s.Width)
		if i == 0 || i == len(row) {
			b.WriteString("|")
		}
		b.WriteString(cell)
		b.WriteString("|")
	}
	b.WriteString("\n")

}

func alignCenter(s string, totalLength int) string {
	padding := totalLength - len(s)
	if padding%2 == 1 {
		s = " " + s
		padding--
	}

	leftPadding := padding / 2
	rightPadding := padding - leftPadding
	return fmt.Sprintf("%s%*s%s", strings.Repeat(" ", leftPadding), len(s), s, strings.Repeat(" ", rightPadding))
}

// maxLength gets the maximum length of each column
func maxLength(cfg *Config) map[string]int {
	l := map[string]int{
		"Status":  len("Status"),
		"Code":    len("Code"),
		"Message": len("Message"),
	}
	for _, value := range cfg.Errors {
		if len(strconv.Itoa(value.Code)) > l["Status"] {
			l["Status"] = len(strconv.Itoa(value.Status))
		}

		if len(strconv.Itoa(value.Status)) > l["Code"] {
			l["Code"] = len(strconv.Itoa(value.Code))
		}

		if len(value.Message) > l["Message"] {
			l["Message"] = len(value.Message)
		}
	}
	l["Status"] += 2
	l["Code"] += 2
	l["Message"] += 2

	return l
}
