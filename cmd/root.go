package cmd

import (
	"log"
	"os"

	"github.com/spf13/cobra"

	"github.com/dinhquockhanh/eg/pkg/errgen"
)

var (
	version        = "1.2.0"
	configFile     string
	outputFile     string
	outputFileType string
	packageName    string
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:     "eg",
	Short:   "Here is an application that generates error variables from a YAML configuration file.",
	Long:    `Here is an application that generates error variables from a YAML configuration file.`,
	Version: version,
	Run: func(cmd *cobra.Command, args []string) {
		if configFile == "" {
			log.Println("please set config file path")
			return
		}
		config, err := errgen.LoadConfig(configFile, packageName)
		if err != nil {
			log.Fatal(err)
		}

		if err := errgen.Generate(config, outputFile, outputFileType); err != nil {
			log.Fatal(err)
		}
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func RootCmd() *cobra.Command {
	return rootCmd
}
func init() {
	rootCmd.PersistentFlags().StringVarP(&configFile, "file", "f", "", "a YAML configuration file")
	rootCmd.PersistentFlags().StringVarP(&outputFile, "output", "o", "./errors/errors.go", "Output file")
	rootCmd.PersistentFlags().StringVarP(&outputFileType, "type", "t", "go", "Output file")
	rootCmd.PersistentFlags().StringVarP(&packageName, "package", "p", "errors", "the Go package name")
}
