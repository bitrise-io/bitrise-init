package builder

import (
	"testing"

	"github.com/bitrise-io/bitrise-init/models"
)

func TestTemplateNode_Execute(t *testing.T) {
	tests := []struct {
		name       string
		node       TemplateNode
		values     map[string]string
		allAnswers ConcreteAnswers
		want       TemplateNode
		wantErr    bool
	}{
		{
			name: "Single Step, no template",
			node: &Step{
				ID: "fastlane",
				Inputs: []Input{
					{Key: "A", Value: "B"},
				},
			},
			want: &Step{
				ID: "fastlane",
				Inputs: []Input{
					{Key: "A", Value: "B"},
				},
			},
		},
		{
			name: "Single Step, template, optional question",
			node: &Step{
				templateID: 1,
				ID:         "fastlane",
				Inputs: []Input{
					{Key: "A", Value: "{{.C}}"},
					{Key: "B", Value: `{{askForInputValue "test_question"}}`},
				},
			},
			allAnswers: ConcreteAnswers{
				AnswerKey{nodeID: 1, NodeKey: "B"}: ConcreteAnswer{
					SelectedAnswer: "",
					Answer: AnswerExpansion{
						Key: AnswerKey{nodeID: 1, NodeKey: "B"},
						Question: &Question{
							ID:     "test_question",
							Title:  "title",
							Type:   models.TypeOptionalUserInput,
							EnvKey: "TEST_KEY",
						},
						SelectionToExpandedValue: map[string]string{
							"": "$TEST_KEY",
						},
					},
				},
			},
			values: map[string]string{"C": "D"},
			want: &Step{
				templateID: 1,
				ID:         "fastlane",
				Inputs: []Input{
					{Key: "A", Value: "D"},
					{Key: "B", Value: "$TEST_KEY"},
				},
			},
		},
		{
			name: "Single Step, template, multiple selection question",
			node: &Step{
				templateID: 1,
				ID:         "fastlane",
				Inputs: []Input{
					{Key: "A", Value: "{{.C}}"},
					{Key: "B", Value: `{{askForInputValue "test_question"}}`},
				},
			},
			allAnswers: ConcreteAnswers{
				AnswerKey{nodeID: 1, NodeKey: "B"}: ConcreteAnswer{
					SelectedAnswer: "C",
					Answer: AnswerExpansion{
						Key: AnswerKey{nodeID: 1, NodeKey: "B"},
						Question: &Question{
							ID:    "test_question",
							Title: "title",
							Type:  models.TypeOptionalSelector,
						},
						SelectionToExpandedValue: map[string]string{
							"C": "C",
						},
					},
				},
			},
			values: map[string]string{"C": "D"},
			want: &Step{
				templateID: 1,
				ID:         "fastlane",
				Inputs: []Input{
					{Key: "A", Value: "D"},
					{Key: "B", Value: "C"},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.node.Execute(tt.values, tt.allAnswers)
			if (err != nil) != tt.wantErr {
				t.Errorf("Step.Execute() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			assertEqual(t, tt.want, got)
		})
	}
}
