package models

type Handler struct {
	FailMsg    string
	SuccessMsg string
	Processor  func(any) (string, error)
}
