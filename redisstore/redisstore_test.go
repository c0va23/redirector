package redisstore_test

import (
	"encoding/json"
	"fmt"
	"testing"

	"github.com/mediocregopher/radix.v2/redis"

	"github.com/icrowley/fake"
	"github.com/stretchr/testify/assert"

	"github.com/c0va23/redirector/models"
	"github.com/c0va23/redirector/redisstore"
	"github.com/c0va23/redirector/store"

	"github.com/c0va23/redirector/test/factories"
	"github.com/c0va23/redirector/test/mocks"
)

func TestRedisStore(t *testing.T) {
	a := assert.New(t)

	a.Implements((*store.Store)(nil), new(redisstore.RedisStore))
}

func TestListHostRules_NotExists(t *testing.T) {
	a := assert.New(t)

	cmder := new(mocks.CmderMock)
	registerHostsScan(cmder, startCursor, finishCursor)

	rs := redisstore.NewRedisStore(cmder)

	listHostRules, err := rs.ListHostRules()
	a.Nil(err)
	a.Equal([]models.HostRules{}, listHostRules)

	cmder.AssertExpectations(t)
}
func TestListHostRules_OneExists(t *testing.T) {
	a := assert.New(t)

	hostRules := factories.HostRulesFactory.MustCreate().(models.HostRules)

	cmder := new(mocks.CmderMock)

	registerHostsScan(cmder, startCursor, finishCursor, hostRules)
	registerHosts(cmder, hostRules)

	rs := redisstore.NewRedisStore(cmder)

	listHostRules, err := rs.ListHostRules()
	a.Nil(err)
	a.Equal([]models.HostRules{hostRules}, listHostRules)

	cmder.AssertExpectations(t)
}

func TestListHostRules_ManyExists(t *testing.T) {
	a := assert.New(t)

	sourceHostRules := buildListHostRules(3)

	cmder := new(mocks.CmderMock)
	registerHostsScan(cmder, startCursor, "15", sourceHostRules[0:2]...)
	registerHostsScan(cmder, "15", finishCursor, sourceHostRules[2:]...)

	registerHosts(cmder, sourceHostRules...)

	rs := redisstore.NewRedisStore(cmder)

	listHostRules, err := rs.ListHostRules()
	a.Nil(err)
	a.Equal(sourceHostRules, listHostRules)

	cmder.AssertExpectations(t)
}

func TestListHostRules_ScannerError(t *testing.T) {
	a := assert.New(t)

	cmder := new(mocks.CmderMock)

	cmder.On("Cmd", "SCAN", []interface{}{startCursor}).
		Return(redis.NewResp([]interface{}{"0"}))

	rs := redisstore.NewRedisStore(cmder)

	_, err := rs.ListHostRules()
	a.EqualError(err, "not enough parts returned")

	cmder.AssertExpectations(t)
}

func TestListHostRules_GetIoError(t *testing.T) {
	a := assert.New(t)

	sourceHostRules := buildListHostRules(1)

	cmder := new(mocks.CmderMock)
	registerHostsScan(cmder, startCursor, finishCursor, sourceHostRules...)

	ioErr := fmt.Errorf("Some io error")
	cmder.On("Cmd", "GET", []interface{}{sourceHostRules[0].Host}).
		Return(redis.NewRespIOErr(ioErr))

	rs := redisstore.NewRedisStore(cmder)

	_, err := rs.ListHostRules()
	a.EqualError(err, ioErr.Error())

	cmder.AssertExpectations(t)
}

func TestListHostRules_GetNotFound(t *testing.T) {
	a := assert.New(t)

	sourceHostRules := buildListHostRules(3)

	cmder := new(mocks.CmderMock)
	registerHostsScan(cmder, startCursor, finishCursor, sourceHostRules...)

	cmder.On("Cmd", "GET", []interface{}{sourceHostRules[0].Host}).
		Return(redis.NewResp(nil))

	rs := redisstore.NewRedisStore(cmder)

	_, err := rs.ListHostRules()
	a.EqualError(err, "response is nil")

	cmder.AssertExpectations(t)
}

func TestListHostRules_DecoreError(t *testing.T) {
	a := assert.New(t)

	sourceHostRules := buildListHostRules(3)

	cmder := new(mocks.CmderMock)
	registerHostsScan(cmder, startCursor, finishCursor, sourceHostRules...)

	cmder.On("Cmd", "GET", []interface{}{sourceHostRules[0].Host}).
		Return(redis.NewRespSimple("error"))

	rs := redisstore.NewRedisStore(cmder)

	_, err := rs.ListHostRules()
	a.EqualError(err, "invalid character 'e' looking for beginning of value")

	cmder.AssertExpectations(t)
}

func TestGetHostRules_Success(t *testing.T) {
	a := assert.New(t)

	sourceHostRule := factories.HostRulesFactory.MustCreate().(models.HostRules)
	hostRuleJson, _ := json.Marshal(sourceHostRule)

	cmder := new(mocks.CmderMock)
	cmder.On("Cmd", "GET", []interface{}{
		sourceHostRule.Host,
	}).Return(redis.NewRespSimple(string(hostRuleJson)))

	rs := redisstore.NewRedisStore(cmder)
	hostRule, err := rs.GetHostRules(sourceHostRule.Host)
	a.Nil(err)
	a.Equal(&sourceHostRule, hostRule)

	cmder.AssertExpectations(t)
}

func TestGetHostRules_NotFound(t *testing.T) {
	a := assert.New(t)

	host := "notexists.org"
	cmder := new(mocks.CmderMock)
	cmder.On("Cmd", "GET", []interface{}{host}).Return(redis.NewResp(nil))

	rs := redisstore.NewRedisStore(cmder)
	hostRule, err := rs.GetHostRules(host)
	a.Nil(err)
	a.Nil(hostRule)

	cmder.AssertExpectations(t)
}

func TestGetHostRules_IoErr(t *testing.T) {
	a := assert.New(t)

	host := "notexists.org"
	cmder := new(mocks.CmderMock)
	ioErr := fmt.Errorf("Some IO error")
	cmder.On("Cmd", "GET", []interface{}{host}).Return(redis.NewRespIOErr(ioErr))

	rs := redisstore.NewRedisStore(cmder)
	_, err := rs.GetHostRules(host)
	a.EqualError(err, ioErr.Error())

	cmder.AssertExpectations(t)
}

func TestGetHostRules_JsonError(t *testing.T) {
	a := assert.New(t)

	host := "notexists.org"

	cmder := new(mocks.CmderMock)
	cmder.On("Cmd", "GET", []interface{}{host}).Return(redis.NewRespSimple("erro"))

	rs := redisstore.NewRedisStore(cmder)
	_, err := rs.GetHostRules(host)
	a.EqualError(err, "invalid character 'e' looking for beginning of value")

	cmder.AssertExpectations(t)
}

func TestCreateHostRules_IoError(t *testing.T) {
	a := assert.New(t)

	hostRule := factories.HostRulesFactory.MustCreate().(models.HostRules)
	hostRuleJson, _ := json.Marshal(hostRule)

	cmder := new(mocks.CmderMock)
	ioErr := fmt.Errorf("Some IO error")
	cmder.On("Cmd", "SETNX", []interface{}{
		hostRule.Host,
		string(hostRuleJson),
	}).Return(redis.NewRespIOErr(ioErr))

	rs := redisstore.NewRedisStore(cmder)

	a.EqualError(
		rs.CreateHostRules(hostRule),
		ioErr.Error(),
	)

	cmder.AssertExpectations(t)
}

func TestCreateHostRules_Exists(t *testing.T) {
	a := assert.New(t)

	hostRule := factories.HostRulesFactory.MustCreate().(models.HostRules)
	hostRuleJson, _ := json.Marshal(hostRule)

	cmder := new(mocks.CmderMock)
	cmder.On("Cmd", "SETNX", []interface{}{
		hostRule.Host,
		string(hostRuleJson),
	}).Return(redis.NewResp(0))

	rs := redisstore.NewRedisStore(cmder)

	a.Equal(
		store.ErrExists,
		rs.CreateHostRules(hostRule),
	)

	cmder.AssertExpectations(t)
}

func TestCreateHostRules_Success(t *testing.T) {
	a := assert.New(t)

	hostRule := factories.HostRulesFactory.MustCreate().(models.HostRules)
	hostRuleJson, _ := json.Marshal(hostRule)

	cmder := new(mocks.CmderMock)
	cmder.On("Cmd", "SETNX", []interface{}{
		hostRule.Host,
		string(hostRuleJson),
	}).Return(redis.NewResp(1))

	rs := redisstore.NewRedisStore(cmder)

	// Not return error
	a.Nil(rs.CreateHostRules(hostRule))

	cmder.AssertExpectations(t)
}

func TestUpdateHostRules_ExistsError(t *testing.T) {
	a := assert.New(t)

	hostRule := factories.HostRulesFactory.MustCreate().(models.HostRules)

	ioErr := fmt.Errorf("Some IO error")

	cmder := new(mocks.CmderMock)
	cmder.On("Cmd", "EXISTS", []interface{}{
		hostRule.Host,
	}).Return(redis.NewRespIOErr(ioErr))

	rs := redisstore.NewRedisStore(cmder)

	a.Equal(
		ioErr,
		rs.UpdateHostRules(hostRule.Host, hostRule),
	)
	cmder.AssertExpectations(t)
}

func TestUpdateHostRules_NotFound(t *testing.T) {
	a := assert.New(t)

	hostRule := factories.HostRulesFactory.MustCreate().(models.HostRules)

	cmder := new(mocks.CmderMock)
	cmder.On("Cmd", "EXISTS", []interface{}{
		hostRule.Host,
	}).Return(redis.NewResp(0))

	rs := redisstore.NewRedisStore(cmder)

	a.Equal(
		store.ErrNotFound,
		rs.UpdateHostRules(hostRule.Host, hostRule),
	)
	cmder.AssertExpectations(t)
}

func TestUpdateHostRules_SetErrorWithoutUpdateHost(t *testing.T) {
	a := assert.New(t)

	hostRule := factories.HostRulesFactory.MustCreate().(models.HostRules)
	hostRuleJson, _ := json.Marshal(&hostRule)

	ioErr := fmt.Errorf("Some IO error")

	cmder := new(mocks.CmderMock)
	cmder.On("Cmd", "EXISTS", []interface{}{
		hostRule.Host,
	}).Return(redis.NewResp(1))

	cmder.On("Cmd", "SET", []interface{}{
		hostRule.Host,
		string(hostRuleJson),
	}).Return(redis.NewRespIOErr(ioErr))

	rs := redisstore.NewRedisStore(cmder)

	a.Equal(
		ioErr,
		rs.UpdateHostRules(hostRule.Host, hostRule),
	)
	cmder.AssertExpectations(t)
}

func TestUpdateHostRules_SuccessWithoutUpdateHost(t *testing.T) {
	a := assert.New(t)

	hostRule := factories.HostRulesFactory.MustCreate().(models.HostRules)
	hostRuleJson, _ := json.Marshal(&hostRule)

	cmder := new(mocks.CmderMock)
	cmder.On("Cmd", "EXISTS", []interface{}{
		hostRule.Host,
	}).Return(redis.NewResp(1))

	cmder.On("Cmd", "SET", []interface{}{
		hostRule.Host,
		string(hostRuleJson),
	}).Return(redis.NewRespSimple("OK"))

	rs := redisstore.NewRedisStore(cmder)

	a.Nil(rs.UpdateHostRules(hostRule.Host, hostRule))

	cmder.AssertExpectations(t)
}

func TestUpdateHostRules_TargetHostExistsError(t *testing.T) {
	a := assert.New(t)

	sourceHost := fake.DomainName()
	hostRule := factories.HostRulesFactory.MustCreate().(models.HostRules)
	ioErr := fmt.Errorf("Some IO error")

	cmder := new(mocks.CmderMock)
	cmder.On("Cmd", "EXISTS", []interface{}{
		sourceHost,
	}).Return(redis.NewResp(1))

	cmder.On("Cmd", "EXISTS", []interface{}{
		hostRule.Host,
	}).Return(redis.NewRespIOErr(ioErr))

	rs := redisstore.NewRedisStore(cmder)

	a.Equal(
		ioErr,
		rs.UpdateHostRules(sourceHost, hostRule),
	)

	cmder.AssertExpectations(t)
}

func TestUpdateHostRules_TargetHostExists(t *testing.T) {
	a := assert.New(t)

	sourceHost := fake.DomainName()
	hostRule := factories.HostRulesFactory.MustCreate().(models.HostRules)

	cmder := new(mocks.CmderMock)
	cmder.On("Cmd", "EXISTS", []interface{}{
		sourceHost,
	}).Return(redis.NewResp(1))

	cmder.On("Cmd", "EXISTS", []interface{}{
		hostRule.Host,
	}).Return(redis.NewResp(1))

	rs := redisstore.NewRedisStore(cmder)

	a.Equal(
		store.ErrExists,
		rs.UpdateHostRules(sourceHost, hostRule),
	)

	cmder.AssertExpectations(t)
}

func TestUpdateHostRules_SourceHostDelError(t *testing.T) {
	a := assert.New(t)

	sourceHost := fake.DomainName()
	hostRule := factories.HostRulesFactory.MustCreate().(models.HostRules)
	hostRuleJSON, _ := json.Marshal(&hostRule)
	ioErr := fmt.Errorf("Some IO error")

	cmder := new(mocks.CmderMock)
	cmder.On("Cmd", "EXISTS", []interface{}{
		sourceHost,
	}).Return(redis.NewResp(1))

	cmder.On("Cmd", "EXISTS", []interface{}{
		hostRule.Host,
	}).Return(redis.NewResp(0))
	cmder.On("Cmd", "SET", []interface{}{
		hostRule.Host,
		string(hostRuleJSON),
	}).Return(redis.NewRespSimple("OK"))
	cmder.On("Cmd", "DEL", []interface{}{
		sourceHost,
	}).Return(redis.NewRespIOErr(ioErr))

	rs := redisstore.NewRedisStore(cmder)

	a.Equal(
		ioErr,
		rs.UpdateHostRules(sourceHost, hostRule),
	)

	cmder.AssertExpectations(t)
}

func TestUpdateHostRules_SuccessWithUpdateHost(t *testing.T) {
	a := assert.New(t)

	sourceHost := fake.DomainName()
	hostRule := factories.HostRulesFactory.MustCreate().(models.HostRules)
	hostRuleJSON, _ := json.Marshal(&hostRule)

	cmder := new(mocks.CmderMock)
	cmder.On("Cmd", "EXISTS", []interface{}{
		sourceHost,
	}).Return(redis.NewResp(1))

	cmder.On("Cmd", "EXISTS", []interface{}{
		hostRule.Host,
	}).Return(redis.NewResp(0))
	cmder.On("Cmd", "SET", []interface{}{
		hostRule.Host,
		string(hostRuleJSON),
	}).Return(redis.NewRespSimple("OK"))
	cmder.On("Cmd", "DEL", []interface{}{
		sourceHost,
	}).Return(redis.NewResp(1))

	rs := redisstore.NewRedisStore(cmder)

	a.Nil(rs.UpdateHostRules(sourceHost, hostRule))

	cmder.AssertExpectations(t)
}

func TestDeleteHostRules_Success(t *testing.T) {
	a := assert.New(t)

	host := fake.DomainName()

	cmder := new(mocks.CmderMock)
	cmder.On("Cmd", "DEL", []interface{}{host}).
		Return(redis.NewResp(1))

	rs := redisstore.NewRedisStore(cmder)

	a.Nil(rs.DeleteHostRules(host))

	cmder.AssertExpectations(t)
}

func TestDeleteHostRules_NotFoundError(t *testing.T) {
	a := assert.New(t)

	host := fake.DomainName()

	cmder := new(mocks.CmderMock)
	cmder.On("Cmd", "DEL", []interface{}{host}).
		Return(redis.NewResp(0))

	rs := redisstore.NewRedisStore(cmder)

	a.Equal(
		store.ErrNotFound,
		rs.DeleteHostRules(host),
	)

	cmder.AssertExpectations(t)
}

func TestDeleteHostRules_IoErrorError(t *testing.T) {
	a := assert.New(t)

	host := fake.DomainName()
	ioErr := fmt.Errorf("DeleteError")

	cmder := new(mocks.CmderMock)
	cmder.On("Cmd", "DEL", []interface{}{host}).
		Return(redis.NewRespIOErr(ioErr))

	rs := redisstore.NewRedisStore(cmder)

	a.Equal(
		ioErr,
		rs.DeleteHostRules(host),
	)

	cmder.AssertExpectations(t)
}
