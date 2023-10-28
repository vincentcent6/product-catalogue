package core

import (
	"context"
	"database/sql/driver"
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/elastic/go-elasticsearch/v8/typedapi/core/search"
	"github.com/elastic/go-elasticsearch/v8/typedapi/core/update"
	"github.com/elastic/go-elasticsearch/v8/typedapi/types"
	"github.com/elastic/go-elasticsearch/v8/typedapi/types/enums/sortorder"
	prdReviewCore "github.com/vincentcent6/product-catalogue/internal/product_review/core"
)

type (
	// Product is an interface for product object
	Product interface {
		CreateProduct(ctx context.Context, data Data) error
		UpdateProduct(ctx context.Context, input UpdateInput) error
		GetProduct(ctx context.Context, inp GetInput) (ProductDataWithoutCreateTime, error)
		GetProducts(ctx context.Context, input GetInput) ([]ProductDataWithoutCreateTime, error)
	}
	product struct{}

	ProductImages []ProductImage

	Data struct {
		SKU         string        `db:"sku" json:"sku,omitempty"`
		Title       string        `db:"title" json:"title,omitempty"`
		Description string        `db:"description" json:"description,omitempty"`
		Category    string        `db:"category" json:"category,omitempty"`
		Etalase     string        `db:"etalase" json:"etalase,omitempty"`
		Images      ProductImages `db:"images" json:"images,omitempty"`
		Weight      float64       `db:"weight" json:"weight,omitempty"`
		Price       int64         `db:"price" json:"price,omitempty"`
		CreateTime  time.Time     `db:"create_time" json:"create_time,omitempty"`
	}

	UpdateInput struct {
		ProductID int64 `db:"product_id" json:"product_id,omitempty"`
		Data      Data
	}

	ProductData struct {
		ProductID   int64         `db:"product_id" json:"product_id,omitempty"`
		SKU         string        `db:"sku" json:"sku,omitempty"`
		Title       string        `db:"title" json:"title,omitempty"`
		Description string        `db:"description" json:"description,omitempty"`
		Category    string        `db:"category" json:"category,omitempty"`
		Etalase     string        `db:"etalase" json:"etalase,omitempty"`
		Images      ProductImages `db:"images" json:"images,omitempty"`
		Weight      float64       `db:"weight" json:"weight,omitempty"`
		Price       int64         `db:"price" json:"price,omitempty"`
		CreateTime  time.Time     `db:"create_time" json:"create_time,omitempty"`
	}

	ProductDataWithoutCreateTime struct {
		ProductID      int64                              `db:"product_id" json:"product_id,omitempty"`
		SKU            string                             `db:"sku" json:"sku,omitempty"`
		Title          string                             `db:"title" json:"title,omitempty"`
		Description    string                             `db:"description" json:"description,omitempty"`
		Category       string                             `db:"category" json:"category,omitempty"`
		Etalase        string                             `db:"etalase" json:"etalase,omitempty"`
		Images         ProductImages                      `db:"images" json:"images,omitempty"`
		Weight         float64                            `db:"weight" json:"weight,omitempty"`
		Price          int64                              `db:"price" json:"price,omitempty"`
		ProductReviews []prdReviewCore.ProductReviewInput `json:"product_reviews,omitempty"`
	}

	ProductImage struct {
		URL         string `db:"url" json:"url,omitempty"`
		Description string `db:"description" json:"description,omitempty"`
	}

	GetInput struct {
		ProductID      int64           `db:"product_id" json:"product_id,omitempty"`
		GetFilter      *GetFilter      `json:"filter,omitempty"`
		GetProductSort *GetProductSort `json:"sort,omitempty"`
	}

	GetFilter struct {
		SKU         string `db:"sku" json:"sku,omitempty"`
		Title       string `db:"title" json:"title,omitempty"`
		Description string `db:"description" json:"description,omitempty"`
		Category    string `db:"category" json:"category,omitempty"`
	}

	GetProductSort struct {
		CreateTime int `json:"sort_create_time,omitempty"`
	}
)

const (
	sortByAsc  = 1
	sortByDesc = 2
)

func (p *product) CreateProduct(ctx context.Context, data Data) error {
	imgJson, err := json.Marshal(data.Images)
	if err != nil {
		return err
	}
	tx, err := db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()
	var productID int64
	var createTime time.Time
	err = tx.QueryRow(qCreateProduct, data.SKU, data.Title, data.Description, data.Category, data.Etalase, imgJson, data.Weight, data.Price).Scan(&productID, &createTime)
	if err != nil {
		return err
	}
	// set create time to be inserted to es
	data.CreateTime = createTime
	err = createToES(ctx, productID, data)
	if err != nil {
		return err
	}
	err = tx.Commit()
	if err != nil {
		return err
	}

	return nil
}

func (p *product) UpdateProduct(ctx context.Context, input UpdateInput) error {
	data := input.Data
	params := []string{"", "", "", "", "", "", "", ""}

	imgJson, err := json.Marshal(data.Images)
	if err != nil {
		return err
	}
	// check if the input is set
	if data.SKU != "" {
		params[0] = fmt.Sprintf("sku='%s'", data.SKU)
	}
	if data.Title != "" {
		params[1] = fmt.Sprintf("title='%s'", data.Title)
	}
	if data.Description != "" {
		params[2] = fmt.Sprintf("description='%s'", data.Description)
	}
	if data.Category != "" {
		params[3] = fmt.Sprintf("category='%s'", data.Category)
	}
	if data.Etalase != "" {
		params[4] = fmt.Sprintf("etalase='%s'", data.Etalase)
	}
	if len(data.Images) > 0 {
		params[5] = fmt.Sprintf("images='%v'", string(imgJson))
	}
	if data.Weight != 0 {
		params[6] = fmt.Sprintf("weight=%f", data.Weight)
	}
	if data.Price != 0 {
		params[7] = fmt.Sprintf("price=%d", data.Price)
	}
	paramsStr := strings.Join(params, ", ")

	tx, err := db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()
	res, err := tx.Exec(fmt.Sprintf(qUpdateProduct, paramsStr), input.ProductID)
	if err != nil {
		return err
	}

	if rows, _ := res.RowsAffected(); rows == 0 {
		return ErrDataNotFound
	}

	err = updateToES(ctx, input)
	if err != nil {
		return err
	}

	err = tx.Commit()
	if err != nil {
		return err
	}
	return nil
}

// for now, get product will return 1 result
func (p *product) GetProduct(ctx context.Context, input GetInput) (ProductDataWithoutCreateTime, error) {
	esRes, err := getFromESByProductID(ctx, input)
	if esRes.ProductID > 0 {
		return esRes, nil
	}

	params := []string{""}
	if input.ProductID != 0 {
		params[0] = fmt.Sprintf("product_id=%d", input.ProductID)
	}
	paramsStr := strings.Join(params, ", ")

	prd := []ProductDataWithoutCreateTime{}
	err = db.Select(&prd, fmt.Sprintf(qGetProduct, paramsStr))
	if err != nil {
		return ProductDataWithoutCreateTime{}, err
	}
	if len(prd) == 0 {
		return ProductDataWithoutCreateTime{}, ErrDataNotFound
	}
	return prd[0], nil
}

func (p *product) GetProducts(ctx context.Context, input GetInput) ([]ProductDataWithoutCreateTime, error) {
	return getProductsFromES(ctx, input)
}

// helper func
func createToES(ctx context.Context, productID int64, data Data) error {
	prdData := ProductData{
		ProductID:   productID,
		SKU:         data.SKU,
		Title:       data.Title,
		Description: data.Description,
		Category:    data.Category,
		Etalase:     data.Etalase,
		Images:      data.Images,
		Weight:      data.Weight,
		Price:       data.Price,
		CreateTime:  data.CreateTime,
	}
	_, err := es.Index(productsIndex).Id(strconv.FormatInt(productID, 10)).Request(prdData).Do(ctx)
	if err != nil {
		return err
	}
	return nil
}

func updateToES(ctx context.Context, input UpdateInput) error {
	jsonData, _ := json.Marshal(ProductDataWithoutCreateTime{
		ProductID:   input.ProductID,
		SKU:         input.Data.SKU,
		Title:       input.Data.Title,
		Description: input.Data.Description,
		Category:    input.Data.Category,
		Etalase:     input.Data.Etalase,
		Images:      input.Data.Images,
		Weight:      input.Data.Weight,
		Price:       input.Data.Price,
	})
	_, err := es.Update(productsIndex, strconv.FormatInt(input.ProductID, 10)).
		Request(&update.Request{
			Doc: json.RawMessage(jsonData),
		}).Do(ctx)
	if err != nil {
		return err
	}
	return nil
}

func getFromESByProductID(ctx context.Context, input GetInput) (ProductDataWithoutCreateTime, error) {
	res, err := es.Get(productsIndex, strconv.FormatInt(input.ProductID, 10)).Do(ctx)
	if err != nil {
		return ProductDataWithoutCreateTime{}, err
	}
	var prdData ProductDataWithoutCreateTime
	err = json.Unmarshal(res.Source_, &prdData)
	if err != nil {
		return ProductDataWithoutCreateTime{}, err
	}
	return prdData, nil
}

func getProductsFromES(ctx context.Context, input GetInput) ([]ProductDataWithoutCreateTime, error) {
	prdsData := []ProductDataWithoutCreateTime{}
	searchFunc := es.Search().Index(productsIndex)

	filterQueries := []types.Query{}
	if input.GetFilter != nil {
		if input.GetFilter.SKU != "" {
			filterQueries = append(filterQueries, types.Query{
				Term: map[string]types.TermQuery{
					"sku": {Value: input.GetFilter.SKU},
				},
			})
		}
		if input.GetFilter.Description != "" {
			filterQueries = append(filterQueries, types.Query{
				Term: map[string]types.TermQuery{
					"description": {Value: input.GetFilter.Description},
				},
			})
		}
		if input.GetFilter.Category != "" {
			filterQueries = append(filterQueries, types.Query{
				Term: map[string]types.TermQuery{
					"category": {Value: input.GetFilter.Category},
				},
			})
		}
		if input.GetFilter.Title != "" {
			filterQueries = append(filterQueries, types.Query{
				Term: map[string]types.TermQuery{
					"title": {Value: input.GetFilter.Title},
				},
			})
		}
	}

	sortCombinations := []types.SortCombinations{}
	if input.GetProductSort != nil {
		if input.GetProductSort.CreateTime == sortByAsc {
			sortCombinations = append(sortCombinations, types.SortOptions{SortOptions: map[string]types.FieldSort{
				"create_time": {Order: &sortorder.Asc},
			}})
		} else if input.GetProductSort.CreateTime == sortByDesc {
			sortCombinations = append(sortCombinations, types.SortOptions{SortOptions: map[string]types.FieldSort{
				"create_time": {Order: &sortorder.Desc},
			}})
		}
	}

	// filter and sort
	if len(filterQueries) > 0 && len(sortCombinations) > 0 {
		searchFunc = searchFunc.Request(&search.Request{
			PostFilter: &types.Query{
				Bool: &types.BoolQuery{
					Must: []types.Query{
						{
							MatchAll: &types.MatchAllQuery{},
						},
					},
					Filter: filterQueries,
				},
			},
			Sort: sortCombinations,
		})
	} else if len(filterQueries) > 0 { // filter only
		searchFunc = searchFunc.Request(&search.Request{
			PostFilter: &types.Query{
				Bool: &types.BoolQuery{
					Must: []types.Query{
						{
							MatchAll: &types.MatchAllQuery{},
						},
					},
					Filter: filterQueries,
				},
			},
		})
	} else if len(sortCombinations) > 0 { // sort only
		searchFunc = searchFunc.Request(&search.Request{
			Sort: sortCombinations,
		})
	}

	res, err := searchFunc.Do(ctx)
	if err != nil {
		return []ProductDataWithoutCreateTime{}, err
	}
	hits := res.Hits.Hits
	for _, hit := range hits {
		var prdData ProductDataWithoutCreateTime
		err := json.Unmarshal(hit.Source_, &prdData)
		if err != nil {
			return []ProductDataWithoutCreateTime{}, err
		}
		prdsData = append(prdsData, prdData)
	}
	return prdsData, nil
}

func (images ProductImages) Value() (driver.Value, error) {
	return json.Marshal(images)
}

func (images *ProductImages) Scan(val interface{}) (err error) {
	switch v := val.(type) {
	case []byte:
		return json.Unmarshal(v, images)
	case string:
		return json.Unmarshal([]byte(v), &images)
	default:
		return errors.New(fmt.Sprintf("Unsupported type: %T", v))
	}
}
