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

	systemDepsInstallScriptStepTitle = "Install system dependencies"

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
	gemName             string
	adapterName         string // Rails adapter name in database.yml (e.g. "postgresql", "mysql2")
	containerName       string
	image               string
	ports               []string
	containerEnvKey     string // env var name the container needs (e.g., POSTGRES_PASSWORD)
	healthCheck         string
	isRelationalDB      bool
	connectionURLEnvKey string   // app-level env var for the service URL (e.g., REDIS_URL)
	connectionURL       string   // value for connectionURLEnvKey (e.g., redis://localhost:6379/0)
	aptPackages         []string // system packages required to compile the gem's native extension
	// hostValue overrides the default "localhost" for DB_HOST. Use "127.0.0.1" for MySQL,
	// which treats "localhost" as a Unix socket path rather than a TCP address.
	hostValue string
}

var knownDatabaseGems = []databaseGem{
	{
		gemName:         "pg",
		adapterName:     "postgresql",
		containerName:   "postgres",
		image:           "postgres:17",
		ports:           []string{"5432:5432"},
		containerEnvKey: "POSTGRES_PASSWORD",
		healthCheck:     `--health-cmd "pg_isready" --health-interval 10s --health-timeout 5s --health-retries 5`,
		isRelationalDB:  true,
	},
	{
		gemName:         "mysql2",
		adapterName:     "mysql2",
		containerName:   "mysql",
		image:           "mysql:8",
		ports:           []string{"3306:3306"},
		containerEnvKey: "MYSQL_ROOT_PASSWORD",
		healthCheck:     `--health-cmd "mysqladmin ping -h 127.0.0.1 -u root --password=$$MYSQL_ROOT_PASSWORD" --health-interval 10s --health-timeout 5s --health-retries 5`,
		isRelationalDB:  true,
		aptPackages:     []string{"libmariadb-dev"},
		hostValue:       "127.0.0.1",
	},
	{
		gemName:             "redis",
		containerName:       "redis",
		image:               "redis:7",
		ports:               []string{"6379:6379"},
		healthCheck:         `--health-cmd "redis-cli ping" --health-interval 10s --health-timeout 5s --health-retries 5`,
		connectionURLEnvKey: "REDIS_URL",
		connectionURL:       "redis://localhost:6379/0",
	},
	{
		gemName:       "mongoid",
		containerName: "mongodb",
		image:         "mongo:8",
		ports:         []string{"27017:27017"},
		healthCheck:   `--health-cmd "mongosh --eval 'db.runCommand({ping:1})'" --health-interval 10s --health-timeout 5s --health-retries 5`,
	},
	{
		gemName:       "mongo",
		containerName: "mongodb",
		image:         "mongo:8",
		ports:         []string{"27017:27017"},
		healthCheck:   `--health-cmd "mongosh --eval 'db.runCommand({ping:1})'" --health-interval 10s --health-timeout 5s --health-retries 5`,
	},
	{
		// SQLite is file-based, no service container needed, but ActiveRecord setup is required
		gemName:        "sqlite3",
		isRelationalDB: true,
	},
}

// databaseYMLInfo holds env var names and defaults extracted from config/database.yml.
type databaseYMLInfo struct {
	adapter        string // e.g. "postgresql", "mysql2", "sqlite3"
	hostEnvVar     databaseEnvVar
	usernameEnvVar databaseEnvVar
	passwordEnvVar databaseEnvVar
}

// mongoidYMLInfo holds connection URL info extracted from config/mongoid.yml.
type mongoidYMLInfo struct {
	connectionURLEnvKey string // e.g. "MONGODB_URL"
	connectionURL       string // e.g. "mongodb://localhost:27017/myapp_test"
}

var (
	gemDeclPattern     = regexp.MustCompile(`^\s*gem\s+['"]([^'"]+)['"]`)
	envFetchPattern    = regexp.MustCompile(`ENV\.fetch\(\s*["'](\w+)["']\s*\)\s*\{\s*["']([^"']*)["']\s*\}`)
	envFetchArgPattern = regexp.MustCompile(`ENV\.fetch\(\s*['"](\w+)['"]\s*,\s*['"]([^'"]*)['"]\s*\)`)
	envBracketPattern  = regexp.MustCompile(`ENV\[["'](\w+)["']\]`)
	// erbTagPattern matches ERB template tags like <%= ... %> that appear in Rails database.yml.
	// It assumes the expression itself does not contain a bare '%>' sequence.
	erbTagPattern = regexp.MustCompile(`<%[^%]*%>`)
)

func detectRails(searchDir string) bool {
	gemfilePath := filepath.Join(searchDir, "Gemfile")
	content, err := fileutil.ReadStringFromFile(gemfilePath)
	if err != nil {
		return false
	}
	for _, line := range strings.Split(content, "\n") {
		match := gemDeclPattern.FindStringSubmatch(line)
		if len(match) >= 2 && match[1] == "rails" {
			return true
		}
	}
	return false
}

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
		dedupKey := dbGem.containerName
		if dedupKey == "" {
			dedupKey = dbGem.gemName
		}
		if declaredGems[dbGem.gemName] && !seen[dedupKey] {
			detected = append(detected, dbGem)
			seen[dedupKey] = true
		}
	}
	return detected
}

func parseDatabaseYML(searchDir string, databases []databaseGem) databaseYMLInfo {
	ymlPath := filepath.Join(searchDir, "config", "database.yml")
	content, err := fileutil.ReadStringFromFile(ymlPath)
	if err != nil {
		log.TPrintf("- config/database.yml - not found or not readable")
		return databaseYMLInfo{}
	}

	log.TPrintf("- config/database.yml - found, parsing credentials")
	return parseDatabaseYMLContent(content, databases)
}

// parseDatabaseYMLContent parses the contents of a database.yml file and extracts
// env-var references for the host, username, and password fields.
// It prefers the "test" environment section, then "default", then any other section.
// YAML anchor merges (<<: *default) are resolved automatically by the YAML parser.
// The adapter field is required: if absent or not matching a detected database gem, the result is empty.
func parseDatabaseYMLContent(content string, databases []databaseGem) databaseYMLInfo {
	preprocessed := preprocessERBForYAML(content)

	var rawYML map[string]map[string]interface{}
	if err := yaml.Unmarshal([]byte(preprocessed), &rawYML); err != nil {
		log.TWarnf("- config/database.yml - failed to parse: %s", err)
		return databaseYMLInfo{}
	}

	// Prefer "test", then "default", then the first available section.
	var section map[string]interface{}
	for _, name := range []string{"test", "default"} {
		if s, ok := rawYML[name]; ok {
			section = s
			break
		}
	}
	if section == nil {
		for _, s := range rawYML {
			section = s
			break
		}
	}
	if section == nil {
		return databaseYMLInfo{}
	}

	info := databaseYMLInfo{
		adapter:        asString(section["adapter"]),
		hostEnvVar:     extractEnvVarFromValue(asString(section["host"])),
		usernameEnvVar: extractEnvVarFromValue(asString(section["username"])),
		passwordEnvVar: extractEnvVarFromValue(asString(section["password"])),
	}

	if info.adapter == "" {
		log.TWarnf("database.yml has no adapter field, skipping database.yml config")
		return databaseYMLInfo{}
	}

	for _, db := range databases {
		if db.adapterName == info.adapter {
			return info
		}
	}

	log.TWarnf("database.yml adapter %q does not match any detected database gem, skipping database.yml config", info.adapter)
	return databaseYMLInfo{}
}

// preprocessERBForYAML wraps ERB template tags (e.g. <%= ENV.fetch(...) %>) in
// single quotes so that the surrounding YAML can be parsed by a standard YAML parser.
func preprocessERBForYAML(content string) string {
	return erbTagPattern.ReplaceAllStringFunc(content, func(match string) string {
		// Escape any single quotes inside the ERB expression (YAML single-quote escaping uses '').
		escaped := strings.ReplaceAll(match, "'", "''")
		return "'" + escaped + "'"
	})
}

// asString converts any value from yaml.Unmarshal to its string representation.
func asString(v interface{}) string {
	if v == nil {
		return ""
	}
	return fmt.Sprintf("%v", v)
}

// extractEnvVarFromValue extracts an env-var name and default from a YAML field value.
// It recognises ENV.fetch("KEY") { "default" }, ENV["KEY"], and plain string values.
func extractEnvVarFromValue(value string) databaseEnvVar {
	// ENV.fetch("KEY") { "default" }
	if match := envFetchPattern.FindStringSubmatch(value); len(match) >= 3 {
		return databaseEnvVar{name: match[1], defaultValue: match[2]}
	}
	// ENV["KEY"]
	if match := envBracketPattern.FindStringSubmatch(value); len(match) >= 2 {
		return databaseEnvVar{name: match[1], defaultValue: ""}
	}
	// Plain value (no ERB reference)
	if value != "" && !strings.Contains(value, "<%") {
		return databaseEnvVar{name: "", defaultValue: value}
	}
	return databaseEnvVar{}
}

func parseMongoidYML(searchDir string, mongoDB databaseGem) mongoidYMLInfo {
	ymlPath := filepath.Join(searchDir, "config", "mongoid.yml")
	content, err := fileutil.ReadStringFromFile(ymlPath)
	if err != nil {
		log.TPrintf("- config/mongoid.yml - not found or not readable")
		return mongoidYMLInfo{}
	}

	log.TPrintf("- config/mongoid.yml - found, parsing connection URL")
	return parseMongoidYMLContent(content)
}

func parseMongoidYMLContent(content string) mongoidYMLInfo {
	// Look for ENV.fetch('KEY', 'mongodb://...') pattern anywhere in the file
	match := envFetchArgPattern.FindStringSubmatch(content)
	if len(match) < 3 {
		return mongoidYMLInfo{}
	}

	envKey := match[1]
	defaultURL := match[2]

	// Script steps run on the host machine, not inside Docker, so they connect to service
	// containers via localhost (ports are mapped to the host).
	// Normalize any IP-based localhost references to the hostname form.
	connectionURL := strings.ReplaceAll(defaultURL, "127.0.0.1", "localhost")

	return mongoidYMLInfo{
		connectionURLEnvKey: envKey,
		connectionURL:       connectionURL,
	}
}

// findMongoDBGem returns the first detected non-relational DB gem that has a container (e.g. mongoid/mongo).
func findMongoDBGem(databases []databaseGem) (databaseGem, bool) {
	for _, db := range databases {
		if !db.isRelationalDB && db.containerName != "" {
			return db, true
		}
	}
	return databaseGem{}, false
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
	hasRails       bool
	isDefault      bool
	databases      []databaseGem
	dbYMLInfo      databaseYMLInfo
	mongoidYMLInfo mongoidYMLInfo
}

func createConfigDescriptor(project project, isDefault bool) configDescriptor {
	descriptor := configDescriptor{
		workdir:        "$" + projectDirInputEnvKey,
		hasBundler:     project.hasBundler,
		hasRakefile:    project.hasRakefile,
		testFramework:  project.testFramework,
		hasRubyVersion: project.hasRubyVersion,
		hasRails:       project.hasRails,
		isDefault:      isDefault,
		databases:      project.databases,
		dbYMLInfo:      project.dbYMLInfo,
		mongoidYMLInfo: project.mongoidYMLInfo,
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

	for _, db := range params.databases {
		if db.containerName != "" {
			name = name + "-" + db.containerName
		}
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

	// Install system dependencies (e.g. native library headers required by some gems)
	if aptPackages := collectAptPackages(descriptor.databases); len(aptPackages) > 0 {
		configBuilder.AppendStepListItemsTo(runTestsWorkflowID, steps.ScriptStepListItem(systemDepsInstallScriptStepTitle, generateSystemDepsScript(aptPackages)))
	}

	// Install dependencies
	if descriptor.hasBundler {
		configBuilder.AppendStepListItemsTo(runTestsWorkflowID, steps.ScriptStepListItem(bundlerInstallScriptStepTitle, bundlerInstallScriptStepContent, workdirInputs(descriptor.workdir)...))
	}

	serviceContainerNames := serviceContainerReferences(descriptor.databases)
	relationalServiceContainerNames := relationalServiceContainerReferences(descriptor.databases)

	// Database setup (only for relational DBs)
	if hasRelationalDB(descriptor.databases) {
		dbSetupScript := generateDBSetupScript(descriptor)
		if len(relationalServiceContainerNames) > 0 {
			configBuilder.AppendStepListItemsTo(runTestsWorkflowID, scriptStepWithServiceContainers("Database setup", dbSetupScript, relationalServiceContainerNames, descriptor.workdir))
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
	appEnvs := buildAppEnvs(descriptor.databases, descriptor.dbYMLInfo, descriptor.mongoidYMLInfo)

	if len(descriptor.databases) > 0 {
		containers := buildContainerDefinitions(descriptor.databases, descriptor.dbYMLInfo)
		if len(containers) > 0 {
			configBuilder.SetContainerDefinitions(containers)
		}
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
		if db.containerName != "" {
			refs = append(refs, db.containerName)
		}
	}
	return refs
}

func relationalServiceContainerReferences(databases []databaseGem) []stepmanModels.ContainerReference {
	var refs []stepmanModels.ContainerReference
	for _, db := range databases {
		if db.isRelationalDB && db.containerName != "" {
			refs = append(refs, db.containerName)
		}
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
		if db.containerName == "" {
			continue
		}
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

func buildAppEnvs(databases []databaseGem, ymlInfo databaseYMLInfo, mongoidInfo mongoidYMLInfo) []envmanModels.EnvironmentItemModel {
	hasRelational := hasRelationalDB(databases)
	hasMongoidURL := mongoidInfo.connectionURLEnvKey != ""
	if !hasRelational && !hasMongoidURL {
		return nil
	}

	var envs []envmanModels.EnvironmentItemModel

	if hasRelational {
		// Host env var: use name from database.yml or default to DB_HOST
		hostEnvName := "DB_HOST"
		if ymlInfo.hostEnvVar.name != "" {
			hostEnvName = ymlInfo.hostEnvVar.name
		}
		// Script steps run on the host machine, not inside Docker, so they connect to service
		// containers via mapped ports. Most databases work with "localhost", but MySQL treats
		// "localhost" as a Unix socket path — "127.0.0.1" forces TCP/IP.
		for _, db := range databases {
			if db.isRelationalDB && db.containerName != "" {
				hostValue := "localhost"
				if db.hostValue != "" {
					hostValue = db.hostValue
				}
				envs = append(envs, envmanModels.EnvironmentItemModel{hostEnvName: hostValue})
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

		// Connection URL env vars for databases with a standard URL convention (e.g. REDIS_URL)
		for _, db := range databases {
			if db.connectionURLEnvKey != "" {
				envs = append(envs, envmanModels.EnvironmentItemModel{db.connectionURLEnvKey: db.connectionURL})
			}
		}
	}

	// MongoDB connection URL parsed from config/mongoid.yml
	if hasMongoidURL {
		envs = append(envs, envmanModels.EnvironmentItemModel{mongoidInfo.connectionURLEnvKey: mongoidInfo.connectionURL})
	}

	return envs
}

// collectAptPackages returns the deduplicated list of apt packages required by the detected database gems.
func collectAptPackages(databases []databaseGem) []string {
	seen := map[string]bool{}
	var packages []string
	for _, db := range databases {
		for _, pkg := range db.aptPackages {
			if !seen[pkg] {
				seen[pkg] = true
				packages = append(packages, pkg)
			}
		}
	}
	return packages
}

func generateSystemDepsScript(packages []string) string {
	return "#!/usr/bin/env bash\nset -euxo pipefail\n\napt-get install -y " + strings.Join(packages, " ") + "\n"
}

func generateDBSetupScript(descriptor configDescriptor) string {
	runner := "rake"
	if descriptor.hasRails {
		runner = "rails"
	}
	dbCommand := runner + " db:create db:schema:load"
	if descriptor.hasBundler {
		dbCommand = "bundle exec " + runner + " db:create db:schema:load"
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
		if descriptor.hasRails {
			if descriptor.hasBundler {
				testCommand = "bundle exec rails test"
			} else {
				testCommand = "rails test"
			}
		} else if descriptor.hasRakefile {
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
