package gen

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"io"
	"path"
	"strconv"
	"strings"
)

func FindIfaceDefn(src io.Reader, typeName string) (pkgName string, imports []*ast.ImportSpec, defn *IfaceDefn, err error) {
	fset := token.NewFileSet()
	f, err := parser.ParseFile(fset, "", src, parser.ParseComments)
	if err != nil {
		return "", nil, nil, fmt.Errorf("can't parse: %v", err)
	}

	var selectors []string
	ast.Inspect(f, func(n ast.Node) bool {

		// not a type definition
		ifaceSpec, ok := n.(*ast.TypeSpec)
		if !ok {
			return true
		}

		// not a definition of the name we want
		if ifaceSpec.Name.Name != typeName {
			return true
		}

		// not an interface
		iface, ok := ifaceSpec.Type.(*ast.InterfaceType)
		if !ok {
			err = fmt.Errorf("type %s is not an interface", typeName)
			return false // this was a user error, we stop
		}
		defn, selectors, err = getIfaceDefn(ifaceSpec, iface)
		if err != nil {
			err = fmt.Errorf("type %s is not a well formed interface: %v", typeName, err)
			return false
		}
		return true
	})

	uniqSelector := make(map[string]struct{})
	for _, selector := range selectors {
		uniqSelector[selector] = struct{}{}
	}

	for _, impSpec := range f.Imports {
		// is this imported with an alias?
		if impSpec.Name != nil {
			if _, ok := uniqSelector[impSpec.Name.Name]; !ok {
				continue // alias not refered to by this interface
			}
		}

		pkgPath, err := strconv.Unquote(impSpec.Path.Value)
		if err != nil {
			return "", nil, nil, err
		}
		selector := path.Base(pkgPath)

		if _, ok := uniqSelector[selector]; !ok {
			continue // not refered to by this interface
		}
		imports = append(imports, impSpec)
	}

	return f.Name.Name, imports, defn, err
}

type IfaceDefn struct {
	Name    string
	Methods []*FuncDefn
}

func getIfaceDefn(ifaceSpec *ast.TypeSpec, ifaceType *ast.InterfaceType) (*IfaceDefn, []string, error) {
	tgt := &IfaceDefn{
		Name: ifaceSpec.Name.Name,
	}
	var imports []string
	for _, m := range ifaceType.Methods.List {
		method, imps, err := getFuncDefn(m)
		if err != nil {
			return nil, nil, err
		}
		method.ReceiverName = strings.ToLower(tgt.Name)
		if len(method.ReceiverName) > 3 {
			method.ReceiverName = method.ReceiverName[:3]
		}
		method.ReceiverType = tgt.Name
		imports = append(imports, imps...)
		tgt.Methods = append(tgt.Methods, method)
	}
	return tgt, imports, nil
}

type FuncDefn struct {
	Doc, Name    string
	ReceiverName string
	ReceiverType string
	Params       []ArgDefn
	Results      []ArgDefn
}

func getFuncDefn(m *ast.Field) (*FuncDefn, []string, error) {
	fnType, ok := m.Type.(*ast.FuncType)
	if !ok {
		return nil, nil, fmt.Errorf("type %s is not a function", m.Names[0].Name)
	}
	method := &FuncDefn{
		Doc:  m.Doc.Text(),
		Name: m.Names[0].Name,
	}
	var selectors []string
	for i, param := range fnType.Params.List {
		arg, selector, err := getArgDefn(param)
		if err != nil {
			return nil, nil, err
		}

		if arg.Name == "" {
			arg.Name = fmt.Sprintf("arg%d", i)
		}
		if len(selector) != 0 {
			selectors = append(selectors, selector...)
		}
		method.Params = append(method.Params, *arg)
	}

	if fnType.Results != nil {
		for i, result := range fnType.Results.List {
			arg, selector, err := getArgDefn(result)
			if err != nil {
				return nil, nil, err
			}
			if arg.Name == "" {
				arg.Name = fmt.Sprintf("out%d", i)
			}
			if len(selector) != 0 {
				selectors = append(selectors, selector...)
			}
			method.Results = append(method.Results, *arg)
		}
	}
	return method, selectors, nil
}

type ArgDefn struct {
	Name     string
	TypeName string
	Nullable bool
}

func getArgDefn(field *ast.Field) (*ArgDefn, []string, error) {
	var name string
	if len(field.Names) == 1 {
		name = field.Names[0].Name
	}
	tn, selector, nullable, err := getTypeName(field.Type)
	if err != nil {
		return nil, selector, err
	}
	return &ArgDefn{
		Name:     name,
		TypeName: tn,
		Nullable: nullable,
	}, selector, nil
}

func getTypeName(expr ast.Expr) (name string, selector []string, nullable bool, err error) {
	switch xt := expr.(type) {
	case *ast.SelectorExpr:
		pkgName, pkgs, nullable, err := getTypeName(xt.X)
		return pkgName + "." + xt.Sel.Name, append(pkgs, pkgName), nullable, err
	case *ast.Ident:
		return xt.Name, nil, false, nil
	case *ast.StarExpr:
		name, pkg, _, err := getTypeName(xt.X)
		return "*" + name, pkg, true, err
	case *ast.ArrayType:
		name, pkg, _, err := getTypeName(xt.Elt)
		return "[]" + name, pkg, true, err
	case *ast.MapType:
		keyName, keyPkg, _, err := getTypeName(xt.Key)
		valName, valPkg, _, err := getTypeName(xt.Value)
		return "map[" + keyName + "]" + valName, append(keyPkg, valPkg...), true, err
	default:
		return "", nil, false, fmt.Errorf("not a valid ast expression: %T", expr)
	}
}
