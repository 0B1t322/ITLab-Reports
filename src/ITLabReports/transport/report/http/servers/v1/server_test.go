package servers_test

import (
	"context"
	"testing"

	"github.com/RTUITLab/ITLab-Reports/config"
	"github.com/RTUITLab/ITLab-Reports/pkg/errors"
	"github.com/RTUITLab/ITLab-Reports/service/reports/reportservice"
	"github.com/RTUITLab/ITLab-Reports/transport/middlewares"
	mcontext "github.com/RTUITLab/ITLab-Reports/transport/middlewares/context"
	"github.com/RTUITLab/ITLab-Reports/transport/report"
	"github.com/RTUITLab/ITLab-Reports/transport/report/http/dto/v1"
	"github.com/RTUITLab/ITLab-Reports/transport/report/http/endpoints/v1"
	serr "github.com/RTUITLab/ITLab-Reports/transport/report/http/errors"
	"github.com/RTUITLab/ITLab-Reports/transport/report/http/servers/v1"
	"github.com/golang-jwt/jwt/v4"

	"github.com/stretchr/testify/require"
)

func TestFunc_Func(t *testing.T) {
	cfg := config.GetConfigFrom(
		"./../../../../../.env",
	)

	s, err := reportservice.New(
		reportservice.WithMongoRepository(cfg.MongoDB.TestURI),
	)
	require.NoError(t, err)

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

	adminToken, err := jwt.NewWithClaims(
		jwt.SigningMethodHS256, jwt.MapClaims{
			"roles": []any{
				"reports.admin",
			},
			"sub": "admin_1_id",
			"aud": "itlab",
			"nbf": 1650536405,
  			"exp": 1650540005,	
		},
	).SignedString([]byte("test"))
	require.NoError(t, err)

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


	// CreateReport
	creartReportFunc := func(mctx mcontext.MiddlewareContext) (*dto.CreateReportResp, error) {
		name := "report_rest"
		resp, err := httpEnds.CreateReport(
			mctx,
			&dto.CreateReportReq{
				Name: &name,
				Text: "test_test",
			},
		)
		if err != nil {
			return nil, err
		}

		return resp, nil
	}

	t.Run(
		"GetReport",
		func(t *testing.T) {
			t.Run(
				"Failure",
				func(t *testing.T) {
					t.Run(
						"Don't find role",
						func(t *testing.T) {
							mctx := mcontext.New(context.Background())
							mctx.SetToken(notAuthUser)
		
							_, err := httpEnds.GetReport(
								mctx,
								&dto.GetReportReq{
									ID: "some_id",
								},
							)
							require.Error(t, err)
						},
					)

					t.Run(
						"UserCantGetReportOfAnotherUser",
						func(t *testing.T) {
							mctxUser := mcontext.New(context.Background())
							mctxUser.SetToken(user1Token)
							created, err := creartReportFunc(mctxUser)
							require.NoError(t, err)

							id := created.ID
							mctxAnotherUser := mcontext.New(context.Background())
							mctxAnotherUser.SetToken(user2Token)

							resp, err := httpEnds.GetReport(
								mctxAnotherUser,
								&dto.GetReportReq{
									ID: id,
								},
							)
							require.Error(t, err)

							require.Nil(t, resp)
						},
					)
				},
			)

			t.Run(
				"Success",
				func(t *testing.T) {
					t.Run(
						"GetCreatedReport",
						func(t *testing.T) {
							mctx := mcontext.New(context.Background())
							mctx.SetToken(user1Token)
							created, err := creartReportFunc(mctx)
							require.NoError(t, err)

							id := created.ID

							get, err := httpEnds.GetReport(
								mctx,
								&dto.GetReportReq{
									ID: id,
								},
							)
							require.NoError(t, err)

							require.Equal(
								t,
								dto.GetReportResp(*created),
								*get,
							)
						},
					)

					t.Run(
						"AdminGetCreatedReportFromUser",
						func(t *testing.T) {
							mctxUser := mcontext.New(context.Background())
							mctxUser.SetToken(user1Token)
							created, err := creartReportFunc(mctxUser)
							require.NoError(t, err)

							id := created.ID
							mctxAdmin := mcontext.New(context.Background())
							mctxAdmin.SetToken(adminToken)

							get, err := httpEnds.GetReport(
								mctxAdmin,
								&dto.GetReportReq{
									ID: id,
								},
							)
							require.NoError(t, err)

							require.Equal(
								t,
								dto.GetReportResp(*created),
								*get,
							)
						},
					)
				},
			)
		},
	)

	t.Run(
		"CreateReport",
		func(t *testing.T) {
			t.Run(
				"Failure",
				func(t *testing.T) {
					t.Run(
						"DontFindRole",
						func(t *testing.T) {
							mcxt := mcontext.New(context.Background())
							mcxt.SetToken(notAuthUser)

							name := "some_name"
							resp, err := httpEnds.CreateReport(
								mcxt,
								&dto.CreateReportReq{
									Name: &name,
									Text: "some_text",
									Implementor: "",
									Reporter: "",
								},
							)
							require.Error(t, err)

							require.Nil(t, resp)
						},
					)
				},
			)

			t.Run(
				"Success",
				func(t *testing.T) {
					t.Run(
						"CreateReportForCurrentUser",
						func(t *testing.T) {
							t.Run(
								"NameAndText",
								func(t *testing.T) {
									mctx := mcontext.New(context.Background())
									mctx.SetToken(user1Token)

									name := "some_name"
									resp, err := httpEnds.CreateReport(
										mctx,
										&dto.CreateReportReq{
											Name: &name,
											Text: "some_text",
										},
									)
									require.NoError(t, err)
									userId, _ := mctx.GetUserID()

									require.Equal(
										t,
										dto.CreateReportResp(
											dto.GetReportResp{
												ID: resp.ID,
												Date: resp.Date,
												Name: "some_name",
												Text: "some_text",
												Assignes: dto.GetAssigneesResp{
													Reporter: userId,
													Implementer: userId,
												},
											},
										),
										*resp,
									)
								},
							)

							t.Run(
								"NameInText",
								func(t *testing.T) {
									mctx := mcontext.New(context.Background())
									mctx.SetToken(user1Token)

									resp, err := httpEnds.CreateReport(
										mctx,
										&dto.CreateReportReq{
											Text: "some_name@\n\t\n@some_text",
										},
									)
									require.NoError(t, err)
									userId, _ := mctx.GetUserID()

									require.Equal(
										t,
										dto.CreateReportResp(
											dto.GetReportResp{
												ID: resp.ID,
												Date: resp.Date,
												Name: "some_name",
												Text: "some_text",
												Assignes: dto.GetAssigneesResp{
													Reporter: userId,
													Implementer: userId,
												},
											},
										),
										*resp,
									)
								},
							)
						},
					)
				},
			)
		},
	)

	t.Run(
		"GetReportsForEmployee",
		func(t *testing.T) {
			t.Run(
				"Failure",
				func(t *testing.T) {
					t.Run(
						"Don't find role",
						func(t *testing.T) {
							mctx := mcontext.New(context.Background())
							mctx.SetToken(notAuthUser)
		
							resp, err := httpEnds.GetReportsForEmployee(
								mctx,
								&dto.GetReportsForEmployeeReq{
									Employee: "some_empl",
								},
							)
							require.Error(t, err)
							require.Nil(t, resp)
						},
					)

					t.Run(
						"Employee can't be empty",
						func(t *testing.T) {
							mctx := mcontext.New(context.Background())
							mctx.SetToken(user1Token)

							resp, err := httpEnds.GetReportsForEmployee(
								mctx,
								&dto.GetReportsForEmployeeReq{
								},
							)
							require.Condition(
								t,
								func() (success bool) {
									return errors.Is(err, serr.ValidationError)
								},
							)
							require.Nil(t, resp)
						},
					)

					t.Run(
						"UserGetNotTheirReports",
						func(t *testing.T) {
							mctxUser1 := mcontext.New(context.Background())
							mctxUser1.SetToken(user1Token)

							mctxUser2 := mcontext.New(context.Background())
							mctxUser2.SetUserID("user_id_2")

							iduser2, _ := mctxUser2.GetUserID()

							resp, err := httpEnds.GetReportsForEmployee(
								mctxUser1,
								&dto.GetReportsForEmployeeReq{
									Employee: iduser2,
								},
							)
							require.Error(t, err)
							require.Nil(t, resp)

						},
					)
				},
			)

			t.Run(
				"Success",
				func(t *testing.T) {
					t.Run(
						"UserGetTheirReports",
						func(t *testing.T) {
							mctx := mcontext.New(context.Background())
							mctx.SetToken(user1Token)

							_, err := creartReportFunc(
								mctx,
							)
							require.NoError(t, err)

							id, _ := mctx.GetUserID()

							_, err = httpEnds.GetReportsForEmployee(
								mctx,
								&dto.GetReportsForEmployeeReq{
									Employee: id,
								},
							)
							require.NoError(t, err)
						},
					)

					t.Run(
						"AdminGetUserReports",
						func(t *testing.T) {
							mctxAdmin := mcontext.New(context.Background())
							mctxAdmin.SetToken(adminToken)

							mctxUser2 := mcontext.New(context.Background())
							mctxUser2.SetUserID("user_id_2")

							iduser2, _ := mctxUser2.GetUserID()

							_, err := httpEnds.GetReportsForEmployee(
								mctxAdmin,
								&dto.GetReportsForEmployeeReq{
									Employee: iduser2,
								},
							)
							require.NoError(t, err)
						},
					)
				},
			)
		},
	)

	t.Run(
		"GetReports",
		func(t *testing.T) {
			t.Run(
				"Failure",
				func(t *testing.T) {
					mctx := mcontext.New(context.Background())
					mctx.SetToken(notAuthUser)

					resp, err := httpEnds.GetReports(
						mctx,
						&dto.GetReportsReq{
						},
					)
					require.Error(t, err)
					require.Nil(t, resp)
				},
			)

			t.Run(
				"Success",
				func(t *testing.T) {
					t.Run(
						"UserGetReportsWithTheir",
						func(t *testing.T) {
							mctx := mcontext.New(context.Background())
							mctx.SetToken(user1Token)

							id := "user_1_id"
							resp, err := httpEnds.GetReports(
								mctx,
								&dto.GetReportsReq{
									Implementer: id,
									Reporter: id,
								},
							)
							require.NoError(t, err)
							require.NotNil(t, resp)
						},
					)

					t.Run(
						"UserGetReportsWithNotTheir",
						func(t *testing.T) {
							mctx := mcontext.New(context.Background())
							mctx.SetToken(user1Token)

							id := "user_2_id"
							resp, err := httpEnds.GetReports(
								mctx,
								&dto.GetReportsReq{
									Implementer: id,
									Reporter: id,
								},
							)
							require.NoError(t, err)
							require.NotNil(t, resp)
						},
					)

					t.Run(
						"AdminGetReportsForUser",
						func(t *testing.T) {
							mctx := mcontext.New(context.Background())
							mctx.SetToken(adminToken)

							id := "user_2_id"
							resp, err := httpEnds.GetReports(
								mctx,
								&dto.GetReportsReq{
									Implementer: id,
									Reporter: id,
								},
							)
							require.NoError(t, err)
							require.NotNil(t, resp)
						},
					)
				},
			)
		},
	)
}