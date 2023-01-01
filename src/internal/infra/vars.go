package infra

import (
	"time"

	"github.com/briandowns/spinner"
)

var s = spinner.New(spinner.CharSets[36], 100*time.Millisecond)
