package builder

const defaultAnswer = ""

type Input struct {
	Key   string
	Value TemplateNode
}

type Step struct {
	ID     string
	Title  string
	RunIf  string
	Inputs []Input

	templateID int
}

type Steps struct {
	Steps []TemplateNode

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
		expandedValue, err := input.Value.Execute(values, allAnswers)
		if err != nil {
			return nil, err
		}

		expandedInputs = append(expandedInputs, Input{
			Key:   input.Key,
			Value: expandedValue,
		})
	}

	return &Step{
		ID:     s.ID,
		Title:  s.Title,
		Inputs: expandedInputs,
		RunIf:  s.RunIf,
	}, nil
}

func (s *Step) GetAnswers(questions map[string]Question, context []interface{}) (*AnswerTree, error) {
	var allAnswers []*AnswerTree

	for _, input := range s.Inputs {
		answers, err := input.Value.GetAnswers(questions, context)
		if err != nil {
			return nil, err
		}

		if answers != nil {
			allAnswers = append(allAnswers, answers)
		}
	}

	return mergeAnswerTrees(allAnswers), nil
}

func (s *Step) SetID(templateIDCounter int) int {
	templateIDCounter++
	s.templateID = templateIDCounter

	return templateIDCounter
}

func (s *Step) GetID() int {
	return s.templateID
}

func (s *Steps) GetAnswers(questions map[string]Question, context []interface{}) (*AnswerTree, error) {
	var answerTrees []*AnswerTree

	for _, step := range s.Steps {
		answerTree, err := step.GetAnswers(questions, context)
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
		templateIDCounter = step.SetID(templateIDCounter)
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

func (s *Steps) Append(step Step) *Steps {
	s.Steps = append(s.Steps, &step)

	return s
}
