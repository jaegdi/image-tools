package ocrequest

import (
	"log"
	"os"
)

var (
	WarningLogger *log.Logger
	InfoLogger    *log.Logger
	ErrorLogger   *log.Logger
)

func Init() {
	file, err := os.OpenFile("logs_ocimagetools.txt", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		log.Fatal(err)
	}

	InfoLogger = log.New(file, "INFO: ", log.Ldate|log.Ltime|log.Llongfile)
	WarningLogger = log.New(file, "WARNING: ", log.Ldate|log.Ltime|log.Llongfile)
	ErrorLogger = log.New(file, "ERROR: ", log.Ldate|log.Ltime|log.Llongfile)

	InfoLogger.Println("------------------------------------------------------------")
	InfoLogger.Println("Starting execution of clean-istags")
}
