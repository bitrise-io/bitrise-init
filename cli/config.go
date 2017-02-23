package cli

import (
	"errors"
	"fmt"
	"os"
	"path"
	"path/filepath"

	"github.com/bitrise-core/bitrise-init/models"
	"github.com/bitrise-core/bitrise-init/output"
	"github.com/bitrise-core/bitrise-init/scanner"
	"github.com/bitrise-io/go-utils/colorstring"
	"github.com/bitrise-io/go-utils/command"
	"github.com/bitrise-io/go-utils/log"
	"github.com/bitrise-io/go-utils/pathutil"
	"github.com/urfave/cli"
)

const (
	defaultScanResultDir = "_scan_result"
)

var configCommand = cli.Command{
	Name:  "config",
	Usage: "Generates a bitrise config files based on your project.",
	Action: func(c *cli.Context) error {
		if err := initConfig(c); err != nil {
			log.Errorft(err.Error())
			os.Exit(1)
		}
		return nil
	},
	Flags: []cli.Flag{
		cli.StringFlag{
			Name:  "dir",
			Usage: "Directory to scan.",
			Value: "./",
		},
		cli.StringFlag{
			Name:  "output-dir",
			Usage: "Directory to save scan results.",
			Value: "./_scan_result",
		},
		cli.StringFlag{
			Name:  "format",
			Usage: "Output format, options [json, yaml].",
			Value: "yaml",
		},
	},
}

func writeOutput(outputDir string, scanResult models.ScanResultModel, format output.Format) (string, error) {
	if exist, err := pathutil.IsDirExists(outputDir); err != nil {
		return "", fmt.Errorf("Failed to check if path (%s) exists, error: %s", outputDir, err)
	} else if !exist {
		if err := os.MkdirAll(outputDir, 0700); err != nil {
			return "", fmt.Errorf("Failed to create (%s), error: %s", outputDir, err)
		}
	}

	pth := path.Join(outputDir, "result")
	outputPth, err := output.WriteToFile(scanResult, format, pth)
	if err != nil {
		return "", fmt.Errorf("Failed to write result, error: %s", err)
	}

	return outputPth, nil
}

func initConfig(c *cli.Context) error {
	// Config
	isCI := c.GlobalBool("ci")
	searchDir := c.String("dir")
	outputDir := c.String("output-dir")
	formatStr := c.String("format")

	if isCI {
		log.Infoft(colorstring.Yellow("CI mode"))
	}
	log.Infoft(colorstring.Yellowf("scan dir: %s", searchDir))
	log.Infoft(colorstring.Yellowf("output dir: %s", outputDir))
	log.Infoft(colorstring.Yellowf("output format: %s", formatStr))
	fmt.Println()

	currentDir, err := pathutil.AbsPath("./")
	if err != nil {
		return fmt.Errorf("Failed to expand path (%s), error: %s", outputDir, err)
	}

	if searchDir == "" {
		searchDir = currentDir
	}
	searchDir, err = pathutil.AbsPath(searchDir)
	if err != nil {
		return fmt.Errorf("Failed to expand path (%s), error: %s", outputDir, err)
	}

	if outputDir == "" {
		outputDir = filepath.Join(currentDir, defaultScanResultDir)
	}
	outputDir, err = pathutil.AbsPath(outputDir)
	if err != nil {
		return fmt.Errorf("Failed to expand path (%s), error: %s", outputDir, err)
	}
	if exist, err := pathutil.IsDirExists(outputDir); err != nil {
		return err
	} else if !exist {
		if err := os.MkdirAll(outputDir, 0700); err != nil {
			return fmt.Errorf("Failed to create (%s), error: %s", outputDir, err)
		}
	}

	if formatStr == "" {
		formatStr = output.YAMLFormat.String()
	}
	format, err := output.ParseFormat(formatStr)
	if err != nil {
		return fmt.Errorf("Failed to parse format (%s), error: %s", formatStr, err)
	}
	if format != output.JSONFormat && format != output.YAMLFormat {
		return fmt.Errorf("Not allowed output format (%s), options: [%s, %s]", format.String(), output.YAMLFormat.String(), output.JSONFormat.String())
	}
	// ---

	scanResult := scanner.Config(searchDir)

	platforms := []string{}
	for platform := range scanResult.OptionsMap {
		platforms = append(platforms, platform)
	}

	if len(platforms) == 0 {
		cmd := command.New("which", "tree")
		out, err := cmd.RunAndReturnTrimmedCombinedOutput()
		if err != nil || out == "" {
			log.Errorf("tree not installed, can not list files")
		} else {
			fmt.Println()
			cmd := command.NewWithStandardOuts("tree", ".", "-L", "3")
			log.Printf("$ %s", cmd.PrintableCommandArgs())
			if err := cmd.Run(); err != nil {
				log.Errorf("Failed to list files in current directory, error: %s", err)
			}
		}

		scanResult.AddError("", "No known platform detected")
		if _, err := writeOutput(outputDir, scanResult, format); err != nil {
			log.Errorf("Failed to write output, error: %s", err)
		}
		return errors.New("No known platform detected")
	}

	// Write output to files
	if isCI {
		if _, err := writeOutput(outputDir, scanResult, format); err != nil {
			return fmt.Errorf("Failed to write output, error: %s", err)
		}
		return nil
	}
	// ---

	// Select option
	log.Infoft(colorstring.Blue("Collecting inputs:"))

	config, err := scanner.AskForConfig(scanResult)
	if err != nil {
		return err
	}

	if exist, err := pathutil.IsDirExists(outputDir); err != nil {
		return err
	} else if !exist {
		if err := os.MkdirAll(outputDir, 0700); err != nil {
			return fmt.Errorf("Failed to create (%s), error: %s", outputDir, err)
		}
	}

	pth := path.Join(outputDir, "bitrise.yml")
	outputPth, err := output.WriteToFile(config, format, pth)
	if err != nil {
		return fmt.Errorf("Failed to print result, error: %s", err)
	}
	log.Infoft("  bitrise.yml template: %s", colorstring.Blue(outputPth))
	fmt.Println()
	// ---

	return nil
}
