package guardian

import (
	"github.com/AdOnWeb/postmanq/common"
	"github.com/AdOnWeb/postmanq/logger"
	yaml "gopkg.in/yaml.v2"
)

var (
	service *Service
	events  = make(chan *common.SendEvent)
)

type Service struct {
	Hostnames      []string `yaml:"exclude"`
	hostnameLen    int
	GuardiansCount int      `yaml:"workers"`
}

func Inst() common.SendingService {
	if service == nil {
		service = new(Service)
	}
	return service
}

func (s *Service) OnInit(event *common.ApplicationEvent) {
	logger.Debug("init guardians...")
	err := yaml.Unmarshal(event.Data, s)
	if err == nil {
		s.hostnameLen = len(s.Hostnames)
		if s.GuardiansCount == 0 {
			s.GuardiansCount = common.DefaultWorkersCount
		}
	} else {
		logger.FailExitWithErr(err)
	}
}

func (s *Service) OnRun() {
	for i := 0; i < s.GuardiansCount; i++ {
		go newGuardian(i + 1)
	}
}

func (s *Service) Events() chan *common.SendEvent {
	return events
}

func (s *Service) OnFinish() {
	close(events)
}