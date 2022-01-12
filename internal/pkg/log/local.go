package log

//Debug output a debug message
func Debug(msg string, keyvals ...interface{}) {
	Log.WithFields(argsToFields(keyvals...)).Debug(msg)
}

//Info output a info message
func Info(msg string, keyvals ...interface{}) {
	Log.WithFields(argsToFields(keyvals...)).Info(msg)
}

//Error output a error message
func Error(msg string, keyvals ...interface{}) {
	Log.WithFields(argsToFields(keyvals...)).Error(msg)
}

//With return a logger with keyvals
func With(keyvals ...interface{}) Logger {
	return entry{
		Log.WithFields(argsToFields(keyvals...)),
	}
}
