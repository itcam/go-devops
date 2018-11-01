package openldap

import (
	"github.com/gin-gonic/gin"
	"github.com/itcam/go-devops/api/config"
	//"gopkg.in/openldap.v2"
)

var db config.DBPool

func Routes(r *gin.Engine) {
	db = config.Con()
	u := r.Group("/api/v1/ldap")
	u.POST("/connect", Connect)
	u.POST("/list", ListUser)
	u.POST("/DelByUid", DelUserByUid)
	u.POST("/userAdd", AddUser)
	u.POST("/findByUid", FindUserByUid)
	u.POST("/findByMail", FindUserByMail)
	u.POST("/findByMobile", FindUserByMobile)

}
