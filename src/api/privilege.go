package api

import (
	"aura/src/api/dto"
	services "aura/src/services"
	"net/http"

	"github.com/gin-gonic/gin"
)

type PrivilegeHandler struct {
	ps services.PrivilegeServiceInterface
}

func NewPrivilegeHandler(pc services.PrivilegeServiceInterface) *PrivilegeHandler {
	return &PrivilegeHandler{pc}
}

func (h *PrivilegeHandler) checkPrivilege(c *gin.Context) {
	var req dto.CheckPrivilegeRequest
	if err := c.BindJSON(&req); err != nil {
		h.writeError(c, http.StatusBadRequest, "bad_request", "Failed to parse request")
		return
	}

	if err := h.valiateInputs(req.User, req.User, req.Action, req.Resource); err != nil {
		h.writeError(c, http.StatusBadRequest, "bad_request", "Invalid input parameters")
		return
	}

	allowed, err := h.ps.CheckRole(req.User, req.Action, req.Resource)
	if err != nil {
		h.writeError(c, http.StatusInternalServerError, "internal_server_error", err.Error())
		return
	}

	if allowed {
		h.writeJSON(c, http.StatusOK, dto.CheckPrivilegeResponse{Allowed: true})
		return
	}

	c.JSON(http.StatusOK, dto.CheckPrivilegeResponse{Allowed: false})
}

/*
func (h *PrivilegeHandler) assignPrivilege(c *gin.Context) {
	userID := c.GetString("userID")
	h.valiateInputs(userID)

	c.JSON(http.StatusOK, gin.H{"success": true})
}

func (h *PrivilegeHandler) revokePrivilege(c *gin.Context) {
	userID := c.GetString("userID")
	h.valiateInputs(userID)

	c.JSON(http.StatusOK, gin.H{"success": true})
}
*/

func (h *PrivilegeHandler) valiateInputs(input ...string) error {
	for _, s := range input {
		if s == "" {
			return http.ErrBodyNotAllowed
		}
	}
	return nil

}

func (h *PrivilegeHandler) writeError(c *gin.Context, status int, code string, message string) {
	error := dto.ErrorResponse{
		Error:   code,
		Code:    status,
		Message: message,
	}

	c.JSON(status, error)
}

func (h *PrivilegeHandler) writeJSON(c *gin.Context, status int, data interface{}) {
	c.Header("Content-Type", "application/json")
	c.JSON(status, data)
}
