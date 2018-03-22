package memstore_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/icrowley/fake"

	"github.com/c0va23/redirector/memstore"
	"github.com/c0va23/redirector/models"
	"github.com/c0va23/redirector/store"
	"github.com/c0va23/redirector/test/factories"
)

func TestMemStore(t *testing.T) {
	a := assert.New(t)

	a.Implements((*store.Store)(nil), new(memstore.MemStore))
}

func TestListHostRules_Empty(t *testing.T) {
	a := assert.New(t)
	s := memstore.NewMemStore()

	listHostRules, err := s.ListHostRules()
	a.Nil(err)
	a.Equal(listHostRules, []models.HostRules{})
}

func TestReplaceHostRules_AddOne(t *testing.T) {
	a := assert.New(t)
	s := memstore.NewMemStore()

	hostRules := factories.HostRulesFactory.MustCreate().(models.HostRules)
	err := s.ReplaceHostRules(hostRules)
	a.Nil(err)

	listHostRules, err := s.ListHostRules()
	a.Equal(listHostRules, []models.HostRules{hostRules})
}

func TestReplaceHostRules_AddMany(t *testing.T) {
	a := assert.New(t)
	s := memstore.NewMemStore()

	listHostRules := make([]models.HostRules, 0, 3)
	for i := 0; i < cap(listHostRules); i++ {
		hostRules := factories.HostRulesFactory.MustCreate().(models.HostRules)
		a.Nil(s.ReplaceHostRules(hostRules))
	}

	listHostRules, err := s.ListHostRules()
	a.Nil(err)
	a.Equal(listHostRules, listHostRules)
}

func TestReplaceHostRules_Replace(t *testing.T) {
	a := assert.New(t)
	s := memstore.NewMemStore()

	listHostRules := make([]models.HostRules, 0, 3)
	for i := 0; i < cap(listHostRules); i++ {
		hostRules := factories.HostRulesFactory.MustCreate().(models.HostRules)
		listHostRules = append(listHostRules, hostRules)
		a.Nil(s.ReplaceHostRules(hostRules))
	}

	updatedHostRules := factories.HostRulesFactory.
		MustCreateWithOption(map[string]interface{}{
			"Host": listHostRules[1].Host,
		}).(models.HostRules)

	a.Nil(s.ReplaceHostRules(updatedHostRules))

	updatedListHostRules, err := s.ListHostRules()
	a.Nil(err)
	a.Equal(updatedListHostRules[0], listHostRules[0])
	a.Equal(updatedListHostRules[1], updatedHostRules)
	a.Equal(updatedListHostRules[2], listHostRules[2])
}

func TestGetHostRule_Success(t *testing.T) {
	a := assert.New(t)
	s := memstore.NewMemStore()

	sourceHostRules := factories.HostRulesFactory.MustCreate().(models.HostRules)
	a.Nil(s.ReplaceHostRules(sourceHostRules))

	hostRule, err := s.GetHostRules(sourceHostRules.Host)
	a.Nil(err)
	a.Equal(hostRule, &sourceHostRules)
}

func TestGetHostRule_NotFound(t *testing.T) {
	a := assert.New(t)
	s := memstore.NewMemStore()

	hostRule, err := s.GetHostRules(fake.DomainName())
	a.Nil(err)
	a.Nil(hostRule)
}
