package config

import (
	"fmt"
	"os"
	"time"

	"github.com/cs161079/monorepo/common/db"
	"github.com/cs161079/monorepo/common/mapper"
	"github.com/cs161079/monorepo/common/repository"
	"github.com/cs161079/monorepo/common/service"
	logger "github.com/cs161079/monorepo/common/utils/goLogger"
	"github.com/joho/godotenv"
	"go.uber.org/dig"
	"gorm.io/gorm"
)

type App struct {
	logger      logger.OpswLogger
	syncService SyncService
}

func NewApp(db *gorm.DB, logger logger.OpswLogger, syncSrv SyncService) *App {
	// db.AutoMigrate(&models.Line{}, &models.Route{}, &models.Stop{}, &models.Route01{}, &models.Route02{}, &models.ScheduleMaster{}, &models.ScheduleTime{})
	return &App{
		logger:      logger,
		syncService: syncSrv,
	}
}

func (a App) Boot() {
	start := time.Now()
	if err := a.syncService.SyncData(); err != nil {
		a.logger.ERROR(fmt.Sprintf("Κάτι πήγε στραβά με την λήψη των δεδομένων. %s\n", err.Error()))
		// fmt.Printf("Κάτι πήγε στραβά με την λήψη των δεδομένων.")
		return
	}
	if err := a.syncService.DeleteAll(); err != nil {

		a.logger.INFO(fmt.Sprintf("Κάτι πήγε στραβά με την διαγραφή των δεδομένων από την βάση δεδομένων. %s\n", err.Error()))
		return
	}
	if err := a.syncService.InserttoDatabase(); err != nil {
		a.logger.ERROR(fmt.Sprintf("Κάτι πήγε στραβά με την εισαγωγή των δεδομένων στην βάση δεδομένων. %s\n", err.Error()))
		return
	}

	end := time.Now()
	// Calculate the duration
	duration := end.Sub(start)
	// Output the duration
	logger.INFO(fmt.Sprintf("Duration: %v\n", duration))
}

func InitializeApplication() {
	// Load the .env file
	err := godotenv.Load()
	if err != nil {
		fmt.Println("Error loading .env file")
	}
	logger.CreateLogger()
	// originalStdout := os.Stdout

	logger.Logger.Out.Write([]byte(fmt.Sprintf(`
  .    ___    _    ____    _      
 /\\  / _ \  / \  / ___|  / \   
( ( )| | | |/ _ \ \___ \ / _ \   
 \\/ | |_| / ___ \ ___) / ___ \     
  '   \___/_/   \_\____/_/   \_\ 		    	
             ____                     _       _                                                                __ _ _
            / ___|_ __ ___  _ __     | | ___ | |__           / \   _ __  _ __ | (_) ___ __ _| |_(_) ___  _ __  \ \ \ \
            | |   | '__/ _ \| '_ \ _  | |/ _ \| '_ \        / _ \ | '_ \| '_ \| | |/ __/ _  | __| |/ _ \| '_ \  \ \ \ \
            | |___| | | (_) | | | | |_| | (_) | |_) |      / ___ \| |_) | |_) | | | (_| (_| | |_| | (_) | | | |  ) ) ) )
             \____|_|  \___/|_| |_|\___/ \___/|_.__/      /_/   \_\ .__/| .__/|_|_|\___\__,_|\__|_|\___/|_| |_| / / / / 
                                                                  |_|   |_|                                    /_/_/_/
	:: OASA Synchtonization Data application (v%s) ::`+"\n\n", os.Getenv("application.version"))))
	//os.Stdout = originalStdout
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
