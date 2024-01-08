package api

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"github.com/JorgeEmanoel/money-keeper-backend/model"
	"github.com/gorilla/mux"
)

type PlanRepository interface {
	Store(name, description, status string, ownerId int) (int, error)
	Delete(id int) error
	GetByUserId(userId int) ([]model.Plan, error)
}

type PlanController struct {
	repo            PlanRepository
	transactionRepo TransactionRepository
	r               *Router
}

func MakePlanController(repo PlanRepository, transactionRepo TransactionRepository, r *Router) *PlanController {
	return &PlanController{
		repo:            repo,
		transactionRepo: transactionRepo,
		r:               r,
	}
}

type PlanJson struct {
	Id          int    `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Status      string `json:"status"`
}

func (c *PlanController) HandleList(w http.ResponseWriter, req *http.Request) {
	plans, err := c.repo.GetByUserId(req.Context().Value("user.id").(int))
	if err != nil {
		log.Printf("Failed to fetch user plans. User Id: %d, err: %v", req.Context().Value("user.id"), err)
		c.r.json(w, map[string]string{"message": "Failed to fetch plans. Please, try again later"}, http.StatusInternalServerError)
		return
	}

	var plansResponse []PlanJson

	for _, plan := range plans {
		p := PlanJson{
			Id:          plan.Id,
			Name:        plan.Name,
			Description: plan.Description,
			Status:      plan.Status,
		}

		plansResponse = append(plansResponse, p)
	}

	response := map[string][]PlanJson{
		"plans": plansResponse,
	}

	c.r.json(w, response, http.StatusOK)
}

type CreatePlanBody struct {
	Name        string `json:"name"`
	Description string `json:"description"`
}

func (c *PlanController) HandleCreate(w http.ResponseWriter, req *http.Request) {
	var body CreatePlanBody

	err := json.NewDecoder(req.Body).Decode(&body)
	if err != nil {
		c.r.json(w, map[string]string{"message": err.Error()}, http.StatusBadRequest)
	}

	id, err := c.repo.Store(body.Name, body.Description, "enabled", req.Context().Value("user.id").(int))
	if err != nil {
		log.Printf("Failed to store plan. Body %v, Err: %v\n", body, err)
		c.r.json(w, map[string]string{"message": "Creation failed, Please, try again later"}, http.StatusInternalServerError)
		return
	}

	c.r.json(w, map[string]int{"id": id}, http.StatusCreated)
}

func (c *PlanController) HandleDelete(w http.ResponseWriter, req *http.Request) {
	params := mux.Vars(req)

	id, _ := strconv.Atoi(params["id"])

	err := c.repo.Delete(id)
	if err != nil {
		log.Printf("Failed to delete plan. Id: %d, err: %v", id, err)
		c.r.json(w, map[string]string{"message": "Failed to delete plan. Please, try again later"}, http.StatusInternalServerError)
		return
	}

	c.r.json(w, map[string]any{}, http.StatusOK)
}

func (c *PlanController) HandleSummary(w http.ResponseWriter, req *http.Request) {
	reference := mux.Vars(req)["reference"]

	transactions, err := c.transactionRepo.GetByUserIdFromReference(
		req.Context().Value("user.id").(int),
		reference,
	)
	if err != nil {
		c.r.json(w, map[string]string{"message": "Failed to retrieve summary"}, http.StatusInternalServerError)
		return
	}

	var (
		totalIncomings  = 0
		totalOutcomings = 0
	)

	for _, transaction := range transactions {
		if transaction.Direction == "income" {
			totalIncomings += transaction.Value
		}

		if transaction.Direction == "outcome" {
			totalOutcomings += transaction.Value
		}
	}

	balance := totalIncomings - totalOutcomings

	c.r.json(w, map[string]any{
		"totalIncomings":  totalIncomings / 100,
		"totalOutcomings": totalOutcomings / 100,
		"balance":         balance / 100,
	}, http.StatusOK)
}
