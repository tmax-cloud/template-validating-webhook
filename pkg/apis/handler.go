package apis

import (
	"encoding/json"
	"fmt"
	"net/http"

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
	var newParamNameVal, oldParamNameVal string
	var newParameters []interface{}
	var oldParameters []interface{}

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

	oldObject := req["request"].(map[string]interface{})["oldObject"].(map[string]interface{})
	if oldObject["spec"].(map[string]interface{})[scope] != nil {
		oldTemplateName = oldObject["spec"].(map[string]interface{})[scope].(map[string]interface{})["metadata"].(map[string]interface{})["name"].(string)
		oldParameters = oldObject["spec"].(map[string]interface{})[scope].(map[string]interface{})["parameters"].([]interface{})
	}
	checkTemplateName := newTemplateName == oldTemplateName

	newNameParam := GetNameParameterAsMap(newParameters)
	oldNameParam := GetNameParameterAsMap(oldParameters)

	if _, exist := newNameParam["NAME"]; exist {
		newParamNameVal = newNameParam["NAME"]
		oldParamNameVal = oldNameParam["NAME"]
		checkParamName := newParamNameVal == oldParamNameVal
		return checkTemplateName && checkParamName
	}

	return checkTemplateName
}

func GetNameParameterAsMap(parameters []interface{}) map[string]string {
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

	nameParam := make(map[string]string)
	for _, param := range Params {
		nameParam[param.Name] = param.Value.StrVal
	}
	return nameParam
}
