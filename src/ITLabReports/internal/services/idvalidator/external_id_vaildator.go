package idvalidator

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/RTUITLab/ITLab-Reports/pkg/errors"
)

type externalRestIdValidator struct {
	baseUrl string

	client *http.Client
}

// ExternalRestIDValidator create id validator that validate by external rest interface
// 
// client param can be nil
func ExternalRestIDValidator(
	baseUrl string,
	// Can be nil
	client *http.Client,
) IdsValidator {
	v := &externalRestIdValidator{
		baseUrl: baseUrl,
	}

	if client == nil {
		v.client = &http.Client{}
	} else {
		v.client = client
	}

	return v
}

func (s *externalRestIdValidator) ValidateIds(
	ctx context.Context,
	token string,
	ids []string,
) (error) {
	body := &bytes.Buffer{}
	{
		if err := json.NewEncoder(body).Encode(ids); err != nil {
			return err
		}
	}

	req, err := http.NewRequest(
		http.MethodPost,
		s.baseUrl + "/api/User/checkUserIds",
		body,
	)
	if err != nil {
		return err
	}

	req.Header.Add("Authorization", token)
	req.Header.Add("Content-Type", "application/json")
	
	resp, err := s.client.Do(
		req,
	)
	if err != nil {
		return err
	}

	if resp.StatusCode == http.StatusOK {
		return nil
	} else if resp.StatusCode != http.StatusBadRequest {
		return fmt.Errorf("Unexpected responce status code: %v", resp.StatusCode)
	}

	var invalidIds []string
	{
		if err := json.NewDecoder(resp.Body).Decode(&invalidIds); err != nil {
			return err
		}
	}

	return errors.Wrap(fmt.Errorf(invalidIds[0]), ErrIdInvalid)
}