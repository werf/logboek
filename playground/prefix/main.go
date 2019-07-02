package main

import (
	"strings"

	"github.com/flant/logboek"
)

//prefix: ┌ process with prefix
//prefix: │ long sentence      ↵
//prefix: │ long sentence      ↵
//prefix: │ long sentence      ↵
//prefix: │ long sentence      ↵
//prefix: │ long sentence
//prefix: └ p ... (0.00 seconds)
//
//┌ process without prefix
//│ long sentence long         ↵
//│ sentence long sentence     ↵
//│ long sentence long         ↵
//│ sentence
//└ process w ... (0.00 seconds)
func main() {
	twidth := 30

	_ = logboek.Init()
	logboek.EnableFitMode()
	logboek.SetWidth(twidth)

	logboek.SetPrefix("prefix: ", logboek.ColorizeSuccess)

	_ = logboek.LogProcess("process with prefix", logboek.LogProcessOptions{}, func() error {
		logboek.LogInfoLn(strings.Repeat("long sentence ", 5))

		return nil
	})

	logboek.ResetPrefix()

	_ = logboek.LogProcess("process without prefix", logboek.LogProcessOptions{}, func() error {
		logboek.LogInfoLn(strings.Repeat("long sentence ", 5))

		return nil
	})
}
