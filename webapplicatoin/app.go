package main

import "github.com/gin-gonic/gin"

type WebApplication struct {
	engine *gin.Engine
}

func (a *WebApplication) Start() {
	a.engine.Run(":8080")
}

func NewEngine() *gin.Engine {
	return gin.Default()
}

func NewRouter() {
	
}
