package main

import (
	"embed"
	"fmt"
	"html/template"
	"net/http"
)

type templateData struct {
	StringMap    map[string]string
	IntMap       map[string]int
	FloatMap     map[string]float32
	Data         map[string]interface{}
	CSRFToken    string
	Flash        map[string]string
	Warning      string
	Error        string
	Port         int
	Host         string
	IsLogged     bool
	Quantity     int
	Page         string
	ConstantPage map[string]string
}

var functions = template.FuncMap{
	"formatCurrency":    formatCurrency,
	"getSizeOfProducts": getSizeOfProducts,
	"pagination":        pagination,
	"renderHeader":      renderHeader,
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

func renderHeader(str string) template.HTML {
	html := ""

	// tạo constant để lưu giá trị của các pages
	type PageURL struct {
		Page string
		Url  string
	}
	pages := []PageURL{
		{Page: "HOME", Url: "/"},
		{Page: "PRODUCT", Url: "/product"},
		{Page: "CONTACT", Url: "/contact"},
	}

	// duyệt key-value trong page
	for _, value := range pages {
		if str == value.Page {
			html += fmt.Sprintf(`<a href="%s" class="nav-item nav-link active">%s</a>`, value.Url, value.Page) // được active class
		} else {
			html += fmt.Sprintf(`<a href="%s" class="nav-item nav-link">%s</a>`, value.Url, value.Page) // không có active class
		}
	}

	return template.HTML(html)
}

//go:embed templates
var templateFS embed.FS

func (app *application) addDefaultData(td *templateData, r *http.Request) *templateData {
	td.Port = app.config.port
	td.Host = app.config.host

	session, err := app.Redis.Get(r, app.config.session_key)
	if err != nil {
		app.errorLog.Println(err.Error())
		panic(err)
	}

	// lấy thông tin về trạng thái đăng nhập của người dùng
	if value, ok := session.Values["IsLogged"].(bool); !ok {
		td.IsLogged = false // value không phải là bool (chưa có giá trị trong session) --> gán giá trị default là false
	} else {
		td.IsLogged = value // value có type là bool
	}

	/*
	 * lấy thông tin số lượng loại sản phẩm trong giỏ hàng
	 * nếu không có đăng nhập hoặc giỏ hàng không có sản phẩm --> quantity mặc định là 0
	 * ngược lại, lấy tổng số loại sản phẩm trong giỏ hàng
	 */

	userEmail, _ := session.Values["UserEmail"].(string)

	if !td.IsLogged { // nếu người dùng chưa đăng nhập
		td.Quantity = 0
	} else {
		// xử lí truy vấn và gán giá trị quantity
		quantity, _ := app.DB.CountNumofTypeProductInCartByUserEmail(userEmail)
		td.Quantity = quantity
	}

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
