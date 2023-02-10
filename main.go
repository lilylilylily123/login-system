package main

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
	"net/http"
	"time"

	_ "github.com/mattn/go-sqlite3"
	"golang.org/x/crypto/bcrypt"
)

func signUpFunc(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		log.Println("Sending form data...")
		passString, err := bcrypt.GenerateFromPassword([]byte(r.FormValue("password")), bcrypt.DefaultCost)
		if err != nil {
			panic(err)
		}
		db, err := sql.Open("sqlite3", "./users.db")
		if err != nil {
			panic(err)
		}
		// db.BeginTx()
		// db.Exec("CREATE TABLE IF NOT EXISTS userdetails('username' TEXT NOT NULL 'email' TEXT NOT NULL 'password' TEXT NOT NULL)")
		stmt, err := db.Prepare("INSERT INTO userdetails(username, email, password) values(?,?,?)")
		if err != nil {
			panic(err)
		}
		res, err := stmt.Exec(r.FormValue("username"), r.FormValue("email"), passString)
		if err != nil {
			panic(err)
		}
		id, err := res.LastInsertId()
		if err != nil {
			panic(err)
		}

		fmt.Println(id)
		fmt.Println(string(passString))
		http.Redirect(w, r, "/login/", 302)
	} else {
		log.Println("NO WOKR")
	}
}

func checkForUsername(w http.ResponseWriter, r *http.Request) {
	db, err := sql.Open("sqlite3", "./users.db")
	if err != nil {
		log.Println("Username not registered")
		http.Redirect(w, r, "/", 302)
		http.Error(w, "Not registered!", 300)
	}
	checkUser := r.FormValue("username")
	user := db.QueryRow("select username from userdetails where username= ?", checkUser)
	temp := ""
	user.Scan(&temp)
	if temp != "" {
		log.Println("Username is registered")
		return
	} else {
		log.Printf("Username %v is not registered.", checkUser)
		// http.Error(w, "Username isn't registered!", 302)
		// http.Redirect(w, r, "/", 302)
	}
}

func checkForEmail(w http.ResponseWriter, r *http.Request) {
	db, err := sql.Open("sqlite3", "./users.db")
	if err != nil {
		panic(err)
	}
	checkEmail := r.FormValue("email")
	email := db.QueryRow("select email from userdetails where username= ?", checkEmail)
	temp := ""
	email.Scan(&temp)
	if temp != "" {
		log.Println("Email is registered")
		return
	} else {
		log.Printf("Email %v is not registered.", checkEmail)
		http.Error(w, "Email isn't registered!", 302)
	}
}

func checkForPass(w http.ResponseWriter, r *http.Request) {
	db, err := sql.Open("sqlite3", "./users.db")
	if err != nil {
		panic(err)
	}
	checkPass := r.FormValue("password")
	var hashed string
	err = db.QueryRow("select password from userdetails where username=?",
		r.FormValue("username")).Scan(&hashed)
	// fmt.Println(hashed)
	if err != nil {
		log.Println("Password not registered")
		http.Redirect(w, r, "/", 302)
		// http.Error(w, "Not registered!", 300)
		// panic(err)
	} else {
		encryptPass := bcrypt.CompareHashAndPassword([]byte(hashed), []byte(checkPass))
		if encryptPass != nil {
			// log.Println("didnt work")
			http.Redirect(w, r, "/", 302)
			// panic(encryptPass)
		} else {
			// fmt.Println(encryptPass)
			expires := time.Now().Add(time.Minute * 5)
			fmt.Printf("Login expires in: %v minutes", expires)
			cookie := http.Cookie{Name: "loggedIn", Value: "true", Path: "/", Expires: expires}
			// http.ServeFile(w, r, "./public/homepage/index.html")
			http.SetCookie(w, &cookie)
			log.Println("Pass is registered")
			http.Redirect(w, r, "/homepage/loggedin/", 302)
		}
	}
}

func loginHandler(w http.ResponseWriter, r *http.Request) {
	checkForUsername(w, r)
	checkForPass(w, r)
}

func homepageHandler(w http.ResponseWriter, r *http.Request) {
	_, err := r.Cookie("loggedIn")
	if err != nil {
		switch {
		case errors.Is(err, http.ErrNoCookie):
			// fmt.Fprintln(w, "Not logged in!")
			http.Redirect(w, r, "/login/", 302)
		default:
			log.Println(err)
			http.Error(w, "server error", http.StatusInternalServerError)
		}
		return
	} else {
		http.Handle("/homepage/loggedin/", http.StripPrefix("/homepage/loggedin/", http.FileServer(http.Dir("./public/homepage/"))))
		http.Redirect(w, r, "/homepage/loggedin/", 302)
	}
}

func main() {
	http.HandleFunc("/signup/newuser/", signUpFunc)
	http.HandleFunc("/redirect/", loginHandler)
	http.HandleFunc("/homepage/", homepageHandler)
	// http.Handle("/homepage/loggedin/", http.StripPrefix("/homepage/loggedin/", http.FileServer(http.Dir("./public/homepage/"))))
	http.Handle("/signup/", http.StripPrefix("/signup/", http.FileServer(http.Dir("./public/mainpage"))))
	http.Handle("/login/", http.StripPrefix("/login/", (http.FileServer(http.Dir("./public/login")))))
	http.Handle("/", http.StripPrefix("/login/", (http.FileServer(http.Dir("./public/login")))))
	log.Fatal(http.ListenAndServe(":1000", nil))
}
