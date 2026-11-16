package application

import "kpo_hw2/domain"

// GradebookRepository — исходящий порт для хранения ведомостей студентов.
type GradebookRepository interface {
	// Add создаёт ведомость. При занятом ID возвращает ошибку без изменения данных.
	Add(gradebook *domain.StudentGradebook) error

	// List возвращает независимые копии всех ведомостей.
	List() ([]domain.StudentGradebook, error)

	// Get возвращает независимую ведомость для изменения.
	// Если ведомость не найдена, возвращает ошибку.
	Get(studentID string) (*domain.StudentGradebook, error)

	// Save сохраняет ведомость. При ошибке сохранённые данные остаются прежними.
	Save(gradebook *domain.StudentGradebook) error
}
