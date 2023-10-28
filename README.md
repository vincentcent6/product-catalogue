# product-catalogue
Product Catalogue

## Prerequisite
1. Golang
2. Docker

## Tech Stack
1. Golang
2. Docker
3. Postgres
4. ElasticSearch

## Step to run
1. Run Docker compose up
`docker compose up --build -d`

2. Please execute Rest API - InitNestedReviews in Postman before proceeding to the next step

3. Make run or run below commands
```
go build -o cmd/product-catalogue/product-catalogue cmd/product-catalogue/main.go
./cmd/product-catalogue/product-catalogue
```
