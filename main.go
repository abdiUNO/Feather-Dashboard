package main

import (
	"fmt"
	"net/http"
	"os"

	"github.com/abdullahi/go-api/controllers"

	"github.com/gorilla/mux"
)

func main() {

	router := mux.NewRouter()

	fs := http.FileServer(http.Dir("assets/"))
	router.PathPrefix("/static/").Handler(http.StripPrefix("/static/", fs))

	router.HandleFunc("/posts", controllers.PostsIndex).Methods("GET")
	router.HandleFunc("/posts/{id}", controllers.ShowPost).Methods("GET")

	router.HandleFunc("/posts/new", controllers.CreatePost).Methods("POST")
	router.HandleFunc("/posts/{id}", controllers.AddComment).Methods("POST")

	router.HandleFunc("/users", controllers.ShowFakeUsers).Methods("GET")
	router.HandleFunc("/users/new", controllers.CreateFakeUser).Methods("POST")

	// router.HandleFunc("/api/user/new", controllers.CreateAccount).Methods("POST")
	// router.HandleFunc("/api/user/login", controllers.Authenticate).Methods("POST")

	//router.Use(app.JwtAuthentication)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8000"
	}

	fmt.Println(port)

	err := http.ListenAndServe(":"+port, router)

	if err != nil {
		fmt.Print(err)
	}
}
