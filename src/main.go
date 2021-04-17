package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	// "os"
	"time"

	"github.com/emblemaa/Carie/src/model"
	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var client *mongo.Client

type Driver struct {
	Id        primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	Userid    string             `json:"_userid,omitempty" bson:"_userid,omitempty"`
	Firstname string             `json:"firstname,omitempty" bson:"firstname,omitempty"`
	Lastname  string             `json:"lastname,omitempty" bson:"lastname,omitempty"`
}

type Time struct {
	Hour   string `json:"hour,omitempty" bson:"hour, omitempty"`
	Minute string `json:"minute,omitempty" bson:"minute, omitempty"`
}

type Schedule struct {
	Content     string `json:"content,omitempty" bson:"content, omitempty"`
	Description string `json:"description,omitempty" bson:"description, omitempty"`
	PickTime    Time   `json:"picktime,omitempty" bson:"picktime, omitempty"`
	DropTime    Time   `json:"droptime,omitempty" bson:"droptime, omitempty"`
	From        string `json:"from,omitempty" bson:"from,omitempty"`
	To          string `json:"to,omitempty" bson:"to,omitempty"`
	DaysInWeek  string `json:"daysinweek,omitempty" bson:"daysinweek, omitempty"`
	IsEnabled   bool   `json:"isenabled,omitempty" bson:"isenabled, omitempty"`
}

type User struct {
	PhoneNumber  string     `json:"phonenumber,omitempty" bson:"phonenumber,omitempty"`
	ScheduleList []Schedule `json:"schedulelist,omitempty" bson:"schedulelist,omitempty"`
}

func GetProjectByName(response http.ResponseWriter, request *http.Request) {
	response.Header().Add("content-type", "application/json")
	var project model.Project
	params := request.URL.Query()
	dbString := os.Getenv("DB_STRING")
	if dbString == "" {
		dbString = "mongodb://localhost:27017"
	}

	ctx, _ := context.WithTimeout(context.Background(), time.Second*10)
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(dbString))

	if err != nil {
		log.Fatal(err)
	}
	collection := client.Database("carie").Collection("tanbinh_location")

	if len(params) != 0 {
		name := params.Get("name")
		err = collection.FindOne(ctx, model.Project{Name: name}).Decode(&project)
		if err != nil {
			response.WriteHeader(http.StatusInternalServerError)
			response.Write([]byte("{message: " + err.Error() + "}"))
			return
		}
		json.NewEncoder(response).Encode(project)
		return
	} else {
		response.WriteHeader(http.StatusBadRequest)
		response.Write([]byte("{message: Cần cung cấp tên địa điểm" + "}"))
	}
}

func AddDriverEndpoint(response http.ResponseWriter, request *http.Request) {

	response.Header().Add("content-type", "application/json")
	var person Driver
	json.NewDecoder(request.Body).Decode(&person)

	dbString := os.Getenv("DB_STRING")
	if dbString == "" {
		dbString = "mongodb://localhost:27017"
	}

	ctx, _ := context.WithTimeout(context.Background(), time.Second*10)
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(dbString))

	if err != nil {
		log.Fatal(err)
	}

	collection := client.Database("carie").Collection("drivers")
	result, _ := collection.InsertOne(ctx, person)
	json.NewEncoder(response).Encode(result)
}

func GetDriverEndpoint(response http.ResponseWriter, request *http.Request) {
	response.Header().Add("content-type", "application/json")
	dbString := os.Getenv("DB_STRING")
	if dbString == "" {
		dbString = "mongodb://localhost:27017"
	}
	ctx, _ := context.WithTimeout(context.Background(), time.Second*10)
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(dbString))
	params := request.URL.Query()

	if err != nil {
		log.Fatal(err)
	}

	collection := client.Database("carie").Collection("drivers")

	if len(params) == 0 {
		var people []Driver
		ctx, _ := context.WithTimeout(context.Background(), time.Second*10)
		cursor, err := collection.Find(ctx, bson.M{})
		if err != nil {
			response.WriteHeader(http.StatusInternalServerError)
			response.Write([]byte("{message: " + err.Error() + " }"))
			return
		}

		// Close the cursor later when function finishes
		defer cursor.Close(ctx)

		for cursor.Next(ctx) {
			var person Driver
			cursor.Decode(&person)
			people = append(people, person)
		}
		if err := cursor.Err(); err != nil {
			response.WriteHeader(http.StatusInternalServerError)
			response.Write([]byte("{message: " + err.Error() + " }"))
			return
		}
		json.NewEncoder(response).Encode(people)
		return
	}

	var person Driver
	id, _ := primitive.ObjectIDFromHex(params.Get("id"))
	err = collection.FindOne(ctx, Driver{Id: id}).Decode(&person)
	if err != nil {
		response.WriteHeader(http.StatusInternalServerError)
		response.Write([]byte("{message: " + err.Error() + " }"))
		return
	}
	json.NewEncoder(response).Encode(person)
}

func AddUserEndpoint(response http.ResponseWriter, request *http.Request) {

	response.Header().Add("content-type", "application/json")
	var person User
	json.NewDecoder(request.Body).Decode(&person)

	dbString := os.Getenv("DB_STRING")
	if dbString == "" {
		dbString = "mongodb://localhost:27017"
	}

	ctx, _ := context.WithTimeout(context.Background(), time.Second*10)
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(dbString))

	if err != nil {
		log.Fatal(err)
	}

	collection := client.Database("carie").Collection("users-test")

	result, _ := collection.InsertOne(ctx, person)
	json.NewEncoder(response).Encode(result)
}

func GetUserEndpoint(response http.ResponseWriter, request *http.Request) {
	response.Header().Add("content-type", "application/json")
	dbString := os.Getenv("DB_STRING")
	if dbString == "" {
		dbString = "mongodb://localhost:27017"
	}
	ctx, _ := context.WithTimeout(context.Background(), time.Second*10)
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(dbString))
	params := request.URL.Query()

	if err != nil {
		log.Fatal(err)
	}

	collection := client.Database("carie").Collection("users-test")

	if len(params) == 0 {
		var people []User
		ctx, _ := context.WithTimeout(context.Background(), time.Second*10)
		cursor, err := collection.Find(ctx, bson.M{})
		if err != nil {
			response.WriteHeader(http.StatusInternalServerError)
			response.Write([]byte("{message: " + err.Error() + " }"))
			return
		}

		// Close the cursor later when function finishes
		defer cursor.Close(ctx)

		for cursor.Next(ctx) {
			var person User
			cursor.Decode(&person)
			people = append(people, person)
		}
		if err := cursor.Err(); err != nil {
			response.WriteHeader(http.StatusInternalServerError)
			response.Write([]byte("{message: " + err.Error() + " }"))
			return
		}
		json.NewEncoder(response).Encode(people)
		return
	}

	var person User
	phone := (params.Get("phone"))
	err = collection.FindOne(ctx, User{PhoneNumber: phone}).Decode(&person)
	if err != nil {
		response.WriteHeader(http.StatusInternalServerError)
		response.Write([]byte("{message: " + err.Error() + " }"))
		return
	}
	json.NewEncoder(response).Encode(person)
}

// func GetPeopleEndpoint(response http.ResponseWriter, request *http.Request) {
// 	response.Header().Add("content-type", "application/json")
// 	var people []Person

// 	ctx, _ := context.WithTimeout(context.Background(), time.Second*10)
// 	cursor, err := collection.Find(ctx, bson.M{})
// 	if err != nil {
// 		response.WriteHeader(http.StatusInternalServerError)
// 		response.Write([]byte("{message: " + err.Error() + " }"))
// 		return
// 	}

// 	// Close the cursor later when function finishes
// 	defer cursor.Close(ctx)

// 	for cursor.Next(ctx) {
// 		var person Person
// 		cursor.Decode(&person)
// 		people = append(people, person)
// 	}
// 	if err := cursor.Err(); err != nil {
// 		response.WriteHeader(http.StatusInternalServerError)
// 		response.Write([]byte("{message: " + err.Error() + " }"))
// 		return
// 	}
// 	json.NewEncoder(response).Encode(people)
// }

// func GetPersonEndpoint(response http.ResponseWriter, request *http.Request) {
// 	response.Header().Add("content-type", "application/json")
// 	var person Person

// 	vals := request.URL.Query()
// 	fmt.Println(vals["id"])
// 	//Get ID param from request
// 	params := mux.Vars(request)
// 	id, _ := primitive.ObjectIDFromHex(params["id"])

// 	ctx, _ := context.WithTimeout(context.Background(), time.Second*10)
// 	err := collection.FindOne(ctx, Person{ID: id}).Decode(&person)
// 	if err != nil {
// 		response.WriteHeader(http.StatusInternalServerError)
// 		response.Write([]byte("{message: " + err.Error() + " }"))
// 		return
// 	}
// 	json.NewEncoder(response).Encode(person)
// }

// func main_new() {
// 	app := App{}
// 	a.Initialize()
// }

func main() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Printf("Can't load .env files: ")
		log.Println(err)
	}
	fmt.Println("Running...")

	//get config
	port := os.Getenv("PORT")
	if port == "" {
		port = "8085"
	}
	dbString := os.Getenv("DB_STRING")
	if dbString == "" {
		dbString = "mongodb://localhost:27017"
	}

	url := ":" + port
	// fmt.Print(url)
	//init connection
	//ctx, _ := context.WithTimeout(context.Background(), time.Second*10)
	//client, err := mongo.Connect(ctx, options.Client().ApplyURI(dbString))
	if err != nil {
		log.Fatal(err)
	}

	//setup router
	router := mux.NewRouter()
	router.HandleFunc("/user", AddUserEndpoint).Methods("POST")
	router.HandleFunc("/user", GetUserEndpoint).Methods("GET")
	router.HandleFunc("/driver", AddDriverEndpoint).Methods("POST")
	router.HandleFunc("/driver", GetDriverEndpoint).Methods("GET")
	router.HandleFunc("/location", GetProjectByName).Methods("GET")
	fmt.Println("Serving on: " + url)
	http.ListenAndServe(url, router)

}
