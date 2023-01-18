package main

import (
	"fmt"

	"github.com/bazelbuild/buildtools/build"
)

func Expr(val interface{}) build.Expr {
	switch val := val.(type) {
	case string:
		return &build.StringExpr{Value: val}
	case []interface{}:
		ret := make([]build.Expr, len(val))
		for i, v := range val {
			ret[i] = Expr(v)
		}
		return &build.ListExpr{List: ret}
	case map[interface{}]interface{}:
		ret := make([]*build.KeyValueExpr, 0, len(val))
		for k, v := range val {
			ret = append(ret, &build.KeyValueExpr{Key: Expr(k), Value: Expr(v)})
		}
		return &build.DictExpr{List: ret}
	default:
		panic(fmt.Errorf("unrecognised type: %T", val))
	}
}
