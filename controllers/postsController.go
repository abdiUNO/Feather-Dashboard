package controllers

import (
	"fmt"
	"html/template"
	"net/http"
	"strings"

	"github.com/abdullahi/go-api/models"
	uuid "github.com/satori/go.uuid"

	u "github.com/abdullahi/go-api/utils"

	"github.com/gorilla/mux"
	"github.com/gorilla/sessions"
	"github.com/unrolled/render"
)

var store = sessions.NewCookieStore([]byte("SESSION_KEY"))

type BreadCrumb struct {
	Link  string
	Last  bool
	Value string
}

var GetPost = func(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	id := params["id"]

	data := models.GetPost(string(id))
	resp := u.Message(true, "success")

	resp["data"] = data
	u.Respond(w, resp)
}

var CreatePost = func(w http.ResponseWriter, r *http.Request) {
	session, _ := store.Get(r, "messages")

	session.Options.MaxAge = 1

	session.AddFlash("Success! Created new post")
	session.Save(r, w)

	r.ParseMultipartForm(128 << 20)

	file, header, err := r.FormFile("input-file-preview")
	post := &models.Post{}

	if err != nil && header.Size > 0 {
		filename := uuid.Must(uuid.NewV4()).String()
		defer file.Close()

		u.SaveFile(file, "/tmp/"+filename)
		u.UploadFile(string(filename))
		post.Image = filename
	}

	post.Text = r.Form.Get("text")
	post.Category = r.Form.Get("category")
	post.UserID = r.Form.Get("fakeuser")

	post.Create()

	http.Redirect(w, r, "/posts", 302)
}

var PostsIndex = func(w http.ResponseWriter, r *http.Request) {

	ren := render.New(render.Options{
		Directory: "templates",
		Layout:    "layout",
		Funcs: []template.FuncMap{
			{
				"myCustomFunc": u.TimeElapsed,
			},
		},
	})

	session, _ := store.Get(r, "messages")

	urlPart := strings.Split(r.URL.Path, "/")

	breadCrumbs := make([]*BreadCrumb, 0)

	for key, value := range urlPart {
		crumb := new(BreadCrumb)
		crumb.Link = value
		crumb.Last = (len(urlPart) - 1) == key
		crumb.Value = strings.Title(value)
		breadCrumbs = append(breadCrumbs, crumb)
	}

	data := map[string]interface{}{
		"Posts":           models.GetPosts(),
		"FlashedMessages": session.Flashes(),
		"BreadCrumbs":     breadCrumbs,
		"FakeUsers":       models.GetFakerUsers(),
		"Active":          "posts",
	}

	err := ren.HTML(w, http.StatusOK, "posts/index", data)
	if err != nil {
		u.Respond(w, u.Message(false, "Invalid request"))
		return
	}
}

var ShowPost = func(w http.ResponseWriter, r *http.Request) {

	ren := render.New(render.Options{
		Directory: "templates",
		Layout:    "layout",
		Funcs: []template.FuncMap{
			{
				"myCustomFunc": u.TimeElapsed,
			},
		},
	})

	session, _ := store.Get(r, "messages")

	params := mux.Vars(r)
	id := params["id"]

	urlPart := strings.Split(r.URL.Path, "/")

	breadCrumbs := make([]*BreadCrumb, 0)

	for key, value := range urlPart {
		crumb := new(BreadCrumb)
		crumb.Link = value
		crumb.Last = (len(urlPart) - 1) == key
		crumb.Value = strings.Title(value)
		breadCrumbs = append(breadCrumbs, crumb)
	}

	data := map[string]interface{}{
		"Post":            models.GetPost(string(id)),
		"Comments":        models.GetComments(string(id)),
		"FlashedMessages": session.Flashes(),
		"BreadCrumbs":     breadCrumbs,
		"FakeUsers":       models.GetFakerUsers(),
		"Active":          "posts",
	}

	err := ren.HTML(w, http.StatusOK, "posts/show", data)
	if err != nil {
		fmt.Println(err)
		u.Respond(w, u.Message(false, "Invalid request"))
		return
	}

}

var AddComment = func(w http.ResponseWriter, r *http.Request) {
	session, _ := store.Get(r, "messages")

	session.Options.MaxAge = 1

	session.AddFlash("Success! Added new comment")
	session.Save(r, w)

	params := mux.Vars(r)
	id := params["id"]

	r.ParseForm()

	comment := &models.Comment{}
	comment.Text = r.Form.Get("text")
	comment.UserID = r.Form.Get("fakeuser")
	comment.PostID = string(id)

	comment.Create()

	http.Redirect(w, r, "/posts/"+id, 302)
}
