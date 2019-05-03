package cli

import (
	"fmt"
	"os"
	"path"
	"path/filepath"

	"github.com/bitrise-io/bitrise-init/models"
	"github.com/bitrise-io/bitrise-init/output"
	"github.com/bitrise-io/bitrise-init/scanner"
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
			log.TErrorf(err.Error())
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

func printDirTree() {
	cmd := command.New("which", "tree")
	out, err := cmd.RunAndReturnTrimmedCombinedOutput()
	if err != nil || out == "" {
		log.TErrorf("tree not installed, can not list files")
	} else {
		fmt.Println()
		cmd := command.NewWithStandardOuts("tree", ".", "-L", "3")
		log.TPrintf("$ %s", cmd.PrintableCommandArgs())
		if err := cmd.Run(); err != nil {
			log.TErrorf("Failed to list files in current directory, error: %s", err)
		}
	}
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

	//
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

	if isCI {
		log.TInfof(colorstring.Yellow("CI mode"))
	}
	log.TInfof(colorstring.Yellowf("scan dir: %s", searchDir))
	log.TInfof(colorstring.Yellowf("output dir: %s", outputDir))
	log.TInfof(colorstring.Yellowf("output format: %s", format))
	fmt.Println()

	result, err := GenerateAndWriteResults(searchDir, outputDir, format)
	if err != nil {
		return err
	}

	if !isCI {
		if err := getInteractiveAnswers(result, outputDir, format); err != nil {
			return nil
		}
	}
	return nil
}

// GenerateAndWriteResults runs the scanner and saves results to the given output dir
func GenerateAndWriteResults(searchDir string, outputDir string, format output.Format) (models.ScanResultModel, error) {
	result, detected, err := generateConfig(searchDir, format)
	if err != nil {
		return result, err
	}

	// Write output to files
	log.TInfof("Saving outputs:")
	outputPth, err := writeScanResult(result, outputDir, format)
	if err != nil {
		return result, fmt.Errorf("Failed to write output, error: %s", err)
	}
	log.TPrintf("scan result: %s", outputPth)

	if !detected {
		printDirTree()
		return result, fmt.Errorf("No known platform detected")
	}
	return result, nil
}

func generateConfig(searchDir string, format output.Format) (models.ScanResultModel, bool, error) {
	scanResult := scanner.Config(searchDir)

	platforms := []string{}
	for platform := range scanResult.ScannerToOptionRoot {
		platforms = append(platforms, platform)
	}

	if len(platforms) == 0 {
		scanResult.AddError("general", "No known platform detected")
		return scanResult, false, nil
	}
	return scanResult, true, nil
}

func getInteractiveAnswers(scanResult models.ScanResultModel, outputDir string, format output.Format) error {
	// Select options
	log.TInfof("Collecting inputs:")
	config, err := scanner.AskForConfig(scanResult)
	if err != nil {
		return err
	}

	pth := path.Join(outputDir, "bitrise.yml")
	outputPth, err := output.WriteToFile(config, format, pth)
	if err != nil {
		return fmt.Errorf("Failed to print result, error: %s", err)
	}
	log.TInfof("  bitrise.yml template: %s", outputPth)
	fmt.Println()
	return nil
}
