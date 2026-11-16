package domain

import (
	"errors"
	"math"
)

type Criterion struct {
	id       string
	name     string
	maxScore float64
	weight   float64 // Доля критерия в итоговой оценке. Например, 0.4 означает 40%.
}

func NewCriterion(id, name string, maxScore, weight float64) (Criterion, error) {
	if err := validateMaxScore(maxScore); err != nil {
		return Criterion{}, err
	}

	criterion := Criterion{id: id, name: name, maxScore: maxScore}
	if err := criterion.SetWeight(weight); err != nil {
		return Criterion{}, err
	}
	return criterion, nil
}

func validateMaxScore(maxScore float64) error {
	if maxScore <= 0 || math.IsNaN(maxScore) || math.IsInf(maxScore, 0) {
		return errors.New("максимум баллов должен быть положительным конечным числом")
	}
	return nil
}

func (c *Criterion) ID() string {
	return c.id
}

func (c *Criterion) Name() string {
	return c.name
}

func (c *Criterion) MaxScore() float64 {
	return c.maxScore
}

func (c *Criterion) Weight() float64 {
	return c.weight
}

func (c *Criterion) SetName(name string) {
	c.name = name
}

func (c *Criterion) SetWeight(weight float64) error {
	if !(weight >= 0 && weight <= 1) {
		return errors.New("вес критерия должен быть в диапазоне от 0 до 1")
	}
	c.weight = weight
	return nil
}
