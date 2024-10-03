package types

type Student struct {
	ID   string
	Name string `validate:"required"`
	Email string `validate:"required,email"`
	Age  int `validate:"required"`
}