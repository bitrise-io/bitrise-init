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

func writeScanResult(scanResult models.ScanResultModel, outputDir string, format output.Format) (string, error) {
	pth := path.Join(outputDir, "result")
	return output.WriteToFile(scanResult, format, pth)
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

	// normalize working dir path
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

	// normalize output path
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

		log.Infoft("Saving outputs:")
		scanResult.AddError("general", "No known platform detected")
		outputPth, err := writeScanResult(scanResult, outputDir, format)
		if err != nil {
			log.Errorf("Failed to write output, error: %s", err)
		} else {
			log.Printft("  scan result: %s", outputPth)
		}

		return errors.New("No known platform detected")
	}

	// Write output to files
	if isCI {
		log.Infoft("Saving outputs:")

		outputPth, err := writeScanResult(scanResult, outputDir, format)
		if err != nil {
			return fmt.Errorf("Failed to write output, error: %s", err)
		}

		log.Printft("  scan result: %s", outputPth)
		return nil
	}
	// ---

	// Select option
	log.Infoft("Collecting inputs:")

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
	log.Infoft("  bitrise.yml template: %s", outputPth)
	fmt.Println()
	// ---

	return nil
}
