package logger

import (
	"fmt"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

// L is the ptr to the logger
var L *logrus.Logger

// Init logger
func Init() {
	L = logrus.New()
	L.Formatter = &logrus.JSONFormatter{}
	file, err := os.OpenFile("./logs/logrus.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	fmt.Println(file, err)
	if err == nil {
		L.Out = file
	} else {
		L.Info("Failed to log to file, using default stderr")
	}
}

// GinLogger is a gin logger middleware
func GinLogger() gin.HandlerFunc {
	return func(c *gin.Context) {
		// Start timer
		start := time.Now()
		path := c.Request.URL.Path
		raw := c.Request.URL.RawQuery

		// Process request
		c.Next()

		// Stop timer
		end := time.Now()
		latency := end.Sub(start)

		clientIP := c.ClientIP()
		method := c.Request.Method
		statusCode := c.Writer.Status()

		comment := c.Errors.ByType(gin.ErrorTypePrivate).String()

		if raw != "" {
			path = path + "?" + raw
		}

		s := fmt.Sprintf("[GIN] %v | %3d | %13v | %15s | %-7s %s\n%s",
			end.Format("2006/01/02 - 15:04:05"),
			statusCode,
			latency,
			clientIP,
			method,
			path,
			comment,
		)

		L.Info(s)
	}
}

func SetCtx(p, action string) *logrus.Entry {
	ctx := L.WithFields(logrus.Fields{
		"package": p,
		"action":  action,
	})
	return ctx
}

func Debug(p, action string, message ...interface{}) {
	ctx := L.WithFields(logrus.Fields{
		"package": p,
		"action":  action,
	})
	ctx.Debug(message...)
}

func Info(p, action string, message ...interface{}) {
	ctx := L.WithFields(logrus.Fields{
		"package": p,
		"action":  action,
	})
	ctx.Info(message...)
}

func Warn(p, action string, message ...interface{}) {
	ctx := L.WithFields(logrus.Fields{
		"package": p,
		"action":  action,
	})
	ctx.Warn(message...)
}

func Error(p, action string, message ...interface{}) {
	ctx := L.WithFields(logrus.Fields{
		"package": p,
		"action":  action,
	})
	ctx.Error(message...)
}

func Fatal(p, action string, message ...interface{}) {
	ctx := L.WithFields(logrus.Fields{
		"package": p,
		"action":  action,
	})
	ctx.Fatal(message...)
}

func Panic(p, action string, message ...interface{}) {
	ctx := L.WithFields(logrus.Fields{
		"package": p,
		"action":  action,
	})
	ctx.Panic(message...)
}
