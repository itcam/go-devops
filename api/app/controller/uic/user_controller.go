package uic

import (
	"fmt"
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

type APIUserInput struct {
	UserName string `form:"username" json:"username" binding:"required"`
	FullName string `form:"fullname" json:"fullname"`
	PassWord string `form:"password" json:"password" binding:"required"`
	Email    string `form:"email" json:"email" binding:"required" valid:"email"`
	Phone    string `form:"phone" json:"phone" `
	Role     string `form:"role" json:"role" binding:"required"`
	Active   string `form:"active" json:"active" binding:"required" default:"T"`
}

func CreateUser(c *gin.Context) {
	var inputs APIUserInput
	err := c.ShouldBindJSON(&inputs)
	switch {
	case err != nil:
		h.JSONR(c, http.StatusBadRequest, "", err)
		log.Printf("有错误")
		return
	case utils.HasDangerousCharacters(inputs.UserName):
		log.Printf("名字格式不对")
		h.JSONR(c, http.StatusBadRequest, "", "name pattern is invalid")
		return

	}
	if _, err := govalidator.ValidateStruct(inputs); err != nil {
		h.JSONR(c, http.StatusBadRequest, "", err)
		return
	}
	user := new(uic.User)
	//如果不存在表，就建表
	if !db.Uic.HasTable(&user) {
		log.Debug("表不存在")
		db.Uic.CreateTable(&user)
	}

	//查询用户名是否已经存在

	db.Uic.Table(user.TableName()).Where("username = ?", inputs.UserName).Scan(&user)

	response := map[string]string{}
	response["id"] = strconv.Itoa(int(user.ID))
	response["username"] = user.Username

	if user.ID != 0 {
		h.JSONR(c, http.StatusBadRequest, response, "用户名已经存在")
		return
	}

	//查询邮箱是否已经存在
	db.Uic.Table(user.TableName()).Where("email = ?", inputs.Email).Scan(&user)
	for k := range response {
		delete(response, k)
	}
	response["id"] = strconv.Itoa(int(user.ID))
	response["email"] = user.Email

	if user.ID != 0 {
		h.JSONR(c, http.StatusBadRequest, response, "邮箱已经存在")
		return
	}

	//查询手机号是否已经存在
	db.Uic.Table(user.TableName()).Where("phone = ?", inputs.Phone).Scan(&user)
	for k := range response {
		delete(response, k)
	}
	response["id"] = strconv.Itoa(int(user.ID))
	response["phone"] = user.Phone
	if user.ID != 0 {
		h.JSONR(c, http.StatusBadRequest, response, "手机号已经存在")
		return
	}

	if inputs.Active == "F" {
		log.Printf("用户被禁用")
	}
	password := utils.HashIt(inputs.PassWord)

	user = &uic.User{
		Username: inputs.UserName,
		FullName: inputs.FullName,
		PassWord: password,
		Email:    inputs.Email,
		Phone:    inputs.Phone,
		Role:     inputs.Role,
		Active:   inputs.Active,
	}

	//开始创建用户
	dt := db.Uic.Table(user.TableName()).Create(&user)
	if dt.Error != nil {
		h.JSONR(c, http.StatusBadRequest, "", dt.Error)
		return
	}

	for k := range response {
		delete(response, k)
	}
	response["id"] = strconv.Itoa(int(user.ID))
	response["username"] = user.Username
	response["fullname"] = user.FullName
	response["email"] = user.Email
	response["phone"] = user.Phone
	h.JSONR(c, http.StatusOK, response, "用户创建成功")
	return

}

type APIUserUpdateFullNameInput struct {
	ID       int    `json:"id" binding:"required"`
	FullName string `form:"fullname" json:"fullname" binding:"required"`
}

func UpdateUserFullName(c *gin.Context) {
	var inputs APIUserUpdateFullNameInput
	err := c.BindJSON(&inputs)
	if err != nil {
		h.JSONR(c, http.StatusExpectationFailed, "", err)
		return
	}

	//if _, err := govalidator.ValidateStruct(inputs); err != nil {
	//	h.JSONR(c, http.StatusBadRequest, "", err)
	//	return
	//}

	user := new(uic.User)

	db.Uic.Table(user.TableName()).Where("id = ?", inputs.ID).Scan(&user)
	if user.ID == 0 {
		h.JSONR(c, http.StatusBadRequest, "", "用户ID不存在")
		return
	}

	uuser := map[string]interface{}{
		"fullname": inputs.FullName,
	}

	dt := db.Uic.Table(user.TableName()).Where("id = ?", inputs.ID).Update(uuser)
	if dt.Error != nil {
		h.JSONR(c, http.StatusExpectationFailed, "", dt.Error)
		return
	}

	response := map[string]string{}
	response["id"] = strconv.Itoa(int(user.ID))
	response["fullname"] = inputs.FullName
	h.JSONR(c, http.StatusOK, response, "姓名修改成功")
	return
}

type APIUserUpdateMailInput struct {
	ID    int    `json:"id" binding:"required"`
	Email string `form:"email" json:"email" binding:"required" valid:"email"`
}

func UpdateUserEmail(c *gin.Context) {
	var inputs APIUserUpdateMailInput
	err := c.BindJSON(&inputs)
	if err != nil {
		h.JSONR(c, http.StatusExpectationFailed, "", err)
		return
	}

	if _, err := govalidator.ValidateStruct(inputs); err != nil {
		h.JSONR(c, http.StatusBadRequest, "", err)
		return
	}

	user := new(uic.User)

	db.Uic.Table(user.TableName()).Where("id = ?", inputs.ID).Scan(&user)
	if user.ID == 0 {
		h.JSONR(c, http.StatusBadRequest, "", "用户ID不存在")
		return
	}

	uuser := map[string]interface{}{
		"Email": inputs.Email,
	}

	dt := db.Uic.Table(user.TableName()).Where("id = ?", inputs.ID).Update(uuser)
	if dt.Error != nil {
		h.JSONR(c, http.StatusExpectationFailed, "", dt.Error)
		return
	}

	response := map[string]string{}
	response["id"] = strconv.Itoa(int(user.ID))
	response["email"] = inputs.Email
	h.JSONR(c, http.StatusOK, response, "邮箱修改成功")
	return
}

type APIUserUpdatePhoneInput struct {
	ID    int    `json:"id" binding:"required"`
	Phone string `form:"phone" json:"phone" binding:"required"`
}

func UpdateUserPhone(c *gin.Context) {
	var inputs APIUserUpdatePhoneInput
	fmt.Println("id:", inputs.ID)
	fmt.Println("phone:", inputs.Phone)
	err := c.BindJSON(&inputs)
	if err != nil {
		h.JSONR(c, http.StatusExpectationFailed, "", err)
		return
	}

	if _, err := govalidator.ValidateStruct(inputs); err != nil {
		h.JSONR(c, http.StatusBadRequest, "", err)
		return
	}

	user := uic.User{}
	db.Uic.Table(user.TableName()).Where("id = ?", inputs.ID).Scan(&user)
	if user.ID == 0 {
		h.JSONR(c, http.StatusBadRequest, "", "用户ID不存在")
		return
	}
	uuser := map[string]interface{}{
		"Phone": inputs.Phone,
	}
	dt := db.Uic.Table(user.TableName()).Where("id = ?", inputs.ID).Update(uuser)
	if dt.Error != nil {
		h.JSONR(c, http.StatusExpectationFailed, "", dt.Error)
		return
	}
	response := map[string]string{}
	response["id"] = strconv.Itoa(int(user.ID))
	response["phone"] = inputs.Phone
	h.JSONR(c, http.StatusOK, response, "手机号修改成功")
	return
}

//修改密码

type APIUserUpdatePassInput struct {
	ID       int    `json:"id" binding:"required"`
	PassWord string `form:"password" json:"password" binding:"required"`
}

func UpdateUserPass(c *gin.Context) {
	var inputs APIUserUpdatePassInput
	err := c.BindJSON(&inputs)
	if err != nil {
		h.JSONR(c, http.StatusExpectationFailed, "", err)
		return
	}

	if _, err := govalidator.ValidateStruct(inputs); err != nil {
		h.JSONR(c, http.StatusBadRequest, "", err)
		return
	}

	user := uic.User{}
	db.Uic.Table(user.TableName()).Where("id = ?", inputs.ID).Scan(&user)
	if user.ID == 0 {
		h.JSONR(c, http.StatusBadRequest, "", "用户ID不存在")
		return
	}

	password := utils.HashIt(inputs.PassWord)
	uuser := map[string]interface{}{
		"password": password,
	}
	dt := db.Uic.Table(user.TableName()).Where("id = ?", inputs.ID).Update(uuser)
	if dt.Error != nil {
		h.JSONR(c, http.StatusExpectationFailed, "", dt.Error)
		return
	}
	response := map[string]string{}
	response["id"] = strconv.Itoa(int(user.ID))
	response["password"] = "*********"
	h.JSONR(c, http.StatusOK, response, "密码修改成功")
	return
}

//修改角色
type APIUserUpdateRoleInput struct {
	ID int `json:"id" binding:"required"`

	Role string `form:"role" json:"role" binding:"required"`
}

func UpdateUserRole(c *gin.Context) {
	var inputs APIUserUpdateRoleInput
	err := c.BindJSON(&inputs)
	if err != nil {
		h.JSONR(c, http.StatusExpectationFailed, "", err)
		return
	}

	if _, err := govalidator.ValidateStruct(inputs); err != nil {
		h.JSONR(c, http.StatusBadRequest, "", err)
		return
	}

	user := uic.User{}
	db.Uic.Table(user.TableName()).Where("id = ?", inputs.ID).Scan(&user)
	if user.ID == 0 {
		h.JSONR(c, http.StatusBadRequest, "", "用户ID不存在")
		return
	}

	fmt.Println("inputs", inputs)
	uuser := map[string]interface{}{
		"Role": inputs.Role,
	}

	dt := db.Uic.Table(user.TableName()).Where("id = ?", inputs.ID).Update(uuser)
	if dt.Error != nil {
		h.JSONR(c, http.StatusExpectationFailed, "", dt.Error)
		return
	}
	response := map[string]string{}
	response["id"] = strconv.Itoa(int(user.ID))
	response["role"] = inputs.Role
	h.JSONR(c, http.StatusOK, response, "角色修改成功")
	return
}

//列出所有用户

func Listuser(c *gin.Context) {

	type userlist struct {
		ID        uint `gorm:"primary_key"`
		CreatedAt time.Time
		UpdatedAt time.Time
		UserName  string `gorm:"type: varchar(64); not null; column:username"`
		FullName  string `gorm:"type: varchar(64); not null; column:fullname"`
		Email     string `gorm:"type: varchar(64); not null; column:email"`
		Phone     string `gorm:"type: varchar(64); not null; column:phone"`
		Role      string `gorm:"type: varchar(64); not null; column:role"`
		Active    string `gorm:"type: char(2); not null; column:Active"`
	}

	var result []userlist
	user := new(uic.User)
	response := make(map[string][]userlist)

	dt := db.Uic.Table(user.TableName()).Select("id,created_at,updated_at,username,fullname,email,phone,role,Active").Scan(&result)
	if dt.Error != nil {
		h.JSONR(c, http.StatusBadRequest, response, dt.Error)
		return
	}
	response["result"] = result

	h.JSONR(c, http.StatusOK, response, "查询完成")

	return
}

//查找用户by Id

func GetUserById(c *gin.Context) {

	type userlist struct {
		ID        uint `gorm:"primary_key"`
		CreatedAt time.Time
		UpdatedAt time.Time
		UserName  string `gorm:"type: varchar(64); not null; column:username"`
		FullName  string `gorm:"type: varchar(64); not null; column:fullname"`
		Email     string `gorm:"type: varchar(64); not null; column:email"`
		Phone     string `gorm:"type: varchar(64); not null; column:phone"`
		Role      string `gorm:"type: varchar(64); not null; column:role"`
		Active    string `gorm:"type: char(2); not null; column:Active"`
	}
	id := c.Param("id")
	user := new(uic.User)
	var result userlist
	response := make(map[string]userlist)
	dt := db.Uic.Table(user.TableName()).Select("id,created_at,updated_at,username,fullname,email,phone,role,Active").Where("id = ?", id).Scan(&result)
	if dt.Error != nil {
		h.JSONR(c, http.StatusBadRequest, response, dt.Error)
		return
	}
	response["result"] = result
	h.JSONR(c, http.StatusOK, response, "查询完成")
	return

}

//查找用户by UserName

func GetUserByUserName(c *gin.Context) {

	type userlist struct {
		ID        uint `gorm:"primary_key"`
		CreatedAt time.Time
		UpdatedAt time.Time
		UserName  string `gorm:"type: varchar(64); not null; column:username"`
		FullName  string `gorm:"type: varchar(64); not null; column:fullname"`
		Email     string `gorm:"type: varchar(64); not null; column:email"`
		Phone     string `gorm:"type: varchar(64); not null; column:phone"`
		Role      string `gorm:"type: varchar(64); not null; column:role"`
		Active    string `gorm:"type: char(2); not null; column:Active"`
	}
	username := c.Param("username")
	user := new(uic.User)
	var result userlist
	response := make(map[string]userlist)
	dt := db.Uic.Table(user.TableName()).Select("id,created_at,updated_at,username,fullname,email,phone,role,Active").Where("username = ?", username).Scan(&result)
	if dt.Error != nil {
		h.JSONR(c, http.StatusBadRequest, response, dt.Error)
		return
	}
	response["result"] = result
	h.JSONR(c, http.StatusOK, response, "查询完成")
	return
}

//查找用户by email

func GetUserByEmail(c *gin.Context) {

	type userlist struct {
		ID        uint `gorm:"primary_key"`
		CreatedAt time.Time
		UpdatedAt time.Time
		UserName  string `gorm:"type: varchar(64); not null; column:username"`
		FullName  string `gorm:"type: varchar(64); not null; column:fullname"`
		Email     string `gorm:"type: varchar(64); not null; column:email"`
		Phone     string `gorm:"type: varchar(64); not null; column:phone"`
		Role      string `gorm:"type: varchar(64); not null; column:role"`
		Active    string `gorm:"type: char(2); not null; column:Active"`
	}
	email := c.Param("email")
	user := new(uic.User)
	var result userlist
	response := make(map[string]userlist)
	dt := db.Uic.Table(user.TableName()).Select("id,created_at,updated_at,username,fullname,email,phone,role,Active").Where("email = ?", email).Scan(&result)
	if dt.Error != nil {
		h.JSONR(c, http.StatusBadRequest, response, dt.Error)
		return
	}
	response["result"] = result
	h.JSONR(c, http.StatusOK, response, "查询完成")
	return
}

//查找用户by phone

func GetUserByPhone(c *gin.Context) {

	type userlist struct {
		ID        uint `gorm:"primary_key"`
		CreatedAt time.Time
		UpdatedAt time.Time
		UserName  string `gorm:"type: varchar(64); not null; column:username"`
		FullName  string `gorm:"type: varchar(64); not null; column:fullname"`
		Email     string `gorm:"type: varchar(64); not null; column:email"`
		Phone     string `gorm:"type: varchar(64); not null; column:phone"`
		Role      string `gorm:"type: varchar(64); not null; column:role"`
		Active    string `gorm:"type: char(2); not null; column:Active"`
	}
	phone := c.Param("phone")
	user := new(uic.User)
	var result userlist
	response := make(map[string]userlist)
	dt := db.Uic.Table(user.TableName()).Select("id,created_at,updated_at,username,fullname,email,phone,role,Active").Where("phone = ?", phone).Scan(&result)
	if dt.Error != nil {
		h.JSONR(c, http.StatusBadRequest, response, dt.Error)
		return
	}
	response["result"] = result
	h.JSONR(c, http.StatusOK, response, "查询完成")
	return
}

//查找用户by role

func GetUserByRole(c *gin.Context) {

	type userlist struct {
		ID        uint `gorm:"primary_key"`
		CreatedAt time.Time
		UpdatedAt time.Time
		UserName  string `gorm:"type: varchar(64); not null; column:username"`
		FullName  string `gorm:"type: varchar(64); not null; column:fullname"`
		Email     string `gorm:"type: varchar(64); not null; column:email"`
		Phone     string `gorm:"type: varchar(64); not null; column:phone"`
		Role      string `gorm:"type: varchar(64); not null; column:role"`
		Active    string `gorm:"type: char(2); not null; column:Active"`
	}
	role := c.Param("role")
	user := new(uic.User)
	var result []userlist
	response := make(map[string][]userlist)
	dt := db.Uic.Table(user.TableName()).Select("id,created_at,updated_at,username,fullname,email,phone,role,Active").Where("role = ?", role).Scan(&result)
	if dt.Error != nil {
		h.JSONR(c, http.StatusBadRequest, response, dt.Error)
		return
	}
	response["result"] = result
	h.JSONR(c, http.StatusOK, response, "查询完成")
	return
}

//删除用户

func DelUser(c *gin.Context) {

	type userlist struct {
		ID        uint `gorm:"primary_key"`
		CreatedAt time.Time
		UpdatedAt time.Time
		UserName  string `gorm:"type: varchar(64); not null; column:username"`
		FullName  string `gorm:"type: varchar(64); not null; column:fullname"`
		Email     string `gorm:"type: varchar(64); not null; column:email"`
		Phone     string `gorm:"type: varchar(64); not null; column:phone"`
		Role      string `gorm:"type: varchar(64); not null; column:role"`
		Active    string `gorm:"type: char(2); not null; column:Active"`
	}
	id := c.Param("id")
	user := new(uic.User)
	var result userlist
	response := make(map[string]userlist)
	dt := db.Uic.Table(user.TableName()).Where("id = ?", id).Scan(&result)
	if dt.Error != nil {
		h.JSONR(c, http.StatusBadRequest, response, dt.Error)
		return
	}

	res := db.Uic.Table(user.TableName()).Where("id = ?", id).Delete(&result)
	if res.Error != nil {
		h.JSONR(c, http.StatusBadRequest, response, res.Error)
		return
	}
	response["result"] = result
	h.JSONR(c, http.StatusOK, response, "删除成功")
	return
}
