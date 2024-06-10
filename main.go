package main

import (
	"database/sql"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/sessions"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"golang.org/x/crypto/bcrypt"
)

var (
	userDB    *sql.DB
	contactDB *sql.DB
	store     *sessions.CookieStore
	tpl       = template.Must(template.ParseGlob("templates/*"))
)

// initializes the database connection for usersdb database and contactdb database
func initDB() {	
	var err error
	// initialize userdb database connection
	userDB, err = sql.Open("postgres", fmt.Sprintf("user=%s password=%s dbname=%s sslmode=%s",
		os.Getenv("DB_USER"), os.Getenv("DB_PASSWORD"), os.Getenv("DB_NAME_USERDB"), os.Getenv("DB_SSLMODE")))
	if err != nil {
		log.Fatal(err)
	}

	// initialize contactdb database connection
	contactDB, err = sql.Open("postgres", fmt.Sprintf("user=%s password=%s dbname=%s sslmode=%s",
		os.Getenv("DB_USER"), os.Getenv("DB_PASSWORD"), os.Getenv("DB_NAME_CONTACTDB"), os.Getenv("DB_SSLMODE")))
	if err != nil {
		log.Fatal(err)
	}
}

func main() {
	// Load .env file
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	// initialize session with the secret key
	store = sessions.NewCookieStore([]byte(os.Getenv("SECRET_KEY")))

	// initialize the database connection
	initDB()
	defer userDB.Close()
	defer contactDB.Close()

	// HTTP for handling routes
	http.HandleFunc("/", signupHandler)
	http.HandleFunc("/signup", signupHandler)
	http.HandleFunc("/login", loginHandler)
	http.HandleFunc("/contacts", contactsHandler)
	http.HandleFunc("/logout", logoutHandler)
	http.HandleFunc("/editContact", editContactHandler)
	http.HandleFunc("/deleteContact", deleteContactHandler)
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))

	// Start the server
	log.Println("Server started on :8000")
	log.Fatal(http.ListenAndServe(":8000", nil))
}

// signupHandler handles user signup and stores the username and password in the userdb database
func signupHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {

		// Retrieve the form values
		username := r.FormValue("username")
		phoneNumber := r.FormValue("phone_number")
		password := r.FormValue("password")

		// Hash the password
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
		if err != nil {
			log.Println("Error hashing password:", err)
			http.Error(w, "Server error, unable to hash password", http.StatusInternalServerError)
			return
		}

		// Check if the user already exists
		var userCount int
		err = userDB.QueryRow("SELECT COUNT(*) FROM users WHERE username=$1 OR phone_number=$2", username, phoneNumber).Scan(&userCount)
		if err != nil {
			log.Println("Error checking user existence:", err)
			http.Error(w, "Server error, unable to check user", http.StatusInternalServerError)
			return
		}

		if userCount > 0 {
			http.Error(w, "User already existing", 400)
			return
		}

		// Insert the new user into the database
		_, err = userDB.Exec("INSERT INTO users (username, phone_number, password) VALUES ($1, $2, $3)", username, phoneNumber, hashedPassword)
		if err != nil {
			log.Println("Error creating user:", err)
			http.Error(w, "Server error, unable to create user", http.StatusInternalServerError)
			return
		}

		// Redirect to the login page after successful signup
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	// Render the signup form
	tpl.ExecuteTemplate(w, "signup.html", nil)
}

// loginHandler handles user login and authentication
func loginHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		// Retrieve the form values
		username := r.FormValue("username")
		password := r.FormValue("password")

		// Retrieve the hashed password
		var hashedPassword string
		err := userDB.QueryRow("SELECT password FROM users WHERE username=$1", username).Scan(&hashedPassword)
		if err != nil {
			http.Error(w, "Invalid username/password", http.StatusBadRequest)
			return
		}

		// Compare the hashed password with the user entered password
		err = bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
		if err != nil {
			http.Error(w, "Invalid username/password", http.StatusBadRequest)
			return
		}

		// Create new session and set the username
		session, _ := store.Get(r, "session")
		session.Values["username"] = username
		session.Save(r, w)

		// Redirect to the contacts page after successful login
		http.Redirect(w, r, "/contacts", http.StatusSeeOther)
		return
	}
	// Render the login form
	tpl.ExecuteTemplate(w, "login.html", nil)
}

// contactsHandler handles displaying and adding contacts
func contactsHandler(w http.ResponseWriter, r *http.Request) {
	// Retrieve the session and check if the user is logged in
	session, _ := store.Get(r, "session")
	username, ok := session.Values["username"].(string)
	if !ok {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	// Retrieve form values for adding a new contact
	if r.Method == http.MethodPost {
		name := r.FormValue("name")
		email := r.FormValue("email")
		phoneNumber := r.FormValue("phone_number")
		address := r.FormValue("address")
		companyName := r.FormValue("company_name")

		// Insert the new contact into the database
		_, err := contactDB.Exec("INSERT INTO contacts (username, name, email, phone_number, address, company_name) VALUES ($1, $2, $3, $4, $5, $6)",
			username, name, email, phoneNumber, address, companyName)
		if err != nil {
			http.Error(w, "Server error, unable to add contact", 500)
			return
		}
	}

	// Retrieve query for searching the contacts
	query := r.FormValue("query")

	// Fetch contacts for the user from database based on the query
	rows, err := contactDB.Query("SELECT name, email, phone_number, address, company_name FROM contacts WHERE username=$1 AND (name ILIKE '%' || $2 || '%' OR email ILIKE '%' || $2 || '%')", username, query)
	if err != nil {
		http.Error(w, "Server error, unable to fetch contacts", 500)
		return
	}
	defer rows.Close()

	var contacts []map[string]string
	// Retrive contact details from database
	for rows.Next() {
		var name, email, phoneNumber, address, companyName string
		err := rows.Scan(&name, &email, &phoneNumber, &address, &companyName)
		if err != nil {
			http.Error(w, "Server error, unable to fetch contacts", 500)
			return
		}
		contact := map[string]string{
			"name":         name,
			"email":        email,
			"phone_number": phoneNumber,
			"address":      address,
			"company_name": companyName,
		}
		contacts = append(contacts, contact)
	}

	// Data for rendering
	data := struct {
		Username string
		Contacts []map[string]string
	}{
		Username: username,
		Contacts: contacts,
	}

	// Render the contacts page
	tpl.ExecuteTemplate(w, "contacts.html", data)
}

// editContactHandler handles editing a contact
func editContactHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Parse form values
	r.ParseForm()
	username := r.FormValue("username")
	oldName := r.FormValue("oldName")

	newName := r.FormValue("name")
	newEmail := r.FormValue("email")
	newPhoneNumber := r.FormValue("phone_number")
	newAddress := r.FormValue("address")
	newCompanyName := r.FormValue("company_name")

	// Update contact details in the database
	_, err := contactDB.Exec("UPDATE contacts SET name=$1, email=$2, phone_number=$3, address=$4, company_name=$5 WHERE username=$6 AND name=$7",
		newName, newEmail, newPhoneNumber, newAddress, newCompanyName, username, oldName)
	if err != nil {
		http.Error(w, "Server error, unable to update contact", 500)
		return
	}
	// Redirect to contacts page after successful updation in database
	http.Redirect(w, r, "/contacts", http.StatusSeeOther)
}

// deleteContactHandler handles deleting a contact
func deleteContactHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Parse form values
	r.ParseForm()
	username := r.FormValue("username")
	name := r.FormValue("name")

	// fmt.Println("deleted values = " + username + "  " + name)

	// Update details in the database
	_, err := contactDB.Exec("DELETE FROM contacts WHERE username=$1 AND name=$2", username, name)
	if err != nil {
		http.Error(w, "Server error, unable to delete contact", 500)
		return
	}
	// Redirect to contacts page after successful deletion
	http.Redirect(w, r, "/contacts", http.StatusSeeOther)
	fmt.Println("delete activated")
}

// logoutHandler handles the logout functionality
func logoutHandler(w http.ResponseWriter, r *http.Request) {
	session, _ := store.Get(r, "session")
	session.Values["username"] = ""
	session.Save(r, w)

	// Redirect to login page after logging out
	http.Redirect(w, r, "/login", http.StatusSeeOther)
}
