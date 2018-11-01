package openldap

import (
	"fmt"
	"gopkg.in/ldap.v2"
	"os"
	"strconv"
	//"github.com/sethvargo/go-password/password"
	log "github.com/Sirupsen/logrus"
	"github.com/gin-gonic/gin"
	h "github.com/itcam/go-devops/api/app/helper"
	"github.com/itcam/go-devops/api/app/utils"
	"net/http"
	"strings"
)

var (
	ldapHost = "10.30.10.11"
	ldapPort = 389
	BaseDN   = "dc=bihu,dc=inet"
	UserDN   = "ou=people,dc=bihu,dc=inet"
	AdUser   = "cn=root,dc=bihu,dc=inet"
	AdPass   = "B1hu12345"
	logFile  = "goldap.log"
	letters  = []rune("absdfdsfsdfsdf2dsfsdfXXDDsdfsdf")
	//logger   = log.New()
)

//type Fields log.Fields

func SliceMax(s []int) int {
	max := 0
	for _, i := range s {
		if max < i {
			max = i
		}
	}
	return max

}

func HandleError(err error) {
	log.Error("error: ", err)
	os.Exit(-1)
}

type EnTryInput struct {
	UserCN   string `json:"cn" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type ListInput struct {
	Filter string `json:"filter" binding:"required"`
}

type DelInput struct {
	Filter string `json:"filter" binding:"required"`
	Uid    string `json:"uid" binding:"required"`
}

type UserAttr struct {
	CN                   string `json:"cn" `
	Uid                  string `json:"uid" binding:"required"`
	Mail                 string `json:"mail" binding:"required"`
	UidNumber            int    `json:"uidnumber"`
	Mobile               string `json:"mobile" binding:"required"`
	ShadowExpire         string `json:"shadowexpire"`
	DockerRegistryActive string `json:"dockerRegistryActive"`
	OpenvpnbhActive      string `json:"openvpnbhActive"`
	UserPassword         string `json:"password"`
}

type FindByUidInput struct {
	Uid string `json:"uid" binding:"required"`
}
type FindByMailInput struct {
	Mail string `json:"mail" binding:"required"`
}
type FindByMobileInput struct {
	Mobile string `json:"mobile" binding:"required"`
}

func Connect(c *gin.Context) {
	l, err := ldap.Dial("tcp", fmt.Sprintf("%s:%d", ldapHost, ldapPort))
	if err != nil {
		h.JSONR(c, http.StatusExpectationFailed, "连接LDAP服务失败", err)
		return
	}

	var inputs EnTryInput
	err = c.BindJSON(&inputs)
	if err != nil {
		h.JSONR(c, http.StatusExpectationFailed, "", err)
		return
	}

	username := fmt.Sprintf("cn=%s,%s", inputs.UserCN, UserDN)
	err = l.Bind(username, inputs.Password)
	if err != nil {
		log.Error(err)
		h.JSONR(c, http.StatusExpectationFailed, "绑定用户失败", err)
		return
	}

	h.JSONR(c, http.StatusOK, inputs.UserCN, "绑定用户成功")
	log.Info("绑定用户成功: ", username)
}

func ListUser(c *gin.Context) {
	l, err := ldap.Dial("tcp", fmt.Sprintf("%s:%d", ldapHost, ldapPort))
	if err != nil {
		h.JSONR(c, http.StatusExpectationFailed, "连接LDAP服务失败", err)
		return
	}
	var inputs ListInput
	err = c.BindJSON(&inputs)
	if err != nil {
		h.JSONR(c, http.StatusExpectationFailed, "", err)
		return
	}

	userAttrList := []string{"dn", "cn", "mail", "uid", "uidNumber", "Mobile", "shadowExpire", "dockerRegistryActive", "openvpnbhActive"}
	searchRequest := ldap.NewSearchRequest(
		BaseDN,
		ldap.ScopeWholeSubtree, ldap.NeverDerefAliases, 0, 0, false,
		inputs.Filter,
		//"(&(objectClass=person))",
		userAttrList,
		nil,
	)
	sr, err := l.Search(searchRequest)
	if err != nil {
		log.Fatal(err)
		h.JSONR(c, http.StatusExpectationFailed, "", err)
	}

	var user UserAttr
	var result []UserAttr
	for _, entry := range sr.Entries {
		user.CN = entry.GetAttributeValue("cn")
		user.Mail = entry.GetAttributeValue("mail")
		user.UidNumber, _ = strconv.Atoi(entry.GetAttributeValue("uidNumber"))
		user.Uid = entry.GetAttributeValue("uid")
		user.Mobile = entry.GetAttributeValue("mobile")
		user.ShadowExpire = entry.GetAttributeValue("shadowExpire")
		user.DockerRegistryActive = entry.GetAttributeValue("dockerRegistryActive")
		user.OpenvpnbhActive = entry.GetAttributeValue("openvpnbhActive")
		result = append(result, user)
	}
	h.JSONR(c, http.StatusOK, result, "用户信息")

}

func BindAdUser(l *ldap.Conn, user, password string) error {
	err := l.Bind(user, password)
	if err != nil {
		log.Error(err)
		return err
	}
	log.Info("绑定管理员成功")
	return nil
}

func DelUserByUid(c *gin.Context) {
	l, err := ldap.Dial("tcp", fmt.Sprintf("%s:%d", ldapHost, ldapPort))
	if err != nil {
		log.Error(err)
		h.JSONR(c, http.StatusExpectationFailed, "连接LDAP服务失败", err)
		return
	}

	err = BindAdUser(l, AdUser, AdPass)
	if err != nil {
		log.Error(err)
		h.JSONR(c, http.StatusExpectationFailed, "绑定管理员失败", err)
		return
	}

	var inputs DelInput
	err = c.BindJSON(&inputs)
	if err != nil {
		log.Error(err)
		h.JSONR(c, http.StatusExpectationFailed, "", err)
		return
	}
	userAttrList := []string{"cn", "mail", "uid", "uidNumber"}

	searchRequest := ldap.NewSearchRequest(
		BaseDN, // The base dn to search
		ldap.ScopeWholeSubtree, ldap.NeverDerefAliases, 0, 0, false,
		inputs.Filter,
		userAttrList, // 列出要查询的属性
		nil,
	)
	sr, err := l.Search(searchRequest)
	if err != nil {
		log.Error(err)
		h.JSONR(c, http.StatusExpectationFailed, "", err)
		return
	}
	var matchUser []string
	for _, entry := range sr.Entries {
		if inputs.Uid == entry.GetAttributeValue("uid") {
			matchUser = append(matchUser, inputs.Uid)
			log.Info(fmt.Sprintf("找到用户%v,开始执行删除...", entry.DN))
			delete := ldap.NewDelRequest(entry.DN, nil)
			err = l.Del(delete)
			if err != nil {
				log.Error(fmt.Sprintf("删除用户%v失败", entry.DN))
				log.Error(err)
				h.JSONR(c, http.StatusExpectationFailed, "", err)
				return
			} else {
				log.Info(fmt.Sprintf("删除用户%v成功", entry.DN))
			}
		}
	}
	if len(matchUser) == 0 {
		log.Info(fmt.Sprintf("没有找到uid 为 %v 的用户", inputs.Uid))
		h.JSONR(c, http.StatusNotFound, "uid="+inputs.Uid, "没有找到用户")
		return
	} else {
		log.Info(fmt.Sprintf("总共找到%d个用户", len(matchUser)))
		h.JSONR(c, http.StatusOK, matchUser, "删除用户成功")
		return

	}

}

func FindMaxUid(l *ldap.Conn) int {
	//l, err := ldap.Dial("tcp", fmt.Sprintf("%s:%d", ldapHost, ldapPort))
	//if err != nil {
	//	log.Fatal(err)
	//}

	userAttrList := []string{"cn", "mail", "uid", "uidNumber"}

	searchRequest := ldap.NewSearchRequest(
		BaseDN, // The base dn to search
		ldap.ScopeWholeSubtree, ldap.NeverDerefAliases, 0, 0, false,
		"(&(objectClass=person))",
		userAttrList, // 列出要查询的属性
		nil,
	)
	sr, err := l.Search(searchRequest)
	if err != nil {
		log.Fatal(err)
	}

	var s []int
	for _, entry := range sr.Entries {
		uidNumber, err := strconv.Atoi(entry.GetAttributeValue("uidNumber"))
		if err != nil {
			log.Error(err)
		}
		s = append(s, uidNumber)
	}

	maxUid := SliceMax(s)
	return maxUid
}

func AddUser(c *gin.Context) {
	l, err := ldap.Dial("tcp", fmt.Sprintf("%s:%d", ldapHost, ldapPort))
	if err != nil {
		log.Error(err)
		h.JSONR(c, http.StatusExpectationFailed, "连接LDAP服务失败", err)
		return
	}

	err = BindAdUser(l, AdUser, AdPass)
	if err != nil {
		log.Error(err)
		h.JSONR(c, http.StatusExpectationFailed, "绑定管理员失败", err)
		return
	}

	var inputs UserAttr
	err = c.BindJSON(&inputs)
	if err != nil {
		h.JSONR(c, http.StatusExpectationFailed, "", err)
		return
	}
	userAttrList := []string{"dn", "cn", "mail", "uid", "uidNumber", "Mobile", "shadowExpire", "dockerRegistryActive", "openvpnbhActive"}

	searchRequest := ldap.NewSearchRequest(
		BaseDN, // The base dn to search
		ldap.ScopeWholeSubtree, ldap.NeverDerefAliases, 0, 0, false,
		"(&(objectClass=person))",
		userAttrList, // 列出要查询的属性
		nil,
	)
	sr, err := l.Search(searchRequest)
	if err != nil {
		log.Error(err)
		h.JSONR(c, http.StatusExpectationFailed, "搜索失败", err)
		return
	}
	for _, entry := range sr.Entries {
		cn := entry.GetAttributeValue("cn")
		uid := entry.GetAttributeValue("uid")
		mail := entry.GetAttributeValue("mail")

		if cn == inputs.Uid && uid == inputs.Uid {
			log.Error(fmt.Sprintf("用户名已经存在,该用户是%v", entry.DN))
			h.JSONR(c, http.StatusExpectationFailed, fmt.Sprintf("用户名已经存在,该用户是%v", entry.DN), err)
			return
		}
		if mail == inputs.Mail {
			log.Error(fmt.Sprintf("用户邮箱已经存在,该用户是%v", entry.DN))
			h.JSONR(c, http.StatusExpectationFailed, fmt.Sprintf("用户邮箱已经存在,该用户是%v", entry.DN), err)
			return
		}
	}
	maxUid := FindMaxUid(l)
	user := fmt.Sprintf("cn=%v,%v", inputs.Uid, UserDN)
	homeDirectory := "/tmp"
	log.Info(user)

	var userPass string
	if inputs.UserPassword == "" {
		log.Info("请求body没有password参数，生成随机密码")
		userPass = utils.GeneratePass(15)
	} else {
		userPass = inputs.UserPassword
	}
	add := ldap.NewAddRequest(user)
	add.Attribute("objectClass", []string{"top", "shadowAccount", "posixAccount", "person", "organizationalPerson", "inetOrgPerson", "inetLocalMailRecipient", "hostObject", "openVPNUser"})
	add.Attribute("uid", []string{inputs.Uid})
	add.Attribute("uidNumber", []string{strconv.Itoa(maxUid + 1)})
	add.Attribute("gidNumber", []string{strconv.Itoa(11111)})
	add.Attribute("homeDirectory", []string{homeDirectory})
	add.Attribute("sn", []string{inputs.Uid})
	add.Attribute("dockerRegistryActive", []string{"FALSE"})
	add.Attribute("openvpnbhActive", []string{"FALSE"})
	add.Attribute("description", []string{inputs.Uid})
	add.Attribute("mail", []string{inputs.Mail})
	add.Attribute("mobile", []string{inputs.Mobile})
	add.Attribute("loginShell", []string{"/bin/false"})
	add.Attribute("shadowExpire", []string{"-1"})
	add.Attribute("shadowFlag", []string{"0"})

	add.Attribute("userPassword", []string{userPass})

	result := []string{
		fmt.Sprintf("用户名: %s", inputs.Uid),
		fmt.Sprintf("邮箱: %s", inputs.Mail),
		fmt.Sprintf("手机号: %s", inputs.Mobile),
		fmt.Sprintf("密码: %s", userPass),
	}

	err = l.Add(add)
	if err != nil {
		log.Error(err)
		h.JSONR(c, http.StatusExpectationFailed, "添加用户失败", err)
		return
	} else {
		h.JSONR(c, http.StatusOK, user, "添加用户成功,请查收邮件")
		utils.SendEmail("smtp.exmail.qq.com:587", "git@gittab.com", "KZQbZ7iM3v1txauP", "LDAP通知", "LDAP用户创建成功", strings.Join(result, "\n"), []string{inputs.Mail})
		return
	}
}

func FindUserByUid(c *gin.Context) {
	l, err := ldap.Dial("tcp", fmt.Sprintf("%s:%d", ldapHost, ldapPort))
	if err != nil {
		log.Error(err)
		h.JSONR(c, http.StatusExpectationFailed, "连接LDAP服务失败", err)
		return
	}
	err = BindAdUser(l, AdUser, AdPass)
	if err != nil {
		log.Error(err)
		h.JSONR(c, http.StatusExpectationFailed, "绑定管理员失败", err)
		return
	}

	var inputs FindByUidInput
	err = c.BindJSON(&inputs)
	if err != nil {
		h.JSONR(c, http.StatusExpectationFailed, "", err)
		return
	}
	userAttrList := []string{"dn", "cn", "mail", "uid", "uidNumber", "Mobile", "shadowExpire", "dockerRegistryActive", "openvpnbhActive"}

	searchRequest := ldap.NewSearchRequest(
		BaseDN, // The base dn to search
		ldap.ScopeWholeSubtree, ldap.NeverDerefAliases, 0, 0, false,
		"(&(objectClass=person))",
		userAttrList, // 列出要查询的属性
		nil,
	)
	sr, err := l.Search(searchRequest)
	if err != nil {
		log.Error(err)
		h.JSONR(c, http.StatusExpectationFailed, "搜索失败", err)
		return
	}
	var user UserAttr
	var result []UserAttr
	fmt.Println(inputs.Uid)
	for _, entry := range sr.Entries {

		if inputs.Uid == entry.GetAttributeValue("uid") {
			user.CN = entry.GetAttributeValue("cn")
			user.Mail = entry.GetAttributeValue("mail")
			user.UidNumber, _ = strconv.Atoi(entry.GetAttributeValue("uidNumber"))
			user.Uid = entry.GetAttributeValue("uid")
			user.Mobile = entry.GetAttributeValue("mobile")
			user.ShadowExpire = entry.GetAttributeValue("shadowExpire")
			user.DockerRegistryActive = entry.GetAttributeValue("dockerRegistryActive")
			user.OpenvpnbhActive = entry.GetAttributeValue("openvpnbhActive")
			log.Error(fmt.Sprintf("用户名存在,该用户是%v", entry.DN))
			result = append(result, user)

		}

	}
	if len(result) == 0 {
		log.Info(fmt.Sprintf("没有找到uid 为 %v 的用户", inputs.Uid))
		h.JSONR(c, http.StatusNotFound, "", "没有找到用户")
		return
	} else {
		log.Info(fmt.Sprintf("总共找到%d个用户", len(result)))
		h.JSONR(c, http.StatusOK, result, "找到"+strconv.Itoa(len(result))+"个用户")
		return
	}

}

func FindUserByMail(c *gin.Context) {
	l, err := ldap.Dial("tcp", fmt.Sprintf("%s:%d", ldapHost, ldapPort))
	if err != nil {
		log.Error(err)
		h.JSONR(c, http.StatusExpectationFailed, "连接LDAP服务失败", err)
		return
	}
	err = BindAdUser(l, AdUser, AdPass)
	if err != nil {
		log.Error(err)
		h.JSONR(c, http.StatusExpectationFailed, "绑定管理员失败", err)
		return
	}

	var inputs FindByMailInput
	err = c.BindJSON(&inputs)
	if err != nil {
		h.JSONR(c, http.StatusExpectationFailed, "", err)
		return
	}
	userAttrList := []string{"dn", "cn", "mail", "uid", "uidNumber", "Mobile", "shadowExpire", "dockerRegistryActive", "openvpnbhActive"}

	searchRequest := ldap.NewSearchRequest(
		BaseDN, // The base dn to search
		ldap.ScopeWholeSubtree, ldap.NeverDerefAliases, 0, 0, false,
		"(&(objectClass=person))",
		userAttrList, // 列出要查询的属性
		nil,
	)
	sr, err := l.Search(searchRequest)
	if err != nil {
		log.Error(err)
		h.JSONR(c, http.StatusExpectationFailed, "搜索失败", err)
		return
	}
	var user UserAttr
	var result []UserAttr
	for _, entry := range sr.Entries {

		if inputs.Mail == entry.GetAttributeValue("mail") {
			user.CN = entry.GetAttributeValue("cn")
			user.Mail = entry.GetAttributeValue("mail")
			user.UidNumber, _ = strconv.Atoi(entry.GetAttributeValue("uidNumber"))
			user.Uid = entry.GetAttributeValue("uid")
			user.Mobile = entry.GetAttributeValue("mobile")
			user.ShadowExpire = entry.GetAttributeValue("shadowExpire")
			user.DockerRegistryActive = entry.GetAttributeValue("dockerRegistryActive")
			user.OpenvpnbhActive = entry.GetAttributeValue("openvpnbhActive")
			log.Error(fmt.Sprintf("用户名存在,该用户是%v", entry.DN))
			result = append(result, user)

		}

	}
	if len(result) == 0 {
		log.Info(fmt.Sprintf("没有找到mail 为 %v 的用户", inputs.Mail))
		h.JSONR(c, http.StatusNotFound, "", "没有找到用户")
		return
	} else {
		log.Info(fmt.Sprintf("总共找到%d个用户", len(result)))
		h.JSONR(c, http.StatusOK, result, "找到"+strconv.Itoa(len(result))+"个用户")
		return
	}

}

func FindUserByMobile(c *gin.Context) {
	l, err := ldap.Dial("tcp", fmt.Sprintf("%s:%d", ldapHost, ldapPort))
	if err != nil {
		log.Error(err)
		h.JSONR(c, http.StatusExpectationFailed, "连接LDAP服务失败", err)
		return
	}
	err = BindAdUser(l, AdUser, AdPass)
	if err != nil {
		log.Error(err)
		h.JSONR(c, http.StatusExpectationFailed, "绑定管理员失败", err)
		return
	}

	var inputs FindByMobileInput
	err = c.BindJSON(&inputs)
	if err != nil {
		h.JSONR(c, http.StatusExpectationFailed, "", err)
		return
	}
	userAttrList := []string{"dn", "cn", "mail", "uid", "uidNumber", "Mobile", "shadowExpire", "dockerRegistryActive", "openvpnbhActive"}

	searchRequest := ldap.NewSearchRequest(
		BaseDN, // The base dn to search
		ldap.ScopeWholeSubtree, ldap.NeverDerefAliases, 0, 0, false,
		"(&(objectClass=person))",
		userAttrList, // 列出要查询的属性
		nil,
	)
	sr, err := l.Search(searchRequest)
	if err != nil {
		log.Error(err)
		h.JSONR(c, http.StatusExpectationFailed, "搜索失败", err)
		return
	}
	var user UserAttr
	var result []UserAttr
	for _, entry := range sr.Entries {

		if inputs.Mobile == entry.GetAttributeValue("mobile") {
			user.CN = entry.GetAttributeValue("cn")
			user.Mail = entry.GetAttributeValue("mail")
			user.UidNumber, _ = strconv.Atoi(entry.GetAttributeValue("uidNumber"))
			user.Uid = entry.GetAttributeValue("uid")
			user.Mobile = entry.GetAttributeValue("mobile")
			user.ShadowExpire = entry.GetAttributeValue("shadowExpire")
			user.DockerRegistryActive = entry.GetAttributeValue("dockerRegistryActive")
			user.OpenvpnbhActive = entry.GetAttributeValue("openvpnbhActive")
			log.Error(fmt.Sprintf("用户名存在,该用户是%v", entry.DN))
			result = append(result, user)

		}

	}
	if len(result) == 0 {
		log.Info(fmt.Sprintf("没有找到mobile 为 %v 的用户", inputs.Mobile))
		h.JSONR(c, http.StatusNotFound, "", "没有找到用户")
		return
	} else {
		log.Info(fmt.Sprintf("总共找到%d个用户", len(result)))
		h.JSONR(c, http.StatusOK, result, "找到"+strconv.Itoa(len(result))+"个用户")
		return
	}

}
