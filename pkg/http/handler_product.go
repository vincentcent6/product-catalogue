package http

import (
	"context"
	"encoding/json"
	"net/http"

	productCtrl "github.com/vincentcent6/product-catalogue/internal/controller/product"
	"github.com/vincentcent6/product-catalogue/pkg/response"
)

type (
	ProductInput struct {
		ProductInput productCtrl.ProductInput `json:"product"`
	}
)

func CreateProduct(w http.ResponseWriter, r *http.Request) {
	ctx := context.TODO()
	// prepare input obj
	var input productCtrl.ProductInput

	// decode JSON input
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&input)
	if err != nil {
		response.WriteError(w, http.StatusUnprocessableEntity, err.Error())
		return
	}

	res, err := prdCtrl.CreateProduct(ctx, input)
	if err != nil {
		response.WriteError(w, http.StatusBadRequest, err.Error())
		return
	}

	response.WriteSuccess(w, http.StatusOK, res)

	return
}
