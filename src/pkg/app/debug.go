package app

import (
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"
)

var debugFile *os.File

func initTeaDebug() {
	f, err := tea.LogToFile("debug.log", "debug")
	if err == nil {
		debugFile = f
	}
}

func writeDebugFLn(fmtStr string, args ...any) {
	if debugFile == nil {
		return
	}
	debugFile.Write([]byte(fmt.Sprintf(fmtStr+"\n", args...)))
}
