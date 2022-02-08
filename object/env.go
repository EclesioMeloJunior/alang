package object

func NewEnv() *Env {
	return &Env{
		store: make(map[string]Representation),
	}
}

type Env struct {
	store map[string]Representation
}

func (e *Env) Get(variable string) (rep Representation, has bool) {
	storedValue, has := e.store[variable]
	return storedValue, has
}

func (e *Env) Set(name string, value Representation) {
	e.store[name] = value
}
