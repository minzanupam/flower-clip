build:
	@go build cmd/main.go

run:
	@go run cmd/main.go

air:
	@air &
templ:
	@templ generate -watch src/templates &
tailwind:
	@tailwindcss -w -i src/templates/main.css -o assets/main.css

dev: templ air
