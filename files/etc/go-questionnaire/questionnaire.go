package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/bluele/gforms"
)

func check(e error) {
	if e != nil {
		panic(e)
	}
}

//StringToHTML takes a string and returns HTML
func StringToHTML(s string) template.HTML {
	return template.HTML(s)
}

//given a question q, will define its validators from given fields
func defineValidators(q map[string]interface{}, locale string) gforms.Validators {
	dat, err := ioutil.ReadFile("errorMessages.json")
	check(err)
	var errors map[string]interface{}
	err = json.Unmarshal(dat, &errors)
	check(err)

	var validators []gforms.Validator

	errorMessages := errors[locale].(map[string]interface{})
	if q["required"].(bool) {
		validators = append(validators, gforms.Required(errorMessages["required"].(string)))
	}
	if v, ok := q["maxLength"].(float64); ok {
		validators = append(validators, gforms.MaxLengthValidator(int(v), fmt.Sprintf(errorMessages["maxChars"].(string), int(v))))
	}

	return gforms.Validators(validators)
}

//given an input string 'locale' returns a valid locale
func getLocale(locale string) string {
	validLocales := [...]string{"en", "ar"}

	for _, l := range validLocales {
		if l == locale {
			return l
		}
	}
	return "en"
}

func questionnaireHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	locale := strings.SplitAfterN(r.URL.Path, "/", 2)[1]
	locale = getLocale(locale)

	//load the json file with questions and create fields list
	dat, err := ioutil.ReadFile("questions.json")
	check(err)
	var questions map[string]interface{}
	err = json.Unmarshal(dat, &questions)
	check(err)

	var fields []gforms.Field

	for _, question := range questions[locale].([]interface{}) {
		q := question.(map[string]interface{})
		switch q["type"].(string) {
		case "textBoxQuestion":
			fields = append(fields, gforms.NewTextField(
				q["question"].(string),
				defineValidators(q, locale),
			))
		case "numberQuestion":
			fields = append(fields, gforms.NewIntegerField(
				q["question"].(string),
				defineValidators(q, locale),
			))
		case "multipleChoiceQuestion":
			fields = append(fields, gforms.NewMultipleTextField(q["question"].(string),
				defineValidators(q, locale),
				gforms.CheckboxMultipleWidget(
					map[string]string{
						"class": q["name"].(string),
					},
					func() gforms.CheckboxOptions {
						var retval [][]string
						for _, v := range q["choices"].([]interface{}) {
							retval = append(retval, []string{
								v.(string), v.(string), "false", "false",
							})
						}
						return gforms.StringCheckboxOptions(retval)
					}),
			))
		case "singleChoiceQuestion":
			fields = append(fields, gforms.NewTextField(q["question"].(string),
				defineValidators(q, locale),
				gforms.RadioSelectWidget(
					map[string]string{
						"class": q["name"].(string),
					},
					func() gforms.RadioOptions {
						var retval [][]string
						for _, v := range q["choices"].([]interface{}) {
							retval = append(retval, []string{
								v.(string), v.(string), "false", "false",
							})
						}
						return gforms.StringRadioOptions(retval)
					}),
			))
		}
	}

	//prepare the template
	funcMap := template.FuncMap{
		"stringToHTML": StringToHTML,
	}
	t := template.Must(template.New("questionnaire.html").Funcs(funcMap).ParseFiles("questionnaire.html"))

	// prepare the form
	userForm := gforms.DefineForm(gforms.NewFields(fields...))
	form := userForm(r)

	//parse the requestsk
	if r.Method == "GET" {
		t.Execute(w, form)
		return
	}
	if !form.IsValid() {
		t.Execute(w, form)
		return
	}

	//dump the question answers into a json
	fmt.Println(form.CleanedData)
	jsonString, err := json.Marshal(form.CleanedData)
	check(err)
	timeString := fmt.Sprintf("%d-%02d-%02dT%02d:%02d:%02d",
		time.Now().Year(), time.Now().Month(), time.Now().Day(),
		time.Now().Hour(), time.Now().Minute(), time.Now().Second())
	f, err := os.Create("questionnaire-answers/" + timeString + ".json")
	defer f.Close()
	check(err)
	f.WriteString(string(jsonString))

	//TODO: when an answer is recorded send it somewhere for Elpis to access off-camp
}

func main() {
	http.HandleFunc("/", questionnaireHandler)
	http.ListenAndServe(":8080", nil)
}
