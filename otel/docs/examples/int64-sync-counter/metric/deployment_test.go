package metric

import (
	"testing"

	"github.com/golang/mock/gomock"

	"github.com/Scalingo/go-utils/otel/otelmock"
	"github.com/Scalingo/go-utils/otel/oteltest"
)

func TestWriteDeploymentMetric(t *testing.T) {
	t.Run("WriteDeploymentMetric", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		meterProvider := oteltest.InitMockMeterProvider(ctrl)
		meter := otelmock.NewMockMeter(ctrl)
		int64counter := otelmock.NewMockInt64Counter(ctrl)

		meterProvider.EXPECT().Meter(
			"deployment",
		).Return(meter).Times(1)

		meter.EXPECT().Int64Counter(
			"deployment_count", gomock.Any(),
		).Return(int64counter, nil).Times(1)

		int64counter.EXPECT().Add(
			gomock.Any(), int64(10), gomock.Any(),
		).Times(1)

		int64counter.EXPECT().Add(
			gomock.Any(), int64(42), gomock.Any(),
		).Times(1)

		result := WriteDeploymentMetric()
		if !result {
			t.Errorf("WriteDeploymentMetric() failed")
		}
	})
}
