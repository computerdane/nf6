package lib

import "os"

type StateSubDir struct {
	P    *string
	Name string
}

var (
	stateDir = ""
	subDirs  = []*StateSubDir{}
)

func SetStateDir(dir string) {
	stateDir = dir
}

func SetHomeStateDir(dir string) {
	home, err := os.UserHomeDir()
	if err != nil {
		Crash(err)
	}
	SetStateDir(home + "/.local/share/" + dir)
}

func SetSystemStateDir(dir string) {
	SetStateDir("/var/lib/" + dir + "/state")
}

func InitStateDir() {
	if err := os.MkdirAll(stateDir, os.ModePerm); err != nil {
		Crash(err)
	}
	for _, subDir := range subDirs {
		*(subDir.P) = stateDir + "/" + subDir.Name
		if err := os.MkdirAll(*subDir.P, os.ModePerm); err != nil {
			Crash(err)
		}
	}
}

func AddStateSubDir(subDir *StateSubDir) {
	subDirs = append(subDirs, subDir)
}
