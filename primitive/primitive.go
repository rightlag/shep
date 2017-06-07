package primitive

import (
	"fmt"
	"math"
	"net/url"
	"reflect"
	"regexp"
)

type Visitor interface {
	VisitPrimitive(p *Primitive) error
	VisitString(s *String) error
	VisitRecord(r *Record) error
	VisitInteger(i *Integer) error
	VisitProperties(p Properties) error
	VisitArray(a *Array) error
	VisitAllOf(allOf *AllOf) error
	VisitAnyOf(anyOf *AnyOf) error
	VisitOneOf(oneOf *OneOf) error
	VisitReference(r *Reference) error
}

type Component interface {
	Accept(visitor Visitor) error
}

type Primitive struct {
	Enum        []interface{} `json:"enum"`
	Const       interface{}   `json:"const"`
	Type        interface{}   `json:"type"`
	AllOf       *AllOf        `json:"allOf"`
	AnyOf       *AnyOf        `json:"anyOf"`
	OneOf       *OneOf        `json:"oneOf"`
	Not         Component     `json:"not"`
	Definitions Definitions   `json:"definitions"`
	Title       string        `json:"title"`
	Description string        `json:"description"`
	Default     interface{}   `json:"default"`
	Examples    []interface{} `json:"examples"`
}

type Properties map[string]Component

type Definitions map[string]Component

type AllOf struct {
	elements []Component
}

type AnyOf struct {
	elements []Component
}

type OneOf struct {
	elements []Component
}

type String struct {
	Primitive
	MaxLength uint   `json:"maxLength"`
	MinLength uint   `json:"minLength"`
	Pattern   string `json:"pattern"`
}

func NewString(options ...func(*String)) *String {
	var s String
	for _, option := range options {
		option(&s)
	}
	s.Type = "string"
	return &s
}

type Record struct {
	MaxProperties uint       `json:"maxProperties"`
	MinProperties uint       `json:"minProperties"`
	Required      []string   `json:"required"`
	Properties    Properties `json:"properties"`
}

type Integer struct {
	MultipleOf       float64 `json:"multipleOf"`
	Maximum          float64 `json:"maximum"`
	ExclusiveMaximum bool    `json:"exclusiveMaximum"`
	Minimum          float64 `json:"minimum"`
	ExclusiveMinimum bool    `json:"exclusiveMinimum"`
}

type Array struct {
	Items           interface{} `json:"items"`
	AdditionalItems Component   `json:"additionalItems"`
	MaxItems        uint        `json:"maxItems"`
	MinItems        uint        `json:"minItems"`
	UniqueItems     bool        `json:"uniqueItems"`
	Contains        Component   `json:"contains"`
}

type Reference struct {
	Value string `json:"$ref"`
}

type ValidationVisitor struct {
	Instance interface{}
}

type ReferenceVisitor struct {
	Instance interface{}
}

func (v *ValidationVisitor) VisitPrimitive(p *Primitive) error {
	if p.Enum != nil {
		var found bool
		for _, element := range p.Enum {
			if v.Instance == element {
				found = true
			}
		}
		if !found {
			return fmt.Errorf("instance is not valid")
		}
	}
	if p.Const != nil {
		if reflect.TypeOf(v.Instance) == reflect.TypeOf(p.Const) {
			if !(reflect.ValueOf(v.Instance) == reflect.ValueOf(p.Const)) {
				return fmt.Errorf("instance is not valid")
			}
		}
	}
	if p.Type != nil {
		// types := []string{"null", "boolean", "object", "array", "number", "string", "integer"}
		// TODO: add validation for `type`
	}
	if p.AllOf != nil {
		if err := p.AllOf.Accept(v); err != nil {
			return err
		}
	}
	if p.Not != nil {
		if err := p.Not.Accept(&ValidationVisitor{v.Instance}); err == nil {
			return fmt.Errorf("instance is not valid")
		}
	}
	return nil
}

func (v *ValidationVisitor) VisitString(s *String) error {
	instance, ok := v.Instance.(string)
	if !ok {
		// http://json-schema.org/latest/json-schema-validation.html#rfc.section.4.1
		return nil
	}
	if err := s.Primitive.Accept(v); err != nil {
		return err
	}
	if s.MaxLength > 0 {
		if !(uint(len(instance)) <= s.MaxLength) {
			return fmt.Errorf("string instance is not valid")
		}
	}
	if !(uint(len(instance)) >= s.MinLength) {
		return fmt.Errorf("string instance is not valid")
	}
	if s.Pattern != "" {
		r, err := regexp.Compile(s.Pattern)
		if err != nil {
			return err
		}
		matches := r.FindAllString(instance, -1)
		if len(matches) < 1 {
			return fmt.Errorf("string instance is not valid")
		}
	}
	return nil
}

func (v *ValidationVisitor) VisitRecord(r *Record) error {
	instance := v.Instance.(map[string]interface{})
	if r.MaxProperties > 0 {
		if !(uint(len(instance)) <= r.MaxProperties) {
			return fmt.Errorf("object instance is not valid")
		}
	}
	if !(uint(len(instance)) >= r.MinProperties) {
		return fmt.Errorf("object instance is not valid")
	}
	for _, item := range r.Required {
		if _, ok := instance[item]; !ok {
			return fmt.Errorf("object instance is not valid")
		}
	}
	if err := r.Properties.Accept(v); err != nil {
		return err
	}
	// TODO: `patternProperties`, `additionalProperties`, `dependencies`, `propertyNames`
	return nil
}

func (v *ValidationVisitor) VisitProperties(properties Properties) error {
	instance := v.Instance.(map[string]interface{})
	for name, member := range instance {
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
	instance, ok := v.Instance.(int)
	if !ok {
		return nil
	}
	if i.MultipleOf > 0 {
		remainder := math.Remainder(float64(instance), i.MultipleOf)
		if !(remainder == 0.0) {
			return fmt.Errorf("numeric instance is not valid")
		}
	}
	if !(float64(instance) <= i.Maximum) {
		return fmt.Errorf("numeric instance is not valid")
	}
	if i.ExclusiveMaximum {
		if !(float64(instance) < i.Maximum) {
			return fmt.Errorf("numeric instance is not valid")
		}
	}
	if !(float64(instance) >= i.Minimum) {
		return fmt.Errorf("numeric instance is not valid")
	}
	if i.ExclusiveMinimum {
		if !(float64(instance) > i.Minimum) {
			return fmt.Errorf("numeric instance is not valid")
		}
	}
	return nil
}

func (v *ValidationVisitor) VisitArray(a *Array) error {
	instance := v.Instance.([]interface{})
	switch a.Items.(type) {
	case Component:
		for _, element := range instance {
			if err := a.Items.(Component).Accept(&ValidationVisitor{element}); err != nil {
				return err
			}
		}
	case []Component:
		for i, element := range instance {
			if err := a.Items.([]Component)[i].Accept(&ValidationVisitor{element}); err != nil {
				return err
			}
		}
	}
	if a.MaxItems > 0 {
		if !(uint(len(instance)) <= a.MaxItems) {
			return fmt.Errorf("array instance is not valid")
		}
	}
	if !(uint(len(instance)) >= a.MinItems) {
		return fmt.Errorf("array instance is not valid")
	}
	// TODO: `additionalItems`, `uniqueItems`, `contains`
	return nil
}

func (v *ValidationVisitor) VisitAllOf(allOf *AllOf) error {
	for _, element := range allOf.elements {
		if err := element.Accept(&ValidationVisitor{v.Instance}); err != nil {
			return err
		}
	}
	return nil
}

func (v *ValidationVisitor) VisitAnyOf(anyOf *AnyOf) error {
	// TODO: add validation for `anyOf`
	var count int
	for _, element := range anyOf.elements {
		if err := element.Accept(&ValidationVisitor{v.Instance}); err == nil {
			count++
		}
	}
	return nil
}

func (v *ValidationVisitor) VisitOneOf(oneOf *OneOf) error {
	return nil
}

func (v *ValidationVisitor) VisitReference(r *Reference) error {
	u, err := url.Parse(r.Value)
	if err != nil {
		return err
	}
	r.Value = u.String()
	return nil
}

func (p *Primitive) Accept(visitor Visitor) error {
	return visitor.VisitPrimitive(p)
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

func (allOf *AllOf) Accept(visitor Visitor) error {
	return visitor.VisitAllOf(allOf)
}

func (anyOf *AnyOf) Accept(visitor Visitor) error {
	return visitor.VisitAnyOf(anyOf)
}

func (oneOf *OneOf) Accept(visitor Visitor) error {
	return visitor.VisitOneOf(oneOf)
}

func (r *Reference) Accept(visitor Visitor) error {
	return visitor.VisitReference(r)
}
