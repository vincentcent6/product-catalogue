package core

import (
	"context"
	"encoding/json"
	"strconv"

	"github.com/elastic/go-elasticsearch/v8/typedapi/core/update"
	"github.com/elastic/go-elasticsearch/v8/typedapi/types"
	"github.com/elastic/go-elasticsearch/v8/typedapi/types/enums/scriptlanguage"
)

type (
	// ProductReview is an interface for product review object
	ProductReview interface {
		CreateProductReview(ctx context.Context, data Data) error
	}
	productReview struct{}

	Data struct {
		ProductID     int64  `db:"product_id" json:"product_id,omitempty"`
		ReviewComment string `db:"review_commment" json:"review_comment,omitempty"`
		Rating        int    `db:"rating" json:"rating,omitempty"`
	}

	ProductReviewInput struct {
		ReviewID      int64  `db:"review_id" json:"review_id,omitempty"`
		ReviewComment string `db:"review_commment" json:"review_comment,omitempty"`
		Rating        int    `db:"rating" json:"rating,omitempty"`
	}
)

func (p *productReview) CreateProductReview(ctx context.Context, data Data) error {
	tx, err := db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()
	var reviewID int64
	err = tx.QueryRow(qCreateProductReview, data.ProductID, data.Rating, data.ReviewComment).Scan(&reviewID)
	if err != nil {
		return err
	}
	err = upsertReviewtoES(ctx, reviewID, data)
	if err != nil {
		return err
	}
	err = tx.Commit()
	if err != nil {
		return err
	}

	return nil
}

func upsertReviewtoES(ctx context.Context, reviewID int64, data Data) error {
	jsonData, _ := json.Marshal(ProductReviewInput{
		ReviewID:      reviewID,
		Rating:        data.Rating,
		ReviewComment: data.ReviewComment,
	})
	_, err := es.Update(productsIndex, strconv.FormatInt(data.ProductID, 10)).
		Request(&update.Request{
			Script: types.InlineScript{
				Source: "if (ctx._source.containsKey('product_reviews')) {ctx._source['product_reviews'].add(params.review)} else { ctx._source['product_reviews'] = [params.review] }",
				Params: map[string]json.RawMessage{
					"review": json.RawMessage(jsonData),
				},
				Lang: &scriptlanguage.Painless,
			},
		}).Do(ctx)
	if err != nil {
		return err
	}
	return nil
}
