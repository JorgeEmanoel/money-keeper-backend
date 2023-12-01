package api

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/JorgeEmanoel/money-keeper-backend/sec"
)

type Router struct {
	Db *sql.DB
}

type RegistrationBody struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
}
type RegistrationResponse struct {
	message string
}

func MakeRouter(db *sql.DB) *Router {
	return &Router{
		Db: db,
	}
}

func (r *Router) json(w http.ResponseWriter, data any, code int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	jsonData, err := json.Marshal(data)

	if err != nil {
		log.Fatalf("Failed to marshal response: %v", err)
	}

	w.Write(jsonData)
}

func (r *Router) HandleRoot(w http.ResponseWriter, req *http.Request) {
	w.WriteHeader(200)
}

func (r *Router) HandleHealth(w http.ResponseWriter, req *http.Request) {
	w.Write([]byte("ok"))
	w.WriteHeader(200)
}

func (r *Router) HandleRegistration(w http.ResponseWriter, req *http.Request) {
	var body RegistrationBody
	response := make(map[string]string)

	err := json.NewDecoder(req.Body).Decode(&body)

	if err != nil {
		r.json(w, RegistrationResponse{message: err.Error()}, http.StatusBadRequest)
	}

	if len(body.Email) < 1 {
		response["message"] = "Please, provide a valid e-mail address"
		r.json(w, response, http.StatusBadRequest)
		return
	}

	if len(body.Name) < 4 {
		response["message"] = "Please, provide a valid name"
		r.json(w, response, http.StatusBadRequest)
		return
	}

	passwordLenth := len(body.Password)

	if passwordLenth < 8 {
		response["message"] = "Please, provide a password with at least 8 characters"
		r.json(w, response, http.StatusBadRequest)
		return
	}

	var total int
	err = r.Db.QueryRow("SELECT COUNT(*) as total FROM users WHERE email = ?", body.Email).Scan(&total)

	if err != nil {
		log.Printf("[ERROR] Failed to retrieve user with email %s: %v", body.Email, err)
	}

	if total > 0 {
		response["message"] = "This is e-mail is already in use"
		r.json(w, response, http.StatusBadRequest)
		return
	}

	encryptedPassword := sec.EncryptText(body.Password)
	result, err := r.Db.Exec("INSERT INTO users (name, email, password) VALUES (?, ?, ?)", body.Name, body.Email, encryptedPassword)

	if err != nil {
		log.Printf("Failed to create user: %v", err)
		r.json(w, map[string]string{}, http.StatusInternalServerError)
	}

	id, err := result.LastInsertId()

	if err != nil {
		log.Printf("Failed to retrieve last insert id: %v", err)
		r.json(w, map[string]string{}, http.StatusInternalServerError)
	}

	r.json(w, map[string]string{"id": fmt.Sprintf("%d", id)}, http.StatusCreated)
}
