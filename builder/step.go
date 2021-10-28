package builder

import (
	"bytes"
	"fmt"
	"text/template"

	"github.com/bitrise-io/bitrise-init/models"
)

const defaultAnswer = ""

type Input struct {
	Key, Value string
}

type Step struct {
	ID     string
	Title  string
	RunIf  string
	Inputs []Input

	templateID int
}

type Steps struct {
	IncludeIfTrue Question
	Steps         []TemplateNode

	templateID int
}

func (s *Step) Execute(values map[string]string, allAnswers ConcreteAnswers) (TemplateNode, error) {
	var (
		expandedInputs = []Input{}
	)

	if s.ID == "" {
		panic("No Step ID provided")
	}

	for _, input := range s.Inputs {
		template, err := template.New("Input").Funcs(template.FuncMap{
			"askForInputValue": func(_ string) string {
				answer, ok := allAnswers[AnswerKey{
					nodeID:  s.GetID(),
					NodeKey: input.Key,
				}]
				if !ok {
					panic(fmt.Sprintf("answer not found for node (%+v) and key (%s) in list (%+v)", s, input.Key, allAnswers))
				}

				return answer.Answer.SelectionToExpandedValue[answer.SelectedAnswer]
			},
		}).Parse(input.Value)
		if err != nil {
			return nil, err
		}

		expandedValue := new(bytes.Buffer)
		if err := template.Execute(expandedValue, values); err != nil {
			return nil, err
		}

		expandedInputs = append(expandedInputs, Input{
			Key:   input.Key,
			Value: expandedValue.String(),
		})
	}

	return &Step{
		ID:     s.ID,
		Title:  s.Title,
		Inputs: expandedInputs,
		RunIf:  s.RunIf,
	}, nil
}

func (s *Step) GetAnswers(questions map[string]Question) (*AnswerTree, error) {
	var allAnswers []AnswerExpansion

	for _, input := range s.Inputs {
		var inputQuestionExpansions *AnswerExpansion

		template, err := template.New("Input").Funcs(template.FuncMap{
			"askForInputValue": func(questionID string) string {
				question, ok := questions[questionID]
				if !ok {
					panic(fmt.Sprintf("Question (%s) unknown", questionID))
				}
				question.ID = questionID // Used for the config map key generation

				selectedValueToTemplateExpansion := make(map[string]string)
				if question.Type == models.TypeSelector || question.Type == models.TypeOptionalSelector {
					for _, answer := range question.Selections {
						if question.EnvKey == "" {
							selectedValueToTemplateExpansion[answer] = answer

							continue
						}

						selectedValueToTemplateExpansion[answer] = "$" + question.EnvKey
					}
				}
				if question.Type == models.TypeUserInput ||
					question.Type == models.TypeOptionalUserInput ||
					question.Type == models.TypeOptionalSelector {
					selectedValueToTemplateExpansion[defaultAnswer] = "$" + question.EnvKey
				}

				inputQuestionExpansions = &AnswerExpansion{
					Key: AnswerKey{
						nodeID:  s.GetID(),
						NodeKey: input.Key,
					},
					Question:                 &question,
					SelectionToExpandedValue: selectedValueToTemplateExpansion,
				}

				// Return value is ignored
				return ""
			},
		}).Parse(input.Value)
		if err != nil {
			return nil, fmt.Errorf("error parsing template `%s`: %s", input.Value, err)
		}

		// Using go templates to call custom expansion logic, ignoring the text output.
		if err := template.Execute(new(bytes.Buffer), nil); err != nil {
			return nil, err
		}

		if inputQuestionExpansions != nil {
			allAnswers = append(allAnswers, *inputQuestionExpansions)
		}
	}

	return newAnswerTree(allAnswers), nil
}

func (s *Step) SetID(templateIDCounter int) int {
	templateIDCounter++
	s.templateID = templateIDCounter

	return templateIDCounter
}

func (s *Step) GetID() int {
	return s.templateID
}

func (s *Steps) GetAnswers(questions map[string]Question) (*AnswerTree, error) {
	var answerTrees []*AnswerTree

	for _, step := range s.Steps {
		answerTree, err := step.GetAnswers(questions)
		if err != nil {
			return nil, err
		}

		if answerTree != nil {
			answerTrees = append(answerTrees, answerTree)
		}
	}

	return mergeAnswerTrees(answerTrees), nil
}

func (s *Steps) SetID(templateIDCounter int) int {
	templateIDCounter++
	s.templateID = templateIDCounter

	for _, step := range s.Steps {
		step.SetID(templateIDCounter)
	}

	return templateIDCounter
}

func (s *Steps) GetID() int {
	return s.templateID
}

func (s *Steps) Execute(values map[string]string, allAnswers ConcreteAnswers) (TemplateNode, error) {
	allSteps := []TemplateNode{}

	for _, step := range s.Steps {
		evaluatedStep, err := step.Execute(values, allAnswers)
		if err != nil {
			return nil, err
		}

		allSteps = append(allSteps, evaluatedStep)
	}

	return &Steps{
		Steps: allSteps,
	}, nil
}

func (s *Steps) TemplateID() uintptr {
	return 0
}

func (s *Steps) Append(step Step) Steps {
	s.Steps = append(s.Steps, &step)

	return *s
}
