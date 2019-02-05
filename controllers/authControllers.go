package controllers

import (
	"encoding/json"
	"html/template"
	"net/http"
	"strings"

	u "github.com/abdullahi/go-api/utils"
	uuid "github.com/satori/go.uuid"
	"github.com/unrolled/render"

	"github.com/abdullahi/go-api/models"
)

var CreateAccount = func(w http.ResponseWriter, r *http.Request) {

	user := &models.User{}
	err := json.NewDecoder(r.Body).Decode(user) //decode the request body into struct and failed if any error occur
	if err != nil {
		u.Respond(w, u.Message(false, "Invalid request"))
		return
	}

	resp := user.Create() //Create account
	u.Respond(w, resp)
}

var Authenticate = func(w http.ResponseWriter, r *http.Request) {

	user := &models.User{}
	err := json.NewDecoder(r.Body).Decode(user) //decode the request body into struct and failed if any error occur
	if err != nil {
		u.Respond(w, u.Message(false, "Invalid request"))
		return
	}

	resp := models.Login(user.Email, user.Password)
	u.Respond(w, resp)
}

var ShowFakeUsers = func(w http.ResponseWriter, r *http.Request) {

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
		"FakeUsers":       models.GetFakerUsers(),
		"FlashedMessages": session.Flashes(),
		"BreadCrumbs":     breadCrumbs,
		"Active":          "users",
	}

	err := ren.HTML(w, http.StatusOK, "fakeusers/index", data)
	if err != nil {
		u.Respond(w, u.Message(false, "Invalid request"))
		return
	}
}

var CreateFakeUser = func(w http.ResponseWriter, r *http.Request) {

	session, _ := store.Get(r, "messages")

	session.Options.MaxAge = 1

	session.AddFlash("Success! Created new user")
	session.Save(r, w)

	r.ParseMultipartForm(128 << 20)

	file, header, err := r.FormFile("input-file-preview")
	user := &models.User{}

	if err != nil || header.Size > 0 {
		filename := uuid.Must(uuid.NewV4()).String()
		if err != nil {
			panic(err)
		}
		defer file.Close()

		u.SaveFile(file, "/tmp/"+filename)
		u.UploadFile(string(filename))
		user.Image = filename
	}

	user.Username = r.Form.Get("username")
	user.Email = strings.ToLower(r.Form.Get("username")) + "fakeuser@unomaha.edu"
	user.Password = "secret"
	user.Subscription = "Gaming,Fitness,Sports,Music,Movies,News,Conservative,Liberal,Business,Art,Science and Engineering,Stories,Anonymous"

	user.Create()

	fakeuser := &models.FakeUser{}
	fakeuser.Username = r.Form.Get("username")
	fakeuser.UserId = user.ID

	fakeuser.Create()

	http.Redirect(w, r, "/users", 302)
}
