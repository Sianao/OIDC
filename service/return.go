package service

import "github.com/gin-gonic/gin"

func ErrorReturn(c *gin.Context, EInfo interface{}) {
	c.JSON(200, gin.H{
		"state": false,
		"msg":   EInfo,
	})
	return
}
func NormalReturn(c *gin.Context, Info interface{}) {
	c.JSON(200, gin.H{
		"state": true,
		"msg":   Info,
	})
	return
}
