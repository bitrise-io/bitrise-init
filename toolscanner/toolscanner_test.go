package toolscanner

import (
	"reflect"
	"testing"

	"github.com/bitrise-core/bitrise-init/models"
)

func TestAddProjectTypeToConfig(t *testing.T) {
	const detectedProject = "ios"
	const configName = "fastlane-config"
	type args struct {
		scannerConfigMap     models.BitriseConfigMap
		detectedProjectTypes []string
	}
	tests := []struct {
		name    string
		args    args
		want    models.BitriseConfigMap
		wantErr bool
	}{
		{
			name: "Ok case",
			args: args{
				scannerConfigMap: models.BitriseConfigMap{
					"fastlane-config": `format_version: "6" 
					default_step_lib_source: https://github.com/bitrise-io/bitrise-steplib.git
					project_type: '` + "[[.PROJECT_TYPE]]" + `'
					app:
					envs:
					- FASTLANE_XCODE_LIST_TIMEOUT: "120"
					trigger_map:
					- push_branch: '*'
					workflow: primary
					- pull_request_source_branch: '*'
					workflow: primary
					workflows:
					primary:
						steps:
						- activate-ssh-key@4.0.3:
							run_if: '{{getenv "SSH_RSA_PRIVATE_KEY" | ne ""}}'
						- git-clone@4.0.14: {}
						- script@1.1.5:
							title: Do anything with Script step
						- certificate-and-profile-installer@1.10.1: {}
						- fastlane@2.3.12:
							inputs:
							- lane: $FASTLANE_LANE
							- work_dir: $FASTLANE_WORK_DIR
						- deploy-to-bitrise-io@1.3.19: {}`,
				},
				detectedProjectTypes: []string{detectedProject},
			},
			want: models.BitriseConfigMap{
				configName + "_" + detectedProject: `format_version: "6" 
					default_step_lib_source: https://github.com/bitrise-io/bitrise-steplib.git
					project_type: '` + detectedProject + `'
					app:
					envs:
					- FASTLANE_XCODE_LIST_TIMEOUT: "120"
					trigger_map:
					- push_branch: '*'
					workflow: primary
					- pull_request_source_branch: '*'
					workflow: primary
					workflows:
					primary:
						steps:
						- activate-ssh-key@4.0.3:
							run_if: '{{getenv "SSH_RSA_PRIVATE_KEY" | ne ""}}'
						- git-clone@4.0.14: {}
						- script@1.1.5:
							title: Do anything with Script step
						- certificate-and-profile-installer@1.10.1: {}
						- fastlane@2.3.12:
							inputs:
							- lane: $FASTLANE_LANE
							- work_dir: $FASTLANE_WORK_DIR
						- deploy-to-bitrise-io@1.3.19: {}`,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := AddProjectTypeToConfig(tt.args.scannerConfigMap, tt.args.detectedProjectTypes)
			if (err != nil) != tt.wantErr {
				t.Errorf("AddProjectTypeToConfig() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("AddProjectTypeToConfig() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestAddProjectTypeToOptions(t *testing.T) {
	const detectedProjectType = "ios"
	type args struct {
		scannerOptionTree    models.OptionNode
		detectedProjectTypes []string
	}
	tests := []struct {
		name string
		args args
		want models.OptionNode
	}{
		{
			name: "Ok case",
			args: args{
				scannerOptionTree: models.OptionNode{
					Title:  "Working directory",
					EnvKey: "FASTLANE_WORK_DIR",
					ChildOptionMap: map[string]*models.OptionNode{
						"BitriseFastlaneSample": &models.OptionNode{
							Title:  "Fastlane lane",
							EnvKey: "FASTLANE_LANE",
							ChildOptionMap: map[string]*models.OptionNode{
								"ios test": &models.OptionNode{
									Config: "fastlane-config",
								},
							},
						},
					},
				},
				detectedProjectTypes: []string{detectedProjectType},
			},
			want: models.OptionNode{
				Title:  "Project type",
				EnvKey: "PROJECT_TYPE",
				ChildOptionMap: map[string]*models.OptionNode{
					detectedProjectType: &models.OptionNode{
						Title:  "Working directory",
						EnvKey: "FASTLANE_WORK_DIR",
						ChildOptionMap: map[string]*models.OptionNode{
							"BitriseFastlaneSample": &models.OptionNode{
								Title:  "Fastlane lane",
								EnvKey: "FASTLANE_LANE",
								ChildOptionMap: map[string]*models.OptionNode{
									"ios test": &models.OptionNode{
										Config: "fastlane-config" + "_" + detectedProjectType,
									},
								},
							},
						},
					},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := AddProjectTypeToOptions(tt.args.scannerOptionTree, tt.args.detectedProjectTypes); !reflect.DeepEqual(got.String(), tt.want.String()) {
				t.Errorf("AddProjectTypeToOptions() = %+v, want %v", got, tt.want)
			}
		})
	}
}
