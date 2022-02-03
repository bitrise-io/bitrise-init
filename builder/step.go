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
		// expandFunc := func(_ string) string {
		// 	answer, ok := allAnswers[AnswerKey{
		// 		nodeID:  s.GetID(),
		// 		NodeKey: input.Key,
		// 	}]
		// 	if !ok {
		// 		panic(fmt.Sprintf("answer missing for node (%+v) and key (%s) in list (%+v)", s, input.Key, allAnswers))
		// 	}

		// 	expandedValue, ok := answer.Answer.SelectionToExpandedValue[answer.SelectedAnswer]
		// 	if !ok {
		// 		panic(fmt.Sprintf("expanded value missing for selected answer (%s), in list (%+v)", answer.SelectedAnswer, answer.Answer.SelectionToExpandedValue))
		// 	}

		// 	return expandedValue
		// }

		// template, err := template.New("Input").Funcs(template.FuncMap{
		// 	"askForInputValue": expandFunc,
		// 	"selectFromContext": func(p, _ string) string {
		// 		return expandFunc(p)
		// 	},
		// }).Parse(input.Value)
		// if err != nil {
		// 	return nil, err
		// }

		// expandedValue := new(bytes.Buffer)
		// if err := template.Execute(expandedValue, values); err != nil {
		// 	return nil, fmt.Errorf("Execute() failed: %s", err)
		// }
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
		// var inputQuestionExpansions *AnswerExpansion

		/*template, err := template.New("Input").Funcs(template.FuncMap{
			"askForInputValue": func(questionID string) string {
				question, ok := questions[questionID]
				if !ok {
					panic(fmt.Sprintf("Question (%s) undefined", questionID))
				}
				question.ID = questionID // Used for the config map key generation

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
			"selectFromContext": func(questionID string, sourceTag string) string {
				question, ok := questions[questionID]
				if !ok {
					panic(fmt.Sprintf("Question (%s) undefined", questionID))
				}
				question.ID = questionID // Used for the config map key generation

				if question.Type != models.TypeSelector &&
					question.Type != models.TypeOptionalSelector {
					panic("selectFromContext supported for selector type questions")
				}

				selectableAnswers := []string{}
			tagSearch:
				for _, c := range context {
					val := reflect.ValueOf(c)
					if val.Kind() != reflect.Slice {
						panic("Unsupported context variable, expected slice")
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
							if ok && tag == sourceTag {
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

				inputQuestionExpansions = &AnswerExpansion{
					Key: AnswerKey{
						nodeID:  s.GetID(),
						NodeKey: input.Key,
					},
					Question:                 &question,
					SelectionToExpandedValue: selectedValueToTemplateExpansion,
				}

				return ""
			},
		}).Parse(input.Value)
		if err != nil {
			return nil, fmt.Errorf("error parsing template `%s`: %s", input.Value, err)
		}

		// Using go templates to call custom expansion logic, ignoring the text output.
		if err := template.Execute(new(bytes.Buffer), nil); err != nil {
			return nil, err
		}*/

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
