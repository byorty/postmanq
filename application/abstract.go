package application

import (
	"github.com/AdOnWeb/postmanq/common"
	"github.com/AdOnWeb/postmanq/logger"
	"io/ioutil"
	"time"
)

type AbstractApplication struct {
	// путь до конфигурационного файла
	configFilename string

	// сервисы приложения, отправляющие письма
	services []interface{}

	// канал событий приложения
	events chan *common.ApplicationEvent

	// флаг, сигнализирующий окончание работы приложения
	done chan bool
}

func (a *AbstractApplication) IsValidConfigFilename(filename string) bool {
	return len(filename) > 0 && filename != common.ExampleConfigYaml
}

func (a *AbstractApplication) run(app common.Application, event *common.ApplicationEvent) {
	app.SetDone(make(chan bool))
	// создаем каналы для событий
	app.SetEvents(make(chan *common.ApplicationEvent, 3))
	go func() {
		for {
			select {
			case event := <-app.Events():
				if event.Kind == common.InitApplicationEventKind {
					// пытаемся прочитать конфигурационный файл
					bytes, err := ioutil.ReadFile(a.configFilename)
					if err == nil {
						event.Data = bytes
					} else {
						logger.FailExit("application can't read configuration file, error -  %v", err)
					}
				}

				for _, service := range app.Services() {
					switch event.Kind {
					case common.InitApplicationEventKind:
						app.FireInit(event, service)
					case common.RunApplicationEventKind:
						app.FireRun(event, service)
					case common.FinishApplicationEventKind:
						app.FireFinish(event, service)
					}
				}

				switch event.Kind {
				case common.InitApplicationEventKind:
					event.Kind = common.RunApplicationEventKind
					app.Events() <- event
				case common.FinishApplicationEventKind:
					time.Sleep(2 * time.Second)
					app.Done() <- true
				}
			}
		}
		close(app.Events())
	}()
	app.Events() <- event
	<-app.Done()
}

func (a *AbstractApplication) SetConfigFilename(configFilename string) {
	a.configFilename = configFilename
}

func (a *AbstractApplication) SetEvents(events chan *common.ApplicationEvent) {
	a.events = events
}

func (a *AbstractApplication) Events() chan *common.ApplicationEvent {
	return a.events
}

func (a *AbstractApplication) SetDone(done chan bool) {
	a.done = done
}

func (a *AbstractApplication) Done() chan bool {
	return a.done
}

func (a *AbstractApplication) Services() []interface{} {
	return a.services
}

func (a *AbstractApplication) FireInit(event *common.ApplicationEvent, abstractService interface{}) {
	service := abstractService.(common.Service)
	service.OnInit(event)
}

func (a *AbstractApplication) Run() {}

func (a *AbstractApplication) RunWithArgs(args ...interface{}) {}

func (a *AbstractApplication) FireRun(event *common.ApplicationEvent, abstractService interface{}) {}

func (a *AbstractApplication) FireFinish(event *common.ApplicationEvent, abstractService interface{}) {
}