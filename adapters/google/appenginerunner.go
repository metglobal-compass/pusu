package google

import (
	"github.com/metglobal-compass/pusu"
	"google.golang.org/appengine"
)

// App engine http handler which implements pusu.Runner interface
type appEngineRunner struct {
}

func (a *appEngineRunner) Run(subscription pusu.Subscription) error {
	appengine.Main()
	return nil
}
