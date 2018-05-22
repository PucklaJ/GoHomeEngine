package gohome

import (
	"log"
	"os"
)

const (
	ERROR_LEVEL_LOG     uint8 = iota
	ERROR_LEVEL_ERROR   uint8 = iota
	ERROR_LEVEL_WARNING uint8 = iota
	ERROR_LEVEL_FATAL   uint8 = iota
)

type GoHomeError struct {
	errorString string
}

type ErrorMessage struct {
	ErrorLevel uint8
	Tag        string
	ObjectName string
	Err        error
}

type ErrorManager struct {
	ErrorLevel        uint8
	DuplicateMessages bool
	ShowMessageBoxes  bool

	messages []ErrorMessage
}

func (this *GoHomeError) Error() string {
	return this.errorString
}

func (this *ErrorMessage) Error() string {
	return this.errorLevelToString() + "\t: " + this.Tag + "\t: " + this.ObjectName + "\t: " + this.Err.Error()
}

func (this *ErrorMessage) errorLevelToString() string {
	switch this.ErrorLevel {
	case ERROR_LEVEL_LOG:
		return "LOG"
	case ERROR_LEVEL_ERROR:
		return "ERROR"
	case ERROR_LEVEL_WARNING:
		return "WARNING"
	default:
		return "MESSAGE"
	}
}

func (this *ErrorMessage) Equals(other ErrorMessage) bool {
	return this.ErrorLevel == other.ErrorLevel && this.Tag == other.Tag && this.ObjectName == other.ObjectName && this.Err.Error() == other.Err.Error()
}

func (this *ErrorManager) Init() {
	this.ErrorLevel = ERROR_LEVEL_ERROR
	this.DuplicateMessages = false
	this.ShowMessageBoxes = true
}

func (this *ErrorManager) Message(errorLevel uint8, tag string, objectName string, err string) {
	this.MessageError(errorLevel, tag, objectName, &GoHomeError{err})
}

func (this *ErrorManager) MessageError(errorLevel uint8, tag string, objectName string, err error) {
	if errorLevel > this.ErrorLevel {
		return
	}
	errMsg := ErrorMessage{
		ErrorLevel: errorLevel,
		Tag:        tag,
		ObjectName: objectName,
		Err:        err,
	}
	if errorLevel != ERROR_LEVEL_FATAL && !this.DuplicateMessages {
		for i := 0; i < len(this.messages); i++ {
			if this.messages[i].Equals(errMsg) {
				return
			}
		}
		this.messages = append(this.messages, errMsg)
	}
	this.onNewError(errMsg)
	if errorLevel == ERROR_LEVEL_FATAL {
		panic(errMsg)
	}
}

func (this *ErrorManager) onNewError(errMsg ErrorMessage) {
	defer func() {
		rec := recover()
		if rec != nil {
			log.Println("Recovered:", rec)
		}
	}()
	log.Println(errMsg.Error())
	if this.ShowMessageBoxes && (errMsg.ErrorLevel == ERROR_LEVEL_ERROR || errMsg.ErrorLevel == ERROR_LEVEL_FATAL) {
		if Framew.ShowYesNoDialog("An Error accoured", errMsg.Error()+"\nContinue?") == DIALOG_NO {
			MainLop.Quit()
			os.Exit(int(errMsg.ErrorLevel))
		}
	}
}

func (this *ErrorManager) Reset() {
	if len(this.messages) > 0 {
		this.messages = this.messages[:0]
	}
}

func (this *ErrorManager) Terminate() {
	this.Reset()
}

var ErrorMgr ErrorManager
