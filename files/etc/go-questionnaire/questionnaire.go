package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"io/ioutil"
	"net/http"
	"os"
	"time"

	"github.com/bluele/gforms"
)

func check(e error) {
	if e != nil {
		panic(e)
	}
}

// type questionnaire struct {
// 	Questions []question
// 	Message   string
// }

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

	//dump the question answers into a json
	fmt.Println(form.CleanedData)
	jsonString, err := json.Marshal(form.CleanedData)
	check(err)

	f, err := os.Create("questionnaire-answers/" + time.Now().Format(time.RFC850) + ".json")
	defer f.Close()
	check(err)
	f.WriteString([]byte(jsonString))

	// for k, v := range form.CleanedData {
	// os.Create("questionnaire-answers/" + time.Now().Format(time.RFC850) + ".json")
	// 	fmt.Println(string(k))
	// 	fmt.Println(fmt.Sprintf("%v", v))
	// }
}

func main() {
	http.HandleFunc("/", questionnaireHandler)
	http.ListenAndServe(":8080", nil)
}
