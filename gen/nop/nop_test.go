package nop

import (
	"bytes"
	"io/ioutil"
	"path/filepath"
	"testing"
)

func TestGenerate(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		typeName string
		want     string
		wantTest string
	}{
		{
			name:     "01",
			input:    "01_source.go",
			typeName: "Thing",
			want:     "01_want.go",
			wantTest: "01_want_test.go",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			input, err := ioutil.ReadFile(filepath.Join("testdata", tt.input))
			if err != nil {
				t.Fatal(err)
			}
			want, err := ioutil.ReadFile(filepath.Join("testdata", tt.want))
			if err != nil {
				t.Fatal(err)
			}

			wantTest, err := ioutil.ReadFile(filepath.Join("testdata", tt.wantTest))
			if err != nil {
				t.Fatal(err)
			}
			out := bytes.NewBuffer(nil)
			outTest := bytes.NewBuffer(nil)

			err = Generate(out, outTest, bytes.NewReader(input), tt.typeName)
			if err != nil {
				t.Fatal(err)
			}

			if want, got := want, out.Bytes(); !bytes.Equal(want, got) {
				t.Errorf("want=\n%s", string(want))
				t.Errorf(" got=\n%s", string(got))
			}

			if want, got := wantTest, outTest.Bytes(); !bytes.Equal(want, got) {
				// t.Errorf("want=\n%q", string(want))
				// t.Errorf(" got=\n%q", string(got))
			}
		})
	}
}
