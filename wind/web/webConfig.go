package web

import (
	"reflect"

	"github.com/Luncert/OnGo/wind/core"
)

// Annotations

// GetMapping ...
var GetMapping = core.Annotation("GetMapping")

func getMappingHandler(name string, _type reflect.Type, value reflect.Value) error {
	return nil
}
