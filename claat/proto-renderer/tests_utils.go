package genrenderer

import (
	"github.com/googlecodelabs/tools/third_party"
)

// unsupportedType is a dummy type used to showcase the failures of rendering
// non-proto custom types since we take in "any" type as rendering input.
type UnsupportedType struct{}

// Simple proto constructor
func NewDummyProto(in string) *tutorial.StylizedText {
	return &tutorial.StylizedText{
		Text: in,
	}
}
