package config

import (
	"fmt"
	"os"

	"github.com/cs161079/monorepo/common/db"
	"github.com/cs161079/monorepo/common/mapper"
	"github.com/cs161079/monorepo/common/repository"
	"github.com/cs161079/monorepo/common/service"
	logger "github.com/cs161079/monorepo/common/utils/goLogger"
	"github.com/cs161079/monorepo/webApplication/controllers"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"go.uber.org/dig"
	"gorm.io/gorm"
)

type App struct {
	engine *gin.Engine
}

// ErrorHandler is a custom middleware for handling errors
// When panic occured from programm ErrorHandler is here to catch it.
func ErrorHandler(c *gin.Context, err any) {
	var httpResponse = map[string]any{"error": "Internal server error", "code": -1}
	c.AbortWithStatusJSON(500, httpResponse)
}

func NewApp(db *gorm.DB, lineCtrl controllers.LineController, rtCtr controllers.RouteController,
	stopCtrl controllers.StopController, schedCtrl controllers.ScheduleController, compCtrl controllers.ComponentController,
	testCtrl controllers.TestController, oasaCtrl controllers.OasaNativeController, notifyCtrl controllers.NotificationController) *App {
	gin.SetMode(gin.ReleaseMode)
	eng := gin.New()
	eng.Use(cors.Default())
	gin.DefaultWriter = logger.Logger.Out
	gin.DefaultErrorWriter = logger.Logger.Out
	eng.Use(gin.Logger(), gin.CustomRecovery(ErrorHandler))

	lineCtrl.AddRouters(eng)
	rtCtr.AddRouters(eng)
	stopCtrl.AddRouters(eng)
	schedCtrl.AddRouters(eng)
	compCtrl.AddRouters(eng)
	testCtrl.AddRoutes(eng)
	oasaCtrl.AddRouters(eng)
	notifyCtrl.AddRouters(eng)

	//db.AutoMigrate(&models.Line{}, &models.Route{}, &models.Stop{}, &models.Route01{}, &models.Route02{}, &models.ScheduleMaster{}, &models.ScheduleTime{})

	return &App{
		engine: eng,
	}
}

func (a App) Boot() {
	var port = os.Getenv("SERVER_PORT")
	if port == "" {
		port = "8080"
	}
	logger.INFO(fmt.Sprintf("Application server start on port %s", port))
	a.engine.Run(fmt.Sprintf(":%s", port))
}

func InitializeApplication() {
	// Load the .env file
	err := godotenv.Load()
	if err != nil {
		fmt.Println("Error loading .env file")
	}
	logger.CreateLogger()

	logger.Logger.Out.Write([]byte(fmt.Sprintf(`
  .    ___    _    ____    _       
 /\\  / _ \  / \  / ___|  / \    
( ( )| | | |/ _ \ \___ \ / _ \   
 \\/ | |_| / ___ \ ___) / ___ \  
  '   \___/_/   \_\____/_/   \_\ 
                                                                                                  __ _ _    
                                                                                                  \ \ \ \   
		__        _______ ____       _    ____  ____  _     ___ ____    _  _____ ___ ___  _   _    \ \ \ \  
		\ \      / / ____| __ )     / \  |  _ \|  _ \| |   |_ _/ ___|  / \|_   _|_ _/ _ \| \ | |    ) ) ) ) 
		 \ \ /\ / /|  _| |  _ \    / _ \ | |_) | |_) | |    | | |     / _ \ | |  | | | | |  \| |   / / / /  
		  \ V  V / | |___| |_) |  / ___ \|  __/|  __/| |___ | | |___ / ___ \| |  | | |_| | |\  |  / / / /   
		   \_/\_/  |_____|____/  /_/   \_\_|   |_|   |_____|___\____/_/   \_\_| |___\___/|_| \_| /_/_/_/    
                                                                                         


:: OASA WEB APPLICATION (v%s) :: `+"\n\n", os.Getenv("application.version"))))

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
		mapper.NewOasaMapper,
		service.NewRouteService,
		service.NewShedule01Service,
		service.NewLineService,
		service.NewSheduleService,
		service.NewStopService,
		service.NewuVersionService,
		service.NewNotificationService,
		service.NewRestService,
		service.NewOasaService,
		controllers.NewLineController,
		controllers.NewRouteController,
		controllers.NewStopController,
		controllers.NewScheduleController,
		controllers.NewComponentController,
		controllers.TestControllerConstructor,
		controllers.NewOasaNativeController,
		controllers.NewNotifcationController,
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
