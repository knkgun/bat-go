// Code generated by MockGen. DO NOT EDIT.
// Source: ./grant/datastore.go

// Package grant is a generated GoMock package.
package grant

import (
	wallet "github.com/brave-intl/bat-go/utils/wallet"
	v4 "github.com/golang-migrate/migrate/v4"
	gomock "github.com/golang/mock/gomock"
	sqlx "github.com/jmoiron/sqlx"
	reflect "reflect"
)

// MockDatastore is a mock of Datastore interface
type MockDatastore struct {
	ctrl     *gomock.Controller
	recorder *MockDatastoreMockRecorder
}

// MockDatastoreMockRecorder is the mock recorder for MockDatastore
type MockDatastoreMockRecorder struct {
	mock *MockDatastore
}

// NewMockDatastore creates a new mock instance
func NewMockDatastore(ctrl *gomock.Controller) *MockDatastore {
	mock := &MockDatastore{ctrl: ctrl}
	mock.recorder = &MockDatastoreMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockDatastore) EXPECT() *MockDatastoreMockRecorder {
	return m.recorder
}

// RawDB mocks base method
func (m *MockDatastore) RawDB() *sqlx.DB {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "RawDB")
	ret0, _ := ret[0].(*sqlx.DB)
	return ret0
}

// RawDB indicates an expected call of RawDB
func (mr *MockDatastoreMockRecorder) RawDB() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "RawDB", reflect.TypeOf((*MockDatastore)(nil).RawDB))
}

// NewMigrate mocks base method
func (m *MockDatastore) NewMigrate() (*v4.Migrate, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "NewMigrate")
	ret0, _ := ret[0].(*v4.Migrate)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// NewMigrate indicates an expected call of NewMigrate
func (mr *MockDatastoreMockRecorder) NewMigrate() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "NewMigrate", reflect.TypeOf((*MockDatastore)(nil).NewMigrate))
}

// Migrate mocks base method
func (m *MockDatastore) Migrate(currentMigrationVersion uint) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Migrate", currentMigrationVersion)
	ret0, _ := ret[0].(error)
	return ret0
}

// Migrate indicates an expected call of Migrate
func (mr *MockDatastoreMockRecorder) Migrate(currentMigrationVersion interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Migrate", reflect.TypeOf((*MockDatastore)(nil).Migrate), currentMigrationVersion)
}

// RollbackTxAndHandle mocks base method
func (m *MockDatastore) RollbackTxAndHandle(tx *sqlx.Tx) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "RollbackTxAndHandle", tx)
	ret0, _ := ret[0].(error)
	return ret0
}

// RollbackTxAndHandle indicates an expected call of RollbackTxAndHandle
func (mr *MockDatastoreMockRecorder) RollbackTxAndHandle(tx interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "RollbackTxAndHandle", reflect.TypeOf((*MockDatastore)(nil).RollbackTxAndHandle), tx)
}

// RollbackTx mocks base method
func (m *MockDatastore) RollbackTx(tx *sqlx.Tx) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "RollbackTx", tx)
}

// RollbackTx indicates an expected call of RollbackTx
func (mr *MockDatastoreMockRecorder) RollbackTx(tx interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "RollbackTx", reflect.TypeOf((*MockDatastore)(nil).RollbackTx), tx)
}

// GetGrantsOrderedByExpiry mocks base method
func (m *MockDatastore) GetGrantsOrderedByExpiry(wallet wallet.Info, promotionType string) ([]Grant, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetGrantsOrderedByExpiry", wallet, promotionType)
	ret0, _ := ret[0].([]Grant)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetGrantsOrderedByExpiry indicates an expected call of GetGrantsOrderedByExpiry
func (mr *MockDatastoreMockRecorder) GetGrantsOrderedByExpiry(wallet, promotionType interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetGrantsOrderedByExpiry", reflect.TypeOf((*MockDatastore)(nil).GetGrantsOrderedByExpiry), wallet, promotionType)
}

// MockReadOnlyDatastore is a mock of ReadOnlyDatastore interface
type MockReadOnlyDatastore struct {
	ctrl     *gomock.Controller
	recorder *MockReadOnlyDatastoreMockRecorder
}

// MockReadOnlyDatastoreMockRecorder is the mock recorder for MockReadOnlyDatastore
type MockReadOnlyDatastoreMockRecorder struct {
	mock *MockReadOnlyDatastore
}

// NewMockReadOnlyDatastore creates a new mock instance
func NewMockReadOnlyDatastore(ctrl *gomock.Controller) *MockReadOnlyDatastore {
	mock := &MockReadOnlyDatastore{ctrl: ctrl}
	mock.recorder = &MockReadOnlyDatastoreMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockReadOnlyDatastore) EXPECT() *MockReadOnlyDatastoreMockRecorder {
	return m.recorder
}

// RawDB mocks base method
func (m *MockReadOnlyDatastore) RawDB() *sqlx.DB {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "RawDB")
	ret0, _ := ret[0].(*sqlx.DB)
	return ret0
}

// RawDB indicates an expected call of RawDB
func (mr *MockReadOnlyDatastoreMockRecorder) RawDB() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "RawDB", reflect.TypeOf((*MockReadOnlyDatastore)(nil).RawDB))
}

// NewMigrate mocks base method
func (m *MockReadOnlyDatastore) NewMigrate() (*v4.Migrate, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "NewMigrate")
	ret0, _ := ret[0].(*v4.Migrate)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// NewMigrate indicates an expected call of NewMigrate
func (mr *MockReadOnlyDatastoreMockRecorder) NewMigrate() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "NewMigrate", reflect.TypeOf((*MockReadOnlyDatastore)(nil).NewMigrate))
}

// Migrate mocks base method
func (m *MockReadOnlyDatastore) Migrate(currentMigrationVersion uint) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Migrate", currentMigrationVersion)
	ret0, _ := ret[0].(error)
	return ret0
}

// Migrate indicates an expected call of Migrate
func (mr *MockReadOnlyDatastoreMockRecorder) Migrate(currentMigrationVersion interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Migrate", reflect.TypeOf((*MockReadOnlyDatastore)(nil).Migrate), currentMigrationVersion)
}

// RollbackTxAndHandle mocks base method
func (m *MockReadOnlyDatastore) RollbackTxAndHandle(tx *sqlx.Tx) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "RollbackTxAndHandle", tx)
	ret0, _ := ret[0].(error)
	return ret0
}

// RollbackTxAndHandle indicates an expected call of RollbackTxAndHandle
func (mr *MockReadOnlyDatastoreMockRecorder) RollbackTxAndHandle(tx interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "RollbackTxAndHandle", reflect.TypeOf((*MockReadOnlyDatastore)(nil).RollbackTxAndHandle), tx)
}

// RollbackTx mocks base method
func (m *MockReadOnlyDatastore) RollbackTx(tx *sqlx.Tx) {
	m.ctrl.T.Helper()
	m.ctrl.Call(m, "RollbackTx", tx)
}

// RollbackTx indicates an expected call of RollbackTx
func (mr *MockReadOnlyDatastoreMockRecorder) RollbackTx(tx interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "RollbackTx", reflect.TypeOf((*MockReadOnlyDatastore)(nil).RollbackTx), tx)
}

// GetGrantsOrderedByExpiry mocks base method
func (m *MockReadOnlyDatastore) GetGrantsOrderedByExpiry(wallet wallet.Info, promotionType string) ([]Grant, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetGrantsOrderedByExpiry", wallet, promotionType)
	ret0, _ := ret[0].([]Grant)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetGrantsOrderedByExpiry indicates an expected call of GetGrantsOrderedByExpiry
func (mr *MockReadOnlyDatastoreMockRecorder) GetGrantsOrderedByExpiry(wallet, promotionType interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetGrantsOrderedByExpiry", reflect.TypeOf((*MockReadOnlyDatastore)(nil).GetGrantsOrderedByExpiry), wallet, promotionType)
}
