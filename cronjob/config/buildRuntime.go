package config

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/cs161079/monorepo/common/db"
	"github.com/cs161079/monorepo/common/mapper"
	"github.com/cs161079/monorepo/common/repository"
	"github.com/cs161079/monorepo/common/service"
	logger "github.com/cs161079/monorepo/common/utils/goLogger"
	"github.com/joho/godotenv"
	"go.uber.org/dig"
)

type App struct {
	logger      logger.OpswLogger
	syncService SyncService
}

func NewApp(logger logger.OpswLogger, syncSrv SyncService) *App {
	return &App{
		logger:      logger,
		syncService: syncSrv,
	}
}

func (a App) Boot() {
	if err := a.syncService.SyncData(); err != nil {
		a.logger.ERROR(fmt.Sprintf("Κάτι πήγε στραβά με την λήψη των δεδομένων. %s", err.Error()))
		// fmt.Printf("Κάτι πήγε στραβά με την λήψη των δεδομένων.")
		return
	}
	if err := a.syncService.DeleteAll(); err != nil {
		a.logger.INFO(fmt.Sprintf("Κάτι πήγε στραβά με την διαγραφή των δεδομένων από την βάση δεδομένων. %s", err.Error()))
		return
	}
	if err := a.syncService.InserttoDatabase(); err != nil {
		a.logger.ERROR(fmt.Sprintf("Κάτι πήγε στραβά με την εισαγωγή των δεδομένων στην βάση δεδομένων. %s", err.Error()))
		return
	}
	if err := a.syncService.SyncSchedule(); err != nil {
		a.logger.ERROR(fmt.Sprintf("Κάτι πήγε στραβά με την λήψη των δρομολογίων. %s", err.Error()))
	}
}

func InitializeApplication() {
	// Redirect fmt's output to a file
	file, err := os.OpenFile(filepath.Join("C:\\logs", "goSyncApplication", "oasaLogs.log"), os.O_APPEND|os.O_CREATE|os.O_RDWR, 0666)
	if err != nil {
		fmt.Println("Error creating file:", err)
		return
	}
	defer file.Close()

	os.Stdout = file // Set output destination

	fmt.Printf(`
  .    ___    _    ____    _    __ _ _  
 /\\  / _ \  / \  / ___|  / \   \ \ \ \
( ( )| | | |/ _ \ \___ \ / _ \   \ \ \ \ 
 \\/ | |_| / ___ \ ___) / ___ \   ) ) ) )  
  '   \___/_/   \_\____/_/   \_\ / / / /
						    	/_/_/_/


:: OASA Synchtonization Data application (v1.0.0) ::

				   `)
	// Load the .env file
	err = godotenv.Load()
	if err != nil {
		fmt.Println("Error loading .env file")
	}
	logger.InitLogger("goSyncApplication")
}

func BuildInRuntime() (*App, error) {
	c := dig.New()
	servicesConstructors := []interface{}{
		logger.CreateLogger,
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

	InitializeApplication()

	var result *App
	err := c.Invoke(func(a *App) {
		result = a
	})
	return result, err
}
