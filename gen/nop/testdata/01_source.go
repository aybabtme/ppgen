package testdata

import (
	"context"

	oldctx "golang.org/x/net/context"
)

type Composite struct {
	A string
	B int
}

type Thing interface {
	// MyFunction does shits.
	MyFunction(context.Context, string, int, Composite, *Composite) (string, int, Composite, *Composite, oldctx.Context, error)

	// MyFunction2 does shits.
	MyFunction2(ctx context.Context, str string, i int, comp Composite, comptr *Composite) (string, int, Composite, *Composite, oldctx.Context, error)

	// MyFunction3 does shits.
	MyFunction3(context.Context, string, int, Composite, *Composite) (str string, i int, comp Composite, comptr *Composite, ctx oldctx.Context, err error)

	// MyFunction4 does shits.
	MyFunction4()

	MyFunction5()
	MyFunction6()
	MyFunction7()

	MyFunction8(string) error
	MyFunction9(string)
	MyFunction10() error

	MyFunction11(a, b string) (c, d string, err error)
}
