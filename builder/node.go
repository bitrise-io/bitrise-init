package builder

import (
	"github.com/bitrise-io/bitrise-init/models"
)

type Question struct {
	ID         string
	Title      string
	Summary    string
	EnvKey     string
	Type       models.Type
	Selections []string
}

type AnswerKey struct {
	nodeID  int
	NodeKey string
}

type AnswerExpansion struct {
	Key                      AnswerKey
	Question                 *Question
	SelectionToExpandedValue map[string]string
}

type ConcreteAnswer struct {
	Answer           AnswerExpansion
	SelectedAnswer   string
	HasEqualChildren bool
}

type ConcreteAnswers map[AnswerKey]ConcreteAnswer

type AnswerTree struct {
	Answer   AnswerExpansion
	Children map[string]*AnswerTree
}

type TemplateNode interface {
	GetAnswers(questions map[string]Question, context []interface{}) (*AnswerTree, error)
	Execute(values map[string]string, answers ConcreteAnswers) (TemplateNode, error)
	Export() (ExportFragment, error)
	SetID(templateIDCounter int) int
	GetID() int
}

func (a *AnswerTree) HasEqualChildren() bool {
	firstSelectedAnswer := a.Answer.Question.Selections[0]
	if len(a.Children) == 0 || a.Children[firstSelectedAnswer] == nil {
		firstValue := a.Answer.SelectionToExpandedValue[firstSelectedAnswer]
		for _, expandedValue := range a.Answer.SelectionToExpandedValue {
			if expandedValue != firstValue {
				return false
			}
		}

		return true
	}

	firstChild := a.Children[firstSelectedAnswer]
	for _, child := range a.Children {
		if child != firstChild {
			return false
		}
	}

	return true
}

func newAnswerTree(answerList []AnswerExpansion) *AnswerTree {
	if len(answerList) == 0 {
		return nil
	}

	root := &AnswerTree{
		Answer: answerList[0],
	}
	current := root

	for _, answer := range answerList[1:] {
		next := &AnswerTree{
			Answer: answer,
		}
		current.AddChild(next)
		current = next
	}

	current.AddChild(nil)

	return root
}

func (current *AnswerTree) AddChild(next *AnswerTree) {
	if current.Children == nil {
		current.Children = make(map[string]*AnswerTree)
	}

	for selectedAnswer := range current.Answer.SelectionToExpandedValue {
		current.Children[selectedAnswer] = next
	}
}

func (current *AnswerTree) append(next *AnswerTree) {
	if next == nil {
		return
	}

	if len(current.Children) == 0 {
		current.AddChild(next)

		return
	}

	for selectedAnswer, child := range current.Children {
		if child == nil {
			current.Children[selectedAnswer] = next

			continue
		}

		child.append(next)
	}
}

func mergeAnswerTrees(answerTrees []*AnswerTree) *AnswerTree {
	if len(answerTrees) == 0 {
		return nil
	}

	result := answerTrees[0]
	for _, answerTree := range answerTrees[1:] {
		result.append(answerTree)
	}

	return result
}
