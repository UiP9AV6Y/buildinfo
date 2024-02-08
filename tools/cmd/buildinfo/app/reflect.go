package app

import (
	"reflect"
	"strings"
)

type lazyReflect struct {
	v interface{}
}

func (r *lazyReflect) String() string {
	v := reflect.TypeOf(r.v).String()
	p := strings.Split(v, ".")

	return strings.ToLower(p[len(p)-1])
}
