package main

import (
"net/http"

"log"

"context"

"os"

"encoding/json"

"go.mongodb.org/mongo-driver/v2/bson"

"go.mongodb.org/mongo-driver/v2/mongo"

"go.mongodb.org/mongo-driver/v2/mongo/options"
)



var collect *mongo.Collection
type User struct {

		ID bson.ObjectID `bson:"_id,omitempty" json:"ID"`

		Text string `bson:"text" json:"text"`
	}


func sendToMe(write http.ResponseWriter, read *http.Request) {

	
 write.Header().Set("Access-Control-Allow-Origin", "*")

  write.Header().Set("Access-Control-Allow-Methods", "POST,OPTIONS")

   write.Header().Set("Access-Control-Allow-Headers", "Content-Type")

   
	if read.Method == "OPTIONS"{
		write.WriteHeader(http.StatusOK)
		return
	}
	defer read.Body.Close();
   var resp User

   err := json.NewDecoder(read.Body).Decode(&resp)



	
   if err != nil {

   	write.WriteHeader(http.StatusBadRequest)
	json.NewEncoder(write).Encode(map[string]string{"message": "Invalid JSON type, send a valid JSON!"})

	return

   }

		
   err = collect.InsertOne(context.TODO(), resp)

   if err != nil {

   	write.WriteHeader(http.StatusInternalServerError)
	json.NewEncoder(write).Encode(map[string]string{"message": "Failed to send message. Check your internet connection or try again later"})

	return
   }

		write.WriteHeader(http.StatusOK)
   json.NewEncoder(write).Encode(map[string]string{"message": "Message successfully sent"})

}



func recieveText(write http.ResponseWriter, read *http.Request) {
	write.Header().Set("Access-Control-Allow-Origin", "*")

  write.Header().Set("Access-Control-Allow-Methods", "GET,OPTIONS")

   write.Header().Set("Access-Control-Allow-Headers", "Content-Type")
	if read.Method == "OPTIONS" {
		write.WriteHeader(http.StatusOK)
		return
	}
	cursor, err := collect.Find(context.TODO(), bson.M{})
	if err != nil {
		write.WriteHeader(http.StatusInternalServerError)
			return
	}
		defer cursor.Close(context.TODO())
		var allUsers []User
		for cursor.Next(context.TODO()) {
				var eachUser User
				err := cursor.Decode(&eachUser)
				if err != nil {
						write.WriteHeader(http.StatusBadRequest)
						return
				}
				
				allUsers = append(allUsers, eachUser)
		}
		
		if err = cursor.Err(); err != nil {
						write.WriteHeader(http.StatusInternalServerError)
						return
				}
		write.Header().Set("Content-Type","application/json")
		write.WriteHeader(http.StatusOK)
		json.NewEncoder(write).Encode(allUsers)
}

func deleteUser(write http.ResponseWriter, read *http.Request) {
		write.Header().Set("Access-Control-Allow-Origin", "*")
		write.Header().Set("Access-Control-Allow-Methods", "DELETE, OPTIONS")
write.Header().Set("Access-Control-Allow-Headers", "Content-Type")
		if read.Method == "OPTIONS" {
				write.WriteHeader(http.StatusOK)
				return
		}
 id := read.URL.Query().Get("id")
		if id == "" {
				write.WriteHeader(http.StatusBadRequest)
				return
		}
		obj, err := bson.ObjectIDFromHex(id)
		if err != nil {
				http.Error(write, "Invalid ID, check your ID", http.StatusBadRequest)
				return
		}
		
		err = collect.DeleteOne(context.TODO(), bson.M{"_id":obj})
		if err != nil {
				write.WriteHeader(http.StatusInternalServerError)
				return
		}
		write.WriteHeader(http.StatusOK)
		write.Header().Set("Content-Type","application/json")
		json.NewEncoder(write).Encode(map[string]string {"message": "successfully deleted"})
}

func main() {
client, fail := mongo.Connect(context.TODO(), options.Client().ApplyURI(os.Getenv("MONGO")))

if fail != nil {

	log.Fatal(fail)
}

defer client.Disconnect(context.TODO())

ers := client.Ping(context.TODO(), nil)

if ers != nil{

	log.Fatal(ers)
}

collect = client.Database("anonymous").Collection("users")


http.HandleFunc("/user", sendToMe)

http.HandleFunc("/admin", recieveText)

http.HandleFunc("/del", deleteUser)
	port := os.Getenv("PORT")
	if port == "" {
		port = "8082"
	}
	http.ListenAndServe(":"+port, nil)
}
