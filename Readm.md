### It package is one wrapper to "go.mau.fi/whatsmeow/util/log" and implement new options 

# Example basic of use
```
package main

import "logecs/log"

func main() {
	// Logecs := log.NewLoggerEcs("Modulo", "debug", true, false)
	Logecs := log.NewLoggerEcs(log.EcsLogger{
		Mod: "ModuleName", Color: true,
		Path: "output.log", OutPut: true,
	})
	Logecs.Debugf("Modulo iniciado")
	Logecs.Warnf("Warning %s", "Modulo iniciado")
	Logecs.Errorf("Error %s", "Modulo iniciado")
	Logecs.Infof("Info %s", "Modulo iniciado")
}
```
