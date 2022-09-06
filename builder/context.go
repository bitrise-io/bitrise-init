package builder

import "fmt"

type Context struct {
	SelectFrom interface{}
	Questions  map[string]Question

	Template TemplateNode
}

func (c *Context) GetAnswers(questions map[string]Question, context []interface{}) (*AnswerTree, error) {
	for ID, question := range c.Questions {
		if previous, exists := questions[ID]; exists {
			return nil, fmt.Errorf("duplicate question ID (%s), previously seen: (%s), current: (%s)", ID, previous, question)
		}

		questions[ID] = question
	}

	if c.SelectFrom != nil {
		context = append([]interface{}{c.SelectFrom}, context...)
	}

	return c.Template.GetAnswers(questions, context)
}

func (c *Context) Execute(values map[string]string, answers ConcreteAnswers) (TemplateNode, error) {
	return c.Template.Execute(values, answers)
}

func (c *Context) Export() (ExportFragment, error) {
	return c.Template.Export()
}
