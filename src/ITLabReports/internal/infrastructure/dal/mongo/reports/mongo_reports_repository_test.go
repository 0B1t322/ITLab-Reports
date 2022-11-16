package reports_test

import (
	"context"
	"testing"

	"github.com/RTUITLab/ITLab-Reports/internal/config"
	dal "github.com/RTUITLab/ITLab-Reports/internal/infrastructure/dal/mongo"
	. "github.com/RTUITLab/ITLab-Reports/internal/infrastructure/dal/mongo/reports"
	"github.com/RTUITLab/ITLab-Reports/internal/models/aggregate"
	"github.com/stretchr/testify/require"
)

func TestFunc_MongoReportsRepository(t *testing.T) {
	config.InitGlobalConfig()
	db, err := dal.ConnectToMongoDB(config.GlobalConfig.MongoDB.TestURI)
	require.NoError(t, err)

	repo := NewMongoReportsRepository(db)

	t.Run(
		"Create",
		func(t *testing.T) {
			t.Run(
				"Success",
				func(t *testing.T) {
					report, err := aggregate.NewReport(
						"test",
						"test",
						"test",
						"test",
					)
					require.NoError(t, err)

					err = repo.CreateReport(context.Background(), &report)
					require.NoError(t, err)

					require.NotEmpty(t, report.ID)
				},
			)
		},
	)

	t.Run(
		"Get",
		func(t *testing.T) {
			t.Run(
				"Success",
				func(t *testing.T) {
					report, err := aggregate.NewReport(
						"test",
						"test",
						"test",
						"test",
					)
					require.NoError(t, err)

					err = repo.CreateReport(context.Background(), &report)
					require.NoError(t, err)

					report, err = repo.GetReport(context.Background(), report.ID)
					require.NoError(t, err)

					require.NotEmpty(t, report.ID)
				},
			)
		},
	)
}
