run:
	@templ generate
	@go run cmd/api/main.go

build:
	@templ generate
	@sqlc generate
	@go build -o tmp/main cmd/api/main.go

setup:
	@command -v templ >/dev/null 2>&1 || { \
		echo "Installing templ..."; \
		go install github.com/a-h/templ/cmd/templ@latest; \
	}
	@command -v sqlc >/dev/null 2>&1 || { \
		echo "Installing sqlc..."; \
		go install github.com/sqlc-dev/sqlc/cmd/sqlc@latest; \
	}
	@curl -o cmd/web/assets/js/datastar.js https://cdn.jsdelivr.net/gh/starfederation/datastar@main/bundles/datastar.js
	@curl -o cmd/web/assets/css/pico.min.css https://cdn.jsdelivr.net/npm/@picocss/pico@2/css/pico.min.css
	@curl -o cmd/web/assets/css/fonts/jersey15.woff2 https://fonts.gstatic.com/s/jersey15/v3/_6_9EDzuROGsUuk2TWjiZYAg.woff2
