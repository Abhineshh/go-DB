package main

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	"os"
)

func main() {
	// we load the .env file to the inbuild env-map of golang
	if err := godotenv.Load(); err != nil {  
		log.Println("No .env file found")
	}
	uri := os.Getenv("MONGODB_URI")  // here we access the contents of the env file via the env-map
	if uri == "" {
		log.Fatal("You must set your 'MONGODB_URI' environment variable")
	}
	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(uri)) // here we connect to the mongoDB  database using the mongodb uri
	if err != nil {
		panic(err)
	}
	// here we have defer function that will execute at the end to disconnect from the mongodb server
	defer func() {
		if err := client.Disconnect(context.TODO()); err != nil {
			panic(err)
		}
	}()

	coll := client.Database("Rango").Collection("drivers") // getting access to the database and collection that we want
	title := "Vinsmoke"

	var result bson.M    // specify the type of datastore
	// using the findone function to get one json instance from the collection 
	//and storing it into the specified data store
	err = coll.FindOne(context.TODO(), bson.D{{"driverName", title}}).Decode(&result)  
	if err == mongo.ErrNoDocuments {
		fmt.Printf("No document was found with the title%s \n", title)
		return
	}
	if err != nil {
		panic(err)
	}
	//converting the bson.M data to json i.e byte[]
	jsonData, err := json.MarshalIndent(result, "", "   ")  
	if err != nil {
		panic(err)
	}
	fmt.Printf("%s \n", jsonData)

	// using the updateoOne function to set one json instance from the collection

	results,errs := coll.UpdateOne(context.TODO(),bson.D{{"driverName","ghidorah"}}, bson.D{{"$set",bson.D{{"driverName","Vinsmoke"}}}},options.Update().SetUpsert(true))
	if errs != nil {
		log.Fatal(errs)
	}
	if results.MatchedCount != 0 {
		fmt.Println("matched and replaced an existing document")
	}
	if results.UpsertedCount != 0 {
		fmt.Println("inserted a new document with ID %V \n",results.UpsertedID)
	}
}
