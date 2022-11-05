package internal

import (
	"fmt"
	"reflect"
)

type RowParsingError struct {
	Column  string
	RowType reflect.Type
}

func (r *RowParsingError) Error() string {
	return fmt.Sprintf("Error when parsing %s in %s", r.Column, r.RowType)
}
