package application

import (
	"errors"
	"fmt"
	"strings"
	"time"

	"kpo_hw2/domain"
)

// Application предоставляет сценарии приложения через общий фасад.
type Application struct {
	gradebooks  GradebookRepository
	rubrics     RubricRepository
	calculators ScoreCalculators
}

func NewApplication(gradebooks GradebookRepository, rubrics RubricRepository,
	calculators ScoreCalculators, ) *Application {return &Application{gradebooks:  gradebooks, rubrics:
		rubrics, calculators: calculators,}
}

// AddStudent создаёт студента вместе с его пустой ведомостью.
func (app *Application) AddStudent(id, fullName string) error {
	id = strings.TrimSpace(id)
	fullName = strings.TrimSpace(fullName)
	if id == "" || fullName == "" {
		return errors.New("ID и ФИО студента должны быть заполнены")
	}
	student := domain.NewStudent(id, fullName)
	return app.gradebooks.Add(domain.NewStudentGradebook(student))
}

func (app *Application) ListStudents() ([]domain.Student, error) {
	gradebooks, err := app.gradebooks.List()
	if err != nil {
		return nil, err
	}
	students := make([]domain.Student, 0, len(gradebooks))
	for _, gradebook := range gradebooks {
		students = append(students, gradebook.Student())
	}
	return students, nil
}

func (app *Application) AddCriterion(id, name string, maxScore, weight float64) error {
	id = strings.TrimSpace(id)
	name = strings.TrimSpace(name)
	if id == "" || name == "" {
		return errors.New("ID и название критерия должны быть заполнены")
	}
	criterion, err := domain.NewCriterion(id, name, maxScore, weight)
	if err != nil {
		return err
	}
	rubric, err := app.rubrics.Get()
	if err != nil {
		return err
	}
	if err := rubric.AddCriterion(criterion); err != nil {
		return err
	}
	return app.rubrics.Save(rubric)
}

func (app *Application) ListCriteria() ([]domain.Criterion, error) {
	rubric, err := app.rubrics.Get()
	if err != nil {
		return nil, err
	}
	return rubric.Criteria(), nil
}

func (app *Application) GetGrades(studentID string) ([]domain.Grade, error) {
	gradebook, err := app.gradebooks.Get(studentID)
	if err != nil {
		return nil, err
	}
	return gradebook.Grades(), nil
}

func (app *Application) ChangeMark(studentID, criterionID string, mark float64) error {
	gradebook, err := app.gradebooks.Get(studentID)
	if err != nil {
		return err
	}
	rubric, err := app.rubrics.Get()
	if err != nil {
		return err
	}
	if err := gradebook.ChangeMark(criterionID, mark, rubric); err != nil {
		return err
	}
	return app.gradebooks.Save(gradebook)
}

func (app *Application) CancelGrade(studentID, criterionID string) error {
	gradebook, err := app.gradebooks.Get(studentID)
	if err != nil {
		return err
	}
	if err := gradebook.CancelGrade(criterionID); err != nil {
		return err
	}
	return app.gradebooks.Save(gradebook)
}

func (app *Application) CalculateFinalScore(studentID string, mode ScoreMode) (domain.FinalScore, error) {
	calculator, err := app.calculators.forMode(mode)
	if err != nil {
		return domain.FinalScore{}, err
	}
	gradebook, err := app.gradebooks.Get(studentID)
	if err != nil {
		return domain.FinalScore{}, err
	}
	rubric, err := app.rubrics.Get()
	if err != nil {
		return domain.FinalScore{}, err
	}
	return calculator.Calculate(gradebook, rubric)
}

// ExportResults сначала рассчитывает все итоги, затем передаёт их адаптеру экспорта.
func (app *Application) ExportResults(exporter ResultsExporter, mode ScoreMode) error {
	calculator, err := app.calculators.forMode(mode)
	if err != nil {
		return err
	}
	gradebooks, err := app.gradebooks.List()
	if err != nil {
		return err
	}
	rubric, err := app.rubrics.Get()
	if err != nil {
		return err
	}
	results := make([]domain.FinalScore, 0, len(gradebooks))
	for _, gradebook := range gradebooks {
		result, err := calculator.Calculate(&gradebook, rubric)
		if err != nil {
			student := gradebook.Student()
			return fmt.Errorf("итог студента %q: %w", student.ID(), err)
		}
		results = append(results, result)
	}
	return exporter.Export(results)
}

// AssignGrade выставляет оценку с текущей датой и сохраняет изменённую ведомость.
func (app *Application) AssignGrade(
	studentID string,
	criterionID string,
	mark float64,
) error {
	// Получаем ведомость студента.
	gradebook, err := app.gradebooks.Get(studentID)
	if err != nil {
		return err
	}

	// Получаем действующий рубрикатор.
	rubric, err := app.rubrics.Get()
	if err != nil {
		return err
	}

	// Домен проверяет правила и добавляет оценку.
	gradedAt := time.Now()
	if err := gradebook.AssignGrade(criterionID, mark, rubric, &gradedAt); err != nil {
		return err
	}

	// Сохраняем изменённую ведомость.
	return app.gradebooks.Save(gradebook)
}
