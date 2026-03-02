package ruby

import (
	"fmt"
	"path/filepath"
	"regexp"
	"strings"

	"gopkg.in/yaml.v2"

	bitriseModels "github.com/bitrise-io/bitrise/v2/models"
	envmanModels "github.com/bitrise-io/envman/v2/models"
	"github.com/bitrise-io/go-utils/fileutil"
	"github.com/bitrise-io/go-utils/log"
	"github.com/bitrise-io/go-utils/pointers"
	stepmanModels "github.com/bitrise-io/stepman/models"

	"github.com/bitrise-io/bitrise-init/models"
	"github.com/bitrise-io/bitrise-init/steps"
	"github.com/bitrise-io/bitrise-init/utility"
)

const (
	gemCachePaths = "vendor/bundle"
	gemCacheKey   = `gem-{{ checksum "Gemfile.lock" }}`
)

const (
	rubyInstallScriptStepTitle   = "Install Ruby"
	rubyInstallScriptStepContent = `#!/usr/bin/env bash
set -euxo pipefail

# Bitrise stacks come with asdf pre-installed to help auto-switch between various software versions
# asdf looks for the Ruby version in these files: .tool-versions, .ruby-version
# See: https://github.com/asdf-vm/asdf-ruby
asdf install ruby
`

	bundlerInstallScriptStepTitle   = "Install dependencies"
	bundlerInstallScriptStepContent = `#!/usr/bin/env bash
set -euxo pipefail

bundle config set --local path vendor/bundle
bundle install
`
)

func checkBundler(searchDir string) bool {
	log.TPrintf("Checking for Bundler")
	hasGemfileLock := utility.FileExists(filepath.Join(searchDir, "Gemfile.lock"))

	if !hasGemfileLock {
		log.TPrintf("- Gemfile.lock - not found")
		return false
	}

	log.TPrintf("- Gemfile.lock - found")
	log.TPrintf("Bundler: detected")
	return true
}

func checkRakefile(searchDir string) bool {
	log.TPrintf("Checking for Rakefile")
	hasRakefile := utility.FileExists(filepath.Join(searchDir, "Rakefile"))

	if !hasRakefile {
		log.TPrintf("- Rakefile - not found")
		return false
	}

	log.TPrintf("- Rakefile - found")
	return true
}

func checkRubyVersion(searchDir string) bool {
	log.TPrintf("Checking for Ruby version file")

	versionFiles := []string{".ruby-version", ".tool-versions"}
	for _, versionFile := range versionFiles {
		if utility.FileExists(filepath.Join(searchDir, versionFile)) {
			log.TPrintf("- %s - found", versionFile)
			return true
		}
	}

	log.TPrintf("- Ruby version file - not found")
	return false
}

func detectTestFramework(searchDir string) string {
	log.TPrintf("Checking test framework")

	for _, fw := range testFrameworks {
		for _, detectionFile := range fw.detectionFiles {
			if utility.FileExists(filepath.Join(searchDir, detectionFile)) {
				log.TPrintf("- %s - found (%s)", fw.name, detectionFile)
				return fw.name
			}
		}
	}

	log.TPrintf("- test framework - not detected")
	return ""
}

// Database detection

// databaseEnvVar represents an environment variable with its name and default value.
type databaseEnvVar struct {
	name         string
	defaultValue string
}

// databaseGem represents a detected database dependency and its container configuration.
type databaseGem struct {
	gemName         string
	containerName   string
	image           string
	ports           []string
	containerEnvKey string // env var name the container needs (e.g., POSTGRES_PASSWORD)
	healthCheck     string
	isRelationalDB  bool
}

var knownDatabaseGems = []databaseGem{
	{
		gemName:         "pg",
		containerName:   "postgres",
		image:           "postgres:17",
		ports:           []string{"5432:5432"},
		containerEnvKey: "POSTGRES_PASSWORD",
		healthCheck:     `--health-cmd "pg_isready" --health-interval 10s --health-timeout 5s --health-retries 5`,
		isRelationalDB:  true,
	},
	{
		gemName:         "mysql2",
		containerName:   "mysql",
		image:           "mysql:8",
		ports:           []string{"3306:3306"},
		containerEnvKey: "MYSQL_ROOT_PASSWORD",
		healthCheck:     `--health-cmd "mysqladmin ping -h localhost" --health-interval 10s --health-timeout 5s --health-retries 5`,
		isRelationalDB:  true,
	},
	{
		gemName:       "redis",
		containerName: "redis",
		image:         "redis:7",
		ports:         []string{"6379:6379"},
		healthCheck:   `--health-cmd "redis-cli ping" --health-interval 10s --health-timeout 5s --health-retries 5`,
	},
	{
		gemName:       "mongoid",
		containerName: "mongo",
		image:         "mongo:7",
		ports:         []string{"27017:27017"},
		healthCheck:   `--health-cmd "mongosh --eval 'db.runCommand({ping:1})'" --health-interval 10s --health-timeout 5s --health-retries 5`,
	},
	{
		gemName:       "mongo",
		containerName: "mongo",
		image:         "mongo:7",
		ports:         []string{"27017:27017"},
		healthCheck:   `--health-cmd "mongosh --eval 'db.runCommand({ping:1})'" --health-interval 10s --health-timeout 5s --health-retries 5`,
	},
}

// databaseYMLInfo holds env var names and defaults extracted from config/database.yml.
type databaseYMLInfo struct {
	hostEnvVar     databaseEnvVar
	usernameEnvVar databaseEnvVar
	passwordEnvVar databaseEnvVar
}

var (
	gemDeclPattern    = regexp.MustCompile(`^\s*gem\s+['"]([^'"]+)['"]`)
	envFetchPattern   = regexp.MustCompile(`ENV\.fetch\(\s*["'](\w+)["']\s*\)\s*\{\s*["']([^"']*)["']\s*\}`)
	envBracketPattern = regexp.MustCompile(`ENV\[["'](\w+)["']\]`)
)

func detectDatabases(searchDir string) []databaseGem {
	gemfilePath := filepath.Join(searchDir, "Gemfile")
	content, err := fileutil.ReadStringFromFile(gemfilePath)
	if err != nil {
		log.TWarnf("Failed to read Gemfile: %s", err)
		return nil
	}

	databases := detectDatabaseGemsFromContent(content)
	return databases
}

func detectDatabaseGemsFromContent(content string) []databaseGem {
	declaredGems := map[string]bool{}
	for _, line := range strings.Split(content, "\n") {
		match := gemDeclPattern.FindStringSubmatch(line)
		if len(match) >= 2 {
			declaredGems[match[1]] = true
		}
	}

	var detected []databaseGem
	seen := map[string]bool{}
	for _, dbGem := range knownDatabaseGems {
		if declaredGems[dbGem.gemName] && !seen[dbGem.containerName] {
			detected = append(detected, dbGem)
			seen[dbGem.containerName] = true
		}
	}
	return detected
}

func parseDatabaseYML(searchDir string) databaseYMLInfo {
	info := databaseYMLInfo{}

	ymlPath := filepath.Join(searchDir, "config", "database.yml")
	content, err := fileutil.ReadStringFromFile(ymlPath)
	if err != nil {
		log.TPrintf("- config/database.yml - not found or not readable")
		return info
	}

	log.TPrintf("- config/database.yml - found, parsing credentials")

	// Extract the test section if available, otherwise use the whole content.
	// We look for lines after "test:" until the next top-level section.
	section := extractYMLSection(content, "test")
	if section == "" {
		section = content
	}

	info.hostEnvVar = extractEnvVarFromYMLField(section, "host")
	info.usernameEnvVar = extractEnvVarFromYMLField(section, "username")
	info.passwordEnvVar = extractEnvVarFromYMLField(section, "password")

	return info
}

// extractYMLSection extracts the content of a top-level YAML section (e.g., "test:").
func extractYMLSection(content, sectionName string) string {
	lines := strings.Split(content, "\n")
	var result []string
	inSection := false

	for _, line := range lines {
		trimmed := strings.TrimSpace(line)

		// Check if this is the target section start
		if !inSection {
			if trimmed == sectionName+":" || strings.HasPrefix(trimmed, sectionName+":") {
				inSection = true
			}
			continue
		}

		// If we hit another top-level key (not indented, not a comment, not empty), stop
		if len(line) > 0 && line[0] != ' ' && line[0] != '\t' && line[0] != '#' {
			break
		}

		result = append(result, line)
	}

	return strings.Join(result, "\n")
}

// extractEnvVarFromYMLField finds an env var reference in a YAML field line.
func extractEnvVarFromYMLField(section, fieldName string) databaseEnvVar {
	for _, line := range strings.Split(section, "\n") {
		trimmed := strings.TrimSpace(line)
		if !strings.HasPrefix(trimmed, fieldName+":") {
			continue
		}

		// Try ENV.fetch("KEY") { "default" }
		if match := envFetchPattern.FindStringSubmatch(line); len(match) >= 3 {
			return databaseEnvVar{name: match[1], defaultValue: match[2]}
		}

		// Try ENV["KEY"] (no default)
		if match := envBracketPattern.FindStringSubmatch(line); len(match) >= 2 {
			return databaseEnvVar{name: match[1], defaultValue: ""}
		}

		// Plain value (no env var reference) — extract the value after the colon
		parts := strings.SplitN(trimmed, ":", 2)
		if len(parts) == 2 {
			val := strings.TrimSpace(parts[1])
			if val != "" && !strings.Contains(val, "<%") {
				return databaseEnvVar{name: "", defaultValue: val}
			}
		}
	}

	return databaseEnvVar{}
}

// hasRelationalDB returns true if any detected database is relational (pg, mysql).
func hasRelationalDB(databases []databaseGem) bool {
	for _, db := range databases {
		if db.isRelationalDB {
			return true
		}
	}
	return false
}

// Options & Configs
type configDescriptor struct {
	workdir        string
	hasBundler     bool
	hasRakefile    bool
	testFramework  string
	hasRubyVersion bool
	isDefault      bool
	databases      []databaseGem
	dbYMLInfo      databaseYMLInfo
}

func createConfigDescriptor(project project, isDefault bool) configDescriptor {
	descriptor := configDescriptor{
		workdir:        "$" + projectDirInputEnvKey,
		hasBundler:     project.hasBundler,
		hasRakefile:    project.hasRakefile,
		testFramework:  project.testFramework,
		hasRubyVersion: project.hasRubyVersion,
		isDefault:      isDefault,
		databases:      project.databases,
		dbYMLInfo:      project.dbYMLInfo,
	}

	// Gemfile placed in the search dir, no need to change-dir
	if project.projectRelDir == "." {
		descriptor.workdir = ""
	}

	return descriptor
}

func createDefaultConfigDescriptor() configDescriptor {
	return createConfigDescriptor(project{
		projectRelDir:  "$" + projectDirInputEnvKey,
		hasBundler:     true,
		hasRakefile:    true,
		testFramework:  "rspec",
		hasRubyVersion: true,
	}, true)
}

func configName(params configDescriptor) string {
	name := "ruby"

	if params.isDefault {
		return "default-" + name + "-config"
	}

	if params.workdir == "" {
		name = name + "-root"
	}

	if params.hasBundler {
		name = name + "-bundler"
	}

	if params.testFramework != "" {
		name = name + "-" + params.testFramework
	}

	if len(params.databases) > 0 {
		name = name + "-" + params.databases[0].containerName
	}

	return name + "-config"
}

func generateOptions(projects []project) (models.OptionNode, models.Warnings, models.Icons, error) {
	if len(projects) == 0 {
		return models.OptionNode{}, nil, nil, fmt.Errorf("no Gemfile files found")
	}

	projectRootOption := models.NewOption(projectDirInputTitle, projectDirInputSummary, projectDirInputEnvKey, models.TypeSelector)
	for _, project := range projects {
		descriptor := createConfigDescriptor(project, false)
		configOption := models.NewConfigOption(configName(descriptor), nil)
		projectRootOption.AddConfig(project.projectRelDir, configOption)
	}

	return *projectRootOption, nil, nil, nil
}

func generateConfigs(projects []project, sshKeyActivation models.SSHKeyActivation) (models.BitriseConfigMap, error) {
	configs := models.BitriseConfigMap{}

	if len(projects) == 0 {
		return models.BitriseConfigMap{}, fmt.Errorf("no Gemfile files found")
	}

	for _, project := range projects {
		descriptor := createConfigDescriptor(project, false)
		config, err := generateConfigBasedOn(descriptor, sshKeyActivation)
		if err != nil {
			return nil, err
		}
		configs[configName(descriptor)] = config
	}

	return configs, nil
}

func generateConfigBasedOn(descriptor configDescriptor, sshKey models.SSHKeyActivation) (string, error) {
	configBuilder := models.NewDefaultConfigBuilder()
	prepareSteps := steps.DefaultPrepareStepList(steps.PrepareListParams{SSHKeyActivation: sshKey})
	configBuilder.AppendStepListItemsTo(runTestsWorkflowID, prepareSteps...)

	// Install Ruby
	if descriptor.hasRubyVersion {
		configBuilder.AppendStepListItemsTo(runTestsWorkflowID, steps.ScriptStepListItem(rubyInstallScriptStepTitle, rubyInstallScriptStepContent, workdirInputs(descriptor.workdir)...))
	}

	// Restore gem cache
	configBuilder.AppendStepListItemsTo(runTestsWorkflowID, steps.RestoreCache(gemCacheKey))

	// Install dependencies
	if descriptor.hasBundler {
		configBuilder.AppendStepListItemsTo(runTestsWorkflowID, steps.ScriptStepListItem(bundlerInstallScriptStepTitle, bundlerInstallScriptStepContent, workdirInputs(descriptor.workdir)...))
	}

	serviceContainerNames := serviceContainerReferences(descriptor.databases)

	// Database setup (only for relational DBs)
	if hasRelationalDB(descriptor.databases) {
		dbSetupScript := generateDBSetupScript(descriptor)
		if len(serviceContainerNames) > 0 {
			configBuilder.AppendStepListItemsTo(runTestsWorkflowID, scriptStepWithServiceContainers("Database setup", dbSetupScript, serviceContainerNames, descriptor.workdir))
		} else {
			configBuilder.AppendStepListItemsTo(runTestsWorkflowID, steps.ScriptStepListItem("Database setup", dbSetupScript, workdirInputs(descriptor.workdir)...))
		}
	}

	// Run tests based on detected framework
	testScript := generateTestScript(descriptor)
	if testScript != "" {
		if len(serviceContainerNames) > 0 {
			configBuilder.AppendStepListItemsTo(runTestsWorkflowID, scriptStepWithServiceContainers("Run tests", testScript, serviceContainerNames, descriptor.workdir))
		} else {
			configBuilder.AppendStepListItemsTo(runTestsWorkflowID, steps.ScriptStepListItem("Run tests", testScript, workdirInputs(descriptor.workdir)...))
		}
	}

	// Save gem cache
	configBuilder.AppendStepListItemsTo(runTestsWorkflowID, steps.SaveCache(gemCacheKey, gemCachePaths))

	// Deploy steps
	configBuilder.AppendStepListItemsTo(runTestsWorkflowID, steps.DefaultDeployStepList()...)

	// Build app-level env vars for database connections
	appEnvs := buildAppEnvs(descriptor.databases, descriptor.dbYMLInfo)

	if len(descriptor.databases) > 0 {
		containers := buildContainerDefinitions(descriptor.databases, descriptor.dbYMLInfo)
		configBuilder.SetContainerDefinitions(containers)
	}

	config, err := configBuilder.Generate(ScannerName, appEnvs...)
	if err != nil {
		return "", err
	}

	data, err := yaml.Marshal(config)
	if err != nil {
		return "", err
	}

	return string(data), nil
}

func serviceContainerReferences(databases []databaseGem) []stepmanModels.ContainerReference {
	var refs []stepmanModels.ContainerReference
	for _, db := range databases {
		refs = append(refs, db.containerName)
	}
	return refs
}

func workdirInputs(workdir string) []envmanModels.EnvironmentItemModel {
	if workdir == "" {
		return nil
	}
	return []envmanModels.EnvironmentItemModel{{"working_dir": workdir}}
}

func scriptStepWithServiceContainers(title, content string, serviceContainerRefs []stepmanModels.ContainerReference, workdir string) bitriseModels.StepListItemModel {
	stepID := steps.ScriptID + "@" + steps.ScriptVersion
	inputs := []envmanModels.EnvironmentItemModel{{"content": content}}
	inputs = append(inputs, workdirInputs(workdir)...)
	step := stepmanModels.StepModel{
		Title:             pointers.NewStringPtr(title),
		Inputs:            inputs,
		ServiceContainers: serviceContainerRefs,
	}
	return bitriseModels.StepListItemModel{stepID: step}
}

func buildContainerDefinitions(databases []databaseGem, ymlInfo databaseYMLInfo) map[string]bitriseModels.Container {
	containers := map[string]bitriseModels.Container{}
	for _, db := range databases {
		def := bitriseModels.Container{
			Type:    "service",
			Image:   db.image,
			Ports:   db.ports,
			Options: db.healthCheck,
		}

		// Set container env var referencing the app-level env var
		if db.containerEnvKey != "" && ymlInfo.passwordEnvVar.name != "" {
			def.Envs = []envmanModels.EnvironmentItemModel{
				{db.containerEnvKey: "$" + ymlInfo.passwordEnvVar.name},
			}
		}

		containers[db.containerName] = def
	}
	return containers
}

func buildAppEnvs(databases []databaseGem, ymlInfo databaseYMLInfo) []envmanModels.EnvironmentItemModel {
	var envs []envmanModels.EnvironmentItemModel

	hasRelational := false
	for _, db := range databases {
		if db.isRelationalDB {
			hasRelational = true
			break
		}
	}

	if !hasRelational {
		return nil
	}

	// Host env var: use name from database.yml or default to DB_HOST
	hostEnvName := "DB_HOST"
	if ymlInfo.hostEnvVar.name != "" {
		hostEnvName = ymlInfo.hostEnvVar.name
	}
	// Default value is the container name of the first relational DB
	for _, db := range databases {
		if db.isRelationalDB {
			envs = append(envs, envmanModels.EnvironmentItemModel{hostEnvName: db.containerName})
			break
		}
	}

	// Username env var
	if ymlInfo.usernameEnvVar.name != "" {
		envs = append(envs, envmanModels.EnvironmentItemModel{ymlInfo.usernameEnvVar.name: ymlInfo.usernameEnvVar.defaultValue})
	}

	// Password env var
	if ymlInfo.passwordEnvVar.name != "" {
		envs = append(envs, envmanModels.EnvironmentItemModel{ymlInfo.passwordEnvVar.name: ymlInfo.passwordEnvVar.defaultValue})
	}

	return envs
}

func generateDBSetupScript(descriptor configDescriptor) string {
	dbCommand := "rake db:create db:schema:load"
	if descriptor.hasBundler {
		dbCommand = "bundle exec rake db:create db:schema:load"
	}

	return fmt.Sprintf(`#!/usr/bin/env bash
set -euxo pipefail

%s`, dbCommand)
}

func generateTestScript(descriptor configDescriptor) string {
	testCommand := ""
	switch descriptor.testFramework {
	case "rspec":
		if descriptor.hasBundler {
			testCommand = "bundle exec rspec"
		} else {
			testCommand = "rspec"
		}
	case "minitest":
		if descriptor.hasRakefile {
			if descriptor.hasBundler {
				testCommand = "bundle exec rake test"
			} else {
				testCommand = "rake test"
			}
		} else {
			if descriptor.hasBundler {
				testCommand = "bundle exec ruby -Itest test/**/*_test.rb"
			} else {
				testCommand = "ruby -Itest test/**/*_test.rb"
			}
		}
	default:
		// Default to rake if Rakefile exists
		if descriptor.hasRakefile {
			if descriptor.hasBundler {
				testCommand = "bundle exec rake test"
			} else {
				testCommand = "rake test"
			}
		} else {
			return ""
		}
	}

	return fmt.Sprintf(`#!/usr/bin/env bash
set -euxo pipefail

%s`, testCommand)
}
