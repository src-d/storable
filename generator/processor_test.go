package generator

import (
	"go/ast"
	"go/parser"
	"go/token"
	"reflect"
	"testing"

	"golang.org/x/tools/go/types"
	. "gopkg.in/check.v1"
)

func Test(t *testing.T) { TestingT(t) }

type ProcessorSuite struct{}

var _ = Suite(&ProcessorSuite{})

func (s *ProcessorSuite) TestTags(c *C) {
	fixtureSrc := `
	package fixture

	import 	"gopkg.in/tyba/storable.v1"

	type Foo struct {
		storable.Document
		Int int "foo"
	}
	`

	pkg := s.processFixture(fixtureSrc)
	c.Assert(pkg.Models[0].Fields[1].Tag, Equals, reflect.StructTag("foo"))
}

func (s *ProcessorSuite) TestRecursiveStruct(c *C) {
	fixtureSrc := `
	package fixture

	import 	"gopkg.in/tyba/storable.v1"

	type Recur struct {
		storable.Document
		Foo string
		R *Recur
	}
	`

	pkg := s.processFixture(fixtureSrc)

	c.Assert(
		pkg.Models[0].Fields[2].Fields[2].CheckedNode,
		Equals,
		pkg.Models[0].Fields[2].CheckedNode,
		Commentf("direct type recursivity not handled correctly."),
	)

	c.Assert(len(pkg.Models[0].Fields[2].Fields[2].Fields), Equals, 0)
}

func (s *ProcessorSuite) TestDeepRecursiveStruct(c *C) {
	fixtureSrc := `
	package fixture

	import 	"gopkg.in/tyba/storable.v1"

	type Recur struct {
		storable.Document
		Foo string
		R *Other
	}

	type Other struct {
		R *Recur
	}
	`

	pkg := s.processFixture(fixtureSrc)

	c.Assert(pkg.Models[0].Fields[2].Fields[0].Fields[2].CheckedNode, Equals, pkg.Models[0].Fields[2].CheckedNode, Commentf("direct type recursivity not handled correctly."))
	c.Assert(len(pkg.Models[0].Fields[2].Fields[0].Fields[2].Fields), Equals, 0)
}

func (s *ProcessorSuite) processFixture(source string) *Package {
	fset := &token.FileSet{}
	astFile, err := parser.ParseFile(fset, "fixture.go", source, 0)
	if err != nil {
		panic(err)
	}

	cfg := &types.Config{}
	p, _ := cfg.Check("foo", fset, []*ast.File{astFile}, nil)
	if err != nil {
		panic(err)
	}

	prc := NewProcessor("fixture", nil)
	prc.TypesPkg = p
	pkg, err := prc.processTypesPkg()
	if err != nil {
		panic(err)
	}

	return pkg
}
