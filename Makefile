setup:
	@docker-compose up -d
run:
	@go build
	@./go-ekyc