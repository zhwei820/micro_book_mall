package utils

import (
	"context"
	"fmt"
	"github.com/micro/go-micro/client"
	"net/http"

	auth "github.com/micro-in-cn/tutorials/microservice-in-micro/part5/auth/proto/auth"
	"github.com/micro-in-cn/tutorials/microservice-in-micro/part5/basic/common"
	"github.com/micro-in-cn/tutorials/microservice-in-micro/part5/plugins/session"
	logzap "github.com/micro-in-cn/tutorials/microservice-in-micro/part5/plugins/zap"
)

var (
	log        = logzap.GetLogger()
	authClient = auth.NewService("mu.micro.book.srv.auth", client.DefaultClient)
)

// AuthWrapper 认证wrapper
func AuthWrapper(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ck, _ := r.Cookie(common.RememberMeCookieName)
		// token不存在，则状态异常，无权限
		if ck == nil {
			log.Info(fmt.Sprintf("token不存在"))
			http.Error(w, "非法请求", 400)
			return
		}

		sess := session.GetSession(w, r)
		if sess.ID != "" {
			// 检测是否通过验证
			if sess.Values["valid"] != nil {
				h.ServeHTTP(w, r)
				return
			} else {
				userId := sess.Values["userId"].(int64)
				if userId != 0 {
					rsp, err := authClient.GetCachedAccessToken(context.TODO(), &auth.Request{
						UserId: userId,
					})
					if err != nil {
						log.Info(fmt.Sprintf("[AuthWrapper]，err：%s", err))
						http.Error(w, "非法请求", 400)
						return
					}

					// token不一致
					if rsp.Token != ck.Value {
						log.Info(fmt.Sprintf("[AuthWrapper]，token不一致"))
						http.Error(w, "非法请求", 400)
						return
					}
				} else {
					log.Info(fmt.Sprintf("[AuthWrapper]，session不合法，无用户id"))
					http.Error(w, "非法请求", 400)
					return
				}
			}
		} else {
			http.Error(w, "非法请求", 400)
			return
		}

		h.ServeHTTP(w, r)
	})
}
