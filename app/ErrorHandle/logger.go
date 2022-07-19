package ErrorHandle

import (
	"fmt"
	"io"
	"log"
)

var (
	Info    *log.Logger
	Warning *log.Logger
	Error   *log.Logger
)

func Init(
	infoHandle io.Writer,
	warningHandle io.Writer,
	ErrorHandle io.Writer) {

	Info = log.New(infoHandle,
		"INFO: ",
		log.Ldate|log.Ltime|log.Lshortfile)

	Warning = log.New(warningHandle,
		"WARNING: ",
		log.Ldate|log.Ltime|log.Lshortfile)

	Error = log.New(ErrorHandle,
		"ERROR: ",
		log.Ldate|log.Ltime|log.Lshortfile)
}

func PrintErr(err error, str string) {
	if err != nil {
		Error.Println(str)
		fmt.Println("=errText=" + err.Error())
	}
}
