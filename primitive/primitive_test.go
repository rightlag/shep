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
			},
		}},
	}
	instance = map[string]interface{}{
		"firstName": "John",
		"lastName":  "Doe",
		"age":       45,
	}
	if err := p.Accept(&ValidationVisitor{instance}); err == nil {
		t.Error("")
	}
}
