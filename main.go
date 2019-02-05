package main

import (
	"fmt"
	"html/template"
	"net/http"
	"os"
	"strings"

	"github.com/abdullahi/go-api/controllers"
	"github.com/unrolled/render"

	u "github.com/abdullahi/go-api/utils"
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
	router.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		ren := render.New(render.Options{
			Directory: "templates",
			Layout:    "layout",
			Funcs: []template.FuncMap{
				{
					"myCustomFunc": u.TimeElapsed,
				},
			},
		})

		urlPart := strings.Split(r.URL.Path, "/")

		breadCrumbs := make([]*controllers.BreadCrumb, 0)

		for key, value := range urlPart {
			crumb := new(controllers.BreadCrumb)
			crumb.Link = value
			crumb.Last = (len(urlPart) - 1) == key
			crumb.Value = strings.Title(value)
			breadCrumbs = append(breadCrumbs, crumb)
		}

		data := map[string]interface{}{
			"BreadCrumbs": breadCrumbs,
			"Active":      "",
		}

		err := ren.HTML(w, http.StatusOK, "home", data)
		fmt.Println(err)
		if err != nil {
			u.Respond(w, u.Message(false, "Invalid request"))
			return
		}
	})

	router.NotFoundHandler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ren := render.New(render.Options{
			Directory: "templates",
			Layout:    "layout",
			Funcs: []template.FuncMap{
				{
					"myCustomFunc": u.TimeElapsed,
				},
			},
		})

		breadCrumbs := make([]*controllers.BreadCrumb, 0)

		data := map[string]interface{}{
			"BreadCrumbs": breadCrumbs,
			"Active":      "",
		}

		err := ren.HTML(w, http.StatusOK, "404", data)
		fmt.Println(err)
		if err != nil {
			u.Respond(w, u.Message(false, "Invalid request"))
			return
		}
	})

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
