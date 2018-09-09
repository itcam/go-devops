package uic

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/itcam/go-devops/api/config"
)

var db config.DBPool

const badstatus = http.StatusBadRequest

func Routes(r *gin.Engine) {
	db = config.Con()
	//session
	u := r.Group("/api/v1/user")

	//user route
	u.POST("/create", CreateUser)
	u.POST("/updatefullname", UpdateUserFullName)
	u.POST("/updatemail", UpdateUserEmail)
	u.POST("/updatephone", UpdateUserPhone)
	u.POST("/updatepass", UpdateUserPass)
	u.POST("/updaterole", UpdateUserRole)
	u.GET("/list", Listuser)
	u.GET("/getbyid/:id", GetUserById)
	u.GET("/getbyusername/:username", GetUserByUserName)
	u.GET("/getbyemail/:email", GetUserByEmail)
	u.GET("/getbyphone/:phone", GetUserByPhone)
	u.GET("/getbyrole/:role", GetUserByRole)
	u.GET("/deluser/:id", DelUser)

	//group route
	g := r.Group("/api/v1/group")
	g.POST("/create", CreateGroup)
	g.POST("/updategroupname", UpdateGroupName)
	g.POST("/list", ListUserGroup)
	u.GET("/getbyid/:id", GetGroupById)
	u.GET("/getbygroupname/:groupname", GetUserByUserGroupName)
	u.GET("/delgroup/:id", DelUserGroup)

}
