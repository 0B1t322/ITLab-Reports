reports_http = ./transport/report/http/dto/v1,./transport/report/http/dto/v2,./transport/report/http/handlers
draft_http = ./transport/draft/http/dto/v1,./transport/draft/http/dto/v2,./transport/draft/http/handlers

generateDocs:
	swag i  --md ./docs/md -d ./cmd/main,$(reports_http),$(draft_http) --parseDependency --parseInternal --parseDepth 1