package builder

import (
	"errors"
	"fmt"
	"strings"

	"github.com/bitrise-io/bitrise-init/models"
	"github.com/bitrise-io/bitrise-init/steps"
	bitriseModels "github.com/bitrise-io/bitrise/models"
	envmanModels "github.com/bitrise-io/envman/models"
)

type Result struct {
	Config    *bitriseModels.BitriseDataModel
	Artifacts Artifacts
}

type ExportFragmentType int

const (
	StepFragment ExportFragmentType = iota
	StepListFragment
)

type ExportFragment struct {
	Type     ExportFragmentType
	Step     bitriseModels.StepListItemModel
	StepList []bitriseModels.StepListItemModel
}

type Artifacts struct {
	Icons    []string
	Warnings []string
}

func Append(a, b Artifacts) *Artifacts {
	merged := &Artifacts{}

	merged.Icons = append(a.Icons, b.Icons...)
	merged.Warnings = append(a.Warnings, b.Warnings...)

	return merged
}

func Export(template TemplateNode, answerTree *AnswerTree, values map[string]string) (*models.OptionNode, map[string]Result, error) {
	if answerTree == nil {
		return exportAnswerTreeLeaf(template, values, nil, "")
	}

	var (
		allOptions    *models.OptionNode = newOptionNode(*answerTree.Answer.Question)
		mergedResults                    = make(map[string]Result)
	)

	if err := walkPaths(answerTree, []ConcreteAnswer{}, "", func(answers []ConcreteAnswer, treePathHash string) error {
		configOption, results, err := exportAnswerTreeLeaf(template, values, answers, treePathHash)
		if err != nil {
			return err
		}

		currentNode := allOptions
		for i, concreteAnswer := range answers {
			key := concreteAnswer.SelectedAnswer

			if i == len(answers)-1 {
				currentNode.AddConfig(key, configOption)

				break
			}

			_, hasChild := currentNode.ChildOptionMap[key]
			if !hasChild {
				newOption := newOptionNode(*concreteAnswer.Answer.Question)
				currentNode.AddOption(key, newOption)
				currentNode = newOption

				continue
			}

			currentNode = currentNode.ChildOptionMap[key]
		}

		for pathHash, result := range results {
			if _, exists := mergedResults[pathHash]; exists {
				panic(fmt.Sprintf("duplicate config result key (%s), unique expected", pathHash))
			}
			mergedResults[pathHash] = result
		}

		return nil
	}); err != nil {
		return nil, nil, err
	}

	return allOptions, mergedResults, nil
}

func exportAnswerTreeLeaf(template TemplateNode, values map[string]string, answers []ConcreteAnswer, pathHash string) (*models.OptionNode, map[string]Result, error) {
	allAnswers := make(ConcreteAnswers)
	for _, answer := range answers {
		allAnswers[answer.Answer.Key] = answer
	}

	output, err := template.Execute(values, allAnswers)
	if err != nil {
		return nil, nil, err
	}

	result, err := exportResult(output)
	if err != nil {
		return nil, nil, err
	}
	resultMap := map[string]Result{
		pathHash: result,
	}

	return models.NewConfigOption(pathHash, result.Artifacts.Icons), resultMap, err
}

func newOptionNode(question Question) *models.OptionNode {
	return models.NewOption(question.Title, question.Summary, question.EnvKey, question.Type)
}

func walkPaths(tree *AnswerTree, concreteAnswers []ConcreteAnswer, pathHash string, callbackFn func(answers []ConcreteAnswer, pathHash string) error) error {
	if tree == nil {
		return nil
	}

	for _, selectedAnswer := range tree.Answer.Question.Selections { // Preserve answer order
		var (
			nextConcreteAnswers = append(concreteAnswers, ConcreteAnswer{
				Answer:         tree.Answer,
				SelectedAnswer: selectedAnswer,
			})
			nextPathHash = pathHash
		)

		if !tree.HasEqualChildren() {
			nextPathHash = fmt.Sprintf("%s/%s.%s",
				nextPathHash,
				tree.Answer.Question.ID,
				selectedAnswer,
			)
		}

		if len(tree.Children) == 0 {
			if err := callbackFn(nextConcreteAnswers, nextPathHash); err != nil {
				return err
			}

			continue
		}

		if _, ok := tree.Answer.SelectionToExpandedValue[selectedAnswer]; !ok {
			panic(fmt.Sprintf("selected answer (%s) missing", selectedAnswer))
		}

		nextAnswer, ok := tree.Children[selectedAnswer]
		if !ok {
			panic(fmt.Sprintf("selected answer (%s) missing", selectedAnswer))
		}
		if nextAnswer == nil {
			if err := callbackFn(nextConcreteAnswers, nextPathHash); err != nil {
				return err
			}

			continue
		}

		if err := walkPaths(tree.Children[selectedAnswer], nextConcreteAnswers, nextPathHash, callbackFn); err != nil {
			return err
		}
	}

	return nil
}

func exportResult(node TemplateNode) (Result, error) {
	exportFragment, err := node.Export()
	if err != nil {
		return Result{}, err
	}

	switch exportFragment.Type {
	case StepListFragment:
		{
			configBuilder := models.NewDefaultConfigBuilder()
			configBuilder.AppendStepListItemsTo(models.PrimaryWorkflowID, exportFragment.StepList...)

			config, err := configBuilder.Generate("")
			if err != nil {
				return Result{}, err
			}

			return Result{
				Config: &config,
			}, nil
		}
	default:
		panic("missing implementation")
	}
}

func (s *Steps) Export() (ExportFragment, error) {
	steps := []bitriseModels.StepListItemModel{}

	for _, step := range s.Steps {
		stepFragment, err := step.Export()
		if err != nil {
			return ExportFragment{}, err
		}

		switch stepFragment.Type {
		case StepFragment:
			steps = append(steps, stepFragment.Step)
		case StepListFragment:
			steps = append(steps, stepFragment.StepList...)
		default:
			return ExportFragment{}, errors.New("failed to export step list, contaning unsupported fragment")
		}
	}

	return ExportFragment{
		Type:     StepListFragment,
		StepList: steps,
	}, nil
}

func (s *Step) Export() (ExportFragment, error) {
	return ExportFragment{
		Type: StepFragment,
		Step: newStep(s.ID, s.Title, s.Inputs, s.RunIf),
	}, nil
}

func newStep(id string, title string, inputs []Input, runIf string) bitriseModels.StepListItemModel {
	if !strings.Contains(id, "@") {
		version, ok := steps.StepIDToVersion[id]
		if !ok {
			panic(fmt.Sprintf("Unknown Step (%s) version", id))
		}

		id = id + "@" + version
	}

	inputEnvs := []envmanModels.EnvironmentItemModel{}
	for _, input := range inputs {
		inputEnvs = append(inputEnvs, envmanModels.EnvironmentItemModel{
			input.Key: input.Value,
		})
	}

	return steps.StepListItem(id, title, runIf, inputEnvs...)
}
