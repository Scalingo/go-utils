package oteltest

import (
	"github.com/golang/mock/gomock"
	otelsdk "go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/metric"

	"github.com/Scalingo/go-utils/otel/otelmock"
)

type MockMeterProviderWrapper struct {
	metric.MeterProvider

	meterProvider *otelmock.MockMeterProvider
}

func InitMockMeterProvider(ctrl *gomock.Controller) *otelmock.MockMeterProvider {
	meterProvider := otelmock.NewMockMeterProvider(ctrl)

	mockMeterProvider := &MockMeterProviderWrapper{
		meterProvider: meterProvider,
	}

	// Ensure OpenTelemetry uses the mocked meter provider
	otelsdk.SetMeterProvider(mockMeterProvider)

	return meterProvider
}

func (m *MockMeterProviderWrapper) Meter(name string, options ...metric.MeterOption) metric.Meter {
	rawMeter := m.meterProvider.Meter(name, options...)

	return &MockMeterWrapper{
		meter: rawMeter,
	}
}
