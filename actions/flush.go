package actions

import (
	"bytes"
	"fmt"
	"html/template"
	"io"
	"net/http"

	"github.com/kyoto-framework/kyoto"
)

// Flush is a function to immediately update a component layout during action.
// It is useful for updating a component layout multiple times during long action.
func Flush(b *kyoto.Core) {
	// Extract context
	rw := b.Context.GetResponseWriter()
	rwf := rw.(http.Flusher)
	// Render
	buffer := bytes.NewBufferString("")
	var err error
	if b.State.Get("internal:render:writer") != nil {
		renderer := b.State.Get("internal:render:writer").(func(io.Writer) error)
		err = renderer(buffer)
	} else {
		tmpl := b.Context.Get("internal:render:template").(*template.Template)
		tmplclone, _ := tmpl.Clone()
		err = tmplclone.Execute(buffer, b.State.Export())
	}
	if err != nil {
		rw.WriteHeader(500)
		panic(err)
	}
	html := buffer.String()
	// Write to stream
	_, err = fmt.Fprint(rw, html)
	if err != nil {
		panic(err)
	}
	// Flush
	rwf.Flush()
}
