package clienttest

import (
	"github.com/ONSdigital/go-ns/log"
)

// ClientParams struct encapsulating the parameters of a single ReportError call
type ClientParams struct {
	InstanceID string
	ErrContext string
	Err        error
	Data       log.Data
}

// ReportErrorFunc returns a reportErrorFunc which returns the error provided when called.
func ReportErrorFunc(err error) func(instanceID string, errContext string, err error, data log.Data) error {
	return func(instanceID string, errContext string, err error, data log.Data) error {
		return err
	}
}

// NewReporterClientMock create a new ReporterClientMock
// reportErrorFunc enables you to customise the result of ReportError func to meet the needs of your test case
func NewReporterClientMock(reportErrorFunc func(string, string, error, log.Data) error) *ReporterClientMock {
	return &ReporterClientMock{
		reportErrorFunc: reportErrorFunc,
		params:          make([]ClientParams, 0),
	}
}

type ReporterClientMock struct {
	params          []ClientParams
	reportErrorFunc func(string, string, error, log.Data) error
}

func (m *ReporterClientMock) ReportError(instanceID string, errContext string, err error, data log.Data) error {
	m.params = append(m.params, ClientParams{
		InstanceID: instanceID,
		ErrContext: errContext,
		Err:        err,
		Data:       data,
	})
	return m.reportErrorFunc(instanceID, errContext, err, data)
}

// ReportErrorCalls return a slice of the ClientParams passed into each invocation of ReportError
func (m *ReporterClientMock) ReportErrorCalls() []ClientParams {
	return m.params
}
