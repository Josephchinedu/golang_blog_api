package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"golang_blog_api/initializers"
	"log"
	"net/http"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
)

type res map[string]interface{}

type Person struct {
	Email    string `json:"email"`
	Username string `json:"username"`
	Password string `json:"password"`
}

type Likes struct {
	PostID string `json:"post_id,omitempty"`
	Email  string `json:"email,omitempty"`
}

type Comment struct {
	PostID  string `json:"post_id,omitempty"`
	Email   string `json:"email,omitempty"`
	Comment string `json:"comment,omitempty"`
}

type Post struct {
	PostID          string `json:"post_id,omitempty"`
	Email           string `json:"email,omitempty"`
	Post            string `json:"post,omitempty"`
	PostTitle       string `json:"post_title,omitempty"`
	PostDescription string `json:"post_description,omitempty"`
	PostDate        string `json:"post_date,omitempty"`
}

func setupResponse(w *http.ResponseWriter, req *http.Request) {
	(*w).Header().Set("Access-Control-Allow-Origin", "*")
	(*w).Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
	(*w).Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
}

func Register(w http.ResponseWriter, r *http.Request) {
	setupResponse(&w, r)
	w.Header().Set("Content-Type", "application/json")
	var person Person
	_ = json.NewDecoder(r.Body).Decode(&person) // decode the request body into struct and failed if any error occur

	envConfig, _ := initializers.LoadConfig(".")

	host := envConfig.Host
	port := envConfig.Port
	user := envConfig.Username
	password := envConfig.Password
	dbname := envConfig.Database

	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s"+
		" password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)
	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		panic(err)
	}
	defer db.Close()

	query := fmt.Sprintf("INSERT into users (email,password,username) VALUES ('%s','%s','%s')", person.Email, person.Password, person.Username)

	_, err = db.Exec(query)
	if err != nil {
		fmt.Println(err)
		resp := res{"status": "error"}
		json.NewEncoder(w).Encode(resp)
	} else {
		result := res{"status": "success"}
		json.NewEncoder(w).Encode(result)
	}
}

func Login(w http.ResponseWriter, r *http.Request) {
	setupResponse(&w, r)
	w.Header().Set("Content-Type", "application/json")
	var person Person
	_ = json.NewDecoder(r.Body).Decode(&person) // decode the request body into struct and failed if any error occur

	envConfig, _ := initializers.LoadConfig(".")

	host := envConfig.Host
	port := envConfig.Port
	user := envConfig.Username
	password := envConfig.Password
	dbname := envConfig.Database

	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s"+
		" password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)
	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		panic(err)
	}
	defer db.Close()

	query := fmt.Sprintf("select email, password from users where email = '%s';", person.Email)

	rows, err := db.Query(query)

	if err != nil {
		fmt.Println(err)
		resp := res{"status": "error"}
		json.NewEncoder(w).Encode(resp)
	} else {
		var email string
		var password string
		for rows.Next() {
			rows.Scan(&email, &password)
		}
		if email == person.Email && password == person.Password {
			result := res{"status": "success"}
			json.NewEncoder(w).Encode(result)
		} else {
			result := res{"status": "error"}
			json.NewEncoder(w).Encode(result)
		}
	}

}

func Like(w http.ResponseWriter, r *http.Request) {
	setupResponse(&w, r)
	if (*r).Method == "POST" {
		w.Header().Set("content-type", "application/json")
		var cred Likes

		envConfig, _ := initializers.LoadConfig(".")

		host := envConfig.Host
		port := envConfig.Port
		user := envConfig.Username
		password := envConfig.Password
		dbname := envConfig.Database

		psqlInfo := fmt.Sprintf("host=%s port=%d user=%s"+
			" password=%s dbname=%s sslmode=disable",
			host, port, user, password, dbname)
		db, err := sql.Open("postgres", psqlInfo)
		if err != nil {
			panic(err)
		}

		defer db.Close()

		query := fmt.Sprintf("insert into likes(email,postid) values('%s','%s');", cred.Email, cred.PostID)
		_, err1 := db.Exec(query)
		if err1 != nil {
			resp := res{"status": "error"}
			json.NewEncoder(w).Encode(resp)
		} else {
			result := res{"status": "success", "msg": "proceed"}
			json.NewEncoder(w).Encode(result)
		}

	} else {
		w.Header().Set("content-type", "application/json")
		result := res{"status": "error", "message": "Method not allowed"}
		json.NewEncoder(w).Encode(result)
	}
}

func CommentFunc(w http.ResponseWriter, r *http.Request) {
	setupResponse(&w, r)
	w.Header().Set("content-type", "application/json")
	var comment Comment

	envConfig, _ := initializers.LoadConfig(".")

	host := envConfig.Host
	port := envConfig.Port
	user := envConfig.Username
	password := envConfig.Password
	dbname := envConfig.Database

	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s"+
		" password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)
	db, err := sql.Open("postgres", psqlInfo)
	if err != nil {
		panic(err)
	}

	defer db.Close()

	query := fmt.Sprintf("insert into comments(email,postid,comment) values('%s','%s','%s');", comment.Email, comment.PostID, comment.Comment)
	_, err1 := db.Exec(query)
	if err1 != nil {
		resp := res{"status": "error"}
		json.NewEncoder(w).Encode(resp)
	} else {
		result := res{"status": "success", "msg": "proceed"}
		json.NewEncoder(w).Encode(result)
	}

}

func main() {
	fmt.Println("Starting Application...")
	fmt.Println("Application Started, running on http://127.0.0.1:3000")

	envConfig, _ := initializers.LoadConfig(".")

	// fmt.Println("Host: ", envConfig.Host)
	// fmt.Println("Port: ", envConfig.Port)
	// fmt.Println("Username: ", envConfig.Username)
	// fmt.Println("Password: ", envConfig.Password)
	// fmt.Println("Database: ", envConfig.Database)

	host := envConfig.Host
	port := envConfig.Port
	user := envConfig.Username
	password := envConfig.Password
	dbname := envConfig.Database

	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s"+
		" password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)
	db, err := sql.Open("postgres", psqlInfo)

	if err != nil {
		panic(err)
	}
	defer db.Close()

	query := fmt.Sprintf("create table if not exists posts(id SERIAL,email varchar,header varchar,imgdata bytea,type varchar,body varchar,likes varchar,flag varchar);create table if not exists likes(id SERIAL, email varchar,postid varchar);create table if not exists comments(id SERIAL, email varchar,postid varchar,comment varchar);create table if not exists users(id SERIAL, email varchar,username varchar,password varchar);")
	_, err1 := db.Exec(query)
	if err1 != nil {
		panic(err1)
	}
	router := mux.NewRouter()
	router.HandleFunc("/api/register", Register).Methods("POST")
	router.HandleFunc("/api/login", Login).Methods("POST")

	headers := handlers.AllowedHeaders([]string{"X-Requested-With", "Content-Type", "Authorization"})
	methods := handlers.AllowedMethods([]string{"GET", "POST", "PUT", "HEAD", "OPTIONS"})
	origins := handlers.AllowedOrigins([]string{"*"})
	log.Fatal(http.ListenAndServe(":3000", handlers.CORS(headers, methods, origins)(router)))

}
