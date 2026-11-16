package domain

type Student struct {
	id       string
	fullName string
}

func NewStudent(id, fullName string) Student {
	return Student{id: id, fullName: fullName}
}

func (s *Student) ID() string {
	return s.id
}

func (s *Student) FullName() string {
	return s.fullName
}

func (s *Student) SetFullName(fullName string) {
	s.fullName = fullName
}
