package application

import "kpo_hw2/domain"

// RubricRepository — исходящий порт для получения действующего рубрикатора.
type RubricRepository interface {
	// Get возвращает независимую копию рубрикатора или ошибку.
	// Для нового курса возвращает пустой рубрикатор.
	Get() (domain.Rubric, error)

	// Save сохраняет рубрикатор. При ошибке сохранённые данные остаются прежними.
	Save(rubric domain.Rubric) error
}
