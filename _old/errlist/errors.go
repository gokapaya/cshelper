package errlist

import "fmt"

type Err string
type ErrList []error

func (e Err) Error() string {
	return string(e)
}

func (el ErrList) NotEmpty() bool {
	if len(el) != 0 {
		return true
	}
	return false
}

func (el ErrList) Error() string {
	return fmt.Sprintf("Found %v erros.", len(el))
}
