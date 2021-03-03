package controllers

import (
	"back/lib"
	"back/src/services"
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"io/ioutil"
	"log"
	"net/http"
)

func GetRoomController (w http.ResponseWriter, req *http.Request) {
	id := mux.Vars(req)["id"]
	objId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		log.Fatal(err)
	}
	room, getError := services.RetrieveRoom(objId)
	if getError != nil {
		w.WriteHeader(404)
		return
	}
	value, _ := json.Marshal(room)
	w.Write(value)
}

func GetRoomsController (w http.ResponseWriter, req *http.Request) {
	fmt.Println("id received for get all")
	rooms, err := services.RetrieveAllRooms()
	if err != nil {
		w.WriteHeader(404)
		return
	}
	a := lib.IRoomsResponse{Rooms: *rooms}
	value, merr := json.Marshal(a)
	value = bytes.Replace(value, []byte(":NaN"), []byte(":null"), -1)
	fmt.Println("value", value)
	if merr != nil {
		panic(merr)
	}
	fmt.Println("rooms", a)
	w.Write(value)
}

func PostRoomController (w http.ResponseWriter, req *http.Request) {
	var room lib.IRoom
	bodyRoom,_ := ioutil.ReadAll(req.Body)
	err := json.Unmarshal(bodyRoom, &room)
	if err != nil {
		fmt.Println("error")
		log.Fatal(err)
	}
	room.Resource = lib.NewResource()
	returnId, serror := services.AddRoomData(room)
	if serror != nil {
		panic(serror)
	}
	w.Header().Add("Location", "http://" +req.Host + req.RequestURI + "/" + returnId.Hex())
	w.WriteHeader(201)
}

func PutRoomController (w http.ResponseWriter, req *http.Request) {
	id := mux.Vars(req)["id"]
	fmt.Println("id received for put", id)
	objId, errId := primitive.ObjectIDFromHex(id)
	if errId != nil {
		log.Fatal(errId)
	}
	var room lib.IRoom
	bodyTrack,_ := ioutil.ReadAll(req.Body)
	err := json.Unmarshal(bodyTrack, &room)
	if err != nil {
		fmt.Println("error")
		log.Fatal(err)
	}
	room.Resource.ID = objId
	returnId, dbError := services.UpdateRoom(room)
	if dbError != nil{
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(`{"message":"`+dbError.Error()+`"}`))
		return
	}
	w.WriteHeader(204)
	w.Header().Add("Location", req.Host + req.RequestURI + "/" + returnId.Hex())
}

func DeleteRoomController (w http.ResponseWriter, req *http.Request) {
	id := mux.Vars(req)["id"]
	fmt.Println("id received for put", id)
	objId, errId := primitive.ObjectIDFromHex(id)
	if errId != nil {
		log.Fatal(errId)
	}
	_, _ = services.RemoveRoom(objId)
	w.WriteHeader(204)
}
