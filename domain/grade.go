package domain

import (
	"fmt"
	"time"
)

type Grade struct {
	criterionID string
	mark        float64
	gradedAt    *time.Time
}

func (g *Grade) CriterionID() string {
	return g.criterionID
}

func (g *Grade) Mark() float64 {
	return g.mark
}

func (g *Grade) GradedAt() (time.Time, bool) {
	if g.gradedAt == nil {
		return time.Time{}, false
	}
	return *g.gradedAt, true
}

func (g *Grade) changeMark(mark, maxScore float64) error {
	if !(mark >= 0 && mark <= maxScore) {
		return fmt.Errorf("балл должен быть в диапазоне от 0 до %g", maxScore)
	}
	g.mark = mark
	return nil
}
