package main

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var client *mongo.Client
var collection *mongo.Collection

type Person struct {
	ID        primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"`
	Firstname string             `json:"firstname,omitempty" bson:"firstname,omitempty"`
	Lastname  string             `json:"lastname,omitempty" bson:"lastname,omitempty"`
}

func AddPersonEndpoint(response http.ResponseWriter, request *http.Request) {
	response.Header().Add("content-type", "application/json")
	var person Person
	json.NewDecoder(request.Body).Decode(&person)

	ctx, _ := context.WithTimeout(context.Background(), time.Second*10)
	result, _ := collection.InsertOne(ctx, person)
	json.NewEncoder(response).Encode(result)
}

func GetPeopleEndpoint(response http.ResponseWriter, request *http.Request) {
	response.Header().Add("content-type", "application/json")
	var people []Person

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
		var person Person
		cursor.Decode(&person)
		people = append(people, person)
	}
	if err := cursor.Err(); err != nil {
		response.WriteHeader(http.StatusInternalServerError)
		response.Write([]byte("{message: " + err.Error() + " }"))
		return
	}
	json.NewEncoder(response).Encode(people)
}

func GetPersonEndpoint(response http.ResponseWriter, request *http.Request) {
	response.Header().Add("content-type", "application/json")
	var person Person

	//Get ID param from request
	params := mux.Vars(request)
	id, _ := primitive.ObjectIDFromHex(params["id"])

	ctx, _ := context.WithTimeout(context.Background(), time.Second*10)
	err := collection.FindOne(ctx, Person{ID: id}).Decode(&person)
	if err != nil {
		response.WriteHeader(http.StatusInternalServerError)
		response.Write([]byte("{message: " + err.Error() + " }"))
		return
	}
	json.NewEncoder(response).Encode(person)
}

func main() {
	fmt.Println("Running...")

	ctx, _ := context.WithTimeout(context.Background(), time.Second*10)
	client, _ = mongo.Connect(ctx, options.Client().ApplyURI("mongodb://localhost:27017"))
	collection = client.Database("emblema").Collection("people")

	router := mux.NewRouter()
	router.HandleFunc("/insertPerson", AddPersonEndpoint).Methods("POST")
	router.HandleFunc("/get", GetPeopleEndpoint).Methods("GET")
	router.HandleFunc("/getOne", GetPersonEndpoint).Methods("GET")
	http.ListenAndServe("localhost:12345", router)
}
