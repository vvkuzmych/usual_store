package main

import (
	"embed"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
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
	UserID               int
	API                  string
	CSSVersion           string
	StripeSecretKey      string
	StripePublishableKey string
	Nonce 				 string
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
		td.UserID = app.Session.GetInt(r.Context(), "userID")
	} else {
		td.IsAuthenticated = 0
		td.UserID = 0
	}

	return td
}

func (app *application) renderTemplate(w http.ResponseWriter, r *http.Request, page string, td *templateData, partials ...string) error {
	// Determine the template to render
	templateToRender := fmt.Sprintf("templates/%s.page.gohtml", page)

	// Retrieve or parse the template
	t, exists := app.templateCache[templateToRender]
	if !exists {
		var err error
		t, err = app.parseTemplate(partials, page, templateToRender)
		if err != nil {
			app.errorLog.Printf("Error parsing template %s: %v", templateToRender, err)
			return err
		}
	}

	// Add default data if templateData is nil
	if td == nil {
		td = &templateData{}
	}
	td = app.addDefaultData(td, r)

	// Execute the template and handle errors
	if err := t.Execute(w, td); err != nil {
		app.errorLog.Printf("Error executing template %s: %v", templateToRender, err)
		return err
	}

	return nil
}

func (app *application) parseTemplate(partials []string, page, templateToRender string) (*template.Template, error) {
	// Prepend the path and extension to each partial
	for i, partial := range partials {
		partials[i] = fmt.Sprintf("templates/%s.partial.gohtml", partial)
	}

	// Base templates to parse
	baseTemplates := []string{"templates/base.layout.gohtml", templateToRender}
	if len(partials) > 0 {
		baseTemplates = append(baseTemplates, partials...)
	}

	// Parse the templates
	templateName := fmt.Sprintf("%s.page.gohtml", page)
	t, err := template.New(templateName).Funcs(functions).ParseFS(templateFS, baseTemplates...)
	if err != nil {
		app.errorLog.Printf("Error parsing template %s: %v", templateToRender, err)
		return nil, err
	}

	// Cache the parsed template
	app.templateCache[templateToRender] = t
	return t, nil
}

func (app *application) getEnvData() map[string]string {
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
