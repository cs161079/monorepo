package main

import (
	"fmt"

	"github.com/cs161079/monorepo/common/db"
	"github.com/cs161079/monorepo/common/repository"
	"github.com/cs161079/monorepo/common/service"
	"github.com/cs161079/monorepo/webApplication/controllers"
	"github.com/gin-gonic/gin"
	"go.uber.org/dig"
)

type App struct {
	engine *gin.Engine
}

func NewApp(lineCtrl controllers.LineControllerImplementation) *App {
	eng := gin.Default()
	lineCtrl.AddRouters(eng)

	return &App{
		engine: eng,
	}
}

func (a App) Boot() {
	a.engine.Run(":8080")
}

func BuildInRuntime() (*App, error) {
	c := dig.New()
	servicesConstructors := []interface{}{
		db.NewOpswConnection,
		repository.NewLineRepository,
		repository.NewRouteRepository,
		repository.NewRoute01Repository,
		repository.NewRoute02Repository,
		repository.NewSchedule01Repository,
		repository.NewScheduleRepository,
		repository.NewStopRepository,
		repository.NewUversionRepository,
		service.NewLineService,
		service.NewRouteService,
		service.NewShedule01Service,
		service.NewSheduleService,
		service.NewStopService,
		service.NewuVersionService,
		controllers.NewLineController,
		NewApp,
	}

	for _, service := range servicesConstructors {
		if err := c.Provide(service); err != nil {
			fmt.Printf("Error on Providing %v", err)
			return nil, err
		}
	}

	var result *App
	err := c.Invoke(func(a *App) {
		result = a
	})
	return result, err
}
