package main

import (
	"fmt"
	"io"
	"os"

	"go.uber.org/dig"

	"kpo_hw2/adapters/console"
	"kpo_hw2/adapters/memory"
	"kpo_hw2/application"
	"kpo_hw2/domain"
)

func run() error {
	container := dig.New()
	constructors := []any{
		func() application.GradebookRepository { return &memory.GradebookRepository{} },
		func() application.RubricRepository { return &memory.RubricRepository{} },
		func() application.ScoreCalculators {
			return application.ScoreCalculators{
				Official:    domain.NewFinalScoreCalculator(domain.OfficialScoreStrategy{}),
				Preliminary: domain.NewFinalScoreCalculator(domain.PreliminaryScoreStrategy{}),
			}
		},
		application.NewApplication,
		func() io.Reader { return os.Stdin },
		func() io.Writer { return os.Stdout },
		console.NewConsole,
	}

	for _, constructor := range constructors {
		if err := container.Provide(constructor); err != nil {
			return fmt.Errorf("регистрация зависимости: %w", err)
		}
	}

	return container.Invoke(func(ui *console.Console) error {
		return ui.Run()
	})
}
