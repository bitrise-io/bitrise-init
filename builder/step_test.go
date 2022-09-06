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
					{Key: "A", Value: &Text{Contents: "B"}},
				},
			},
			want: &Step{
				ID: "fastlane",
				Inputs: []Input{
					{Key: "A", Value: &Text{Contents: "B"}},
				},
			},
		},
		{
			name: "Single Step, template, optional question",
			node: &Step{
				ID: "fastlane",
				Inputs: []Input{
					{Key: "A", Value: &Text{Contents: "C"}},
					{Key: "B", Value: &InputFreeForm{QuestionID: "test_question"}},
				},
			},
			allAnswers: ConcreteAnswers{
				"test_question": ConcreteAnswer{
					SelectedAnswer: "",
					Answer: AnswerExpansion{
						Key: "test_question",
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
			// values: map[string]string{"C": "D"},
			want: &Step{
				ID: "fastlane",
				Inputs: []Input{
					{Key: "A", Value: &Text{Contents: "C"}},
					{Key: "B", Value: &Text{Contents: "$TEST_KEY"}},
				},
			},
		},
		{
			name: "Single Step, template, multiple selection question",
			node: &Step{
				ID: "fastlane",
				Inputs: []Input{
					{Key: "A", Value: &Text{Contents: "C"}},
					{Key: "B", Value: &InputSelect{QuestionID: "test_question"}},
				},
			},
			allAnswers: ConcreteAnswers{
				"test_question": ConcreteAnswer{
					SelectedAnswer: "C",
					Answer: AnswerExpansion{
						Key: "test_question",
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
			// values: map[string]string{"C": "D"},
			want: &Step{
				ID: "fastlane",
				Inputs: []Input{
					{Key: "A", Value: &Text{Contents: "D"}},
					{Key: "B", Value: &Text{Contents: "C"}},
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
