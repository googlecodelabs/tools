package devrel_tutorial

// TODO: Same generic TODO as for oneofs +
//       Once that's figured out, name all repeated fields 'children'

func (el *TextBlock) GetInnerContent() []ProtoRenderer {
	desambiguatedContents := make([]ProtoRenderer, len(el.Children))
	for i, c := range el.Children {
		// TODO: Parrallel/concurrent
		// 			 Perfomance change: O(h) vs O(n)
		desambiguatedContents[i] = c.GetInnerContent()
	}
	return desambiguatedContents
	// renderedContent
}
