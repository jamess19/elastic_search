package middleware

import (
	"bytes"
	"encoding/json"
	"business/conf"
	"business/pkg/utils"
	"fmt"
	"io/ioutil"
	"net/http"
	"runtime/debug"
	"github.com/gin-gonic/gin"
	"gitlab.com/goxp/cloud0/logger"
)

func LoggingRequest() gin.HandlerFunc {
	return func(c *gin.Context) {
		log := logger.WithCtx(c, "Request Detail")
		defer func() {
			if r := recover(); r != nil {
				log.Error(r)
				debug.PrintStack()
				panic(r)
			}
		}()
		r := c.Request
		header := c.Request.Header
		buf, err := ioutil.ReadAll(r.Body)
		if err != nil {
			log.WithError(err).Errorf("Error reading request body: %v", err.Error())
			http.Error(c.Writer, err.Error(), http.StatusInternalServerError)
			return
		}

		var obj interface{}
		_ = json.Unmarshal(buf, &obj)

		data, err := json.Marshal(obj)
		if conf.LoadEnv().AppEnv != utils.ENV_PRD {
			log.WithField("body", fmt.Sprintf("%s", data)).WithField("header", header).Info("uri: ", c.Request.RequestURI)
		} else {
			log.WithField("body", fmt.Sprintf("%s", data)).Info("uri: ", c.Request.RequestURI)
		}
		reader := ioutil.NopCloser(bytes.NewBuffer(buf))
		c.Request.Body = reader
		c.Next()
	}
}
