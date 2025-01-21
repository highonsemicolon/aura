package api

import (
	"github.com/gin-gonic/gin"
)

func Register(router *gin.Engine, h *PrivilegeHandler) {
	router.POST("/check", h.checkPrivilege)
	// router.POST(("assign"), h.assignPrivilege)
	// router.POST("/revoke", h.revokePrivilege)
}
