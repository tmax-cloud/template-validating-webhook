package apis

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/tmax-cloud/template-validating-webhook/pkg/schemas"
)

func CheckInstanceUpdatable(w http.ResponseWriter, r *http.Request) {
	var req map[string]interface{}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		fmt.Println("request json decoding error")
		return
	}

	result := Validate(req)

	body := schemas.ResponseBody{
		Kind:       "AdmissionReview",
		ApiVersion: "admission.k8s.io/v1beta1",
		Response: schemas.Response{
			Allowed: result,
			Status: schemas.Status{
				Reason: "TemplateInstance with the same name already exists",
			},
		},
	}

	if err := json.NewEncoder(w).Encode(body); err != nil {
		fmt.Println("response json encoding error")
	}
}

func Validate(req map[string]interface{}) bool {

	var scope, newTemplateName, oldTemplateName string
	var newParamNameVal, oldParamNameVal, paramName string
	var oldObject map[string]interface{}
	var newParameters []interface{}
	var oldParameters []interface{}
	var oldStatusObj []interface{}

	object := req["request"].(map[string]interface{})["object"].(map[string]interface{})

	if object["spec"].(map[string]interface{})["clustertemplate"] != nil {
		scope = "clustertemplate"
		newTemplateName = object["spec"].(map[string]interface{})["clustertemplate"].(map[string]interface{})["metadata"].(map[string]interface{})["name"].(string)
		newParameters = object["spec"].(map[string]interface{})["clustertemplate"].(map[string]interface{})["parameters"].([]interface{})
	} else {
		scope = "template"
		newTemplateName = object["spec"].(map[string]interface{})["template"].(map[string]interface{})["metadata"].(map[string]interface{})["name"].(string)
		newParameters = object["spec"].(map[string]interface{})["template"].(map[string]interface{})["parameters"].([]interface{})
	}

	if req["request"].(map[string]interface{})["oldObject"] != nil {
		oldObject = req["request"].(map[string]interface{})["oldObject"].(map[string]interface{})
		if oldObject["spec"].(map[string]interface{})[scope] != nil {
			oldTemplateName = oldObject["spec"].(map[string]interface{})[scope].(map[string]interface{})["metadata"].(map[string]interface{})["name"].(string)
			oldParameters = oldObject["spec"].(map[string]interface{})[scope].(map[string]interface{})["parameters"].([]interface{})
			oldStatusObj = oldObject["status"].(map[string]interface{})[scope].(map[string]interface{})["objects"].([]interface{})
		}
	}

	if oldObject == nil { // create templateInstance for the first time
		return true
	}

	checkTemplateName := newTemplateName == oldTemplateName

	newParam := GetParameterAsMap(newParameters)
	oldParam := GetParameterAsMap(oldParameters)

	if len(oldStatusObj) != 0 {
		paramName = CheckObjectNameParameter(oldStatusObj)
		if paramName == "false" {
			fmt.Println("There are multiple NAME parameters")
			return false
		}
		if _, exist := newParam[paramName]; exist {
			newParamNameVal = newParam[paramName]
			oldParamNameVal = oldParam[paramName]
			checkParamName := newParamNameVal == oldParamNameVal
			return checkTemplateName && checkParamName
		}
	}
	return checkTemplateName
}

func CheckObjectNameParameter(objs []interface{}) string {
	var name string
	var names []string
	var ref string

	// Get objects.metadata.name and extract string between { }
	for _, o := range objs {
		obj := o.(map[string]interface{})
		fullName := obj["metadata"].(map[string]interface{})["name"].(string)
		left := strings.IndexAny(fullName, "{")
		right := strings.IndexAny(fullName, "}")

		if left == -1 { // In case of No parameter in objects.metadata.name
			name = ""
		} else {
			s1 := strings.Split(fullName, "")
			s2 := s1[left+1 : right]
			name = strings.Join(s2, "")
		}
		names = append(names, name)
	}

	for _, n := range names {
		if n != "" {
			ref = n
			break
		}
	}

	for _, n := range names {
		if n == ref || n == "" {
			continue
		} else {
			return "false"
		}
	}
	return ref
}

func GetParameterAsMap(parameters []interface{}) map[string]string {
	Params := []schemas.ParamSpec{}

	for _, p := range parameters {
		m := p.(map[string]interface{})
		param := schemas.ParamSpec{}
		if name, ok := m["name"].(string); ok {
			param.Name = name
		}
		if value, ok := m["value"].(string); ok {
			param.Value.StrVal = value
		}
		Params = append(Params, param)
	}

	totalParam := make(map[string]string)
	for _, param := range Params {
		totalParam[param.Name] = param.Value.StrVal
	}
	return totalParam
}
