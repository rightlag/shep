package primitive

import (
	"testing"
)

var instance interface{}

func TestStringValidation(t *testing.T) {
	s := &String{}
	if err := s.Accept(&ValidationVisitor{42}); err != nil {
		t.Error("")
	}
	s = &String{
		MinLength: 2,
		MaxLength: 3,
	}
	instance = "a"
	if err := s.Accept(&ValidationVisitor{instance}); err == nil {
		t.Error("")
	}
	instance = "abcd"
	if err := s.Accept(&ValidationVisitor{instance}); err == nil {
		t.Error("")
	}
	instance = "(888)555-1212 ext. 532"
	s = &String{
		Pattern: `^(\\([0-9]{3}\\))?[0-9]{3}-[0-9]{4}$`,
	}
	if err := s.Accept(&ValidationVisitor{instance}); err == nil {
		t.Error("")
	}
}

func TestIntegerValidation(t *testing.T) {
	i := &Integer{}
	instance = "42"
	if err := i.Accept(&ValidationVisitor{instance}); err != nil {
		t.Error("")
	}
}

func TestEnumValidation(t *testing.T) {
	s := &String{
		Primitive: Primitive{
			Type: "string",
			Enum: []interface{}{"red", "amber", "green"},
		},
	}
	instance = "blue"
	if err := s.Accept(&ValidationVisitor{instance}); err == nil {
		t.Error("")
	}
}

func TestValidationVisitor(t *testing.T) {
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
		"firstName": "John",
		"lastName":  "Doe",
	}
	if err := r.Accept(&ValidationVisitor{instance}); err == nil {
		t.Error("")
	}
	a := &Array{
		Items: []Component{
			&String{MaxLength: 10},
			&String{MinLength: 4},
		},
		MinItems: 1,
	}
	instance = []interface{}{"cold", "ice"}
	if err := a.Accept(&ValidationVisitor{instance}); err == nil {
		t.Error("")
	}
	a = &Array{
		Items: &String{MaxLength: 2},
	}
	if err := a.Accept(&ValidationVisitor{instance}); err == nil {
		t.Error("")
	}
	p := &Primitive{
		AllOf: &AllOf{[]Component{
			&Record{
				Properties: Properties(map[string]Component{
					"firstName": &String{},
					"lastName":  &String{},
				}),
			},
			&Record{
				Properties: Properties(map[string]Component{
					"age": &Integer{
						MultipleOf: 5,
						Maximum:    40,
					},
				}),
				Required: []string{"age"},
			},
		}},
	}
	instance = map[string]interface{}{
		"firstName": "John",
		"lastName":  "Doe",
	}
	if err := p.Accept(&ValidationVisitor{instance}); err == nil {
		t.Error("")
	}
}
