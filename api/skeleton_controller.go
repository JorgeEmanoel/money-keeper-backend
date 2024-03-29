package api

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/JorgeEmanoel/money-keeper-backend/model"
	"github.com/gorilla/mux"
)

type SkeletonRepository interface {
	Store(name, description, direction, frequency, currency string, value float64, planId, ownerId int) (int, error)
	Delete(id int) error
	GetByUserId(userId int) ([]model.Skeleton, error)
	GetIncomingsByUserId(userId int) ([]model.Skeleton, error)
	GetOutcomingsByUserId(userId int) ([]model.Skeleton, error)
}

type SkeletonController struct {
	repo SkeletonRepository
	r    *Router
}

func MakeSkeletonController(repo SkeletonRepository, r *Router) *SkeletonController {
	return &SkeletonController{
		repo: repo,
		r:    r,
	}
}

type SkeletonJson struct {
	Id          int     `json:"id"`
	Name        string  `json:"name"`
	Description string  `json:"description"`
	Direction   string  `json:"direction"`
	Frequency   string  `json:"frequency"`
	Value       float64 `json:"value"`
	Currency    string  `json:"currency"`
}

func (c *SkeletonController) HandleList(w http.ResponseWriter, req *http.Request) {
	skeletons, err := c.repo.GetByUserId(req.Context().Value("user.id").(int))
	if err != nil {
		log.Printf("Failed to fetch user skeletons. User Id: %d, err: %v", req.Context().Value("user.id"), err)
		c.r.json(w, map[string]string{"message": "Failed to fetch skeletons. Please, try again later"}, http.StatusInternalServerError)
		return
	}

	var skeletonsResponse []SkeletonJson
	var total float64

	for _, skeleton := range skeletons {
		p := SkeletonJson{
			Id:          skeleton.Id,
			Name:        skeleton.Name,
			Description: skeleton.Description,
			Direction:   skeleton.Direction,
			Frequency:   skeleton.Frequency,
			Value:       skeleton.Value / 100,
			Currency:    skeleton.Currency,
		}

		skeletonsResponse = append(skeletonsResponse, p)
		total += skeleton.Value
	}

	response := map[string]any{
		"skeletons": skeletonsResponse,
		"total":     total,
	}

	c.r.json(w, response, http.StatusOK)
}

func (c *SkeletonController) HandleIncomingList(w http.ResponseWriter, req *http.Request) {
	skeletons, err := c.repo.GetIncomingsByUserId(req.Context().Value("user.id").(int))
	if err != nil {
		log.Printf("Failed to fetch user skeletons. User Id: %d, err: %v", req.Context().Value("user.id"), err)
		c.r.json(w, map[string]string{"message": "Failed to fetch skeletons. Please, try again later"}, http.StatusInternalServerError)
		return
	}

	var skeletonsResponse []SkeletonJson
	var total float64

	for _, skeleton := range skeletons {
		p := SkeletonJson{
			Id:          skeleton.Id,
			Name:        skeleton.Name,
			Description: skeleton.Description,
			Direction:   skeleton.Direction,
			Frequency:   skeleton.Frequency,
			Value:       skeleton.Value / 100,
			Currency:    skeleton.Currency,
		}

		skeletonsResponse = append(skeletonsResponse, p)
		total += skeleton.Value / 100
	}

	response := map[string]any{
		"skeletons": skeletonsResponse,
		"total":     total,
	}

	c.r.json(w, response, http.StatusOK)
}

func (c *SkeletonController) HandleOutocomingList(w http.ResponseWriter, req *http.Request) {
	skeletons, err := c.repo.GetOutcomingsByUserId(req.Context().Value("user.id").(int))
	if err != nil {
		log.Printf("Failed to fetch user skeletons. User Id: %d, err: %v", req.Context().Value("user.id"), err)
		c.r.json(w, map[string]string{"message": "Failed to fetch skeletons. Please, try again later"}, http.StatusInternalServerError)
		return
	}

	var skeletonsResponse []SkeletonJson
	var total float64

	for _, skeleton := range skeletons {
		p := SkeletonJson{
			Id:          skeleton.Id,
			Name:        skeleton.Name,
			Description: skeleton.Description,
			Direction:   skeleton.Direction,
			Frequency:   skeleton.Frequency,
			Value:       skeleton.Value / 100,
			Currency:    skeleton.Currency,
		}

		skeletonsResponse = append(skeletonsResponse, p)
		total += skeleton.Value / 100
	}

	response := map[string]any{
		"skeletons": skeletonsResponse,
		"total":     total,
	}

	c.r.json(w, response, http.StatusOK)
}

type CreateSkeletonBody struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	Direction   string `json:"direction"`
	Frequency   string `json:"frequency"`
	Value       string `json:"value"`
	Currency    string `json:"currency"`
}

func (c *SkeletonController) HandleCreate(w http.ResponseWriter, req *http.Request) {
	var body CreateSkeletonBody
	planId, _ := strconv.Atoi(mux.Vars(req)["planId"])

	err := json.NewDecoder(req.Body).Decode(&body)
	if err != nil {
		c.r.json(w, map[string]string{"message": err.Error()}, http.StatusBadRequest)
		return
	}

	cleanValue := strings.ReplaceAll(body.Value, ".", "")
	cleanValue = strings.ReplaceAll(cleanValue, ",", ".")

	value, _ := strconv.ParseFloat(cleanValue, 64)

	id, err := c.repo.Store(
		body.Name,
		body.Description,
		body.Direction,
		body.Frequency,
		body.Currency,
		value,
		planId,
		req.Context().Value("user.id").(int),
	)
	if err != nil {
		log.Printf("Failed to store skeleton. Body %v, Err: %v\n", body, err)
		c.r.json(w, map[string]string{"message": "Creation failed, Please, try again later"}, http.StatusInternalServerError)
		return
	}

	c.r.json(w, map[string]int{"id": id}, http.StatusCreated)
}

func (c *SkeletonController) HandleDelete(w http.ResponseWriter, req *http.Request) {
	params := mux.Vars(req)

	id, _ := strconv.Atoi(params["id"])

	err := c.repo.Delete(id)
	if err != nil {
		log.Printf("Failed to delete skeleton. Id: %d, err: %v", id, err)
		c.r.json(w, map[string]string{"message": "Failed to delete skeleton. Please, try again later"}, http.StatusInternalServerError)
		return
	}

	c.r.json(w, map[string]any{}, http.StatusOK)
}
