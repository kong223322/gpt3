package common

import (
	"bytes"
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
	"io/ioutil"
	"time"
)

type BodyLogWriter struct {
	gin.ResponseWriter
	body *bytes.Buffer
}

func (w BodyLogWriter) Write(b []byte) (int, error) {
	w.body.Write(b)
	return w.ResponseWriter.Write(b)
}
func (w BodyLogWriter) WriteString(s string) (int, error) {
	w.body.WriteString(s)
	return w.ResponseWriter.WriteString(s)
}

// 打印日志
func Logger() gin.HandlerFunc {
	return func(c *gin.Context) {
		requestBody, _ := c.GetRawData()
		c.Request.Body = ioutil.NopCloser(bytes.NewBuffer(requestBody))

		bodyLogWriter := &BodyLogWriter{body: bytes.NewBufferString(""), ResponseWriter: c.Writer}
		c.Writer = bodyLogWriter

		start := time.Now()
		//handler
		c.Next()
		//log
		end := time.Now()
		responseBody := bodyLogWriter.body.String()
		log.WithFields(log.Fields{
			"uri":             c.Request.URL.Path,
			"raw_query":       c.Request.URL.RawQuery,
			"start_timestamp": start.Format("2006-01-02 15:04:05"),
			"end_timestamp":   end.Format("2006-01-02 15:04:05"),
			"server_name":     c.Request.Host,
			"remote_addr":     c.ClientIP(),
			"proto":           c.Request.Proto,
			"referer":         c.Request.Referer(),
			"request_method":  c.Request.Method,
			"response_time":   end.Sub(start).Milliseconds(), // 毫秒
			"content_type":    c.Request.Header.Get("Content-Type"),
			"status":          c.Writer.Status(),
			"user_agent":      c.Request.UserAgent(),
			"trace_id":        c.Request.Header.Get("X-Request-Trace-Id"),
			"response":        responseBody,
			"request":         string(requestBody),
			"response_err":    c.Errors.Errors(),
		}).Log(log.InfoLevel)
	}
}
