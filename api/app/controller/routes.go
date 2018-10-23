package controller

import (
	"github.com/gin-gonic/gin"
	"github.com/itcam/go-devops/api/app/controller/openldap"
	"github.com/itcam/go-devops/api/app/controller/uic"
	"github.com/itcam/go-devops/api/app/utils"
)

func StartGin(port string, r *gin.Engine) {
	r.Use(utils.CORS())
	uic.Routes(r)
	openldap.Routes(r)
	r.Run(port)
}
