package main

import (
	"fmt"
	"html/template"
	"net/http"
	"strconv"

	mgo "gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"

	"github.com/gorilla/mux"
)

type Item struct {
	Id          bson.ObjectId `json:"id" bson:"_id"`
	Title       string        `json:"title" bson:"title"`
	Description string        `json:"description" bson:"description"`
	Updated     string        `json:"updated" bson:"updated"`
}

type Handler struct {
	Sess  *mgo.Session
	Items *mgo.Collection
	Tmpl  *template.Template
}

func (h *Handler) List(w http.ResponseWriter, r *http.Request) {

	items := []*Item{}

	// bson.M{} - это типа условия для поиска
	err := h.Items.Find(bson.M{}).All(&items)
	__err_panic(err)

	err = h.Tmpl.ExecuteTemplate(w, "index.html", struct {
		Items []*Item
	}{
		Items: items,
	})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (h *Handler) AddForm(w http.ResponseWriter, r *http.Request) {
	err := h.Tmpl.ExecuteTemplate(w, "create.html", nil)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (h *Handler) Add(w http.ResponseWriter, r *http.Request) {

	newItem := bson.M{
		"_id":         bson.NewObjectId(),
		"title":       r.FormValue("title"),
		"description": r.FormValue("description"),
		"some_filed":  123,
	}
	err := h.Items.Insert(newItem)
	__err_panic(err)

	fmt.Println("Insert - LastInsertId:", newItem["id"])

	http.Redirect(w, r, "/", http.StatusFound)
}

func (h *Handler) Edit(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	if !bson.IsObjectIdHex(vars["id"]) {
		http.Error(w, "bad id", 500)
		return
	}
	id := bson.ObjectIdHex(vars["id"])

	post := &Item{}
	err := h.Items.Find(bson.M{"_id": id}).One(&post)

	err = h.Tmpl.ExecuteTemplate(w, "edit.html", post)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func (h *Handler) Update(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	if !bson.IsObjectIdHex(vars["id"]) {
		http.Error(w, "bad id", 500)
		return
	}
	id := bson.ObjectIdHex(vars["id"])

	post := &Item{}
	err := h.Items.Find(bson.M{"_id": id}).One(&post)

	post.Title = r.FormValue("title")
	post.Description = r.FormValue("description")
	post.Updated = "rvasily"

	// err = h.Items.Update(bson.M{"_id": id}, &post)
	err = h.Items.Update(
		bson.M{"_id": id},
		bson.M{
			"title":       r.FormValue("title"),
			"description": r.FormValue("description"),
			"updated":     "rvasily",
			"newField":    123,
		})
	affected := 1
	if err == mgo.ErrNotFound {
		affected = 0
	} else if err != nil {
		__err_panic(err)
	}

	fmt.Println("Update - RowsAffected", affected)

	http.Redirect(w, r, "/", http.StatusFound)
}

func (h *Handler) Delete(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	if !bson.IsObjectIdHex(vars["id"]) {
		http.Error(w, "bad id", 500)
		return
	}
	id := bson.ObjectIdHex(vars["id"])

	err := h.Items.Remove(bson.M{"_id": id})
	affected := 1
	if err == mgo.ErrNotFound {
		affected = 0
	} else if err != nil {
		__err_panic(err)
	}

	w.Header().Set("Content-type", "application/json")
	resp := `{"affected": ` + strconv.Itoa(int(affected)) + `}`
	w.Write([]byte(resp))
}

func main() {
	sess, err := mgo.Dial("mongodb://localhost")
	__err_panic(err)

	// если коллекции не будет, то она создасться автоматически
	collection := sess.DB("coursera").C("items")

	// для монги нет такого красивого дампа SQL, так что я вставляю демо-запись если коллекция пуста
	if n, _ := collection.Count(); n == 0 {
		collection.Insert(&Item{
			bson.NewObjectId(),
			"mongodb",
			"Рассказать про монгу",
			"",
		})
		collection.Insert(&Item{
			bson.NewObjectId(),
			"redis",
			"Рассказать про redis",
			"rvasily",
		})
	}

	handlers := &Handler{
		Items: collection,
		Tmpl:  template.Must(template.ParseGlob("./templates/*")),
	}

	// в целям упрощения примера пропущена авторизация и csrf
	r := mux.NewRouter()
	r.HandleFunc("/", handlers.List).Methods("GET")
	r.HandleFunc("/items", handlers.List).Methods("GET")
	r.HandleFunc("/items/new", handlers.AddForm).Methods("GET")
	r.HandleFunc("/items/new", handlers.Add).Methods("POST")
	r.HandleFunc("/items/{id}", handlers.Edit).Methods("GET")
	r.HandleFunc("/items/{id}", handlers.Update).Methods("POST")
	r.HandleFunc("/items/{id}", handlers.Delete).Methods("DELETE")

	fmt.Println("starting server at :8080")
	http.ListenAndServe(":8080", r)
}

// не используйте такой код в прошакшене
// ошибка должна всегда явно обрабатываться
func __err_panic(err error) {
	if err != nil {
		panic(err)
	}
}
