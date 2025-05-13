package controllers

import (
	"net/http"
	"time"

	"github.com/coreos/go-oidc"
	"github.com/cs161079/monorepo/common/models"
	"github.com/cs161079/monorepo/common/service"
	"github.com/cs161079/monorepo/common/utils"
	"github.com/cs161079/monorepo/webApplication/keycloak"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type AdminController interface {
	AddRouters(eng *gin.Engine)
}

type adminControllerImpl struct {
	connection *gorm.DB
	ntSvc      service.NotificationService
	verifier   *oidc.IDTokenVerifier
}

func NewAdminController(dbConnection *gorm.DB, svc service.NotificationService, verif *oidc.IDTokenVerifier) AdminController {
	return &adminControllerImpl{
		connection: dbConnection,
		ntSvc:      svc,
		verifier:   verif,
	}
}

func (c *adminControllerImpl) AddRouters(eng *gin.Engine) {
	apiGroup := eng.Group("/admin")
	apiGroup.Use(keycloak.AuthMiddleware(c.verifier))
	apiGroup.POST("/push", keycloak.RequireRole("oasaAdmin"), c.pushNotification)
	apiGroup.GET("/jobs", c.jobList)

}

func (c *adminControllerImpl) pushNotification(ctx *gin.Context) {
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

type PaginatedResult struct {
	Data       interface{} `json:"data"`
	Total      int64       `json:"total"`
	Page       int         `json:"page"`
	PageSize   int         `json:"pageSize"`
	TotalPages int         `json:"totalPages"`
}

func (c *adminControllerImpl) jobList(ctx *gin.Context) {
	start := time.Now()
	var pagePtr, err = utils.StrToInt32(ctx.Query("page"))
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	var page = int(*pagePtr)
	var pageSize = 10
	var total int64
	offset := (page - 1) * pageSize

	var recs []models.OpswCronRuns
	// Count total records
	c.connection.Model(&recs).Count(&total)

	// Fetch page data
	result := c.connection.Order("runtime DESC").Limit(pageSize).Offset(offset).Find(&recs)
	if result.Error != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": result.Error.Error()})
		return
	}

	ctx.JSON(http.StatusOK, map[string]any{"duration": time.Since(start).Seconds(), "pagingData": PaginatedResult{
		Data:       recs,
		Total:      total,
		Page:       page,
		PageSize:   pageSize,
		TotalPages: int((total + int64(pageSize) - 1) / int64(pageSize)),
	}})
}
