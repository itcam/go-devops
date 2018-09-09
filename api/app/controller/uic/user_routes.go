package uic

import (
	"github.com/gin-gonic/gin"
	"go-devops/api/config"
)

var db config.DBPool

func Routes(r *gin.Engine) {
	db = config.Con()
	u := r.Group("/api/v1/user")
	g := r.Group("/api/v1/group")

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

	g.POST("/create", CreateGroup)
	g.POST("/updategroupname", UpdateGroupName)
	g.Any("/list", ListUserGroup)
	g.GET("/getbyid/:id", GetGroupById)
	g.GET("/getbygroupname/:groupname", GetUserByUserGroupName)
	g.GET("/delgroup/:id", DelUserGroup)

}
