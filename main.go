package main

import (
	"bytes"
	"context"
	"crypto/sha256"
	"errors"
	"fmt"
	"io"
	"os"
)

func main() {
	ctx := context.Background()
	if err := run(ctx); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}

func run(ctx context.Context) error {
	config, err := getConfig()
	if err != nil {
		return err
	}

	if err := checkBackupsDontExist(config.OutputPath, config.MetaPath); err != nil {
		return err
	}

	outputExists, err := checkExists(config.OutputPath)
	if err != nil {
		return fmt.Errorf("error checking output file: %w", err)
	}

	metaExists, err := checkExists(config.MetaPath)
	if err != nil {
		return fmt.Errorf("error checking metadata file: %w", err)
	}

	var inputContent bytes.Buffer
	for i, in := range config.InputPath {
		if i > 0 {
			inputContent.WriteString("\n\n")
		}

		if err := readInputOrStdin(&inputContent, in); err != nil {
			return err
		}
	}

	var outputContent []byte
	if outputExists {
		outputContent, err = os.ReadFile(config.OutputPath)
		if err != nil {
			return fmt.Errorf("error reading output file: %w", err)
		}
	}

	prompt, err := readPromptOrStdin(config.Prompt)
	if err != nil {
		return err
	}

	forceGenerate := config.Force

	var previousMd *metadata
	if metaExists {
		md, err := metadataRead(config.MetaPath)
		if err != nil {
			return fmt.Errorf("error reading metadata file: %w", err)
		}
		previousMd = md

		metaPrompt := md.Input.Prompt

		if prompt == "" {
			prompt = metaPrompt
		} else if prompt != metaPrompt {
			fmt.Fprintf(os.Stderr, "Warning: provided prompt differs from metadata prompt.\n")
			forceGenerate = true
			prompt = metaPrompt
		}

		if outputExists {
			sum := sha256.Sum256(outputContent)
			outputHash := fmt.Sprintf("%x", sum)
			if outputHash != md.Output.Sha256 {
				fmt.Fprintf(os.Stderr, "Warning: output file's SHA-256 differs from metadata.\n")
			}
		}
	}

	var previousIn, previousOut *string

	if metaExists {
		previousIn = &previousMd.Input.Content
	}
	if outputExists {
		previousOut = stringPtr(string(outputContent))
	}

	currentIn := inputContent.String()

	if (previousIn != nil) && *previousIn == currentIn && !forceGenerate {
		fmt.Println("Previous and current inputs are same. Generation aborted.")
		return nil
	}

	generated, err := generate(
		ctx,
		&config.Service,
		prompt,
		previousIn,
		previousOut,
		currentIn,
		config.OutputPath,
	)
	if err != nil {
		return fmt.Errorf("error generating output: %w", err)
	}

	outputSha256 := sha256.Sum256([]byte(generated))
	outputSha256Str := fmt.Sprintf("%x", outputSha256)

	newMetadata := &metadata{}
	newMetadata.Input.Content = currentIn
	newMetadata.Input.Prompt = prompt
	newMetadata.Output.Sha256 = outputSha256Str

	bs := make(backupSet, 0, 2)

	var backupErr error

	if outputExists {
		if _, err := bs.createBackup(config.OutputPath); err != nil {
			backupErr = fmt.Errorf("error backing up output file: %w", err)
		}
	}

	if metaExists {
		if _, err := bs.createBackup(config.MetaPath); err != nil {
			backupErr = fmt.Errorf("error backing up metadata file: %w", err)
		}
	}

	if backupErr != nil {
		if rerr := bs.restoreBackups(); rerr != nil {
			return fmt.Errorf("backup error: %v; then failed to restore backups: %v", backupErr, rerr)
		}
		return fmt.Errorf("backup error: %v; backups restored", backupErr)
	}

	var writeErr error

	if err := metadataWrite(config.MetaPath, newMetadata); err != nil {
		writeErr = fmt.Errorf("error writing metadata file: %w", err)
	} else if err := atomicWrite(config.OutputPath, []byte(generated)); err != nil {
		writeErr = fmt.Errorf("error writing output file: %w", err)
	}

	if writeErr != nil {
		if rerr := bs.restoreBackups(); rerr != nil {
			return fmt.Errorf("write error: %v; then failed to restore backups: %v", backupErr, rerr)
		}
		return fmt.Errorf("write error: %v; backups restored", backupErr)
	}

	bs.removeBackups()

	return nil
}

func readPromptOrStdin(prompt string) (string, error) {
	if prompt == "-" {
		b, err := io.ReadAll(os.Stdin)
		if err != nil {
			return "", fmt.Errorf("error reading prompt from stdin: %w", err)
		}

		return string(b), nil
	}

	return prompt, nil
}

func readInputOrStdin(dst *bytes.Buffer, src string) error {
	if src == "-" {
		if _, err := dst.ReadFrom(os.Stdin); err != nil {
			return fmt.Errorf("error reading input from stdin: %w", err)
		}

		return nil
	}

	f, err := os.Open(src)
	if err != nil {
		return fmt.Errorf("error reading input file: %w", err)
	}
	defer f.Close()

	if _, err := dst.ReadFrom(f); err != nil {
		return fmt.Errorf("error reading input file: %w", err)
	}

	return nil
}

func checkBackupsDontExist(outputPath, metaPath string) error {
	if backupExists, err := checkBackupExists(outputPath); err != nil {
		return fmt.Errorf("error checking backup for output: %w", err)
	} else if backupExists {
		return errors.New("backup file for output exists; aborting")
	}

	if backupExists, err := checkBackupExists(metaPath); err != nil {
		return fmt.Errorf("error checking backup for metadata: %w", err)
	} else if backupExists {
		return errors.New("backup file for metadata exists; aborting")
	}

	return nil
}
