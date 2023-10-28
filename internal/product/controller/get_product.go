package controller

import (
	"context"

	core "github.com/vincentcent6/product-catalogue/internal/product/core"
)

func (p *ctrl) GetProduct(ctx context.Context, input GetProductInput) (result GetProductResult, err error) {
	res, err := prd.GetProduct(ctx, convertGetData(input))
	if err != nil {
		return GetProductResult{}, err
	}

	return parseGetResult(res), err
}

func (p *ctrl) GetProducts(ctx context.Context, input GetProductInput) (result []GetProductResult, err error) {
	prdsData, err := prd.GetProducts(ctx, convertGetData(input))
	if err != nil {
		return []GetProductResult{}, err
	}

	return parseGetResults(prdsData), err
}

func convertGetData(input GetProductInput) core.GetInput {
	var getPrdFilter *core.GetFilter
	if input.Filter != nil {
		getPrdFilter = &core.GetFilter{}
		if input.Filter.SKU != "" {
			getPrdFilter.SKU = input.Filter.SKU
		}
		if input.Filter.Description != "" {
			getPrdFilter.Description = input.Filter.Description
		}
		if input.Filter.Title != "" {
			getPrdFilter.Title = input.Filter.Title
		}
		if input.Filter.Category != "" {
			getPrdFilter.Category = input.Filter.Category
		}
	}

	var getPrdSort *core.GetProductSort
	if input.Sort != nil {
		getPrdSort = &core.GetProductSort{
			CreateTime: input.Sort.CreateTime,
			AvgRating:  input.Sort.AvgRating,
		}
	}

	return core.GetInput{
		ProductID:      input.ProductID,
		GetFilter:      getPrdFilter,
		GetProductSort: getPrdSort,
	}
}

func parseGetResult(prdData core.ProductDataWithoutCreateTime) GetProductResult {
	images := []ProductImage{}
	for _, img := range prdData.Images {
		images = append(images, ProductImage{
			URL:         img.URL,
			Description: img.Description,
		})
	}

	return GetProductResult{
		ProductID: prdData.ProductID,
		Attribute: ProductInput{
			SKU:            prdData.SKU,
			Title:          prdData.Title,
			Description:    prdData.Description,
			Category:       prdData.Category,
			Etalase:        prdData.Etalase,
			Images:         images,
			Weight:         prdData.Weight,
			Price:          prdData.Price,
			ProductReviews: prdData.ProductReviews,
		},
	}
}

func parseGetResults(prdData []core.ProductDataWithoutCreateTime) []GetProductResult {
	res := []GetProductResult{}
	for _, prdData := range prdData {
		images := []ProductImage{}
		for _, img := range prdData.Images {
			images = append(images, ProductImage{
				URL:         img.URL,
				Description: img.Description,
			})
		}

		res = append(res, GetProductResult{
			ProductID: prdData.ProductID,
			Attribute: ProductInput{
				SKU:            prdData.SKU,
				Title:          prdData.Title,
				Description:    prdData.Description,
				Category:       prdData.Category,
				Etalase:        prdData.Etalase,
				Images:         images,
				Weight:         prdData.Weight,
				Price:          prdData.Price,
				ProductReviews: prdData.ProductReviews,
			},
		})
	}
	return res
}
