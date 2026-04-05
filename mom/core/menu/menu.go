package menu

import "github.com/manifoldco/promptui"

type Handler func() error

type Option struct {
	Label   string
	Handler Handler
}

type Loop struct {
	Label     string
	Options   []Option
	ExitLabel string
}

func NewLoop(label string, options []Option) *Loop {
	return &Loop{
		Label:     label,
		Options:   options,
		ExitLabel: "Sair",
	}
}

func (l *Loop) Run() error {
	for {
		items := make([]string, 0, len(l.Options)+1)
		for _, option := range l.Options {
			items = append(items, option.Label)
		}
		items = append(items, l.ExitLabel)

		prompt := promptui.Select{
			Label: l.Label,
			Items: items,
			Size:  len(items),
		}

		index, _, err := prompt.Run()
		if err != nil {
			if err == promptui.ErrInterrupt || err == promptui.ErrEOF {
				return nil
			}

			return err
		}

		if index == len(l.Options) {
			return nil
		}

		if l.Options[index].Handler == nil {
			continue
		}

		if err := l.Options[index].Handler(); err != nil {
			return err
		}
	}
}

func Select(label string, items []string) (int, string, error) {
	prompt := promptui.Select{
		Label: label,
		Items: items,
		Size:  len(items),
	}

	index, value, err := prompt.Run()
	if err != nil {
		return 0, "", err
	}

	return index, value, nil
}
