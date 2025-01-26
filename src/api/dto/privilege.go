package dto

type CheckPrivilegeRequest struct {
	User     string `json:"user"`     // The user whose privileges are being checked
	Action   string `json:"action"`   // The action the user wants to perform
	Resource string `json:"resource"` // The resource the action is performed on
}

type CheckPrivilegeResponse struct {
	Allowed bool `json:"allowed"` // Whether the user is allowed to perform the action
}

type AssignPrivilegeRequest struct {
	User     string `json:"user"`
	Action   string `json:"action"`
	Resource string `json:"resource"`
}

type AssignPrivilegeResponse struct {
	Success bool `json:"success"`
}

// ErrorResponse represents a structured error response returned by the API
//
//	@Description	Standard error response format used by the API
//	@Example		{"error": "bad_request", "code": 400, "message": "Failed to parse request"}
type ErrorResponse struct {
	//	@Description	Error type, like "bad_request", "internal_server_error", etc.
	//	@Example		"bad_request"
	Error string `json:"error"`

	//	@Description	HTTP status code for the error
	//	@Example		400
	Code int `json:"code"`

	//	@Description	Detailed error message explaining the issue
	//	@Example		"Failed to parse request"
	Message string `json:"message"`
}
