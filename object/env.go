package object

func NewEnclosedEnv(outer *Env) *Env {
	env := NewEnv()
	env.outer = outer
	return env
}

func NewEnv() *Env {
	return &Env{
		store: make(map[string]Representation),
	}
}

type Env struct {
	outer *Env
	store map[string]Representation
}

func (e *Env) Get(variable string) (rep Representation, has bool) {
	rep, has = e.store[variable]

	if !has && e.outer != nil {
		rep, has = e.outer.Get(variable)
	}

	return rep, has
}

func (e *Env) Set(name string, value Representation) {
	e.store[name] = value
}
