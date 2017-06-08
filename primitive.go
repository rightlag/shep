package shep

type Component interface {
	Accept(visitor Visitor) error
}

type Primitive struct {
	Enum        []interface{} `json:"enum"`
	Const       interface{}   `json:"const"`
	Type        interface{}   `json:"type"`
	AllOf       []Component   `json:"allOf"`
	AnyOf       []Component   `json:"anyOf"`
	OneOf       []Component   `json:"oneOf"`
	Not         Component     `json:"not"`
	Definitions Definitions   `json:"definitions"`
	Title       string        `json:"title"`
	Description string        `json:"description"`
	Default     interface{}   `json:"default"`
	Examples    []interface{} `json:"examples"`
}

type Properties map[string]Component

type Definitions map[string]Component

type String struct {
	Primitive
	MaxLength uint   `json:"maxLength"`
	MinLength uint   `json:"minLength"`
	Pattern   string `json:"pattern"`
}

func NewString(options ...func(*String)) *String {
	// Proposed by Rob Pike: https://commandcenter.blogspot.nl/2014/01/self-referential-functions-and-design.html
	var s String
	for _, option := range options {
		option(&s)
	}
	s.Type = "string"
	return &s
}

type Record struct {
	Primitive
	MaxProperties uint       `json:"maxProperties"`
	MinProperties uint       `json:"minProperties"`
	Required      []string   `json:"required"`
	Properties    Properties `json:"properties"`
}

func NewRecord(options ...func(*Record)) *Record {
	var r Record
	for _, option := range options {
		option(&r)
	}
	r.Type = "object"
	return &r
}

type Integer struct {
	Primitive
	MultipleOf       float64 `json:"multipleOf"`
	Maximum          float64 `json:"maximum"`
	ExclusiveMaximum bool    `json:"exclusiveMaximum"`
	Minimum          float64 `json:"minimum"`
	ExclusiveMinimum bool    `json:"exclusiveMinimum"`
}

func NewInteger(options ...func(*Integer)) *Integer {
	var i Integer
	for _, option := range options {
		option(&i)
	}
	i.Type = "integer"
	return &i
}

type Array struct {
	Primitive
	Items           interface{} `json:"items"`
	AdditionalItems Component   `json:"additionalItems"`
	MaxItems        uint        `json:"maxItems"`
	MinItems        uint        `json:"minItems"`
	UniqueItems     bool        `json:"uniqueItems"`
	Contains        Component   `json:"contains"`
}

func NewArray(options ...func(*Array)) *Array {
	var a Array
	for _, option := range options {
		option(&a)
	}
	a.Type = "array"
	return &a
}

type Boolean struct {
	Primitive
}

func NewBoolean(options ...func(*Boolean)) *Boolean {
	var b Boolean
	for _, option := range options {
		option(&b)
	}
	b.Type = "boolean"
	return &b
}

type Null struct {
	Primitive
}

func NewNull(options ...func(*Null)) *Null {
	var n Null
	for _, option := range options {
		option(&n)
	}
	n.Type = "null"
	return &n
}

type Reference struct {
	Value string `json:"$ref"`
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

func (b *Boolean) Accept(visitor Visitor) error {
	return visitor.VisitBoolean(b)
}

func (n *Null) Accept(visitor Visitor) error {
	return visitor.VisitNull(n)
}

func (r *Reference) Accept(visitor Visitor) error {
	return visitor.VisitReference(r)
}
