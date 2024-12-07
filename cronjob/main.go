package main

import (
	"context"
	"fmt"
	"os"
	"sync"
	"time"

	"github.com/cs161079/monorepo/common/db"
	models "github.com/cs161079/monorepo/common/models"
	logger "github.com/cs161079/monorepo/common/utils/goLogger"
	"github.com/joho/godotenv"
	"gorm.io/gorm"
)

func initEnviroment() {
	// loads values from .env into the system
	if err := godotenv.Load(".env"); err != nil {
		logger.ERROR("No .env file found")
	}
}

func getEnv(key string, defaultVal string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}

	return defaultVal
}

func worker(id int, connection *gorm.DB, wg *sync.WaitGroup) {
	defer func() {
		fmt.Printf("--------WORKER %d DONE-----------------------\n\n", id)
		wg.Done()
	}()
	fmt.Printf("------------------------------------------\n")
	fmt.Printf("Created worker with id %d \n", id)
	var res models.Line = models.Line{}
	cn := connection.Table("line").Where("line_code=?", 1151).Find(&res)
	if cn.Error != nil {
		panic(cn.Error.Error())
	}
	c, err := connection.DB()
	if err != nil {
		panic(err.Error())
	}
	openConns := c.Stats().OpenConnections // Total number of open connections
	inUseConns := c.Stats().InUse          // Connections currently in use
	idleConns := c.Stats().Idle            // Idle (unused) connections in the pool

	fmt.Println("----------- Connection Pool Stats ---------------")
	fmt.Printf("Open connections: %d\n", openConns)
	fmt.Printf("In-use connections: %d\n", inUseConns)
	fmt.Printf("Idle connections: %d\n", idleConns)
	fmt.Println("----------- End Connection Pool Stats -----------")
	time.Sleep(10 * time.Second)
}

// Αυτό είναι ένα τεστ που είχα κάνει με Worker
// Δεν θα το χρησιμοποιήσω.
func main_worker() {
	var wg sync.WaitGroup // Create a WaitGroup

	emf, err := db.CreateConnection()
	if err != nil {
		logger.ERROR(err.Error())
		return
	} // Increment the WaitGroup counter
	logger.InitLogger("goSyncApplication")

	// Launch several goroutines
	for i := 1; i <= 7; i++ {
		wg.Add(1)

		// em, err := emf.DB()
		if err != nil {
			logger.ERROR(err.Error())
			// return
		} else {
			go worker(i, emf, &wg) // Start a goroutine
		}

	}

	// Wait for all goroutines to finish
	wg.Wait()
	fmt.Println("All workers done")
}

func main() {
	logger.InitLogger("goSyncApplication")

	// ********* Κάνουμε Create το connection Με την βάση δεδομένων **************
	dbConnection, err := db.CreateConnection()
	if err != nil {
		logger.ERROR(fmt.Sprintf("Database not established. [%s]", err.Error()))
		return
	}
	// ***************************************************************************
	// ******************* Δημιουργόυμε ένα parent Context ***********************
	var parentContext = context.Background()
	// ******* Δημιουργούμε ένα Context με Value το connection στην Βάση *********
	var applicationContext = context.WithValue(parentContext, db.CONNECTIONVAR, dbConnection)
	// ***************************************************************************

	// ******** Δημιουργία Service for Sychronization *************
	syncSrv := NewSyncService()

	err = syncSrv.SyncData(applicationContext)
	if err != nil {
		logger.ERROR(fmt.Sprintf("An error occurred on data sychronization from Server \n[%s]", err.Error()))
	}

}
