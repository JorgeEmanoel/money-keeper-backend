package api

import (
	"log"
	"net/http"
	"strconv"

	"github.com/JorgeEmanoel/money-keeper-backend/model"
	"github.com/gorilla/mux"
)

var (
	TRANSACTION_STATUS_PENDING  = "pending"
	TRANSACTION_STATUS_PAID     = "paid"
	TRANSACTION_STATUS_CANCELED = "canceled"
)

type TransactionRepository interface {
	Store(name, description, direction, reference, currency, status string, value, planId, ownerId int) (int, error)
	Delete(id int) error
	GetByUserId(userId int) ([]model.Transaction, error)
	ChangeStatus(id int, status string) error
}

type TransactionController struct {
	repo TransactionRepository
	r    *Router
}

func MakeTransactionController(repo TransactionRepository, r *Router) *TransactionController {
	return &TransactionController{
		repo: repo,
		r:    r,
	}
}

type TransactionJson struct {
	Id          int    `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Direction   string `json:"direction"`
	Reference   string `json:"reference"`
	Value       int    `json:"value"`
	Currency    string `json:"currency"`
	Status      string `json:"status"`
}

func (c *TransactionController) HandleList(w http.ResponseWriter, req *http.Request) {
	transactions, err := c.repo.GetByUserId(req.Context().Value("user.id").(int))
	if err != nil {
		log.Printf("Failed to fetch user transactions. User Id: %d, err: %v", req.Context().Value("user.id"), err)
		c.r.json(w, map[string]string{"message": "Failed to fetch transactions. Please, try again later"}, http.StatusInternalServerError)
		return
	}

	var transactionsResponse []TransactionJson

	for _, transaction := range transactions {
		p := TransactionJson{
			Id:          transaction.Id,
			Name:        transaction.Name,
			Description: transaction.Description,
			Direction:   transaction.Direction,
			Reference:   transaction.Reference,
			Value:       transaction.Value,
			Currency:    transaction.Currency,
			Status:      transaction.Status,
		}

		transactionsResponse = append(transactionsResponse, p)
	}

	response := map[string][]TransactionJson{
		"transactions": transactionsResponse,
	}

	c.r.json(w, response, http.StatusOK)
}

func (c *TransactionController) HandleChangeStatus(w http.ResponseWriter, req *http.Request) {
	params := mux.Vars(req)

	id, _ := strconv.Atoi(params["id"])
	status := params["status"]

	if status != TRANSACTION_STATUS_PAID || status != TRANSACTION_STATUS_PENDING || status != TRANSACTION_STATUS_CANCELED {
		c.r.json(w, map[string]string{"message": "Invalid status"}, http.StatusBadRequest)
	}

	err := c.repo.ChangeStatus(id, status)
	if err != nil {
		log.Printf("Failed to update status: transaction = %d, status = %s", id, status)
		c.r.json(w, map[string]string{"message": "Failed to update transaction. Please, try again later"}, http.StatusInternalServerError)
		return
	}

	c.r.json(w, map[string]any{}, http.StatusNoContent)
}

func (c *TransactionController) HandleDelete(w http.ResponseWriter, req *http.Request) {
	params := mux.Vars(req)

	id, _ := strconv.Atoi(params["id"])

	err := c.repo.Delete(id)
	if err != nil {
		log.Printf("Failed to delete transaction. Id: %d, err: %v", id, err)
		c.r.json(w, map[string]string{"message": "Failed to delete transaction. Please, try again later"}, http.StatusInternalServerError)
		return
	}

	c.r.json(w, map[string]any{}, http.StatusOK)
}
