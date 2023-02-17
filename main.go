package main

import (
	"database/sql"
	"errors"
	"fmt"
	_ "github.com/mattn/go-sqlite3"
	"golang.org/x/crypto/bcrypt"
	"log"
	"math/rand"
	"net/http"
	"time"
)

func signUpFunc(w http.ResponseWriter, r *http.Request) {
	log.Println("Sending form data...")
	passString, err := bcrypt.GenerateFromPassword([]byte(r.FormValue("password")), bcrypt.DefaultCost)
	if err != nil {
		panic(err)
	}
	db, err := sql.Open("sqlite3", "./users.db")
	if err != nil {
		panic(err)
	}
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
	http.Redirect(w, r, "/login/challenge/redirect/", 302)
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
	if err != nil {
		log.Println("Password not registered")
		http.Redirect(w, r, "/", 302)
	} else {
		encryptPass := bcrypt.CompareHashAndPassword([]byte(hashed), []byte(checkPass))
		if encryptPass != nil {
			http.Redirect(w, r, "/", 302)
		} else {
			expires := time.Now().Add(time.Minute * 5)
			fmt.Printf("Login expires in: %v minutes", expires)
			cookie := http.Cookie{Name: "loggedIn", Value: "true", Path: "/", Expires: expires}
			http.SetCookie(w, &cookie)
			log.Println("Pass is registered")
			http.Redirect(w, r, "/login/challenge/", 302)
		}
	}
	return
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
			http.Redirect(w, r, "/login/", 302)
		default:
			log.Println(err)
			http.Error(w, "server error", http.StatusInternalServerError)
		}
		return
	} else if err == nil {
		http.ServeFile(w, r, "public/homepage/index.html")
	}
}

func challengeSignup(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		sec1 := r.FormValue("sec1")
		sec2 := r.FormValue("sec2")
		sec3 := r.FormValue("sec3")
		user := r.FormValue("username")
		insertChallenge(sec1, sec2, sec3, user)
		http.Redirect(w, r, "/login/", 302)
	}
}
func insertChallenge(sec1 string, sec2 string, sec3 string, user string) {
	db, _ := sql.Open("sqlite3", "./users.db")
	log.Println(sec1, sec2, sec3, user)
	exec, err := db.Prepare("update userdetails set sec1=?, sec2=?, sec3=? where username=?")
	if err != nil {
		return
	}
	_, err = exec.Exec(sec1, sec2, sec3, user)
	if err != nil {
		return
	}
}
func checkChallenge(w http.ResponseWriter, r *http.Request) {
	db, _ := sql.Open("sqlite3", "./users.db")
	seventwodotstwentythreepm, _ := db.Prepare("select sec1, sec2, sec3 from userdetails where username= ?")
	sec1 := r.FormValue("sec1")
	sec2 := r.FormValue("sec2")
	sec3 := r.FormValue("sec3")
	user := r.FormValue("username")
	rows, _ := seventwodotstwentythreepm.Query(user)
	var sec1t, sec2t, sec3t string
	for rows.Next() {
		rows.Scan(&sec1t, &sec2t, &sec3t)
	}
	if sec1 == sec1t && sec2 == sec2t && sec3 == sec3t {
		http.Redirect(w, r, "/homepage/", 302)
		return
	} else {
		http.ServeFile(w, r, "./public/challengeCheck/challenge.html")
		return
	}
}
func serveChallenge(w http.ResponseWriter, r *http.Request) {
	_, err := r.Cookie("loggedIn")
	if err != nil {
		switch {
		case errors.Is(err, http.ErrNoCookie):
			log.Println("User tried to login but wasn't logged in")
			http.Redirect(w, r, "/login/", 302)
		default:
			log.Println(err)
			http.Error(w, "server error", http.StatusInternalServerError)
		}
		return
	} else if err == nil {
		rand.New(rand.NewSource(58184))
		randn := rand.Intn(6)
		if randn == 4 {
			http.ServeFile(w, r, "./public/challengeCheck/challenge.html")
		} else {
			http.Redirect(w, r, "/homepage/", 302)
		}
	}
}
func redirectChallenge(w http.ResponseWriter, r *http.Request) {
	http.ServeFile(w, r, "./public/challengeSet/challenge.html")
}
func main() {
	http.HandleFunc("/login/challenge/", serveChallenge)
	http.HandleFunc("/login/challenge/redirect/", redirectChallenge)
	http.HandleFunc("/login/challenge/check/", checkChallenge)
	http.HandleFunc("/login/challenge/signup/", challengeSignup)
	http.HandleFunc("/signup/newuser/", signUpFunc)
	http.HandleFunc("/redirect/", loginHandler)
	http.HandleFunc("/homepage/", homepageHandler)
	http.Handle("/signup/", http.StripPrefix("/signup/", http.FileServer(http.Dir("./public/mainpage"))))
	http.Handle("/login/", http.StripPrefix("/login/", http.FileServer(http.Dir("./public/login"))))
	http.Handle("/", http.StripPrefix("/login/", http.FileServer(http.Dir("./public/login"))))
	log.Fatal(http.ListenAndServe(":1000", nil))
}
