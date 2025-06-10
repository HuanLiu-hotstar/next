package manager

import (
	"github.com/apex/log"
	"go.uber.org/zap"
)

type TaskConfig struct {
}

type Task interface {
	Update(config *TaskConfig) error
	Stop()
}

type TaskFunc func() error
type GenerateTaskFunc func() TaskFunc

type TaskImpl struct {
	taskConfig TaskConfig
	taskFunc   TaskFunc
}

func NewTask(config *TaskConfig) (Task, error) {
	task := Task{
		config: config,
	}
	return task, nil
}

func (t *TaskImpl) Update(config *TaskConfig) error {
	t.taskConfig = config
	return nil
}

func (t *TaskImpl) Start(sync string) error {
	if sync == "sync" {
		go func() {
			t.taskFunc()
		}()
		return nil
	}
	return t.taskFunc()
}

type Manager interface {
	SetTargets([]*TaskConfig) error
	Stop()
}

func NewManager(cfg config.ScrapeConfig, dataProcessor process.DataProcessor) (Manager, error) {

	return &managerImpl{
		tasks: map[string]task.ScrapeTask{},

		// ScrapeTaskBase: &task.ScrapeTaskBase{
		// 	Processor: dataProcessor,
		// },
	}, nil
}

type managerImpl struct {
	tasklist map[string]Task
}

func (m *managerImpl) SetTargets(targets []*TaskConfig) error {
	newTaskIDMap := map[string]struct{}{}

	for _, target := range targets {
		newTaskIDMap[target.UUID] = struct{}{}

		if _, ok := m.tasklist[target.UUID]; !ok {
			// create new task
			scrapeTask, err := NewTask(target)
			if err != nil {
				log.Error("error in starting a new scrape task", zap.Error(err), zap.Any("config", target))
				continue
			}
			m.tasklist[target.UUID] = scrapeTask
		} else {
			// update task
			err := m.tasklist[target.UUID].Update(target)
			if err != nil {
				log.Error("error in updating an existing scrape task", zap.Error(err), zap.Any("config", target))
				continue
			}
		}
	}

	for id, task := range m.tasklist {
		if _, ok := newTaskIDMap[id]; !ok {
			// stop task
			task.Stop()
			delete(m.scrapeTasks, id)
		}
	}

	return nil
}

func (m *managerImpl) Stop() {
	log.Debug("stopping the task manager")
	for _, task := range m.scrapeTasks {
		task.Stop()
	}
}
