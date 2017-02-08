package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"io/ioutil"
	"net/http"

	"github.com/bluele/gforms"
)

func check(e error) {
	if e != nil {
		panic(e)
	}
}

type question interface {
	RenderQuestion() string
}

type multipleQuestion struct {
	question string
	choices  []string
	isRadio  bool
}

type textBoxQuestion struct {
	question  string
	charLimit int
}

type numberTextBoxQuestion struct {
	question string
}

func (q multipleQuestion) RenderQuestion() string {
	return "Multiple question."
}

func (q textBoxQuestion) RenderQuestion() string {
	return "Text Box question."
}

func (q numberTextBoxQuestion) RenderQuestion() string {
	return "Number Text Box question."
}

type questionnaire struct {
	Questions []question
	Message   string
}

//StringToHTML takes a string and returns HTML
func StringToHTML(s string) template.HTML {
	return template.HTML(s)
}

func questionnaireHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")

	//load the json file with questions and create fields list
	dat, err := ioutil.ReadFile("questions-en.json")
	check(err)
	var questions map[string]interface{}
	err = json.Unmarshal(dat, &questions)
	check(err)

	var fields []gforms.Field

	for _, question := range questions["questions"].([]interface{}) {
		q := question.(map[string]interface{})
		switch q["type"].(string) {
		case "textBoxQuestion":
			fields = append(fields, gforms.NewTextField(
				q["question"].(string),
				gforms.Validators{
					gforms.Required(),
					gforms.MaxLengthValidator(int(q["maxChars"].(float64))),
				},
				gforms.TextInputWidget(
					map[string]string{
						"name": q["name"].(string),
					},
				),
			))
		case "numberTextBoxQuestion":
			fields = append(fields, gforms.NewFloatField(
				q["question"].(string),
				gforms.Validators{
					gforms.Required(),
				},
				gforms.TextInputWidget(
					map[string]string{
						"name": q["name"].(string),
					},
				),
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

	//parse the request
	if r.Method == "GET" {
		t.Execute(w, form)
		return
	}
	if !form.IsValid() {
		t.Execute(w, form)
		return
	}

	//parse the form fields
	for k, v := range form.Data {
		fmt.Println(string(k))
		fmt.Println(fmt.Sprintf("%v", v.Value))
	}
}

func main() {
	http.HandleFunc("/", questionnaireHandler)
	http.ListenAndServe(":8080", nil)
}
