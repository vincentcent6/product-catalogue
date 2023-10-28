package http

import (
	"context"
	"encoding/json"
	"net/http"

	productReviewCtrl "github.com/vincentcent6/product-catalogue/internal/product_review/controller"
	"github.com/vincentcent6/product-catalogue/pkg/response"
)

func CreateProductReview(w http.ResponseWriter, r *http.Request) {
	ctx := context.TODO()
	// prepare input obj
	var input productReviewCtrl.ProductReviewInput

	// decode JSON input
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&input)
	if err != nil {
		response.WriteError(w, http.StatusUnprocessableEntity, err.Error())
		return
	}

	res, err := prdReviewCtrl.CreateProductReview(ctx, input)
	if err != nil {
		response.WriteError(w, http.StatusBadRequest, err.Error())
		return
	}

	response.WriteSuccess(w, http.StatusOK, res)

	return
}
