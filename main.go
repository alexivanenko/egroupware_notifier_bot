package main

import (
	"time"

	"github.com/alexivanenko/egroupware_notifier_bot/config"
	"github.com/alexivanenko/egroupware_notifier_bot/model"
	"github.com/alexivanenko/egroupware_notifier_bot/notifier"
)

func main() {
	defer model.GetDB().Close()
	config.Log("Run GMail Notifier")

	doEvery(5*time.Minute, run)
}

//doEvery executes given function by duration period
func doEvery(d time.Duration, f func(time.Time)) {
	for x := range time.Tick(d) {
		f(x)
	}
}

func run(time.Time) {
	notifier.LoadMailsAndNotify()
}
