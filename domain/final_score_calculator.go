package domain

type FinalScoreCalculator struct {
	strategy ScoreStrategy
}

func NewFinalScoreCalculator(strategy ScoreStrategy) FinalScoreCalculator {
	return FinalScoreCalculator{strategy: strategy}
}

func (c FinalScoreCalculator) Calculate(gradebook *StudentGradebook, rubric Rubric) (FinalScore, error) {
	strategy := c.strategy
	if strategy == nil {
		strategy = PreliminaryScoreStrategy{}
	}

	value, err := strategy.Calculate(gradebook, rubric)
	if err != nil {
		return FinalScore{}, err
	}
	return FinalScore{student: gradebook.Student(), value: value}, nil
}
