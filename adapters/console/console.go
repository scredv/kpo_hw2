package console

import (
	"bufio"
	"fmt"
	"io"
	"math"
	"strconv"
	"strings"

	"kpo_hw2/adapters/csvexport"
	"kpo_hw2/application"
)

// Console читает команды, вызывает сценарии приложения и выводит результат.
type Console struct {
	app    *application.Application
	input  io.Reader
	output io.Writer
}

func NewConsole(app *application.Application, input io.Reader, output io.Writer) *Console {
	return &Console{app: app, input: input, output: output}
}

const help = `Команды:
student-add <id> <ФИО> — добавить студента
students — показать студентов
criterion-add <id> <максимум> <вес> <название> — добавить критерий
criteria — показать критерии
grade <id студента> <id критерия> <балл> — выставить оценку
grades <id студента> — показать оценки
grade-change <id студента> <id критерия> <балл> — исправить балл
grade-cancel <id студента> <id критерия> — отменить оценку
score <id студента> [official|preliminary] — рассчитать итог
export [--mode official|preliminary] <путь.csv> — экспортировать итоги всех студентов
help — показать команды
exit — завершить работу
Вес задаётся долей от 0 до 1, например 0.4 означает 40%.
Дробные числа вводятся через точку. ФИО, название и путь вводятся без кавычек.
Дата и время выставления оценки сохраняются автоматически.
По умолчанию используется official — окончательный итог, требующий всех оценок.
preliminary — предварительный итог; пропущенные оценки дают нулевой вклад.`

// Run обрабатывает команды до exit или конца ввода.
// Ошибка команды выводится пользователю, после чего можно ввести следующую.
func (c *Console) Run() error {
	if _, err := fmt.Fprintln(c.output, help); err != nil {
		return err
	}
	scanner := bufio.NewScanner(c.input)
	scanner.Buffer(make([]byte, 4096), 1024*1024)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "exit" {
			return nil
		}
		if line == "" {
			continue
		}
		result, err := c.Execute(line)
		if err != nil {
			result = "Ошибка: " + err.Error()
		}
		if _, err := fmt.Fprintln(c.output, result); err != nil {
			return err
		}
	}
	return scanner.Err()
}

// Execute разбирает одну команду и возвращает текст для пользователя.
func (c *Console) Execute(line string) (string, error) {
	args := strings.Fields(line)
	if len(args) == 0 {
		return "", nil
	}
	switch args[0] {
	case "help":
		if err := checkArguments(args, 1); err != nil {
			return "", err
		}
		return help, nil
	case "student-add":
		if len(args) < 3 {
			return "", fmt.Errorf("формат: student-add <id> <ФИО>")
		}
		return "Студент добавлен.", c.app.AddStudent(args[1], strings.Join(args[2:], " "))
	case "students":
		if err := checkArguments(args, 1); err != nil {
			return "", err
		}
		return c.listStudents()
	case "criterion-add":
		return c.addCriterion(args)
	case "criteria":
		if err := checkArguments(args, 1); err != nil {
			return "", err
		}
		return c.listCriteria()
	case "grade", "grade-change":
		return c.setGrade(args)
	case "grades":
		if err := checkArguments(args, 2); err != nil {
			return "", err
		}
		return c.listGrades(args[1])
	case "grade-cancel":
		if err := checkArguments(args, 3); err != nil {
			return "", err
		}
		return "Оценка отменена.", c.app.CancelGrade(args[1], args[2])
	case "score":
		if len(args) < 2 || len(args) > 3 {
			return "", fmt.Errorf("формат: score <id студента> [official|preliminary]")
		}
		mode := application.OfficialMode
		if len(args) == 3 {
			mode = application.ScoreMode(args[2])
		}
		result, err := c.app.CalculateFinalScore(args[1], mode)
		if err != nil {
			return "", err
		}
		student := result.Student()
		return fmt.Sprintf("%s (%s): %s итог %g", student.FullName(), student.ID(), scoreModeLabel(mode), result.Value()), nil
	case "export":
		// Остаток строки сохраняем целиком, чтобы поддерживать пробелы в пути.
		path := strings.TrimSpace(strings.TrimPrefix(strings.TrimSpace(line), args[0]))
		mode := application.OfficialMode
		if len(args) > 1 && args[1] == "--mode" {
			if len(args) < 4 {
				return "", fmt.Errorf("формат: export [--mode official|preliminary] <путь.csv>")
			}
			mode = application.ScoreMode(args[2])
			path = strings.TrimSpace(strings.TrimPrefix(path, "--mode"))
			path = strings.TrimSpace(strings.TrimPrefix(path, args[2]))
		}
		if path == "" {
			return "", fmt.Errorf("формат: export [--mode official|preliminary] <путь.csv>")
		}
		if err := c.app.ExportResults(csvexport.NewExporter(path), mode); err != nil {
			return "", err
		}
		return "Экспорт завершён (" + scoreModeLabel(mode) + " итог): " + path, nil
	default:
		return "", fmt.Errorf("неизвестная команда %q; список команд доступен через help", args[0])
	}
}

func scoreModeLabel(mode application.ScoreMode) string {
	if mode == application.PreliminaryMode {
		return "предварительный"
	}
	return "окончательный"
}

func (c *Console) addCriterion(args []string) (string, error) {
	if len(args) < 5 {
		return "", fmt.Errorf("формат: criterion-add <id> <максимум> <вес> <название>")
	}
	maxScore, err := parseNumber(args[2])
	if err != nil {
		return "", err
	}
	weight, err := parseNumber(args[3])
	if err != nil {
		return "", err
	}
	return "Критерий добавлен.", c.app.AddCriterion(args[1], strings.Join(args[4:], " "), maxScore, weight)
}

func (c *Console) setGrade(args []string) (string, error) {
	if err := checkArguments(args, 4); err != nil {
		return "", err
	}
	mark, err := parseNumber(args[3])
	if err != nil {
		return "", err
	}
	if args[0] == "grade-change" {
		return "Балл изменён.", c.app.ChangeMark(args[1], args[2], mark)
	}
	return "Оценка выставлена.", c.app.AssignGrade(args[1], args[2], mark)
}

func (c *Console) listStudents() (string, error) {
	students, err := c.app.ListStudents()
	if err != nil {
		return "", err
	}
	if len(students) == 0 {
		return "Студентов пока нет.", nil
	}
	var output strings.Builder
	for _, student := range students {
		fmt.Fprintf(&output, "%s\t%s\n", student.ID(), student.FullName())
	}
	return strings.TrimSuffix(output.String(), "\n"), nil
}

func (c *Console) listCriteria() (string, error) {
	criteria, err := c.app.ListCriteria()
	if err != nil {
		return "", err
	}
	if len(criteria) == 0 {
		return "Критериев пока нет.", nil
	}
	var output strings.Builder
	for _, criterion := range criteria {
		fmt.Fprintf(&output, "%s\t%s\tмаксимум %g\tвес %g%%\n",
			criterion.ID(), criterion.Name(), criterion.MaxScore(), criterion.Weight()*100)
	}
	return strings.TrimSuffix(output.String(), "\n"), nil
}

func (c *Console) listGrades(studentID string) (string, error) {
	grades, err := c.app.GetGrades(studentID)
	if err != nil {
		return "", err
	}
	if len(grades) == 0 {
		return "Оценок пока нет.", nil
	}
	var output strings.Builder
	for _, grade := range grades {
		date := "—"
		if gradedAt, exists := grade.GradedAt(); exists {
			date = gradedAt.Format("2006-01-02 15:04:05")
		}
		fmt.Fprintf(&output, "%s\tбалл %g\tдата %s\n",
			grade.CriterionID(), grade.Mark(), date)
	}
	return strings.TrimSuffix(output.String(), "\n"), nil
}

func checkArguments(args []string, count int) error {
	if len(args) != count {
		return fmt.Errorf("команда %s ожидает аргументов: %d; формат доступен через help", args[0], count-1)
	}
	return nil
}

func parseNumber(value string) (float64, error) {
	number, err := strconv.ParseFloat(value, 64)
	if err != nil || math.IsNaN(number) || math.IsInf(number, 0) {
		return 0, fmt.Errorf("%q должно быть конечным числом, например 8.5", value)
	}
	return number, nil
}
