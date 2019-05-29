package codelab_renderer



func Render(el interface{}) string {
	// Derefence: proto struct addr (*devrel_tutorial.Type) => proto struct (Type)
	// This is backtracts that to get a "*"less-name for template mappings (Type)
	originalType := reflect.TypeOf(el).Elem()
	// might need to do other "magic" stuff here...
	// fmt.Printf("%+v\n", el)
	// fmt.Printf("%+v\n", originalType)
	// can the below be done concurrently?? hmm...
	// required the protos to have the "Text" field, right?
	// renderRepeated
	// if render_base.IsComposite(typ) {
	// 	tmplData := NewCompositeData()
	// // renderOneOf
	// } else 
	if render_base.IsOneof(originalType) {
		tmplData := &el
		// => access the 'Text'     => struct
		// => access the 'Content'  => interface
	} else {
		tmplData := &el
	}

	// ce := &render_base.CodelabElement{&el}
	// return render_base.ExecuteTemplate(&originalType, tmplName, t)
	// need for generic:
	// Content.(reflect.TypeOf(devrel_tutorial.TextBlock{}))
	//
	// else
	// c, isComposite := originalType.FieldByName("Children")
	// if isComposite {
	// 	fmt.Println(originalType, c)
	tmplName := originalType.Name()
	// pass thru a safehtml filter... etc!
	return render_base.ExecuteTemplate(tmplData, tmplName, t)
}

// func renderOneof(el *render_base.CodelabElement) string {
// 	// mapping.. etc
// 	el.Content
// 	// switch mapping... ugh
// 	// for i, c := range el.Children {
// 	// 	// TODO: Parrallel/concurrent
// 	// 	// 			 Perfomance change: O(h) vs O(n)
// 	// 	rendered_content[i] = Render(c)
// 	// }
// 	return strings.Join(rendered_content, "&#10;")
// }