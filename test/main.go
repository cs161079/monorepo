package main

import (
	"context"
	"fmt"
	"log"
	"sort"
	"strings"
	"time"

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

	// Depedencies
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	restSrv := service.NewRestService()
	opswLogger = logger.CreateLogger()
	lineSrv := service.NewLineService(nil)

	// Connect to MongoDB
	client, err := mongo.Connect(ctx, options.Client().ApplyURI("mongodb://localhost:27017"))
	if err != nil {
		log.Fatal(err)
	}
	defer client.Disconnect(ctx)

	// Select database and collection
	collection := client.Database("testdb").Collection("line")

	response := restSrv.OasaRequestApi00("webGetLinesWithMLInfo", nil)
	if response.Error != nil {
		opswLogger.ERROR(response.Error.Error())
		return
	}
	// TODO: Το έκοψα γιατί δεν θα κάνω εδώ το Delete
	opswLogger.INFO("\tFetch lines data completed successfully.")
	var lineArray []interface{} = make([]interface{}, 0)
	for _, ln := range response.Data.([]any) {
		lineOasa := lineSrv.GetMapper().GenDtLineOasa(ln.(map[string]interface{}))
		line := lineSrv.GetMapper().OasaToLine(lineOasa)
		lineArray = append(lineArray, line)
	}

	opswLogger.INFO("Going to insert to Database.")
	// Εισαγωγή στην MongoDb
	_, err = collection.InsertMany(ctx, lineArray)
	if err != nil {
		opswLogger.ERROR(err.Error())
	}
	opswLogger.INFO("Inserting to Database completed successfully.")

	response = restSrv.OasaRequestApi02("getRoutes")
	if response.Error != nil {
		opswLogger.ERROR(response.Error.Error())
		return
	}

	var routesRec []models.Route = make([]models.Route, 0)
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
			opswLogger.ERROR("Η γραμμή του Record  είναι ελλειπής.")
			return
		}
		rt := models.Route{}
		num, err := utils.StrToInt32(recordArr[1])
		if err != nil {
			opswLogger.ERROR(err.Error())
			return
		}
		rt.LnCode = *num
		num, err = utils.StrToInt32(recordArr[0])
		if err != nil {
			opswLogger.ERROR(err.Error())
			return
		}
		rt.RouteCode = *num
		// if _, ok := s.routeKeys[rt.RouteCode]; !ok {
		// 	s.routeKeys[rt.RouteCode] = rt.RouteCode
		// }

		rt.RouteDescr = recordArr[2]
		rt.RouteDescrEng = recordArr[3]
		num, err = utils.StrToInt32(recordArr[4])
		if err != nil {
			opswLogger.ERROR(err.Error())
			return
		}
		rt.RouteType = int8(*num)
		fl32 := utils.StrToFloat32(recordArr[5])
		rt.RouteDistance = fl32

		// s.HelpRoute = append(s.HelpRoute, rt)
		routesRec = append(routesRec, rt)
	}

	sort.Slice(routesRec, func(i, j int) bool {
		return routesRec[i].LnCode < routesRec[j].LnCode
	})

	var lnCodCurr int32 = routesRec[0].LnCode
	var lineRoutes []interface{} = make([]interface{}, 0)
	for i, rec := range routesRec {
		if rec.LnCode != int32(lnCodCurr) {
			err = insertToMongo(ctx, collection, lineRoutes, int32(lnCodCurr))
			if err != nil {
				opswLogger.ERROR(err.Error())
				return
			}
			lineRoutes = make([]interface{}, 0)
			lnCodCurr = rec.LnCode
		}

		lineRoutes = append(lineRoutes, routesRec[i])
	}

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
