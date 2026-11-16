package domain

import "fmt"

// ScoreStrategy задаёт сменяемое правило расчёта итогового балла.
type ScoreStrategy interface {
	Calculate(gradebook *StudentGradebook, rubric Rubric) (float64, error)
}

func calculateWeightedScore(gradebook *StudentGradebook, rubric Rubric) (float64, error) {
	var weightedScore float64
	for _, grade := range gradebook.Grades() {
		criterionID := grade.CriterionID()
		criterion, exists := rubric.FindCriterion(criterionID)
		if !exists {
			return 0, fmt.Errorf("критерий %q отсутствует в рубрикаторе", criterionID)
		}

		maxScore := criterion.MaxScore()
		if maxScore <= 0 {
			return 0, fmt.Errorf("у критерия %q максимум баллов должен быть больше нуля", criterionID)
		}
		weightedScore += grade.Mark() / maxScore * criterion.Weight()
	}
	return weightedScore * 100, nil
}
