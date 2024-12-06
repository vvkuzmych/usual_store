package main

import (
	"embed"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"strings"
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
	"formatCurrency": formatCurrency,
}

func formatCurrency(n int) string {
	f := float32(n) / float32(100)
	return fmt.Sprintf("$%.2f", f)
}

//go:embed templates
var templateFS embed.FS

func (app *application) addDefaultData(td *templateData, r *http.Request) *templateData {
	td.API = app.config.api
	env := app.getEnvData()
	app.config.stripe.secret = env["stripe_secret"]
	app.config.stripe.key = env["publishable_key"]
	td.StripeSecretKey = app.config.stripe.secret
	td.StripePublishableKey = app.config.stripe.key
	if app.Session.Exists(r.Context(), "userID") {
		td.IsAuthenticated = 1
	} else {
		td.IsAuthenticated = 0
	}

	return td
}

func (app *application) renderTemplate(w http.ResponseWriter, r *http.Request, page string, td *templateData, partials ...string) error {
	var t *template.Template
	var err error
	templateToRender := fmt.Sprintf("templates/%s.page.gohtml", page)
	_, templateInMap := app.templateCache[templateToRender]

	if templateInMap {
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
		for i, partial := range partials {
			partials[i] = fmt.Sprintf("templates/%s.partial.gohtml", partial)
		}
	}
	if len(partials) > 0 {
		t, err = template.New(fmt.Sprintf("%s.page.gohtml", page)).Funcs(functions).ParseFS(templateFS, "templates/base.layout.gohtml", strings.Join(partials, ","), templateToRender)
	} else {
		t, err = template.New(fmt.Sprintf("%s.page.gohtml", page)).Funcs(functions).ParseFS(templateFS, "templates/base.layout.gohtml", templateToRender)
	}
	if err != nil {
		app.errorLog.Println(err)
		return nil, err
	}
	app.templateCache[templateToRender] = t

	return t, nil
}

func (app *application) getEnvData() map[string]string {
	//err := godotenv.Load()
	//if err != nil {
	//	log.Printf("Error loading .env file: %v", err)
	//	log.Fatal("Exiting due to missing .env file")
	//}

	// Get the publishable key from the environment variable
	publishableKey := os.Getenv("PUBLISHABLE_KEY")
	if publishableKey == "" {
		log.Fatalf("PUBLISHABLE_KEY not set in .env file")
	}
	stringMap := make(map[string]string)
	stringMap["publishable_key"] = publishableKey

	//Get the secret key from the environment variable
	stripeSecret := os.Getenv("SECRET")
	if stripeSecret == "" {
		log.Fatalf("SECRET not set in .env file")
	}
	stringMap["stripe_secret"] = stripeSecret
	return stringMap
}
