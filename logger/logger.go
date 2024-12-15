package logger

import (
	"os"
	"path"
	"runtime"
	"strconv"

	"github.com/sirupsen/logrus"
)

var (
	logJson    *logrus.Logger
	logConsole *logrus.Logger
)

func jsonLog() {
	logJson = logrus.New()
	logJson.SetFormatter(&logrus.JSONFormatter{
		CallerPrettyfier: func(frame *runtime.Frame) (function string, file string) {
			pc, file, line, _ := runtime.Caller(8)

			fileName := path.Dir(file) + "/" + path.Base(file) + ", lineNumber:" + strconv.Itoa(line)
			funcName := runtime.FuncForPC(pc).Name()

			return funcName, fileName
		},
	})
	logJson.SetReportCaller(true)
	logJson.SetLevel(logrus.TraceLevel)
	file, err := os.OpenFile("application.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		logrus.Fatal("failed open log file : %s", err.Error())
	}
	logJson.SetOutput(file)
}

func consoleLog() {
	logConsole = logrus.New()
	logConsole.SetFormatter(&logrus.TextFormatter{
		CallerPrettyfier: func(frame *runtime.Frame) (function string, file string) {
			pc, file, line, _ := runtime.Caller(8)

			fileName := path.Dir(file) + "/" + path.Base(file) + ", lineNumber:" + strconv.Itoa(line)
			funcName := runtime.FuncForPC(pc).Name()

			return funcName, fileName
		},
	})
	logConsole.SetReportCaller(true)
	logConsole.SetLevel(logrus.InfoLevel)
}

func init() {
	jsonLog()
	consoleLog()
}

// func fileInfo() string {
// 	_, file, line, ok := runtime.Caller(3)
// 	if !ok {
// 		file = "<???>"
// 		line = 1
// 	} else {
// 		slash := strings.LastIndex(file, "/")
// 		if slash >= 0 {
// 			file = file[slash+1:]
// 		}
// 	}
// 	return fmt.Sprintf("%s:%d", file, line)

// }

// func caller() func(*runtime.Frame) (function string, file string) {
// 	return func(f *runtime.Frame) (function string, file string) {
// 		p, _ := os.Getwd()

// 		return f.Function, fmt.Sprintf("%s:%d", strings.TrimPrefix(f.File, p), f.Line)
// 	}
// }

// func entryLogs(fields map[string]interface{}) *logrus.Entry {
// 	if len(fields) > 0 {
// 		entry = logger.WithFields(fields)
// 	} else {
// 		entry = logger.WithFields(logrus.Fields{})
// 	}
// 	entry.Data["file"] = fileInfo()
// 	return entry
// }

func Info(message string, fields map[string]interface{}) {
	if len(fields) > 0 {
		logJson.WithFields(fields).Info(message)
		logConsole.WithFields(fields).Info(message)
		return
	}
	logJson.Info(message)
	logConsole.Info(message)
}

func Error(message string, fields map[string]interface{}) {
	if len(fields) > 0 {
		logJson.WithFields(fields).Error(message)
		logConsole.WithFields(fields).Error(message)
		return
	}
	logJson.Error(message)
	logConsole.Error(message)
}

func Debug(message string, fields map[string]interface{}) {
	if len(fields) > 0 {
		logJson.WithFields(fields).Debug(message)
		logConsole.WithFields(fields).Debug(message)
		return
	}
	logJson.Debug(message)
	logConsole.Debug(message)
}

func Warning(message string, fields map[string]interface{}) {
	if len(fields) > 0 {
		logJson.WithFields(fields).Warning(message)
		logConsole.WithFields(fields).Warning(message)
		return
	}
	logJson.Warning(message)
	logConsole.Warning(message)
}

func Fatal(message string, fields map[string]interface{}) {
	if len(fields) > 0 {
		defer func() {
			logConsole.WithFields(fields).Fatal(message)
		}()
		logJson.WithFields(fields).Fatal(message)
		return
	}

	logConsole.Info(message)
	logJson.Fatal(message)
}
