package main

import (
	"embed"
	"fmt"
	"html/template"
	"net/http"
)

type templateData struct {
	StringMap            map[string]string
	IntMap               map[string]int
	FloatMap             map[string]float32
	Data                 map[string]interface{}
	CSRFToken            string
	Flash                string
	Warning              string
	Error                string
	IsAuthenticated      int
	API                  string
	CSSVersion           string
	StripeSecretKey      string
	StripePublishableKey string
}

var functions = template.FuncMap{
	"formatCurrency":    formatCurrency,
	"getSizeOfProducts": getSizeOfProducts,
	"pagination":        pagination,
}

func pagination(currentpage int, pagesize int) template.HTML {
	html := ""
	// previous button
	// vd: pre 1 2 nex
	if currentpage-1 > 0 {
		html += fmt.Sprintf(`<li onclick="previous(%d)" class="page-item">
                                <a class="page-link">Previous</a>
                            </li>`, currentpage)
	} else {
		html += `<li class="page-item disabled">
                                <a class="page-link">Previous</a>
				</li>`
	}

	// currentpage - 1
	if currentpage-1 >= 1 {
		html += fmt.Sprintf(`<li onclick="previous(%d)" class="page-item">
                                    <a class="page-link">%d</a>
                                </li>`, currentpage, currentpage-1)
	}

	// currentpage
	html += fmt.Sprintf(`<li class="page-item active">
                                <a class="page-link">%d</a>
						</li>`, currentpage)

	// current page + 1
	if currentpage+1 <= pagesize {
		html += fmt.Sprintf(`<li onclick="next(%d, %d)" class="page-item">
								<a class="page-link">%d</a>
							</li>`, currentpage, pagesize, currentpage+1)
	}

	// Next
	// pre 2 3 nex
	if currentpage+1 <= pagesize {
		html += fmt.Sprintf(`<li onclick="next(%d, %d)" class="page-item">
                                    <a class="page-link">Next</a>
				</li>`, currentpage, pagesize)
	} else if currentpage+1 > pagesize {
		html += `<li class="page-item disabled">
					<a class="page-link">Next</a>
				</li>`
	}

	return template.HTML(html)
}

func formatCurrency(n int) string {
	f := float32(n / 1)
	return fmt.Sprintf("$%.2f", f)
}

func getSizeOfProducts(products []string) bool {
	if len(products) > 0 {
		return true
	} else {
		return false
	}
}

//go:embed templates
var templateFS embed.FS

func (app *application) addDefaultData(td *templateData, r *http.Request) *templateData {
	td.API = app.config.api
	td.StripeSecretKey = app.config.stripe.secret
	td.StripePublishableKey = app.config.stripe.key
	return td
}

func (app *application) renderTemplate(w http.ResponseWriter, r *http.Request, page string, td *templateData, partials ...string) error {
	var t *template.Template
	var err error

	templateToRender := fmt.Sprintf("templates/%s.page.gohtml", page)
	_, templateInMap := app.templateCache[templateToRender]

	if app.config.env == "production" && templateInMap {
		t = app.templateCache[templateToRender]
	} else {

		t, err = app.parseTemplate(partials, page, templateToRender)
		if err != nil {
			app.errorLog.Println(err)
			return err
		}
	}

	if td == nil {
		td = &templateData{}
	}

	td = app.addDefaultData(td, r)

	err = t.Execute(w, td)
	if err != nil {
		app.errorLog.Println(err)
		return err
	}

	return nil
}

func (app *application) parseTemplate(partials []string, page, templateToRender string) (*template.Template, error) {
	var t *template.Template
	var err error

	// build partials
	if len(partials) > 0 {
		for i, x := range partials {
			partials[i] = fmt.Sprintf("templates/%s.partial.gohtml", x)
		}
	}

	partials = append(partials, templateToRender, "templates/base.layout.gohtml")

	t, err = template.New(fmt.Sprintf("%s.page.gohtml", page)).Funcs(functions).ParseFS(templateFS, partials...)

	if err != nil {
		app.errorLog.Println(err)
		return nil, err
	}

	app.templateCache[templateToRender] = t
	return t, nil
}
