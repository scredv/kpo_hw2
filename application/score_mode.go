package application

import (
	"fmt"

	"kpo_hw2/domain"
)

type ScoreMode string

const (
	OfficialMode    ScoreMode = "official"
	PreliminaryMode ScoreMode = "preliminary"
)

type ScoreCalculators struct {
	Official    domain.FinalScoreCalculator
	Preliminary domain.FinalScoreCalculator
}

func (c ScoreCalculators) forMode(mode ScoreMode) (domain.FinalScoreCalculator, error) {
	switch mode {
	case OfficialMode:
		return c.Official, nil
	case PreliminaryMode:
		return c.Preliminary, nil
	default:
		return domain.FinalScoreCalculator{}, fmt.Errorf("неизвестный режим %q; используйте official или preliminary", mode)
	}
}
