package logger

import (
	"os"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

// L is the ptr to the logger
var L *logrus.Logger

// Init logger
func Init(f *os.File, debugLevel bool) {
	// logrus.SetFormatter(&logrus.JSONFormatter{})
	L = logrus.New()
	if debugLevel { // beta环境用debug也记录
		L.SetLevel(logrus.DebugLevel)
	}
	L.Formatter = &logrus.JSONFormatter{}

	if f != nil {
		L.Out = f
	} else {
		L.Info("Failed to log to file, using default stderr")
	}
}

func InitFromLogger(input *logrus.Logger) {
	L = input
}

func InitForTest() {
	L = logrus.New()
}

func SetCtx(c *gin.Context, p, action string) *logrus.Entry {
	var adminID uint
	if c != nil {
		cadminID, exists := c.Get("current_admin_id")
		if exists {
			adminID = cadminID.(uint)
		}
	}
	ctx := L.WithFields(logrus.Fields{
		"admin_id": adminID,
		"package":  p,
		"action":   action,
	})
	return ctx
}

func Info(c *gin.Context, p, action string, message ...interface{}) {
	ctx := SetCtx(c, p, action)
	ctx.Info(message...)
}

func Debug(c *gin.Context, p, action string, message ...interface{}) {
	ctx := SetCtx(c, p, action)
	ctx.Debug(message...)
}

func Warn(c *gin.Context, p, action string, message ...interface{}) {
	ctx := SetCtx(c, p, action)
	ctx.Warn(message...)
}

func Error(c *gin.Context, p, action string, message ...interface{}) {
	ctx := SetCtx(c, p, action)
	ctx.Error(message...)
}

func Fatal(c *gin.Context, p, action string, message ...interface{}) {
	ctx := SetCtx(c, p, action)
	ctx.Fatal(message...)
}

func Panic(c *gin.Context, p, action string, message ...interface{}) {
	ctx := SetCtx(c, p, action)
	ctx.Panic(message...)
}
