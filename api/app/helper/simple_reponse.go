package helper

import (
	"encoding/json"
	log "github.com/Sirupsen/logrus"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
)

type RespJson struct {
	statusCode int    `json:"statusCode"`
	Error      string `json:"error,omitempty"`
	Msg        string `json:"message,omitempty"`
}

func JSONR(c *gin.Context, code int, dat interface{}, message interface{}) (werror error) {
	var (
		wcode int
		data  interface{}
		msg   interface{}
	)
	wcode = code
	data = dat
	msg = message

	need_doc := viper.GetBool("gen_doc")
	var body interface{}
	defer func() {
		if need_doc {
			ds, _ := json.Marshal(body)
			bodys := string(ds)
			log.Debugf("body: %v, bodys: %v ", body, bodys)
			c.Set("body_doc", bodys)
		}
	}()

	switch msg.(type) {

	case string:
		body = gin.H{"statusCode": wcode, "data": data, "message": msg.(string)}
		c.JSON(wcode, body)
	case error:
		body = gin.H{"statusCode": wcode, "data": data, "message": msg.(error).Error()}
		c.JSON(wcode, body)
	default:
		body = gin.H{"statusCode": wcode, "data": data, "message": "unknow err"}
		c.JSON(wcode, body)
	}
	return

	//need_doc := viper.GetBool("gen_doc")
	//var body interface{}
	//defer func() {
	//	if need_doc {
	//		ds, _ := json.Marshal(body)
	//		bodys := string(ds)
	//		log.Debugf("body: %v, bodys: %v ", body, bodys)
	//		c.Set("body_doc", bodys)
	//	}
	//}()

	//if wcode == 200 {
	//	switch msg.(type) {
	//	case string:
	//		//body = RespJson{statusCode: wcode, Msg: msg.(string)}
	//		body = gin.H{"statusCode": wcode, "data": data, "message": msg.(string)}
	//		c.JSON(http.StatusOK, body)
	//	default:
	//		c.JSON(http.StatusOK, msg)
	//		body = msg
	//	}
	//} else {
	//	switch msg.(type) {
	//	case string:
	//		body = RespJson{statusCode: wcode, Error: msg.(string)}
	//		c.JSON(wcode, body)
	//	case error:
	//		body = RespJson{statusCode: wcode, Error: msg.(error).Error()}
	//		c.JSON(wcode, body)
	//	default:
	//		body = RespJson{statusCode: wcode, Error: "system type error. please ask admin for help"}
	//		c.JSON(wcode, body)
	//	}
	//}

}
