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
		{"multi-package", "https://github.com/bitrise-io/nestjs-sample-01-cats-app", "multi-package", multiPackageYml, multiPackageYmlVersions},
	}

	helper.Execute(t, testCases)
}

// Expected results
var multiPackageYmlVersions = []interface{}{
	models.FormatVersion,
	steps.ActivateSSHKeyVersion,
	steps.GitCloneVersion,
	steps.ScriptVersion,
	steps.CacheRestoreNPMVersion,
	steps.NpmVersion,
	steps.NpmVersion,
	steps.NpmVersion,
	steps.CacheSaveNPMVersion,

	models.FormatVersion,
	steps.ActivateSSHKeyVersion,
	steps.GitCloneVersion,
	steps.ScriptVersion,
	steps.CacheRestoreNPMVersion,
	steps.NpmVersion,
	steps.NpmVersion,
	steps.NpmVersion,
	steps.CacheSaveNPMVersion,

	models.FormatVersion,
	steps.ActivateSSHKeyVersion,
	steps.GitCloneVersion,
	steps.ScriptVersion,
	steps.CacheRestoreNPMVersion,
	steps.NpmVersion,
	steps.NpmVersion,
	steps.NpmVersion,
	steps.CacheSaveNPMVersion,

	models.FormatVersion,
	steps.ActivateSSHKeyVersion,
	steps.GitCloneVersion,
	steps.ScriptVersion,
	steps.CacheRestoreNPMVersion,
	steps.YarnVersion,
	steps.YarnVersion,
	steps.YarnVersion,
	steps.CacheSaveNPMVersion,

	models.FormatVersion,
	steps.ActivateSSHKeyVersion,
	steps.GitCloneVersion,
	steps.ScriptVersion,
	steps.CacheRestoreNPMVersion,
	steps.YarnVersion,
	steps.YarnVersion,
	steps.YarnVersion,
	steps.CacheSaveNPMVersion,

	models.FormatVersion,
	steps.ActivateSSHKeyVersion,
	steps.GitCloneVersion,
	steps.NpmVersion,
	steps.CacheRestoreNPMVersion,
	steps.YarnVersion,
	steps.YarnVersion,
	steps.YarnVersion,
	steps.CacheSaveNPMVersion,
}

var multiPackageYml = fmt.Sprintf(`options:
  node-js:
    title: Project Directory
    summary: The directory containing the package.json file
    env_key: NODEJS_PROJECT_DIR
    type: selector
    value_map:
      .:
        title: Node Version
        summary: The version of Node.js used in the project. Leave it empty to use
          the latest Node version
        env_key: NODEJS_VERSION
        type: user_input_optional
        value_map:
          "":
            title: Package Manager
            summary: The package manager used in the project
            type: selector
            value_map:
              npm:
                config: node-js-npm-root-build-lint-test-config
              yarn:
                config: node-js-yarn-root-build-lint-test-config
      other-projects/node-version:
        title: Node Version
        summary: The version of Node.js used in the project. Leave it empty to use
          the latest Node version
        env_key: NODEJS_VERSION
        type: selector
        value_map:
          "22":
            title: Package Manager
            summary: The package manager used in the project
            type: selector
            value_map:
              npm:
                config: node-js-npm-build-lint-test-config
      other-projects/nvmrc:
        title: Node Version
        summary: The version of Node.js used in the project. Leave it empty to use
          the latest Node version
        env_key: NODEJS_VERSION
        type: selector
        value_map:
          lts:
            title: Package Manager
            summary: The package manager used in the project
            type: selector
            value_map:
              npm:
                config: node-js-npm-nvm-build-lint-test-config
      other-projects/tool-versions:
        title: Node Version
        summary: The version of Node.js used in the project. Leave it empty to use
          the latest Node version
        env_key: NODEJS_VERSION
        type: selector
        value_map:
          "20":
            title: Package Manager
            summary: The package manager used in the project
            type: selector
            value_map:
              npm:
                config: node-js-npm-build-lint-test-config
      other-projects/yarn:
        title: Node Version
        summary: The version of Node.js used in the project. Leave it empty to use
          the latest Node version
        env_key: NODEJS_VERSION
        type: selector
        value_map:
          "20":
            title: Package Manager
            summary: The package manager used in the project
            type: selector
            value_map:
              yarn:
                config: node-js-yarn-build-lint-test-config
      other-projects/yarn-nvmrc:
        title: Node Version
        summary: The version of Node.js used in the project. Leave it empty to use
          the latest Node version
        env_key: NODEJS_VERSION
        type: selector
        value_map:
          lts:
            title: Package Manager
            summary: The package manager used in the project
            type: selector
            value_map:
              yarn:
                config: node-js-yarn-nvm-build-lint-test-config
configs:
  node-js:
    node-js-npm-build-lint-test-config: |
      format_version: "%s"
      default_step_lib_source: https://github.com/bitrise-io/bitrise-steplib.git
      project_type: node-js
      workflows:
        run_tests:
          steps:
          - activate-ssh-key@%s: {}
          - git-clone@%s: {}
          - script@%s:
              inputs:
              - title: Install Node.js
              - content: |
                  #!/usr/bin/env bash
                  set -euxo pipefail

                  export ASDF_NODEJS_LEGACY_FILE_DYNAMIC_STRATEGY=latest_installed
                  envman add --key ASDF_NODEJS_LEGACY_FILE_DYNAMIC_STRATEGY --value latest_installed

                  pushd "${NODEJS_PROJECT_DIR:-.}" > /dev/null

                  # Bitrise stacks come with asdf pre-installed to help auto-switch between various software versions
                  # asdf looks for the Node.js version in these files: .tool-versions, .nvmrc, .node-version
                  # so it should work out-of-the-box even if the project uses another Node.js manager
                  # See: https://github.com/asdf-vm/asdf-nodejs
                  asdf install nodejs

                  popd > /dev/null
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
    node-js-npm-nvm-build-lint-test-config: |
      format_version: "%s"
      default_step_lib_source: https://github.com/bitrise-io/bitrise-steplib.git
      project_type: node-js
      workflows:
        run_tests:
          steps:
          - activate-ssh-key@%s: {}
          - git-clone@%s: {}
          - script@%s:
              inputs:
              - title: Install Node.js
              - content: |
                  #!/usr/bin/env bash
                  set -euxo pipefail

                  export ASDF_NODEJS_LEGACY_FILE_DYNAMIC_STRATEGY=latest_installed
                  envman add --key ASDF_NODEJS_LEGACY_FILE_DYNAMIC_STRATEGY --value latest_installed

                  pushd "${NODEJS_PROJECT_DIR:-.}" > /dev/null

                  # Bitrise stacks come with asdf pre-installed to help auto-switch between various software versions
                  # asdf looks for the Node.js version in these files: .tool-versions, .nvmrc, .node-version
                  # so it should work out-of-the-box even if the project uses another Node.js manager
                  # See: https://github.com/asdf-vm/asdf-nodejs
                  asdf install nodejs

                  popd > /dev/null
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
    node-js-npm-root-build-lint-test-config: |
      format_version: "%s"
      default_step_lib_source: https://github.com/bitrise-io/bitrise-steplib.git
      project_type: node-js
      workflows:
        run_tests:
          steps:
          - activate-ssh-key@%s: {}
          - git-clone@%s: {}
          - script@%s:
              inputs:
              - title: Install Node.js
              - content: |
                  #!/usr/bin/env bash
                  set -euxo pipefail

                  export ASDF_NODEJS_LEGACY_FILE_DYNAMIC_STRATEGY=latest_installed
                  envman add --key ASDF_NODEJS_LEGACY_FILE_DYNAMIC_STRATEGY --value latest_installed

                  pushd "${NODEJS_PROJECT_DIR:-.}" > /dev/null

                  # Bitrise stacks come with asdf pre-installed to help auto-switch between various software versions
                  # asdf looks for the Node.js version in these files: .tool-versions, .nvmrc, .node-version
                  # so it should work out-of-the-box even if the project uses another Node.js manager
                  # See: https://github.com/asdf-vm/asdf-nodejs
                  asdf install nodejs

                  popd > /dev/null
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
    node-js-yarn-build-lint-test-config: |
      format_version: "%s"
      default_step_lib_source: https://github.com/bitrise-io/bitrise-steplib.git
      project_type: node-js
      workflows:
        run_tests:
          steps:
          - activate-ssh-key@%s: {}
          - git-clone@%s: {}
          - script@%s:
              inputs:
              - title: Install Node.js
              - content: |
                  #!/usr/bin/env bash
                  set -euxo pipefail

                  export ASDF_NODEJS_LEGACY_FILE_DYNAMIC_STRATEGY=latest_installed
                  envman add --key ASDF_NODEJS_LEGACY_FILE_DYNAMIC_STRATEGY --value latest_installed

                  pushd "${NODEJS_PROJECT_DIR:-.}" > /dev/null

                  # Bitrise stacks come with asdf pre-installed to help auto-switch between various software versions
                  # asdf looks for the Node.js version in these files: .tool-versions, .nvmrc, .node-version
                  # so it should work out-of-the-box even if the project uses another Node.js manager
                  # See: https://github.com/asdf-vm/asdf-nodejs
                  asdf install nodejs

                  popd > /dev/null
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
              title: yarn run test
              inputs:
              - workdir: $NODEJS_PROJECT_DIR
              - command: run test
          - save-npm-cache@%s: {}
    node-js-yarn-nvm-build-lint-test-config: |
      format_version: "%s"
      default_step_lib_source: https://github.com/bitrise-io/bitrise-steplib.git
      project_type: node-js
      workflows:
        run_tests:
          steps:
          - activate-ssh-key@%s: {}
          - git-clone@%s: {}
          - script@%s:
              inputs:
              - title: Install Node.js
              - content: |
                  #!/usr/bin/env bash
                  set -euxo pipefail

                  export ASDF_NODEJS_LEGACY_FILE_DYNAMIC_STRATEGY=latest_installed
                  envman add --key ASDF_NODEJS_LEGACY_FILE_DYNAMIC_STRATEGY --value latest_installed

                  pushd "${NODEJS_PROJECT_DIR:-.}" > /dev/null

                  # Bitrise stacks come with asdf pre-installed to help auto-switch between various software versions
                  # asdf looks for the Node.js version in these files: .tool-versions, .nvmrc, .node-version
                  # so it should work out-of-the-box even if the project uses another Node.js manager
                  # See: https://github.com/asdf-vm/asdf-nodejs
                  asdf install nodejs

                  popd > /dev/null
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
              title: yarn run test
              inputs:
              - workdir: $NODEJS_PROJECT_DIR
              - command: run test
          - save-npm-cache@%s: {}
    node-js-yarn-root-build-lint-test-config: |
      format_version: "%s"
      default_step_lib_source: https://github.com/bitrise-io/bitrise-steplib.git
      project_type: node-js
      workflows:
        run_tests:
          steps:
          - activate-ssh-key@%s: {}
          - git-clone@%s: {}
          - script@%s:
              inputs:
              - title: Install Node.js
              - content: |
                  #!/usr/bin/env bash
                  set -euxo pipefail

                  export ASDF_NODEJS_LEGACY_FILE_DYNAMIC_STRATEGY=latest_installed
                  envman add --key ASDF_NODEJS_LEGACY_FILE_DYNAMIC_STRATEGY --value latest_installed

                  pushd "${NODEJS_PROJECT_DIR:-.}" > /dev/null

                  # Bitrise stacks come with asdf pre-installed to help auto-switch between various software versions
                  # asdf looks for the Node.js version in these files: .tool-versions, .nvmrc, .node-version
                  # so it should work out-of-the-box even if the project uses another Node.js manager
                  # See: https://github.com/asdf-vm/asdf-nodejs
                  asdf install nodejs

                  popd > /dev/null
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
              title: yarn run test
              inputs:
              - command: run test
          - save-npm-cache@%s: {}
warnings:
  node-js: []
warnings_with_recommendations:
  node-js: []`, multiPackageYmlVersions...)
