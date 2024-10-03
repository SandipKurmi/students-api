package types

type Student struct {
	ID   int64
	Name string `validate:"required"`
	Email string `validate:"required,email"`
	Age  int `validate:"required"`
}