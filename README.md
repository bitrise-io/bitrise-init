# Bitrise Init Tool

Initialize bitrise config, step template or plugin template

## How to build this project 
Project is written in [Go](https://golang.org/) language and 
uses [godep](https://github.com/tools/godep) as dependency management tool.

You can build this project using sequence of `go` commands or refer to [bitrise.yml](./bitrise.yml) file,
which contains workflows for this project.

You can run `bitrise` workflows on your local machine using [bitrise CLI](https://www.bitrise.io/cli).

Before you start, make sure 
- `$HOME/go/bin` (or `$GOPATH/bin` in case of custom go workspace) is added to `$PATH`
- `Ruby >= 2.2.2` version is installed
- `bundler` gem installed

**How to build the project using bitrise workflows**

Please check available workflows in [bitrise.yml](./bitrise.yml). 
`bitrise --ci run ci` will execute `ci` workflow which consists of `prepare/build/run tests` stages.

**How to build the project using Go commands**
- `go build` command builds the project and generates `bitrise-init` in project root folder
- `go install` command installs `bitrise-init` binary at `$HOME/go/bin/bitrise-init` (or `$GOPATH/bin/bitrise-init` in case of custom go workspace).
- `go test ./...` command runs unit tests in every project folder/subfolder.
- `go test -v ./_tests/integration/...` command runs integration tests. This command requires `INTEGRATION_TEST_BINARY_PATH=$HOME/go/bin/bitrise-init` (or `INTEGRATION_TEST_BINARY_PATH=$GOPATH/bin/bitrise-init` in case of custom go workspace) environment variable.

**How to write integration test for a scanner**
Once new scanner is introduced there are at least 2 mandatory integration tests which should be added:
- integration test for new scanner on test project
- integration test for new scanner in manual config integration test.
This one already exists and requires adjustments for new scanner.

Standard integration for a scanner consists of verification of generated config with expected config on a test project.
Choose any project identifiable by a new scanner and create integration test for it under `_tests/integration/{SCANNER_NAME}_test.go`.
Any existing integration test for other scanners can be used as a starting point.
Integration test should contain expected config YML and it will be compared with a generated config YML (`bitrise-init --ci config` command) .
You can run this command on a test project on your local machine to understand test expectations.
Steps versions in expected config YML should be replaced with format placeholders like `%s` (take `_tests/integration/ionic_test.go` as an example).
If generated locally config on a test project doesn't match with expected config YML, then easiest way to fix this will be to use generated config YML file from integration test.
`gitClone(t, dir, URL)` command prints into console the path to a temp folder where `result.yml` will be generated.
Just copy and paste the content of `result.yml` to integration test and use placeholders for steps versions.

In case of a manual config integration test, the procedure is very similar: generate the config YML by running `bitrise-init --ci manual-config` command.
Alternatively, copy generated config YML from failed manual config integration test and update `_tests/integration/manual_config_test.go`

## How to release new bitrise-init version

- update the step versions in steps/const.go
- bump `RELEASE_VERSION` in bitrise.yml
- commit these changes
- call `bitrise run create-release`
- check and update the generated CHANGELOG.md
- test the generated binaries in _bin/ directory
- push these changes to the master branch
- once `create-release` workflow finishes on bitrise.io test the build generated binaries
- create a github release with the build generated binaries

__Update manual config on website__

- use the generated binaries in `./_bin/` directory to generate the manual config by calling: `BIN_PATH --ci manual-config` this will generate the manual.config.yml at: `CURRENT_DIR/_defaults/result.yml`
- throw the generated `result.yml` to the frontend team, to update the manual-config on the website
- once they put the new config in the website project, check the git changes to make sure, everything looks great

__Update the [project-scanner step](https://github.com/bitrise-steplib/steps-project-scanner)__

- update bitrise-init dependency
- share a new version into the steplib (check the [README.md](https://github.com/bitrise-steplib/steps-project-scanner/blob/master/README.md))

__Update the [bitrise init plugin]((https://github.com/bitrise-core/bitrise-plugins-init))__

- update bitrise-init dependency
- release a new version (check the [README.md](https://github.com/bitrise-core/bitrise-plugins-init/blob/master/README.md))