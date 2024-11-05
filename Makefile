build:
	@go build cmd/main.go

run:
	@go run cmd/main.go

dev:
	@tailwindcss -w -i src/templates/main.css -o assets/main.css &
	@templ generate -watch src/templates &
	@air
