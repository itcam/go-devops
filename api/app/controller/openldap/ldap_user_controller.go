package openldap

import (
	"crypto/md5"
	"fmt"
	"gopkg.in/ldap.v2"
	"io"
	"math/rand"
	"os"
	"strconv"
	"time"
	//"github.com/sethvargo/go-password/password"
	log "github.com/Sirupsen/logrus"
	"github.com/gin-gonic/gin"
	h "github.com/itcam/go-devops/api/app/helper"
	"net/http"
)

var (
	ldapHost = "10.30.10.11"
	ldapPort = 389
	BaseDN   = "dc=bihu,dc=inet"
	UserDN   = "ou=people,dc=bihu,dc=inet"
	AdUser   = "cn=root,dc=bihu,dc=inet"
	AdPass   = ""
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

type UserAttr struct {
	CN        string `json:"cn" binding:"required"`
	Uid       string `json:"uid" binding:"required"`
	Mail      string `json:"mail" binding:"required"`
	UidNumber int    `json:"uidnumber" binding:"required"`
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

	searchRequest := ldap.NewSearchRequest(
		BaseDN,
		ldap.ScopeWholeSubtree, ldap.NeverDerefAliases, 0, 0, false,
		inputs.Filter,
		//"(&(objectClass=person))",
		[]string{"dn", "cn", "mail", "uid", "uidNumber"}, // 列出要查询的属性
		nil,
	)
	sr, err := l.Search(searchRequest)
	if err != nil {
		log.Fatal(err)
	}

	//var result []map[string]interface{}
	for _, entry := range sr.Entries {
		fmt.Printf("%s %v %s %v\n", entry.DN, entry.GetAttributeValue("cn"), entry.GetAttributeValue("mail"), entry.GetAttributeValue("uidNumber"))
	}
	//h.JSONR(c, http.StatusOK, inputs.UserCN, "绑定用户成功")

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

func DelUserByUid(l *ldap.Conn, BaseDN, filter string, attr []string, uid int) {
	searchRequest := ldap.NewSearchRequest(
		BaseDN, // The base dn to search
		ldap.ScopeWholeSubtree, ldap.NeverDerefAliases, 0, 0, false,
		filter,
		attr, // 列出要查询的属性
		nil,
	)
	sr, err := l.Search(searchRequest)
	if err != nil {
		log.Fatal(err)
	}
	var matchUser []int
	for _, entry := range sr.Entries {
		id, err := strconv.Atoi(entry.GetAttributeValue("uidNumber"))
		if err != nil {
			HandleError(err)
		}
		if id == uid {
			matchUser = append(matchUser, id)
			log.Info(fmt.Sprintf("找到用户%v,开始执行删除...", entry.DN))
			delete := ldap.NewDelRequest(entry.DN, nil)
			err = l.Del(delete)
			if err != nil {
				HandleError(err)
				log.Info(fmt.Sprintf("删除用户%v失败", entry.DN))
			} else {
				log.Info(fmt.Sprintf("删除用户%v成功", entry.DN))
			}
		}
	}
	if len(matchUser) == 0 {
		log.Println(fmt.Sprintf("没有找到uidNumber 为 %v 的用户", uid))
	} else {
		log.Println(fmt.Sprintf("总共找到%d个用户", len(matchUser)))
	}

}

func FindMaxUid(l *ldap.Conn, BaseDN, filter string, attr []string) int {
	searchRequest := ldap.NewSearchRequest(
		BaseDN, // The base dn to search
		ldap.ScopeWholeSubtree, ldap.NeverDerefAliases, 0, 0, false,
		filter,
		attr, // 列出要查询的属性
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
			HandleError(err)
		}
		s = append(s, uidNumber)
	}

	maxUid := SliceMax(s)
	return maxUid

}

func AddUser(l *ldap.Conn, uid int, BaseDN, filter, username, email, phone, pass string, attr []string) {
	searchRequest := ldap.NewSearchRequest(
		BaseDN, // The base dn to search
		ldap.ScopeWholeSubtree, ldap.NeverDerefAliases, 0, 0, false,
		filter,
		[]string{"uidNumber", "mail", "uid", "cn", "phone"}, // 列出要查询的属性
		nil,
	)
	sr, err := l.Search(searchRequest)
	if err != nil {
		log.Fatal(err)
	}
	for _, entry := range sr.Entries {
		cn := entry.GetAttributeValue("cn")
		uid := entry.GetAttributeValue("uid")
		mail := entry.GetAttributeValue("mail")

		if cn == username && uid == username {
			log.Fatal(fmt.Sprintf("用户名已经存在,该用户是%v", entry.DN))
			HandleError(err)
		}
		if mail == email {
			log.Fatal(fmt.Sprintf("用户邮箱已经存在,该用户是%v", entry.DN))
			HandleError(err)
		}
	}
	user := fmt.Sprintf("cn=%v,ou=people,dc=bihu,dc=inet", username)
	homeDirectory := "/tmp"
	log.Info(user)
	add := ldap.NewAddRequest(user)
	add.Attribute("objectClass", []string{"top", "shadowAccount", "posixAccount", "person", "organizationalPerson", "inetOrgPerson", "inetLocalMailRecipient", "hostObject", "openVPNUser"})
	add.Attribute("uid", []string{username})
	add.Attribute("uidNumber", []string{strconv.Itoa(uid)})
	add.Attribute("gidNumber", []string{strconv.Itoa(uid)})
	add.Attribute("homeDirectory", []string{homeDirectory})
	add.Attribute("sn", []string{username})
	add.Attribute("dockerRegistryActive", []string{"FALSE"})
	add.Attribute("openvpnbhActive", []string{"FALSE"})
	add.Attribute("description", []string{username})
	add.Attribute("mail", []string{email})
	add.Attribute("mobile", []string{phone})
	add.Attribute("loginShell", []string{"/bin/false"})
	add.Attribute("shadowExpire", []string{"-1"})
	add.Attribute("shadowFlag", []string{"0"})
	add.Attribute("userPassword", []string{pass})

	err = l.Add(add)
	if err != nil {
		HandleError(err)
		log.Error("添加用户失败", user)
	}

}

func randSeq(n int) string {
	b := make([]rune, n)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b)
}

func createPasswd() string {
	t := time.Now()
	h := md5.New()
	io.WriteString(h, "crazyo")
	io.WriteString(h, t.String())
	passwd := fmt.Sprintf("%x", h.Sum(nil))
	return passwd
}

//func main() {
//
//	log.SetFormatter(&log.TextFormatter{})
//	file, err := os.OpenFile(logFile, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
//	if err != nil {
//		log.Fatal(err.Error())
//	}
//	defer file.Close()
//	log.SetOutput(file)
//	log.SetLevel(log.DebugLevel)
//
//	l, err := openldap.Dial("tcp", fmt.Sprintf("%s:%d", ldapHost, ldapPort))
//	if err != nil {
//		log.Fatal(err)
//	}
//	defer l.Close()
//
//	testConnect(l, "cn=arrow,ou=people,dc=bihu,dc=inet", "123456")
//	//列出所有用户
//	ListUser(l, BaseDN, "(&(objectClass=person))", []string{"dn", "cn", "mail", "uid", "uidNumber"})
//	err = BindAdUser(l, AdUser, AdPass)
//	if err != nil {
//		log.Fatal("绑定管理员失败")
//	}
//	//DelUserByUid(l, 10099)
//
//	maxId := FindMaxUid(l, BaseDN, "(&(objectClass=person))", []string{"dn", "cn", "mail", "uid", "uidNumber"})
//
//}c
