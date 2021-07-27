package schemas

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
