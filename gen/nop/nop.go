package nop

import (
	"bytes"
	"fmt"
	"go/format"
	"io"
	"strings"

	"github.com/aybabtme/ppgen/gen"
)

func Generate(dst, testDst io.Writer, src io.Reader, typeName string) error {
	pkgName, imports, iface, err := gen.FindIfaceDefn(src, typeName)
	if err != nil {
		return err
	}

	buf := bytes.NewBuffer(nil)
	g := gen.NewFile(buf, pkgName, imports)

	g.WritePrelude()
	g.WritePkgName()
	g.WriteImports()

	genTypeName := "nop" + typeName

	funcName := strings.Title(genTypeName)
	g.WriteFunction(
		&gen.FuncDefn{
			Doc:     fmt.Sprintf("%s returns a %s that does nothing.\n", funcName, typeName),
			Name:    funcName,
			Results: []gen.ArgDefn{{TypeName: typeName}},
		},
		func(body *gen.FuncGenerator) {
			body.Inline(
				fmt.Sprintf("return %s{}", genTypeName),
			)
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
				body.Inline("return")
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
