// Package logecs is a wrapper for the standard log package and .
// Installation:
//  1. You can install the package by running the following command in your terminal in folder of project go.mod:
//  2. go get -u "github.com/ecsavigne/logecs" or go get github.com/ecsavigne/logecs@latest
//  3. go mod vendor

// Usage:
// file: main.go
// Code:
// package main
//
// import logecs "github.com/ecsavigne/logecs/log"
//
//	func main() {
//		// Logecs := log.NewLoggerEcs("Modulo", "debug", true, false)
//		Logecs := logecs.NewLoggerEcs(logecs.EcsLogger{
//			Mod: "ModuleName", Color: true,
//			Path: "output.log", OutPut: true,
//		})
//		Logecs.Debugf("Modulo iniciado")
//		Logecs.Warnf("Warning %s", "Modulo iniciado")
//		Logecs.Errorf("Error %s", "Modulo iniciado")
//		Logecs.Infof("Info %s", "Modulo iniciado")
//	}
package log

import (
	"fmt"
	"log"
	"os"
	"strings"
	"time"
)

// Logger is a simple Logger implementation that outputs to stdout.
type Logger interface {
	// outputs a log message with the given level off info and message.
	// example:
	// 	Logecs := logecs.NewLoggerEcs(logecs.EcsLogger{
	// 			Mod: "ModuleName", Color: true,
	// 			Path: "output.log", OutPut: true,
	// 		})
	// Logecs.Warnf("Warning %s", "Module initilized")
	Warnf(msg string, args ...interface{})
	// outputs a log message with the given level off error and message.
	// example:
	// 	Logecs := logecs.NewLoggerEcs(logecs.EcsLogger{
	// 			Mod: "ModuleName", Color: true,
	// 			Path: "output.log", OutPut: true,
	// 		})
	// Logecs.Errorf("Error %s", "Module initilized")
	Errorf(msg string, args ...interface{})
	// outputs a log message with the given level off info and message.
	// example:
	// 	Logecs := logecs.NewLoggerEcs(logecs.EcsLogger{
	// 			Mod: "ModuleName", Color: true,
	// 			Path: "output.log", OutPut: true,
	// 		})
	// Logecs.Infof("Info %s", "Module initilized")
	Infof(msg string, args ...interface{})
	// outputs a log message with the given level off info and message.
	// example:
	// 	Logecs := logecs.NewLoggerEcs(logecs.EcsLogger{
	// 			Mod: "ModuleName", Color: true,
	// 			Path: "output.log", OutPut: true,
	// 		})
	// Logecs.Debugf("Debug %s", "Module initilized")
	Debugf(msg string, args ...interface{})
	// func (s *EcsLogger) Sub(mod string) Logger
	Sub(module string) Logger
}

type EcsLogger struct {
	Mod      string // Name of module of the log, exemple: "NameService"
	Color    bool   // If true, then info, warn and error logs will be colored cyan, yellow and red respectively using ANSI color escape codes
	OutPut   bool   // If true, then the log will be output to stdout
	min      int
	LevelMin string
	logger   *log.Logger
	Path     string // path to the log file
}

var colors = map[string]string{
	"INFO":  "\033[36m",
	"WARN":  "\033[33m",
	"ERROR": "\033[31m",
}

// levelToInt is a map that maps log levels to integers.
var levelToInt = map[string]int{
	"":      -1,
	"DEBUG": 0,
	"INFO":  1,
	"WARN":  2,
	"ERROR": 3,
}

func (s *EcsLogger) outputf(level, msg string, args ...interface{}) {
	if levelToInt[level] < s.min {
		return
	}
	var colorStart, colorReset string
	if s.Color {
		colorStart = colors[level]
		colorReset = "\033[0m"
	}
	if s.Mod == "" {
		if s.logger != nil {
			s.logger.Printf("%s%s [%s] %s%s\n", time.Now().Format("15:04:05.000"), colorStart, level, fmt.Sprintf(msg, args...), colorReset)
		}
		fmt.Printf("%s%s [%s] %s%s\n", time.Now().Format("15:04:05.000"), colorStart, level, fmt.Sprintf(msg, args...), colorReset)
	} else {
		if s.logger != nil {
			s.logger.Printf("%s%s [%s %s] %s%s\n", time.Now().Format("15:04:05.000"), colorStart, s.Mod, level, fmt.Sprintf(msg, args...), colorReset)
		}
		fmt.Printf("%s%s [%s %s] %s%s\n", time.Now().Format("15:04:05.000"), colorStart, s.Mod, level, fmt.Sprintf(msg, args...), colorReset)
	}
}

func (s *EcsLogger) Errorf(msg string, args ...interface{}) { s.outputf("ERROR", msg, args...) }
func (s *EcsLogger) Warnf(msg string, args ...interface{})  { s.outputf("WARN", msg, args...) }
func (s *EcsLogger) Infof(msg string, args ...interface{})  { s.outputf("INFO", msg, args...) }
func (s *EcsLogger) Debugf(msg string, args ...interface{}) { s.outputf("DEBUG", msg, args...) }
func (s *EcsLogger) Sub(mod string) Logger {
	return &EcsLogger{Mod: fmt.Sprintf("%s/%s", s.Mod, mod), Color: s.Color, min: s.min}
}

// Stdout is a simple Logger implementation that outputs to stdout. The module name given is included in log lines.
//
// minLevel specifies the minimum log level to output. An empty string will output all logs.
//
// If color is true, then info, warn and error logs will be colored cyan, yellow and red respectively using ANSI color escape codes.
func stdout(module string, minLevel, path string, color, output bool) Logger {
	var (
		fileLog *os.File
		err     error
	)
	if path != "" && output {
		fileLog, err = os.OpenFile(path, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0777)
		if err != nil {
			panic(fmt.Sprintln("Error opening file:", err))
		}
	}
	return &EcsLogger{
		Mod:    module,
		Color:  color,
		OutPut: output,
		Path:   path,
		min:    levelToInt[strings.ToUpper(minLevel)],
		logger: log.New(fileLog, "", log.Ldate|log.Ltime|log.Lshortfile),
	}
}

// NewLoggerEcs returns a Logger that outputs to stdout.
func NewLoggerEcs(l EcsLogger) Logger {
	return stdout(l.Mod, l.LevelMin, l.Path, l.Color, l.OutPut)
}
