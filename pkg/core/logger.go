package core

import "fmt"

type LoggerFn func(a ...interface{})

type Logger struct {
	info LoggerFn
}

func NewPrintLogger() *Logger {
	return &Logger{
		info: func(a ...interface{}) {
			fmt.Println(a...)
		},
	}
}

func NewNullLogger() *Logger {
	return &Logger{
		info: func(a ...interface{}) {},
	}
}
