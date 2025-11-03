// Package logecs is a wrapper for the standard log package, providing enhanced
// structured logging capabilities with support for colors, output to files,
// and different log levels.
//
// Example usage:
//
// package main
//
// import logecs "github.com/ecsavigne/logecs/log"
//
//	func main() {
//	    Logecs := logecs.NewLoggerEcs(logecs.EcsLogger{
//	        Mod: "ModuleName", Color: true,
//	        Path: "output.log", OutPut: true,
//	    })
//	    Logecs.Debugf("Modulo iniciado")
//	    Logecs.Warnf("Warning %s", "Modulo iniciado")
//	    Logecs.Errorf("Error %s", "Modulo iniciado")
//	    Logecs.Infof("Info %s", "Modulo iniciado")
//	}
package log

import (
	"fmt"
	"log"
	"os"
	"path"
	"runtime"
	"strings"
	"time"
)

// Logger is a simple interface Logger implementation that outputs to stdout.
type Logger interface {
	// outputs a log message with the given level off info and message.
	// example:
	//
	//	Logecs := logecs.NewLoggerEcs(logecs.EcsLogger{
	//			Mod: "ModuleName", Color: true,
	//			Path: "output.log", OutPut: true,
	//		})
	//
	// Logecs.Warnf("Warning %s", "Module initilized")
	Warnf(msg string, args ...any)
	// outputs a log message with the given level off error and message.
	// example:
	//
	//	Logecs := logecs.NewLoggerEcs(logecs.EcsLogger{
	//			Mod: "ModuleName", Color: true,
	//			Path: "output.log", OutPut: true,
	//		})
	//
	// Logecs.Errorf("Error %s", "Module initilized")
	Errorf(msg string, args ...any)
	// outputs a log message with the given level off info and message.
	// example:
	//
	//	Logecs := logecs.NewLoggerEcs(logecs.EcsLogger{
	//			Mod: "ModuleName", Color: true,
	//			Path: "output.log", OutPut: true,
	//		})
	//
	// Logecs.Infof("Info %s", "Module initilized")
	Infof(msg string, args ...any)
	// outputs a log message with the given level off info and message.
	// example:
	//
	//	Logecs := logecs.NewLoggerEcs(logecs.EcsLogger{
	//			Mod: "ModuleName", Color: true,
	//			Path: "output.log", OutPut: true,
	//		})
	//
	// Logecs.Debugf("Debug %s", "Module initilized")
	Debugf(msg string, args ...any)
	// func (s *EcsLogger) Sub(mod string) Logger
	Sub(module string) Logger
}

type EcsLogger struct {
	Mod            string // Name of module of the log, exemple: "NameService"
	Color          bool   // If true, then info, warn and error logs will be colored cyan, yellow and red respectively using ANSI color escape codes
	OutPut         bool   // If true, then the log will be output to stdout
	min            int
	NotStandardPut bool // If true, then the log will be to standard output. If not especified do not output to standard output
	LevelMin       string
	logger         *log.Logger
	Path           string // path to the log file
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

func (s *EcsLogger) outputf(level, msg string, args ...any) {
	if levelToInt[level] < s.min {
		return
	}
	var colorStart, colorReset string
	if s.Color {
		colorStart = colors[level]
		colorReset = "\033[0m"
	}

	tracert := ""
	_, file, line, ok := runtime.Caller(2) // Aumenta el número si sigue mostrando `wrapper.go`
	if ok {
		dir, file := path.Split(file)
		file = path.Join(path.Base(dir), file)
		tracert = fmt.Sprintf(" %s:%d", file, line)
	}

	layout := "2006/01/02 15:04:05"
	if s.Mod == "" {
		if s.logger != nil {
			s.logger.Printf("%s%s%s [%s] %s%s", time.Now().Format(layout), tracert, colorStart, level, colorReset, fmt.Sprintf(msg, args...))
		}
		if !s.NotStandardPut {
			fmt.Printf("%s%s%s [%s] %s%s", time.Now().Format(layout), tracert, colorStart, level, colorReset, fmt.Sprintf(msg, args...))
		}
	} else {
		if s.logger != nil {
			s.logger.Printf("%s%s%s [%s %s] %s%s", time.Now().Format(layout), tracert, colorStart, s.Mod, level, colorReset, fmt.Sprintf(msg, args...))
		}
		if !s.NotStandardPut {
			fmt.Printf("%s%s%s [%s %s] %s%s", time.Now().Format(layout), tracert, colorStart, s.Mod, level, colorReset, fmt.Sprintf(msg, args...))
		}
	}
}

func (s *EcsLogger) Errorf(msg string, args ...any) { s.outputf("ERROR", msg, args...) }

func (s *EcsLogger) Warnf(msg string, args ...any) { s.outputf("WARN", msg, args...) }

func (s *EcsLogger) Infof(msg string, args ...any) { s.outputf("INFO", msg, args...) }

func (s *EcsLogger) Debugf(msg string, args ...any) { s.outputf("DEBUG", msg, args...) }

func (s *EcsLogger) Sub(mod string) Logger {
	return &EcsLogger{Mod: fmt.Sprintf("%s/%s", s.Mod, mod), Color: s.Color, min: s.min, logger: s.logger, Path: s.Path, OutPut: s.OutPut}
}

// Stdout is a simple Logger implementation that outputs to stdout. The module name given is included in log lines.
// minLevel specifies the minimum log level to output. An empty string will output all logs.
// If color is true, then info, warn and error logs will be colored cyan, yellow and red respectively using ANSI color escape codes.
func stdout(module string, minLevel, path string, color, output, NotStandardPut bool) Logger {
	var (
		fileLog *os.File
		err     error
		logger  *log.Logger
	)
	if path != "" && output {
		fileLog, err = os.OpenFile(path, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0777)
		if err != nil {
			panic(fmt.Sprintln("Error opening file:", err))
		}
		logger = log.New(fileLog, "", 0)
	}

	return &EcsLogger{
		Mod:            module,
		Color:          color,
		OutPut:         output,
		NotStandardPut: NotStandardPut,
		Path:           path,
		min:            levelToInt[strings.ToUpper(minLevel)],
		logger:         logger,
	}
}

// NewLoggerEcs returns a Logger that outputs to stdout.
func NewLoggerEcs(l EcsLogger) Logger {
	return stdout(l.Mod, l.LevelMin, l.Path, l.Color, l.OutPut, l.NotStandardPut)
}
