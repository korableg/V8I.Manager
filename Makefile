.PHONY: lint
lint:
	docker run --rm -v ${CURDIR}:/app -w /app golangci/golangci-lint:v1.44.0 golangci-lint run -v