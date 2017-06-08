package shep

import (
	"testing"
)

var instance interface{}

func TestStringValidation(t *testing.T) {
	s := NewString()
	if err := s.Accept(&ValidationVisitor{42}); err != nil {
		t.Error("")
	}
	s = NewString(func(s *String) {
		s.MinLength = 2
		s.MaxLength = 3
	})
	instance = "a"
	if err := s.Accept(&ValidationVisitor{instance}); err == nil {
		t.Error("")
	}
	instance = "abcd"
	if err := s.Accept(&ValidationVisitor{instance}); err == nil {
		t.Error("")
	}
	instance = "(888)555-1212 ext. 532"
	s = NewString(func(s *String) {
		s.Pattern = `^(\\([0-9]{3}\\))?[0-9]{3}-[0-9]{4}$`
	})
	if err := s.Accept(&ValidationVisitor{instance}); err == nil {
		t.Error("")
	}
}

func TestIntegerValidation(t *testing.T) {
	i := NewInteger()
	instance = "42"
	if err := i.Accept(&ValidationVisitor{instance}); err != nil {
		t.Error("")
	}
}

func TestBooleanValidation(t *testing.T) {
	b := NewBoolean()
	instance = "true"
	if err := b.Accept(&ValidationVisitor{instance}); err != nil {
		t.Error("")
	}
}

func TestNullValidation(t *testing.T) {
	n := NewNull()
	instance = false
	if err := n.Accept(&ValidationVisitor{instance}); err != nil {
		t.Error("")
	}
}

func TestArrayValidation(t *testing.T) {
	a := NewArray(func(a *Array) {
		a.Items = NewInteger()
	})
	instance = []interface{}{}
	if err := a.Accept(&ValidationVisitor{instance}); err != nil {
		t.Error("")
	}
	instance = []interface{}{1, 2, 3, 4, 5}
	if err := a.Accept(&ValidationVisitor{instance}); err != nil {
		t.Error("")
	}
	a = NewArray(func(a *Array) {
		a.Items = []Component{
			NewInteger(),
			NewString(),
			NewString(func(s *String) {
				s.Enum = []interface{}{"Street", "Avenue", "Boulevard"}
			}),
			NewString(func(s *String) {
				s.Enum = []interface{}{"NW", "NE", "SW", "SE"}
			}),
		}
	})
	instance = []interface{}{24, "Sussex", "Drive"}
	if err := a.Accept(&ValidationVisitor{instance}); err == nil {
		t.Error("")
	}
}

func TestRecordValidation(t *testing.T) {
	r := NewRecord(func(r *Record) {
		r.Properties = Properties(map[string]Component{
			"number":     NewInteger(),
			"streetName": NewString(),
			"streetType": NewString(func(s *String) {
				s.Enum = []interface{}{"Street", "Avenue", "Boulevard"}
			}),
		})
	})
	instance = map[string]interface{}{
		"number":     "1600",
		"streetName": "Pennsylvania",
		"streetType": "Avenue",
	}
	if err := r.Accept(&ValidationVisitor{instance}); err != nil {
		t.Error("")
	}
}

func TestEnumValidation(t *testing.T) {
	s := NewString(func(s *String) {
		s.Enum = []interface{}{"red", "amber", "green"}
	})
	instance = "blue"
	if err := s.Accept(&ValidationVisitor{instance}); err == nil {
		t.Error("")
	}
}

func TestValidationVisitor(t *testing.T) {
	s := NewString(func(s *String) {
		s.MaxLength = 5
	})
	instance = "A green door"
	if err := s.Accept(&ValidationVisitor{instance}); err == nil {
		t.Error("")
	}
	r := NewRecord(func(r *Record) {
		r.Properties = Properties(map[string]Component{
			"firstName": NewString(func(s *String) {
				s.MaxLength = 1
			}),
			"lastName": NewString(),
			"age":      NewInteger(),
		})
		r.Required = []string{"firstName", "lastName"}
	})
	instance = map[string]interface{}{
		"firstName": "John",
		"lastName":  "Doe",
	}
	if err := r.Accept(&ValidationVisitor{instance}); err == nil {
		t.Error("")
	}
	a := NewArray(func(a *Array) {
		a.Items = []Component{
			NewString(func(s *String) {
				s.MaxLength = 10
			}),
			NewString(func(s *String) {
				s.MinLength = 4
			}),
		}
		a.MinItems = 1
	})
	instance = []interface{}{"cold", "ice"}
	if err := a.Accept(&ValidationVisitor{instance}); err == nil {
		t.Error("")
	}
	a = NewArray(func(a *Array) {
		a.Items = NewString(func(s *String) {
			s.MaxLength = 2
		})
	})
	if err := a.Accept(&ValidationVisitor{instance}); err == nil {
		t.Error("")
	}
	r = NewRecord(func(r *Record) {
		r.AllOf = []Component{
			NewRecord(func(r *Record) {
				r.Properties = Properties(map[string]Component{
					"firstName": NewString(),
					"lastName":  NewString(),
				})
			}),
			NewRecord(func(r *Record) {
				r.Properties = Properties(map[string]Component{
					"age": NewInteger(func(i *Integer) {
						i.MultipleOf = 5
						i.Maximum = 40
					}),
				})
				r.Required = []string{"age"}
			}),
		}
	})
	instance = map[string]interface{}{
		"firstName": "John",
		"lastName":  "Doe",
	}
	if err := r.Accept(&ValidationVisitor{instance}); err == nil {
		t.Error("")
	}
}
