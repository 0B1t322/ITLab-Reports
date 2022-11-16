reports_http = ./internal/controllers/http/reports/v1/dto,./internal/controllers/http/reports/v1,./internal/controllers/http/reports/v2/dto,./internal/controllers/http/reports/v2
draft_http = ./internal/controllers/http/drafts/v1/dto,./internal/controllers/http/drafts/v1,./internal/controllers/http/drafts/v2/dto,./internal/controllers/http/drafts/v2
models = "./internal/controllers/http/shared/view"
controllers = $(models),$(reports_http),$(draft_http)

generateDocs:
	swag i  --md ./docs/md -d ./cmd/http,$(controllers) --parseDependency --parseInternal --parseDepth 1