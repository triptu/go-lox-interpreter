package lox

type loxClass struct {
	name        string
	declaration sClass
	closure     *environment
}

func (c loxClass) String() string {
	return c.name
}
