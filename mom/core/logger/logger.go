package logger

import (
	"fmt"
	"log"
	"os"
	"sync"
)

type Logger struct {
	service string
	base    *log.Logger
}

var (
	instance *Logger
	once     sync.Once
)

func Init(service string) *Logger {
	once.Do(func() {
		instance = &Logger{
			service: service,
			base:    log.New(os.Stdout, fmt.Sprintf("[%s] ", service), 0),
		}
	})

	return instance
}

func Get() *Logger {
	if instance == nil {
		return Init("app")
	}

	return instance
}

func (l *Logger) Println(v ...any) {
	l.base.Println(v...)
}

func (l *Logger) Printf(format string, v ...any) {
	l.base.Printf(format, v...)
}

func (l *Logger) Errorf(format string, v ...any) {
	l.base.Printf(format, v...)
}

func (l *Logger) Fatalf(format string, v ...any) {
	l.base.Fatalf(format, v...)
}
