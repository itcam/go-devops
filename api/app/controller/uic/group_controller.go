package uic

import (
	log "github.com/Sirupsen/logrus"
	"github.com/asaskevich/govalidator"
	"github.com/gin-gonic/gin"
	h "github.com/itcam/go-devops/api/app/helper"
	"github.com/itcam/go-devops/api/app/model/uic"
	"github.com/itcam/go-devops/api/app/utils"
	"net/http"
	"strconv"
	"time"
)

type APIGroupInput struct {
	GroupName string `form:"groupname" json:"groupname" binding:"required"`
}

func CreateGroup(c *gin.Context) {
	var inputs APIGroupInput
	err := c.ShouldBindJSON(&inputs)
	switch {
	case err != nil:
		h.JSONR(c, http.StatusBadRequest, "", err)
		log.Printf("有错误")
		return
	case utils.HasDangerousCharacters(inputs.GroupName):
		log.Printf("名字格式不对")
		h.JSONR(c, http.StatusBadRequest, "", "name pattern is invalid")
		return

	}
	if _, err := govalidator.ValidateStruct(inputs); err != nil {
		h.JSONR(c, http.StatusBadRequest, "", err)
		return
	}
	group := new(uic.Group)
	//如果不存在表，就建表
	if !db.Uic.HasTable(&group) {
		log.Debug("表不存在")
		db.Uic.CreateTable(&group)
	}

	//查询用户名是否已经存在

	db.Uic.Table(group.TableName()).Where("groupname = ?", inputs.GroupName).Scan(&group)

	response := map[string]string{}
	response["id"] = strconv.Itoa(int(group.ID))
	response["groupname"] = group.GroupName

	if group.ID != 0 {
		h.JSONR(c, http.StatusBadRequest, response, "组名已经存在")
		return
	}

	//开始创建用户组
	dt := db.Uic.Table(group.TableName()).Create(&group)
	if dt.Error != nil {
		h.JSONR(c, http.StatusBadRequest, "", dt.Error)
		return
	}

	for k := range response {
		delete(response, k)
	}
	response["id"] = strconv.Itoa(int(group.ID))
	response["groupname"] = group.GroupName
	h.JSONR(c, http.StatusOK, response, "用户组创建成功")
	return

}

type APIGroupUpdateGroupNameInput struct {
	ID        int    `json:"id" binding:"required"`
	GroupName string `form:"groupname" json:"groupname" binding:"required"`
}

func UpdateGroupName(c *gin.Context) {
	var inputs APIGroupUpdateGroupNameInput
	err := c.BindJSON(&inputs)
	if err != nil {
		h.JSONR(c, http.StatusExpectationFailed, "", err)
		return
	}

	group := new(uic.Group)

	db.Uic.Table(group.TableName()).Where("id = ?", inputs.ID).Scan(&group)
	if group.ID == 0 {
		h.JSONR(c, http.StatusBadRequest, "", "用户组ID不存在")
		return
	}

	ggroup := map[string]interface{}{
		"groupname": inputs.GroupName,
	}

	dt := db.Uic.Table(group.TableName()).Where("id = ?", inputs.ID).Update(ggroup)
	if dt.Error != nil {
		h.JSONR(c, http.StatusExpectationFailed, "", dt.Error)
		return
	}

	response := map[string]string{}
	response["id"] = strconv.Itoa(int(group.ID))
	response["groupname"] = inputs.GroupName
	h.JSONR(c, http.StatusOK, response, "组名修改成功")
	return
}

//列出所有用户组

func ListUserGroup(c *gin.Context) {

	type grouplist struct {
		ID        uint `gorm:"primary_key"`
		CreatedAt time.Time
		UpdatedAt time.Time
		GroupName string `gorm:"type: varchar(64); not null; column:groupname"`
	}

	var result []grouplist
	group := new(uic.Group)
	response := make(map[string][]grouplist)

	dt := db.Uic.Table(group.TableName()).Select("id,created_at,updated_at,groupname").Scan(&result)
	if dt.Error != nil {
		h.JSONR(c, http.StatusBadRequest, response, dt.Error)
		return
	}
	response["result"] = result

	h.JSONR(c, http.StatusOK, response, "查询完成")

	return
}

//查找组by Id

func GetGroupById(c *gin.Context) {

	type grouplist struct {
		ID        uint `gorm:"primary_key"`
		CreatedAt time.Time
		UpdatedAt time.Time
		GroupName string `gorm:"type: varchar(64); not null; column:groupname"`
	}
	id := c.Param("id")
	group := new(uic.Group)
	var result grouplist
	response := make(map[string]grouplist)
	dt := db.Uic.Table(group.TableName()).Select("id,created_at,updated_at,groupname").Where("id = ?", id).Scan(&result)
	if dt.Error != nil {
		h.JSONR(c, http.StatusBadRequest, response, dt.Error)
		return
	}
	response["result"] = result
	h.JSONR(c, http.StatusOK, response, "查询完成")
	return

}

//查找用户组by GroupName

func GetUserByUserGroupName(c *gin.Context) {

	type grouplist struct {
		ID        uint `gorm:"primary_key"`
		CreatedAt time.Time
		UpdatedAt time.Time
		GroupName string `gorm:"type: varchar(64); not null; column:groupname"`
	}
	groupname := c.Param("groupname")
	group := new(uic.Group)
	var result grouplist
	response := make(map[string]grouplist)
	dt := db.Uic.Table(group.TableName()).Select("id,created_at,updated_at,groupname").Where("groupname = ?", groupname).Scan(&result)
	if dt.Error != nil {
		h.JSONR(c, http.StatusBadRequest, response, dt.Error)
		return
	}
	response["result"] = result
	h.JSONR(c, http.StatusOK, response, "查询完成")
	return
}

//删除组

func DelUserGroup(c *gin.Context) {

	type grouplist struct {
		ID        uint `gorm:"primary_key"`
		CreatedAt time.Time
		UpdatedAt time.Time
		GroupName string `gorm:"type: varchar(64); not null; column:groupname"`
	}
	id := c.Param("id")
	group := new(uic.Group)
	var result grouplist
	response := make(map[string]grouplist)
	dt := db.Uic.Table(group.TableName()).Where("id = ?", id).Scan(&result)
	if dt.Error != nil {
		h.JSONR(c, http.StatusBadRequest, response, dt.Error)
		return
	}

	res := db.Uic.Table(group.TableName()).Where("id = ?", id).Delete(&result)
	if res.Error != nil {
		h.JSONR(c, http.StatusBadRequest, response, res.Error)
		return
	}
	response["result"] = result
	h.JSONR(c, http.StatusOK, response, "删除成功")
	return
}
