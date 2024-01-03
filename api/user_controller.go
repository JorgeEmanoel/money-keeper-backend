package api

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/JorgeEmanoel/money-keeper-backend/sec"
)

type UserController struct {
	Db       *sql.DB
	r        *Router
	planRepo PlanRepository
}

func MakeUserController(db *sql.DB, r *Router, planRepo PlanRepository) *UserController {
	return &UserController{
		Db:       db,
		r:        r,
		planRepo: planRepo,
	}
}

type RegistrationBody struct {
	Name                 string `json:"name"`
	Email                string `json:"email"`
	Password             string `json:"password"`
	PasswordConfirmation string `json:"passwordConfirmation"`
}

func (u *UserController) HandleRegistration(w http.ResponseWriter, req *http.Request) {
	var body RegistrationBody

	err := json.NewDecoder(req.Body).Decode(&body)
	if err != nil {
		u.r.json(w, map[string]string{"message": err.Error()}, http.StatusBadRequest)
		return
	}

	if len(body.Email) < 1 {
		u.r.json(w, map[string]string{"message": "Please, provide a valid e-mail"}, http.StatusBadRequest)
		return
	}

	if len(body.Name) < 4 {
		u.r.json(w, map[string]string{"message": "Plase, provide a valid name"}, http.StatusBadRequest)
		return
	}

	passwordLenth := len(body.Password)

	if passwordLenth < 8 {
		u.r.json(w, map[string]string{"message": "Your password must have at least 8 characters"}, http.StatusBadRequest)
		return
	}

	if body.Password != body.PasswordConfirmation {
		u.r.json(w, map[string]string{"message": "The password and password confirmation do not match"}, http.StatusBadRequest)
		return
	}

	var total int
	err = u.Db.QueryRow("SELECT COUNT(*) as total FROM users WHERE email = ?", body.Email).Scan(&total)

	if err != nil {
		log.Printf("[ERROR] Failed to retrieve user with email %s: %v", body.Email, err)
	}

	if total > 0 {
		u.r.json(w, map[string]string{"message": "This is e-mail is already in use"}, http.StatusBadRequest)
		return
	}

	encryptedPassword := sec.EncryptText(body.Password)
	result, err := u.Db.Exec("INSERT INTO users (name, email, password) VALUES (?, ?, ?)", body.Name, body.Email, encryptedPassword)
	if err != nil {
		log.Printf("Failed to create user: %v", err)
		u.r.json(w, map[string]string{}, http.StatusInternalServerError)
	}

	id, err := result.LastInsertId()
	if err != nil {
		log.Printf("Failed to retrieve last insert id: %v", err)
		u.r.json(w, map[string]string{}, http.StatusInternalServerError)
	}

	u.planRepo.Store("Default Plan", "Your default plan", "enabled", int(id))

	u.r.json(w, map[string]string{"id": fmt.Sprintf("%d", id)}, http.StatusCreated)
}

type LoginBody struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (u *UserController) HandleLogin(w http.ResponseWriter, req *http.Request) {
	var body LoginBody

	err := json.NewDecoder(req.Body).Decode(&body)
	if err != nil {
		u.r.json(w, map[string]string{"message": err.Error()}, http.StatusBadRequest)
	}

	if len(body.Email) < 1 {
		u.r.json(w, map[string]string{"message": "Invalid e-mail"}, http.StatusBadRequest)
		return
	}

	if len(body.Password) < 1 {
		u.r.json(w, map[string]string{"message": "Invalid password"}, http.StatusBadRequest)
		return
	}

	result := u.Db.QueryRow("SELECT id, password FROM users WHERE email = ?", body.Email)

	if result.Err() != nil {
		u.r.json(w, map[string]string{"message": "Invalid e-mail or password"}, http.StatusBadRequest)
		return
	}

	var (
		id       int
		password string
	)
	err = result.Scan(&id, &password)

	if err != nil {
		u.r.json(w, map[string]string{"message": "Internal error. Please, try again later"}, http.StatusInternalServerError)
		return
	}

	passwordDecrypted := sec.DecryptText(password)
	if passwordDecrypted != body.Password {
		u.r.json(w, map[string]string{"message": "Invalid e-mail or password"}, http.StatusBadRequest)
		return
	}

	jwt, err := sec.JWTFromPayload(map[string]any{
		"id": strconv.Itoa(id),
	}, 24*time.Hour)
	if err != nil {
		u.r.json(w, map[string]string{}, http.StatusInternalServerError)
		return
	}

	u.r.json(w, map[string]string{"token": jwt}, http.StatusOK)
}

type UserJson struct {
	Id            int    `json:"id"`
	Name          string `json:"name"`
	Email         string `json:"email"`
	CurrentPlanId int    `json:"currentPlanId"`
}

func (u *UserController) HandleMe(w http.ResponseWriter, req *http.Request) {
	userId := req.Context().Value("user.id")

	result := u.Db.QueryRow("SELECT id, name, email FROM users WHERE id = ?", userId)

	if result.Err() != nil {
		u.r.json(w, map[string]string{"message": "Invalid user. Please, try logging in again"}, http.StatusInternalServerError)
		return
	}

	var user UserJson
	result.Scan(&user.Id, &user.Name, &user.Email)

	result = u.Db.QueryRow("SELECT id FROM plans WHERE user_id = ? LIMIT 1", userId)
	result.Scan(&user.CurrentPlanId)

	u.r.json(w, user, http.StatusOK)
}
