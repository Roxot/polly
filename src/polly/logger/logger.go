package logger

import (
	"fmt"
	"os"
	"polly"
	"time"
)

const (
	cLogFile    = "server.log"
	cBufferSize = 1024
	cTimeFormat = "02/01 15:04:05"
)

type Logger struct {
	logFile  *os.File
	logChan  chan string
	quitChan chan int
}

func New() *Logger {
	logger := Logger{}
	logger.logChan = make(chan string, cBufferSize)
	logger.quitChan = make(chan int)
	return &logger
}

func (logger *Logger) Start() error {
	path, err := polly.GetPollyHome()
	if err != nil {
		return err
	}

	file, err := os.Create(path + cLogFile)
	if err != nil {
		return fmt.Errorf("$%s not set correctly.", polly.POLLY_HOME)
	}

	logger.logFile = file

	go func() {

	Loop:
		for {
			select {
			case <-logger.quitChan:
				break Loop
			case logMessage := <-logger.logChan:
				logger.logFile.WriteString(logMessage)
			}
		}

		logger.logFile.Close()
	}()

	return nil
}

func (logger *Logger) Stop() {
	logger.quitChan <- 1
}

func (logger *Logger) Log(tag, message, origin string) {
	logger.logChan <- fmt.Sprintf("%s: [%s] %s (%s)\n",
		time.Now().Format(cTimeFormat), tag, message, origin)
}
