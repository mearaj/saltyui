package alog

import (
	log "github.com/sirupsen/logrus"
	"os"
	"runtime"
	"strconv"
	"strings"
)

func init() {
	log.SetFormatter(&log.JSONFormatter{})
	log.SetOutput(os.Stdout)
	//log.SetOutput(ioutil.Discard)
}
func Logger() *log.Entry {
	pc, file, line, ok := runtime.Caller(1)
	if !ok {
		panic("Could not get context info for logger!")
	}
	filename := file[strings.LastIndex(file, "/")+1:] + ":" + strconv.Itoa(line)
	funcname := runtime.FuncForPC(pc).Name()
	fn := funcname[strings.LastIndex(funcname, ".")+1:]
	return log.WithField("file", filename).WithField("function", fn)
}

func Fatal(args ...interface{}) {
	Logger().Fatal(args)
}

func Error(args ...interface{}) {
	Logger().Error(args)
}

func Debug(args ...interface{}) {
	Logger().Debug(args)
}

func Print(args ...interface{}) {
	Logger().Print(args)
}

func Println(args ...interface{}) {
	Logger().Println(args)
}

func Info(args ...interface{}) {
	Logger().Info(args)
}

func Warn(args ...interface{}) {
	Logger().Warn(args)
}

func Warning(args ...interface{}) {
	Logger().Warning(args)
}

func Panic(args ...interface{}) {
	Logger().Panic(args)
}
