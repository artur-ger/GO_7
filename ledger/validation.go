package ledger

// Validatable — объекты с методом Validate.
type Validatable interface {
	Validate() error
}

// CheckValid вызывает Validate у объекта.
func CheckValid(v Validatable) error {
	return v.Validate()
}
