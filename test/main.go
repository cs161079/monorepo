package main

import (
	"context"
	"fmt"
	"os"
	"sort"
	"strings"
	"time"

	"github.com/cs161079/monorepo/common/mapper"
	"github.com/cs161079/monorepo/common/models"
	"github.com/cs161079/monorepo/common/service"
	"github.com/cs161079/monorepo/common/utils"
	logger "github.com/cs161079/monorepo/common/utils/goLogger"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

var opswLogger logger.OpswLogger

func main() {
	// Depedencies
	restSrv := service.NewRestService()
	opswLogger = logger.CreateLogger()

	f, err := os.OpenFile("getLines_response.txt", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	defer f.Close()
	if err != nil {
		opswLogger.ERROR(err.Error())
		return
	}

	if _, err := f.WriteString("Appended line.\n"); err != nil {
		opswLogger.ERROR(err.Error())
		return
	}
	intVal, err := utils.StrToInt64("020")
	opswLogger.INFO(fmt.Sprintf("In value %d", *intVal))

	response := restSrv.OasaRequestApi02("getLines")
	if response.Error != nil {
		opswLogger.ERROR(response.Error.Error())
		return
	}

	for _, rec := range response.Data.([]string) {
		f.WriteString(fmt.Sprintf("%v\n", rec))
	}
	return
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
