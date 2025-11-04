run:
	@templ generate
	@go run cmd/api/main.go

build:
	@templ generate
	@sqlc generate
	@go build -o tmp/main cmd/api/main.go

setup:
	@curl -o cmd/web/assets/js/datastar.js https://cdn.jsdelivr.net/gh/starfederation/datastar@main/bundles/datastar.js
	@curl -o cmd/web/assets/css/fonts/jersey15.woff2 https://fonts.gstatic.com/s/jersey15/v3/_6_9EDzuROGsUuk2TWjiZYAg.woff2
