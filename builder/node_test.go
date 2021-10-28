package builder

import (
	"testing"

	"github.com/bitrise-io/bitrise-init/models"
	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
)

func assertEqual(t *testing.T, want, got interface{}) {
	opts := []cmp.Option{
		cmpopts.IgnoreUnexported(models.OptionNode{}),
		// cmpopts.IgnoreUnexported(AnswerTree{}),
		cmpopts.IgnoreUnexported(AnswerExpansion{}),
		cmpopts.IgnoreUnexported(AnswerKey{}),
		// cmp.Comparer(func(x, y AnswerKey) bool { return x.NodeKey == y.NodeKey }),
	}

	if !cmp.Equal(want, got, opts...) {
		t.Fatalf("Not equal:\n%s", cmp.Diff(want, got, opts...))
	}
}

/*
func TestSyntaxNode_Evaluate(t *testing.T) {
	tests := []struct {
		name      string
		node      *SyntaxNode
		values    map[string]string
		questions map[string]Question
		want      *AnswerTree
		wantErr   bool
	}{
		// Empty StepList
		{
			name: "Single Step, no template",
			node: &SyntaxNode{
				Type: StepNode,
				Step: Step{
					ID: "fastlane",
					Inputs: []Input{
						{Key: "A", Value: "B"},
					},
				},
			},
			want: &AnswerTree{
				Content: &SyntaxNode{
					Type: StepNode,
					Step: Step{
						ID: "fastlane",
						Inputs: []Input{
							{Key: "A", Value: "B"},
						},
					},
				},
			},
		},
		{
			name: "Single Step, template, optional question",
			node: &SyntaxNode{
				Type: StepNode,
				Step: Step{
					ID: "fastlane",
					Inputs: []Input{
						{Key: "A", Value: "{{.C}}"},
						{Key: "B", Value: `{{askForInputValue "test_question"}}`},
					},
				},
			},
			values: map[string]string{
				"C": "D",
			},
			questions: map[string]Question{
				"test_question": {
					Title:  "title",
					Type:   models.TypeOptionalUserInput,
					EnvKey: "TEST_KEY",
				},
			},
			want: &AnswerTree{
				Answer: Answer2{
					key: AnswerKey{
						NodeKey: "B",
					},
					Question: &Question{
						ID:     "test_question",
						Title:  "title",
						Type:   models.TypeOptionalUserInput,
						EnvKey: "TEST_KEY",
					},
					SelectedAnswerToExpandedTemplate: map[string]string{
						"": "$TEST_KEY",
					},
				},
				Children: map[string]*AnswerTree{
					"": {
						Content: &SyntaxNode{
							Type: StepNode,
							Step: Step{
								ID: "fastlane",
								Inputs: []Input{
									{Key: "A", Value: "D"},
									{Key: "B", Value: "$TEST_KEY"},
								},
							},
						},
					},
				},
			},
		},
		{
			name: "Single Step, template, multiple selection question",
			node: &SyntaxNode{
				Type: StepNode,
				Step: Step{
					ID: "fastlane",
					Inputs: []Input{
						{
							Key: "A", Value: "{{.C}}",
						},
						{
							Key: "B", Value: `{{askForInputValue "test_question"}}`,
						},
					},
				},
			},
			values: map[string]string{
				"C": "D",
			},
			questions: map[string]Question{
				"test_question": {
					Title:   "title",
					Type:    models.TypeSelector,
					EnvKey:  "TEST_KEY",
					Answers: []string{"development", "app-store"},
				},
			},
			want: &AnswerTree{
				Answer: Answer2{
					key: AnswerKey{
						NodeKey: "B",
					},
					Question: &Question{
						ID:      "test_question",
						Title:   "title",
						Type:    models.TypeSelector,
						EnvKey:  "TEST_KEY",
						Answers: []string{"development", "app-store"},
					},
					SelectedAnswerToExpandedTemplate: map[string]string{
						"development": "$TEST_KEY",
						"app-store":   "$TEST_KEY",
					},
				},
				Children: map[string]*AnswerTree{
					"development": {
						Content: &SyntaxNode{
							Type: StepNode,
							Step: Step{
								ID: "fastlane",
								Inputs: []Input{
									{Key: "A", Value: "D"},
									{Key: "B", Value: "$TEST_KEY"},
								},
							},
						},
					},
					"app-store": {
						Content: &SyntaxNode{
							Type: StepNode,
							Step: Step{
								ID: "fastlane",
								Inputs: []Input{
									{Key: "A", Value: "D"},
									{Key: "B", Value: "$TEST_KEY"},
								},
							},
						},
					},
				},
			},
		},
		{
			name: "2 Steps, no template",
			node: &SyntaxNode{
				Type: StepListNode,
				Steps: Steps{
					Steps: []*SyntaxNode{
						{
							Type: StepNode,
							Step: Step{
								ID: "fastlane",
								Inputs: []Input{
									{Key: "A", Value: "B"},
								},
							},
						},
						{
							Type: StepNode,
							Step: Step{
								ID: "cache-push",
								Inputs: []Input{
									{Key: "C", Value: "D"},
								},
							},
						},
					},
				},
			},
			want: &AnswerTree{
				Content: &SyntaxNode{
					Type: StepListNode,
					Steps: Steps{
						Steps: []*SyntaxNode{
							{
								Type: StepNode,
								Step: Step{
									ID: "fastlane",
									Inputs: []Input{
										{Key: "A", Value: "B"},
									},
								},
							},
							{
								Type: StepNode,
								Step: Step{
									ID: "cache-push",
									Inputs: []Input{
										{Key: "C", Value: "D"},
									},
								},
							},
						},
					},
				},
			},
		},
		{
			name: "2 Steps, 1 Step with template",
			node: &SyntaxNode{
				Type: StepListNode,
				Steps: Steps{
					Steps: []*SyntaxNode{
						{
							Type: StepNode,
							Step: Step{
								ID: "fastlane",
								Inputs: []Input{
									{Key: "A", Value: "B"},
								},
							},
						},
						{
							Type: StepNode,
							Step: Step{
								ID: "xcode-archive",
								Inputs: []Input{
									{Key: "export_method", Value: `{{askForInputValue "export-method"}}`},
								},
							},
						},
					},
				},
			},
			questions: map[string]Question{
				"export-method": {
					Title:   "title",
					Answers: []string{"development", "app-store"},
					Type:    models.TypeSelector,
				},
			},
			want: &AnswerTree{
				Answer: Answer2{
					key: AnswerKey{
						NodeKey: "B",
					},
					Question: &Question{
						ID:      "export-method",
						Title:   "title",
						Type:    models.TypeSelector,
						Answers: []string{"development", "app-store"},
					},
					SelectedAnswerToExpandedTemplate: map[string]string{
						"development": "development",
						"app-store":   "app-store",
					},
				},
				Children: map[string]*AnswerTree{
					"development": {
						Content: &SyntaxNode{
							Type: StepListNode,
							Steps: Steps{
								Steps: []*SyntaxNode{
									{
										Type: StepNode,
										Step: Step{
											ID: "fastlane",
											Inputs: []Input{
												{Key: "A", Value: "B"},
											},
										},
									},
									{
										Type: StepNode,
										Step: Step{
											ID: "xcode-archive",
											Inputs: []Input{
												{Key: "export_method", Value: "development"},
											},
										},
									},
								},
							},
						},
					},
					"app-store": {
						Content: &SyntaxNode{
							Type: StepListNode,
							Steps: Steps{
								Steps: []*SyntaxNode{
									{
										Type: StepNode,
										Step: Step{
											ID: "fastlane",
											Inputs: []Input{
												{Key: "A", Value: "B"},
											},
										},
									},
									{
										Type: StepNode,
										Step: Step{
											ID: "xcode-archive",
											Inputs: []Input{
												{Key: "export_method", Value: "app-store"},
											},
										},
									},
								},
							},
						},
					},
				},
			},
		},
		{
			name: "2 Steps, 2 templated Steps",
			node: &SyntaxNode{
				Type: StepListNode,
				Steps: Steps{
					Steps: []*SyntaxNode{
						{
							Type: StepNode,
							Step: Step{
								ID: "fastlane",
								Inputs: []Input{
									{Key: "A", Value: `{{askForInputValue "question2"}}`},
								},
							},
						},
						{
							Type: StepNode,
							Step: Step{
								ID: "xcode-archive",
								Inputs: []Input{
									{Key: "export_method", Value: `{{askForInputValue "export-method"}}`},
								},
							},
						},
					},
				},
			},
			questions: map[string]Question{
				"export-method": {
					Title:   "title",
					Answers: []string{"development", "app-store"},
					Type:    models.TypeSelector,
				},
				"question2": {
					Title:   "title",
					Answers: []string{"n", "m"},
					Type:    models.TypeSelector,
				},
			},
			want: &AnswerTree{
				Answer: Answer2{
					key: AnswerKey{
						NodeKey: "A",
					},
					Question: &Question{
						ID:      "question2",
						Title:   "title",
						Type:    models.TypeSelector,
						Answers: []string{"n", "m"},
					},
					SelectedAnswerToExpandedTemplate: map[string]string{
						"n": "n",
						"m": "f",
					},
				},
				Children: map[string]*AnswerTree{
					"n": {
						Answer: Answer2{
							key: AnswerKey{
								NodeKey: "export_method",
							},
							Question: &Question{
								ID:      "export-method",
								Title:   "title",
								Type:    models.TypeSelector,
								Answers: []string{"development", "app-store"},
							},
							SelectedAnswerToExpandedTemplate: map[string]string{
								"development": "development",
								"app-store":   "app-store",
							},
						},
						Children: map[string]*AnswerTree{
							"development": {
								Content: &SyntaxNode{
									Type: StepListNode,
									Steps: Steps{
										Steps: []*SyntaxNode{
											{
												Type: StepNode,
												Step: Step{
													ID: "fastlane",
													Inputs: []Input{
														{Key: "A", Value: "n"},
													},
												},
											},
											{
												Type: StepNode,
												Step: Step{
													ID: "xcode-archive",
													Inputs: []Input{
														{Key: "export_method", Value: "development"},
													},
												},
											},
										},
									},
								},
							},
							"app-store": {
								Content: &SyntaxNode{
									Type: StepListNode,
									Steps: Steps{
										Steps: []*SyntaxNode{
											{
												Type: StepNode,
												Step: Step{
													ID: "fastlane",
													Inputs: []Input{
														{Key: "A", Value: "n"},
													},
												},
											},
											{
												Type: StepNode,
												Step: Step{
													ID: "xcode-archive",
													Inputs: []Input{
														{Key: "export_method", Value: "app-store"},
													},
												},
											},
										},
									},
								},
							},
						},
					},
					"m": {
						Answer: Answer2{
							key: AnswerKey{
								NodeKey: "export_method",
							},
							Question: &Question{
								ID:      "export-method",
								Title:   "title",
								Type:    models.TypeSelector,
								Answers: []string{"development", "app-store"},
							},
							SelectedAnswerToExpandedTemplate: map[string]string{
								"development": "development",
								"app-store":   "app-store",
							},
						},
						Children: map[string]*AnswerTree{
							"development": {
								Content: &SyntaxNode{
									Type: StepListNode,
									Steps: Steps{
										Steps: []*SyntaxNode{
											{
												Type: StepNode,
												Step: Step{
													ID: "fastlane",
													Inputs: []Input{
														{Key: "A", Value: "m"},
													},
												},
											},
											{
												Type: StepNode,
												Step: Step{
													ID: "xcode-archive",
													Inputs: []Input{
														{Key: "export_method", Value: "development"},
													},
												},
											},
										},
									},
								},
							},
							"app-store": {
								Content: &SyntaxNode{
									Type: StepListNode,
									Steps: Steps{
										Steps: []*SyntaxNode{
											{
												Type: StepNode,
												Step: Step{
													ID: "fastlane",
													Inputs: []Input{
														{Key: "A", Value: "m"},
													},
												},
											},
											{
												Type: StepNode,
												Step: Step{
													ID: "xcode-archive",
													Inputs: []Input{
														{Key: "export_method", Value: "app-store"},
													},
												},
											},
										},
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
			got, err := tt.node.Step.GetAnswers(tt.values, tt.questions)
			if (err != nil) != tt.wantErr {
				t.Errorf("SyntaxNode.Evaluate() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			// require.Equal(t, tt.want, got)
			assertEqual(t, tt.want, got)
		})
	}
}
*/
