package analyzer

import (
	"go/ast"
	"go/token"
	"go/types"
	"log"
	"regexp"
	"strings"

	"golang.org/x/tools/go/analysis"
	"golang.org/x/tools/go/analysis/passes/inspect"
)

var (
	// Analyzer defines the analyzer for closecheck
	Analyzer = &analysis.Analyzer{
		Name:      "verifyiface",
		Doc:       "check that a interface implementation is verified",
		Run:       run,
		Requires:  []*analysis.Analyzer{inspect.Analyzer},
		FactTypes: []analysis.Fact{new(ifaceVerifier)},
	}

	externalFilesRx = regexp.MustCompile(`(/pkg/mod/|/libexec/|/go-build/)`)

	verbose bool
)

type impl struct {
	st        *types.Struct
	ifaceData *ifaceData
	verified  bool
	stPos     token.Pos
	stName    string
}

type ifaceData struct {
	name  string
	iface *types.Interface
	pos   token.Pos
}

type ifaceVerifier struct {
}

func (c *ifaceVerifier) AFact() {}

func run(pass *analysis.Pass) (interface{}, error) {
	ifaces := findAllInterfaces(pass)
	implementations := findAllImplementations(pass, ifaces)

	setInstantiations(pass, ifaces, implementations)

	for _, impl := range implementations {
		if !impl.verified {
			pass.Reportf(impl.stPos, "struct %s doesn't verify interface compliance for %s", impl.stName, impl.ifaceData.name)
		}
	}

	return nil, nil
}

func findAllInterfaces(pass *analysis.Pass) []*ifaceData {
	ifaces := []*ifaceData{}

	for _, file := range pass.Files {
		tfile := pass.Fset.File(file.Pos())

		printV("## looking for interfaces in %v", tfile.Name())

		ast.Inspect(file, func(n ast.Node) bool {
			switch nt := n.(type) {
			case *ast.File, *ast.GenDecl:
				return true
			case *ast.TypeSpec:
				// debug(pass, n)

				ifaceType, ok := nt.Type.(*ast.InterfaceType)
				if !ok {
					return false
				}

				iface := pass.TypesInfo.TypeOf(ifaceType).(*types.Interface)

				// ignore empty interfaces
				if len(ifaceType.Methods.List) > 0 {
					printV("found interface %v", nt.Name.Name)

					ifaces = append(ifaces, &ifaceData{
						name:  nt.Name.Name,
						iface: iface,
						pos:   nt.Pos(),
					})
				}

				return false
			default:
				return false
			}
		})
	}

	return ifaces
}

func findAllImplementations(pass *analysis.Pass, ifaces []*ifaceData) []*impl {
	implementations := []*impl{}
	visited := map[token.Pos]bool{}

	for _, file := range pass.Files {
		tfile := pass.Fset.File(file.Pos())
		if externalFilesRx.MatchString(tfile.Name()) {
			// ignore external files
			continue
		}

		printV("## looking for structs in %v", tfile.Name())

		ast.Inspect(file, func(n ast.Node) bool {
			switch nt := n.(type) {
			case *ast.File:
				return true
			case *ast.TypeSpec:
				if strings.Contains(nt.Doc.Text(), "#noverifyiface") {
					return false
				}

				return true
			case *ast.GenDecl:
				if strings.Contains(nt.Doc.Text(), "#noverifyiface") {
					return false
				}

				return true
			case *ast.Ident:
				t := pass.TypesInfo.ObjectOf(nt)
				if t == nil || t.Type() == nil || t.Type().Underlying() == nil || nt.Obj == nil {
					return false
				}

				if nt.Obj.Kind != ast.Typ {
					return false
				}

				if visited[t.Pos()] {
					return false
				}
				visited[t.Pos()] = true

				switch tt := t.Type().Underlying().(type) {
				case *types.Struct:
					printV("verifying struct: %s (%T) %v (%T)", t.Name(), t, tt, tt)

					for _, iface := range ifaces {
						implemented := types.Implements(t.Type(), iface.iface) || types.Implements(types.NewPointer(t.Type()), iface.iface)
						printV("%v implements %v? %v", t.Name(), iface.name, implemented)

						if implemented {
							implementations = append(implementations, &impl{t.Type().Underlying().(*types.Struct), iface, false, t.Pos(), t.Name()})
						}
					}

					return false
				default:
					_ = tt

					return false
				}
			}

			return false
		})
	}

	return implementations
}

func setInstantiations(pass *analysis.Pass, ifaces []*ifaceData, implementations []*impl) {
	for _, file := range pass.Files {
		ast.Inspect(file, func(n ast.Node) bool {
			switch nt := n.(type) {
			case *ast.File:
				return true
			case *ast.GenDecl:
				if nt.Tok != token.VAR {
					return false
				}

				if len(nt.Specs) == 0 {
					return false
				}

				for _, spec := range nt.Specs {
					valueSpec, ok := spec.(*ast.ValueSpec)
					if !ok {
						continue
					}

					id := valueSpec.Names[0]
					if id.Name != "_" {
						continue
					}

					typ, ok := valueSpec.Type.(*ast.Ident)
					if !ok {
						continue
					}

					t := pass.TypesInfo.ObjectOf(typ)
					iface, ok := t.Type().Underlying().(*types.Interface)
					if !ok {
						continue
					}

					ifaceFound := false
					for _, v := range ifaces {
						if v.iface == iface {
							ifaceFound = true
							break
						}
					}

					if !ifaceFound {
						continue
					}

					value := valueSpec.Values[0]
					switch vv := value.(type) {
					case *ast.CompositeLit:
						// checks: var _ Iface = Foo{}

						st := getStructFromCompositeLit(pass, vv)
						for _, v := range implementations {
							if st == v.st && v.ifaceData.iface == iface {
								v.verified = true
							}
						}
					case *ast.UnaryExpr:
						// checks: var _ Iface = &Foo{}

						st := getStructFromCompositeLit(pass, vv.X)
						for _, v := range implementations {
							if st == v.st && v.ifaceData.iface == iface {
								v.verified = true
							}
						}

					case *ast.CallExpr:
						// checks: var _ Iface = (*Foo)(nil)
						if len(vv.Args) == 0 {
							continue
						}

						id, ok := vv.Args[0].(*ast.Ident)
						if !ok {
							continue
						}

						if id.Name != "nil" {
							continue
						}

						p, ok := vv.Fun.(*ast.ParenExpr)
						if !ok {
							continue
						}

						s, ok := p.X.(*ast.StarExpr)
						if !ok {
							continue
						}

						id, ok = s.X.(*ast.Ident)
						if !ok {
							continue
						}

						t := pass.TypesInfo.ObjectOf(id)
						st, ok := t.Type().Underlying().(*types.Struct)
						if !ok {
							continue
						}

						for _, v := range implementations {
							if st == v.st && v.ifaceData.iface == iface {
								v.verified = true
							}
						}
					default:
						// log.Printf("unhandled: %T", vv)
					}
				}

				return true
			case *ast.Ident:
				_ = nt
			}

			return false
		})
	}
}

func getStructFromCompositeLit(pass *analysis.Pass, expr ast.Expr) *types.Struct {
	cl, ok := expr.(*ast.CompositeLit)
	if !ok {
		return nil
	}

	typ, ok := cl.Type.(*ast.Ident)
	if !ok {
		return nil
	}

	t := pass.TypesInfo.ObjectOf(typ)
	st, ok := t.Type().Underlying().(*types.Struct)
	if !ok {
		return nil
	}

	return st
}

func printV(format string, args ...interface{}) {
	if !verbose {
		return
	}

	log.Printf(format, args...)
}

var _ = printV

func debug(pass *analysis.Pass, x interface{}) {
	log.Printf("n: %T", x)
	ast.Print(pass.Fset, x)
}

var _ = debug

func init() {
	Analyzer.Flags.BoolVar(&verbose, "verbose", false, "enable verbose mode")
}
