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

//StringToHTML takes a string and returns HTML
func StringToHTML(s string) template.HTML {
	return template.HTML(s)
}

func defineValidators(q map[string]interface{}) gforms.Validators {
	var validators []gforms.Validator

	if q["required"].(bool) {
		validators = append(validators, gforms.Required())
	}
	if v, ok := q["maxLength"].(float64); ok {
		validators = append(validators, gforms.MaxLengthValidator(int(v)))
	}

	return gforms.Validators(validators)
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
				defineValidators(q),
			))
		case "numberQuestion":
			fields = append(fields, gforms.NewIntegerField(
				q["question"].(string),
				defineValidators(q),
			))
		case "multipleChoiceQuestion":
			fields = append(fields, gforms.NewMultipleTextField(q["question"].(string),
				defineValidators(q),
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
	timeString := fmt.Sprintf("%d-%02d-%02dT%02d:%02d:%02d\n",
		time.Now().Year(), time.Now().Month(), time.Now().Day(),
		time.Now().Hour(), time.Now().Minute(), time.Now().Second())
	f, err := os.Create("questionnaire-answers/" + timeString + ".json")
	defer f.Close()
	check(err)
	f.WriteString(string(jsonString))
}

func main() {
	http.HandleFunc("/", questionnaireHandler)
	http.ListenAndServe(":8080", nil)
}
