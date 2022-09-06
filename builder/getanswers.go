package builder

import (
	"fmt"
	"reflect"

	"github.com/bitrise-io/bitrise-init/models"
	"github.com/bitrise-io/go-utils/log"
)

type Text struct {
	Contents string
}

func NewText(contents string) TemplateNode {
	return &Text{Contents: contents}
}

func (t *Text) GetAnswers(questions map[string]Question, context []interface{}) (*AnswerTree, error) {
	return nil, nil
}

func (t *Text) Execute(values map[string]string, allAnswers ConcreteAnswers) (TemplateNode, error) {
	return t, nil
}

func (t *Text) Export() (ExportFragment, error) {
	return ExportFragment{Text: t.Contents}, nil
}

type InputSelect struct {
	QuestionID string
	ContextTag string
}

func (input *InputSelect) Execute(values map[string]string, allAnswers ConcreteAnswers) (TemplateNode, error) {
	return executeInput(input.QuestionID, values, allAnswers)
}

func (input *InputSelect) GetAnswers(questions map[string]Question, context []interface{}) (*AnswerTree, error) {
	question, ok := questions[input.QuestionID]
	if !ok {
		return nil, fmt.Errorf("question (%s) undefined", input.QuestionID)
	}
	question.ID = input.QuestionID // Used for the config map key generation

	if question.Type != models.TypeSelector &&
		question.Type != models.TypeOptionalSelector {
		return nil, fmt.Errorf("InputSelect requires a selector type question")
	}

	selectableAnswers := []string{}
tagSearch:
	for _, c := range context {
		val := reflect.ValueOf(c)
		if val.Kind() != reflect.Slice {
			return nil, fmt.Errorf("unsupported context variable, expected slice")
		}

		tagFound := false
		for i := 0; i < val.Len(); i++ {
			elem := val.Index(i)
			if elem.Kind() != reflect.Struct {
				log.Debugf("Expected  struct")

				continue
			}

			for j := 0; j < elem.NumField(); j++ {
				tag, ok := elem.Type().Field(j).Tag.Lookup("builder")
				if ok && tag == input.ContextTag {
					selectableAnswers = append(selectableAnswers, elem.Field(j).String())

					tagFound = true
				}
			}
		}

		if tagFound {
			break tagSearch
		}
	}

	selectedValueToTemplateExpansion := make(map[string]string)
	for _, answer := range selectableAnswers {
		if question.EnvKey == "" {
			// Actual answer directly written to the config.
			selectedValueToTemplateExpansion[answer] = answer

			continue
		}

		// Indirectly added to the config. Environment variable is added as a global env by the frontend.
		selectedValueToTemplateExpansion[answer] = "$" + question.EnvKey
	}

	if question.Type == models.TypeOptionalSelector {
		if len(question.EnvKey) == 0 {
			panic("Question Environment Key name empty, required for freeform answers")
		}
	}

	question.Selections = selectableAnswers

	return newAnswerTree([]AnswerExpansion{{
		Key:                      input.QuestionID,
		Question:                 &question,
		SelectionToExpandedValue: selectedValueToTemplateExpansion,
	}}), nil
}

func (t *InputSelect) Export() (ExportFragment, error) {
	return ExportFragment{}, fmt.Errorf("Unsupported operation")
}

type InputFreeForm struct {
	QuestionID string
}

func (input *InputFreeForm) Execute(values map[string]string, allAnswers ConcreteAnswers) (TemplateNode, error) {
	return executeInput(input.QuestionID, values, allAnswers)
}

func (input *InputFreeForm) GetAnswers(questions map[string]Question, context []interface{}) (*AnswerTree, error) {
	question, ok := questions[input.QuestionID]
	if !ok {
		panic(fmt.Sprintf("Question (%s) undefined", input.QuestionID))
	}
	question.ID = input.QuestionID // Used for the config map key generation

	selectedValueToTemplateExpansion := make(map[string]string)

	if question.Type != models.TypeUserInput &&
		question.Type != models.TypeOptionalUserInput {
		panic("askForInputValue supported for freeform type questions")
	}

	if len(question.EnvKey) == 0 {
		panic("Question Environment Key name empty, required for freeform answers")
	}
	selectedValueToTemplateExpansion[defaultAnswer] = "$" + question.EnvKey
	question.Selections = []string{defaultAnswer}

	return newAnswerTree([]AnswerExpansion{{
		Key:                      input.QuestionID,
		Question:                 &question,
		SelectionToExpandedValue: selectedValueToTemplateExpansion,
	}}), nil
}

func (t *InputFreeForm) Export() (ExportFragment, error) {
	return ExportFragment{}, fmt.Errorf("Unsupported operation")
}

func executeInput(questionID string, _ map[string]string, allAnswers ConcreteAnswers) (TemplateNode, error) {
	if questionID == "" {
		return nil, fmt.Errorf("no question ID provided")
	}

	answer, ok := allAnswers[questionID]
	if !ok {
		return nil, fmt.Errorf(fmt.Sprintf("answer missing for key (%s) in list (%+v)", questionID, allAnswers))
	}

	expandedValue, ok := answer.Answer.SelectionToExpandedValue[answer.SelectedAnswer]
	if !ok {
		return nil, fmt.Errorf("expanded value missing for selected answer (%s), in list (%+v)", answer.SelectedAnswer, answer.Answer.SelectionToExpandedValue)
	}

	return NewText(expandedValue), nil
}
