package main

import (
    "net/http"
    "io"

    "context"
    "fmt"

    "encoding/json"

    "go.mongodb.org/mongo-driver/bson"
    "go.mongodb.org/mongo-driver/bson/primitive"
    "go.mongodb.org/mongo-driver/mongo"
    "go.mongodb.org/mongo-driver/mongo/options"	
)

const uri = "mongodb://localhost"

type Book struct {
    Id primitive.ObjectID `bson:"_id"`
    Book string `bson:"book"`
    Author_id primitive.ObjectID `bson:"author_id"`
    Published_year int `bson:"published_year"`
    Published_city string `bson:"published_city"`
}

func ShowBooks(writer http.ResponseWriter, request *http.Request) {
    serverAPI := options.ServerAPI(options.ServerAPIVersion1)
    opts := options.Client().ApplyURI(uri).SetServerAPIOptions(serverAPI)

    client, err := mongo.Connect(context.TODO(), opts)
    if err != nil {
        panic(err)
    }

    defer func() {
        if err = client.Disconnect(context.TODO()); err != nil {
            panic(err)
        }
    }()

    coll := client.Database("maindb").Collection("books")

    filter := bson.D{{ "book", bson.D{{ "$exists", true}} }}

    cursor, err := coll.Find(context.TODO(), filter)
    if err != nil {
        panic(err)
    }

    var results []Book
    if err = cursor.All(context.TODO(), &results); err != nil {
        panic(err)
    }

    for _, result := range results {
        cursor.Decode(&result)
        output, err := json.MarshalIndent(result, "", "    ")
        if err != nil {
            panic(err)
        }

        io.WriteString(writer, fmt.Sprintf("%s\n", output))
    }
}

func main() {
    http.HandleFunc("/books/", ShowBooks)

    err := http.ListenAndServe(":81", nil)
    if (err != nil) {
        Printfln("Error: %v", err.Error())
    }
}
