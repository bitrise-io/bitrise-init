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

	result, err := scanner.GenerateAndWriteResults(searchDir, outputDir, format)
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
