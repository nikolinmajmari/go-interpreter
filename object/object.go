package object

import (
	"go/types"
	"interpreter/ast"
)

const FUNCTION_OBJ = "FUNCTION"

func NewEnvironment() *Environment {
	s := make(map[string]types.Object)
	return &Environment{store: s}
}

type Environment struct {
	store map[string]types.Object
}

func (e *Environment) Get(key string) (types.Object, bool) {
	obj, ok := e.store[key]
	return obj, ok
}

func (e *Environment) Set(key string, value types.Object) types.Object {
	e.store[key] = value
	return value
}

type Function struct {
	Parameters  []*ast.Identifier
	Body        *ast.BlockStatement
	Environment Environment
}
