package main

import "github.com/ecsavigne/logecs/log"

func main() {
	// Logecs := log.NewLoggerEcs("Modulo", "debug", true, false)
	Logecs := log.NewLoggerEcs(log.EcsLogger{
		Mod: "ModuleName", Color: true,
		Path: "output.log", OutPut: true,
		NotStandardPut: false,
	})
	Logecs.Debugf("Modulo iniciado\n")
	Logecs.Warnf("Warning %s", "Modulo iniciado\n")
	Logecs.Errorf("Error %s", "Modulo iniciado df sdf \n")
	Logecs.Infof("Info %s", "Modulo iniciado edfdfd dfd \n")

	info := log.InfoLog{
		Type: log.Debug,
		// Sub:  "SubModulo",
		Name: "test_name_logs",
		Content: map[string]any{
			"modulo": "Modulo",
			"sub_modulo": map[string]any{
				"SubModulo": "SubModulo",
			},
			"test":  1,
			"test2": true,
			"test3": 4.3,
			"a_b1":  "annene",
			"0_b1":  "llllll",
		},
	}

	Logecs.Create(info)
}
