package errors

import (
	"encoding/json"

	"github.com/gin-gonic/gin"
)

type ErrorFormatter struct {
}

func (e ErrorFormatter) FormatError(c *gin.Context, err error, status int) {
	if _, ok := err.(json.Marshaler); ok {
		c.JSON(status, gin.H{"error": err})
	} else {
		c.JSON(status, gin.H{"error": err.Error()})
	}

	c.Abort()
	return
}
