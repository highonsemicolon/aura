package api

import (
	"aura/src/api/dto"
	services "aura/src/services"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

type PrivilegeHandler struct {
	ps services.PrivilegeServiceInterface
}

func NewPrivilegeHandler(pc services.PrivilegeServiceInterface) *PrivilegeHandler {
	return &PrivilegeHandler{pc}
}

// checkPrivilege checks if a user has permission to perform a specific action
//
//	@Summary		Check user privilege
//	@Description	Checks if a user has the specified action on a resource
//	@Tags			privileges
//	@Accept			json
//	@Produce		json
//	@Param			request	body		dto.CheckPrivilegeRequest	true	"Privilege check request"
//	@Success		200		{object}	dto.CheckPrivilegeResponse
//	@Failure		400		{object}	dto.ErrorResponse
//	@Failure		500		{object}	dto.ErrorResponse
//	@Router			/check [post]
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

// assignPrivilege assigns a privilege to a user.
//
//	@Summary		Assign privilege
//	@Description	Allows an assigner to assign a privilege to a user
//	@Tags			privileges
//	@Accept			json
//	@Produce		json
//	@Param			userID	header		string						true	"Assigner's User ID"
//	@Param			body	body		dto.AssignPrivilegeRequest	true	"Request payload to assign privilege"
//	@Success		201		{object}	dto.AssignPrivilegeResponse	"Successfully assigned privilege"
//	@Failure		400		{object}	dto.AssignPrivilegeResponse	"Invalid input or bad request"
//	@Failure		403		{object}	dto.ErrorResponse			"Unauthorised to perform this action"
//	@Router			/privileges/assign [post]
func (h *PrivilegeHandler) assignPrivilege(c *gin.Context) {
	assigner := c.GetString("userID")

	var req dto.AssignPrivilegeRequest
	if err := c.BindJSON(&req); err != nil {
		h.writeError(c, http.StatusBadRequest, "bad_request", "Failed to parse request")
		return
	}

	if err := h.valiateInputs(assigner, req.User, req.Action, req.Resource); err != nil {
		h.writeError(c, http.StatusBadRequest, "bad_request", "Invalid input parameters")
		return
	}

	if err := h.ps.AssignRole(assigner, req.User, req.Action, req.Resource); err != nil {
		h.writeError(c, http.StatusForbidden, "unauthorised", err.Error())
		log.Println(err)
		return
	}

	h.writeJSON(c, http.StatusCreated, dto.AssignPrivilegeResponse{Success: true})
}

/*
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
