package genrenderer

// TODO: delete these after next PR
// SampleProtoTemplate serves as a dummy supported proto before proto redefinitions
type SampleProtoTemplate struct {
	Value interface{}
}

// NewSampleProtoTemplate pointer constructor for the above
func NewSampleProtoTemplate(el interface{}) *SampleProtoTemplate {
	return &SampleProtoTemplate{Value: el}
}
