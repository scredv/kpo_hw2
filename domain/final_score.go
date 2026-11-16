package domain

type FinalScore struct {
	student Student
	value   float64 // Взвешенный итог в баллах на шкале от 0 до 100 при сумме весов 1.
}

func (s FinalScore) Student() Student {
	return s.student
}

func (s FinalScore) Value() float64 {
	return s.value
}
