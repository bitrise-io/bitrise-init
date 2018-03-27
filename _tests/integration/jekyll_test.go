package integration

import (
	"fmt"
	"path/filepath"
	"strings"
	"testing"

	"github.com/bitrise-core/bitrise-init/models"
	"github.com/bitrise-core/bitrise-init/steps"
	"github.com/bitrise-io/go-utils/command"
	"github.com/bitrise-io/go-utils/fileutil"
	"github.com/bitrise-io/go-utils/pathutil"
	"github.com/stretchr/testify/require"
)

func TestJekyll(t *testing.T) {
	tmpDir, err := pathutil.NormalizedOSTempDirPath("__jekyll__")
	require.NoError(t, err)

	t.Log("jekyll-uno")
	{
		sampleAppDir := filepath.Join(tmpDir, "jekyll")
		sampleAppURL := "https://github.com/vgaidarji/jekyll-uno.git"
		gitClone(t, sampleAppDir, sampleAppURL)

		cmd := command.New(binPath(), "--ci", "config", "--dir", sampleAppDir, "--output-dir", sampleAppDir)
		out, err := cmd.RunAndReturnTrimmedCombinedOutput()
		require.NoError(t, err, out)

		scanResultPth := filepath.Join(sampleAppDir, "result.yml")

		result, err := fileutil.ReadStringFromFile(scanResultPth)
		require.NoError(t, err)
		require.Equal(t, strings.TrimSpace(jekyllResultYML), strings.TrimSpace(result))
	}
}

var jekyllVersions = []interface{}{
	models.FormatVersion,
	steps.ActivateSSHKeyVersion,
	steps.GitCloneVersion,
	steps.ScriptVersion,
	steps.CertificateAndProfileInstallerVersion,
	steps.NpmVersion,
	steps.GenerateCordovaBuildConfigVersion,
	steps.IonicArchiveVersion,
	steps.DeployToBitriseIoVersion,
}

// TODO use correct YML (copy-pasted from ionic as a starting point)
var jekyllResultYML = fmt.Sprintf(`2
          ios,android:
            config: ionic-config
configs:
  ionic:
    ionic-config: |
      format_version: "%s"
      default_step_lib_source: https://github.com/bitrise-io/bitrise-steplib.git
      project_type: ionic
      trigger_map:
      - push_branch: '*'
        workflow: primary
      - pull_request_source_branch: '*'
        workflow: primary
      workflows:
        primary:
          steps:
          - activate-ssh-key@%s:
              run_if: '{{getenv "SSH_RSA_PRIVATE_KEY" | ne ""}}'
          - git-clone@%s: {}
          - script@%s:
              title: Do anything with Script step
          - certificate-and-profile-installer@%s: {}
          - npm@%s:
              inputs:
              - workdir: $IONIC_WORK_DIR
              - command: install
          - generate-cordova-build-configuration@%s: {}
          - ionic-archive@%s:
              inputs:
              - platform: $IONIC_PLATFORM
              - target: emulator
              - workdir: $IONIC_WORK_DIR
          - deploy-to-bitrise-io@%s: {}
warnings:
  ionic: []
`, jekyllVersions...)
