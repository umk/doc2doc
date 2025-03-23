package main

import (
	"errors"
	"flag"
	"fmt"
	"os"
	"slices"
)

var config struct {
	inputs   []string
	output   string
	metadata string

	saveMetaOnly bool // Only save metadata file based on the input and output
	force        bool // Force generation
	autoconfirm  bool // Answer Yes to all confirmations

	service struct {
		baseURL string
		key     string
		model   string

		seed        int64
		temperature float64
		topP        float64
	}
}

func setConfig() error {
	readConfigFromEnv()

	if err := readConfigFromFlags(); err != nil {
		return err
	}

	if len(config.inputs) == 0 {
		return fmt.Errorf("input file path is required")
	}
	if config.output == "" {
		return fmt.Errorf("output file path is required")
	}
	if config.metadata == "" {
		config.metadata = config.output + ".d2d"
	}

	if err := checkRedirects(); err != nil {
		return err
	}

	return nil
}

func checkRedirects() error {
	hasRedirect := false

	for _, v := range config.inputs {
		if v == "-" {
			if hasRedirect {
				return errors.New("cannot use redirect to stdin for more than one input")
			}

			hasRedirect = true
		}
	}

	return nil
}

func readConfigFromEnv() {
	if v, ok := os.LookupEnv("D2D_BASE_URL"); ok {
		config.service.baseURL = v
	}
	if v, ok := os.LookupEnv("D2D_KEY"); ok {
		config.service.key = v

		os.Unsetenv("D2D_KEY")
	}
	if v, ok := os.LookupEnv("D2D_MODEL"); ok {
		config.service.model = v
	}
}

func readConfigFromFlags() error {
	// Define and parse command-line flags
	flag.Func("i", "input file path (required)", func(v string) error {
		config.inputs = append(config.inputs, v)
		return nil
	})
	flag.StringVar(&config.output, "o", "", "output file path (required)")
	flag.StringVar(&config.metadata, "d", "", "metadata file path")

	flag.BoolVar(&config.saveMetaOnly, "meta", false, "only save metadata given input and output")
	flag.BoolVar(&config.force, "force", false, "force generation")
	flag.BoolVar(&config.autoconfirm, "y", false, "confirm automatically")

	flag.StringVar(&config.service.baseURL, "svc.base", "", "service base URL")
	flag.StringVar(&config.service.key, "svc.key", "", "service key")
	flag.StringVar(&config.service.model, "svc.model", "", "service model name")

	flag.Int64Var(&config.service.seed, "gen.seed", 0, "generation seed")
	flag.Float64Var(&config.service.temperature, "gen.t", 0, "generation temperature")
	flag.Float64Var(&config.service.topP, "gen.p", 0, "generation top P")

	flag.Parse()

	slices.Sort(config.inputs)

	if flag.NArg() > 0 {
		return fmt.Errorf("unexpected number of arguments: %d", flag.NArg())
	}

	return nil
}
