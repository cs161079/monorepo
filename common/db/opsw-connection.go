package db

import (
	"fmt"
	"os"
	"strconv"
	"time"

	logger "github.com/cs161079/monorepo/common/utils/goLogger"
	"github.com/joho/godotenv"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	gormlogger "gorm.io/gorm/logger"
)

type Datasource interface {
	DatasourceUrl() (string, error)
}

type datasource struct {
	Address      *string
	Port         *int32
	User         string
	Password     string
	DatabaseName string
}

func (ds datasource) DatasourceUrl() string {
	return fmt.Sprintf(dataSourceFormat, ds.User, ds.Password,
		*ds.Address, *ds.Port, ds.DatabaseName)

}

const (
	// *ΓΙΑ ΤΙΣ ΜΕΤΑΒΛΗΤΕΣ ΠΟΥ ΔΙΑΜΟΙΡΑΖΟΜΑΙ ΜΕΣΩ CONTEXT *
	CONNECTIONVAR = "db_conn"
	//*****************************************************
	// ****** ΟΝΟΜΑΤΑ ΠΙΝΑΚΩΝ ΣΤΗΝ ΒΑΣΗ ΔΕΔΟΜΕΝΩΝ *********
	LINETABLE           = "line"
	ROUTETABLE          = "route"
	STOPTABLE           = "stop"
	ROUTESTOPSTABLE     = "route02"
	SCHEDULEMASTERTABLE = "schedulemaster"
	SCHEDULETIMETABLE   = "scheduletime"
	SCHEDULELINE        = "scheduleline"
	ROUTEDETAILTABLE    = "route01"
	SYNCVERSIONSTABLE   = "syncversions"
	// ****************************************************
)

const dataSourceFormat = "%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=true&loc=Local"

func createDataSource() (*datasource, error) {
	// ******** Φορτώνουμε τις μεταβλητές ************************
	err := godotenv.Load(".env")
	if err != nil {
		return nil, fmt.Errorf("Error loading enviroment variables.[%s]", err.Error())
	}

	// ******** Get Database IP from env, if empty put a default IP ******************
	ip := os.Getenv("database.ip")
	var defaultIp = "127.0.0.1"
	if ip == "" {
		ip = defaultIp
	}
	// *******************************************************************************

	// ************* Get Database port from env, if empty put a default port *********
	port, err := strconv.ParseInt(os.Getenv("database.port"), 10, 32)
	if err != nil {
		return nil, fmt.Errorf("Error converting database.port variable.[%s]", err.Error())
	}

	var defaultPort int32 = 3306

	if port == 0 {
		port = int64(defaultPort)
	}
	// ********************************************************************************

	// ************ Rest information for database connection **************************
	user := os.Getenv("database.user")
	password := os.Getenv("database.password")
	database := os.Getenv("database.dbname")
	// ********************************************************************************

	var port32 = int32(port)
	return &datasource{
		Address:      &ip,
		Port:         &port32,
		User:         user,
		Password:     password,
		DatabaseName: database,
	}, nil
}

func getGormConfig() *gorm.Config {
	gormLogger := logger.GetGormLogger()
	if gormLogger == nil {
		return &gorm.Config{
			Logger: gormlogger.Default.LogMode(gormlogger.Silent),
		}
	}
	return &gorm.Config{
		Logger: gormLogger,
	}
}

// This is core for DB

func NewOpswConnection() (*gorm.DB, error) {
	dataSource, err := createDataSource()
	if err != nil {
		return nil, err
	}

	dialector := mysql.New(mysql.Config{
		DSN:                       dataSource.DatasourceUrl(), // data source name
		DefaultStringSize:         256,                        // default size for string fields
		DisableDatetimePrecision:  true,                       // disable datetime precision, which not supported before MySQL 5.6
		DontSupportRenameIndex:    true,                       // drop & create when rename index, rename index not supported before MySQL 5.7, MariaDB
		DontSupportRenameColumn:   true,                       // `change` when rename column, rename column not supported before MySQL 8, MariaDB
		SkipInitializeWithVersion: false,                      // auto configure based on currently MySQL version
	})

	db, err := gorm.Open(dialector, getGormConfig())

	if err != nil {
		// fmt.Println("An Error occured on creation of database connection")
		return nil, err
	}

	sqlDb, err := db.DB()
	if err != nil {
		return nil, err
	}
	sqlDb.SetMaxIdleConns(5)
	sqlDb.SetConnMaxLifetime(time.Minute)
	sqlDb.SetMaxOpenConns(5)

	return db, nil
}
