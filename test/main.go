package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/cs161079/monorepo/common/mapper"
	"github.com/cs161079/monorepo/common/models"
	"github.com/cs161079/monorepo/common/service"
	"github.com/cs161079/monorepo/common/utils"
	logger "github.com/cs161079/monorepo/common/utils/goLogger"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var opswLogger logger.OpswLogger

func main() {
	start := time.Now()
	// Depedencies
	restSrv := service.NewRestService()
	opswLogger = logger.CreateLogger()
	lineSrv := service.NewLineService(nil)

	response := restSrv.OasaRequestApi00("webGetLinesWithMLInfo", nil)
	if response.Error != nil {
		opswLogger.ERROR(response.Error.Error())
		return
	}
	// TODO: Το έκοψα γιατί δεν θα κάνω εδώ το Delete
	opswLogger.INFO("Fetch Line data completed successfully.")
	var lines []models.LineM
	var _ []interface{} = make([]interface{}, 0)
	var lineOasaArr []models.LineOasa = make([]models.LineOasa, 0)
	for _, ln := range response.Data.([]any) {
		lineOasa := lineSrv.GetMapper().GenDtLineOasa(ln.(map[string]interface{}))
		lineOasaArr = append(lineOasaArr, lineOasa)
	}

	bts, err := json.Marshal(lineOasaArr)
	if err != nil {
		opswLogger.ERROR(err.Error())
		return
	}
	json.Unmarshal(bts, &lines)

	response = restSrv.OasaRequestApi02("getRoutes")
	if response.Error != nil {
		opswLogger.ERROR(response.Error.Error())
		return
	}

	opswLogger.INFO("Going to insert to Database.")

	response = restSrv.OasaRequestApi02("getRoutes")
	if response.Error != nil {
		opswLogger.ERROR(response.Error.Error())
		return
	}

	opswLogger.INFO("Fetch Route data completed successfully.")

	routes, err := analyzeRouteData(response.Data.([]string))
	if err != nil {
		opswLogger.ERROR(err.Error())
		return
	}

	sort.Slice(routes, func(i, j int) bool {
		return routes[i].LnCode < routes[j].LnCode
	})

	var rtMap map[int32][]models.RouteM = make(map[int32][]models.RouteM)
	var rtArr []models.RouteM = make([]models.RouteM, 0)
	var lnCode = routes[0].LnCode
	for idx, rt01 := range routes {
		if lnCode != rt01.LnCode {
			rtMap[lnCode] = rtArr
			lnCode = rt01.LnCode
			rtArr = make([]models.RouteM, 0)
		}
		bts, err := json.Marshal(rt01)
		if err != nil {
			opswLogger.ERROR(err.Error())
			return
		}
		var rtM models.RouteM
		json.Unmarshal(bts, &rtM)
		rtArr = append(rtArr, rtM)
		if idx == (len(routes) - 1) {
			rtMap[lnCode] = rtArr
		}
	}

	opswLogger.INFO("\tFetching Route details & information data...")
	response = restSrv.OasaRequestApi02("getRoute_detail")
	if response.Error != nil {
		opswLogger.ERROR(response.Error.Error())
		return
	}

	routeDetailsArr, err := analyzeRouteDetailsData(response.Data.([]string))
	if err != nil {
		opswLogger.ERROR(err.Error())
		return
	}

	sort.Slice(routeDetailsArr, func(i, j int) bool {
		return routeDetailsArr[i].RtCode < routeDetailsArr[j].RtCode
	})

	var curRouteCode = routeDetailsArr[0].RtCode
	var rtDetailsMap map[int32][]models.Route01M = make(map[int32][]models.Route01M)
	var rtDetail []models.Route01 = make([]models.Route01, 0)
	for indx, rec := range routeDetailsArr {
		if curRouteCode != rec.RtCode {
			bts, err := json.Marshal(rtDetail)
			if err != nil {
				opswLogger.ERROR(err.Error())
				return
			}
			var rtM []models.Route01M
			json.Unmarshal(bts, &rtM)
			rtDetailsMap[curRouteCode] = rtM
			rtDetail = make([]models.Route01, 0)
			curRouteCode = rec.RtCode
		}
		rtDetail = append(rtDetail, rec)
		if indx == len(routeDetailsArr)-1 {
			bts, err := json.Marshal(rtDetail)
			if err != nil {
				opswLogger.ERROR(err.Error())
				return
			}
			var rtM []models.Route01M
			json.Unmarshal(bts, &rtM)
			rtDetailsMap[curRouteCode] = rtM
		}
	}

	// ========= Συγχρονισμός των δεδομένων των στάσεων από το ΟΑΣΑ server ==============
	//               Δημιουργώ ένα Map για γρήγορη επιλογή των στάσεων.
	opswLogger.INFO("Fetching stops data per route ...")
	response = restSrv.OasaRequestApi02("getStops")
	if response.Error != nil {
		opswLogger.ERROR(response.Error.Error())
		return
	}

	stopsArr, err := analyzeStopData(response.Data.([]string))
	if err != nil {
		opswLogger.ERROR(err.Error())
		return
	}
	var stopMap map[int32]models.StopM = make(map[int32]models.StopM)
	for _, stop := range stopsArr {
		bts, err := json.Marshal(stop)
		if err != nil {
			opswLogger.ERROR(err.Error())
			return
		}
		var stopM models.StopM
		json.Unmarshal(bts, &stopM)
		stopMap[stopM.StopCode] = stopM
	}

	opswLogger.INFO("Fetching stops data per route ...")
	response = restSrv.OasaRequestApi02("getRouteStops")
	if response.Error != nil {
		opswLogger.ERROR(response.Error.Error())
		return
	}

	routeStopArr, err := analyzeRouteStopData(response.Data.([]string))
	if err != nil {
		opswLogger.ERROR(err.Error())
		return
	}

	sort.Slice(routeStopArr, func(i, j int) bool {
		if routeStopArr[i].Route.RouteCode == routeStopArr[j].RtCode {
			return routeStopArr[i].Senu < routeStopArr[j].Senu // compare Field2 if Field1 is equal
		}
		return routeStopArr[i].RtCode < routeStopArr[j].RtCode // primary sort by Field1
	})

	var routeStopMap map[int32][]models.StopM = make(map[int32][]models.StopM)
	var routeCurr = routeStopArr[0].RtCode
	var stopArr []models.StopM = make([]models.StopM, 0)
	for idx, rec := range routeStopArr {
		if routeCurr != rec.RtCode {
			routeStopMap[routeCurr] = stopArr
			routeCurr = rec.RtCode
			stopArr = make([]models.StopM, 0)
		}

		if _, ok := stopMap[int32(rec.StpCode)]; ok {
			bts, err := json.Marshal(stopMap[int32(rec.StpCode)])
			if err != nil {
				opswLogger.ERROR(err.Error())
				return
			}
			var stopM models.StopM
			json.Unmarshal(bts, &stopM)
			stopM.StopSenu = rec.Senu
			stopArr = append(stopArr, stopM)
			if idx == (len(routeStopArr) - 1) {
				routeStopMap[routeCurr] = stopArr
			}
		}

	}

	// ============================= Συγχρονισμός Προγραμμάτων και δρομολογίων =========================================
	opswLogger.INFO("Fetchig Route Scedule infromation data...")
	response = restSrv.OasaRequestApi02("getSched_Cats")
	if response.Error != nil {
		opswLogger.ERROR(response.Error.Error())
		return
	}

	scheduleMasterMap, err := analyzeScheduleMaster(response.Data.([]string))
	if err != nil {
		opswLogger.ERROR(err.Error())
		return
	}

	// ============================= Συγχρονισμός Προγραμμάτων και δρομολογίων =========================================
	opswLogger.INFO("Fetching Schedule details and times ...")
	response = restSrv.OasaRequestApi02("getSched_entries")
	if response.Error != nil {
		opswLogger.ERROR(response.Error.Error())
		return
	}

	scheduleEntries, err := analyzeScheduleEntries(response.Data.([]string))
	if err != nil {
		opswLogger.ERROR(err.Error())
		return
	}

	var currLnCode = scheduleEntries[0].LnCode
	var currSdcCode = scheduleEntries[0].SDCCd
	var finalScheduleMap map[int32][]models.ScheduleMasterM = make(map[int32][]models.ScheduleMasterM)
	var scheduleArr []models.ScheduleMasterM = make([]models.ScheduleMasterM, 0)
	var schedule models.ScheduleMasterM
	bts, err = json.Marshal(scheduleMasterMap[int32(currSdcCode)])
	if err != nil {
		opswLogger.ERROR(err.Error())
		return
	}
	schedule = models.ScheduleMasterM{}
	json.Unmarshal(bts, &schedule)
	schedule.Times = make([]models.ScheduleTimeM, 0)
	for _, entry := range scheduleEntries {
		if currLnCode != entry.LnCode {
			finalScheduleMap[int32(currLnCode)] = scheduleArr
			currLnCode = entry.LnCode
			currSdcCode = entry.SDCCd
			scheduleArr = make([]models.ScheduleMasterM, 0)
			bts, err := json.Marshal(scheduleMasterMap[int32(currSdcCode)])
			if err != nil {
				opswLogger.ERROR(err.Error())
				return
			}
			schedule = models.ScheduleMasterM{}
			json.Unmarshal(bts, &schedule)
			schedule.Times = make([]models.ScheduleTimeM, 0)
		} else if currSdcCode != entry.SDCCd {
			scheduleArr = append(scheduleArr, schedule)
			currSdcCode = entry.SDCCd

			bts, err := json.Marshal(scheduleMasterMap[int32(currSdcCode)])
			if err != nil {
				opswLogger.ERROR(err.Error())
				return
			}
			schedule = models.ScheduleMasterM{}
			json.Unmarshal(bts, &schedule)
			schedule.Times = make([]models.ScheduleTimeM, 0)
		}

		btsEntry, err := json.Marshal(entry)
		if err != nil {
			opswLogger.ERROR(err.Error())
			return
		}
		var scheduleTM models.ScheduleTimeM = models.ScheduleTimeM{}
		err = json.Unmarshal(btsEntry, &scheduleTM)
		if err != nil {
			opswLogger.ERROR(err.Error())
		}
		schedule.Times = append(schedule.Times, scheduleTM)
	}

	opswLogger.INFO("Schedule sychronization finished successfully.")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// ====================== Connect to Mongo Database ======================
	clientOpts := options.Client().ApplyURI("mongodb://localhost:27017").
		SetServerSelectionTimeout(5 * time.Second)

	client, err := mongo.Connect(ctx, clientOpts)
	if err != nil {
		log.Fatal(err)
	}
	defer client.Disconnect(ctx)
	// ======================================================================

	// ================ Create or Use exist Database and Create Collection ==============
	database := client.Database("testdb")
	collection := database.Collection("line")
	if err := collection.Drop(ctx); err != nil {
		opswLogger.ERROR(err.Error())
		return
	}
	collection = database.Collection("line")

	for _, line := range lines {
		line.Routes = rtMap[int32(line.LineCode)]
		line.Schedules = finalScheduleMap[int32(line.LineCode)]
		for indx01, route := range line.Routes {
			line.Routes[indx01].Stops = append(line.Routes[indx01].Stops, routeStopMap[route.RouteCode]...)
			line.Routes[indx01].Details = append(line.Routes[indx01].Details, rtDetailsMap[route.RouteCode]...)
		}

		// Εισαγωγή στην MongoDb
		_, err = collection.InsertOne(ctx, line)
		if err != nil {
			opswLogger.ERROR(err.Error())
		}
		opswLogger.INFO("Insert to Database completed successfully. [line_code: " + strconv.Itoa(line.LineCode) + "].")
	}

	opswLogger.INFO("Synchronized Lines and Routes completed successfully.")

	// var lnCodCurr int32 = routesRec[0].LnCode

	// var lineRoutes []interface{} = make([]interface{}, 0)
	// for i, rec := range routesRec {
	// 	if rec.LnCode != int32(lnCodCurr) {
	// 		err = insertToMongo(ctx, collection, lineRoutes, int32(lnCodCurr))
	// 		if err != nil {
	// 			opswLogger.ERROR(err.Error())
	// 			return
	// 		}
	// 		lineRoutes = make([]interface{}, 0)
	// 		lnCodCurr = rec.LnCode
	// 	}

	// 	lineRoutes = append(lineRoutes, routesRec[i])
	// }
	fmt.Printf("=========================================================================================================\n")
	fmt.Printf("============================ Sychronize make %.2f seconds to complete. ==================================\n", time.Since(start).Seconds())
	fmt.Printf("=========================================================================================================\n")

}

func convertStrToTime(strVal string) *time.Time {
	timeVal, err := time.Parse(models.CustomTimeFormat, strings.TrimSpace(strVal))
	if err != nil {
		opswLogger.ERROR(fmt.Sprintf("Σφάλμα κατά την μετατροπή απόσυμβολοσειρά σε Time. Δεν είναι valid τιμή %s. [%v]", strVal, err))
		return nil
	}
	return &timeVal
}

func analyzeScheduleEntries(data []string) ([]models.ScheduleTime, error) {
	var scheduleEntries []models.ScheduleTime = make([]models.ScheduleTime, 0)
	logger.INFO("Get Schedule time entries data from OASA Server...")
	for _, rec := range data {
		// Γράφουμε κάθε γραμμή των δεδομένων στο αρχείο.
		// fmt.Fprintf(file, "%s\n", rec)

		row := strings.Split(recPreparation(rec), ",")
		if len(row) < int(13) {
			return nil, fmt.Errorf("Data row Schedule time data is missing or corrupted.")
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
					return nil, fmt.Errorf("Error occured on sdc_code=%s field convertion from string to number. %v", row[1], err)
				}
				mapVals["sdc_code"] = *num32
			} else if index == 4 {
				num32, err := utils.StrToInt32(recField)
				if err != nil {
					return nil, fmt.Errorf("Error occured on line_code=%s field convertion from string to number. %v", row[4], err)
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
					return nil, fmt.Errorf("Error occured on sort=%s field convertion from string to number. %v", row[12], err)
				}
				mapVals["sort"] = *inSort
			}
		}
		if mapVals["sdc_code"] == -1 {
			opswLogger.INFO("Sdc code is -1")
		}

		if mapVals["sort"].(int32) > 0 && mapVals["line_code"] != int32(-1) && mapVals["sdc_code"] != int32(-1) {
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
				scheduleEntries = append(scheduleEntries, rtt)
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
				scheduleEntries = append(scheduleEntries, rtt)
			}
		}
	}
	sort.Slice(scheduleEntries, func(i, j int) bool {
		if scheduleEntries[i].LnCode == scheduleEntries[j].LnCode {
			return scheduleEntries[i].SDCCd < scheduleEntries[j].SDCCd
		}
		return scheduleEntries[i].LnCode < scheduleEntries[j].LnCode
	})
	return scheduleEntries, nil
}

func analyzeScheduleMaster(data []string) (map[int32]models.ScheduleMaster, error) {
	var scheduleMaster map[int32]models.ScheduleMaster = make(map[int32]models.ScheduleMaster)
	for _, rec := range data {

		// Γράφουμε κάθε γραμμή των δεδομένων στο αρχείο.
		// fmt.Fprintf(file, "%s\n", rec)

		row := strings.Split(recPreparation(rec), ",")
		if len(row) < int(5) {
			return nil, fmt.Errorf("Data is missing for Master Schedule.")
		}

		rt := models.ScheduleMaster{}
		for index, recField := range row {
			recField = strings.TrimSpace(recField)
			if index == 0 {
				num32, err := utils.StrToInt32(recField)
				if err != nil {
					return nil, err
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
		scheduleMaster[rt.SDCCode] = rt
	}
	return scheduleMaster, nil
}

func analyzeRouteDetailsData(data []string) ([]models.Route01, error) {
	var route01Arr []models.Route01 = make([]models.Route01, 0)
	for _, rec := range data {
		// Γράφουμε κάθε γραμμή των δεδομένων στο αρχείο.
		// fmt.Fprintf(file, "%s\n", rec)

		row := strings.Split(recPreparation(rec), ",")
		if len(row) < int(5) {
			return nil, fmt.Errorf("Τα δεδομένα της γραμμής είναι ελλειπής.")
		}

		rt := models.Route01{}
		num32, err := utils.StrToInt32(strings.TrimSpace(row[1]))
		if err != nil {
			return nil, err
		}
		rt.RtCode = *num32
		num16, err := utils.StrToInt16(strings.TrimSpace(row[2]))
		if err != nil {
			return nil, err
		}
		rt.RoutedOrder = *num16
		fl32 := utils.StrToFloat(strings.TrimSpace((row[3])))
		rt.RoutedX = fl32
		fl32 = utils.StrToFloat(strings.TrimSpace(row[4]))
		rt.RoutedY = fl32

		route01Arr = append(route01Arr, rt)

	}
	return route01Arr, nil
}

func analyzeRouteStopData(data []string) ([]models.Route02, error) {
	var routeStopsArr []models.Route02 = make([]models.Route02, 0)
	opswLogger.INFO("Get Stops per Route data from OASA Server...")
	for _, rec := range data {

		// Γράφουμε κάθε γραμμή των δεδομένων στο αρχείο.
		// fmt.Fprintf(file, "%s\n", rec)

		row := strings.Split(recPreparation(rec), ",")
		if len(row) < int(4) {
			return nil, fmt.Errorf("Τα δεδομένα της γραμμής είναι ελλειπής.")
		}

		rt := models.Route02{}
		num32, err := utils.StrToInt32(strings.TrimSpace(row[1]))
		if err != nil {
			return nil, err
		}
		rt.RtCode = *num32

		num32, err = utils.StrToInt32(strings.TrimSpace(row[2]))
		if err != nil {
			return nil, err
		}
		rt.StpCode = *num32
		num16, err := utils.StrToInt16(strings.TrimSpace(row[3]))
		if err != nil {
			return nil, err
		}
		rt.Senu = *num16
		routeStopsArr = append(routeStopsArr, rt)
	}
	return routeStopsArr, nil
}

func analyzeStopData(data []string) ([]models.Stop, error) {
	var stopArr []models.Stop = make([]models.Stop, 0)
	for _, rec := range data {

		// Γράφουμε κάθε γραμμή των δεδομένων στο αρχείο.
		// fmt.Fprintf(file, "%s\n", rec)

		// ************** Κάθε γραμμή την κάνω Split με το κόμμα και γεμίζω τα Record των στάσεων **************
		// ************************* Έλεγχος της γραμμής εάν έχει όλη την πληροφορία *****************************
		recordArr := strings.Split(recPreparation(rec), ",")
		if len(recordArr) < 13 {
			return nil, fmt.Errorf("Τα δεδομένα της γραμμής είναι ελλειπής.")
		}
		st := models.Stop{}
		num32, err := utils.StrToInt32(strings.TrimSpace(recordArr[0]))
		if err != nil {
			return nil, err
		}
		st.StopCode = *num32
		st.StopID = strings.TrimSpace(recordArr[1])
		st.StopDescr = strings.TrimSpace(recordArr[2])
		st.StopDescrEng = strings.TrimSpace(recordArr[3])
		st.StopStreet = strings.TrimSpace(recordArr[4])
		st.StopStreetEng = strings.TrimSpace(recordArr[5])
		num32, err = utils.StrToInt32(strings.TrimSpace(recordArr[6]))
		if err != nil {
			return nil, err
		}
		st.StopHeading = *num32
		st.StopLng = utils.StrToFloat(strings.TrimSpace(recordArr[7]))
		st.StopLat = utils.StrToFloat(strings.TrimSpace(recordArr[8]))
		num8, err := utils.StrToInt8(strings.TrimSpace(recordArr[9]))
		if err != nil {
			return nil, err
		}
		st.StopType = *num8
		num8, err = utils.StrToInt8(strings.TrimSpace(recordArr[10]))
		if err != nil {
			return nil, err
		}
		st.StopAmea = *num8
		st.Destinations = strings.TrimSpace(recordArr[11])
		st.DestinationsEng = strings.TrimSpace(recordArr[12])

		stopArr = append(stopArr, st)
	}
	return stopArr, nil
}

func analyzeRouteData(data []string) ([]models.Route, error) {
	var routesRec []models.Route = make([]models.Route, 0)
	// Δεν θα χρησιμοποιήσουμε τοπικό πίνακα
	// var routeArray []models.Route = make([]models.Route, 0)
	// Εδώ η διαδικασία μας γυρνάει από το API έναν πίνακα με τα Record σε γραμμή χωρισμένα τα πεδία με κόμμα
	for _, rec := range data {

		// Γράφουμε κάθε γραμμή των δεδομένων στο αρχείο.
		// fmt.Fprintf(file, "%s\n", rec)

		// ************** Κάθε γραμμή την κάνω Split με το κόμμα και γεμίζω τα Record των διαδρομών **************
		// ************************* Έλεγχος της γραμμής εάν έχει όλη την πληροφορία *****************************
		recordArr := strings.Split(recPreparation(rec), ",")
		if len(recordArr) < 6 {
			// opswLogger.ERROR("Η γραμμή του Record  είναι ελλειπής.")
			return nil, fmt.Errorf("Η γραμμή του Record  είναι ελλειπής.")
		}

		rt := models.Route{}

		lineCode, err := utils.StrToInt32(strings.TrimSpace(recordArr[1]))
		if err != nil {
			// opswLogger.ERROR(err.Error())
			return nil, err
		}
		rt.LnCode = *lineCode

		num, err := utils.StrToInt32(strings.TrimSpace(recordArr[0]))
		if err != nil {
			// opswLogger.ERROR(err.Error())
			return nil, err
		}
		rt.RouteCode = *num
		// if _, ok := s.routeKeys[rt.RouteCode]; !ok {
		// 	s.routeKeys[rt.RouteCode] = rt.RouteCode
		// }

		rt.RouteDescr = strings.TrimSpace(recordArr[2])
		rt.RouteDescrEng = strings.TrimSpace(recordArr[3])
		num, err = utils.StrToInt32(strings.TrimSpace(recordArr[4]))
		if err != nil {
			opswLogger.ERROR(err.Error())
			return nil, err
		}
		rt.RouteType = int8(*num)
		fl32 := utils.StrToFloat32(strings.TrimSpace(recordArr[5]))
		rt.RouteDistance = fl32

		// s.HelpRoute = append(s.HelpRoute, rt)
		routesRec = append(routesRec, rt)
	}
	return routesRec, nil
}

func insertToMongo(ctx context.Context, collection *mongo.Collection, arr []interface{}, lineCode int32) error {
	filter := bson.M{"linecode": lineCode}
	cursor, err := collection.Find(ctx, filter)
	if err != nil {
		// opswLogger.ERROR(err.Error())
		return err
	}

	for cursor.Next(ctx) {
		var l models.Line
		if err := cursor.Decode(&l); err != nil {
			// opswLogger.ERROR(err.Error())
			return err
		}

		update := bson.M{
			"$set": bson.M{
				"routes01": arr,
			},
		}
		_, err = collection.UpdateOne(ctx, filter, update)
		if err != nil {
			// opswLogger.ERROR(err.Error())
			return err
		}
		opswLogger.INFO(fmt.Sprintf("Line %d updated successfully!", lineCode))
	}
	return nil
}

func recPreparation(recStr string) string {
	// var trimmedSpace = strings.TrimSpace(recStr)
	return strings.ReplaceAll(recStr, "\"", "")
}
