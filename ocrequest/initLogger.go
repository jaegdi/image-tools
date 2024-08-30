package ocrequest

import (
	"fmt"
	"log"
	"os"
	"regexp"
	"runtime"
)

type T_DebugLogger log.Logger

// Global Vars
var (
	WarningLogger       *log.Logger
	InfoLogger          *log.Logger
	VerifyLogger        *log.Logger
	ErrorLogger         *log.Logger
	DebugLogger         *log.Logger
	Multiproc           bool
	regexValidNamespace *regexp.Regexp
	LogfileName         string
)

func InitLogging() {
	var logfile *os.File
	var err error

	if CmdParams.Options.ServerMode {
		logfile = os.Stdout
	} else {
		err = os.Remove(LogFileName)
		if err != nil && !os.IsNotExist(err) {
			log.Fatal(err)
		}
		logfile, err = os.OpenFile(LogFileName, os.O_TRUNC|os.O_CREATE|os.O_WRONLY, 0666)
		if err != nil {
			log.Fatal(err)
		}
	}

	InfoLogger = log.New(logfile, "INFO: ", log.Ldate|log.Ltime|log.Lmicroseconds|log.Lmsgprefix)
	VerifyLogger = log.New(logfile, "VERIFY: ", log.Ldate|log.Ltime|log.Lmicroseconds|log.Lmsgprefix)
	WarningLogger = log.New(logfile, "WARNING: ", log.Ldate|log.Ltime|log.Lmicroseconds|log.Lmsgprefix)
	ErrorLogger = log.New(logfile, "ERROR: ", log.Ldate|log.Ltime|log.Lmicroseconds|log.Lmsgprefix)
	DebugLogger = log.New(logfile, "DEBUG: ", log.Ldate|log.Ltime|log.Lmicroseconds|log.Lmsgprefix)
}

func InfoMsg(p ...interface{}) {
	caller := getCaller()
	p = append([]interface{}{caller}, p...)
	InfoLogger.Println(p...)
}

func VerifyMsg(p ...interface{}) {
	if CmdParams.Options.Verify {
		caller := getCaller()
		p = append([]interface{}{caller}, p...)
		if VerifyLogger == nil {
			InitLogging()
		}
		VerifyLogger.Println(p...)
	}
}

func ErrorMsg(p ...interface{}) {
	caller := getCaller()
	prevCaller := getPreviousCaller()
	p = append([]interface{}{caller}, p...)
	p = append([]interface{}{prevCaller}, p...)
	ErrorLogger.Println(p...)
}

func DebugMsg(p ...interface{}) {
	if CmdParams.Options.Debug {
		caller := getCaller()
		p = append([]interface{}{caller}, p...)
		DebugLogger.Println(p...)
	}
}

func getCaller() string {
	pc, _, line, ok := runtime.Caller(2)
	if !ok {
		return "unknown"
	}
	fn := runtime.FuncForPC(pc)
	if fn == nil {
		return "unknown"
	}
	return fmt.Sprintf("caller: %s:%d| ", fn.Name(), line)
}

func getPreviousCaller() string {
	pc, _, line, ok := runtime.Caller(3)
	if !ok {
		return "unknown"
	}
	fn := runtime.FuncForPC(pc)
	if fn == nil {
		return "unknown"
	}
	return fmt.Sprintf("prevCaller: %s:%d| ", fn.Name(), line)
}
