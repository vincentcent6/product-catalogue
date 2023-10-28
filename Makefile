.PHONY: run
run:
	@echo ">> building product-catalogue binaries."
	@go build -o cmd/product-catalogue/product-catalogue cmd/product-catalogue/main.go
	@echo ">> product-catalogue is built."
	@./cmd/product-catalogue/product-catalogue
	@echo ">> product-catalogue is running."