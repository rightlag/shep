package primitive

import (
	"fmt"
	"regexp"
)

type Visitor interface {
	VisitString(s *String) error
	VisitRecord(r *Record) error
	VisitInteger(i *Integer) error
	VisitProperties(p Properties) error
}

type Component interface {
	Accept(visitor Visitor) error
}

type Primitive struct {
	Type interface{} `json:"type"`
}

type Properties map[string]Component

type String struct {
	*Primitive
	MaxLength uint   `json:"maxLength"`
	MinLength uint   `json:"minLength"`
	Pattern   string `json:"pattern"`
}

type Record struct {
	MaxProperties uint       `json:"maxProperties"`
	MinProperties uint       `json:"minProperties"`
	Required      []string   `json:"required"`
	Properties    Properties `json:"properties"`
}

type Integer struct{}

type ValidationVisitor struct {
	Instance interface{}
}

func (v *ValidationVisitor) VisitString(s *String) error {
	if s.MaxLength > 0 {
		if !(uint(len(v.Instance.(string))) <= s.MaxLength) {
			return fmt.Errorf("string instance is not valid")
		}
	}
	if !(uint(len(v.Instance.(string))) >= s.MinLength) {
		return fmt.Errorf("string instance is not valid")
	}
	if _, err := regexp.Compile(v.Instance.(string)); err != nil {
		return err
	}
	return nil
}

func (v *ValidationVisitor) VisitRecord(r *Record) error {
	for _, item := range r.Required {
		if properties, ok := v.Instance.(map[string]interface{})["properties"]; ok {
			if _, ok := properties.(map[string]interface{})[item]; !ok {
				return fmt.Errorf("object instance is not valid")
			}
			if err := r.Properties.Accept(v); err != nil {
				return fmt.Errorf("")
			}
		}
	}
	return nil
}

func (v *ValidationVisitor) VisitProperties(p Properties) error {
	for name, member := range v.Instance.(map[string]interface{})["properties"].(map[string]interface{}) {
		if _, ok := p[name]; !ok {
			// For each name that appears in both the instance and
			// as a name within this keyword's value, the child
			// instance for that name successfully validates
			// against the corresponding schema.
			continue
		}
		if err := p[name].Accept(&ValidationVisitor{member}); err != nil {
			return err
		}
	}
	return nil
}

func (v *ValidationVisitor) VisitInteger(i *Integer) error {
	return nil
}

func (s *String) Accept(visitor Visitor) error {
	return visitor.VisitString(s)
}

func (r *Record) Accept(visitor Visitor) error {
	return visitor.VisitRecord(r)
}

func (i *Integer) Accept(visitor Visitor) error {
	return visitor.VisitInteger(i)
}

func (p Properties) Accept(visitor Visitor) error {
	return visitor.VisitProperties(p)
}
