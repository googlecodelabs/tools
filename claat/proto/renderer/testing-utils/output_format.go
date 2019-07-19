package testingutils

import (
	"github.com/googlecodelabs/tools/third_party"
)

// UnsupportedType is a dummy type used to showcase the failures of rendering
// non-proto custom types since we take in "any" type as rendering input.
type UnsupportedType struct{}

// NewDummyProto is a simple proto constructor
func NewDummyProto() *tutorial.StylizedText {
	return &tutorial.StylizedText{
		Text: "dummy",
	}
}
