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
	tripPlanner TripPlannerService
}

func NewApp(db *gorm.DB, logger logger.OpswLogger, tripSrv TripPlannerService) *App {
	return &App{
		logger:      logger,
		tripPlanner: tripSrv,
	}
}

func (a App) Boot() {
	start := time.Now()
	a.tripPlanner.IntializeService()
	logger.INFO("Prepare GTFS Data for Stops.")
	// err := a.tripPlanner.StopsData()
	// if err != nil {
	// 	a.logger.ERROR(err.Error())
	// }

	// logger.INFO("Prepare GTFS Data for Route.")
	// err = a.tripPlanner.RoutesData()
	// if err != nil {
	// 	a.logger.ERROR(err.Error())
	// }

	// logger.INFO("Prepare GTFS Data for Calendar Data.")
	// err = a.tripPlanner.CalendarData()
	// if err != nil {
	// 	a.logger.ERROR(err.Error())
	// }

	// logger.INFO("Prepare GTFS Data for Trip Data.")
	// err = a.tripPlanner.TripsData()
	// if err != nil {
	// 	a.logger.ERROR(err.Error())
	// }

	logger.INFO("Prepare GTFS Data for Stop Times Data.")
	err := a.tripPlanner.StopTimesData()
	if err != nil {
		a.logger.ERROR(err.Error())
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
         _____ ____  ___ ____    ____  _             _   _ _   _ _____ ___ 
		|_   _|  _ \|_ _|  _ \  |  _ \| |      / \  | \ | | \ | | ____|  _ \ 
		  | | | |_) || || |_) | | |_) | |     / _ \ |  \| |  \| |  _| | |_) |
		  | | |  _ < | ||  __/  |  __/| |___ / ___ \| |\  | |\  | |___|  _ < 
		  |_| |_| \_\___|_|     |_|   |_____/_/   \_\_| \_|_| \_|_____|_| \_\
	:: OASA Prepare Data for Trip Planner application (v%s) ::`+"\n\n", os.Getenv("application.version"))))
	//os.Stdout = originalStdout

	// Database Migration Proccess
	err = db.DatabaseMigrations()
	if err != nil {
		logger.ERROR(err.Error())
	}
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
