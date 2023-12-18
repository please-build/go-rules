package generate

import "github.com/bazelbuild/buildtools/build"

type Rule struct {
	name          string
	kind          string
	module        string
	subrepo       string
	srcs          []string
	cgoSrcs       []string
	cSrcs         []string
	compilerFlags []string
	linkerFlags   []string
	pkgConfigs    []string
	asmFiles      []string
	hdrs          []string
	deps          []string
	embedPatterns []string
	isCMD         bool
}

func populateRule(r *build.Rule, targetState *Rule) {
	if len(targetState.cgoSrcs) > 0 {
		r.SetAttr("srcs", NewStringList(targetState.cgoSrcs))
		r.SetAttr("go_srcs", NewStringList(targetState.srcs))
	} else {
		r.SetAttr("srcs", NewStringList(targetState.srcs))
	}
	if len(targetState.cSrcs) > 0 {
		r.SetAttr("c_srcs", NewStringList(targetState.cSrcs))
	}
	if len(targetState.deps) > 0 {
		r.SetAttr("deps", NewStringList(targetState.deps))
	}
	if len(targetState.pkgConfigs) > 0 {
		r.SetAttr("pkg_config", NewStringList(targetState.pkgConfigs))
	}
	if len(targetState.compilerFlags) > 0 {
		r.SetAttr("compiler_flags", NewStringList(targetState.compilerFlags))
	}
	if len(targetState.linkerFlags) > 0 {
		r.SetAttr("linker_flags", NewStringList(targetState.linkerFlags))
	}
	if len(targetState.hdrs) > 0 {
		r.SetAttr("hdrs", NewStringList(targetState.hdrs))
	}
	if len(targetState.asmFiles) > 0 {
		r.SetAttr("asm_srcs", NewStringList(targetState.asmFiles))
	}
	if len(targetState.embedPatterns) > 0 {
		r.SetAttr("resources", &build.CallExpr{
			X: &build.Ident{Name: "glob"},
			List: []build.Expr{
				NewStringList(targetState.embedPatterns),
			},
		})
	}
	if !targetState.isCMD {
		r.SetAttr("_module", NewStringExpr(targetState.module))
		r.SetAttr("_subrepo", NewStringExpr(targetState.subrepo))
	}
}
