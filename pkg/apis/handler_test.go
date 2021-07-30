package apis

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type testObj struct {
	admissionReview []byte
	result          bool
}

func TestValidate(t *testing.T) {
	set := []testObj{
		{
			admissionReview: []byte(`{
				"kind": "AdmissionReview", 
				"request": { 
					"object": { 
						"spec" : { 
							"template": { 
								"metadata" : { 
									"name" : "new-template"} } } }, 
					"oldObject": { 
						"spec" : {
							"clustertemplate": { 
								"metadata" : { 
									"name" : "old-template"} } } } } }`),
			result: false,
		},
		{
			admissionReview: []byte(`{
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
									"name" : "old-template"} } } } } }`),
			result: false,
		},
		{
			admissionReview: []byte(`{
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
									"name" : "new-template"} } } } } }`),
			result: true,
		},
	}

	for _, s := range set {
		req := map[string]interface{}{}
		err := json.Unmarshal(s.admissionReview, &req)
		require.NoError(t, err)
		assert.Equal(t, s.result, Validate(req))
	}
}
