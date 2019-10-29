package gohome

import (
	"os"
)

// Used to set which level of messages should be printed
const (
	ERROR_LEVEL_LOG     uint8 = iota
	ERROR_LEVEL_ERROR   uint8 = iota
	ERROR_LEVEL_WARNING uint8 = iota
	ERROR_LEVEL_FATAL   uint8 = iota
)

// An useful struct for error handling
type GoHomeError struct {
	errorString string
}

// An ErrorMessage of the ErrorManager
type ErrorMessage struct {
	// The ErrorLevel of the message
	ErrorLevel uint8
	// The tag of the message (often the type of ObjectName)
	Tag        string
	// The name of the object which this message is about
	ObjectName string
	// The message itself as an error
	Err        error
}

// The ErrorManager used to log messages and errors
type ErrorManager struct {
	// The currently set error level
	// Determines which messages should be printed to the screen
	ErrorLevel        uint8
	// Determines wether everytime a message is written it should be 
	// printed or ignored
	DuplicateMessages bool
	// Wether small dialogues should pop up with errors
	ShowMessageBoxes  bool

	messages []ErrorMessage
}

// Returns the error message
func (this *GoHomeError) Error() string {
	return this.errorString
}

// Returns the whole message as a string
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
	case ERROR_LEVEL_FATAL:
		return "FATAL"
	default:
		return "MESSAGE"
	}
}

// Returns wether one ErrorMessage equals the based on every member
func (this *ErrorMessage) Equals(other ErrorMessage) bool {
	return this.ErrorLevel == other.ErrorLevel && this.Tag == other.Tag && this.ObjectName == other.ObjectName && this.Err.Error() == other.Err.Error()
}

// Initialises the ErrorManager with default values
func (this *ErrorManager) Init() {
	this.ErrorLevel = ERROR_LEVEL_ERROR
	this.DuplicateMessages = false
	this.ShowMessageBoxes = true
}

// Writes an error as a string
func (this *ErrorManager) Message(errorLevel uint8, tag string, objectName string, err string) {
	this.MessageError(errorLevel, tag, objectName, &GoHomeError{err})
}

// Writes an error using the error interface
func (this *ErrorManager) MessageError(errorLevel uint8, tag string, objectName string, err error) {
	if errorLevel > this.ErrorLevel && errorLevel != ERROR_LEVEL_FATAL {
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

type Stringer interface {
	String() string
}

func (this *ErrorManager) onNewError(errMsg ErrorMessage) {
	defer func() {
		rec := recover()
		if rec != nil {
			Framew.Log("Recovered: " + rec.(Stringer).String())
		}
	}()
	Framew.Log(errMsg.Error())
	if this.ShowMessageBoxes && (errMsg.ErrorLevel == ERROR_LEVEL_ERROR || errMsg.ErrorLevel == ERROR_LEVEL_FATAL) {
		if Framew.ShowYesNoDialog("An Error accoured", errMsg.Error()+"\nContinue?") == DIALOG_NO {
			MainLop.Quit()
			os.Exit(int(errMsg.ErrorLevel))
		}
	}
}

// Clears all previous messages
func (this *ErrorManager) Reset() {
	if len(this.messages) > 0 {
		this.messages = this.messages[:0]
	}
}

// Cleans up the manager
func (this *ErrorManager) Terminate() {
	this.Reset()
}

// Writes a message using the log error level
func (this *ErrorManager) Log(tag, objectName, message string) {
	this.Message(ERROR_LEVEL_LOG, tag, objectName, message)
}

// Writes a message using the error error level
func (this *ErrorManager) Error(tag, objectName, message string) {
	this.Message(ERROR_LEVEL_ERROR, tag, objectName, message)
}

// Writes a message using the warning error level
func (this *ErrorManager) Warning(tag, objectName, message string) {
	this.Message(ERROR_LEVEL_WARNING, tag, objectName, message)
}

// Writes a message using the fatal error level
func (this *ErrorManager) Fatal(tag, objectName, message string) {
	this.Message(ERROR_LEVEL_FATAL, tag, objectName, message)
}

// The ErrorManager that should be used for everything
var ErrorMgr ErrorManager
