package servers_test

import (
	"context"
	"net/http/httptest"
	"testing"

	"github.com/RTUITLab/ITLab-Reports/config"
	reportdomain "github.com/RTUITLab/ITLab-Reports/domain/report"
	"github.com/RTUITLab/ITLab-Reports/pkg/adapters/toapprovereportsidgetter"
	"github.com/RTUITLab/ITLab-Reports/pkg/filter"
	"github.com/RTUITLab/ITLab-Reports/pkg/optional"
	"github.com/RTUITLab/ITLab-Reports/pkg/ordertype"
	"github.com/RTUITLab/ITLab-Reports/service/reports/reportservice"
	"github.com/RTUITLab/ITLab-Reports/service/salary"
	"github.com/RTUITLab/ITLab-Reports/transport/middlewares"
	mcontext "github.com/RTUITLab/ITLab-Reports/transport/middlewares/context"
	"github.com/RTUITLab/ITLab-Reports/transport/report"
	"github.com/RTUITLab/ITLab-Reports/transport/report/http/dto/v2"
	"github.com/RTUITLab/ITLab-Reports/transport/report/http/endpoints/v2"
	"github.com/RTUITLab/ITLab-Reports/transport/report/http/servers/v2"
	"github.com/golang-jwt/jwt/v4"
	"github.com/stretchr/testify/require"
)

func TestFunc_Server(t *testing.T) {
	cfg := config.GetConfigFrom(
		"./../../../../../.env",
	)

	s, err := reportservice.New(
		reportservice.WithMongoRepository(cfg.MongoDB.TestURI),
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
			servers.WithApprovedreportsIdGetter(
				toapprovereportsidgetter.ToApproveReportsIdGetter(
					salary.NewTestModeSalaryService(),
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

	// user2Token, err := jwt.NewWithClaims(
	// 	jwt.SigningMethodHS256, jwt.MapClaims{
	// 		"roles": []any{
	// 			"reports.user",
	// 		},
	// 		"sub": "user_2_id",
	// 		"aud": "itlab",
	// 		"nbf": 1650536405,
	// 		"exp": 1650540005,
	// 	},
	// ).SignedString([]byte("test"))
	// require.NoError(t, err)

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
						&dto.GetReportsReq{},
					)
					require.Error(t, err)
					require.Nil(t, resp)
				},
			)

			t.Run(
				"Success",
				func(t *testing.T) {
					t.Run(
						"UserGetReports",
						func(t *testing.T) {
							mctx := mcontext.New(context.Background())
							mctx.SetToken(user1Token)

							req := &dto.GetReportsReq{
								Query: dto.GetReportsQuery{
									Params: &reportdomain.GetReportsParams{
										Filter: &reportdomain.GetReportsFilter{
											GetReportsFilterFieldsWithOrAnd: reportdomain.GetReportsFilterFieldsWithOrAnd{
												GetReportsFilterFields: reportdomain.GetReportsFilterFields{
													Implementer: &filter.FilterField[string]{
														Value:     "some_id",
														Operation: filter.EQ,
													},
													Reporter: &filter.FilterField[string]{
														Value:     "some_id",
														Operation: filter.EQ,
													},
												},
											},
										},
									},
								},
							}

							_, err := httpEnds.GetReports(
								mctx,
								req,
							)
							require.NoError(t, err)

							require.Nil(t, req.Query.Params.Filter.GetReportsFilterFieldsWithOrAnd.GetReportsFilterFields.Implementer)
							require.Nil(t, req.Query.Params.Filter.GetReportsFilterFieldsWithOrAnd.GetReportsFilterFields.Reporter)

						},
					)

					t.Run(
						"Approved",
						func(t *testing.T) {
							t.Run(
								"Admin",
								func(t *testing.T) {
									mctx := mcontext.New(context.Background())
									mctx.SetToken(adminToken)

									req := &dto.GetReportsReq{
										Query: dto.GetReportsQuery{
											ApprovedState: dto.Approved,
											Params: &reportdomain.GetReportsParams{
												Filter: &reportdomain.GetReportsFilter{

												},
											},
										},
									}

									_, err := httpEnds.GetReports(
										mctx,
										req,
									)
									require.NoError(t, err)
								},
							)

							t.Run(
								"User",
								func(t *testing.T) {
									mctx := mcontext.New(context.Background())
									mctx.SetToken(user1Token)

									req := &dto.GetReportsReq{
										Query: dto.GetReportsQuery{
											ApprovedState: dto.Approved,
											Params: &reportdomain.GetReportsParams{
												Filter: &reportdomain.GetReportsFilter{
													
												},
											},
										},
									}

									_, err := httpEnds.GetReports(
										mctx,
										req,
									)
									require.NoError(t, err)
								},
							)
						},
					)

					t.Run(
						"NotApproved",
						func(t *testing.T) {
							t.Run(
								"Admin",
								func(t *testing.T) {
									mctx := mcontext.New(context.Background())
									mctx.SetToken(adminToken)

									req := &dto.GetReportsReq{
										Query: dto.GetReportsQuery{
											ApprovedState: dto.NotApproved,
											Params: &reportdomain.GetReportsParams{
												Filter: &reportdomain.GetReportsFilter{
													
												},
											},
										},
									}

									_, err := httpEnds.GetReports(
										mctx,
										req,
									)
									require.NoError(t, err)
								},
							)

							t.Run(
								"User",
								func(t *testing.T) {
									mctx := mcontext.New(context.Background())
									mctx.SetToken(user1Token)

									req := &dto.GetReportsReq{
										Query: dto.GetReportsQuery{
											ApprovedState: dto.NotApproved,
											Params: &reportdomain.GetReportsParams{
												Filter: &reportdomain.GetReportsFilter{
													
												},
											},
										},
									}

									_, err := httpEnds.GetReports(
										mctx,
										req,
									)
									require.NoError(t, err)
								},
							)
						},
					)
				},
			)
		},
	)
}

func TestFunc_DTO(t *testing.T) {
	t.Run(
		"GetReportsReqDecode",
		func(t *testing.T) {
			httpReq := httptest.NewRequest(
				"GET",
				"/v2/reports/reports?offset=12&offset=14&limit=10&limit=15&dateBegin=2019-10-12T07:20:50Z&dateEnd=2019-10-12T07:20:50Z&sortBy=name:asc&sortBy=date:asc&match=name:dan&match=date:2019-10-12T07:20:50.30Z&match=assignees.implementer:id_1&match=assignees.reporter:id_2",
				nil,
			)

			req, err := dto.DecodeGetReportsReq(
				context.Background(),
				httpReq,
			)
			require.NoError(t, err)	
			expect := &dto.GetReportsReq{
				Query: dto.GetReportsQuery{
					Params: &reportdomain.GetReportsParams{
						Limit:  *optional.NewOptional[int64](10),
						Offset: *optional.NewOptional[int64](12),
						Filter: &reportdomain.GetReportsFilter{
							GetReportsSort: reportdomain.GetReportsSort{
								NameSort: *optional.NewOptional[ordertype.OrderType](ordertype.ASC),
								DateSort: *optional.NewOptional[ordertype.OrderType](ordertype.ASC),
							},
							GetReportsFilterFieldsWithOrAnd: reportdomain.GetReportsFilterFieldsWithOrAnd{
								GetReportsFilterFields: reportdomain.GetReportsFilterFields{
									Name: &filter.FilterField[string]{
										Value:     "dan",
										Operation: filter.LIKE,
									},
									Date: &filter.FilterField[string]{
										Value:     "2019-10-12T07:20:50.3Z",
										Operation: filter.EQ,
									},
									Implementer: &filter.FilterField[string]{
										Value:     "id_1",
										Operation: filter.EQ,
									},
									Reporter: &filter.FilterField[string]{
										Value:     "id_2",
										Operation: filter.EQ,
									},
								},
								And: []*reportdomain.GetReportsFilterFieldsWithOrAnd{
									{
										GetReportsFilterFields: reportdomain.GetReportsFilterFields{
											Date: &filter.FilterField[string]{
												Value:     "2019-10-12T07:20:50Z",
												Operation: filter.GTE,
											},
										},
									},
									{
										GetReportsFilterFields: reportdomain.GetReportsFilterFields{
											Date: &filter.FilterField[string]{
												Value:     "2019-10-12T07:20:50Z",
												Operation: filter.LTE,
											},
										},
									},
								},
							},
						},
					},
				},
			}

			require.Equal(
				t,
				expect.Query.Params,
				req.Query.Params,
				"Failed to decode GetReportsReq",
			)
		},
	)
}
