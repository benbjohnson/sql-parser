package testutil

import (
	log "github.com/cihub/seelog"
)

// ConfigureTestLogger configures the global logger to print to console only
func ConfigureTestLogger() {

	testConfig := `
        <seelog type="sync" minlevel="debug">
            <outputs formatid="main"><console/></outputs>
            <formats><format id="main" format="%Date %Time [%LEVEL] %Msg%n"/></formats>
        </seelog>`

	logger, err := log.LoggerFromConfigAsBytes([]byte(testConfig))
	if err != nil {
		panic(err)
	}

	err = log.ReplaceLogger(logger)
	if err != nil {
		panic(err)
	}
}
