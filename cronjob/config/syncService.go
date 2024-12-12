package config

import (
	"fmt"
	"strings"

	"github.com/cs161079/monorepo/common/mapper"
	"github.com/cs161079/monorepo/common/models"
	"github.com/cs161079/monorepo/common/service"
	"github.com/cs161079/monorepo/common/utils"
	logger "github.com/cs161079/monorepo/common/utils/goLogger"
	"gorm.io/gorm"
)

type SyncService interface {
	// =================================================================================================================
	// Με αυτή τη διαδικασία συγχρονίζουμε τα δεδομένα για τις γραμμές των λεοφωρείων από τον Server του OASA στην δική
	// μας βάση δεδομένων. Καλούμε το API /webGetLinesWithMLInfo το οποίο μας επιστρέφει όλες τις γραμμές με την μορφή
	// JSON, δηλαδή ένα πίνακα από record. Η κάθε γραμμής περιέχει την εξής πληροφορία
	//	{
	//	    "ml_code": "9",
	//	    "sdc_code": "54",
	//	    "line_code": "1151",
	//	    "line_id": "021",
	//	    "line_descr": "ΠΛΑΤΕΙΑ ΚΑΝΙΓΓΟΣ - ΓΚΥΖH (ΚΥΚΛΙΚΗ)",
	//	    "line_descr_eng": "PLATEIA KANIGKOS - GKIZI",
	//	    "mld_master": "1"
	//	}
	// =================================================================================================================
	syncLines() error
	// =================================================================================================================
	// Με αυτή τη διαδικασία συγχρονίζουμε τα δεδομένα των διαδρομών από τον Server του OASA στην δική μας βάση δεδομένων.
	// Καλούμε το API /getRoutes το οποίο μας επιστρέφει όλες τις διαδρομές σε txt μορφή, τα δεδομένα των διαδρομών
	// χωρισμένα με κόμμα και κάθε διαδρομή από την άλλη χωρίζονται με κόμμα. Αυτά είναι τα δεδομένα μίας διαδρομής.
	//
	// (1754,799, "ΕΛ.ΒΕΝΙΖΕΛΟΥ - ΚΑΙΣΑΡΙΑΝΗ", "EL. VENIZELOU - KAISARIANI",2,9889.61)
	// =================================================================================================================
	syncRoutes() error
	// =================================================================================================================
	// Με αυτή τη διαδικασία συγχρονίζουμε τα δεδομένα των στάσεων από τον Server του OASA στην δική μας βάση δεδομένων.
	// Καλούμε το API /getStops το οποίο μας επιστρέφει όλες τις στάσεις σε txt μορφή, τα δεδομένα των στάσεων
	// χωρισμένα με κόμμα και κάθε στάση από την άλλη χωρίζονται ομοιώς με κόμμα. Αυτά είναι τα δεδομένα μίας στάσης.
	//
	// (10001, "010001", "ΣΤΡΟΦΗ", "STROFH", "ΕΛ.ΒΕΝΙΖΕΛΟΥ", "ΕΛ.ΒΕΝΙΖΕΛΟΥ", -1,23.665,37.9986,0,0,
	//                       "| ΑΝΩ ΑΓ. ΒΑΡΒΑΡΑ| ΠΕΙΡΑΙΑΣ ΠΛ. ΚΑΡΑΪΣΚΑΚΗ", "| ANO AG. BARBARA| PEIRAIAS PL. KARAISKAKΗ")
	// =================================================================================================================
	syncStops() error
	// =================================================================================================================
	// Με αυτή τη διαδικασία συγχρονίζουμε τα δεδομένα των στάσεων ανά διαδρομή από τον Server του OASA στην δική μας
	// βάση δεδομένων. Καλούμε το API /getRouteStops το οποίο μας επιστρέφει όλες τις στάσεις σε txt μορφή, τα δεδομένα
	// των στάσεων ανά διαδρομή είναι χωρισμένα με κόμμα και κάθε γραμμή είναι και μία εγγραφή στην βάση.
	// Αυτά είναι τα δεδομένα μίας εγγραφής.
	//
	//	(103406,2081,10373,1)
	// =================================================================================================================
	syncRouteStops() error

	syncRouteDetails() error
	uVersionFromOasa() ([]models.UVersions, error)
	SyncData() error
}

type syncService struct {
	dbConnection     *gorm.DB
	restService      service.RestService
	lineService      service.LineService
	routeService     service.RouteService
	stopService      service.StopService
	uVversionService service.UVersionService
}

func NewSyncService(dbConnection *gorm.DB, restSrv service.RestService, lineSrv service.LineService,
	routeSrv service.RouteService, stopSrv service.StopService, uvVerSrv service.UVersionService) SyncService {
	return syncService{
		dbConnection:     dbConnection,
		restService:      restSrv,
		lineService:      lineSrv,
		routeService:     routeSrv,
		stopService:      stopSrv,
		uVversionService: uvVerSrv,
	}
}

func recPreparation(recStr string) string {
	var trimmedSpace = strings.ReplaceAll(recStr, " ", "")
	return strings.ReplaceAll(trimmedSpace, "\"", "")
}

func (s syncService) uVersionFromOasa() ([]models.UVersions, error) {
	response := s.restService.OasaRequestApi00("getUVersions", nil)
	if response.Error != nil {
		return nil, response.Error
	}
	var mapper = mapper.NewUVersionMapper()
	var result []models.UVersions = make([]models.UVersions, 0)
	for _, rec := range response.Data.([]interface{}) {
		result = append(result, mapper.OasaToUVersions(mapper.GeneralUVersions(rec)))
	}
	return result, nil
}

func (s syncService) SyncData() error {
	// *********** Κάνουμε get το connection της  βάσης από το Context ************
	//var dbConnection *gorm.DB = ctx.Value(db.CONNECTIONVAR).(*gorm.DB)
	// **************************************************************************
	versionsArr, err := s.uVersionFromOasa()
	if err != nil {
		return err
	}
	uvServ := s.uVversionService
	routeDetailMustUpdate := false
	for _, rec := range versionsArr {
		if uvServ == nil {
			logger.INFO("The UVersion Service is Nil!!!")
		}
		dbRec, err := uvServ.Select(rec.Uv_descr)
		if err != nil {
			return nil
		}

		if dbRec == nil || rec.Uv_lastupdatelong > dbRec.Uv_lastupdatelong {
			switch rec.Uv_descr {
			case "LINES":
				logger.INFO("########### Lines will be updated...")
				if err := s.syncLines(); err != nil {
					return err
				}
				// Εδώ θα πρέπει να κάνουμε Update την εγγραφή στον πίνακα με το νέο Version.
				uvServ.Post(&rec)
			case "ROUTE STOPS":
				routeDetailMustUpdate = true
				logger.INFO("########### Route Stops will be updated...")
				if err := s.syncRouteStops(); err != nil {
					return err
				}
				// Εδώ θα πρέπει να κάνουμε Update την εγγραφή στον πίνακα με το νέο Version.
				uvServ.Post(&rec)
			case "ROUTES":
				routeDetailMustUpdate = true
				logger.INFO("########### Routes will be updated...")
				if err := s.syncRoutes(); err != nil {
					return err
				}
				// Εδώ θα πρέπει να κάνουμε Update την εγγραφή στον πίνακα με το νέο Version.
				uvServ.Post(&rec)
			case "STOPS":
				logger.INFO("########### Stops will be updated...")
				if err := s.syncStops(); err != nil {
					return err
				}
				// Εδώ θα πρέπει να κάνουμε Update την εγγραφή στον πίνακα με το νέο Version.
				uvServ.Post(&rec)
			}
		}
	}
	if routeDetailMustUpdate {
		logger.INFO("########### Route Details will be updated...")
		if err := s.syncRouteDetails(); err != nil {
			return err
		}
	}
	return nil

}

func (s syncService) syncLines() error {
	// *********** Κάνουμε get το connection της  βάσης από το Context ************
	//var dbConnection *gorm.DB = ctx.Value(db.CONNECTIONVAR).(*gorm.DB)
	// **************************************************************************

	lineSrv := s.lineService //service.NewLineService(repository.NewLineRepository(dbConnection))
	var restSrv = s.restService

	response := restSrv.OasaRequestApi00("webGetLinesWithMLInfo", nil)
	if response.Error != nil {
		return response.Error
	}
	txt := s.dbConnection.Begin()
	if err := lineSrv.WithTrx(txt).DeleteAll(); err != nil {
		txt.Rollback()
	}
	logger.INFO("Delete all data from Line table in database succesfully.")
	var lineArray []models.Line = make([]models.Line, 0)
	logger.INFO("Start sychronize data from OASA Server...")
	for _, ln := range response.Data.([]any) {
		lineOasa := lineSrv.GetMapper().GeneralLine(ln.(map[string]interface{}))
		line := lineSrv.GetMapper().OasaToLine(lineOasa)

		lineArray = append(lineArray, line)
		if len(lineArray) == 1000 {
			_, err := lineSrv.WithTrx(txt).InsertArray(lineArray)
			if err != nil {
				txt.Rollback()
				return err
			}
			logger.INFO(fmt.Sprintf("Batch of data size %d saved succesfully.", len(lineArray)))
			lineArray = make([]models.Line, 0)
		}

	}

	if len(lineArray) > 0 {
		_, err := lineSrv.WithTrx(txt).InsertArray(lineArray)
		if err != nil {
			txt.Rollback()
			return err
		}
		logger.INFO(fmt.Sprintf("Batch of data size %d saved succesfully.", len(lineArray)))
	}

	txt.Commit()
	logger.INFO("Finished sychronize data from OASA Server.")
	return nil
}

func (s syncService) syncRoutes() error {
	// *********** Κάνου get το connection της  βάσης από το Context ************
	//var dbConnection *gorm.DB = ctx.Value(db.CONNECTIONVAR).(*gorm.DB)
	// **************************************************************************
	var restSrv = s.restService //service.NewRestService()

	var routeSrv = s.routeService
	response := restSrv.OasaRequestApi02("getRoutes")
	if response.Error != nil {
		return response.Error
	}

	tx := s.dbConnection.Begin()

	err := routeSrv.WithTrx(tx).DeleteAll()
	if err != nil {
		tx.Rollback()
		return err
	}
	var routeArray []models.Route = make([]models.Route, 0)
	// Εδώ η διαδικασία μας γυρνάει από το API έναν πίνακα με τα Record σε γραμμή χωρισμένα τα πεδία με κόμμα
	for _, rec := range response.Data.([]string) {
		// ************** Κάθε γραμμή την κάνω Split με το κόμμα και γεμίζω τα Record των διαδρομών **************
		// ************************* Έλεγχος της γραμμής εάν έχει όλη την πληροφορία *****************************
		recordArr := strings.Split(recPreparation(rec), ",")
		if len(recordArr) < 6 {
			return fmt.Errorf("Η γραμμή του Record  είναι ελλειπής.")
		}
		rt := models.Route{}
		num, err := utils.StrToInt32(recordArr[0])
		if err != nil {
			return err
		}
		rt.Route_Code = *num
		num, err = utils.StrToInt32(recordArr[1])
		if err != nil {
			return err
		}
		rt.Line_Code = *num
		rt.Route_Descr = recordArr[2]
		rt.Route_Descr_eng = recordArr[3]
		num, err = utils.StrToInt32(recordArr[4])
		if err != nil {
			return err
		}
		rt.Route_Type = int8(*num)
		fl32 := utils.StrToFloat32(recordArr[5])
		rt.Route_Distance = fl32

		routeArray = append(routeArray, rt)

		if len(routeArray) == 10000 {
			// Εδώ Θα καλούμε την Insert για να κάνουμε εγγραφή στην βάση
			_, err = routeSrv.WithTrx(tx).InsertArray(routeArray)
			if err != nil {
				tx.Rollback()
				return err
			}
			logger.INFO(fmt.Sprintf("Batch of data size %d saved succesfully.", len(routeArray)))
		}

	}

	if len(routeArray) > 0 {
		// Εδώ Θα καλούμε την Insert για να κάνουμε εγγραφή στην βάση
		_, err = routeSrv.WithTrx(tx).InsertArray(routeArray)
		if err != nil {
			tx.Rollback()
			return err
		}
		logger.INFO(fmt.Sprintf("Batch of data size %d saved succesfully.", len(routeArray)))
	}
	tx.Commit()

	return nil
}

func (s syncService) syncStops() error {
	// *********** Κάνου get το connection της  βάσης από το Context ************
	//var dbConnection *gorm.DB = ctx.Value(db.CONNECTIONVAR).(*gorm.DB)
	// **************************************************************************
	var restSrv = s.restService //service.NewRestService()

	var stopSrv = s.stopService
	response := restSrv.OasaRequestApi02("getStops")
	if response.Error != nil {
		return response.Error
	}

	tx := s.dbConnection.Begin()

	err := stopSrv.WithTrx(tx).DeleteAll()
	if err != nil {
		tx.Rollback()
		return err
	}
	logger.INFO("Η διαγραφή των στάσεων έγινε με επιτυχία.")

	var stopArr []models.Stop = make([]models.Stop, 0)
	// Εδώ η διαδικασία μας γυρνάει από το API έναν πίνακα με τα Record σε γραμμή χωρισμένα τα πεδία με κόμμα
	logger.INFO("Έναρξη συγχρονισμού δεδομένων...")
	for _, rec := range response.Data.([]string) {
		// ************** Κάθε γραμμή την κάνω Split με το κόμμα και γεμίζω τα Record των στάσεων **************
		// ************************* Έλεγχος της γραμμής εάν έχει όλη την πληροφορία *****************************
		recordArr := strings.Split(recPreparation(rec), ",")
		if len(recordArr) < 13 {
			return fmt.Errorf("Τα δεδομένα της γραμμής είναι ελλειπής.")
		}
		rt := models.Stop{}
		num32, err := utils.StrToInt32(recordArr[0])
		if err != nil {
			return err
		}
		rt.Stop_code = *num32
		rt.Stop_id = recordArr[1]
		rt.Stop_descr = recordArr[2]
		rt.Stop_descr_eng = recordArr[3]
		rt.Stop_street = recordArr[4]
		rt.Stop_street_eng = recordArr[5]
		num32, err = utils.StrToInt32(recordArr[6])
		if err != nil {
			return err
		}
		rt.Stop_heading = *num32
		rt.Stop_lng = utils.StrToFloat(recordArr[7])
		rt.Stop_lat = utils.StrToFloat(recordArr[8])
		num8, err := utils.StrToInt8(recordArr[9])
		if err != nil {
			return err
		}
		rt.Stop_type = *num8
		num8, err = utils.StrToInt8(recordArr[10])
		if err != nil {
			return err
		}
		rt.Stop_amea = *num8
		rt.Destinations = recordArr[11]
		rt.Destinations_Eng = recordArr[12]

		if len(stopArr) == 1000 {
			// Εδώ Θα καλούμε την Insert για να κάνουμε εγγραφή στην βάση
			_, err = stopSrv.WithTrx(tx).InsertArray(stopArr)
			if err != nil {
				tx.Rollback()
				return err
			}
			logger.INFO(fmt.Sprintf("Batch of data size %d saved succesfully.", len(stopArr)))
			stopArr = make([]models.Stop, 0)
		}

	}

	if len(stopArr) > 0 {
		// Εδώ Θα καλούμε την Insert για να κάνουμε εγγραφή στην βάση
		_, err = stopSrv.WithTrx(tx).InsertArray(stopArr)
		if err != nil {
			tx.Rollback()
			return err
		}
		logger.INFO(fmt.Sprintf("Batch of data size %d saved succesfully.", len(stopArr)))
	}

	tx.Commit()

	return nil
}

func (s syncService) syncRouteStops() error {
	// *********** Παίρνουμε το connection από το Context της εφαρμογής ************
	//var dbConnection *gorm.DB = ctx.Value(db.CONNECTIONVAR).(*gorm.DB)
	// *****************************************************************************

	// Δημιουργία ενός Rest Service για να κάνω την κλήση στον Server
	var restSrv = s.restService

	var routeSrv = s.routeService
	response := restSrv.OasaRequestApi02("getRouteStops")
	if response.Error != nil {
		return response.Error
	}

	tx := s.dbConnection.Begin()
	err := routeSrv.WithTrx(tx).DeleteRoute02()
	if err != nil {
		tx.Rollback()
		return err
	}
	logger.INFO("Delete from Route02 Table succesfully.")

	var route02Arr []models.Route02 = make([]models.Route02, 0)
	logger.INFO("Start sychronization Route02 data from OASA Server...")
	for _, rec := range response.Data.([]string) {
		row := strings.Split(recPreparation(rec), ",")
		if len(row) < int(4) {
			return fmt.Errorf("Τα δεδομένα της γραμμής είναι ελλειπής.")
		}

		rt := models.Route02{}
		num32, err := utils.StrToInt32(row[1])
		if err != nil {
			return err
		}
		rt.Route_code = *num32
		num64, err := utils.StrToInt64(row[2])
		rt.Stop_code = *num64
		num16, err := utils.StrToInt16(row[3])
		rt.Senu = *num16

		route02Arr = append(route02Arr, rt)
		if len(route02Arr) == 10000 {
			_, err = routeSrv.WithTrx(tx).Route02InsertArr(route02Arr)
			if err != nil {
				tx.Rollback()
				return err
			}
			logger.INFO(fmt.Sprintf("Batch of data size %d saved succesfully.", len(route02Arr)))
			route02Arr = make([]models.Route02, 0)
		}
	}
	if len(route02Arr) > 0 {
		_, err = routeSrv.WithTrx(tx).Route02InsertArr(route02Arr)
		if err != nil {
			tx.Rollback()
			return err
		}
		logger.INFO(fmt.Sprintf("Batch of data size %d saved succesfully.", len(route02Arr)))
	}
	tx.Commit()
	logger.INFO("Finished sychronization Route02 data from OASA Server.")
	return nil
}

func (s syncService) syncRouteDetails() error {
	// *********** Παίρνουμε το connection από το Context της εφαρμογής ************
	//var dbConnection *gorm.DB = ctx.Value(db.CONNECTIONVAR).(*gorm.DB)
	// *****************************************************************************

	// Δημιουργία ενός Rest Service για να κάνω την κλήση στον Server
	var restSrv = s.restService

	var routeSrv = s.routeService
	var allRoutes, err = routeSrv.List01()
	if err != nil {
		return err
	}

	var allRouteDetails []models.Route01 = make([]models.Route01, 0)

	var tx = s.dbConnection.Begin()

	if err := routeSrv.WithTrx(tx).DeleteRoute01(); err != nil {
		tx.Rollback()
		return err
	}

	for _, rec := range allRoutes {
		response := restSrv.OasaRequestApi00("webRouteDetails",
			map[string]interface{}{
				"p1": int64(rec.Route_Code),
			},
		)
		if response.Error != nil {
			return response.Error
		}

		// Είναι Array από interfaced{} τα οποία είναι map[stirng]interface{}
		for _, j := range response.Data.([]interface{}) {
			var route01Oasa = routeSrv.GetMapper01().GeneralRoute01(j.(map[string]interface{}))
			var route01 = routeSrv.GetMapper01().OasaToRoute01Dto(route01Oasa)
			route01.Route_code = rec.Route_Code
			allRouteDetails = append(allRouteDetails, route01)
			if len(allRouteDetails) == 10000 {
				if _, err := routeSrv.WithTrx(tx).Route01InsertArr(allRouteDetails); err != nil {
					tx.Rollback()
					return err
				}
				logger.INFO(fmt.Sprintf("Batch of data size %d saved succesfully.", len(allRouteDetails)))
				allRouteDetails = make([]models.Route01, 0)
			}
		}
	}
	if len(allRouteDetails) > 0 {
		if _, err := routeSrv.Route01InsertArr(allRouteDetails); err != nil {
			tx.Rollback()
			return err
		}
		logger.INFO(fmt.Sprintf("Batch of data size %d saved succesfully.", len(allRouteDetails)))
	}

	tx.Commit()
	return nil
}
