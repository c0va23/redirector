package factories

import (
	"math/rand"

	"github.com/bluele/factory-go/factory"
	"github.com/icrowley/fake"

	"github.com/c0va23/redirector/models"
)

// GeneratePath generate random path
func GeneratePath() string {
	return "/" + fake.Word() + "/" + fake.Word()
}

func pathFactory(args factory.Args) (interface{}, error) {
	return GeneratePath(), nil
}

// TargetFactory generate models.Target
var TargetFactory = factory.NewFactory(models.Target{}).
	Attr(
		"HTTPCode",
		func(args factory.Args) (interface{}, error) {
			return 300 + rand.Int31n(100), nil
		},
	).
	Attr("TargetPath", pathFactory)

// RuleFactory generate models.Rule
var RuleFactory = factory.NewFactory(models.Rule{}).
	SubFactory("Target", TargetFactory).
	Attr("SourcePathPattern", pathFactory)

// HostRuleFactory generate HostRule
var HostRuleFactory = factory.NewFactory(models.HostRule{}).
	Attr(
		"Host",
		func(args factory.Args) (interface{}, error) {
			return fake.DomainName(), nil
		},
	).
	SubFactory("DefaultTarget", TargetFactory).
	SubSliceFactory(
		"Rules",
		RuleFactory,
		func() int {
			return rand.Intn(5)
		},
	)
