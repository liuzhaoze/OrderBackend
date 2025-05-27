package middleware

import (
	"bytes"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"io"
	"time"
)

func LogHttpRequest(entry *logrus.Entry) gin.HandlerFunc {
	return func(context *gin.Context) {
		requestIn(context, entry)
		defer requestOut(context, entry)

		context.Next()
	}
}

func requestIn(context *gin.Context, entry *logrus.Entry) {
	requestArrival := time.Now()
	requestBody, _ := io.ReadAll(context.Request.Body)
	// 读取 request body 后缓冲区为空，需要写回
	context.Request.Body = io.NopCloser(bytes.NewBuffer(requestBody))

	var compactedRequest bytes.Buffer
	_ = json.Compact(&compactedRequest, requestBody)

	context.Set("arrival", requestArrival)
	entry.WithContext(context.Request.Context()).WithFields(logrus.Fields{
		"arrival": requestArrival.Unix(),
		"from":    context.RemoteIP(),
		"uri":     context.Request.RequestURI,
		"request": compactedRequest.String(),
	}).Infoln("http request arrived")
}

func requestOut(context *gin.Context, entry *logrus.Entry) {
	response, _ := context.Get("response")
	arrival, _ := context.Get("arrival")

	entry.WithContext(context.Request.Context()).WithFields(logrus.Fields{
		"cost":     time.Since(arrival.(time.Time)).Milliseconds(),
		"response": response,
	}).Infoln("http request finished")
}
