package golang

import (
	"path/filepath"

	"github.com/bitrise-core/bitrise-init/models"
	"github.com/bitrise-core/bitrise-init/steps"
	"github.com/pkg/errors"
	yaml "gopkg.in/yaml.v2"
)

const scannerName = "go"

const (
	configName = "go-config"
)

// Scanner ...
type Scanner struct {
}

// NewScanner ...
func NewScanner() *Scanner {
	return &Scanner{}
}

// Name ...
func (Scanner) Name() string {
	return scannerName
}

// DetectPlatform ...
func (scanner *Scanner) DetectPlatform(searchDir string) (bool, error) {
	matches, err := filepath.Glob(filepath.Clean(searchDir) + "/*.go")
	if err != nil {
		return false, errors.WithStack(err)
	}
	anyGoFileFound := matches != nil
	return anyGoFileFound, nil
}

// ExcludedScannerNames ...
func (*Scanner) ExcludedScannerNames() []string {
	return []string{}
}

// Options ...
func (scanner *Scanner) Options() (models.OptionModel, models.Warnings, error) {
	return scanner.DefaultOptions(), models.Warnings{}, nil
}

// DefaultOptions ...
func (*Scanner) DefaultOptions() models.OptionModel {
	return models.OptionModel{
		Title:  "_",
		EnvKey: "_",
		Config: configName,
	}
}

// Configs ...
func (*Scanner) Configs() (models.BitriseConfigMap, error) {
	return confGen()
}

// DefaultConfigs ...
func (*Scanner) DefaultConfigs() (models.BitriseConfigMap, error) {
	return confGen()
}

func confGen() (models.BitriseConfigMap, error) {
	configBuilder := models.NewDefaultConfigBuilder()

	configBuilder.AppendStepListItemsTo(models.PrimaryWorkflowID, steps.DefaultPrepareStepList(false)...)

	config, err := configBuilder.Generate(scannerName)
	if err != nil {
		return models.BitriseConfigMap{}, err
	}

	data, err := yaml.Marshal(config)
	if err != nil {
		return models.BitriseConfigMap{}, err
	}

	return models.BitriseConfigMap{
		configName: string(data),
	}, nil
}
