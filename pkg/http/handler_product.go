package http

import (
	"context"
	"encoding/json"
	"net/http"
	"strconv"

	productCtrl "github.com/vincentcent6/product-catalogue/internal/controller/product"
	"github.com/vincentcent6/product-catalogue/pkg/response"
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

func UpdateProduct(w http.ResponseWriter, r *http.Request) {
	ctx := context.TODO()
	// prepare input obj
	var input productCtrl.UpdateProductInput

	// decode JSON input
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&input)
	if err != nil {
		response.WriteError(w, http.StatusUnprocessableEntity, err.Error())
		return
	}

	err = prdCtrl.UpdateProduct(ctx, input)
	if err != nil {
		response.WriteError(w, http.StatusBadRequest, err.Error())
		return
	}

	response.WriteSuccess(w, http.StatusOK, nil)

	return
}

func GetProduct(w http.ResponseWriter, r *http.Request) {
	ctx := context.TODO()

	productIDStr := r.FormValue("product_id")
	productID, err := strconv.ParseInt(productIDStr, 10, 64)
	if err != nil {
		response.WriteError(w, http.StatusBadRequest, err.Error())
		return
	}

	res, err := prdCtrl.GetProduct(ctx, productCtrl.GetProductInput{
		ProductID: productID,
	})
	if err != nil {
		response.WriteError(w, http.StatusBadRequest, err.Error())
		return
	}

	response.WriteSuccess(w, http.StatusOK, res)

	return
}

func GetProducts(w http.ResponseWriter, r *http.Request) {
	ctx := context.TODO()

	sku := r.FormValue("sku")
	title := r.FormValue("title")
	description := r.FormValue("description")
	category := r.FormValue("category")

	var filter *productCtrl.GetProductFilter
	if sku != "" || title != "" || description != "" || category != "" {
		filter = &productCtrl.GetProductFilter{}
		if sku != "" {
			filter.SKU = sku
		}
		if title != "" {
			filter.Title = title
		}
		if description != "" {
			filter.Description = description
		}
		if category != "" {
			filter.Category = category
		}
	}

	var sort *productCtrl.GetProductSort
	sortByCreateTimeStr := r.FormValue("sort_ct")
	if sortByCreateTimeStr != "" {
		sortByCreateTime, _ := strconv.Atoi(sortByCreateTimeStr)
		if sortByCreateTime > 0 {
			sort = &productCtrl.GetProductSort{}
			if sortByCreateTime > 0 {
				sort.CreateTime = sortByCreateTime
			}
		}
	}

	res, err := prdCtrl.GetProducts(ctx, productCtrl.GetProductInput{
		Filter: filter,
		Sort:   sort,
	})
	if err != nil {
		response.WriteError(w, http.StatusBadRequest, err.Error())
		return
	}

	response.WriteSuccess(w, http.StatusOK, res)

	return
}
