package api

import (
	"aura/src/middleware"

	"github.com/gin-gonic/gin"
)

func Register(router *gin.Engine, h *PrivilegeHandler) {

	useridMiddleware := middleware.UserIDMiddleware

	router.POST("/check", h.checkPrivilege)
	router.POST("role", useridMiddleware, h.assignPrivilege)
	// router.DELETE(("role"), h.revokePrivilege)
}
