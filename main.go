package main

import (
	"context"
	"fmt"
	"html/template"
	"log"
	"myproject-page/connection"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gorilla/mux"
	"github.com/gorilla/sessions"
	"golang.org/x/crypto/bcrypt"
)

func main() {
	route := mux.NewRouter()

	connection.DatabaseConnect()

	route.PathPrefix("/public/").Handler(http.StripPrefix("/public", http.FileServer((http.Dir("./public")))))

	route.HandleFunc("/", home).Methods("GET")
	route.HandleFunc("/form-project", formProject).Methods("GET")
	route.HandleFunc("/contact", contact).Methods("GET")
	route.HandleFunc("/project-detail/{id}", projectDetail).Methods("GET")
	route.HandleFunc("/add-project", addProject).Methods("POST")
	route.HandleFunc("/delete-project/{index}", deleteProject).Methods("GET")
	route.HandleFunc("/edit-project/{index}", formEditProject).Methods("GET")
	route.HandleFunc("/edit-project/{index}", editProject).Methods("POST")

	route.HandleFunc("/register", formRegister).Methods("GET")
	route.HandleFunc("/register", register).Methods("POST")

	route.HandleFunc("/login", formLogin).Methods("GET")
	route.HandleFunc("/login", login).Methods("POST")

	route.HandleFunc("/logout", logout).Methods("GET")

	fmt.Println(("server berjalan di port 5000"))
	http.ListenAndServe("localhost:5000", route)
}

// type MetaData struct {
// }

// var DataFlash = MetaData{
// 	TitleSessions: "Project Web",
// }

type User struct {
	Id       int
	Name     string
	Email    string
	Password string
}

type Project struct {
	// sessions struct
	TitleSessions string
	IsLogin       bool
	UserName      string
	FlashData     string

	//card project struct
	ID                     int
	Title                  string
	DateStart              time.Time
	DateEnd                time.Time
	Format_date_start      string
	Format_date_start_edit string
	Format_date_end        string
	Description            string
	Technologies           []string
	NodeJs                 string
	ReactJs                string
	NextJs                 string
	Javascript             string
}

var DataFlash = Project{
	TitleSessions: "Project Web",
}

func home(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-type", "text/html; charset=utf-8")
	tmpt, err := template.ParseFiles("views/index.html")

	if err != nil {
		w.Write([]byte("Message :" + err.Error()))
		return
	}

	dataProject, errQuery := connection.Conn.Query(context.Background(), "SELECT id, title, start_date, end_date, description, technologies FROM tb_projects")
	if errQuery != nil {
		fmt.Println("Message : " + errQuery.Error())
		return
	}

	var result []Project

	for dataProject.Next() {
		var each = Project{}

		err := dataProject.Scan(&each.ID, &each.Title, &each.DateStart, &each.DateEnd, &each.Description, &each.Technologies)
		if err != nil {
			fmt.Println("Message : " + err.Error())
			return
		}

		// for i := 0; i < len(each.Technologies); i++ {
		// if each.Technologies[0] == "true" {
		// 	each.NodeJs = "nodejs.svg"
		// }
		// if each.Technologies[1] == "true" {
		// 	each.ReactJs = "react.svg"
		// }
		// if each.Technologies[2] == "true" {
		// 	each.NextJs = "nextjs.svg"
		// }
		// if each.Technologies[3] == "true" {
		// 	each.Javascript = "javascript.svg"
		// }
		//  = each.Technologies[0]
		// }

		// fmt.Println(each.Technologies[0])

		each.Format_date_start = each.DateStart.Format("2 January 2006")
		each.Format_date_end = each.DateEnd.Format("2 January 2006")
		result = append(result, each)

	}

	// sessions
	var store = sessions.NewCookieStore([]byte("SESSIONS_ID"))
	session, _ := store.Get(r, "SESSIONS_ID")

	if session.Values["IsLogin"] != true {
		DataFlash.IsLogin = false
	} else {
		DataFlash.IsLogin = session.Values["IsLogin"].(bool)
		DataFlash.UserName = session.Values["Names"].(string)
	}

	fm := session.Flashes("message")

	var flashes []string

	if len(fm) > 0 {
		session.Save(r, w)

		for _, fl := range fm {
			flashes = append(flashes, fl.(string))
		}
	}

	DataFlash.FlashData = strings.Join(flashes, "")

	resData := map[string]interface{}{
		"Projects":  result,
		"DataFlash": DataFlash,
	}

	tmpt.Execute(w, resData)

}

func formProject(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-type", "text/html; charset=utf-8")
	tmpt, err := template.ParseFiles("views/addProject.html")

	if err != nil {
		w.Write([]byte("Message :" + err.Error()))
		return
	}

	// sessions
	var store = sessions.NewCookieStore([]byte("SESSIONS_ID"))
	session, _ := store.Get(r, "SESSIONS_ID")

	if session.Values["IsLogin"] != true {
		DataFlash.IsLogin = false
	} else {
		DataFlash.IsLogin = session.Values["IsLogin"].(bool)
		DataFlash.UserName = session.Values["Names"].(string)
	}

	fm := session.Flashes("message")

	var flashes []string

	if len(fm) > 0 {
		session.Save(r, w)

		for _, fl := range fm {
			flashes = append(flashes, fl.(string))
		}
	}

	DataFlash.FlashData = strings.Join(flashes, "")

	Data := map[string]interface{}{
		"DataFlash": DataFlash,
		// "DataFlash": DataFlash,
	}

	tmpt.Execute(w, Data)
}

func addProject(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()

	if err != nil {
		log.Fatal(err)
	}

	title := r.PostForm.Get("title")
	content := r.PostForm.Get("content")
	dateStart := r.PostForm.Get("date-start")
	dateEnd := r.PostForm.Get("date-end")

	nodeJs := r.PostForm.Get("nodeJs")
	nextJs := r.PostForm.Get("nextJs")
	reactJs := r.PostForm.Get("reactJs")
	javascript := r.PostForm.Get("javascript")

	checked := []string{

		nodeJs,
		nextJs,
		reactJs,
		javascript,
	}

	_, errQuery := connection.Conn.Exec(context.Background(), "INSERT INTO public.tb_projects(title, start_date, end_date, description, technologies) VALUES ($1, $2, $3, $4, $5)", title, dateStart, dateEnd, content, checked)
	if errQuery != nil {
		fmt.Println("Message : " + errQuery.Error())
		return
	}

	// var newProject = Project{
	// 	DateStartEdit: dateStart,
	// 	DateEndEdit: dateEnd,
	// }

	// projects = append(projects, newProject)

	fmt.Println(dateStart)

	http.Redirect(w, r, "/", http.StatusMovedPermanently)
}

func formEditProject(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-type", "text/html; charset=utf-8")
	tmpt, err := template.ParseFiles("views/editProject.html")

	if err != nil {
		w.Write([]byte("Message :" + err.Error()))
		return
	}

	index, _ := strconv.Atoi(mux.Vars(r)["index"])

	var ProjectEdit = Project{}

	err = connection.Conn.QueryRow(context.Background(), "SELECT id, title, start_date, end_date, description FROM tb_projects WHERE id = $1", index).Scan(&ProjectEdit.ID, &ProjectEdit.Title, &ProjectEdit.DateStart, &ProjectEdit.DateEnd, &ProjectEdit.Description)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("message : " + err.Error()))
	}

	// ProjectEdit.nodeJs = ProjectEdit.Technologies[0]

	ProjectEdit.Format_date_start = ProjectEdit.DateStart.Format("2 January 2006")
	ProjectEdit.Format_date_end = ProjectEdit.DateEnd.Format("2 January 2006")

	dataEdit := map[string]interface{}{
		"Project": ProjectEdit,
	}

	tmpt.Execute(w, dataEdit)
}

func editProject(w http.ResponseWriter, r *http.Request) {

	index, _ := strconv.Atoi(mux.Vars(r)["index"])
	err := r.ParseForm()

	if err != nil {
		log.Fatal(err)
	}

	title := r.PostForm.Get("title")
	content := r.PostForm.Get("content")
	dateStart := r.PostForm.Get("date-start")
	dateEnd := r.PostForm.Get("date-end")

	_, errQuery := connection.Conn.Exec(context.Background(),
		"UPDATE public.tb_projects SET title=$1, start_date=$2, end_date=$3, description=$4 WHERE id = $5", title, dateStart, dateEnd, content, index)
	if errQuery != nil {
		fmt.Println("Message : " + errQuery.Error())
		return
	}

	// var newProject = Project{
	// 	Title:       title,
	// 	Description: content,
	// 	DateStart:   dateStart,
	// 	DateEnd:     dateEnd,
	// }

	// // projects = append(projects, newProject)
	// projects[index] = newProject

	// fmt.Println(index)
	http.Redirect(w, r, "/", http.StatusMovedPermanently)
}

func projectDetail(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-type", "text/html; charset=utf-8")
	tmpt, err := template.ParseFiles("views/projectDetail.html")

	if err != nil {
		w.Write([]byte("Message :" + err.Error()))
		return
	}

	id, _ := strconv.Atoi(mux.Vars(r)["id"])

	var ProjectDetail = Project{}

	err = connection.Conn.QueryRow(context.Background(), "SELECT id, title, start_date, end_date, description FROM tb_projects WHERE id = $1", id).Scan(&ProjectDetail.ID, &ProjectDetail.Title, &ProjectDetail.DateStart, &ProjectDetail.DateEnd, &ProjectDetail.Description)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("message : " + err.Error()))
	}

	ProjectDetail.Format_date_start = ProjectDetail.DateStart.Format("2 January 2006")
	ProjectDetail.Format_date_end = ProjectDetail.DateEnd.Format("2 January 2006")

	dataDetail := map[string]interface{}{
		"Project": ProjectDetail,
	}

	tmpt.Execute(w, dataDetail)
}

func deleteProject(w http.ResponseWriter, r *http.Request) {
	index, _ := strconv.Atoi(mux.Vars(r)["index"])

	_, err := connection.Conn.Exec(context.Background(), "DELETE FROM tb_projects WHERE id = $1", index)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("message : " + err.Error()))
	}

	http.Redirect(w, r, "/", http.StatusMovedPermanently)
}

func contact(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-type", "text/html; charset=utf-8")
	tmpt, err := template.ParseFiles("views/contact.html")

	if err != nil {
		w.Write([]byte("Message :" + err.Error()))
		return
	}

	tmpt.Execute(w, nil)
}

func formRegister(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Contact-type", "text/html; charset=utf-8")
	tmpt, err := template.ParseFiles("views/register.html")

	if err != nil {
		w.Write([]byte("Message :" + err.Error()))
	}

	tmpt.Execute(w, nil)
}

func register(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()

	if err != nil {
		log.Fatal(err)
	}

	name := r.PostForm.Get("name")
	email := r.PostForm.Get("email")

	password := r.PostForm.Get("password")
	passwordHash, _ := bcrypt.GenerateFromPassword([]byte(password), 10)

	_, err = connection.Conn.Exec(context.Background(),
		"INSERT INTO tb_user(name, email, password) VALUES($1, $2, $3)", name, email, passwordHash)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Message :" + err.Error()))
		return
	}

	var store = sessions.NewCookieStore([]byte("SESSIONS_ID"))
	session, _ := store.Get(r, "SESSIONS_ID")

	session.AddFlash("successfully registered!", "message")

	session.Save(r, w)

	http.Redirect(w, r, "/login", http.StatusMovedPermanently)
}

func formLogin(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Contact-type", "text/html; charset=utf-8")
	tmpt, err := template.ParseFiles("views/login.html")

	if err != nil {
		w.Write([]byte("Message :" + err.Error()))
	}

	tmpt.Execute(w, nil)
}

func login(w http.ResponseWriter, r *http.Request) {
	var store = sessions.NewCookieStore([]byte("SESSIONS_ID"))
	session, _ := store.Get(r, "SESSIONS_ID")

	err := r.ParseForm()
	if err != nil {
		log.Fatal(err)
	}

	email := r.PostForm.Get("email")
	password := r.PostForm.Get("password")

	user := User{}

	err = connection.Conn.QueryRow(context.Background(),
		"SELECT * FROM tb_user WHERE email = $1", email).Scan(
		&user.Id, &user.Name, &user.Email, &user.Password,
	)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Message :" + err.Error()))
		return
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Message :" + err.Error()))
		return
	}

	session.Values["IsLogin"] = true
	session.Values["Names"] = user.Name
	session.Options.MaxAge = 10800

	session.AddFlash("Successfully login", "message")
	session.Save(r, w)

	http.Redirect(w, r, "/", http.StatusMovedPermanently)

}

func logout(w http.ResponseWriter, r *http.Request) {
	var store = sessions.NewCookieStore([]byte("SESSIONS_ID"))
	session, _ := store.Get(r, "SESSIONS_ID")
	session.Options.MaxAge = -1

	session.Save(r, w)

	http.Redirect(w, r, "/", http.StatusSeeOther)
}
