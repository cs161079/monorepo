package controllers

import (
	"net/http"

	"github.com/cs161079/monorepo/common/models"
	"github.com/cs161079/monorepo/common/service"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type NotificationController interface {
	AddRouters(eng *gin.Engine)
}

type notificationControllerImpl struct {
	connection *gorm.DB
	ntSvc      service.NotificationService
}

func NewNotifcationController(svc service.NotificationService) NotificationController {
	return &notificationControllerImpl{
		ntSvc: svc,
	}
}

func (c notificationControllerImpl) AddRouters(eng *gin.Engine) {
	apiGroup := eng.Group("/notification")
	apiGroup.POST("/push", c.pushNotification)
}

func (c notificationControllerImpl) pushNotification(ctx *gin.Context) {
	var notification models.Notification
	if err := ctx.ShouldBindJSON(&notification); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := c.ntSvc.SendPushNotification(notification); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"resposne": "Notificatin Pushed Successfully!"})
}
