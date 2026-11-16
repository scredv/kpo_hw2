package domain

import (
	"fmt"
	"slices"
	"time"
)

type StudentGradebook struct {
	student Student
	grades  []Grade
}

func NewStudentGradebook(student Student) *StudentGradebook {
	return &StudentGradebook{student: student}
}

func (b *StudentGradebook) Student() Student {
	return b.student
}

// Clone создаёт независимую копию ведомости для хранения.
func (b *StudentGradebook) Clone() StudentGradebook {
	return StudentGradebook{student: b.student, grades: b.Grades()}
}

func (b *StudentGradebook) Grades() []Grade {
	return slices.Clone(b.grades)
}

func (b *StudentGradebook) FindGrade(criterionID string) (Grade, bool) {
	index := b.gradeIndex(criterionID)
	if index < 0 {
		return Grade{}, false
	}
	return b.grades[index], true
}

func (b *StudentGradebook) AssignGrade(
	criterionID string, mark float64, rubric Rubric, gradedAt *time.Time,
) error {
	criterion, exists := rubric.FindCriterion(criterionID)
	if !exists {
		return fmt.Errorf("критерий %q отсутствует в рубрикаторе", criterionID)
	}
	if b.gradeIndex(criterionID) >= 0 {
		return fmt.Errorf("по критерию %q уже выставлена оценка", criterionID)
	}

	grade := Grade{criterionID: criterionID}
	if err := grade.changeMark(mark, criterion.MaxScore()); err != nil {
		return err
	}
	if gradedAt != nil {
		date := *gradedAt
		grade.gradedAt = &date
	}

	b.grades = append(b.grades, grade)
	return nil
}

func (b *StudentGradebook) ChangeMark(criterionID string, mark float64, rubric Rubric) error {
	index := b.gradeIndex(criterionID)
	if index < 0 {
		return fmt.Errorf("по критерию %q ещё нет оценки", criterionID)
	}
	criterion, exists := rubric.FindCriterion(criterionID)
	if !exists {
		return fmt.Errorf("критерий %q отсутствует в рубрикаторе", criterionID)
	}
	return b.grades[index].changeMark(mark, criterion.MaxScore())
}

func (b *StudentGradebook) CancelGrade(criterionID string) error {
	index := b.gradeIndex(criterionID)
	if index < 0 {
		return fmt.Errorf("по критерию %q ещё нет оценки", criterionID)
	}
	b.grades = slices.Delete(b.grades, index, index+1)
	return nil
}

func (b *StudentGradebook) gradeIndex(criterionID string) int {
	for index, grade := range b.grades {
		if grade.CriterionID() == criterionID {
			return index
		}
	}
	return -1
}
