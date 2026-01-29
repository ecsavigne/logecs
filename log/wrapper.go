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
	"maps"
	"os"
	"path"
	"runtime"
	"slices"
	"sort"
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
	Create(info InfoLog)
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
	runtimeCaller  int    // defautl = 3 if log is creating from Craete = 3
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
	_, file, line, ok := runtime.Caller(s.runtimeCaller) // Aumenta el número si sigue mostrando `wrapper.go`
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
	if mod == "" {
		return s
	}

	return &EcsLogger{Mod: fmt.Sprintf("%s/%s", s.Mod, mod), Color: s.Color, min: s.min, logger: s.logger, Path: s.Path, OutPut: s.OutPut, runtimeCaller: s.runtimeCaller}
}

type TYPE_LOG int

const (
	Unknown = iota
	Info
	Error
	Warn
	Debug
)

func (t TYPE_LOG) ENUM() string {
	return []string{"Unknown", "Info", "Error", "Warn", "Debug"}[t]
}

// ej: event_service_logs {"tag": "value", .... "tagN": "valueN"}
type InfoLog struct {
	Type    TYPE_LOG // "Debug" | "Info" | "Warm" | "Error"
	Sub     string
	Name    string // Name of the log ej: "event_service_logs"
	Content map[string]any
}

func content(info InfoLog) string {
	cont := info.Content
	str := &strings.Builder{}
	str.WriteString(info.Name + " {")
	count := 1
	length := len(info.Content)

	keys := slices.Collect(maps.Keys(info.Content))
	sort.Strings(keys)

	for _, k := range keys {
		fmt.Fprintf(str, `"%s": `, k)

		switch v := info.Content[k].(type) {
		case string:
			fmt.Fprintf(str, `"%s"`, v)
		case int, int8, int16, int32, int64, uint, uint8, uint16, uint32, uint64:
			fmt.Fprintf(str, `%d`, v)
		case float64, float32:
			fmt.Fprintf(str, `%f`, v)
		case bool:
			fmt.Fprintf(str, `%t`, v)
		default:
			fmt.Fprintf(str, `%v`, v)
		}

		if length != count {
			count++
			str.WriteString(", ")
		}
	}

	str.WriteString("}\n")

	if cont != nil {
		return str.String()
	}

	return ""
}

func (s *EcsLogger) Create(info InfoLog) {
	s.runtimeCaller = 3

	cont := content(info)

	switch info.Type {
	case Info:
		s.Sub(info.Sub).Infof(cont)
	case Error:
		s.Sub(info.Sub).Errorf(cont)
	case Warn:
		s.Sub(info.Sub).Warnf(cont)
	case Debug:
		s.Sub(info.Sub).Debugf(cont)
	}
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
		runtimeCaller:  2,
	}
}

// NewLoggerEcs returns a Logger that outputs to stdout.
func NewLoggerEcs(l EcsLogger) Logger {
	return stdout(l.Mod, l.LevelMin, l.Path, l.Color, l.OutPut, l.NotStandardPut)
}
