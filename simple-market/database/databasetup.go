package database

//
import (
	"context"
	"fmt"
	"log"
	"time"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func DBSet() *mongo.Client {

	serverAPIOptions := options.ServerAPI(options.ServerAPIVersion1)
clientOptions := options.Client().
    ApplyURI("mongodb+srv://faiksdb:PassIsStrong1@cluster0.hiuvuol.mongodb.net/?retryWrites=true&w=majority").
    SetServerAPIOptions(serverAPIOptions)
ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
defer cancel()
client, err := mongo.Connect(ctx, clientOptions)
if err != nil {
    log.Fatal(err)
}
	err = client.Ping(context.TODO(), nil)
	if err != nil {
		log.Println("failed to connect to mongodb")
		return nil
	}
	fmt.Println("Successfully Connected to the mongodb")
	return client 
}
var Client *mongo.Client = DBSet()
func UserData(client *mongo.Client, CollectionName string) *mongo.Collection {
	var collection *mongo.Collection = client.Database("Market").Collection(CollectionName)
	return collection
}
func ProductData(client *mongo.Client, CollectionName string) *mongo.Collection {
	var productcollection *mongo.Collection = client.Database("Market").Collection(CollectionName)
	return productcollection
}
