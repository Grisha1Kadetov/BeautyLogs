package analyzer

import (
	"go/ast"
	"go/constant"
	"go/token"
	"regexp"
	"strconv"

	"github.com/Grisha1Kadetov/BeautyLogs/pkg/stringchecker"
	"golang.org/x/tools/go/analysis"
)

func getIdent(expr ast.Expr) *ast.Ident {
	switch v := expr.(type) {
	case *ast.Ident:
		return v
	case *ast.SelectorExpr:
		return v.Sel
	}
	return nil
}

func run(pass *analysis.Pass, config Config) (interface{}, error) {
	config = prepareConfig(config)
	cm := catchMap(config)
	for _, f := range pass.Files {
		ast.Inspect(f, func(n ast.Node) bool {
			call, ok := n.(*ast.CallExpr)
			if !ok {
				return true
			}

			ident := getIdent(call.Fun)
			if ident == nil {
				return true
			}

			obj := pass.TypesInfo.Uses[ident]
			if obj == nil {
				return true
			}

			if obj.Pkg() != nil && cm[obj.Pkg().Path()] != nil &&
				cm[obj.Pkg().Path()][obj.Name()] != nil {
				catch(call, pass, config)
			}

			return true
		})
	}
	return nil, nil
}

func catch(call *ast.CallExpr, pass *analysis.Pass, config Config) {
	if len(call.Args) == 0 {
		return
	}

	first := true
	for _, arg := range call.Args {
		walkStrings(pass, arg, func(start token.Pos, end token.Pos, s string) {
			if first && config.Lowercase {
				first = false
				suggest, ok := stringchecker.CheckFirstLowercase(s)
				if !ok {
					pass.Report(analysis.Diagnostic{
						Pos:     start,
						End:     end,
						Message: "first letter should be lowercase",
						SuggestedFixes: []analysis.SuggestedFix{
							{
								Message: "change to " + strconv.Quote(suggest),
								TextEdits: []analysis.TextEdit{
									{
										Pos:     start,
										End:     end,
										NewText: []byte(strconv.Quote(suggest)),
									},
								},
							},
						},
					})
				}
			}

			var formatRe = regexp.MustCompile(`%(\[[0-9]+\])?[-+0-9.#]*[a-zA-Z]`) //ignore e.g. %v %s
			s = formatRe.ReplaceAllString(s, "")

			if config.OnlyEng {
				if !stringchecker.CheckEnglish(s) {
					pass.Report(analysis.Diagnostic{
						Pos:     start,
						End:     end,
						Message: "logs should contain only English letters",
					})
				}
			}
			if config.Special {
				if !stringchecker.CheckSpecial(s, config.IgnoreSpecial) {
					pass.Report(analysis.Diagnostic{
						Pos:     start,
						End:     end,
						Message: "logs should not contain special characters",
					})
				}
			}
			if config.Sensitive {
				if !stringchecker.CheckSensitive(s, config.SensitiveKeys) {
					pass.Report(analysis.Diagnostic{
						Pos:     start,
						End:     end,
						Message: "logs should not contain sensitive data",
					})
				}
			}
		})
	}
}

func walkStrings(pass *analysis.Pass, expr ast.Expr, onString func(start token.Pos, end token.Pos, s string)) {
	ast.Inspect(expr, func(n ast.Node) bool {
		e, ok := n.(ast.Expr)
		if !ok {
			return true
		}

		if tv, ok := pass.TypesInfo.Types[e]; ok && tv.Value != nil {
			if tv.Type != nil && tv.Type.String() == "string" {
				onString(e.Pos(), e.End(), constant.StringVal(tv.Value))
				return false
			}
		}

		if lit, ok := e.(*ast.BasicLit); ok && lit.Kind == token.STRING {
			if s, err := strconv.Unquote(lit.Value); err == nil {
				onString(lit.Pos(), lit.End(), s)
			}
		}

		return true
	})
}

func prepareConfig(config Config) Config {
	if config.SensitiveKeys == nil {
		config.SensitiveKeys = DefaultConfig.SensitiveKeys
	}
	if config.IgnoreSpecial == nil {
		config.IgnoreSpecial = DefaultConfig.IgnoreSpecial
	}
	if config.Loggers == nil {
		config.Loggers = DefaultConfig.Loggers
	}
	return config
}

func catchMap(config Config) map[string]map[string]any { //map[pkg]map[func]any
	mapped := make(map[string]map[string]any, len(config.Loggers))
	for _, logger := range config.Loggers {
		funcMap := make(map[string]any, len(logger.Funcs))
		for _, fn := range logger.Funcs {
			funcMap[fn] = true
		}
		mapped[logger.Pkg] = funcMap
	}
	return mapped
}
