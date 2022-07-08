package _type

// ECMABoolean :
// Primitive ECMA boolean.
type ECMABoolean struct {
	*ECMAPrimitive
}

// NewPrimitiveBoolean :
// Create a new primitive boolean.
func NewPrimitiveBoolean(value interface{}) *ECMABoolean {
	return &ECMABoolean{
		ECMAPrimitive: NewPrimitive(value),
	}
}

// ToBool :
// Convert this primitive into a boolean value.
// https://developer.mozilla.org/en-US/docs/Glossary/Falsy
// https://developer.mozilla.org/en-US/docs/Glossary/Truthy
func (p *ECMABoolean) ToBool() (bool, error) {
	return p.ECMAPrimitive.GetValue().(bool), nil
}

type ECMABooleanObject struct {
	*ECMAObject
}