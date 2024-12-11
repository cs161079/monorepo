package controllers

import (
	"net/http"

	"github.com/cs161079/monorepo/common/service"
	"github.com/gin-gonic/gin"
)

type LineController interface {
	GetLineList(*gin.Context)
}

type LineControllerImplementation struct {
	svc service.LineService
}

func NewLineController(svc service.LineService) LineControllerImplementation {
	return LineControllerImplementation{
		svc: svc,
	}
}

func (u LineControllerImplementation) GetLineList(c *gin.Context) {
	data, err := u.svc.GetLineList()
	if err != nil {
		c.Error(err)
		return
	}
	c.JSON(http.StatusOK, map[string]any{
		"lines": data,
	})
}

func (u LineControllerImplementation) AddRouters(eng *gin.Engine) {
	apiGroup := eng.Group("/line")
	apiGroup.GET("/list", u.GetLineList)
}
