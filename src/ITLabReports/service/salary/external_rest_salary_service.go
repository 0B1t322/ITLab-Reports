package salary

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/RTUITLab/ITLab-Reports/pkg/optional"
)

type ExternalRESTSalaryService struct {
	baseUrl string
	client  *http.Client
}

func NewExternalRestSalaryService(
	baseUrl string,
	// Can be nil
	client *http.Client,
) SalaryService {
	s := &ExternalRESTSalaryService{
		baseUrl: baseUrl,
	}

	if client == nil {
		s.client = &http.Client{}
	} else {
		s.client = client
	}

	return s
}

func (e *ExternalRESTSalaryService) GetApprovedReportsIds(
	ctx context.Context,
	token string,
	userId optional.Optional[string],
) ([]string, error) {
	req, err := http.NewRequest(
		http.MethodGet,
		e.baseUrl + "/api/salary/v1/report/approved/reportIds",
		nil,
	)
	if err != nil {
		return nil, err
	}

	if userId.HasValue() {
		req.URL.Query().Add("userId", userId.MustGetValue())
	}

	req.Header.Add("Authorization", token)
	req.Header.Add("Content-Type", "application/json")

	resp, err := e.client.Do(req)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		bodyMessage := &bytes.Buffer{}
		bodyMessage.ReadFrom(resp.Body)
		return nil, fmt.Errorf("Unexpected responce status code: %v, message: %v", resp.StatusCode, bodyMessage.String())
	}

	var approvedIds []string
	{
		if err := json.NewDecoder(resp.Body).Decode(&approvedIds); err != nil {
			return nil, err
		}
	}

	return approvedIds, nil
}