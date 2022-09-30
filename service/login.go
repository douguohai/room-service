package service

import (
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"github.com/douguohai/room-service/base"
	"github.com/douguohai/room-service/entity"
	"io/ioutil"
	"net/http"
)

var MD5 = md5.New()

func CrosMiddleware(handler http.HandlerFunc) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")                                                                                         // 指明哪些请求源被允许访问资源，值可以为 "*"，"null"，或者单个源地址。
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type,AccessToken,X-CSRF-Token, Authorization, Token")                              //对于预请求来说，指明了哪些头信息可以用于实际的请求中。
		w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS")                                                                       //对于预请求来说，哪些请求方式可以用于实际的请求。
		w.Header().Set("Access-Control-Expose-Headers", "Content-Length, Access-Control-Allow-Origin, Access-Control-Allow-Headers, Content-Type") //对于预请求来说，指明哪些头信息可以安全的暴露给 CORS API 规范的 API
		w.Header().Set("Access-Control-Allow-Credentials", "true")                                                                                 //指明当请求中省略 creadentials 标识时响应是否暴露。对于预请求来说，它表明实际的请求中可以包含用户凭证。

		//放行所有OPTIONS方法
		if r.Method == "OPTIONS" {
			w.WriteHeader(200)
			return
		}
		handler.ServeHTTP(w, r)
	})
}

//Login 用户登录，返回token信息
func Login(w http.ResponseWriter, r *http.Request) {
	result := entity.Success(base.SuccessMsg)
	body, err := ioutil.ReadAll(r.Body)
	if err == nil {
		loginUser := entity.Login{}
		if err = json.Unmarshal(body, &loginUser); err == nil {
			uid := getUid(loginUser.Username)
			fmt.Printf("用户登录  用户名: %v 相关uid: %v\n", loginUser.Username, uid)
			result.Data = uid
			msg, _ := json.Marshal(result)
			_, _ = w.Write(msg)
			return
		}
	}
	result.Code = base.FailCode
	msg, _ := json.Marshal(entity.Fail(base.ParameterParsingErrorMsg))
	_, _ = w.Write(msg)
}

// 根据字符串生成uid，uid相同代表是同一个人
func getUid(str string) string {
	return hex.EncodeToString(MD5.Sum([]byte(str)))
}
