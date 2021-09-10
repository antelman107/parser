package mocks

import (
	"github.com/antelman107/parser/app"
	"github.com/stretchr/testify/mock"
)

type CsvAppBuilderMock struct {
	mock.Mock
}

func NewCsvAppBuilderMock() *CsvAppBuilderMock {
	return &CsvAppBuilderMock{}
}

func (r *CsvAppBuilderMock) Build(k uint32, cfg app.CsvParserAppConfig) app.CsvApp {
	ret := r.Called(k, cfg)
	return ret.Get(0).(app.CsvApp)
}

type CsvAppMock struct {
	mock.Mock
}

func NewCsvAppMock() *CsvAppMock {
	return &CsvAppMock{}
}

func (r *CsvAppMock) Run() error {
	return r.Called().Error(0)
}
