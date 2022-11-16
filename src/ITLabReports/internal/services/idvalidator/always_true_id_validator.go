package idvalidator

import "context"

// Use for test mode and tests
type alwaysTrueIdValidator struct {}

func AlwaysTrueIdValidator() IdsValidator {
	return &alwaysTrueIdValidator{}
}

func (*alwaysTrueIdValidator) ValidateIds(
	ctx context.Context,
	token string,
	ids []string,
) (error) {
	return nil
}