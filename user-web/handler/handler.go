package handler

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/labstack/echo"
	"github.com/micro-in-cn/tutorials/microservice-in-micro/part5/utils/econtext"
	"github.com/micro-in-cn/tutorials/microservice-in-micro/part5/utils/errormap"
	"net/http"
	"time"

	auth "github.com/micro-in-cn/tutorials/microservice-in-micro/part5/auth/proto/auth"
	"github.com/micro-in-cn/tutorials/microservice-in-micro/part5/plugins/session"
	us "github.com/micro-in-cn/tutorials/microservice-in-micro/part5/user-srv/proto/user"
)

type User struct {
	UserName string `json:"userName" form:"userName" `
	Pwd      string `json:"pwd" form:"pwd" `
}

// Login 登录入口
func Login(c echo.Context) {
	// 返回结果
	response := map[string]interface{}{
		"data": nil,
		"err":  nil,
	}

	ctx := econtext.ParseContext(c)

	user := &User{}
	if err := c.Bind(user); err != nil {
		response["err"] = errormap.GetError(errormap.ErrParam) // 参数错误
		c.JSON(400, response)
		return
	}
	// 调用后台服务
	rsp, err := serviceClient.QueryUserByName(ctx, &us.Request{
		UserName: user.UserName,
	})
	if err != nil {
		response["err"] = errormap.GetError(errormap.ErrQueryUserByName) // 获取用户失败
		c.JSON(500, response)
		return
	}
	respData := map[string]interface{}{}
	if rsp.User.Pwd == user.Pwd {

		// 干掉密码返回
		rsp.User.Pwd = ""
		respData["user"] = rsp.User
		log.Info(ctx, fmt.Sprintf("[Login] 密码校验完成，生成token..."))

		// 生成token
		rsp2, err := authClient.MakeAccessToken(context.TODO(), &auth.Request{
			UserId:   rsp.User.Id,
			UserName: rsp.User.Name,
		})
		if err != nil {
			log.Info(ctx, fmt.Sprintf("[Login] 创建token失败，err：%s", err))
			response["err"] = errormap.GetError(errormap.ErrMakeAccessToken) // 获取用户失败
			c.JSON(500, response)                                            // 生成token失败
			return
		}

		log.Info(ctx, fmt.Sprintf("[Login] token %s", rsp2.Token))
		respData["token"] = rsp2.Token

		// 同步到session中

	} else {
		response["error"] = errormap.Error{
			Detail: "密码错误",
		}
	}
	response["data"] = respData
	c.JSON(200, response)
}

// Logout 退出登录
func Logout(c echo.Context) {

	tokenCookie, err := r.Cookie("remember-me-token")
	if err != nil {
		log.Info(fmt.Sprintf("token获取失败"))
		http.Error(w, "非法请求", 400)
		return
	}

	// 删除token
	_, err = authClient.DelUserAccessToken(context.TODO(), &auth.Request{
		Token: tokenCookie.Value,
	})
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	// 清除cookie
	cookie := http.Cookie{Name: "remember-me-token", Value: "", Path: "/", Expires: time.Now().Add(0 * time.Second), MaxAge: 0}
	http.SetCookie(w, &cookie)

	w.Header().Add("Content-Type", "application/json; charset=utf-8")

	// 返回结果
	response := map[string]interface{}{
		"ref":     time.Now().UnixNano(),
		"success": true,
	}

	// 返回JSON结构
	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
}

func TestSession(w http.ResponseWriter, r *http.Request) {
	sess := session.GetSession(w, r)

	if v, ok := sess.Values["path"]; !ok {
		sess.Values["path"] = r.URL.Query().Get("path")
		log.Info(fmt.Sprintf("path:" + r.URL.Query().Get("path")))
	} else {
		log.Info(fmt.Sprintf(v.(string)))
	}

	log.Info(fmt.Sprintf(sess.ID))
	log.Info(fmt.Sprintf(sess.Name()))

	w.Write([]byte("OK"))
}
