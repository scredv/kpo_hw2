package application

import "kpo_hw2/domain"

// ResultsExporter — исходящий порт для экспорта рассчитанных итогов.
type ResultsExporter interface {
	Export(results []domain.FinalScore) error
}
