package config

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/cs161079/monorepo/common/models"
	"github.com/cs161079/monorepo/common/service"
	logger "github.com/cs161079/monorepo/common/utils/goLogger"
)

type TripPlannerService interface {
	IntializeService()
	AgencyData() error
	// =================================================================================================================
	// Με αυτή τη διαδικασία θα παράξουμε το αντίστοιχο αρχείο με τα δεδομένα των στάσεων για
	// το GTFS. Τα δεδομένα θα τα πάρουμε από την βάση δεδομένων.
	//
	// Τα δεδομένα στο σρχείο πρέπει να έχοθυν την μορφή
	// stop_id,stop_name,stop_lat,stop_lon
	// 2001,Σύνταγμα,37.9755,23.7348
	// 2002,Καλλιθέα,37.9501,23.7003
	// =================================================================================================================
	StopsData() error
	RoutesData() error
	CalendarData() error
	TripsData() error
	StopTimesData() error

	writeStopFile(stopRec []StopGTFS) error
	writeRouteFile(routeRecs []RouteGTFS) error
	writeCalendarFile([]CalendarGTFS) error
	writeTripsFile([]TripGTFS) error
	writeStopTimesFile([]StopTimesGTFS) error
}

type tripPlannerServiceImp struct {
	gtfsFolder    string
	stopSrv       service.StopService
	routeSrv      service.RouteService
	scheduleSrv   service.ScheduleService
	schedule01Srv service.Schedule01Service
}

func NewSyncService(stopSrv service.StopService, routeSrv service.RouteService,
	scheduleSrv service.ScheduleService, sched01Srv service.Schedule01Service) TripPlannerService {
	return &tripPlannerServiceImp{
		stopSrv:       stopSrv,
		routeSrv:      routeSrv,
		scheduleSrv:   scheduleSrv,
		schedule01Srv: sched01Srv,
	}
}

func (s *tripPlannerServiceImp) IntializeService() {
	logger.INFO("Initialization of Trip Planner Service is processed...")
	s.gtfsFolder = os.Getenv("gtfs.location")
	if s.gtfsFolder == "" {
		s.gtfsFolder = "gtfs"
	}
	logger.INFO("GTFS structed file stored in " + s.gtfsFolder + ".")
}

func (s *tripPlannerServiceImp) AgencyData() error {
	// Define folder and file path
	// folderPath := "gtfs"
	fileName := "stops.txt"
	fullPath := filepath.Join(s.gtfsFolder, fileName)

	// Create folder if it doesn't exist
	err := os.MkdirAll(s.gtfsFolder, os.ModePerm)
	if err != nil {
		return err
	}
	// Create or truncate the file
	file, err := os.Create(fullPath)
	if err != nil {
		return err
	}
	defer file.Close()

	writeLine(file, "agency_id,agency_name,agency_url,agency_timezone,agency_lang,agency_phone,agency_fare_url,agency_email")
	writeLine(file, "OASA,OASA - Athens Urban Transport Organization,https://www.oasa.gr,Europe/Athens,el,+30 11185,https://www.oasa.gr/en/tickets/prices/,info@oasa.gr")
	return nil
}

func (s *tripPlannerServiceImp) StopsData() error {

	data, err := s.stopSrv.SelectStops()
	if err != nil {
		return err
	}
	var dto []StopGTFS = make([]StopGTFS, 0)
	dataBytes, err := json.Marshal(data)
	if err != nil {
		return err
	}
	err = json.Unmarshal(dataBytes, &dto)
	if err != nil {
		return err
	}
	s.writeStopFile(dto)
	return nil
}

func (s *tripPlannerServiceImp) RoutesData() error {
	data, err := s.routeSrv.RouteList()
	if err != nil {
		return err
	}

	var dto []RouteGTFS = make([]RouteGTFS, 0)
	dataBytes, err := json.Marshal(data)
	if err != nil {
		return err
	}
	err = json.Unmarshal(dataBytes, &dto)
	if err != nil {
		return err
	}
	if err = s.writeRouteFile(dto); err != nil {
		return err
	}
	return nil
}

func (s *tripPlannerServiceImp) CalendarData() error {
	data, err := s.scheduleSrv.ScheduleMasterList()
	if err != nil {
		return err
	}

	var finalResult []CalendarGTFS = make([]CalendarGTFS, 0)

	for _, rec := range data {
		currRes, err := createCalendar(rec)
		if err != nil {
			return err
		}
		finalResult = append(finalResult, currRes...)
	}

	if err = s.writeCalendarFile(finalResult); err != nil {
		return err
	}

	return nil

}

func prepareDaysData(input string) string {
	var result []string

	var sunday = "-1"
	for i, ch := range input {
		if i == 1 {
			sunday = string(ch)
		} else {
			result = append(result, string(ch))
		}
	}
	result = append(result, sunday)
	return strings.Join(result, ",")
}

func LastDayOfMonth(year int, month time.Month) time.Time {
	// Set to the first day of the next month, then subtract 1 day
	t := time.Date(year, month+1, 0, 0, 0, 0, 0, time.UTC)
	return t
}

func createCalendar(inputRec models.ScheduleMaster) ([]CalendarGTFS, error) {
	var result []CalendarGTFS = make([]CalendarGTFS, 0)

	var layout = "2006-01-02"
	var dateNow = time.Now()

	var serviceIdPart string = strconv.Itoa(int(inputRec.SDCCode))
	var daysStrData = prepareDaysData(inputRec.SDCDays)

	var input = inputRec.SDCMonths

	var start *time.Time
	for i, ch := range input {
		if ch == '1' {
			if start == nil {
				// start = i + 1 // months are 1-indexed (January = 1)
				strDate, err := time.Parse(layout, fmt.Sprintf("%d-%02d-%02d", dateNow.Year(), i+1, 1))
				if err != nil {
					return nil, err
				}
				start = &strDate
			}
		} else {
			if start != nil {
				// endDate, err := time.Parse(layout, fmt.Sprintf("%d-%02d-%02d", dateNow.Year(), i, 30))
				// if err != nil {
				// 	// panic(err.Error())
				// 	return nil, err
				// }
				endDate := LastDayOfMonth(dateNow.Year(), time.Month(i))
				var end *time.Time = &endDate
				result = append(result, CalendarGTFS{
					ServiceId: fmt.Sprintf("%s_%d", serviceIdPart, len(result)),
					Days:      daysStrData,
					StartDate: *start,
					EndDate:   *end,
				})
				start = nil
			}
		}
	}

	// Check for a run that ends at the last month
	if start != nil {
		// endDate, err := time.Parse(layout, fmt.Sprintf("%d-%02d-%02d", dateNow.Year(), len(input), 1))
		// if err != nil {
		// 	return nil, err
		// }
		lastDate := LastDayOfMonth(dateNow.Year(), time.Month(len(input)))
		result = append(result, CalendarGTFS{
			ServiceId: fmt.Sprintf("%s_%d", serviceIdPart, len(result)),
			Days:      daysStrData,
			StartDate: *start,
			EndDate:   lastDate,
		})
	}
	return result, nil
}

func countConsecutiveOnesGroups(s string) int {
	count := 0
	inGroup := false

	for _, ch := range s {
		if ch == '1' {
			if !inGroup {
				count++
				inGroup = true
			}
		} else {
			inGroup = false
		}
	}

	return count
}

func (s *tripPlannerServiceImp) TripsData() error {
	routeData, err := s.routeSrv.RouteList()
	if err != nil {
		return err
	}

	var result []TripGTFS = make([]TripGTFS, 0)

	// var currenLine = -1
	for _, rec := range routeData {
		parts := strings.Split(rec.RouteDescr, " - ")
		var head = ""
		if len(parts) >= 1 {
			head = strings.Trim(parts[0], " ")
		}

		schedules, err := s.scheduleSrv.ScheduleMasterDistinct(rec.LnCode)
		if err != nil {
			return err
		}

		for _, rec02 := range schedules {
			var count = countConsecutiveOnesGroups(rec02.SDCMonths)
			var direction = rec.RouteType
			if direction == 2 {
				direction = 0
			}
			scheduleTimes, err := s.schedule01Srv.ScheduleTimeList(rec02.LnCode, rec02.SDCCd, int(direction))
			if err != nil {
				return err
			}
			for _, rec03 := range scheduleTimes {
				for k := 0; k <= count-1; k++ {
					var srvId = fmt.Sprintf("%d_%d", rec02.SDCCd, k)
					result = append(result, TripGTFS{
						RouteId:   int(rec.RouteCode),
						ServiceId: srvId,
						TripId:    fmt.Sprintf("%d_%s_%d_%d", int(rec.RouteCode), srvId, rec03.Sort, k),
						TripHead:  head,
					})
				}
			}
		}
	}

	if err := s.writeTripsFile(result); err != nil {
		return err
	}

	return nil
}

func (s *tripPlannerServiceImp) StopTimesData() error {
	var timeLayout = "15:04:05"
	var result []StopTimesGTFS = make([]StopTimesGTFS, 0)

	routeData, err := s.routeSrv.RouteList()
	if err != nil {
		return err
	}

	var currenLine int32 = -1
	for _, routeRec := range routeData {
		var scheduleMaster []models.ScheduleTimeDto
		if currenLine != routeRec.LnCode {
			currenLine = routeRec.LnCode
			scheduleMaster, err = s.scheduleSrv.ScheduleMasterDistinct(int32(currenLine))
			if err != nil {
				return err
			}
		}

		routeStops, err := s.routeSrv.RouteStopList(int32(routeRec.RouteCode))
		if err != nil {
			return err
		}

		var direction int = int(routeRec.RouteType)
		if direction == 2 {
			direction = 0
		}
		for _, schedMasterRec := range scheduleMaster {
			var count = countConsecutiveOnesGroups(schedMasterRec.SDCMonths)

			schedTime, err := s.schedule01Srv.ScheduleTimeList(int(routeRec.LnCode), schedMasterRec.SDCCd, direction)
			if err != nil {
				return err
			}

			for k := 0; k <= count-1; k++ {
				var srvId = fmt.Sprintf("%d_%d", schedMasterRec.SDCCd, k)
				for _, schedTimeRec := range schedTime {
					duration := time.Time(schedTimeRec.EndTime).Sub(time.Time(schedTimeRec.StartTime))
					intervals := len(routeStops) - 1
					var timeSpace = duration.Seconds() / float64(intervals)
					for _, stop := range routeStops {
						result = append(result, StopTimesGTFS{
							TripId:        fmt.Sprintf("%d_%s_%d_%d", int(routeRec.RouteCode), srvId, schedTimeRec.Sort, k),
							StopSeq:       strconv.Itoa(int(stop.Senu)),
							StopId:        strconv.Itoa(int(stop.StpCode)),
							ArrivalTime:   time.Time(schedTimeRec.StartTime).Add(time.Duration(k*int(timeSpace)) * time.Second).Format(timeLayout),
							DepartureTime: time.Time(schedTimeRec.StartTime).Add(time.Duration(k*int(timeSpace)) * time.Second).Format(timeLayout),
						})
					}

				}
			}
		}

	}

	// // var currentRoute = -1
	// // var routeRec *models.Route
	// // for i, rec := range routeData {
	// // 	if i == 0 {
	// // 		continue
	// // 	}
	// // 	var routeId, err = strconv.Atoi(strings.Trim(row[0], " "))
	// // 	if err != nil {
	// // 		return err
	// // 	}
	// // 	var tripId = row[2]
	// // 	var routeStops []models.Route02

	// // 	parts := strings.Split(row[1], "_")
	// // 	sdc_code, err := strconv.Atoi(strings.Trim(parts[0], " "))
	// // 	if err != nil {
	// // 		return err
	// // 	}
	// // 	var stopCount int
	// // 	if currentRoute != routeId {
	// // 		currentRoute = routeId
	// // 		routeRec, err = s.routeSrv.RouteSelect(int32(routeId))
	// // 		if err != nil {
	// // 			return err
	// // 		}
	// // 		routeStops, err = s.routeSrv.RouteStopList(int32(routeId))
	// // 		if err != nil {
	// // 			return err
	// // 		}
	// // 		stopCount = len(routeStops)
	// // 	}
	// // 	var direction int
	// // 	if routeRec.RouteType == 1 {
	// // 		direction = 1
	// // 	} else {
	// // 		direction = 0
	// // 	}
	// // 	schedTime, err := s.schedule01Srv.ScheduleTimeList(int(routeRec.LnCode), sdc_code, direction)
	// // 	for _, recTime := range schedTime {
	// // 		duration := time.Time(recTime.EndTime).Sub(time.Time(recTime.StartTime))
	// // 		intervals := stopCount - 1
	// // 		var timeSpace = duration.Seconds() / float64(intervals)
	// // 		for k, stop := range routeStops {
	// // 			result = append(result, StopTimesGTFS{
	// // 				TripId:        tripId,
	// // 				StopSeq:       strconv.Itoa(int(stop.Senu)),
	// // 				StopId:        strconv.Itoa(int(stop.StpCode)),
	// // 				ArrivalTime:   time.Time(recTime.StartTime).Add(time.Duration(k*int(timeSpace)) * time.Second).Format(timeLayout),
	// // 				DepartureTime: time.Time(recTime.StartTime).Add(time.Duration(k*int(timeSpace)) * time.Second).Format(timeLayout),
	// // 			})
	// // 		}
	// // 	}

	// }

	err = s.writeStopTimesFile(result)
	if err != nil {
		return err
	}

	return nil
}

func (s *tripPlannerServiceImp) writeStopFile(stopRec []StopGTFS) error {

	// Define folder and file path
	// folderPath := "gtfs"
	fileName := "stops.txt"
	fullPath := filepath.Join(s.gtfsFolder, fileName)

	// Create folder if it doesn't exist
	err := os.MkdirAll(s.gtfsFolder, os.ModePerm)
	if err != nil {
		return err
	}
	// Create or truncate the file
	file, err := os.Create(fullPath)
	if err != nil {
		return err
	}
	defer file.Close()

	// Write string to the file
	writeLine(file, "stop_id,stop_name,stop_lat,stop_lon")
	for _, rec := range stopRec {
		if rec.StopLat != 0 && rec.StopLng != 0 {
			writeLine(file, fmt.Sprintf("%d, %s, %f, %f", rec.StopCode, rec.StopDescr, rec.StopLat, rec.StopLng))
		}

	}
	return nil
}

func (s tripPlannerServiceImp) writeRouteFile(routeRecs []RouteGTFS) error {
	// Define folder and file path
	// folderPath := g
	fileName := "route.txt"
	fullPath := filepath.Join(s.gtfsFolder, fileName)

	// Create folder if it doesn't exist
	err := os.MkdirAll(s.gtfsFolder, os.ModePerm)
	if err != nil {
		return err
	}
	// Create or truncate the file
	file, err := os.Create(fullPath)
	if err != nil {
		return err
	}
	defer file.Close()

	// Write string to the file
	writeLine(file, "route_id,agency_id,route_short_name,route_long_name,route_type")
	for _, rec := range routeRecs {
		writeLine(file, fmt.Sprintf("%d, %s, %s, %s, %d", rec.RouteCode, "OASA", rec.LineID, rec.RouteDescr, 3))
	}
	return nil
}

func (s tripPlannerServiceImp) writeCalendarFile(calRecs []CalendarGTFS) error {
	// Define folder and file path
	var layout = "20060102"
	// folderPath := "gtfs"
	fileName := "calendar.txt"
	fullPath := filepath.Join(s.gtfsFolder, fileName)

	// Create folder if it doesn't exist
	err := os.MkdirAll(s.gtfsFolder, os.ModePerm)
	if err != nil {
		return err
	}
	// Create or truncate the file
	file, err := os.Create(fullPath)
	if err != nil {
		return err
	}
	defer file.Close()

	// Write string to the file
	writeLine(file, "service_id,monday,tuesday,wednesday,thursday,friday,saturday,sunday,start_date,end_date")
	for _, rec := range calRecs {
		writeLine(file, fmt.Sprintf("%s, %s, %s, %s", rec.ServiceId, rec.Days, rec.StartDate.Format(layout), rec.EndDate.Format(layout)))
	}
	return nil
}

func (s tripPlannerServiceImp) writeTripsFile(recs []TripGTFS) error {
	// folderPath := "gtfs"
	fileName := "trips.txt"
	fullPath := filepath.Join(s.gtfsFolder, fileName)

	// Create folder if it doesn't exist
	err := os.MkdirAll(s.gtfsFolder, os.ModePerm)
	if err != nil {
		return err
	}
	// Create or truncate the file
	file, err := os.Create(fullPath)
	if err != nil {
		return err
	}
	defer file.Close()

	// Write string to the file
	writeLine(file, "route_id,service_id,trip_id,trip_headsign")
	for _, rec := range recs {
		writeLine(file, fmt.Sprintf("%d, %s, %s, %s", rec.RouteId, rec.ServiceId, rec.TripId, rec.TripHead))
	}
	return nil
}

func (s *tripPlannerServiceImp) writeStopTimesFile(recs []StopTimesGTFS) error {
	// folderPath := "gtfs"
	fileName := "stop_times.txt"
	fullPath := filepath.Join(s.gtfsFolder, fileName)

	// Create folder if it doesn't exist
	err := os.MkdirAll(s.gtfsFolder, os.ModePerm)
	if err != nil {
		return err
	}
	// Create or truncate the file
	file, err := os.Create(fullPath)
	if err != nil {
		return err
	}
	defer file.Close()

	// Write string to the file
	writeLine(file, "trip_id,arrival_time,departure_time,stop_id,stop_sequence")
	for _, rec := range recs {
		writeLine(file, fmt.Sprintf("%s, %s, %s, %s, %s", rec.TripId, rec.ArrivalTime, rec.DepartureTime, rec.StopId, rec.StopSeq))
	}
	return nil
}

func writeLine(f *os.File, line string) {
	// Write string to the file
	_, err := f.WriteString(fmt.Sprintf("%s\n", line))
	if err != nil {
		panic(err)
	}
}

type StopGTFS struct {
	StopCode  int32   `json:"stop_code"`
	StopDescr string  `json:"stop_descr" gorm:"column:stop_descr"`
	StopLat   float64 `json:"stop_lat" gorm:"column:stop_lat"`
	StopLng   float64 `json:"stop_lng" gorm:"column:stop_lng"`
}

type RouteGTFS struct {
	RouteCode  int32  `json:"route_code"`
	LineID     string `json:"line_id"`
	RouteDescr string `json:"route_descr"`
}

type CalendarGTFS struct {
	ServiceId string    `json:"service_id"`
	Days      string    `json:"days"`
	StartDate time.Time `json:"start_date"`
	EndDate   time.Time `json:"end_date"`
}

// route_id,service_id,trip_id,trip_headsign
type TripGTFS struct {
	RouteId   int    `json:"route_id"`
	ServiceId string `json:"service_id"`
	TripId    string `json:"trip_id"`
	TripHead  string `json:"trip_head"`
}

// trip_id,arrival_time,departure_time,stop_id,stop_sequence
type StopTimesGTFS struct {
	TripId        string `json:"trip_id"`
	ArrivalTime   string `json:"arrival_time"`
	DepartureTime string `json:"departure_time"`
	StopId        string `json:"stop_id"`
	StopSeq       string `json:"stop_sequence"`
}
