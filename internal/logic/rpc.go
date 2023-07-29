package logic

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/realmicro/realmicro/client"
	"github.com/realmicro/realmicro/errors"
	log "github.com/sirupsen/logrus"
	"net/http"
	"qingyun/common/token"
	"qingyun/services/fishing_gateway/common"
	"qingyun/services/fishing_gateway/internal/config"
	"strconv"
	"strings"
)

var (
	CORS = map[string]bool{"*": true}
)

const (
	HEADER_KEY_SERVICE_NAME  = "servername"
	HEADER_KEY_METHOD_NAME   = "methodname"
	HEADER_KEY_USER_ID       = "user_id"
	HEADER_KEY_VERSION       = "version"
	HEADER_KEY_AUTHORIZATION = "authorization"
)

func Rpc(c *gin.Context) {
	if origin := c.GetHeader("Origin"); CORS[origin] {
		c.Request.Header.Set("Access-Control-Allow-Origin", origin)
	} else if len(origin) > 0 && CORS["*"] {
		c.Request.Header.Set("Access-Control-Allow-Origin", origin)
	}
	c.Request.Header.Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
	c.Request.Header.Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
	c.Request.Header.Set("Access-Control-Allow-Credentials", "true")
	c.Request.Header.Set("Access-Control-Allow-Headers", c.GetHeader("Access-Control-Request-Headers"))
	c.Request.Header.Set("X-Real-Ip", c.ClientIP())
	if c.Request.Method == "OPTIONS" {
		return
	}
	if c.Request.Method != "POST" {
		http.Error(c.Writer, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	badRequest := func(description string) {
		e := errors.BadRequest("real.micro.rpc", description)
		c.Writer.WriteHeader(400)
		c.Writer.Write([]byte(e.Error()))
	}
	var serviceName, method, authorization, userId, version string
	serviceName = c.GetHeader(HEADER_KEY_SERVICE_NAME)
	method = c.GetHeader(HEADER_KEY_METHOD_NAME)
	authorization = c.GetHeader(HEADER_KEY_AUTHORIZATION)
	userId = c.GetHeader(HEADER_KEY_USER_ID)

	if len(serviceName) == 0 {
		badRequest("invalid service")
		return
	}
	if len(method) == 0 {
		badRequest("invalid method")
		return
	}

	if !checkToken(authorization, userId, serviceName, method) {
		log.Infof("[service] invalid token service %s, userid: %s, method: %s, token: %s", serviceName, userId, method, authorization)
		e := errors.Unauthorized("real", "invalid check")
		c.Writer.WriteHeader(401)
		c.Writer.Write([]byte(e.Error()))
		return
	}
	c.Request.Header.Set(HEADER_KEY_USER_ID, userId)
	c.Request.Header.Set(HEADER_KEY_AUTHORIZATION, authorization)
	c.Request.Header.Set(HEADER_KEY_VERSION, version)
	c.Request.Header.Set(HEADER_KEY_SERVICE_NAME, serviceName)
	c.Request.Header.Set(HEADER_KEY_METHOD_NAME, method)
	var request interface{}
	c.BindJSON(&request)
	NewRequest := client.NewRequest(serviceName, method, request, client.WithContentType("application/json"))
	var response gin.H
	ctx := common.RequestToContext(c.Request)
	err := client.Call(ctx, NewRequest, &response)
	if err != nil {
		ce := errors.Parse(err.Error())
		switch ce.Code {
		case 0:
			ce.Code = 500
			ce.Id = "real.micro.rpc"
			ce.Status = http.StatusText(500)
			ce.Detail = "error during request: " + ce.Detail
			c.Writer.WriteHeader(500)
		default:
			c.Writer.WriteHeader(int(ce.Code))
		}
		c.Writer.Write([]byte(ce.Error()))
		return

	}
	c.JSON(200, response)
}

func checkToken(authorization, userid, service, method string) bool {
	white, err := config.GetGoConfig().GetClientGatewayTokenControllerWhiteList()
	if err != nil {
		log.Infof("[gateway] get client white list error: %v", err)
		return false
	}
	if _, ok := white[userid]; ok {
		return true
	}

	serverMap, err := config.GetGoConfig().GetGatewayServerWhiteList()
	fmt.Println(serverMap)
	if err != nil {
		log.Infof("[gateway] get service white list error: %v", err)
		return false
	}
	if _, ok := serverMap[service]; ok {
		return true
	}

	gtMap, err := config.GetGoConfig().GetGatewayTokenController()
	fmt.Println(gtMap,method)
	if err != nil {
		log.Infof("[geteway] get gateway token controller error: %v", err)
		return false
	}
	if _, ok := gtMap[method]; ok {
		return true
	}

	authorization = strings.ReplaceAll(authorization, "Bearer ", "")
	claims, err := token.ParseJwtToken(authorization)
	if err != nil {
		log.Infof("[gateway] ParseJwtToken: %v", err)
		return false
	}
	if strconv.Itoa(int(claims.UserId)) != userid {
		log.Infof("[gateway] UserId: %v", err)
		return false
	}
	return true
}

func Md5(key string) string {
	h := md5.New()
	h.Write([]byte(key))
	return hex.EncodeToString(h.Sum(nil))
}
