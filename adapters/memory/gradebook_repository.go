package memory

import (
	"errors"
	"fmt"
	"slices"

	"kpo_hw2/domain"
)

// GradebookRepository хранит ведомости в памяти. Нулевое значение готово к использованию.
type GradebookRepository struct {
	gradebooks map[string]domain.StudentGradebook
}

func (r *GradebookRepository) Add(gradebook *domain.StudentGradebook) error {
	if gradebook == nil {
		return errors.New("ведомость не передана")
	}
	student := gradebook.Student()
	id := student.ID()
	if _, exists := r.gradebooks[id]; exists {
		return fmt.Errorf("студент с ID %q уже существует", id)
	}
	if r.gradebooks == nil {
		r.gradebooks = make(map[string]domain.StudentGradebook)
	}
	r.gradebooks[id] = gradebook.Clone()
	return nil
}

func (r *GradebookRepository) Get(studentID string) (*domain.StudentGradebook, error) {
	gradebook, exists := r.gradebooks[studentID]
	if !exists {
		return nil, fmt.Errorf("ведомость студента %q не найдена", studentID)
	}
	snapshot := gradebook.Clone()
	return &snapshot, nil
}

func (r *GradebookRepository) List() ([]domain.StudentGradebook, error) {
	ids := make([]string, 0, len(r.gradebooks))
	for id := range r.gradebooks {
		ids = append(ids, id)
	}
	slices.Sort(ids)

	gradebooks := make([]domain.StudentGradebook, 0, len(ids))
	for _, id := range ids {
		gradebook := r.gradebooks[id]
		gradebooks = append(gradebooks, gradebook.Clone())
	}
	return gradebooks, nil
}

func (r *GradebookRepository) Save(gradebook *domain.StudentGradebook) error {
	if gradebook == nil {
		return errors.New("ведомость не передана")
	}
	student := gradebook.Student()
	id := student.ID()
	if _, exists := r.gradebooks[id]; !exists {
		return fmt.Errorf("ведомость студента %q не найдена", id)
	}
	r.gradebooks[id] = gradebook.Clone()
	return nil
}
