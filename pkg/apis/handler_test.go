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
		{ // Case1 : new-Template cannot compare value with old-ClusterTemplate
			admissionReview: []byte(`{
			"kind": "AdmissionReview",
			"request": {
				"object": {
					"spec" : {
						"template": {
							"metadata" : {
								"name" : "new-template"},
							"parameters": [ 
								{
									"name": "NAME",
									"value": "test-name" 
								} ] } } },
				"oldObject": {
					"spec" : {
						"clustertemplate": {
							"metadata" : {
								"name" : "new-template"},
							"parameters": [ 
								{
									"name": "NAME",
									"value": "test-name" 
								} ] } } } } }`),
			result: false,
		},
		{ // Case2 : new-template Name != old-template Name
			admissionReview: []byte(`{
			"kind": "AdmissionReview",
			"request": {
				"object": {
					"spec" : {
						"template": {
							"metadata" : {
								"name" : "new-template"},
							"parameters": [ 
								{
									"name": "NAME",
									"value": "test-name" 
								} ] } } },
				"oldObject": {
					"spec" : {
						"template": {
							"metadata" : {
								"name" : "old-template"},
							"parameters": [ 
								{
									"name": "NAME",
									"value": "test-name" 
								} ] } } } } }`),
			result: false,
		},
		{ // Case3 : new-Parameter Value of "NAME" != old-Parameter Value of "NAME"
			admissionReview: []byte(`{
				"kind": "AdmissionReview",
				"request": {
				  "object": {
					"spec": {
					  "template": {
						"metadata": {
						  "name": "new-template" },
						"parameters": [ 
							{
								"name": "NAME",
								"value": "new-name" 
							},
							{
								"name": "IMAGE",
								"value": "new-image" 
							} ] } } },
				  "oldObject": {
					"spec": {
					  "template": {
						"metadata": {
						  "name": "new-template" },
						"parameters": [ 
							{
								"name": "NAME",
								"value": "old-name"  
							},
							{
								"name": "IMAGE",
								"value": "old-image"  
							} ] } } } } }`),
			result: false,
		},
		{ // Case4 : Template Name && Parameter "NAME" values are same
			admissionReview: []byte(`{
				"kind": "AdmissionReview",
				"request": {
				  "object": {
					"spec": {
					  "template": {
						"metadata": {
						  "name": "new-template" },
						"parameters": [ 
							{
								"name": "NAME",
								"value": "test-name" 
							},
							{
								"name": "IMAGE",
								"value": "new-image" 
							} ] } } },
				  "oldObject": {
					"spec": {
					  "template": {
						"metadata": {
						  "name": "new-template" },
						"parameters": [ 
							{
								"name": "NAME",
								"value": "test-name"  
							},
							{
								"name": "IMAGE",
								"value": "old-image"  
							} ] } } } } }`),
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
