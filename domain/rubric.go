package domain

import (
	"errors"
	"fmt"
	"slices"
)

type Rubric struct {
	criteria []Criterion
}

// Clone создаёт независимую копию рубрикатора по паттерну «Прототип».
func (r *Rubric) Clone() Rubric {
	return Rubric{criteria: r.Criteria()}
}

func (r *Rubric) Criteria() []Criterion {
	return slices.Clone(r.criteria)
}

func (r *Rubric) FindCriterion(id string) (Criterion, bool) {
	for _, criterion := range r.criteria {
		if criterion.ID() == id {
			return criterion, true
		}
	}
	return Criterion{}, false
}

func (r *Rubric) AddCriterion(criterion Criterion) error {
	if err := validateMaxScore(criterion.MaxScore()); err != nil {
		return err
	}
	id := criterion.ID()
	if _, exists := r.FindCriterion(id); exists {
		return fmt.Errorf("критерий с идентификатором %q уже существует", id)
	}
	if err := r.validateTotalWeight(criterion); err != nil {
		return err
	}
	r.criteria = append(r.criteria, criterion)
	return nil
}

func (r *Rubric) ChangeWeight(criterionID string, weight float64) error {
	for index, criterion := range r.criteria {
		if criterion.ID() != criterionID {
			continue
		}
		if err := criterion.SetWeight(weight); err != nil {
			return err
		}
		if err := r.validateTotalWeight(criterion); err != nil {
			return err
		}
		r.criteria[index] = criterion
		return nil
	}
	return fmt.Errorf("критерий %q отсутствует в рубрикаторе", criterionID)
}

func (r *Rubric) validateTotalWeight(candidate Criterion) error {
	totalWeight := candidate.Weight()
	for _, criterion := range r.criteria {
		if criterion.ID() != candidate.ID() {
			totalWeight += criterion.Weight()
		}
	}
	// Допуск учитывает погрешность сложения весов типа float64.
	const tolerance = 1e-12
	if totalWeight > 1+tolerance {
		return errors.New("суммарный вес критериев не должен превышать 100%")
	}
	return nil
}
