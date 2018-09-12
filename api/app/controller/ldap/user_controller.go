package ldap

import "C"
import (
	"fmt"
	"gopkg.in/ldap.v2"
	"log"
	"strconv"
)

type Conn struct {
	Host string
	Port int
}

func (c *Conn) setConn(host string, port int) {
	c.Host = host
	c.Port = port
}

func (c *Conn) getConnHost() string {
	return c.Host
}

func (c *Conn) getConnPort() int {
	return c.Port
}

func (c *Conn) getConn() (string, int) {
	return c.Host, c.Port
}

func SliceMax(s []int) int {
	max := 0
	for _, i := range s {
		if max < i {
			max = i
		}
	}
	return max

}

func (c *Conn) ConnSearch() {
	l, err := ldap.Dial("tcp", fmt.Sprintf("%s:%d", c.getConnHost(), c.getConnPort()))
	if err != nil {
		log.Fatal(err)
	}
	defer l.Close()
	err = l.Bind("cn=arrow,ou=people,dc=bihu,dc=inet", "123456")
	if err != nil {
		log.Fatal(err)
	}

	searchRequest := ldap.NewSearchRequest(
		"dc=bihu,dc=inet", // The base dn to search
		ldap.ScopeWholeSubtree, ldap.NeverDerefAliases, 0, 0, false,
		"(&(objectClass=person))",
		[]string{"dn", "cn", "mail", "uid", "uidNumber"}, // 列出要查询的属性
		nil,
	)

	sr, err := l.Search(searchRequest)
	if err != nil {
		log.Fatal(err)
	}

	var ss []int
	for _, entry := range sr.Entries {
		fmt.Printf("%s: %v %v\n", entry.DN, entry.GetAttributeValue("cn"), entry.GetAttributeValue("uidNumber"))

		uid, err := strconv.Atoi(entry.GetAttributeValue("uidNumber"))
		if err != nil {
			log.Println("错误")
		}
		ss = append(ss, uid)
	}
	fmt.Println(SliceMax(ss))

}

type EntryAttribute struct {
	// Name is the name of the attribute
	Name string
	// Values contain the string values of the attribute
	Values []string
	// ByteValues contain the raw values of the attribute
	ByteValues [][]byte
}

func (c *Conn) AddEntry(username string, email string) {
	l, err := ldap.Dial("tcp", fmt.Sprintf("%s:%d", c.getConnHost(), c.getConnPort()))
	if err != nil {
		log.Fatal(err)
	}
	defer l.Close()
	err = l.Bind("cn=root,dc=bihu,dc=inet", "B1hu12345")
	if err != nil {
		log.Fatal(err)
	}

	dn := "cn=test01,ou=people,dc=bihu,dc=inet"
	attributes := map[string][]interface{}{
		"objectClass": {"top", "shadowAccount", "sambaSamAccount", "person", "posixAccount", "organizationalPerson",
			"openVPNUser", "ldapPublicKey", "inetOrgPerson", "inetLocalMailRecipient", "hostObject"},
		"dockerRegistryActive": {"FALSE"},
		"openvpnbhActive":      {"FALSE"},
		"mail":                 {email},
		"shadowExpire":         {},
	}

	modify := ldap.NewEntry(dn, attributes)

}
