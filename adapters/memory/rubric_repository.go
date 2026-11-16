package memory

import (
	"kpo_hw2/domain"
)

// RubricRepository хранит текущий рубрикатор в памяти.
type RubricRepository struct {
	rubric domain.Rubric
}

func (r *RubricRepository) Get() (domain.Rubric, error) {
	return r.rubric.Clone(), nil
}

func (r *RubricRepository) Save(rubric domain.Rubric) error {
	r.rubric = rubric.Clone()
	return nil
}
