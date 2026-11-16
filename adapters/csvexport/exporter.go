package csvexport

import (
	"encoding/csv"
	"os"
	"path/filepath"
	"strconv"

	"kpo_hw2/domain"
)

// Exporter сохраняет итоги студентов в CSV-файл.
type Exporter struct {
	path string
}

func NewExporter(path string) *Exporter {
	return &Exporter{path: path}
}

func (e *Exporter) Export(results []domain.FinalScore) error {
	// Записываем временный файл в том же каталоге для атомарной замены результата.
	file, err := os.CreateTemp(filepath.Dir(e.path), ".results-*.csv")
	if err != nil {
		return err
	}
	defer os.Remove(file.Name())
	defer file.Close()

	writer := csv.NewWriter(file)
	if err := writer.Write([]string{"student_id", "full_name", "final_score"}); err != nil {
		return err
	}
	for _, result := range results {
		student := result.Student()
		if err := writer.Write([]string{
			student.ID(), student.FullName(), strconv.FormatFloat(result.Value(), 'f', -1, 64),
		}); err != nil {
			return err
		}
	}
	writer.Flush()
	if err := writer.Error(); err != nil {
		return err
	}
	if err := file.Close(); err != nil {
		return err
	}
	return os.Rename(file.Name(), e.path)
}
