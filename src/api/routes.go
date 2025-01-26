package api

import (
	"aura/src/middleware"

	"github.com/gin-gonic/gin"
)

func Register(router *gin.Engine, h *PrivilegeHandler) {

	useridMiddleware := middleware.UserIDMiddleware

	router.POST("aura", useridMiddleware, h.assignPrivilege)
	router.POST("aura/check", h.checkPrivilege)
	// router.DELETE(("role"), h.revokePrivilege)
}
