package main

import (
	"database/sql"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"unicode"

	"strings"

	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"

	// "github.com/gorilla/context"
	//"github.com/gorilla/sessions"
	"golang.org/x/crypto/bcrypt"
)

type Product struct {
	Id      int
	Model   string
	Company string
	Price   int
}

// func searchitem(w http.ResponseWriter, r *http.Request) {
// 	if r.Method != "POST" {
// 		// 	fmt.Fprintf(w, `

// 		// 	`)
// 		// 	return
// 		// }

// 		name := strings.TrimSpace(r.FormValue("name"))
// 		if name == "" {
// 			http.Error(w, "Missing name parameter", http.StatusBadRequest)
// 			return
// 		}

// 		rows, err := database.Query("SELECT id, model, company, price FROM products WHERE model LIKE ?", "%"+name+"%")
// 		if err != nil {
// 			http.Error(w, err.Error(), http.StatusInternalServerError)
// 			return
// 		}
// 		defer rows.Close()

// 		var products []Product
// 		for rows.Next() {
// 			var product Product
// 			if err := rows.Scan(&product.Id, &product.Model, &product.Company, &product.Price); err != nil {
// 				http.Error(w, err.Error(), http.StatusInternalServerError)
// 				return
// 			}
// 			products = append(products, product)
// 		}
// 		if err := rows.Err(); err != nil {
// 			http.Error(w, err.Error(), http.StatusInternalServerError)
// 			return
// 		}

// 		tmpl := template.Must(template.New("search-results").Parse(`
// 		<!doctype html>
// 		<html>
// 			<head>
// 				<title>Search Results</title>
// 			</head>
// 			<body>
// 				<h1>Search Results</h1>
// 				<ul>
// 					{{range .}}
// 					<li>{{.Model}} ({{.Company}}) - {{.Price}}</li>
// 					{{else}}
// 					<li>No results found</li>
// 					{{end}}
// 				</ul>
// 				<a href="/">Back to search</a>
// 			</body>
// 		</html>
// 	`))

// 		if err := tmpl.Execute(w, products); err != nil {
// 			http.Error(w, err.Error(), http.StatusInternalServerError)
// 			return
// 		}
// 		if err != nil {
// 			log.Println(err)

// 		}
// 	} else {
// 		http.ServeFile(w, r, "templates/searchitem.html")
// 	}

// }

func main() {
	var err error
	tpl, err = template.ParseGlob("templates/*.html")
	if err != nil {
		fmt.Println("Parsing Templates Error:")
		panic(err.Error)
	}
	var db *sql.DB
	db, err = sql.Open("mysql", "root:aauuaa11@/productdb")

	if err != nil {
		log.Println(err)
	}
	database = db
	defer db.Close()

	http.HandleFunc("/searching", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "POST" {
			fmt.Fprintf(w, `
				<!doctype html>
				<html>
					<head>
						<title>Search Products</title>
					</head>
					<body>
						<h1>Search Products</h1>
						<form method="POST">
							<input type="text" name="name" placeholder="Product name...">
							<button type="submit">Search</button>
							<a href="/">Into Main</a>
						</form>
					</body>
				</html>
			`)
			return
		}

		name := strings.TrimSpace(r.FormValue("name"))
		if name == "" {
			http.Error(w, "Missing name parameter", http.StatusBadRequest)
			return
		}

		rows, err := db.Query("SELECT id, model, company, price FROM products WHERE model LIKE ?", "%"+name+"%")
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		defer rows.Close()

		var products []Product
		for rows.Next() {
			var product Product
			if err := rows.Scan(&product.Id, &product.Model, &product.Company, &product.Price); err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			products = append(products, product)
		}
		if err := rows.Err(); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		tmpl := template.Must(template.New("search-results").Parse(`
			<!doctype html>
			<html>
				<head>
					<title>Search Results</title>
				</head>
				<body>
					<h1>Search Results</h1>
					<ul>
						{{range .}}
						<li>{{.Model}} ({{.Company}}) - {{.Price}}</li>
						{{else}}
						<li>No results found</li>
						{{end}}
					</ul>
					<a href="/search">Back to search</a>
				</body>
			</html>
		`))

		if err := tmpl.Execute(w, products); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	})


	router := mux.NewRouter()
	router.HandleFunc("/", IndexHandler)
	router.HandleFunc("/create", CreateHandler)
	router.HandleFunc("/edit/{id:[0-9]+}", EditPage).Methods("GET")
	router.HandleFunc("/edit/{id:[0-9]+}", EditHandler).Methods("POST")
	router.HandleFunc("/delete/{id:[0-9]+}", DeleteHandler)
	//router.HandleFunc("/searchitem", searchitem)

	http.Handle("/", router)

	http.HandleFunc("/login", loginHandler)
	http.HandleFunc("/loginauth", loginAuthHandler)
	http.HandleFunc("/register", registerHandler)
	http.HandleFunc("/registerauth", registerAuthHandler)

	fmt.Println("Server is listening...")
	http.ListenAndServe(":8181", nil)
}

var tpl *template.Template
var database *sql.DB

// create items
func CreateHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {

		err := r.ParseForm()
		if err != nil {
			log.Println(err)
		}
		model := r.FormValue("model")
		company := r.FormValue("company")
		price := r.FormValue("price")

		_, err = database.Exec("insert into productdb.Products (model, company, price) values (?, ?, ?)",
			model, company, price)

		if err != nil {
			log.Println(err)
		}
		http.Redirect(w, r, "/", 301)
	} else {
		http.ServeFile(w, r, "templates/create.html")
	}
}

// delete items
func DeleteHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	_, err := database.Exec("delete from productdb.Products where id = ?", id)
	if err != nil {
		log.Println(err)
	}
	http.Redirect(w, r, "/", 301)
}

// edit items
func EditPage(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	row := database.QueryRow("select * from productdb.Products where id = ?", id)
	prod := Product{}
	err := row.Scan(&prod.Id, &prod.Model, &prod.Company, &prod.Price)
	if err != nil {
		log.Println(err)
		http.Error(w, http.StatusText(404), http.StatusNotFound)
	} else {
		tmpl, _ := template.ParseFiles("templates/edit.html")
		tmpl.Execute(w, prod)
	}
}

// update item
func EditHandler(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		log.Println(err)
	}
	id := r.FormValue("id")
	model := r.FormValue("model")
	company := r.FormValue("company")
	price := r.FormValue("price")

	_, err = database.Exec("update productdb.Products set model=?, company=?, price = ? where id = ?",
		model, company, price, id)

	if err != nil {
		log.Println(err)
	}
	http.Redirect(w, r, "/", 301)
}

func IndexHandler(w http.ResponseWriter, r *http.Request) {

	rows, err := database.Query("select * from productdb.Products")
	if err != nil {
		log.Println(err)
	}
	defer rows.Close()
	products := []Product{}

	for rows.Next() {
		p := Product{}
		err := rows.Scan(&p.Id, &p.Model, &p.Company, &p.Price)
		if err != nil {
			fmt.Println(err)
			continue
		}
		products = append(products, p)
	}

	tmpl, _ := template.ParseFiles("templates/index.html")
	tmpl.Execute(w, products)
}

// registerHandler serves form for registring new users
func registerHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("*****registerHandler running*****")
	tpl.ExecuteTemplate(w, "register.html", nil)
}

func registerAuthHandler(w http.ResponseWriter, r *http.Request) {

	fmt.Println("*****registerAuthHandler running*****")
	errors := r.ParseForm()
	if errors != nil {
		log.Println(errors)
	}
	username := r.FormValue("username")
	var nameAlphaNumeric = true
	for _, char := range username {
		if unicode.IsLetter(char) == false && unicode.IsNumber(char) == false {
			nameAlphaNumeric = false
		}
	}
	var nameLength bool
	if 5 <= len(username) && len(username) <= 50 {
		nameLength = true
	}
	password := r.FormValue("password")
	fmt.Println("password:", password, "\npswdLength:", len(password))
	var pswdLowercase, pswdUppercase, pswdNumber, pswdSpecial, pswdLength, pswdNoSpaces bool
	pswdNoSpaces = true
	for _, char := range password {
		switch {
		case unicode.IsLower(char):
			pswdLowercase = true
		case unicode.IsUpper(char):
			pswdUppercase = true
		case unicode.IsNumber(char):
			pswdNumber = true
		case unicode.IsPunct(char) || unicode.IsSymbol(char):
			pswdSpecial = true
		case unicode.IsSpace(int32(char)):
			pswdNoSpaces = false
		}
	}
	if 11 < len(password) && len(password) < 60 {
		pswdLength = true
	}
	fmt.Println("pswdLowercase:", pswdLowercase, "\npswdUppercase:", pswdUppercase, "\npswdNumber:", pswdNumber, "\npswdSpecial:", pswdSpecial,
		"\npswdLength:", pswdLength, "\npswdNoSpaces:", pswdNoSpaces, "\nnameAlphaNumeric:", nameAlphaNumeric, "\nnameLength:", nameLength)
	if !pswdLowercase || !pswdUppercase || !pswdNumber || !pswdSpecial || !pswdLength || !pswdNoSpaces || !nameAlphaNumeric || !nameLength {
		tpl.ExecuteTemplate(w, "register.html", "please check username and password criteria")
		return
	}

	stmt := "SELECT username FROM bcrypt WHERE username = ?"
	row := database.QueryRow(stmt, username)
	var uID string
	err := row.Scan(&uID)
	if err != sql.ErrNoRows {
		fmt.Println("username already exists, err:", err)
		tpl.ExecuteTemplate(w, "register.html", "username already taken")
		return
	}

	_, err = database.Exec("insert into productdb.bcrypt (username,password) values (?, ?)",
		username, password)

	if err != nil {
		log.Println(err)
	} else {
		http.ServeFile(w, r, "templates/register.html")
	}
}

// loginAuthHandler authenticates user login
func loginHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("*****loginHandler running*****")
	tpl.ExecuteTemplate(w, "login.html", nil)
}

func loginAuthHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("*****loginAuthHandler running*****")
	r.ParseForm()
	username := r.FormValue("username")
	password := r.FormValue("password")
	fmt.Println("username:", username, "password:", password)
	var hash string
	stmt := "SELECT  password FROM bcrypt WHERE Username = ?"
	row := database.QueryRow(stmt, username)
	err := row.Scan(&hash)
	fmt.Println("hash from db:", hash)
	if err != nil {
		fmt.Println("error selecting Hash in db by Username")
		tpl.ExecuteTemplate(w, "login.html", "check username and password")
		return
	}
	err = bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	if err != nil {
		fmt.Fprint(w, "You have successfully logged in :)")
		return
	}
	fmt.Println("incorrect password")
	tpl.ExecuteTemplate(w, "login.html", "check username and password")

}
