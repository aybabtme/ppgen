package interceptor

import (
	"bytes"
	"fmt"
	"go/ast"
	"go/format"
	"io"
	"log"
	"strings"

	"github.com/aybabtme/ppgen/gen"
)

func Generate(dst, testDst io.Writer, src io.Reader, typeName string) error {
	pkgName, imports, iface, err := gen.FindIfaceDefn(src, typeName)
	if err != nil {
		return err
	}
	genTypeName := "nop" + typeName
	if err := generateImpl(dst, typeName, genTypeName, pkgName, imports, iface); err != nil {
		return err
	}

	if err := generateImplTest(testDst, typeName, genTypeName, pkgName, imports, iface); err != nil {
		return err
	}

	return err
}

func generateImpl(dst io.Writer, typeName, genTypeName, pkgName string, imports []*ast.ImportSpec, iface *gen.IfaceDefn) error {

	buf := bytes.NewBuffer(nil)
	g := gen.NewFile(buf, pkgName, imports)

	g.WritePrelude()
	g.WritePkgName()
	g.WriteImports()

	funcName := strings.Title(genTypeName)
	g.WriteFunction(
		&gen.FuncDefn{
			Doc:     fmt.Sprintf("%s returns a %s that does nothing.\n", funcName, typeName),
			Name:    funcName,
			Results: []gen.ArgDefn{{TypeName: typeName}},
		},
		func(body *gen.FuncGenerator) {
			body.Return([]gen.ArgDefn{
				{Name: fmt.Sprintf("%s{}", genTypeName)},
			})
		},
	)

	// make the impl. ignore it's arguments and return arbitrary result values
	for _, metDefn := range iface.Methods {
		metDefn.ReceiverName = ""
		for i := range metDefn.Params {
			metDefn.Params[i].Name = "_"
		}
	}

	g.WriteIfaceImpl(
		iface,
		genTypeName,
		nil,
		func(defn *gen.FuncDefn, body *gen.FuncGenerator) {
			if len(defn.Results) != 0 {
				body.Return(defn.Results)
			}
		},
	)

	out, err := format.Source(buf.Bytes())
	if err != nil {
		return err
	}

	_, err = dst.Write(out)
	return err
}

func generateImplTest(dst io.Writer, typeName, genTypeName, pkgName string, imports []*ast.ImportSpec, iface *gen.IfaceDefn) error {

	buf := bytes.NewBuffer(nil)
	g := gen.NewFile(buf, pkgName, imports)

	g.WritePrelude()
	g.WritePkgName()
	g.WriteImports("testing")

	g.WriteEmptyLine()

	funcName := "Test" + strings.Title(genTypeName)
	g.WriteFunction(
		&gen.FuncDefn{
			Name: funcName,
			Params: []gen.ArgDefn{{
				Name:     "t",
				TypeName: "*testing.T",
			}},
		},
		func(body *gen.FuncGenerator) {
			body.Inline(`tests := []struct{`)
			body.Inline(`	name  string`)
			body.Inline(`	check func(%s)`, typeName)
			body.Inline(`}{`)
			for _, fn := range iface.Methods {
				fnName := fn.Name
				recv := strings.ToLower(typeName)
				args := strings.Join(defaultNulArgs(fn.Params), ", ")
				body.Inline(`	{name: %q, check: func(%s %s) { %s.%s(%s) }},`,
					fnName,
					recv, typeName,
					recv, fnName, args,
				)
			}
			body.Inline(`}`)
			funcName := strings.Title(genTypeName)
			body.Inline(`for _, tt := range tests {`)
			body.Inline(`	t.Run(tt.name, func(t *testing.T) {`)
			body.Inline(`		tt.check(%s())`, funcName)
			body.Inline(`	})`)
			body.Inline(`}`)
		},
	)

	out, err := format.Source(buf.Bytes())
	if err != nil {
		log.Print(err)
		log.Fatal(buf.String())
		return err
	}

	_, err = dst.Write(out)
	return err
}

func defaultNulArgs(params []gen.ArgDefn) []string {
	var out []string
	for _, p := range params {
		out = append(out, nulForTypename(p))
	}
	return out
}

func nulForTypename(arg gen.ArgDefn) string {
	if arg.Nullable {
		return "nil"
	}
	typeName := arg.TypeName
	if strings.HasPrefix(typeName, "[") {
		return "nil"
	}
	if strings.HasPrefix(typeName, "map[") {
		return "nil"
	}
	if strings.HasPrefix(typeName, "chan") || strings.HasPrefix(typeName, "<-") {
		return "nil"
	}
	switch typeName {
	case "int", "int8", "int16", "int32", "int64",
		"uint", "uint8", "uint16", "uint32", "uint64":
		return "0"
	case "float32", "float64":
		return "0.0"
	case "string":
		return `""`
	case "byte":
		return `byte(0)`
	case "rune":
		return `rune(0)`
	case "bool":
		return `false`
	default:
		return "func() (out " + typeName + ") { return out }()"
	}
}
