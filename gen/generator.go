package gen

import (
	"fmt"
	"go/ast"
	"io"
)

type FileGenerator struct {
	w io.Writer

	pkgName string
	imports []*ast.ImportSpec
}

func NewFile(w io.Writer, pkgName string, imports []*ast.ImportSpec) *FileGenerator {
	return &FileGenerator{
		w: w, pkgName: pkgName, imports: imports,
	}
}

func (g *FileGenerator) WritePrelude() {
	fmt.Fprintf(g.w, "// Code generated by ppgen (github.com/aybabtme/ppgen).\n")
	fmt.Fprintf(g.w, "// DO NOT EDIT!\n")
}

func (g *FileGenerator) WritePkgName() {
	fmt.Fprintf(g.w, "\n")
	fmt.Fprintf(g.w, "package %s\n", g.pkgName)
}

func (g *FileGenerator) WriteImports() {
	fmt.Fprintf(g.w, "\n")

	if len(g.imports) == 1 {
		fmt.Fprintf(g.w, "import ")
	} else {
		fmt.Fprintf(g.w, "import (\n")
	}

	for _, impSpec := range g.imports {

		path := impSpec.Path.Value
		if impSpec.Name != nil {
			alias := impSpec.Name.Name
			fmt.Fprintf(g.w, "\t%s %s\n", alias, path)
		} else {
			fmt.Fprintf(g.w, "\t%s\n", path)
		}
	}

	if len(g.imports) == 1 {
		fmt.Fprintf(g.w, "\n")
	} else {
		fmt.Fprintf(g.w, ")\n")
	}
}

func (g *FileGenerator) WriteFunction(def *FuncDefn, body func(gen *FuncGenerator)) {
	if def.Doc != "" {
		fmt.Fprintf(g.w, "// %s", def.Doc)
	}
	fmt.Fprintf(g.w, "func ")
	if def.ReceiverType != "" {
		fmt.Fprintf(g.w, "(")
		if def.ReceiverName != "" {
			fmt.Fprintf(g.w, "%s ", def.ReceiverName)
		}
		fmt.Fprintf(g.w, "%s", def.ReceiverType)
		fmt.Fprintf(g.w, ") ")
	}

	fmt.Fprintf(g.w, "%s", def.Name)

	fmt.Fprintf(g.w, "(")
	for i, param := range def.Params {
		if i != 0 {
			fmt.Fprintf(g.w, ", ")
		}
		fmt.Fprintf(g.w, "%s %s", param.Name, param.TypeName)
	}
	fmt.Fprintf(g.w, ")")

	if len(def.Results) != 0 {
		fmt.Fprintf(g.w, " ")
		fmt.Fprintf(g.w, "(")
		for i, result := range def.Results {
			if i != 0 {
				fmt.Fprintf(g.w, ", ")
			}
			if result.Name != "" {
				fmt.Fprintf(g.w, "%s %s", result.Name, result.TypeName)
			} else {
				fmt.Fprintf(g.w, "%s", result.TypeName)
			}
		}
		fmt.Fprintf(g.w, ")")

	}
	fmt.Fprintf(g.w, " {")

	if body != nil {
		body(&FuncGenerator{w: g.w, def: def})
	} else {
		fmt.Fprintf(g.w, " return")
		if len(def.Results) != 0 {
			for i, result := range def.Results {
				if i != 0 {
					fmt.Fprintf(g.w, ", ")
				}

				if result.Name != "" {
					fmt.Fprintf(g.w, " %s", result.Name)
				}
			}
		}
		fmt.Fprintf(g.w, " ")
	}
	fmt.Fprintf(g.w, "}\n")
}

func (g *FileGenerator) WriteIfaceImpl(def *IfaceDefn, name string, fields []string, method func(*FuncDefn, *FuncGenerator)) {
	fmt.Fprintf(g.w, "type %s struct{", name)
	for _, field := range fields {
		fmt.Fprintf(g.w, "\t%s\n", field)
	}
	fmt.Fprintf(g.w, "}\n")

	for _, methodDef := range def.Methods {
		methodDef.ReceiverType = name
		g.WriteFunction(methodDef, func(body *FuncGenerator) {
			method(methodDef, body)
		})
	}
}

type FuncGenerator struct {
	w   io.Writer
	def *FuncDefn
}

func (g *FuncGenerator) Inline(code string) {

	fmt.Fprintf(g.w, "\n\t%s\n", code)
}