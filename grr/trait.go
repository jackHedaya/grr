package grr

type Trait string

func NewTrait(name string) Trait {
	return Trait(name)
}

func (t Trait) String() string {
	return string(t)
}
