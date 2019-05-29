package devrel_tutorial

// TODO: How to make the following cleaner/more generic?
//       { [typString] = func(x) { return x.[Type] } ?
//       using: https://github.com/golang/protobuf/issues/261#issuecomment-430496210

func (el *InlineContent) GetInnerContent() ProtoRenderer {
	switch x := el.Content.(type) {
	case *InlineContent_Text:
		return x.Text
	}
	return nil
}

func (el *BlockContent) GetInnerContent() ProtoRenderer {
	switch x := el.Content.(type) {
	case *BlockContent_Text:
		return x.Text
	}
	return nil
}
