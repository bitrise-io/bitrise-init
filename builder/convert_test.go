package builder

import (
	"testing"

	"github.com/bitrise-io/bitrise-init/models"
	bitriseModels "github.com/bitrise-io/bitrise/models"
)

func TestSyntaxNode_Export(t *testing.T) {
	tests := []struct {
		name       string
		node       TemplateNode
		answerTree *AnswerTree
		want       *models.OptionNode
		want1      map[string]Result
		wantErr    bool
	}{
		{
			name: "One Step, no question",
			node: &Steps{
				Steps: []TemplateNode{
					&Step{
						ID: "fastlane",
						Inputs: []Input{
							{Key: "A", Value: &Text{Contents: "B"}},
						},
					},
				},
			},
			answerTree: nil,
			want: &models.OptionNode{
				Config:         "prefix",
				ChildOptionMap: map[string]*models.OptionNode{},
			},
			want1: map[string]Result{
				"prefix": {
					Config: &bitriseModels.BitriseDataModel{
						FormatVersion:        "11",
						DefaultStepLibSource: "https://github.com/bitrise-io/bitrise-steplib.git",
						TriggerMap: bitriseModels.TriggerMapModel{
							bitriseModels.TriggerMapItemModel{
								PushBranch: "*",
								WorkflowID: string(models.PrimaryWorkflowID),
							},
							bitriseModels.TriggerMapItemModel{
								PullRequestSourceBranch: "*",
								WorkflowID:              string(models.PrimaryWorkflowID),
							},
						},
						Workflows: map[string]bitriseModels.WorkflowModel{
							string(models.PrimaryWorkflowID): {
								Steps: []bitriseModels.StepListItemModel{
									newStep("fastlane", "", []Input{{Key: "A", Value: &Text{Contents: "B"}}}, ""),
								},
							},
						},
					},
				},
			},
		},
		{
			name: "2 questions",
			node: &Steps{
				Steps: []TemplateNode{
					&Step{
						ID: "fastlane",
						Inputs: []Input{
							{Key: "A", Value: &InputSelect{QuestionID: "question2"}},
						},
						// templateID: 1,
					},
					&Step{
						ID: "xcode-archive",
						Inputs: []Input{
							{Key: "export-method", Value: &InputSelect{QuestionID: "export-method"}},
						},
						// templateID: 2,
					},
				},
			},
			answerTree: &AnswerTree{
				Answer: AnswerExpansion{
					Key: "question2",
					Question: &Question{
						ID:         "question2",
						Title:      "title2",
						Type:       models.TypeSelector,
						Selections: []string{"n", "m"},
					},
					SelectionToExpandedValue: map[string]string{
						"n": "n",
						"m": "m",
					},
				},
				Children: map[string]*AnswerTree{
					"n": {
						Answer: AnswerExpansion{
							Key: "export-method",
							Question: &Question{
								ID:         "export-method",
								Title:      "title2",
								Type:       models.TypeSelector,
								Selections: []string{"development", "app-store"},
							},
							SelectionToExpandedValue: map[string]string{
								"development": "development",
								"app-store":   "app-store",
							},
						},
						Children: map[string]*AnswerTree{
							"development": nil,
							"app-store":   nil,
						},
					},
					"m": {
						Answer: AnswerExpansion{
							Key: "export-method",
							Question: &Question{
								ID:         "export-method",
								Title:      "title2",
								Type:       models.TypeSelector,
								Selections: []string{"development", "app-store"},
							},
							SelectionToExpandedValue: map[string]string{
								"development": "development",
								"app-store":   "app-store",
							},
						},
						Children: map[string]*AnswerTree{
							"development": nil,
							"app-store":   nil,
						},
					},
				},
			},
			want: &models.OptionNode{
				Title: "title2",
				Type:  models.TypeSelector,
				ChildOptionMap: map[string]*models.OptionNode{
					"n": {
						Title: "title2",
						Type:  models.TypeSelector,
						ChildOptionMap: map[string]*models.OptionNode{
							"development": {
								Config:         "prefix/question2.n/export-method.development",
								ChildOptionMap: map[string]*models.OptionNode{},
							},
							"app-store": {
								Config:         "prefix/question2.n/export-method.app-store",
								ChildOptionMap: map[string]*models.OptionNode{},
							},
						},
					},
					"m": {
						Title: "title2",
						Type:  models.TypeSelector,
						ChildOptionMap: map[string]*models.OptionNode{
							"development": {
								Config:         "prefix/question2.m/export-method.development",
								ChildOptionMap: map[string]*models.OptionNode{},
							},
							"app-store": {
								Config:         "prefix/question2.m/export-method.app-store",
								ChildOptionMap: map[string]*models.OptionNode{},
							},
						},
					},
				},
			},
			want1: map[string]Result{
				"prefix/question2.n/export-method.development": {
					Config: &bitriseModels.BitriseDataModel{
						FormatVersion:        "11",
						DefaultStepLibSource: "https://github.com/bitrise-io/bitrise-steplib.git",
						TriggerMap: bitriseModels.TriggerMapModel{
							bitriseModels.TriggerMapItemModel{
								PushBranch: "*",
								WorkflowID: string(models.PrimaryWorkflowID),
							},
							bitriseModels.TriggerMapItemModel{
								PullRequestSourceBranch: "*",
								WorkflowID:              string(models.PrimaryWorkflowID),
							},
						},
						Workflows: map[string]bitriseModels.WorkflowModel{
							string(models.PrimaryWorkflowID): {
								Steps: []bitriseModels.StepListItemModel{
									newStep("fastlane", "", []Input{{Key: "A", Value: &Text{Contents: "n"}}}, ""),
									newStep("xcode-archive", "", []Input{{Key: "export_method", Value: &Text{Contents: "development"}}}, ""),
								},
							},
						},
					},
				},
				"prefix/question2.n/export-method.app-store": {
					Config: &bitriseModels.BitriseDataModel{
						FormatVersion:        "11",
						DefaultStepLibSource: "https://github.com/bitrise-io/bitrise-steplib.git",
						TriggerMap: bitriseModels.TriggerMapModel{
							bitriseModels.TriggerMapItemModel{
								PushBranch: "*",
								WorkflowID: string(models.PrimaryWorkflowID),
							},
							bitriseModels.TriggerMapItemModel{
								PullRequestSourceBranch: "*",
								WorkflowID:              string(models.PrimaryWorkflowID),
							},
						},
						Workflows: map[string]bitriseModels.WorkflowModel{
							string(models.PrimaryWorkflowID): {
								Steps: []bitriseModels.StepListItemModel{
									newStep("fastlane", "", []Input{{Key: "A", Value: &Text{Contents: "n"}}}, ""),
									newStep("xcode-archive", "", []Input{{Key: "export_method", Value: &Text{Contents: "app-store"}}}, ""),
								},
							},
						},
					},
				},
				"prefix/question2.m/export-method.development": {
					Config: &bitriseModels.BitriseDataModel{
						FormatVersion:        "11",
						DefaultStepLibSource: "https://github.com/bitrise-io/bitrise-steplib.git",
						TriggerMap: bitriseModels.TriggerMapModel{
							bitriseModels.TriggerMapItemModel{
								PushBranch: "*",
								WorkflowID: string(models.PrimaryWorkflowID),
							},
							bitriseModels.TriggerMapItemModel{
								PullRequestSourceBranch: "*",
								WorkflowID:              string(models.PrimaryWorkflowID),
							},
						},
						Workflows: map[string]bitriseModels.WorkflowModel{
							string(models.PrimaryWorkflowID): {
								Steps: []bitriseModels.StepListItemModel{
									newStep("fastlane", "", []Input{{Key: "A", Value: &Text{Contents: "m"}}}, ""),
									newStep("xcode-archive", "", []Input{{Key: "export_method", Value: &Text{Contents: "development"}}}, ""),
								},
							},
						},
					},
				},
				"prefix/question2.m/export-method.app-store": {
					Config: &bitriseModels.BitriseDataModel{
						FormatVersion:        "11",
						DefaultStepLibSource: "https://github.com/bitrise-io/bitrise-steplib.git",
						TriggerMap: bitriseModels.TriggerMapModel{
							bitriseModels.TriggerMapItemModel{
								PushBranch: "*",
								WorkflowID: string(models.PrimaryWorkflowID),
							},
							bitriseModels.TriggerMapItemModel{
								PullRequestSourceBranch: "*",
								WorkflowID:              string(models.PrimaryWorkflowID),
							},
						},
						Workflows: map[string]bitriseModels.WorkflowModel{
							string(models.PrimaryWorkflowID): {
								Steps: []bitriseModels.StepListItemModel{
									newStep("fastlane", "", []Input{{Key: "A", Value: &Text{Contents: "m"}}}, ""),
									newStep("xcode-archive", "", []Input{{Key: "export_method", Value: &Text{Contents: "app-store"}}}, ""),
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
			got, got1, err := Export(tt.node, tt.answerTree, nil, "prefix")
			if (err != nil) != tt.wantErr {
				t.Errorf("SyntaxNode.Export() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			assertEqual(t, tt.want, got)
			assertEqual(t, tt.want1, got1)
		})
	}
}
