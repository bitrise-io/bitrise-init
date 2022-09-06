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
		cmpopts.IgnoreUnexported(Step{}),
		cmpopts.IgnoreUnexported(Text{}),
		// cmp.Comparer(func(x, y AnswerKey) bool { return x.NodeKey == y.NodeKey }),
	}

	if !cmp.Equal(want, got, opts...) {
		t.Fatalf("Not equal:\n%s", cmp.Diff(want, got, opts...))
	}
}

func TestTemplateNode_GetAnswers(t *testing.T) {
	tests := []struct {
		name      string
		node      TemplateNode
		values    map[string]string
		questions map[string]Question
		context   []interface{}
		want      *AnswerTree
		wantErr   bool
	}{
		// Empty StepList
		{
			name: "Single Step, no template",
			node: &Step{
				ID: "fastlane",
				Inputs: []Input{
					{Key: "A", Value: &Text{Contents: "B"}},
				},
			},
			want: nil,
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
			questions: map[string]Question{
				"test_question": {
					Title:  "title",
					Type:   models.TypeOptionalUserInput,
					EnvKey: "TEST_KEY",
				},
			},
			want: &AnswerTree{
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
				Children: map[string]*AnswerTree{
					"": nil,
				},
			},
		},
		{
			name: "Single Step, template, question with answers from context",
			node: &Step{
				ID: "fastlane",
				Inputs: []Input{
					{Key: "A", Value: &InputSelect{QuestionID: "test_question", ContextTag: "projectPath"}},
				},
			},
			questions: map[string]Question{
				"test_question": {
					Title: "title",
					Type:  models.TypeSelector,
				},
			},
			context: []interface{}{
				[]struct {
					projectPath string `builder:"projectPath"`
				}{
					{projectPath: "path-1"},
					{projectPath: "path-2"},
				},
			},
			want: &AnswerTree{
				Answer: AnswerExpansion{
					Key: "test_question",
					Question: &Question{
						ID:    "test_question",
						Title: "title",
						Type:  models.TypeSelector,
					},
					SelectionToExpandedValue: map[string]string{
						"path-1": "path-1",
						"path-2": "path-2",
					},
				},
				Children: map[string]*AnswerTree{
					"path-1": nil,
					"path-2": nil,
				},
			},
		},
		{
			name: "2 Steps, 2 templated Steps",
			node: &Steps{
				Steps: []TemplateNode{
					&Step{
						ID: "fastlane",
						Inputs: []Input{
							{Key: "A", Value: &InputSelect{QuestionID: "question2"}},
						},
					},
					&Step{
						ID: "xcode-archive",
						Inputs: []Input{
							{Key: "export_method", Value: &InputSelect{QuestionID: "export-method"}},
						},
					},
				},
			},
			questions: map[string]Question{
				"export-method": {
					Title:      "title",
					Selections: []string{"development", "app-store"},
					Type:       models.TypeSelector,
				},
				"question2": {
					Title:      "title2",
					Selections: []string{"n", "m"},
					Type:       models.TypeSelector,
				},
			},
			want: &AnswerTree{
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
							Key: "export_method",
							Question: &Question{
								ID:         "export-method",
								Title:      "title",
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
							Key: "export_method",
							Question: &Question{
								ID:         "export-method",
								Title:      "title",
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
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.node.GetAnswers(tt.questions, tt.context)
			if (err != nil) != tt.wantErr {
				t.Errorf("TemplateNode.GetAnswers() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			// require.Equal(t, tt.want, got)
			assertEqual(t, tt.want, got)
		})
	}
}
