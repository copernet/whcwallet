package routers

import (
	"fmt"
	"time"

	"github.com/copernet/whcwallet/api"
	"github.com/gin-gonic/gin"
)

var supportTxType = map[string]int{
	"0":  0,
	"1":  1,
	"3":  3,
	"4":  4,
	"50": 50,
	"51": 51,
	"53": 53,
	"54": 54,
	"55": 55,
	"56": 56,
	"68": 68,
	"70": 70,
	"185": 185,
	"186": 186,
}

func IsSupportTxType() gin.HandlerFunc {
	return func(c *gin.Context) {
		reqtype := c.Param("txtype")
		txtype, ok := supportTxType[reqtype]
		if !ok {
			c.AbortWithStatusJSON(200, gin.H{
				"code":    api.ErrUnSupportTxType,
				"message": "UnSupported tx type: " + reqtype,
				"result":  "",
			})
		}

		c.Set("txtype", txtype)

		c.Next()
	}
}

func FixedFormParam() gin.HandlerFunc {
	return func(c *gin.Context) {
		// check tx version
		txversion := c.PostForm("transaction_version")
		if txversion != "0" {
			c.AbortWithStatusJSON(200, gin.H{
				"code":    api.ErrUnSupportTxVersion,
				"message": "UnSupported tx version: " + txversion,
				"result":  "",
			})
		}

		// check ecosystem
		// only support ecosystem = 1 at the current version
		eco := c.PostForm("ecosystem")
		if eco != "" && eco != "1" {
			c.AbortWithStatusJSON(200, gin.H{
				"code":    api.ErrUnSupportEcosystem,
				"message": "UnSupported ecosystem: " + eco,
				"result":  "",
			})
		}

		// check asset type for participating crowd sale
		desiredSP := c.PostForm("currency_identifier_desired")
		if desiredSP != "" && desiredSP != "1" {
			c.AbortWithStatusJSON(200, gin.H{
				"code":    api.ErrUnSupportCrowdSaleSP,
				"message": "UnSupported SP for crowd sale: " + desiredSP,
				"result":  "",
			})
		}

		c.Next()
	}
}

func ginLogger() gin.HandlerFunc {
	// not log the followings router
	notlogged := []string{"/static"}

	out := gin.DefaultWriter
	var skip map[string]struct{}

	if length := len(notlogged); length > 0 {
		skip = make(map[string]struct{}, length)

		for _, path := range notlogged {
			skip[path] = struct{}{}
		}
	}

	return func(c *gin.Context) {
		// Start timer
		start := time.Now()
		path := c.Request.URL.Path
		raw := c.Request.URL.RawQuery

		// Process request
		c.Next()

		// Log only when path is not being skipped
		if _, ok := skip[path]; !ok {
			// Stop timer
			end := time.Now()
			latency := end.Sub(start)

			clientIP := c.ClientIP()
			method := c.Request.Method
			statusCode := c.Writer.Status()
			var statusColor, methodColor, resetColor string
			comment := c.Errors.ByType(gin.ErrorTypePrivate).String()

			if raw != "" {
				path = path + "?" + raw
			}

			fmt.Fprintf(out, "[GIN] %v |%s %3d %s| %13v | %15s |%s %-7s %s %s\n%s",
				end.Format("2006/01/02 - 15:04:05"),
				statusColor, statusCode, resetColor,
				latency,
				clientIP,
				methodColor, method, resetColor,
				path,
				comment,
			)
		}
	}
}
