package domain

// PreliminaryScoreStrategy считает итог по выставленным оценкам.
type PreliminaryScoreStrategy struct{}

func (PreliminaryScoreStrategy) Calculate(gradebook *StudentGradebook, rubric Rubric) (float64, error) {
	return calculateWeightedScore(gradebook, rubric)
}
