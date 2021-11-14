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
							{Key: "A", Value: "B"},
						},
					},
				},
			},
			answerTree: nil,
			want: &models.OptionNode{
				Config:         "",
				ChildOptionMap: map[string]*models.OptionNode{},
			},
			want1: map[string]Result{
				"": {
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
									newStep("fastlane", "", []Input{{Key: "A", Value: "B"}}, ""),
								},
							},
						},
					},
				},
			},
		},
		/*{
			name: "One question",
			node: &SyntaxNode{
				Type: QuestionNode,
				Question: Question{
					ID:     "question",
					Title:  "title",
					Type:   models.TypeOptionalUserInput,
					EnvKey: "TEST_KEY",
				},
				AnswerToNode: map[string]*SyntaxNode{
					"": {
						Type: StepListNode,
						Steps: Steps{
							Steps: []*SyntaxNode{{
								Type: StepNode,
								Step: Step{
									ID: "fastlane",
									Inputs: []Input{
										{Key: "A", Value: "D"},
										{Key: "B", Value: "$TEST_KEY"},
									},
								},
							}},
						},
					},
				},
			},
			want: &models.OptionNode{
				Title:  "title",
				Type:   models.TypeOptionalUserInput,
				EnvKey: "TEST_KEY",
				ChildOptionMap: map[string]*models.OptionNode{
					"": {
						Config:         "config/question.",
						ChildOptionMap: map[string]*models.OptionNode{},
					},
				},
			},
			want1: map[string]Result{
				"config/question.": {
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
									newStep("fastlane", "",
										[]Input{
											{Key: "A", Value: "D"},
											{Key: "B", Value: "$TEST_KEY"},
										}, ""),
								},
							},
						},
					},
				},
			},
		},*/
		{
			name: "2 questions",
			node: &Steps{
				Steps: []TemplateNode{
					&Step{
						ID: "fastlane",
						Inputs: []Input{
							{Key: "A", Value: `{{askForInputValue "question2"}}`},
						},
						templateID: 1,
					},
					&Step{
						ID: "xcode-archive",
						Inputs: []Input{
							{Key: "export_method", Value: `{{askForInputValue "export-method"}}`},
						},
						templateID: 2,
					},
				},
			},
			answerTree: &AnswerTree{
				Answer: AnswerExpansion{
					Key: AnswerKey{
						nodeID:  1,
						NodeKey: "A",
					},
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
							Key: AnswerKey{
								nodeID:  2,
								NodeKey: "export_method",
							},
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
							Key: AnswerKey{
								nodeID:  2,
								NodeKey: "export_method",
							},
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
								Config:         "/question2.n/export-method.development",
								ChildOptionMap: map[string]*models.OptionNode{},
							},
							"app-store": {
								Config:         "/question2.n/export-method.app-store",
								ChildOptionMap: map[string]*models.OptionNode{},
							},
						},
					},
					"m": {
						Title: "title2",
						Type:  models.TypeSelector,
						ChildOptionMap: map[string]*models.OptionNode{
							"development": {
								Config:         "/question2.m/export-method.development",
								ChildOptionMap: map[string]*models.OptionNode{},
							},
							"app-store": {
								Config:         "/question2.m/export-method.app-store",
								ChildOptionMap: map[string]*models.OptionNode{},
							},
						},
					},
				},
			},
			want1: map[string]Result{
				"/question2.n/export-method.development": {
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
									newStep("fastlane", "", []Input{{Key: "A", Value: "n"}}, ""),
									newStep("xcode-archive", "", []Input{{Key: "export_method", Value: "development"}}, ""),
								},
							},
						},
					},
				},
				"/question2.n/export-method.app-store": {
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
									newStep("fastlane", "", []Input{{Key: "A", Value: "n"}}, ""),
									newStep("xcode-archive", "", []Input{{Key: "export_method", Value: "app-store"}}, ""),
								},
							},
						},
					},
				},
				"/question2.m/export-method.development": {
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
									newStep("fastlane", "", []Input{{Key: "A", Value: "m"}}, ""),
									newStep("xcode-archive", "", []Input{{Key: "export_method", Value: "development"}}, ""),
								},
							},
						},
					},
				},
				"/question2.m/export-method.app-store": {
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
									newStep("fastlane", "", []Input{{Key: "A", Value: "m"}}, ""),
									newStep("xcode-archive", "", []Input{{Key: "export_method", Value: "app-store"}}, ""),
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
			got, got1, err := Export(tt.node, tt.answerTree, nil)
			if (err != nil) != tt.wantErr {
				t.Errorf("SyntaxNode.Export() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			assertEqual(t, tt.want, got)
			assertEqual(t, tt.want1, got1)
		})
	}
}
