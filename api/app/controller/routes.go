package controller

import (
	"github.com/gin-gonic/gin"
	"go-devops/api/app/controller/uic"
	"go-devops/api/app/utils"
)

func StartGin(port string, r *gin.Engine) {
	r.Use(utils.CORS())
	uic.Routes(r)
	r.Run(port)
}
