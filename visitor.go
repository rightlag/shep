package shep

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
	VisitReference(r *Reference) error
	VisitBoolean(b *Boolean) error
	VisitNull(n *Null) error
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
	for _, element := range p.AllOf {
		if err := element.Accept(v); err != nil {
			return err
		}
	}
	if p.Not != nil {
		if err := p.Not.Accept(v); err == nil {
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
	instance, ok := v.Instance.(map[string]interface{})
	if !ok {
		return nil
	}
	if err := r.Primitive.Accept(v); err != nil {
		return err
	}
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
	if err := i.Primitive.Accept(v); err != nil {
		return err
	}
	if i.MultipleOf > 0 {
		remainder := math.Remainder(float64(instance), i.MultipleOf)
		if !(remainder == 0.0) {
			return fmt.Errorf("numeric instance is not valid")
		}
	}
	if i.Maximum > 0 {
		if !(float64(instance) <= i.Maximum) {
			return fmt.Errorf("numeric instance is not valid")
		}
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
	instance, ok := v.Instance.([]interface{})
	if !ok {
		return nil
	}
	if err := a.Primitive.Accept(v); err != nil {
		return err
	}
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

/*
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
*/

func (v *ValidationVisitor) VisitBoolean(b *Boolean) error {
	_, ok := v.Instance.(bool)
	if !ok {
		return nil
	}
	if err := b.Primitive.Accept(v); err != nil {
		return err
	}
	return nil
}

func (v *ValidationVisitor) VisitNull(n *Null) error {
	if v.Instance != nil {
		return nil
	}
	if err := n.Primitive.Accept(v); err != nil {
		return err
	}
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
