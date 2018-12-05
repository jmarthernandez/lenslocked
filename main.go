package main

import (
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/lenslocked/controllers"
	"github.com/lenslocked/models"
)

const (
	host     = "localhost"
	port     = 5432
	user     = "jhernandez2"
	password = ""
	dbname   = "lenslocked_dev"
)

func main() {
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s "+
		"password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)

	us, err := models.NewUserService(psqlInfo)
	if err != nil {
		panic(err)
	}
	us.AutoMigrate()
	defer us.Close()
	staticC := controllers.NewStatic()
	usersC := controllers.NewUsers(us)

	r := mux.NewRouter()
	r.Handle("/", staticC.Home).Methods("GET")
	r.Handle("/contact", staticC.Contact).Methods("GET")
	r.HandleFunc("/signup", usersC.New).Methods("GET")
	r.HandleFunc("/signup", usersC.Create).Methods("POST")
	fmt.Println("starting on 3000")
	http.ListenAndServe(":3030", r)
}
