package tests

import (
	"math"
	"testing"

	"kpo_hw2/domain"
)

func TestAssignGradeBounds(t *testing.T) {
	cases := []struct {
		name    string
		mark    float64
		wantErr bool
	}{
		{"zero", 0, false},
		{"maximum", 10, false},
		{"fraction", 8.5, false},
		{"negative", -1, true},
		{"above maximum", 12, true},
		{"NaN", math.NaN(), true},
		{"infinity", math.Inf(1), true},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			book, rubric := newGradebookAndRubric(t)
			err := book.AssignGrade("exam", tc.mark, rubric, nil)
			if (err != nil) != tc.wantErr {
				t.Fatalf("AssignGrade(%g): ошибка %v, ожидалась ошибка: %t", tc.mark, err, tc.wantErr)
			}
			if tc.wantErr {
				if len(book.Grades()) != 0 {
					t.Fatal("отклонённая оценка изменила ведомость")
				}
				return
			}
			grade, exists := book.FindGrade("exam")
			if !exists || grade.Mark() != tc.mark {
				t.Fatalf("оценка не сохранена: %v", book.Grades())
			}
		})
	}
}

func TestAssignGradeRejectsUnknownAndDuplicateCriterion(t *testing.T) {
	book, rubric := newGradebookAndRubric(t)
	if err := book.AssignGrade("unknown", 8, rubric, nil); err == nil {
		t.Fatal("ожидалась ошибка для неизвестного критерия")
	}
	if len(book.Grades()) != 0 {
		t.Fatal("оценка неизвестного критерия попала в ведомость")
	}
	if err := book.AssignGrade("exam", 8, rubric, nil); err != nil {
		t.Fatal(err)
	}
	if err := book.AssignGrade("exam", 9, rubric, nil); err == nil {
		t.Fatal("ожидалась ошибка для повторной оценки")
	}
	grade, exists := book.FindGrade("exam")
	if !exists || grade.Mark() != 8 || len(book.Grades()) != 1 {
		t.Fatal("повторное выставление изменило исходную оценку")
	}
}

func newGradebookAndRubric(t *testing.T) (*domain.StudentGradebook, domain.Rubric) {
	t.Helper()
	var rubric domain.Rubric
	for _, spec := range []struct {
		id       string
		maxScore float64
		weight   float64
	}{
		{"exam", 10, 0.6},
		{"project", 20, 0.4},
	} {
		criterion, err := domain.NewCriterion(spec.id, spec.id, spec.maxScore, spec.weight)
		if err != nil {
			t.Fatal(err)
		}
		if err := rubric.AddCriterion(criterion); err != nil {
			t.Fatal(err)
		}
	}
	return domain.NewStudentGradebook(domain.NewStudent("s1", "Анна")), rubric
}
