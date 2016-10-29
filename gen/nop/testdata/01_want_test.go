// Code generated by ppgen (github.com/aybabtme/ppgen).
// DO NOT EDIT!

package testdata

import "testing"

func TestNopThing(t *testing.T) {
	tests := []struct {
		name  string
		check func(Thing)
	}{
		{name: "MyFunction", check: func(thing Thing) { thing.MyFunction(nil, "", 0, Composite{}, nil) }},
		{name: "MyFunction2", check: func(thing Thing) { thing.MyFunction2(nil, "", 0, Composite{}, nil) }},
		{name: "MyFunction3", check: func(thing Thing) { thing.MyFunction3(nil, "", 0, Composite{}, nil) }},
		{name: "MyFunction4", check: func(thing Thing) { thing.MyFunction4() }},
		{name: "MyFunction5", check: func(thing Thing) { thing.MyFunction5() }},
		{name: "MyFunction6", check: func(thing Thing) { thing.MyFunction6() }},
		{name: "MyFunction7", check: func(thing Thing) { thing.MyFunction7() }},
		{name: "MyFunction8", check: func(thing Thing) { thing.MyFunction8("") }},
		{name: "MyFunction9", check: func(thing Thing) { thing.MyFunction9("") }},
		{name: "MyFunction10", check: func(thing Thing) { thing.MyFunction10() }},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.check(NopThing())
		})
	}
}
