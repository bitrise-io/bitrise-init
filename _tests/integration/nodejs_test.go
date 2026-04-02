package integration

import (
	"fmt"
	"testing"

	"github.com/bitrise-io/bitrise-init/_tests/integration/helper"
	"github.com/bitrise-io/bitrise-init/models"
	"github.com/bitrise-io/bitrise-init/steps"
)

func TestNodeJs(t *testing.T) {
	testCases := []helper.TestCase{
		{
			Name:              "nextjs-npm",
			RepoURL:           "https://github.com/bitrise-io/nodejs-samples.git",
			RelativeSearchDir: "nextjs-npm",
			Branch:            "main",
			ExpectedResult:    nextjsNpmResultYML,
			ExpectedVersions:  nextjsNpmResultVersions,
		},
		{
			Name:              "nextjs-yarn",
			RepoURL:           "https://github.com/bitrise-io/nodejs-samples.git",
			RelativeSearchDir: "nextjs-yarn",
			Branch:            "main",
			ExpectedResult:    nextjsYarnResultYML,
			ExpectedVersions:  nextjsYarnResultVersions,
		},
		{
			Name:              "nestjs-cats-app",
			RepoURL:           "https://github.com/bitrise-io/nodejs-samples.git",
			RelativeSearchDir: "nestjs-cats-app",
			Branch:            "main",
			ExpectedResult:    nestjsCatsAppResultYML,
			ExpectedVersions:  nestjsCatsAppResultVersions,
		},
		{
			Name:             "nodejs-samples",
			RepoURL:          "https://github.com/bitrise-io/nodejs-samples.git",
			Branch:           "main",
			ExpectedResult:   nodejsSamplesResultYML,
			ExpectedVersions: nodejsSamplesResultVersions,
		},
	}

	helper.Execute(t, testCases)
}

// nextjs-npm: Next.js with npm, .nvmrc for Node version, has lint + test scripts.

var nextjsNpmResultVersions = []interface{}{
	models.FormatVersion,
	steps.ActivateSSHKeyVersion,
	steps.GitCloneVersion,
	steps.CacheRestoreNPMVersion,
	steps.NpmVersion,
	steps.NpmVersion,
	steps.NpmVersion,
	steps.CacheSaveNPMVersion,
}

var nextjsNpmResultYML = fmt.Sprintf(`options:
  node-js:
    title: Project Directory
    summary: The directory containing the package.json file
    env_key: NODEJS_PROJECT_DIR
    type: selector
    value_map:
      .:
        title: Package Manager
        summary: The package manager used in the project
        type: selector
        value_map:
          npm:
            config: node-js-nextjs-npm-root-lint-test-config
configs:
  node-js:
    node-js-nextjs-npm-root-lint-test-config: |
      format_version: "%s"
      default_step_lib_source: https://github.com/bitrise-io/bitrise-steplib.git
      project_type: node-js
      workflows:
        run_tests:
          steps:
          - activate-ssh-key@%s: {}
          - git-clone@%s: {}
          - restore-npm-cache@%s: {}
          - npm@%s:
              title: npm install
              inputs:
              - command: install
          - npm@%s:
              title: npm run lint
              inputs:
              - command: run lint
          - npm@%s:
              title: npm run test
              inputs:
              - command: run test
          - save-npm-cache@%s: {}
      tools:
        node: "22"
warnings:
  node-js: []
warnings_with_recommendations:
  node-js: []`, nextjsNpmResultVersions...)

// nextjs-yarn: Next.js with Yarn, engines.node for Node version, has lint + build scripts (no test).

var nextjsYarnResultVersions = []interface{}{
	models.FormatVersion,
	steps.ActivateSSHKeyVersion,
	steps.GitCloneVersion,
	steps.CacheRestoreNPMVersion,
	steps.YarnVersion,
	steps.YarnVersion,
	steps.YarnVersion,
	steps.CacheSaveNPMVersion,
}

var nextjsYarnResultYML = fmt.Sprintf(`options:
  node-js:
    title: Project Directory
    summary: The directory containing the package.json file
    env_key: NODEJS_PROJECT_DIR
    type: selector
    value_map:
      .:
        title: Package Manager
        summary: The package manager used in the project
        type: selector
        value_map:
          yarn:
            config: node-js-nextjs-yarn-root-lint-build-config
configs:
  node-js:
    node-js-nextjs-yarn-root-lint-build-config: |
      format_version: "%s"
      default_step_lib_source: https://github.com/bitrise-io/bitrise-steplib.git
      project_type: node-js
      workflows:
        run_tests:
          steps:
          - activate-ssh-key@%s: {}
          - git-clone@%s: {}
          - restore-npm-cache@%s: {}
          - yarn@%s:
              title: yarn install
              inputs:
              - command: install
          - yarn@%s:
              title: yarn run lint
              inputs:
              - command: run lint
          - yarn@%s:
              title: yarn run build
              inputs:
              - command: run build
          - save-npm-cache@%s: {}
      tools:
        node: 22.0.0
warnings:
  node-js: []
warnings_with_recommendations:
  node-js: []`, nextjsYarnResultVersions...)

// nestjs-cats-app: NestJS with npm, .tool-versions for Node version, has lint + test scripts.

var nestjsCatsAppResultVersions = []interface{}{
	models.FormatVersion,
	steps.ActivateSSHKeyVersion,
	steps.GitCloneVersion,
	steps.CacheRestoreNPMVersion,
	steps.NpmVersion,
	steps.NpmVersion,
	steps.NpmVersion,
	steps.CacheSaveNPMVersion,
}

var nestjsCatsAppResultYML = fmt.Sprintf(`options:
  node-js:
    title: Project Directory
    summary: The directory containing the package.json file
    env_key: NODEJS_PROJECT_DIR
    type: selector
    value_map:
      .:
        title: Package Manager
        summary: The package manager used in the project
        type: selector
        value_map:
          npm:
            config: node-js-nestjs-npm-root-build-lint-test-config
configs:
  node-js:
    node-js-nestjs-npm-root-build-lint-test-config: |
      format_version: "%s"
      default_step_lib_source: https://github.com/bitrise-io/bitrise-steplib.git
      project_type: node-js
      workflows:
        run_tests:
          steps:
          - activate-ssh-key@%s: {}
          - git-clone@%s: {}
          - restore-npm-cache@%s: {}
          - npm@%s:
              title: npm install
              inputs:
              - command: install
          - npm@%s:
              title: npm run lint
              inputs:
              - command: run lint
          - npm@%s:
              title: npm run test
              inputs:
              - command: run test
          - save-npm-cache@%s: {}
      tools:
        node: 22.14.0
warnings:
  node-js: []
warnings_with_recommendations:
  node-js: []`, nestjsCatsAppResultVersions...)

// nodejs-samples: full repo scan — all 4 projects.
// nestjs-cats-app and nestjs-node-version share the same config name; the last-written config
// (nestjs-node-version, node: "22") is what ends up in the configs map.

var nodejsSamplesResultVersions = []interface{}{
	// node-js-nestjs-npm-build-lint-test-config (nestjs-node-version, node: "22")
	models.FormatVersion,
	steps.ActivateSSHKeyVersion,
	steps.GitCloneVersion,
	steps.CacheRestoreNPMVersion,
	steps.NpmVersion,
	steps.NpmVersion,
	steps.NpmVersion,
	steps.CacheSaveNPMVersion,
	// node-js-nextjs-npm-lint-test-config
	models.FormatVersion,
	steps.ActivateSSHKeyVersion,
	steps.GitCloneVersion,
	steps.CacheRestoreNPMVersion,
	steps.NpmVersion,
	steps.NpmVersion,
	steps.NpmVersion,
	steps.CacheSaveNPMVersion,
	// node-js-nextjs-yarn-lint-build-config
	models.FormatVersion,
	steps.ActivateSSHKeyVersion,
	steps.GitCloneVersion,
	steps.CacheRestoreNPMVersion,
	steps.YarnVersion,
	steps.YarnVersion,
	steps.YarnVersion,
	steps.CacheSaveNPMVersion,
}

var nodejsSamplesResultYML = fmt.Sprintf(`options:
  node-js:
    title: Project Directory
    summary: The directory containing the package.json file
    env_key: NODEJS_PROJECT_DIR
    type: selector
    value_map:
      nestjs-cats-app:
        title: Package Manager
        summary: The package manager used in the project
        type: selector
        value_map:
          npm:
            config: node-js-nestjs-npm-build-lint-test-config
      nestjs-node-version:
        title: Package Manager
        summary: The package manager used in the project
        type: selector
        value_map:
          npm:
            config: node-js-nestjs-npm-build-lint-test-config
      nextjs-npm:
        title: Package Manager
        summary: The package manager used in the project
        type: selector
        value_map:
          npm:
            config: node-js-nextjs-npm-lint-test-config
      nextjs-yarn:
        title: Package Manager
        summary: The package manager used in the project
        type: selector
        value_map:
          yarn:
            config: node-js-nextjs-yarn-lint-build-config
configs:
  node-js:
    node-js-nestjs-npm-build-lint-test-config: |
      format_version: "%s"
      default_step_lib_source: https://github.com/bitrise-io/bitrise-steplib.git
      project_type: node-js
      workflows:
        run_tests:
          steps:
          - activate-ssh-key@%s: {}
          - git-clone@%s: {}
          - restore-npm-cache@%s: {}
          - npm@%s:
              title: npm install
              inputs:
              - workdir: $NODEJS_PROJECT_DIR
              - command: install
          - npm@%s:
              title: npm run lint
              inputs:
              - workdir: $NODEJS_PROJECT_DIR
              - command: run lint
          - npm@%s:
              title: npm run test
              inputs:
              - workdir: $NODEJS_PROJECT_DIR
              - command: run test
          - save-npm-cache@%s: {}
      tools:
        node: 22.14.0
    node-js-nextjs-npm-lint-test-config: |
      format_version: "%s"
      default_step_lib_source: https://github.com/bitrise-io/bitrise-steplib.git
      project_type: node-js
      workflows:
        run_tests:
          steps:
          - activate-ssh-key@%s: {}
          - git-clone@%s: {}
          - restore-npm-cache@%s: {}
          - npm@%s:
              title: npm install
              inputs:
              - workdir: $NODEJS_PROJECT_DIR
              - command: install
          - npm@%s:
              title: npm run lint
              inputs:
              - workdir: $NODEJS_PROJECT_DIR
              - command: run lint
          - npm@%s:
              title: npm run test
              inputs:
              - workdir: $NODEJS_PROJECT_DIR
              - command: run test
          - save-npm-cache@%s: {}
      tools:
        node: "22"
    node-js-nextjs-yarn-lint-build-config: |
      format_version: "%s"
      default_step_lib_source: https://github.com/bitrise-io/bitrise-steplib.git
      project_type: node-js
      workflows:
        run_tests:
          steps:
          - activate-ssh-key@%s: {}
          - git-clone@%s: {}
          - restore-npm-cache@%s: {}
          - yarn@%s:
              title: yarn install
              inputs:
              - workdir: $NODEJS_PROJECT_DIR
              - command: install
          - yarn@%s:
              title: yarn run lint
              inputs:
              - workdir: $NODEJS_PROJECT_DIR
              - command: run lint
          - yarn@%s:
              title: yarn run build
              inputs:
              - workdir: $NODEJS_PROJECT_DIR
              - command: run build
          - save-npm-cache@%s: {}
      tools:
        node: 22.0.0
warnings:
  node-js: []
warnings_with_recommendations:
  node-js: []`, nodejsSamplesResultVersions...)
