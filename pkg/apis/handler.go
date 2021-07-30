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

	request := req["request"].(map[string]interface{})
	object := request["object"].(map[string]interface{})
	newTemplateName := object["spec"].(map[string]interface{})["template"].(map[string]interface{})["metadata"].(map[string]interface{})["name"].(string)

	oldObject := request["oldObject"].(map[string]interface{})
	oldTemplateName := oldObject["spec"].(map[string]interface{})["template"].(map[string]interface{})["metadata"].(map[string]interface{})["name"].(string)

	return newTemplateName == oldTemplateName
}
