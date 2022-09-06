package builder

type Workflow struct {
	ID          string
	Steps       TemplateNode
	Description string
}

type Workflows struct {
	ProjectType string
	Workflows   []Workflow
}

func (w *Workflows) GetAnswers(questions map[string]Question, context []interface{}) (*AnswerTree, error) {
	var answerTrees []*AnswerTree

	for _, wf := range w.Workflows {
		answerTree, err := wf.Steps.GetAnswers(questions, context)
		if err != nil {
			return nil, err
		}

		if answerTree != nil {
			answerTrees = append(answerTrees, answerTree)
		}
	}

	return mergeAnswerTrees(answerTrees), nil
}

func (w *Workflows) Execute(values map[string]string, allAnswers ConcreteAnswers) (TemplateNode, error) {
	var allWorkflows []Workflow

	for _, wf := range w.Workflows {
		evaluatedSteps, err := wf.Steps.Execute(values, allAnswers)
		if err != nil {
			return nil, err
		}

		allWorkflows = append(allWorkflows, Workflow{
			ID:          wf.ID,
			Steps:       evaluatedSteps,
			Description: wf.Description,
		})
	}

	return &Workflows{
		ProjectType: w.ProjectType,
		Workflows:   allWorkflows,
	}, nil
}
