package config

import (
	"fmt"
	"sort"
	"strings"

	"github.com/cs161079/monorepo/common/mapper"
	"github.com/cs161079/monorepo/common/models"
	"github.com/cs161079/monorepo/common/service"
	"github.com/cs161079/monorepo/common/utils"
	logger "github.com/cs161079/monorepo/common/utils/goLogger"
	"github.com/cs161079/monorepo/cronjob/dao"
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

	SyncSchedule() error

	syncRouteDetails() error
	uVersionFromOasa() ([]dao.UVersion01, error)
	SyncData() error
	DeleteAll() error
	InserttoDatabase() error
}

type syncService struct {
	HelpLine          []models.Line
	HelpRoute         []models.Route
	HelpRoute01       []models.Route01
	HelpRoute02       []models.Route02
	HelpStop          []models.Stop
	HelpSchedule      []models.Schedule
	HelpScheduletime  []models.Scheduletime
	HelpScheduleline  []models.Scheduleline
	dbConnection      *gorm.DB
	restService       service.RestService
	lineService       service.LineService
	routeService      service.RouteService
	stopService       service.StopService
	uVversionService  service.UVersionService
	scheduleService   service.ScheduleService
	schedule01Service service.Schedule01Service
	// Εδώ κρατάμε τα κλειδιά των διαδρομών που έχουν συγχρονιστεί
	// γιατί μου φέρνει Detail διαδρομών οι οποίες δεν υπάρχουν.
	routeKeys         map[int32]int32
	scheduleMasterKey map[int32]int32
}

func NewSyncService(dbConnection *gorm.DB, restSrv service.RestService, lineSrv service.LineService,
	routeSrv service.RouteService, stopSrv service.StopService, uvVerSrv service.UVersionService,
	schedule service.ScheduleService, schedule01 service.Schedule01Service) SyncService {
	return &syncService{
		dbConnection:      dbConnection,
		restService:       restSrv,
		lineService:       lineSrv,
		routeService:      routeSrv,
		stopService:       stopSrv,
		scheduleService:   schedule,
		schedule01Service: schedule01,
		uVversionService:  uvVerSrv,
		routeKeys:         make(map[int32]int32),
		scheduleMasterKey: make(map[int32]int32),
	}
}

func recPreparation(recStr string) string {
	var trimmedSpace = strings.ReplaceAll(recStr, " ", "")
	return strings.ReplaceAll(trimmedSpace, "\"", "")
}

func FixRecordOrder(rec *dao.UVersion01) {
	rec.Orderd = 9999
	if rec.UVersion.Uv_descr == "LINES" {
		rec.Orderd = 0
	} else if rec.UVersion.Uv_descr == "ROUTES" {
		rec.Orderd = 1
	} else if rec.UVersion.Uv_descr == "STOPS" {
		rec.Orderd = 2
	} else if rec.UVersion.Uv_descr == "ROUTE STOPS" {
		rec.Orderd = 3
	} else if rec.UVersion.Uv_descr == "ROUTE DETAIL" {
		rec.Orderd = 4
	} else if rec.UVersion.Uv_descr == "SCHED_CATS" {
		rec.Orderd = 5
	}
}

func (s *syncService) InserttoDatabase() error {
	// Εισαγωγή γραμμών
	txt := s.dbConnection.Begin()
	if err := s.lineService.WithTrx(txt).InsertChunkArray(1000, s.HelpLine); err != nil {
		txt.Rollback()
		return err
	}

	if err := txt.Commit().Error; err != nil {
		return err
	}
	logger.INFO("Η εισαγωγή γραμμών στην βάση δεδομένων ολοκληρώθηκε.")
	// Εισαγωγή διαδρομών
	txt = s.dbConnection.Begin()
	if err := s.routeService.WithTrx(txt).InserChunkArray(10000, s.HelpRoute); err != nil {
		txt.Rollback()
		return err
	}
	logger.INFO("Η εισαγωγή διαδρομών στην βάση δεδομένων ολοκληρώθηκε.")
	if err := txt.Commit().Error; err != nil {
		return err
	}
	txt = s.dbConnection.Begin()
	if err := s.routeService.WithTrx(txt).Route01InsertChunkArray(10000, s.HelpRoute01); err != nil {
		txt.Rollback()
		return err
	}
	logger.INFO("Η εισαγωγή λεπτομερειών διαδρομών στην βάση δεδομένων ολοκληρώθηκε.")
	if err := txt.Commit().Error; err != nil {
		return err
	}
	// Εισαγωγή Στάσεων
	txt = s.dbConnection.Begin()
	if err := s.stopService.WithTrx(txt).InsertChunkArray(1000, s.HelpStop); err != nil {
		txt.Rollback()
		return err
	}
	logger.INFO("Η εισαγωγή στάσεων στην βάση δεδομένων ολοκληρώθηκε.")
	if err := txt.Commit().Error; err != nil {
		return err
	}
	// Εισαγωγή συχετισμένω στάσεων ανα διαδρομή.
	txt = s.dbConnection.Begin()
	if err := s.routeService.WithTrx(txt).Route02InsertChunkArray(15000, s.HelpRoute02); err != nil {
		txt.Rollback()
		return err
	}
	logger.INFO("Η εισαγωγή εγγραφών στάσεων ανά διαδρομή στην βάση δεδομένων ολοκληρώθηκε.")
	if err := txt.Commit().Error; err != nil {
		return err
	}

	txt = s.dbConnection.Begin()
	if err := s.scheduleService.WithTrx(txt).InsertScheduleChunkArray(10000, s.HelpSchedule); err != nil {
		return err
	}
	logger.INFO("Η εισαγωγή το δρομολογίων στην βάση δεδομένων ολοκληρώθηκε.")
	if err := txt.Commit().Error; err != nil {
		return err
	}

	txt = s.dbConnection.Begin()
	if err := s.lineService.WithTrx(txt).InsertChunkSchedulesArray(10000, s.HelpScheduleline); err != nil {
		return err
	}
	logger.INFO("Η εισαγωγή των δρομολογίων ανά γραμμή στην βάση δεδομένων ολοκληρώθηκε.")
	if err := txt.Commit().Error; err != nil {
		return err
	}

	txt = s.dbConnection.Begin()
	if err := s.schedule01Service.WithTrx(txt).InsertSchedule01ChunkArray(10000, s.HelpScheduletime); err != nil {
		return err
	}
	logger.INFO("Η εισαγωγή ωραρίων στην βάση δεδομένων ολοκληρώθηκε.")
	if err := txt.Commit().Error; err != nil {
		return err
	}
	return nil
}

func (s *syncService) DeleteAll() error {
	txt := s.dbConnection.Begin()
	// Διαγραφή εγγραφών σε πίνακα που συχετίζει στάσεις ανα διαδρομή
	if err := s.routeService.WithTrx(txt).DeleteRoute02(); err != nil {
		txt.Rollback()
		return err
	}

	// Διαγραφή ΣΤΑΣΕΩΝ
	if err := s.stopService.WithTrx(txt).DeleteAll(); err != nil {
		txt.Rollback()
		return err
	}

	// Διαγραφή Λεπτομερειών ΔΙΑΔΡΟΜΗΣ
	if err := s.routeService.WithTrx(txt).DeleteRoute01(); err != nil {
		txt.Rollback()
		return err
	}
	// Διαγραφή ΔΙΑΔΡΟΜΩΝ
	if err := s.routeService.WithTrx(txt).DeleteAll(); err != nil {
		txt.Rollback()
		return err
	}

	// Διαγραφή ΩΡΑΡΙΩΝ
	if err := s.schedule01Service.WithTrx(txt).DeleteAll(); err != nil {
		txt.Rollback()
		return err
	}

	if err := s.lineService.WithTrx(txt).DeleteAllLineSchedules(); err != nil {
		txt.Rollback()
		return err
	}

	// Διαγραφή ΔΡΟΜΟΛΟΓΙΩΝ
	if err := s.scheduleService.WithTrx(txt).DeleteAll(); err != nil {
		txt.Rollback()
		return err
	}

	//Διαγραφή ΓΡΑΜΜΩΝ
	if err := s.lineService.WithTrx(txt).DeleteAll(); err != nil {
		txt.Rollback()
		return err
	}
	return txt.Commit().Error
}

func (s *syncService) uVersionFromOasa() ([]dao.UVersion01, error) {
	response := s.restService.OasaRequestApi00("getUVersions", nil)
	if response.Error != nil {
		return nil, response.Error
	}
	var mapper = mapper.NewUVersionMapper()
	var result []dao.UVersion01 = make([]dao.UVersion01, 0)
	for _, rec := range response.Data.([]interface{}) {
		appRec := mapper.OasaToUVersions(mapper.GeneralUVersions(rec))
		appRec01 := dao.UVersion01{
			UVersion: appRec,
		}
		FixRecordOrder(&appRec01)
		result = append(result, appRec01)
	}
	sort.Slice(result, func(i, j int) bool {
		return result[i].Orderd < result[j].Orderd
	})
	return result, nil
}

func filterRecords(records []models.UVersions, condition func(models.UVersions) bool) *models.UVersions {
	var result models.UVersions
	for _, record := range records {
		if condition(record) {
			result = record
			return &result
		}
	}
	return nil
}

func (s *syncService) SyncData() error {
	// *********** Κάνουμε get το connection της  βάσης από το Context ************
	//var dbConnection *gorm.DB = ctx.Value(db.CONNECTIONVAR).(*gorm.DB)
	// **************************************************************************
	versionsArr, err := s.uVersionFromOasa()
	if err != nil {
		return err
	}
	uvServ := s.uVversionService
	// routeDetailMustUpdate := false
	for _, rec := range versionsArr {
		dbRec, err := uvServ.Select(rec.UVersion.Uv_descr)
		if err != nil {
			return nil
		}

		var frServeVers = rec.UVersion.Uv_lastupdatelong
		if dbRec == nil || frServeVers > dbRec.Uv_lastupdatelong {
			switch rec.UVersion.Uv_descr {
			case "LINES":
				logger.INFO("Lines will be updated...")
				if err := s.syncLines(); err != nil {
					return err
				}
				// Εδώ θα πρέπει να κάνουμε Update την εγγραφή στον πίνακα με το νέο Version.
				//uvServ.Post(&rec.UVersion)
			case "ROUTES":
				// routeDetailMustUpdate = true
				logger.INFO("Routes will be updated...")
				if err := s.syncRoutes(); err != nil {
					return err
				}
				// Εδώ θα πρέπει να κάνουμε Update την εγγραφή στον πίνακα με το νέο Version.
				//uvServ.Post(&rec.UVersion)
			case "STOPS":
				logger.INFO("Stops will be updated...")
				if err := s.syncStops(); err != nil {
					return err
				}
				// Εδώ θα πρέπει να κάνουμε Update την εγγραφή στον πίνακα με το νέο Version.
				//uvServ.Post(&rec.UVersion)
			case "ROUTE STOPS":
				// routeDetailMustUpdate = true
				logger.INFO("Stops per Route will be updated...")
				if err := s.syncRouteStops(); err != nil {
					return err
				}
				// Εδώ θα πρέπει να κάνουμε Update την εγγραφή στον πίνακα με το νέο Version.
				//uvServ.Post(&rec.UVersion)
			case "ROUTE DETAIL":
				// routeDetailMustUpdate = false
				logger.INFO("Route details be updated....")
				if err := s.syncRouteDetails(); err != nil {
					return err
				}
			}

		}
	}
	if err = s.SyncSchedule(); err != nil {
		return err
	}
	return nil

}

// func writeResponseToFile(filename string, content []string) {
// 	file, err := os.Create(filename)
// 	if err != nil {
// 		logger.ERROR(err.Error())
// 		return
// 	}
// 	// file := utils.NewOpswfile("SCHED_ENTRIES.txt")
// 	defer file.Close()
// 	for _, sched := range content {
// 		file.WriteString(fmt.Sprintf("%s\n", sched))
// 	}
// }

func (s *syncService) syncLines() error {
	// *********** Κάνουμε get το connection της  βάσης από το Context ************
	//var dbConnection *gorm.DB = ctx.Value(db.CONNECTIONVAR).(*gorm.DB)
	// **************************************************************************

	lineSrv := s.lineService //service.NewLineService(repository.NewLineRepository(dbConnection))
	var restSrv = s.restService

	response := restSrv.OasaRequestApi00("webGetLinesWithMLInfo", nil)
	if response.Error != nil {
		return response.Error
	}
	// TODO: Το έκοψα γιατί δεν θα κάνω εδώ το Delete
	// txt := s.dbConnection.Begin()
	// if err := lineSrv.WithTrx(txt).DeleteAll(); err != nil {
	// 	txt.Rollback()
	// }
	// logger.INFO("Delete all data from Line table in database succesfully.")
	// var lineArray []models.Line = make([]models.Line, 0)
	logger.INFO("Get Route data from OASA Server...")
	s.HelpLine = make([]models.Line, 0)
	logger.INFO("Get line data from OASA Server...")
	for _, ln := range response.Data.([]any) {
		lineOasa := lineSrv.GetMapper().GeneralLine(ln.(map[string]interface{}))
		line := lineSrv.GetMapper().OasaToLine(lineOasa)

		// if _, ok := s.lineKeys[line.Line_Code]; !ok {
		// 	s.lineKeys[line.Line_Code] = line.Line_Code
		// }
		s.HelpLine = append(s.HelpLine, line)
	}

	// TODO: Δεν θα κάνουμε εδώ Insert στην βάση.
	// if len(lineArray) > 0 {
	// 	_, err := lineSrv.WithTrx(txt).InsertArray(lineArray)
	// 	if err != nil {
	// 		txt.Rollback()
	// 		return err
	// 	}
	// 	logger.INFO(fmt.Sprintf("Batch of data size %d saved succesfully.", len(lineArray)))
	// }

	// txt.Commit()
	logger.INFO("Finished sychronize data from OASA Server.")
	return nil
}

func (s *syncService) syncRoutes() error {
	// *********** Κάνου get το connection της  βάσης από το Context ************
	//var dbConnection *gorm.DB = ctx.Value(db.CONNECTIONVAR).(*gorm.DB)
	// **************************************************************************
	var restSrv = s.restService //service.NewRestService()

	// var routeSrv = s.routeService
	response := restSrv.OasaRequestApi02("getRoutes")
	if response.Error != nil {
		return response.Error
	}

	// TODO: Δεν θα κάνουμε εδώ Delete από τον πίνακα
	// tx := s.dbConnection.Begin()

	// err := routeSrv.WithTrx(tx).DeleteAll()
	// if err != nil {
	// 	tx.Rollback()
	// 	return err
	// }
	// Δεν θα χρησιμοποιήσουμε τοπικό πίνακα
	// var routeArray []models.Route = make([]models.Route, 0)
	// Εδώ η διαδικασία μας γυρνάει από το API έναν πίνακα με τα Record σε γραμμή χωρισμένα τα πεδία με κόμμα
	for _, rec := range response.Data.([]string) {
		// ************** Κάθε γραμμή την κάνω Split με το κόμμα και γεμίζω τα Record των διαδρομών **************
		// ************************* Έλεγχος της γραμμής εάν έχει όλη την πληροφορία *****************************
		recordArr := strings.Split(recPreparation(rec), ",")
		if len(recordArr) < 6 {
			return fmt.Errorf("Η γραμμή του Record  είναι ελλειπής.")
		}
		rt := models.Route{}
		num, err := utils.StrToInt32(recordArr[1])
		if err != nil {
			return err
		}
		rt.Ln_Code = *num
		num, err = utils.StrToInt32(recordArr[0])
		if err != nil {
			return err
		}
		rt.Route_Code = *num
		if _, ok := s.routeKeys[rt.Route_Code]; !ok {
			s.routeKeys[rt.Route_Code] = rt.Route_Code
		}

		rt.Route_Descr = recordArr[2]
		rt.Route_Descr_eng = recordArr[3]
		num, err = utils.StrToInt32(recordArr[4])
		if err != nil {
			return err
		}
		rt.Route_Type = int8(*num)
		fl32 := utils.StrToFloat32(recordArr[5])
		rt.Route_Distance = fl32

		s.HelpRoute = append(s.HelpRoute, rt)
		//logger.INFO(fmt.Sprintf("Η διαδρομή [routr %d, line %d] προστεθηκε", rt.Route_Code, rt.Line_Code))

		// TODO: Αλλαγή
		// if len(routeArray) == 10000 {
		// 	// Εδώ Θα καλούμε την Insert για να κάνουμε εγγραφή στην βάση
		// 	_, err = routeSrv.WithTrx(tx).InsertArray(routeArray)
		// 	if err != nil {
		// 		tx.Rollback()
		// 		return err
		// 	}
		// 	logger.INFO(fmt.Sprintf("Batch of data size %d saved succesfully.", len(routeArray)))
		// }

	}

	// if len(routeArray) > 0 {
	// 	// Εδώ Θα καλούμε την Insert για να κάνουμε εγγραφή στην βάση
	// 	_, err = routeSrv.WithTrx(tx).InsertArray(routeArray)
	// 	if err != nil {
	// 		tx.Rollback()
	// 		return err
	// 	}
	// 	logger.INFO(fmt.Sprintf("Batch of data size %d saved succesfully.", len(routeArray)))
	// }
	// tx.Commit()

	return nil
}

func (s *syncService) syncStops() error {
	// *********** Κάνου get το connection της  βάσης από το Context ************
	//var dbConnection *gorm.DB = ctx.Value(db.CONNECTIONVAR).(*gorm.DB)
	// **************************************************************************
	var restSrv = s.restService //service.NewRestService()

	// var stopSrv = s.stopService
	response := restSrv.OasaRequestApi02("getStops")
	if response.Error != nil {
		return response.Error
	}

	//TODO: Δεν θα κάνουμε διαγραφή
	// tx := s.dbConnection.Begin()

	// err := stopSrv.WithTrx(tx).DeleteAll()
	// if err != nil {
	// 	tx.Rollback()
	// 	return err
	// }
	// logger.INFO("Η διαγραφή των στάσεων έγινε με επιτυχία.")

	// var stopArr []models.Stop = make([]models.Stop, 0)
	s.HelpStop = make([]models.Stop, 0)
	// Εδώ η διαδικασία μας γυρνάει από το API έναν πίνακα με τα Record σε γραμμή χωρισμένα τα πεδία με κόμμα
	logger.INFO("Get Stop data from OASA Server...")
	for _, rec := range response.Data.([]string) {
		// ************** Κάθε γραμμή την κάνω Split με το κόμμα και γεμίζω τα Record των στάσεων **************
		// ************************* Έλεγχος της γραμμής εάν έχει όλη την πληροφορία *****************************
		recordArr := strings.Split(recPreparation(rec), ",")
		if len(recordArr) < 13 {
			return fmt.Errorf("Τα δεδομένα της γραμμής είναι ελλειπής.")
		}
		st := models.Stop{}
		num32, err := utils.StrToInt32(recordArr[0])
		if err != nil {
			return err
		}
		st.Stop_code = *num32
		st.Stop_id = recordArr[1]
		st.Stop_descr = recordArr[2]
		st.Stop_descr_eng = recordArr[3]
		st.Stop_street = recordArr[4]
		st.Stop_street_eng = recordArr[5]
		num32, err = utils.StrToInt32(recordArr[6])
		if err != nil {
			return err
		}
		st.Stop_heading = *num32
		st.Stop_lng = utils.StrToFloat(recordArr[7])
		st.Stop_lat = utils.StrToFloat(recordArr[8])
		num8, err := utils.StrToInt8(recordArr[9])
		if err != nil {
			return err
		}
		st.Stop_type = *num8
		num8, err = utils.StrToInt8(recordArr[10])
		if err != nil {
			return err
		}
		st.Stop_amea = *num8
		st.Destinations = recordArr[11]
		st.Destinations_Eng = recordArr[12]

		s.HelpStop = append(s.HelpStop, st)

		// TODO: Αλλαγή δεν θα γίνετα
		// if len(stopArr) == 1000 {
		// 	// Εδώ Θα καλούμε την Insert για να κάνουμε εγγραφή στην βάση
		// 	_, err = stopSrv.WithTrx(tx).InsertArray(stopArr)
		// 	if err != nil {
		// 		tx.Rollback()
		// 		return err
		// 	}
		// 	logger.INFO(fmt.Sprintf("Batch of data size %d saved succesfully.", len(stopArr)))
		// 	stopArr = make([]models.Stop, 0)
		// }

	}
	logger.INFO("Finished sycronization from OASA Server...")

	// if len(stopArr) > 0 {
	// 	// Εδώ Θα καλούμε την Insert για να κάνουμε εγγραφή στην βάση
	// 	_, err = stopSrv.WithTrx(tx).InsertArray(stopArr)
	// 	if err != nil {
	// 		tx.Rollback()
	// 		return err
	// 	}
	// 	logger.INFO(fmt.Sprintf("Batch of data size %d saved succesfully.", len(stopArr)))
	// }

	// tx.Commit()

	return nil
}

func (s *syncService) syncRouteStops() error {
	// *********** Παίρνουμε το connection από το Context της εφαρμογής ************
	//var dbConnection *gorm.DB = ctx.Value(db.CONNECTIONVAR).(*gorm.DB)
	// *****************************************************************************

	// Δημιουργία ενός Rest Service για να κάνω την κλήση στον Server
	var restSrv = s.restService

	// var routeSrv = s.routeService
	response := restSrv.OasaRequestApi02("getRouteStops")
	if response.Error != nil {
		return response.Error
	}

	// TODO: Δεν θα γίνεται
	// tx := s.dbConnection.Begin()
	// err := routeSrv.WithTrx(tx).DeleteRoute02()
	// if err != nil {
	// 	tx.Rollback()
	// 	return err
	// }
	// logger.INFO("Delete from Route02 Table succesfully.")
	// var route02Arr []models.Route02 = make([]models.Route02, 0)
	s.HelpRoute02 = make([]models.Route02, 0)
	logger.INFO("Get Stops per Route data from OASA Server...")
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
		rt.Rt_code = *num32
		if _, ok := s.routeKeys[rt.Rt_code]; ok {
			num64, err := utils.StrToInt64(row[2])
			if err != nil {
				return err
			}
			rt.Stp_code = *num64
			num16, err := utils.StrToInt16(row[3])
			if err != nil {
				return err
			}
			rt.Senu = *num16
			s.HelpRoute02 = append(s.HelpRoute02, rt)
		}

		// if len(route02Arr) == 10000 {
		// 	_, err = routeSrv.WithTrx(tx).Route02InsertArr(route02Arr)
		// 	if err != nil {
		// 		tx.Rollback()
		// 		return err
		// 	}
		// 	logger.INFO(fmt.Sprintf("Batch of data size %d saved succesfully.", len(route02Arr)))
		// 	route02Arr = make([]models.Route02, 0)
		// }
	}
	// if len(route02Arr) > 0 {
	// 	_, err = routeSrv.WithTrx(tx).Route02InsertArr(route02Arr)
	// 	if err != nil {
	// 		tx.Rollback()
	// 		return err
	// 	}
	// 	logger.INFO(fmt.Sprintf("Batch of data size %d saved succesfully.", len(route02Arr)))
	// }
	// tx.Commit()
	logger.INFO("Finished sychronization data from OASA Server.")
	return nil
}

func (s *syncService) syncRouteDetails() error {
	// *********** Παίρνουμε το connection από το Context της εφαρμογής ************
	//var dbConnection *gorm.DB = ctx.Value(db.CONNECTIONVAR).(*gorm.DB)
	// *****************************************************************************

	// Δημιουργία ενός Rest Service για να κάνω την κλήση στον Server
	var restSrv = s.restService

	// var routeSrv = s.routeService
	response := restSrv.OasaRequestApi02("getRoute_detail")
	if response.Error != nil {
		return response.Error
	}

	// TODO: Δεν θα κάνουμε εδώ διαγραφή από την Βάση Δεδομένων
	// tx := s.dbConnection.Begin()
	// err := routeSrv.WithTrx(tx).DeleteRoute01()
	// if err != nil {
	// 	tx.Rollback()
	// 	return err
	// }
	// logger.INFO("Delete from Route01 Table succesfully.")

	// var route01Arr []models.Route01 = make([]models.Route01, 0)
	s.HelpRoute01 = make([]models.Route01, 0)
	logger.INFO("Get details for Route data from OASA Server...")
	for _, rec := range response.Data.([]string) {
		row := strings.Split(recPreparation(rec), ",")
		if len(row) < int(5) {
			return fmt.Errorf("Τα δεδομένα της γραμμής είναι ελλειπής.")
		}

		rt := models.Route01{}
		num32, err := utils.StrToInt32(row[1])
		if err != nil {
			return err
		}
		rt.Rt_code = *num32
		if _, ok := s.routeKeys[rt.Rt_code]; ok {
			num16, err := utils.StrToInt16(row[2])
			if err != nil {
				return err
			}
			rt.Routed_order = *num16
			fl32 := utils.StrToFloat32(row[3])
			rt.Routed_x = fl32
			fl32 = utils.StrToFloat32(row[4])
			rt.Routed_y = fl32

			s.HelpRoute01 = append(s.HelpRoute01, rt)
		}
		// if len(route01Arr) == 15000 {
		// 	_, err = routeSrv.WithTrx(tx).Route01InsertArr(route01Arr)
		// 	if err != nil {
		// 		tx.Rollback()
		// 		return err
		// 	}
		// 	logger.INFO(fmt.Sprintf("Batch of data size %d saved succesfully.", len(route01Arr)))
		// 	route01Arr = make([]models.Route01, 0)
		// }
	}
	// if len(route01Arr) > 0 {
	// 	_, err = routeSrv.WithTrx(tx).Route01InsertArr(route01Arr)
	// 	if err != nil {
	// 		tx.Rollback()
	// 		return err
	// 	}
	// 	logger.INFO(fmt.Sprintf("Batch of data size %d saved succesfully.", len(route01Arr)))
	// }
	// tx.Commit()
	logger.INFO("Finished sychronization detail for Route data from OASA Server.")
	return nil
}

func (s *syncService) SyncSchedule() error {
	lines, err := s.lineService.GetLineList()
	if err != nil {
		return err
	}
	// s.HelpSchedule = make([]models.Schedule, 0)
	// s.HelpScheduletime = make([]models.Scheduletime, 0)
	// s.HelpScheduleline = make([]models.Scheduleline, 0)

	logger.INFO("Get Schedule data per line from OASA server...")
	for _, recLine := range lines {
		// Με αυτό το Request φέρνουμε τα header των δρομολογίων
		response := s.restService.OasaRequestApi00("getScheduleDaysMasterline", map[string]interface{}{
			"p1": recLine.Line_Code,
		})
		if response.Error != nil {
			return response.Error
		}

		s.addScheduleInArray(response.Data.([]interface{}), recLine.Line_Code, int32(recLine.Ml_Code))
	}
	logger.INFO("Finished sychronization Schedule data per line from OASA server.")
	return nil
}

func (s *syncService) addScheduleInArray(currArr []interface{}, line_code int32, ml_code int32) {
	if currArr == nil {
		return
	}

	if s.HelpSchedule == nil {
		s.HelpSchedule = make([]models.Schedule, 0) // Initialize if nil
	}

	if s.HelpScheduletime == nil {
		s.HelpScheduletime = make([]models.Scheduletime, 0) // Initialize if nil
	}

	if s.HelpScheduleline == nil {
		s.HelpScheduleline = make([]models.Scheduleline, 0)
	}

	var scheduleMapper = mapper.NewScheduleMapper()
	var scheduleMappertime = mapper.NewScheduletimeMapper()
	scheduleArr, err := scheduleMapper.MapDto(currArr)
	if err != nil {
		logger.ERROR(err.Error())
	}
	for _, sched := range scheduleArr {

		if _, ok := s.scheduleMasterKey[int32(sched.Sdc_Code)]; !ok {
			s.scheduleMasterKey[int32(sched.Sdc_Code)] = int32(sched.Sdc_Code)
			s.HelpSchedule = append(s.HelpSchedule, scheduleMapper.MapDtoToSchedule(sched))
		}
		s.HelpScheduleline = append(s.HelpScheduleline, models.Scheduleline{Sdc_Cd: int64(sched.Sdc_Code), Ln_Code: line_code})
		// TODO: Εδώ πρέπει να φέρνουμε για κάθε συνδιασμό line_code, ml_code για κάθε sdc_code
		response := s.restService.OasaRequestApi00("getSchedLines", map[string]interface{}{
			"p1": ml_code,
			"p2": sched.Sdc_Code,
			"p3": line_code,
		})
		if response.Error != nil {
			return
		}
		scheduletimeDto, err := scheduleMappertime.DtoToScheduleTime(response.Data.(map[string]interface{}))
		// Εδώ φτιάχνουμε τα δεδομένα για τα δρομολογία το προορισμού
		for _, rec := range scheduletimeDto.Go {
			finalRec := scheduleMappertime.ScheduletimeDtoToScheduletime(rec, models.Direction_GO)
			finalRec.Direction = models.Direction_GO
			s.HelpScheduletime = append(s.HelpScheduletime, finalRec)
		}
		// Εδώ φτιάχνουμε τα δεδομένα για τα δρομολογία της επιστροφής
		for _, rec := range scheduletimeDto.Come {
			finalRec := scheduleMappertime.ScheduletimeDtoToScheduletime(rec, models.Direction_COME)
			finalRec.Direction = models.Direction_COME
			s.HelpScheduletime = append(s.HelpScheduletime, finalRec)
		}
		if err != nil {
			logger.ERROR(err.Error())
		}
	}

}
