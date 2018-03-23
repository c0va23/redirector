package redisstore_test

import (
	"encoding/json"

	"github.com/mediocregopher/radix.v2/redis"

	"github.com/c0va23/redirector/test/factories"
	"github.com/c0va23/redirector/test/mocks"

	"github.com/c0va23/redirector/models"
)

const startCursor = ""
const finishCursor = "0"

func buildListHostRules(size int) []models.HostRules {
	listHostRules := make([]models.HostRules, 0, size)
	for i := 0; i < cap(listHostRules); i++ {
		listHostRules = append(
			listHostRules,
			factories.HostRulesFactory.MustCreate().(models.HostRules),
		)
	}
	return listHostRules
}

func registerHostsScan(
	cmderMock *mocks.CmderMock,
	cursor string,
	nextCursor string,
	listHostRules ...models.HostRules,
) {
	hosts := make([]interface{}, 0, len(listHostRules))
	for _, hostRule := range listHostRules {
		hosts = append(hosts, hostRule.Host)
	}
	cmderMock.On("Cmd", "SCAN", []interface{}{cursor}).
		Return(redis.NewResp(
			[]interface{}{
				nextCursor,
				hosts,
			},
		))
}

func registerHosts(
	cmderMock *mocks.CmderMock,
	listHostRules ...models.HostRules,
) {
	for _, hostRules := range listHostRules {
		jsonBytes, err := json.Marshal(hostRules)
		if nil != err {
			panic(err)
		}

		cmderMock.On("Cmd", "GET", []interface{}{hostRules.Host}).Return(
			redis.NewResp(string(jsonBytes)),
		)
	}
}
