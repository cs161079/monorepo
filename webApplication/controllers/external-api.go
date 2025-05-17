package controllers

import (
	"net/http"

	"github.com/coreos/go-oidc"
	"github.com/cs161079/monorepo/common/db"
	"github.com/cs161079/monorepo/webApplication/keycloak"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type ExtApiController interface {
	AddRouters(eng *gin.Engine)
	getLines(*gin.Context)
}

type extApiControllerImpl struct {
	connection *gorm.DB
	verifier   *oidc.IDTokenVerifier
}

func NewExtApiController(dbConnection *gorm.DB, verif *oidc.IDTokenVerifier) ExtApiController {
	return &extApiControllerImpl{
		connection: dbConnection,
		verifier:   verif,
	}
}

func (c *extApiControllerImpl) AddRouters(eng *gin.Engine) {
	apiGroup := eng.Group("/api/ext")
	apiGroup.Use(keycloak.AuthMiddleware(c.verifier))
	apiGroup.GET("/lines", c.getLines)
}

func (c *extApiControllerImpl) getLines(ctx *gin.Context) {
	var results []LineExtDto

	if err := c.connection.Table(db.LINETABLE).Find(&results).Error; err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, results)
}

type LineExtDto struct {
	MLCode       int    `json:"ml_code"`
	SDCCode      int    `json:"sdc_code"`
	LineCode     int    `json:"line_code"`
	LineID       string `json:"line_id"`
	LineDescr    string `json:"line_descr"`
	LineDescrEng string `json:"line_descr_eng"`
	MldMaster    int16  `json:"mld_master"`
}
