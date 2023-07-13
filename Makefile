setup:
	@docker-compose up -d
run:
	@go build
	@./go-ekyc
setup-down:
	@docker-compose down