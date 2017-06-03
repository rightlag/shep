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
	VisitArray(a *Array) error
}

type Component interface {
	Accept(visitor Visitor) error
}

type Primitive struct {
	Type interface{} `json:"type"`
}

type Properties map[string]Component

type String struct {
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

type Array struct {
	Items           interface{} `json:"items"`
	AdditionalItems Component   `json:"additionalItems"`
	MaxItems        uint        `json:"maxItems"`
	MinItems        uint        `json:"minItems"`
	UniqueItems     bool        `json:"uniqueItems"`
	Contains        Component   `json:"contains"`
}

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
			if r.MaxProperties > 0 {
				if !(uint(len(properties.(map[string]interface{}))) <= r.MaxProperties) {
					return fmt.Errorf("object instance is not valid")
				}
			}
			if !(uint(len(properties.(map[string]interface{}))) >= r.MinProperties) {
				return fmt.Errorf("object instance is not valid")
			}
			if _, ok := properties.(map[string]interface{})[item]; !ok {
				return fmt.Errorf("object instance is not valid")
			}
			if err := r.Properties.Accept(v); err != nil {
				return fmt.Errorf("")
			}
		}
	}
	// TODO: `patternProperties`, `additionalProperties`, `dependencies`, `propertyNames`
	return nil
}

func (v *ValidationVisitor) VisitProperties(properties Properties) error {
	for name, member := range v.Instance.(map[string]interface{})["properties"].(map[string]interface{}) {
		if _, ok := properties[name]; !ok {
			// For each name that appears in both the instance and
			// as a name within this keyword's value, the child
			// instance for that name successfully validates
			// against the corresponding schema.
			continue
		}
		if err := properties[name].Accept(&ValidationVisitor{member}); err != nil {
			return err
		}
	}
	return nil
}

func (v *ValidationVisitor) VisitInteger(i *Integer) error {
	return nil
}

func (v *ValidationVisitor) VisitArray(a *Array) error {
	switch a.Items.(type) {
	case Component:
		for _, element := range v.Instance.([]interface{}) {
			if err := a.Items.(Component).Accept(&ValidationVisitor{element}); err != nil {
				return err
			}
		}
	case []Component:
		for i, element := range v.Instance.([]interface{}) {
			if err := a.Items.([]Component)[i].Accept(&ValidationVisitor{element}); err != nil {
				return err
			}
		}
	}
	if a.MaxItems > 0 {
		if !(uint(len(v.Instance.([]interface{}))) <= a.MaxItems) {
			return fmt.Errorf("array instance is not valid")
		}
	}
	if !(uint(len(v.Instance.([]interface{}))) >= a.MinItems) {
		return fmt.Errorf("array instance is not valid")
	}
	// TODO: `additionalItems`, `uniqueItems`, `contains`
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

func (a *Array) Accept(visitor Visitor) error {
	return visitor.VisitArray(a)
}
