package app

import (
	"github.com/RTUITLab/ITLab-Reports/transport/draft/http/endpoints/v1"
	"github.com/RTUITLab/ITLab-Reports/transport/report/http/servers/v1"
)

func ToDraftService(e DraftEndpoints) servers.DraftService {
	return endpoints.ToDraftService(e)
}