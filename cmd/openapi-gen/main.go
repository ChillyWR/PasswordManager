package main

import (
	"encoding/json"
	"flag"
	"os"
	"path"

	"github.com/okutsen/PasswordManager/config"
	"github.com/okutsen/PasswordManager/internal/api"
	"github.com/okutsen/PasswordManager/internal/log"

	"github.com/getkin/kin-openapi/openapi3"
	"github.com/invopop/yaml"
)

const (
	defaultOutputPath     = "."
	defaultOutputFileType = "yaml"
)

func main() {
	var logger log.Logger = log.NewLogrusLogger()

	var outputPath string
	var outputFileType string

	flag.StringVar(&outputPath, "path", defaultOutputPath, "Path to use for generating OpenAPIv3 spec")
	flag.StringVar(&outputFileType, "file-type", defaultOutputFileType, "File type of an output: \"json\", \"yaml\" or \"both\"")
	flag.Parse()

	if outputFileType != "json" && outputFileType != "yaml" && outputFileType != "both" {
		logger.Fatalf("Invalid file_type flag: %s", outputFileType)
	}

	cfg, err := config.NewConfig()
	if err != nil {
		logger.Fatalf("Failed to initialize config: %v", err)
	}

	spec := api.NewOpenAPIv3(&api.Config{Port: cfg.Port})

	switch outputFileType {
	case "json":
		genJSON(spec, outputPath, logger)
	case "yaml":
		genYAML(spec, outputPath, logger)
	case "both":
		genJSON(spec, outputPath, logger)
		genYAML(spec, outputPath, logger)
	}
	logger.Infof("Successfully generated. Path: %s", outputPath)
}

func genJSON(spec *openapi3.T, targetPath string, logger log.Logger) {
	data, err := json.Marshal(&spec)
	if err != nil {
		logger.Fatalf("Couldn't marshal json: %s", err)
	}

	if err := os.WriteFile(path.Join(targetPath, "openapi3.json"), data, 0666); err != nil {
		logger.Fatalf("Couldn't write json: %s", err)
	}
}

func genYAML(spec *openapi3.T, targetPath string, logger log.Logger) {
	data, err := yaml.Marshal(&spec)
	if err != nil {
		logger.Fatalf("Couldn't marshal yaml: %s", err)
	}

	if err := os.WriteFile(path.Join(targetPath, "openapi3.yaml"), data, 0666); err != nil {
		logger.Fatalf("Couldn't write yaml: %s", err)
	}

}
