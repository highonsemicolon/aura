package dto

type CheckPrivilegeRequest struct {
	User     string `json:"user"`
	Action   string `json:"action"`
	Resource string `json:"resource"`
}

type CheckPrivilegeResponse struct {
	Allowed bool `json:"allowed"`
}

type ErrorResponse struct {
	Error   string `json:"error"`
	Code    int    `json:"code"`
	Message string `json:"message"`
}
