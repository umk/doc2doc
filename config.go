package main

import (
	"errors"
	"flag"
	"fmt"
	"os"
	"strings"
)

type config struct {
	InputPath  []string
	OutputPath string
	MetaPath   string

	Prompt string

	Force bool // Force generation

	Service configService
}

type configService struct {
	BaseURL *string
	Key     *string
	Model   *string

	Seed        *int64
	Temperature *float64
	TopP        *float64
}

func getConfig() (*config, error) {
	config := &config{}

	readConfigFromEnv(config)

	if err := readConfigFromFlags(config); err != nil {
		return nil, err
	}

	if len(config.InputPath) == 0 {
		return nil, fmt.Errorf("input file path is required")
	}
	if config.OutputPath == "" {
		return nil, fmt.Errorf("output file path is required")
	}
	if config.MetaPath == "" {
		config.MetaPath = config.OutputPath + ".d2d"
	}

	if err := checkRedirects(config); err != nil {
		return nil, err
	}

	return config, nil
}

func checkRedirects(c *config) error {
	hasRedirect := false

	values := append([]string{c.Prompt}, c.InputPath...)
	for _, v := range values {
		if v == "-" {
			if hasRedirect {
				return errors.New("cannot use redirect to stdin for more than one input")
			}

			hasRedirect = true
		}
	}

	return nil
}

func readConfigFromEnv(c *config) {
	if v, ok := os.LookupEnv("D2D_BASE_URL"); ok {
		c.Service.BaseURL = &v
	}
	if v, ok := os.LookupEnv("D2D_KEY"); ok {
		c.Service.Key = &v

		os.Unsetenv("D2D_KEY")
	}
	if v, ok := os.LookupEnv("D2D_MODEL"); ok {
		c.Service.Model = &v
	}
}

func readConfigFromFlags(c *config) error {
	// Define and parse command-line flags
	flag.Func("i", "input file path (required)", func(v string) error {
		c.InputPath = append(c.InputPath, v)
		return nil
	})
	outputPath := flag.String("o", "", "output file path (required)")
	metaPath := flag.String("d", "", "metadata file path")

	force := flag.Bool("force", false, "force generation")

	baseURL := flag.String("svc.base", "", "service base URL")
	key := flag.String("svc.key", "", "service key")
	model := flag.String("svc.model", "", "service model name")

	seed := flag.Int64("gen.seed", 0, "generation seed")
	temperature := flag.Float64("gen.t", 0, "generation temperature")
	topP := flag.Float64("gen.p", 0, "generation top P")

	flag.Parse()

	// Handle positional arguments and assign values to Config
	var prompt string
	switch flag.NArg() {
	case 0:
		// No additional arguments
	case 1:
		prompt = strings.TrimSpace(flag.Arg(0))
	default:
		return fmt.Errorf("unexpected number of arguments: %d", flag.NArg())
	}

	c.OutputPath = *outputPath
	c.MetaPath = *metaPath

	c.Prompt = prompt

	c.Force = *force

	// Assign service-specific parameters
	if *baseURL != "" {
		c.Service.BaseURL = baseURL
	}
	if *key != "" {
		c.Service.Key = key
	}
	if *model != "" {
		c.Service.Model = model
	}

	// Assign generation parameters
	if *seed != 0 {
		c.Service.Seed = seed
	}
	if *temperature != 0 {
		c.Service.Temperature = temperature
	}
	if *topP != 0 {
		c.Service.TopP = topP
	}

	return nil
}
