package jman_test

import (
	"strings"
	"testing"

	"github.com/akaswenwilk/jman"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockT struct {
	mock.Mock
	gotMsg  string
	wantMsg string
}

func (m *MockT) Fatalf(format string, args ...any) {
	m.Called(format)
	m.gotMsg = format
	panic(format) // panic to simulate Fatalf behavior
}

func (m *MockT) AssertExpectations(t *testing.T) {
	t.Helper()
	m.Mock.AssertExpectations(t)
	gotMsgLines := strings.Split(m.gotMsg, "\n")
	wantMsgLines := strings.Split(m.wantMsg, "\n")
	assert.Equal(t, len(wantMsgLines), len(gotMsgLines), "number of lines in error message should match")
	for _, wantMsgLine := range wantMsgLines {
		assert.Contains(t, m.gotMsg, wantMsgLine, "got message should contain expected line")
	}
}

func newMockT(expectedMsg string) *MockT {
	mockT := &MockT{wantMsg: expectedMsg}
	mockT.On("Fatalf", mock.Anything).Return()
	return mockT
}

func assertFatalf(t *testing.T, expectedMsg string, fn func(t jman.T)) {
	t.Helper()
	mt := newMockT(expectedMsg)
	defer mt.AssertExpectations(t)
	assert.Panics(t, func() { fn(mt) })
}
