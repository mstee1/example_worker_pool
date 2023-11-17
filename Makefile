.PHONY: cover

cover:
	go test internal/config/* -coverprofile=cov.out ./...
	go tool cover -html=cov.out

	rm cov.out