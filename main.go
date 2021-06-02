package main

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
)

var db *gorm.DB
var err error

type Booking struct {
	Id      int    `json:"id"`
	User    string `json:"user"`
	Members int    `json:"members"`
}

func homePage(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Welcome to HomePage!")
	fmt.Println("Endpoint Hit: HomePage")
}

var bookings []Booking

func handleRequests() {
	log.Println("Starting development server at http://127.0.0.1:10000/")

	// creates a new instance of a mux router
	myRouter := mux.NewRouter().StrictSlash(true) //created instance of mux(create routes and http handlers)
	myRouter.HandleFunc("/", homePage)
	myRouter.HandleFunc("/new-booking", createNewBooking).Methods("POST")
	myRouter.HandleFunc("/all-bookings", returnAllBookings).Methods("GET")
	myRouter.HandleFunc("/booking/{id}", returnSingleBooking).Methods("GET")
	myRouter.HandleFunc("/booking/{id}", deletesingleBooking).Methods("DELETE")
	myRouter.HandleFunc("/updatebooking/{id}", updatebooking).Methods("PUT")

	log.Fatal(http.ListenAndServe(":10000", myRouter))
}

func createNewBooking(w http.ResponseWriter, r *http.Request) {
	// get the body of our POST request
	// return the string response containing the request body
	reqBody, _ := ioutil.ReadAll(r.Body) //reqbody has string in bytes

	var booking Booking               //Booking is object
	json.Unmarshal(reqBody, &booking) //json encode to structs,&booking is addres where we stores value
	fmt.Printf("users:%s,members:%d", booking.User, booking.Members)

	db.Create(&booking)
	fmt.Println("Endpoint Hit: Creating New Booking")
	json.NewEncoder(w).Encode(booking)
}

func returnAllBookings(w http.ResponseWriter, r *http.Request) {
	bookings := []Booking{}
	db.Find(&bookings)
	fmt.Println("Endpoint Hit: returnAllBookings")
	json.NewEncoder(w).Encode(bookings)
}

func returnSingleBooking(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	key := vars["id"]
	bookings := []Booking{}
	db.Find(&bookings)
	for _, booking := range bookings {
		// string to int
		s, err := strconv.Atoi(key)
		if err == nil {
			if booking.Id == s {
				fmt.Println(booking)
				fmt.Println("Endpoint Hit: Booking No:", key)
				json.NewEncoder(w).Encode(booking)
			}
		}
	}
}

func deletesingleBooking(w http.ResponseWriter, r *http.Request) {
	//vars := mux.Vars(r)
	//key := vars["id"]
	//bookings := []Booking{}
	//
	//for index, booking := range bookings {
	//	// string to int
	//	s, err := strconv.Atoi(key)
	//	if err == nil {
	//		if booking.Id == s {
	//			bookings = append(bookings[:index], bookings[index+1:]...)
	//			fmt.Println("Endpoint Hit:  Delete Booking")
	//			break
	//		}
	//		json.NewEncoder(w).Encode(booking)
	//	}
	//}

	//	ABOVE CODE IS JUST FOR API DELETE METHOD CALL,WILL NOT DELETE RECORD FROM DB
	vars := mux.Vars(r)
	key := vars["id"]
	var booking Booking
	id, _ := strconv.ParseInt(key, 10, 64)
	print(id)
	db.Where("id = ?", id).Delete(&booking)

	json.NewEncoder(w).Encode(booking)
	w.WriteHeader(http.StatusNoContent)
}

func updatebooking(w http.ResponseWriter, r *http.Request) {

	requestBody, _ := ioutil.ReadAll(r.Body)

	var booking Booking
	json.Unmarshal(requestBody, &booking)
	db.Save(&booking)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(booking)
}

func main() {
	// Please define your username and password for MySQL.
	db, err = gorm.Open("mysql", "root:Mysqlroot123@tcp(localhost:3306)/football?charset=utf8&parseTime=True")
	log.Println(db)
	log.Println(err)

	// NOTE: See we’re using = to assign the global var
	// instead of := which would assign it only in this function

	if err != nil {
		//log.Println(db)
		//log.Println(err)
		log.Println("Connection Failed to Open")
	} else {
		log.Println("Connection Established")
	}
	db.AutoMigrate(&Booking{}) //it create table
	handleRequests()           //handle http requests
}
