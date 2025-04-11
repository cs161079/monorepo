package config

import (
	"fmt"
	"sort"
	"strings"
	"time"

	"github.com/cs161079/monorepo/common/db"
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

	syncScheduleMaster() error
	syncScheduleTime() error

	syncRouteDetails() error
	uVersionFromOasa() ([]dao.UVersion01, error)
	SyncData() error
	DeleteAll() error
	InserttoDatabase() error

	// Αυτή η μέθοδος θα ελέγχει και θα συγχρονίζει μόνο ότι πληροφορία χρειάζεται.
	SyncData02() error
}

type syncService struct {
	HelpLine         []models.Line
	HelpRoute        []models.Route
	HelpRoute01      []models.Route01
	HelpRoute02      []models.Route02
	HelpStop         []models.Stop
	HelpSchedule     []models.ScheduleMaster
	HelpScheduletime []models.ScheduleTime
	// HelpScheduleline  []models.Scheduleline
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
	routeKeys map[int32]int32
	//scheduleMasterKey map[int32]int32
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
		//scheduleMasterKey: make(map[int32]int32),
	}
}

func recPreparation(recStr string) string {
	// var trimmedSpace = strings.TrimSpace(recStr)
	return strings.ReplaceAll(recStr, "\"", "")
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
	} else if rec.UVersion.Uv_descr == "SCHED_ENTRIES" {
		rec.Orderd = 6
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

	// Θα καλέσουμε μία διαδικασία ελέγχου.
	CheckSdc(s.HelpScheduletime, s.dbConnection)

	//s.HelpScheduletime = demoArr

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

func CheckSdc(arrData []models.ScheduleTime, connection *gorm.DB) []models.ScheduleTime {

	// Kάνουμε Select όλα τα sdc_code από το πίνακα ScheduleMaster
	var distinctFromDb []int32
	dbResult := connection.Table(db.SCHEDULEMASTERTABLE).Select("sdc_code").Find(&distinctFromDb)
	if dbResult.Error != nil {
		logger.ERROR("Database error " + dbResult.Error.Error())
	}

	// Γεμίζουμε ένα Map για να μπορούμε να βρίσκουμε εύκολα με το κλειδί που θέλουμε.
	// Στην προκυμένη με το sdc_code.
	var sdcMap map[int32]int32 = make(map[int32]int32)
	for _, rec := range distinctFromDb {
		sdcMap[rec] = rec
	}

	// Κάνουμε μία ταξινόμηση με το sdc_code το πίνακα με τις ώρες
	sort.Slice(arrData, func(i, j int) bool {
		return arrData[i].SDCCd < arrData[j].SDCCd
	})

	var medData []models.ScheduleTime = make([]models.ScheduleTime, 0)

	for i, rec := range arrData {
		if rec.SDCCd != -1 {
			medData = append(medData, arrData[i])
		}
	}

	//Αρχικοποιώ τον τελικό πίνακα
	var lastData []models.ScheduleTime = make([]models.ScheduleTime, 0)

	for i, rec := range medData {
		if _, exists := sdcMap[int32(rec.SDCCd)]; exists {
			lastData = append(lastData, arrData[i])
		}
	}
	return lastData

}

func removeByIndex(slice []int, index int) []int {
	return append(slice[:index], slice[index+1:]...)
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

/*
filterRecords(records []models.UVersions, condition func(models.UVersions) bool) *models.UVersions
filterRecords takes an Array of UVersion Recors and a function for condition and returns
the records that satisfy the condition.

Parameters:

	records ([]models.UVersions): TAn Array of Data versions from OASA.
	condition func (models.UVersions) bool: The second integer to add.

Returns:

	*models.UVersions: Pointer of record that satify the condition which take as parameter.
*/
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
	logger.INFO("Fetching data from OASA Server...")
	for _, rec := range versionsArr {
		dbRec, err := uvServ.Select(rec.UVersion.Uv_descr)
		if err != nil {
			return nil
		}

		var frServeVers = rec.UVersion.Uv_lastupdatelong
		if dbRec == nil || frServeVers > dbRec.Uv_lastupdatelong {
			switch rec.UVersion.Uv_descr {
			case "LINES":
				if err := s.syncLines(); err != nil {
					return err
				}
				// Εδώ θα πρέπει να κάνουμε Update την εγγραφή στον πίνακα με το νέο Version.
				//uvServ.Post(&rec.UVersion)
			case "ROUTES":
				if err := s.syncRoutes(); err != nil {
					return err
				}
				// Εδώ θα πρέπει να κάνουμε Update την εγγραφή στον πίνακα με το νέο Version.
				//uvServ.Post(&rec.UVersion)
			case "STOPS":
				if err := s.syncStops(); err != nil {
					return err
				}
				// Εδώ θα πρέπει να κάνουμε Update την εγγραφή στον πίνακα με το νέο Version.
				//uvServ.Post(&rec.UVersion)
			case "ROUTE STOPS":
				if err := s.syncRouteStops(); err != nil {
					return err
				}
				// Εδώ θα πρέπει να κάνουμε Update την εγγραφή στον πίνακα με το νέο Version.
				//uvServ.Post(&rec.UVersion)
			case "ROUTE DETAIL":
				if err := s.syncRouteDetails(); err != nil {
					return err
				}
			case "SCHED_CATS":
				if err := s.syncScheduleMaster(); err != nil {
					return err
				}
			case "SCHED_ENTRIES":
				if err := s.syncScheduleTime(); err != nil {
					return err
				}
			}

		}
	}
	return nil
}

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
	logger.INFO("\tFetch lines data...")
	s.HelpLine = make([]models.Line, 0)
	for _, ln := range response.Data.([]any) {
		lineOasa := lineSrv.GetMapper().GenDtLineOasa(ln.(map[string]interface{}))
		line := lineSrv.GetMapper().OasaToLine(lineOasa)

		// if _, ok := s.lineKeys[line.Line_Code]; !ok {
		// 	s.lineKeys[line.Line_Code] = line.Line_Code
		// }
		s.HelpLine = append(s.HelpLine, line)
	}
	return nil
}

func (s *syncService) syncRoutes() error {
	// *********** Κάνου get το connection της  βάσης από το Context ************
	//var dbConnection *gorm.DB = ctx.Value(db.CONNECTIONVAR).(*gorm.DB)
	// **************************************************************************
	var restSrv = s.restService //service.NewRestService()

	// var routeSrv = s.routeService
	logger.INFO("\tFetching Routes data...")
	response := restSrv.OasaRequestApi02("getRoutes")
	if response.Error != nil {
		return response.Error
	}

	// Create or overwrite the file (if it exists)
	// file, err := os.Create("tmp/routes.txt")
	// if err != nil {
	// 	//fmt.Println("Error creating file:", err)
	// 	return err
	// }
	// defer file.Close() // Make sure to close the file when done

	// Δεν θα χρησιμοποιήσουμε τοπικό πίνακα
	// var routeArray []models.Route = make([]models.Route, 0)
	// Εδώ η διαδικασία μας γυρνάει από το API έναν πίνακα με τα Record σε γραμμή χωρισμένα τα πεδία με κόμμα
	for _, rec := range response.Data.([]string) {

		// Γράφουμε κάθε γραμμή των δεδομένων στο αρχείο.
		// fmt.Fprintf(file, "%s\n", rec)

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
		rt.LnCode = *num
		num, err = utils.StrToInt32(recordArr[0])
		if err != nil {
			return err
		}
		rt.RouteCode = *num
		if _, ok := s.routeKeys[rt.RouteCode]; !ok {
			s.routeKeys[rt.RouteCode] = rt.RouteCode
		}

		rt.RouteDescr = recordArr[2]
		rt.RouteDescrEng = recordArr[3]
		num, err = utils.StrToInt32(recordArr[4])
		if err != nil {
			return err
		}
		rt.RouteType = int8(*num)
		fl32 := utils.StrToFloat32(recordArr[5])
		rt.RouteDistance = fl32

		s.HelpRoute = append(s.HelpRoute, rt)
	}

	return nil
}

func (s *syncService) syncStops() error {
	// *********** Κάνου get το connection της  βάσης από το Context ************
	//var dbConnection *gorm.DB = ctx.Value(db.CONNECTIONVAR).(*gorm.DB)
	// **************************************************************************
	var restSrv = s.restService //service.NewRestService()

	// var stopSrv = s.stopService
	logger.INFO("\tFetching stops data...")
	response := restSrv.OasaRequestApi02("getStops")
	if response.Error != nil {
		return response.Error
	}

	// Create or overwrite the file (if it exists)
	// file, err := os.Create("tmp/stops.txt")
	// if err != nil {
	// 	//fmt.Println("Error creating file:", err)
	// 	return err
	// }
	// defer file.Close() // Make sure to close the file when done

	s.HelpStop = make([]models.Stop, 0)
	// Εδώ η διαδικασία μας γυρνάει από το API έναν πίνακα με τα Record σε γραμμή χωρισμένα τα πεδία με κόμμα
	logger.INFO("Get Stop data from OASA Server...")
	for _, rec := range response.Data.([]string) {

		// Γράφουμε κάθε γραμμή των δεδομένων στο αρχείο.
		// fmt.Fprintf(file, "%s\n", rec)

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
		st.StopCode = *num32
		st.StopID = recordArr[1]
		st.StopDescr = recordArr[2]
		st.StopDescr = recordArr[3]
		st.StopStreet = recordArr[4]
		st.StopStreetEng = recordArr[5]
		num32, err = utils.StrToInt32(recordArr[6])
		if err != nil {
			return err
		}
		st.StopHeading = *num32
		st.StopLng = utils.StrToFloat(recordArr[7])
		st.StopLat = utils.StrToFloat(recordArr[8])
		num8, err := utils.StrToInt8(recordArr[9])
		if err != nil {
			return err
		}
		st.StopType = *num8
		num8, err = utils.StrToInt8(recordArr[10])
		if err != nil {
			return err
		}
		st.StopAmea = *num8
		st.Destinations = recordArr[11]
		st.DestinationsEng = recordArr[12]

		s.HelpStop = append(s.HelpStop, st)

		// TODO: Αλλαγή δεν θα γίνετα

	}

	return nil
}

func (s *syncService) syncRouteStops() error {
	// *********** Παίρνουμε το connection από το Context της εφαρμογής ************
	//var dbConnection *gorm.DB = ctx.Value(db.CONNECTIONVAR).(*gorm.DB)
	// *****************************************************************************

	// Δημιουργία ενός Rest Service για να κάνω την κλήση στον Server
	var restSrv = s.restService

	// var routeSrv = s.routeService
	logger.INFO("\tFetching stops data per route ...")
	response := restSrv.OasaRequestApi02("getRouteStops")
	if response.Error != nil {
		return response.Error
	}

	// Create or overwrite the file (if it exists)
	// file, err := os.Create("tmp/route_stops.txt")
	// if err != nil {
	// 	//fmt.Println("Error creating file:", err)
	// 	return err
	// }
	// defer file.Close() // Make sure to close the file when done

	s.HelpRoute02 = make([]models.Route02, 0)
	logger.INFO("Get Stops per Route data from OASA Server...")
	for _, rec := range response.Data.([]string) {

		// Γράφουμε κάθε γραμμή των δεδομένων στο αρχείο.
		// fmt.Fprintf(file, "%s\n", rec)

		row := strings.Split(recPreparation(rec), ",")
		if len(row) < int(4) {
			return fmt.Errorf("Τα δεδομένα της γραμμής είναι ελλειπής.")
		}

		rt := models.Route02{}
		num32, err := utils.StrToInt32(row[1])
		if err != nil {
			return err
		}
		rt.RtCode = *num32
		if _, ok := s.routeKeys[rt.RtCode]; ok {
			num32, err := utils.StrToInt32(row[2])
			if err != nil {
				return err
			}
			rt.StpCode = *num32
			num16, err := utils.StrToInt16(row[3])
			if err != nil {
				return err
			}
			rt.Senu = *num16
			s.HelpRoute02 = append(s.HelpRoute02, rt)
		}
	}
	return nil
}

func (s *syncService) syncRouteDetails() error {
	// *********** Παίρνουμε το connection από το Context της εφαρμογής ************
	//var dbConnection *gorm.DB = ctx.Value(db.CONNECTIONVAR).(*gorm.DB)
	// *****************************************************************************

	// Δημιουργία ενός Rest Service για να κάνω την κλήση στον Server
	var restSrv = s.restService

	// var routeSrv = s.routeService
	logger.INFO("\tFetching Route details & information data...")
	response := restSrv.OasaRequestApi02("getRoute_detail")
	if response.Error != nil {
		return response.Error
	}

	// Create or overwrite the file (if it exists)
	// file, err := os.Create("tmp/route_details.txt")
	// if err != nil {
	// 	//fmt.Println("Error creating file:", err)
	// 	return err
	// }
	// defer file.Close() // Make sure to close the file when done

	s.HelpRoute01 = make([]models.Route01, 0)
	for _, rec := range response.Data.([]string) {

		// Γράφουμε κάθε γραμμή των δεδομένων στο αρχείο.
		// fmt.Fprintf(file, "%s\n", rec)

		row := strings.Split(recPreparation(rec), ",")
		if len(row) < int(5) {
			return fmt.Errorf("Τα δεδομένα της γραμμής είναι ελλειπής.")
		}

		rt := models.Route01{}
		num32, err := utils.StrToInt32(row[1])
		if err != nil {
			return err
		}
		rt.RtCode = *num32
		if _, ok := s.routeKeys[rt.RtCode]; ok {
			num16, err := utils.StrToInt16(row[2])
			if err != nil {
				return err
			}
			rt.RoutedOrder = *num16
			fl32 := utils.StrToFloat(row[3])
			rt.RoutedX = fl32
			fl32 = utils.StrToFloat(row[4])
			rt.RoutedY = fl32

			s.HelpRoute01 = append(s.HelpRoute01, rt)
		}

	}
	return nil
}

/*
Αυτό είναι από τα APIs που έχω αναλαλύψει μόνος μου, δεν τα δίνει ο OASA.
Τα χρησιμοποιεί η εφαρμογή για να συγχρονίζει δεδομένα.
*/
func (s *syncService) syncScheduleMaster() error {
	// *********** Παίρνουμε το connection από το Context της εφαρμογής ************
	//var dbConnection *gorm.DB = ctx.Value(db.CONNECTIONVAR).(*gorm.DB)
	// *****************************************************************************

	// Δημιουργία ενός Rest Service για να κάνω την κλήση στον Server
	var restSrv = s.restService

	// var routeSrv = s.routeService
	logger.INFO("\tFetchig Route Scedule infromatino data...")
	response := restSrv.OasaRequestApi02("getSched_Cats")
	if response.Error != nil {
		return response.Error
	}

	// Create or overwrite the file (if it exists)
	// file, err := os.Create("tmp/schedule_master.txt")
	// if err != nil {
	// 	//fmt.Println("Error creating file:", err)
	// 	return err
	// }
	// defer file.Close() // Make sure to close the file when done

	s.HelpSchedule = make([]models.ScheduleMaster, 0)
	for _, rec := range response.Data.([]string) {

		// Γράφουμε κάθε γραμμή των δεδομένων στο αρχείο.
		// fmt.Fprintf(file, "%s\n", rec)

		row := strings.Split(recPreparation(rec), ",")
		if len(row) < int(5) {
			return fmt.Errorf("Data is missing for Master Schedule.")
		}

		rt := models.ScheduleMaster{}
		for index, recField := range row {
			recField = strings.TrimSpace(recField)
			if index == 0 {
				num32, err := utils.StrToInt32(recField)
				if err != nil {
					return err
				}
				rt.SDCCode = *num32
			} else if index == 1 {
				rt.SDCDescr = recField
			} else if index == 2 {
				rt.SDCDescrEng = recField
			} else if index == 3 {
				rt.SDCDays = recField
			} else if index == 4 {
				rt.SDCMonths = recField
			}
		}
		s.HelpSchedule = append(s.HelpSchedule, rt)
	}
	return nil
}

func convertStrToTime(strVal string) *time.Time {
	timeVal, err := time.Parse(models.CustomTimeFormat, strings.TrimSpace(strVal))
	if err != nil {
		logger.Logger.Error("Σφάλμα κατά την μετατροπή απόσυμβολοσειρά σε Time. Δεν είναι valid τιμή %s. [%v]", strVal, err)
		return nil
	}
	return &timeVal
}

func (s *syncService) syncScheduleTime() error {
	// *********** Παίρνουμε το connection από το Context της εφαρμογής ************
	//var dbConnection *gorm.DB = ctx.Value(db.CONNECTIONVAR).(*gorm.DB)
	// *****************************************************************************

	// Δημιουργία ενός Rest Service για να κάνω την κλήση στον Server
	var restSrv = s.restService

	// var routeSrv = s.routeService
	logger.INFO("\tFetching Schedule details and times ...")
	response := restSrv.OasaRequestApi02("getSched_entries")
	if response.Error != nil {
		return response.Error
	}

	// Create or overwrite the file (if it exists)
	// file, err := os.Create("tmp/schedule_times.txt")
	// if err != nil {
	// 	//fmt.Println("Error creating file:", err)
	// 	return err
	// }
	// defer file.Close() // Make sure to close the file when done

	s.HelpScheduletime = make([]models.ScheduleTime, 0)
	logger.INFO("Get Schedule time entries data from OASA Server...")
	for _, rec := range response.Data.([]string) {
		// Γράφουμε κάθε γραμμή των δεδομένων στο αρχείο.
		// fmt.Fprintf(file, "%s\n", rec)

		row := strings.Split(recPreparation(rec), ",")
		if len(row) < int(13) {
			return fmt.Errorf("Data row Schedule time data is missing or corrupted.")
		}

		/*
				Από την γραμμή χρειαζόμαστε
			    index arr 1 -> sdc_code,
						  4 -> line_code,
						  6 -> start1,
						  7 -> end1,
						  10 -> start2,
						  11 -> end2,
						  12 -> sort
		*/
		var mapVals map[string]interface{} = make(map[string]interface{})
		for index, recField := range row {
			recField = strings.TrimSpace(recField)
			if index == 1 {
				num32, err := utils.StrToInt32(recField)
				if err != nil {
					return fmt.Errorf("Error occured on sdc_code=%s field convertion from string to number. %v", row[1], err)
				}
				mapVals["sdc_code"] = *num32
			} else if index == 4 {
				num32, err := utils.StrToInt32(recField)
				if err != nil {
					return fmt.Errorf("Error occured on line_code=%s field convertion from string to number. %v", row[4], err)
				}
				mapVals["line_code"] = *num32
			} else if index == 6 {
				mapVals["Start_time1"] = recField
			} else if index == 7 {
				mapVals["End_time1"] = recField
			} else if index == 10 {
				mapVals["Start_time2"] = recField
			} else if index == 11 {
				mapVals["End_time2"] = recField
			} else if index == 12 {
				inSort, err := utils.StrToInt32(recField)
				if err != nil {
					return fmt.Errorf("Error occured on sort=%s field convertion from string to number. %v", row[12], err)
				}
				mapVals["sort"] = *inSort
			}
		}

		if mapVals["sort"].(int32) > 0 {
			if mapVals["Start_time1"] != "null" && mapVals["End_time1"] != "null" {
				var rtt = models.ScheduleTime{}
				mapper.MapStruct(mapVals, &rtt)
				rtt.Direction = models.Direction_GO
				timeVal := convertStrToTime(mapVals["Start_time1"].(string))
				if timeVal != nil {
					rtt.StartTime = models.OpswTime(*timeVal)
				}

				timeVal = convertStrToTime(mapVals["End_time1"].(string))
				if timeVal != nil {
					rtt.EndTime = models.OpswTime(*timeVal)
				}
				s.HelpScheduletime = append(s.HelpScheduletime, rtt)
			}

			if mapVals["Start_time2"] != "null" && mapVals["End_time2"] != "null" {
				var rtt = models.ScheduleTime{}
				mapper.MapStruct(mapVals, &rtt)
				rtt.Direction = models.Direction_COME
				timeVal := convertStrToTime(mapVals["Start_time2"].(string))
				if timeVal != nil {
					rtt.StartTime = models.OpswTime(*timeVal)
				}

				timeVal = convertStrToTime(mapVals["End_time2"].(string))
				if timeVal != nil {
					rtt.EndTime = models.OpswTime(*timeVal)
				}
				s.HelpScheduletime = append(s.HelpScheduletime, rtt)
			}
		}
	}
	return nil
}

func (s *syncService) SyncData02() error {
	// Πρέπει να φέρουμε πληροφορία από τον OASA Server το max id για κάθέ οντότητα π.χ Line, Route, Route01, Route02,
	// Schedule, Scheduletime

	restResponse := s.restService.OasaRequestApi00("getUversions", nil)
	if restResponse.Error != nil {
		return restResponse.Error
	}
	// Κάνω Cast σε Array από Maps
	restData := restResponse.Data.([]any)
	mapper := mapper.NewUVersionMapper()
	var arrUvers []models.UVersionsOasa = make([]models.UVersionsOasa, 0)
	for _, rec := range restData {
		arrUvers = append(arrUvers, mapper.GeneralUVersions(rec))
	}
	fmt.Printf("These are data %v", arrUvers)

	return nil
}
