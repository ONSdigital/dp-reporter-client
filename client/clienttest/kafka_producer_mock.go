package clienttest

import "context"

// KafkaProducerMock a mock of the go-ns kafka.Producer
type KafkaProducerMock struct {
	output      chan []byte
	outputCalls int
	closeFunc   func(ctx context.Context) error
	closeCalls  []context.Context
}

// NewKafkaProducerMock convenience func for create a new KafkaProducerMock for your tests
func NewKafkaProducerMock(output chan []byte, closeFunc func(ctx context.Context) error) *KafkaProducerMock {
	return &KafkaProducerMock{
		output:      output,
		outputCalls: 0,
		closeFunc:   closeFunc,
	}
}

func (m *KafkaProducerMock) Output() chan []byte {
	m.outputCalls += 1
	return m.output
}

// OutputCalls return the number times Output() was invoked on this mock
func (m *KafkaProducerMock) OutputCalls() int {
	return m.outputCalls
}

func (m *KafkaProducerMock) Close(ctx context.Context) error {
	m.closeCalls = append(m.closeCalls, ctx)
	return m.closeFunc(ctx)
}

// CloseCalls returns a slice of the parameters passed into the Close function in the order they were invoked
func (m *KafkaProducerMock) CloseCalls() []context.Context {
	return m.closeCalls
}
