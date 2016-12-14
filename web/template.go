package web

import (
	"html/template"
	"log"
	"net/http"

	"github.com/oxtoacart/bpool"
)

var (
	bufpool       = bpool.NewBufferPool(64)
	templateFuncs = template.FuncMap{}
	templateMap   = map[string]*template.Template{}
)

func initTemplates(f *Funnel) {
	templateFuncs = template.FuncMap{
		"reverse":    f.reverse,
		"newcontext": newcontext,
	}

	templateMap = map[string]*template.Template{
		"index":      createTemplate("templates/base.html", "templates/sidebar.html", "templates/index.html"),
		"feed-index": createTemplate("templates/base.html", "templates/sidebar.html", "templates/feed-index.html"),
		"feed-new":   createTemplate("templates/base.html", "templates/sidebar.html", "templates/feed-new.html"),
		"feed-view":  createTemplate("templates/base.html", "templates/sidebar.html", "templates/feed-view.html"),
		"setup":      createTemplate("templates/noauth-base.html", "templates/setup.html"),
		"login":      createTemplate("templates/noauth-base.html", "templates/login.html"),
		"404":        createTemplate("templates/noauth-base.html", "templates/404.html"),
	}
}

func renderTemplate(w http.ResponseWriter, name string, data interface{}) error {
	t, ok := templateMap[name]
	if !ok {
		log.Fatalf("missing template %s", name)
	}
	buf := bufpool.Get()
	defer bufpool.Put(buf)

	// Avoid partially written responses by writing to a buffer
	err := t.ExecuteTemplate(buf, "base", data)
	if err != nil {
		return err
	}

	// TODO(kt): Make this changeable when making the API
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	if _, err := buf.WriteTo(w); err != nil {
		return err
	}
	return nil
}

func createTemplate(files ...string) *template.Template {
	var filebytes = []byte{}
	for _, f := range files {
		filebytes = append(filebytes, MustAsset(f)...)
	}
	tmpl := template.New("*").Funcs(templateFuncs)
	return template.Must(tmpl.Parse(string(filebytes)))
}
