package tests

import (
	"math"
	"testing"

	"kpo_hw2/domain"
)

func TestFinalScoreReflectsGradeChanges(t *testing.T) {
	book, rubric := newGradebookAndRubric(t)
	if err := book.AssignGrade("exam", 8, rubric, nil); err != nil {
		t.Fatal(err)
	}
	if err := book.AssignGrade("project", 10, rubric, nil); err != nil {
		t.Fatal(err)
	}
	calculator := domain.NewFinalScoreCalculator(domain.OfficialScoreStrategy{})
	// 8/10 * 60 + 10/20 * 40 = 68.
	assertFinalScore(t, calculator, book, rubric, 68)

	if err := book.ChangeMark("exam", 9, rubric); err != nil {
		t.Fatal(err)
	}
	assertFinalScore(t, calculator, book, rubric, 74)

	if err := book.CancelGrade("project"); err != nil {
		t.Fatal(err)
	}
	if _, exists := book.FindGrade("project"); exists {
		t.Fatal("отменённая оценка осталась в ведомости")
	}
	if _, err := calculator.Calculate(book, rubric); err == nil {
		t.Fatal("официальный итог требует всех оценок")
	}
	preliminary := domain.NewFinalScoreCalculator(domain.PreliminaryScoreStrategy{})
	assertFinalScore(t, preliminary, book, rubric, 54)
}

func assertFinalScore(t *testing.T, calculator domain.FinalScoreCalculator, book *domain.StudentGradebook, rubric domain.Rubric, want float64) {
	t.Helper()
	score, err := calculator.Calculate(book, rubric)
	if err != nil {
		t.Fatal(err)
	}
	if math.IsNaN(score.Value()) || math.Abs(score.Value()-want) > 1e-9 {
		t.Fatalf("итог %g, ожидался %g", score.Value(), want)
	}
	if score.Student() != book.Student() {
		t.Fatal("итог относится к другому студенту")
	}
}
