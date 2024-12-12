package config

import (
	"fmt"

	"github.com/cs161079/monorepo/common/db"
	"github.com/cs161079/monorepo/common/mapper"
	"github.com/cs161079/monorepo/common/repository"
	"github.com/cs161079/monorepo/common/service"
	"go.uber.org/dig"
)

type App struct {
	syncService SyncService
}

func NewApp(syncSrv SyncService) *App {
	return &App{
		syncService: syncSrv,
	}
}

func (a App) Boot() {
	a.syncService.SyncData()
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
		mapper.NewRouteDetailMapper,
		service.NewRouteService,
		service.NewShedule01Service,
		service.NewSheduleService,
		service.NewStopService,
		service.NewuVersionService,
		service.NewRestService,
		NewSyncService,
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
