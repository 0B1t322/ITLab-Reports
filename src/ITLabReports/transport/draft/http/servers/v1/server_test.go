package servers_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/RTUITLab/ITLab-Reports/config"
	"github.com/RTUITLab/ITLab-Reports/pkg/adapters/toidchecker"
	"github.com/RTUITLab/ITLab-Reports/pkg/errors"
	"github.com/RTUITLab/ITLab-Reports/service/idvalidator"
	"github.com/RTUITLab/ITLab-Reports/service/reports/reportservice"
	"github.com/RTUITLab/ITLab-Reports/transport/draft/http/dto/v1"
	"github.com/RTUITLab/ITLab-Reports/transport/draft/http/endpoints/v1"
	derr "github.com/RTUITLab/ITLab-Reports/transport/draft/http/errors"
	"github.com/RTUITLab/ITLab-Reports/transport/draft/http/servers/v1"
	"github.com/RTUITLab/ITLab-Reports/transport/middlewares"
	mcontext "github.com/RTUITLab/ITLab-Reports/transport/middlewares/context"
	"github.com/RTUITLab/ITLab-Reports/transport/report"
	"github.com/golang-jwt/jwt/v4"
	"github.com/stretchr/testify/require"
)

type IdValidatorWithSomeFail struct{}
func(*IdValidatorWithSomeFail) ValidateIds(ctx context.Context, token string, ids []string) error {
	for _, id := range ids {
		if id == "failed" {
			return fmt.Errorf("Failed error")
		} else if id == "invalid_id" {
			return errors.Wrap(fmt.Errorf("invalid_id"), idvalidator.ErrIdInvalid)
		}
	}
	return nil
}

func TestFunc_Server(t *testing.T) {
	cfg := config.GetConfigFrom(
		"./../../../../../.env",
	)

	s, err := reportservice.New(
		reportservice.WithMongoRepositoryAndCollectionName(
			cfg.MongoDB.TestURI, 
			"draft",
		),
	)
	require.NoError(t, err, err)

	ends := report.MakeEndpoints(s)

	httpEnds := endpoints.NewEndpoints(ends)

	httpEnds = servers.BuildMiddlewares(
		httpEnds,
		servers.MergeServerOptions(
			servers.WithAuther(
				middlewares.NewTestAuth(
					middlewares.WithAdminRole("reports.admin"),
					middlewares.WithUserRole("reports.user"),
					middlewares.WithSuperAdminRole("admin"),
					middlewares.WithRoleClaim("roles"),
				),
			),
			servers.WithIdChecker(
				toidchecker.ToIdChecker(
					idvalidator.New(
						&IdValidatorWithSomeFail{},
					),
				),
			),
		),
	)

	user1Token, err := jwt.NewWithClaims(
		jwt.SigningMethodHS256, jwt.MapClaims{
			"roles": []any{
				"reports.user",
			},
			"sub": "user_1_id",
			"aud": "itlab",
			"nbf": 1650536405,
  			"exp": 1650540005,	
		},
	).SignedString([]byte("test"))
	require.NoError(t, err)

	user2Token, err := jwt.NewWithClaims(
		jwt.SigningMethodHS256, jwt.MapClaims{
			"roles": []any{
				"reports.user",
			},
			"sub": "user_2_id",
			"aud": "itlab",
			"nbf": 1650536405,
  			"exp": 1650540005,	
		},
	).SignedString([]byte("test"))
	require.NoError(t, err)

	// adminToken, err := jwt.NewWithClaims(
	// 	jwt.SigningMethodHS256, jwt.MapClaims{
	// 		"roles": []any{
	// 			"reports.admin",
	// 		},
	// 		"sub": "admin_1_id",
	// 		"aud": "itlab",
	// 		"nbf": 1650536405,
  	// 		"exp": 1650540005,	
	// 	},
	// ).SignedString([]byte("test"))
	// require.NoError(t, err)

	notAuthUser, err := jwt.NewWithClaims(
		jwt.SigningMethodHS256, jwt.MapClaims{
			"roles": []any{
				"service.user",
			},
			"sub": "user_2_id",
			"aud": "itlab",
			"nbf": 1650536405,
  			"exp": 1650540005,	
		},
	).SignedString([]byte("test"))
	require.NoError(t, err)

	t.Run(
		"CreateDraft",
		func(t *testing.T) {
			t.Run(
				"Failure",
				func(t *testing.T) {
					t.Run(
						"Don't find role",
						func(t *testing.T) {
							mctx := mcontext.New(context.Background())
							mctx.SetToken(notAuthUser)

							_, err := httpEnds.CreateDraft(
								mctx,
								&dto.CreateDraftReq{
									Name: "some_name",
									Text: "some_text",
									Implementor: "implementer",
									Reporter: "reporter",
								},
							)
							require.ErrorIs(t, err, middlewares.RoleNotFound)
						},
					)

					t.Run(
						"ValidationError",
						func(t *testing.T) {
							t.Run(
								"EmptyFields",
								func(t *testing.T) {
									mctx := mcontext.New(context.Background())
									mctx.SetToken(user1Token)

									_, err := httpEnds.CreateDraft(
										mctx,
										&dto.CreateDraftReq{
											Name: "",
											Text: "",
											Implementor: "implementer",
											Reporter: "reporter",
										},
									)
									require.Condition(
										t,
										func() (success bool) {
											return errors.Is(err, derr.DraftValidationError)
										},
									)
								},
							)

							t.Run(
								"InvalidId",
								func(t *testing.T) {
									t.Run(
										"Implementer",
										func(t *testing.T) {
											mctx := mcontext.New(context.Background())
											mctx.SetToken(user1Token)

											_, err := httpEnds.CreateDraft(
												mctx,
												&dto.CreateDraftReq{
													Name: "some_name",
													Text: "some_text",
													Implementor: "invalid_id",
													Reporter: "some",
												},
											)
											require.Condition(
												t,
												func() (success bool) {
													return errors.Is(err, middlewares.ErrIncorectId)
												},
											)
										},
									)
								},
							)

							t.Run(
								"Failed to validate id",
								func(t *testing.T) {
									mctx := mcontext.New(context.Background())
									mctx.SetToken(user1Token)

									_, err := httpEnds.CreateDraft(
										mctx,
										&dto.CreateDraftReq{
											Name: "some_name",
											Text: "some_text",
											Implementor: "failed",
											Reporter: "some",
										},
									)
									require.Condition(
										t,
										func() (success bool) {
											return errors.Is(err, middlewares.ErrFaieldToValidateId)
										},
									)
								},
							)
						},
					)
				},
			)

			t.Run(
				"Success",
				func(t *testing.T) {
					t.Run(
						"CreateDraft",
						func(t *testing.T) {
							mctx := mcontext.New(context.Background())
							mctx.SetToken(user1Token)

							resp, err := httpEnds.CreateDraft(
								mctx,
								&dto.CreateDraftReq{
									Name: "some_name",
									Text: "somte_text",
									Implementor: "implementer",
									Reporter: "reporter",
								},
							)
							require.NoError(t, err)

							require.NotNil(t, resp)

							require.Equal(
								t,
								resp.Assignes.Reporter,
								"user_1_id",
							)
						},
					)
				},
			)
		},
	)

	t.Run(
		"GetDraft",
		func(t *testing.T) {
			t.Run(
				"Failure",
				func(t *testing.T) {
					t.Run(
						"Role not found",
						func(t *testing.T) {
							mctx := mcontext.New(context.Background())
							mctx.SetToken(notAuthUser)

							_, err := httpEnds.GetDraft(
								mctx,
								&dto.GetDraftReq{
									ID: "some_id",
								},
							)
							require.ErrorIs(t, err, middlewares.RoleNotFound)
						},
					)

					t.Run(
						"IDInvalidOrNotFound",
						func(t *testing.T) {
							mctx := mcontext.New(context.Background())
							mctx.SetToken(user1Token)

							_, err :=httpEnds.GetDraft(
								mctx,
								&dto.GetDraftReq{
									ID: "some_id",
								},
							)
							require.Condition(
								t,
								func() (success bool) {
									return err == derr.DraftIDIsInvalid || err == derr.DraftNotFound
								},
							)
						},
					)

					t.Run(
						"UserTryToGetNotTheirDraft",
						func(t *testing.T) {
							ctx := mcontext.New(context.Background())
							ctx.SetToken(user1Token)

							resp, err := httpEnds.CreateDraft(
								ctx,
								&dto.CreateDraftReq{
									Name: "some_name",
									Text: "somte_text",
									Implementor: "implementer",
									Reporter: "reporter",
								},
							)
							require.NoError(t, err)

							id := resp.ID

							mctx := mcontext.New(context.Background())
							mctx.SetToken(user2Token)

							get, err := httpEnds.GetDraft(
								mctx,
								&dto.GetDraftReq{
									ID: id,
								},
							)
							require.ErrorIs(t, err, middlewares.YouAreNotOwner)

							require.Nil(t, get)
						},
					)
				},
			)

			t.Run(
				"Success",
				func(t *testing.T) {
					mctx := mcontext.New(context.Background())
					mctx.SetToken(user1Token)

					resp, err := httpEnds.CreateDraft(
						mctx,
						&dto.CreateDraftReq{
							Name: "some_name",
							Text: "somte_text",
							Implementor: "implementer",
							Reporter: "reporter",
						},
					)
					require.NoError(t, err)

					id := resp.ID

					get, err := httpEnds.GetDraft(
						mctx,
						&dto.GetDraftReq{
							ID: id,
						},
					)
					require.NoError(t, err)

					require.Equal(
						t,
						&dto.GetDraftResp{
							ID: resp.ID,
							Name: resp.Name,
							Date: resp.Date,
							Text: resp.Text,
							Assignees: resp.Assignes,
						},
						get,
					)
				},
			)
		},
	)

	t.Run(
		"DeleteDraft",
		func(t *testing.T) {
			t.Run(
				"Failure",
				func(t *testing.T) {
					t.Run(
						"Role not found",
						func(t *testing.T) {
							mctx := mcontext.New(context.Background())
							mctx.SetToken(notAuthUser)

							_, err := httpEnds.DeleteDraft(
								mctx,
								&dto.DeleteDraftReq{
									ID: "some_id",
								},
							)
							require.ErrorIs(t, err, middlewares.RoleNotFound)
						},
					)

					t.Run(
						"IDInvalidOrNotFound",
						func(t *testing.T) {
							mctx := mcontext.New(context.Background())
							mctx.SetToken(user1Token)

							_, err :=httpEnds.DeleteDraft(
								mctx,
								&dto.DeleteDraftReq{
									ID: "some_id",
								},
							)
							require.Condition(
								t,
								func() (success bool) {
									return err == derr.DraftIDIsInvalid || err == derr.DraftNotFound
								},
							)
						},
					)

					t.Run(
						"UserTryToDeleteNotTheirDraft",
						func(t *testing.T) {
							ctx := mcontext.New(context.Background())
							ctx.SetToken(user1Token)

							resp, err := httpEnds.CreateDraft(
								ctx,
								&dto.CreateDraftReq{
									Name: "some_name",
									Text: "somte_text",
									Implementor: "implementer",
									Reporter: "reporter",
								},
							)
							require.NoError(t, err)

							id := resp.ID

							mctx := mcontext.New(context.Background())
							mctx.SetToken(user2Token)

							get, err := httpEnds.DeleteDraft(
								mctx,
								&dto.DeleteDraftReq{
									ID: id,
								},
							)
							require.ErrorIs(t, err, middlewares.YouAreNotOwner)

							require.Nil(t, get)
						},
					)
				},
			)

			t.Run(
				"Success",
				func(t *testing.T) {
					t.Run(
						"DeleteOwnDraft",
						func(t *testing.T) {
							ctx := mcontext.New(context.Background())
							ctx.SetToken(user1Token)

							resp, err := httpEnds.CreateDraft(
								ctx,
								&dto.CreateDraftReq{
									Name: "some_name",
									Text: "somte_text",
									Implementor: "implementer",
									Reporter: "reporter",
								},
							)
							require.NoError(t, err)
							get, err := httpEnds.DeleteDraft(
								ctx,
								&dto.DeleteDraftReq{
									ID: resp.ID,
								},
							)
							require.NoError(t, err)

							require.NotNil(t, get)
						},
					)
				},
			)
		},
	)

	t.Run(
		"UpdateReport",
		func(t *testing.T) {
			t.Run(
				"Failure",
				func(t *testing.T) {
					t.Run(
						"Role don't find",
						func(t *testing.T) {
							mctx := mcontext.New(context.Background())
							mctx.SetToken(notAuthUser)

							_, err := httpEnds.UpdateDraft(
								mctx,
								&dto.UpdateDraftReq{
								},
							)
							require.ErrorIs(t, err, middlewares.RoleNotFound)
						},
					)

					t.Run(
						"TryToUpdateNotYourReport",
						func(t *testing.T) {
							ctx := mcontext.New(context.Background())
							ctx.SetToken(user1Token)
		
							resp, err := httpEnds.CreateDraft(
								ctx,
								&dto.CreateDraftReq{
									Name: "some_name",
									Text: "somte_text",
									Implementor: "implementer",
									Reporter: "reporter",
								},
							)
							require.NoError(t, err)
		
							id := resp.ID
		
							mctx := mcontext.New(context.Background())
							mctx.SetToken(user2Token)
		
							get, err := httpEnds.UpdateDraft(
								mctx,
								&dto.UpdateDraftReq{
									ID: id,
								},
							)
							require.ErrorIs(t, err, middlewares.YouAreNotOwner)
		
							require.Nil(t, get)
						},
					)

					t.Run(
						"InvalidId",
						func(t *testing.T) {
							ctx := mcontext.New(context.Background())
							ctx.SetToken(user1Token)
		
							resp, err := httpEnds.CreateDraft(
								ctx,
								&dto.CreateDraftReq{
									Name: "some_name",
									Text: "somte_text",
									Implementor: "implementer",
									Reporter: "reporter",
								},
							)
							require.NoError(t, err)
							

							_, err = httpEnds.UpdateDraft(
								ctx,
								&dto.UpdateDraftReq{
									ID: resp.ID,
									Implementer: "invalid_id",
								},
							)
							require.Condition(
								t,
								func() (success bool) {
									return errors.Is(err, middlewares.ErrIncorectId)
								},
							)
						},
					)
				},
			)

			t.Run(
				"Success",
				func(t *testing.T) {
					t.Run(
						"UserUpdateTheirDraft",
						func(t *testing.T) {
							ctx := mcontext.New(context.Background())
							ctx.SetToken(user1Token)
		
							resp, err := httpEnds.CreateDraft(
								ctx,
								&dto.CreateDraftReq{
									Name: "some_name",
									Text: "somte_text",
									Implementor: "implementer",
									Reporter: "reporter",
								},
							)
							require.NoError(t, err)
							

							get, err := httpEnds.UpdateDraft(
								ctx,
								&dto.UpdateDraftReq{
									ID: resp.ID,
								},
							)
							require.NoError(t, err)
		
							require.NotNil(t, get)
						},
					)
				},
			)

		},
	)

	t.Run(
		"GetDrafts",
		func(t *testing.T) {
			t.Run(
				"Failure",
				func(t *testing.T) {
					t.Run(
						"Role don't find",
						func(t *testing.T) {
							mctx := mcontext.New(context.Background())
							mctx.SetToken(notAuthUser)
							req := &dto.GetDraftsReq{}

							_, err := httpEnds.GetDrafts(
								mctx,
								req,
							)
							require.ErrorIs(t, err, middlewares.RoleNotFound)
						},
					)
				},
			)

			t.Run(
				"Success",
				func(t *testing.T) {
					mctx := mcontext.New(context.Background())
					mctx.SetToken(user1Token)
					req := &dto.GetDraftsReq{}

					_, err := httpEnds.GetDrafts(
						mctx,
						req,
					)
					require.NoError(t, err)

					require.Equal(
						t,
						"user_1_id",
						req.UserID,
					)
				},
			)
		},
	)
}