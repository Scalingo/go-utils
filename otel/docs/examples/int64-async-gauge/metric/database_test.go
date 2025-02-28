package metric

import (
	"testing"

	"github.com/golang/mock/gomock"

	"github.com/Scalingo/go-utils/otel/otelmock"
	"github.com/Scalingo/go-utils/otel/oteltest"
)

func TestRegisterDatabaseAsyncGauge(t *testing.T) {
	t.Run("RegisterDatabasAsyncGauge", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		meterProvider := oteltest.InitMockMeterProvider(ctrl)
		meter := otelmock.NewMockMeter(ctrl)
		int64ObservableGauge := otelmock.NewMockInt64ObservableGauge(ctrl)

		meterProvider.EXPECT().Meter(
			"database",
		).Return(meter).Times(1)

		meter.EXPECT().Int64ObservableGauge(
			"database_count", gomock.Any(),
		).Return(int64ObservableGauge, nil).Times(1)

		result := RegisterDatabaseAsyncGauge()
		if !result {
			t.Errorf("RegisterDatabaseAsyncGauge() failed")
		}
	})
}
