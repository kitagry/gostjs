package test

// Test's docs
type Test struct {
	lowercase

	// Name is name
	Name string `json:"name"`

	ID int `json:"id" yaml:"id"`

	Array []string `json:"array"`

	Child Child `json:"child"`

	Child2 *Child
}

// Child is Test's child
type Child struct {
	Name string
}

type lowercase struct {
	Name string
}
