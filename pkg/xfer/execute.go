package xfer

import (
	"fmt"

	"get.porter.sh/porter/pkg/exec/builder"
	yaml "gopkg.in/yaml.v2"
)

func (m *Mixin) loadAction() (*Action, error) {
	var action Action
	err := builder.LoadAction(m.Context, "", func(contents []byte) (interface{}, error) {
		err := yaml.Unmarshal(contents, &action)
		return &action, err
	})
	return &action, err
}

func (m *Mixin) Execute() error {
	action, err := m.loadAction()
	if err != nil {
		return err
	}
	
	if m.Context.Debug {
		steps := (*action).Steps
		for i, s := range steps {
			fmt.Fprintf(m.Out, "adding debug argument to %s \n", s.Description)
			s.Arguments = append(s.Arguments, "--debug")
			// range creates a local var for s and we lose the pointer
			steps[i].Arguments = s.Arguments
		}
	}
	_, err = builder.ExecuteSingleStepAction(m.Context, action)
	return err
}
