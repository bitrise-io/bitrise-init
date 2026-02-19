package ruby

import (
	"path/filepath"

	"github.com/bitrise-io/bitrise-init/models"
	"github.com/bitrise-io/bitrise-init/utility"
	"github.com/bitrise-io/go-utils/log"
)

// Options
const (
	ScannerName = "ruby"

	runTestsWorkflowID = models.WorkflowID("run_tests")

	projectDirInputTitle   = "Project Directory"
	projectDirInputSummary = "The directory containing the Gemfile"
	projectDirInputEnvKey  = "RUBY_PROJECT_DIR"

	bundlerInputTitle   = "Bundler"
	bundlerInputSummary = "Whether to use Bundler for dependency management"
)

type testFramework struct {
	name           string
	detectionFiles []string
}

var testFrameworks = []testFramework{
	{"rspec", []string{"spec/spec_helper.rb", ".rspec"}},
	{"minitest", []string{"test/test_helper.rb"}},
}

type project struct {
	projectRelDir  string
	hasBundler     bool
	hasRakefile    bool
	testFramework  string
	hasRubyVersion bool
}

// Scanner implements the Scanner interface for Ruby projects
type Scanner struct {
	projects []project
}

// NewScanner creates a new scanner instance.
func NewScanner() *Scanner {
	return &Scanner{}
}

// Name returns the name of the scanner
func (scanner *Scanner) Name() string {
	return ScannerName
}

// DetectPlatform checks if the given search directory contains a Ruby project
func (scanner *Scanner) DetectPlatform(searchDir string) (bool, error) {
	gemfilePaths, err := utility.FindFileInAppDir(searchDir, "Gemfile")
	if err != nil {
		log.TWarnf("%s", err)
		log.TPrintf("Platform not detected")
		return false, nil
	}

	for _, gemfilePath := range gemfilePaths {
		log.TPrintf("Checking: %s", gemfilePath)

		// determine workdir
		gemfileDir := filepath.Dir(gemfilePath)

		hasBundler := checkBundler(gemfileDir)
		hasRakefile := checkRakefile(gemfileDir)
		testFw := detectTestFramework(gemfileDir)
		hasRubyVersion := checkRubyVersion(gemfileDir)

		projectRelDir, err := utility.RelPath(searchDir, gemfileDir)
		if err != nil {
			log.TWarnf("failed to get relative Gemfile dir path: %s", err)
			continue
		}

		project := project{
			projectRelDir:  projectRelDir,
			hasBundler:     hasBundler,
			hasRakefile:    hasRakefile,
			testFramework:  testFw,
			hasRubyVersion: hasRubyVersion,
		}

		scanner.projects = append(scanner.projects, project)
	}

	if len(scanner.projects) == 0 {
		log.TPrintf("Platform not detected")
		return false, nil
	}

	log.TSuccessf("Platform detected")
	return true, nil
}

func (scanner *Scanner) ExcludedScannerNames() []string {
	return []string{}
}

// Options returns the options for the scanner
func (scanner *Scanner) Options() (models.OptionNode, models.Warnings, models.Icons, error) {
	return generateOptions(scanner.projects)
}

// Configs returns the default configurations for the scanner
func (scanner *Scanner) Configs(sshKeyActivation models.SSHKeyActivation) (models.BitriseConfigMap, error) {
	return generateConfigs(scanner.projects, sshKeyActivation)
}

// DefaultOptions returns the default options for the scanner
func (scanner *Scanner) DefaultOptions() models.OptionNode {
	projectRootOption := models.NewOption(projectDirInputTitle, projectDirInputSummary, projectDirInputEnvKey, models.TypeUserInput)

	defaultDescriptor := createDefaultConfigDescriptor()
	configOption := models.NewConfigOption(configName(defaultDescriptor), nil)
	projectRootOption.AddConfig(models.UserInputOptionDefaultValue, configOption)

	return *projectRootOption
}

func (scanner *Scanner) DefaultConfigs() (models.BitriseConfigMap, error) {
	configs := models.BitriseConfigMap{}

	defaultDescriptor := createDefaultConfigDescriptor()
	config, err := generateConfigBasedOn(defaultDescriptor, models.SSHKeyActivationConditional)
	if err != nil {
		return nil, err
	}
	configs[configName(defaultDescriptor)] = config

	return configs, nil
}
