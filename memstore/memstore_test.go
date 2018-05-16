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

func TestCreateHostRules_AddOne(t *testing.T) {
	a := assert.New(t)
	s := memstore.NewMemStore()

	hostRules := factories.
		HostRulesFactory.
		MustCreate().(models.HostRules)
	err := s.CreateHostRules(hostRules)
	a.Nil(err)

	cratedHostRules, err := s.GetHostRules(hostRules.Host)
	a.Nil(err)
	a.Equal(cratedHostRules, &hostRules)
}

func TestCreateHostRules_AddMany(t *testing.T) {
	a := assert.New(t)
	s := memstore.NewMemStore()

	listHostRules := make([]models.HostRules, 0, 3)
	for i := 0; i < cap(listHostRules); i++ {
		hostRules := factories.
			HostRulesFactory.
			MustCreate().(models.HostRules)
		a.Nil(s.CreateHostRules(hostRules))
	}

	listHostRules, err := s.ListHostRules()
	a.Nil(err)
	a.Equal(listHostRules, listHostRules)
}

func TestCreateHostRules_Exists(t *testing.T) {
	a := assert.New(t)
	s := memstore.NewMemStore()

	hostRules := factories.
		HostRulesFactory.
		MustCreate().(models.HostRules)
	a.Nil(s.CreateHostRules(hostRules))

	a.Equal(
		store.ErrExists,
		s.CreateHostRules(hostRules),
	)
}

func TestUpdateHostRules_NotFound(t *testing.T) {
	a := assert.New(t)
	s := memstore.NewMemStore()

	hostRules := factories.
		HostRulesFactory.
		MustCreate().(models.HostRules)

	a.Equal(
		store.ErrNotFound,
		s.UpdateHostRules(hostRules.Host, hostRules),
	)
}

func TestUpdateHostRules_Success(t *testing.T) {
	a := assert.New(t)
	s := memstore.NewMemStore()

	existsHostRules := factories.
		HostRulesFactory.
		MustCreate().(models.HostRules)

	a.Nil(s.CreateHostRules(existsHostRules))

	newHostRules := factories.
		HostRulesFactory.
		MustCreateWithOption(map[string]interface{}{
			"Host": existsHostRules.Host,
		}).(models.HostRules)

	// Not return error
	a.Nil(s.UpdateHostRules(existsHostRules.Host, newHostRules))

	updatedHostRules, err := s.GetHostRules(existsHostRules.Host)
	a.Nil(err)

	// Update host rules
	a.Equal(
		&newHostRules,
		updatedHostRules,
	)
}

func TestGetHostRule_Success(t *testing.T) {
	a := assert.New(t)
	s := memstore.NewMemStore()

	sourceHostRules := factories.HostRulesFactory.MustCreate().(models.HostRules)
	a.Nil(s.CreateHostRules(sourceHostRules))

	hostRule, err := s.GetHostRules(sourceHostRules.Host)
	a.Nil(err)
	a.Equal(hostRule, &sourceHostRules)
}

func TestGetHostRule_NotFound(t *testing.T) {
	a := assert.New(t)
	s := memstore.NewMemStore()

	hostRule, err := s.GetHostRules(fake.DomainName())
	a.Nil(hostRule)
	a.Equal(store.ErrNotFound, err)
}

func TestDeleteHostRules_Success(t *testing.T) {
	a := assert.New(t)
	s := memstore.NewMemStore()

	// Host rules exists
	hostRules := factories.
		HostRulesFactory.
		MustCreate().(models.HostRules)
	a.Nil(s.CreateHostRules(hostRules))

	// Not return error
	a.Nil(s.DeleteHostRules(hostRules.Host))

	// Delete host rules
	_, err := s.GetHostRules(hostRules.Host)
	a.Equal(
		store.ErrNotFound,
		err,
	)
}
