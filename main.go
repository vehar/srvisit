package main

import (
	"fmt"
	"os"
	"runtime"
	"strings"
)

func main() {
	logAdd(MESS_INFO, "Запускается сервер reVisit версии "+REVISIT_VERSION)

	runtime.GOMAXPROCS(runtime.NumCPU())

	loadOptions()

	for _, x := range os.Args {
		if strings.Contains(x, "node") {
			options.Mode = NODE
		} else if strings.Contains(x, "master") {
			options.Mode = MASTER
		}
	}

	if options.Mode != NODE {
		loadVNCList()
		loadCounters()
		loadProfiles()

		go helperThread() //используем для периодических действий(сохранения и т.п.)
		go httpServer()   //обработка веб запросов
		go mainServer()   //обработка основных команд от клиентов и агентов
	}

	if options.Mode != MASTER {
		go dataServer() //обработка потоков данных от клиентов
	}

	if options.Mode == MASTER {
		go masterServer() //общаемся с агентами
	}

	if options.Mode == NODE {
		go nodeClient() //клинет подключающийся к мастеру
	}

	var r string
	for r != "quit" {
		fmt.Scanln(&r)
	}

	logAdd(MESS_INFO, "Завершили работу")
}
