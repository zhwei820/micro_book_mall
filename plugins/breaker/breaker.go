package breaker

import (
	"errors"
	"fmt"
	"github.com/afex/hystrix-go/hystrix"
	statusCode "github.com/micro-in-cn/tutorials/microservice-in-micro/part5/plugins/breaker/http"
	logzap "github.com/micro-in-cn/tutorials/microservice-in-micro/part5/plugins/zap"
	"net/http"
)

var (
	log = logzap.GetLogger()
)

//BreakerWrapper hystrix breaker
func BreakerWrapper(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		name := r.Method + "-" + r.RequestURI
		hystrix.Do(name, func() error {
			sct := &statusCode.StatusCodeTracker{ResponseWriter: w, Status: http.StatusOK}
			h.ServeHTTP(sct.WrappedResponseWriter(), r)

			if sct.Status >= http.StatusInternalServerError {
				str := fmt.Sprintf("status code %d", sct.Status)
				return errors.New(str)
			}
			return nil
		}, func(e error) error {
			log.Info("hystrix BreakerWrapper" + e.Error())
			if e == hystrix.ErrCircuitOpen {
				w.WriteHeader(http.StatusAccepted)
				w.Write([]byte("{\"code\":500,\"detail\":\"稍后重试\"}"))
			} else if e == hystrix.ErrTimeout {
				w.WriteHeader(http.StatusAccepted)
				w.Write([]byte("{\"code\":500,\"detail\":\"稍后重试\"}"))
			}
			//w.Write([]byte(e.Error()))

			return e
		})
	})
}
