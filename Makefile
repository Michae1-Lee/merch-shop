run:
	docker-compose build
	docker-compose up -d
test:
	docker exec -it avito-shop-service go test -v ./...