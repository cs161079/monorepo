package config

import (
	"fmt"
	"os"

	"github.com/cs161079/monorepo/common/db"
	"github.com/cs161079/monorepo/common/repository"
	"github.com/cs161079/monorepo/common/service"
	logger "github.com/cs161079/monorepo/common/utils/goLogger"
	"github.com/cs161079/monorepo/webApplication/controllers"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/go-errors/errors"
	"github.com/joho/godotenv"
	"go.uber.org/dig"
	"gorm.io/gorm"
)

type App struct {
	engine *gin.Engine
}

func ErrorHandler(c *gin.Context, err any) {
	// Wrap the error with stack trace
	var wrappedErr error
	switch e := err.(type) {
	case error:
		wrappedErr = errors.Wrap(e, 1)
	default:
		wrappedErr = errors.New("unknown error occurred")
	}

	// Log the error with context
	logger.ERROR(fmt.Sprintln("Error occurred",
		"error", wrappedErr.Error(),
		"stack", wrappedErr.(*errors.Error).ErrorStack(),
		"path", c.Request.URL.Path,
		"method", c.Request.Method,
		"clientIP", c.ClientIP(),
	))
	var httpResponse = map[string]any{"Message": "Internal server error", "Status": 500}
	c.AbortWithStatusJSON(500, httpResponse)
}

func NewApp(db *gorm.DB, lineCtrl controllers.LineControllerImplementation,
	testCtrl controllers.TestController) *App {
	gin.SetMode(gin.ReleaseMode)
	eng := gin.New()
	eng.Use(cors.Default())
	gin.DefaultWriter = logger.Logger.Out
	gin.DefaultErrorWriter = logger.Logger.Out
	eng.Use(gin.Logger(), gin.CustomRecovery(ErrorHandler))

	lineCtrl.AddRouters(eng)
	testCtrl.AddRoutes(eng)

	return &App{
		engine: eng,
	}
}

func (a App) Boot() {
	var port = os.Getenv("SERVER_PORT")
	if port == "" {
		port = "8080"
	}
	logger.Logger.Info("Application server start on port %s", port)
	a.engine.Run(fmt.Sprintf(":%s", port))
}

func InitializeApplication() {
	// Load the .env file
	err := godotenv.Load()
	if err != nil {
		fmt.Println("Error loading .env file")
	}
	logger.InitLogger("WebApplication")
	originalStdout := os.Stdout

	os.Stdout = logger.Logger.Out.(*os.File) // Set output destination
	fmt.Printf(`
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
                                                                                         


:: OASA WEB APPLICATION (v1.0.0) ::

`)
	os.Stdout = originalStdout

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
		service.NewRouteService,
		service.NewShedule01Service,
		service.NewSheduleService,
		service.NewStopService,
		service.NewuVersionService,
		controllers.NewLineController,
		controllers.TestControllerConstructor,
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
