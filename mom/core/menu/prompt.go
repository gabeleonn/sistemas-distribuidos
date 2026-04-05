package menu

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/manifoldco/promptui"
)

func PromptText(label string) (string, error) {
	prompt := promptui.Prompt{
		Label:    label,
		Validate: requiredValidator,
	}

	result, err := prompt.Run()
	if err != nil {
		return "", err
	}

	return strings.TrimSpace(result), nil
}

func PromptInt(label string) (int, error) {
	prompt := promptui.Prompt{
		Label: label,
		Validate: func(input string) error {
			if requiredValidator(input) != nil {
				return requiredValidator(input)
			}

			if _, err := strconv.Atoi(strings.TrimSpace(input)); err != nil {
				return fmt.Errorf("digite um numero inteiro valido")
			}

			return nil
		},
	}

	result, err := prompt.Run()
	if err != nil {
		return 0, err
	}

	value, err := strconv.Atoi(strings.TrimSpace(result))
	if err != nil {
		return 0, err
	}

	return value, nil
}

func PromptFloat64(label string) (float64, error) {
	prompt := promptui.Prompt{
		Label: label,
		Validate: func(input string) error {
			if requiredValidator(input) != nil {
				return requiredValidator(input)
			}

			if _, err := strconv.ParseFloat(strings.TrimSpace(input), 64); err != nil {
				return fmt.Errorf("digite um numero valido")
			}

			return nil
		},
	}

	result, err := prompt.Run()
	if err != nil {
		return 0, err
	}

	value, err := strconv.ParseFloat(strings.TrimSpace(result), 64)
	if err != nil {
		return 0, err
	}

	return value, nil
}

func PromptConfirm(label string) (bool, error) {
	prompt := promptui.Prompt{
		Label:     label,
		IsConfirm: true,
	}

	_, err := prompt.Run()
	if err != nil {
		if err == promptui.ErrAbort {
			return false, nil
		}

		return false, err
	}

	return true, nil
}

func requiredValidator(input string) error {
	if strings.TrimSpace(input) == "" {
		return fmt.Errorf("campo obrigatorio")
	}

	return nil
}
