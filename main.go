package main

import (
	"bytes"
	"context"
	"crypto/sha256"
	"errors"
	"flag"
	"fmt"
	"os"

	"github.com/sergi/go-diff/diffmatchpatch"
)

func main() {
	ctx := context.Background()
	if err := run(ctx); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}

func run(ctx context.Context) error {
	if err := setConfig(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		flag.Usage()
		os.Exit(2)
	}

	if err := checkBackupsDontExist(config.output, config.metadata); err != nil {
		return err
	}

	outputExists, err := checkExists(config.output)
	if err != nil {
		return fmt.Errorf("error checking output file: %w", err)
	}

	metaExists, err := checkExists(config.metadata)
	if err != nil {
		return fmt.Errorf("error checking metadata file: %w", err)
	}

	var inputContent bytes.Buffer
	for i, in := range config.inputs {
		if i > 0 {
			inputContent.WriteString("\n\n")
		}

		if err := readInputOrStdin(&inputContent, in); err != nil {
			return err
		}
	}

	var outputContent []byte
	if outputExists {
		outputContent, err = os.ReadFile(config.output)
		if err != nil {
			return fmt.Errorf("error reading output file: %w", err)
		}
	}

	forceGenerate := config.force

	var previousMd *metadata
	if metaExists {
		md, err := metadataRead(config.metadata)
		if err != nil {
			return fmt.Errorf("error reading metadata file: %w", err)
		}
		previousMd = md

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

	var currentOut string

	if !config.saveMetaOnly {
	Retries:
		for {
			generated, err := generate(
				ctx,
				previousIn,
				previousOut,
				currentIn,
				config.output,
			)
			if err != nil {
				return fmt.Errorf("error generating output: %w", err)
			}

			currentOut = generated

			if previousOut != nil && !config.autoconfirm {
				d := diffmatchpatch.New()
				diffs := d.DiffMain(*previousOut, generated, false)
				fmt.Println(renderDiff(diffs))

				k, err := readKeyOrDefaultOf("Continue? (Y/n/r) ", 'y', 'n', 'r')
				if err != nil {
					return err
				}

				switch k {
				case 'n':
					return nil
				case 'r':
					continue Retries
				}
			}

			break Retries
		}
	} else {
		if outputExists {
			currentOut = *previousOut
		} else {
			return fmt.Errorf("previous output doesn't exist")
		}
	}

	currentOutBytes := []byte(currentOut)

	outputSha256 := sha256.Sum256(currentOutBytes)
	outputSha256Str := fmt.Sprintf("%x", outputSha256)

	newMetadata := &metadata{}
	newMetadata.Input.Content = currentIn
	newMetadata.Output.Sha256 = outputSha256Str

	bs := make(backupSet, 0, 2)

	var backupErr error

	if outputExists {
		if _, err := bs.createBackup(config.output); err != nil {
			backupErr = fmt.Errorf("error backing up output file: %w", err)
		}
	}

	if metaExists {
		if _, err := bs.createBackup(config.metadata); err != nil {
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

	if err := metadataWrite(config.metadata, newMetadata); err != nil {
		writeErr = fmt.Errorf("error writing metadata file: %w", err)
	} else if !config.saveMetaOnly {
		if err := atomicWrite(config.output, currentOutBytes); err != nil {
			writeErr = fmt.Errorf("error writing output file: %w", err)
		}
	}

	if writeErr != nil {
		if rerr := bs.restoreBackups(); rerr != nil {
			return fmt.Errorf("write error: %v; then failed to restore backups: %v", writeErr, rerr)
		}
		return fmt.Errorf("write error: %v; backups restored", writeErr)
	}

	bs.removeBackups()

	return nil
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
