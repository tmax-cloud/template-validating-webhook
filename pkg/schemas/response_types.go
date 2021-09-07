package schemas

import "k8s.io/apimachinery/pkg/util/intstr"

type ResponseBody struct {
	Kind       string   `json:"kind"`
	ApiVersion string   `json:"apiVersion"`
	Response   Response `json:"response"`
}

type Response struct {
	Allowed bool   `json:"allowed"`
	Status  Status `json:"status"`
}

type Status struct {
	Reason string `json:"reason"`
}

type ParamSpec struct {
	// A description of the parameter.
	// Provide more detailed information for the purpose of the parameter, including any constraints on the expected value.
	// Descriptions should use complete sentences to follow the console’s text standards.
	// Don’t make this a duplicate of the display name.
	Description string `json:"description,omitempty"`
	// The user-friendly name for the parameter. This will be displayed to users.
	DisplayName string `json:"displayName,omitempty"`
	// The name of the parameter. This value is used to reference the parameter within the template.
	Name string `json:"name"`
	// Indicates this parameter is required, meaning the user cannot override it with an empty value.
	// If the parameter does not provide a default or generated value, the user must supply a value.
	Required bool `json:"required,omitempty"`
	// A default value for the parameter which will be used if the user does not override the value when instantiating the template.
	// Avoid using default values for things like passwords, instead use generated parameters in combination with Secrets.
	Value intstr.IntOrString `json:"value,omitempty"`
	// Set the data type of the parameter.
	// You can specify string and number for a string or integer type.
	// If not specified, it defaults to string.
	// +kubebuilder:validation:Enum:=string;number
	ValueType string `json:"valueType,omitempty"`
}
