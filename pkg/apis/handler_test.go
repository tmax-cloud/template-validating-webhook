package apis

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestValidate(t *testing.T) {

	admissionReview := []byte(`{
		"kind": "AdmissionReview", 
		"request": { 
			"object": { 
				"spec" : { 
					"template": { 
						"metadata" : { 
							"name" : "new-template"} } } }, 
			"oldObject": { 
				"spec" : {
					"template": { 
						"metadata" : { 
							"name" : "old-template"} } } } } }`)

	var req map[string]interface{}

	err := json.Unmarshal(admissionReview, &req)
	require.NoError(t, err)

	assert.Equal(t, false, Validate(req))
}
