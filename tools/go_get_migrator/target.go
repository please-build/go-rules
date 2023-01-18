package main

import (
	"github.com/bazelbuild/buildtools/build"
)

type Target struct {
	Kind       string
	Attributes map[string]build.Expr
}

func (t *Target) WithSources(srcs []string) *Target {
	t.Attributes["srcs"] = Expr(srcs)
	return t
}

func (t *Target) WithDeps(deps []string) *Target {
	t.Attributes["deps"] = Expr(deps)
	return t
}

func (t *Target) ToCallExpr() *build.CallExpr {
	rule := build.NewRule(&build.CallExpr{
		X: &build.Ident{Name: t.Kind},
	})
	for k, v := range t.Attributes {
		rule.SetAttr(k, v)
	}
	return rule.Call
}

func NewTarget(kind, name string) *Target {
	return &Target{
		Kind:       kind,
		Attributes: map[string]build.Expr{"name": &build.StringExpr{Value: name}},
	}
}
