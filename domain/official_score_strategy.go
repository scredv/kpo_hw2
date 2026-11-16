package domain

import "fmt"

// OfficialScoreStrategy требует оценку по каждому критерию рубрикатора.
type OfficialScoreStrategy struct{}

func (OfficialScoreStrategy) Calculate(gradebook *StudentGradebook, rubric Rubric) (float64, error) {
	for _, criterion := range rubric.Criteria() {
		if _, exists := gradebook.FindGrade(criterion.ID()); !exists {
			return 0, fmt.Errorf("для окончательного итога нужна оценка по критерию %q", criterion.Name())
		}
	}
	return calculateWeightedScore(gradebook, rubric)
}
