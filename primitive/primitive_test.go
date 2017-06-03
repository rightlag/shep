package primitive

import (
	"testing"
)

func TestValidationVisitor(t *testing.T) {
	var instance interface{}
	s := &String{MaxLength: 5}
	instance = "A green door"
	if err := s.Accept(&ValidationVisitor{instance}); err == nil {
		t.Error("")
	}
	r := &Record{
		Properties: Properties(map[string]Component{
			"firstName": &String{MaxLength: 1},
			"lastName":  &String{},
			"age":       &Integer{},
		}),
		Required: []string{"firstName", "lastName"},
	}
	instance = map[string]interface{}{
		"properties": map[string]interface{}{
			"firstName": "John",
			"lastName":  "Doe",
		},
	}
	if err := r.Accept(&ValidationVisitor{instance}); err == nil {
		t.Error("")
	}
}
