package main

import (
	"bytes"
	"io"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/aybabtme/ppgen/gen/nop"
	"github.com/fatih/camelcase"
	"github.com/urfave/cli"
)

const appname = "ppgen"

func main() {
	log.SetPrefix(appname + ": ")
	log.SetFlags(0)
	var (
		srcFilename string
		typeName    string

		dstPrefix   = ""
		fileModeStr = "0644"

		flags = []cli.Flag{
			cli.StringFlag{Name: "src", Destination: &srcFilename, Usage: "`file` containing the type declaration."},
			cli.StringFlag{Name: "type", Destination: &typeName, Usage: "`name` of the type to generate code for."},
		}
		globalFlags = append(
			flags,
			cli.StringFlag{Name: "dst.prefix", Value: dstPrefix, Destination: &dstPrefix, Usage: "`prefix` to use for the generated files."},
			cli.StringFlag{Name: "file.mode", Value: fileModeStr, Destination: &fileModeStr, Usage: "`mode` to use for the generated files."},
		)
	)

	app := cli.NewApp()
	app.Name = appname
	app.Usage = "Pew Pew generator - wtf is this name?"
	app.Description = "Generates wrapper stuff you're tired of typing by hand."
	app.Flags = globalFlags
	app.Commands = []cli.Command{
		{
			Name:  "nop",
			Usage: "generates a no-op implementation of the type",
			Flags: flags,
			Action: func(ctx *cli.Context) {
				cmdCommon(nop.Generate, "nop", dstPrefix, srcFilename, typeName, mustFilemode(fileModeStr))
			},
		},
	}
	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}

func mustFilemode(fileModeStr string) uint32 {
	u, err := strconv.ParseUint(fileModeStr, 8, 32)
	if err != nil {
		log.Fatalf("invalid filemode: %v", err)
	}
	return uint32(u)
}

func dstFilenames(srcFilename, typeName, suffix, prefix string) (srcName, testName string) {
	dir := filepath.Dir(srcFilename)

	var parts []string
	if prefix != "" {
		parts = append(parts, prefix)
	}
	parts = append(parts, suffix)
	parts = append(parts, camelcase.Split(typeName)...)
	base := strings.ToLower(strings.Join(parts, "_"))

	srcName = filepath.Join(dir, base+".go")
	testName = filepath.Join(dir, base+"_test.go")
	return srcName, testName
}

func cmdCommon(
	gen func(dst, dsttest io.Writer, src io.Reader, typename string) error,
	suffix, dstPrefix, srcFilename, typeName string, filemode uint32,
) {
	src, err := os.Open(srcFilename)
	if err != nil {
		log.Fatalf("can't open src, %v", err)
	}
	defer src.Close()

	dstFile, dstTestfile := dstFilenames(srcFilename, typeName, suffix, dstPrefix)

	dstBuf := bytes.NewBuffer(nil)
	dstTestBuf := bytes.NewBuffer(nil)

	if err := gen(dstBuf, dstTestBuf, src, typeName); err != nil {
		log.Fatalf("can't generate implementation: %v", err)
	}

	log.Printf("generated implementation in %q", dstFile)

	err = ioutil.WriteFile(dstFile, dstBuf.Bytes(), os.FileMode(filemode))
	if err != nil {
		log.Fatalf("can't create dst file, %v", err)
	}

	log.Printf("generated tests in %q", dstTestfile)
	err = ioutil.WriteFile(dstTestfile, dstTestBuf.Bytes(), os.FileMode(filemode))
	if err != nil {
		log.Fatalf("can't create dst test file, %v", err)
	}

}
