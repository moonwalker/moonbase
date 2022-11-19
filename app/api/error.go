package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type Error struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Error   string `json:"error,omitempty"`
}

func httpError(c *gin.Context, code int, msg string, err error) {
	var e string
	if err != nil {
		e = err.Error()
	}
	c.JSON(http.StatusInternalServerError, Error{code, msg, e})
}
