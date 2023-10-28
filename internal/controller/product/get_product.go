package controller

import (
	"context"

	"github.com/vincentcent6/product-catalogue/internal/core"
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
			SKU:         prdData.SKU,
			Title:       prdData.Title,
			Description: prdData.Description,
			Category:    prdData.Category,
			Etalase:     prdData.Etalase,
			Images:      images,
			Weight:      prdData.Weight,
			Price:       prdData.Price,
		},
	}
}

func parseGetResults(prdDatasData []core.ProductDataWithoutCreateTime) []GetProductResult {
	res := []GetProductResult{}
	for _, prdData := range prdDatasData {
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
				SKU:         prdData.SKU,
				Title:       prdData.Title,
				Description: prdData.Description,
				Category:    prdData.Category,
				Etalase:     prdData.Etalase,
				Images:      images,
				Weight:      prdData.Weight,
				Price:       prdData.Price,
			},
		})
	}
	return res
}
