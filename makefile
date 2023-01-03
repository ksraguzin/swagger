generate-server:
	@mkdir -p ./server/internal/handler
	@swagger generate server -f ./server/api/api.yml -t ./server -A "service-pdf-compose" --server-package=internal/handler
generate-client:
	@swagger generate client -f ./server/api/api.yml -t ./controller -A "service-pdf-compose"
