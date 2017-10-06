package jlog

import (
	"context"
	"fmt"
	"path/filepath"
	"reflect"
	"runtime"
	"strings"
	"time"

	"github.com/fsnotify/fsnotify"
	"github.com/jerejones/jlog/target"
	"github.com/pkg/errors"
)

type Manager struct {
	cfg         *Config
	loggers     map[string]*logger
	routers     []*router
	targets     map[string]target.Target
	cancelWatch func()
}

func NewManager(cfg *Config) (*Manager, error) {
	m := Manager{
		cfg:         nil,
		loggers:     make(map[string]*logger),
		routers:     []*router{},
		targets:     make(map[string]target.Target),
		cancelWatch: func() {},
	}
	err := m.ApplyConfig(cfg)
	if err != nil {
		return nil, err
	}
	return &m, nil
}

func (m *Manager) ApplyConfig(cfg *Config) error {
	if cfg == nil {
		return nil
	}
	m.cancelWatch()

	targets := make(map[string]target.Target)

	for _, spec := range cfg.Targets {
		if _, exists := targets[spec.Name]; exists {
			return errors.Errorf("Duplicate target name: %s", spec.Name)
		}
		t, err := target.New(spec)
		if err != nil {
			if _, ok := err.(target.UnknownTargetError); ok {
				continue
			}
			return errors.Wrap(err, "Error applying config")
		}
		targets[spec.Name] = t
	}

	routers := []*router{}
	for _, spec := range cfg.Routes {
		router := newRouter(spec)
		targetNames := strings.Split(spec.WriteTo, ",")
		for _, targetName := range targetNames {
			t, exists := targets[targetName]
			if !exists {
				continue
			}
			router.targets = append(router.targets, t)
		}
		routers = append(routers, router)
	}

	m.cfg = cfg
	m.routers = routers
	m.targets = targets
	for _, logger := range m.loggers {
		m.updateLogger(logger)
	}

	if cfg.AutoReload {
		go m.watch(cfg.sourceFileName)
	}

	return nil
}

func (m *Manager) GetNamedLogger(name string) *logger {
	l, exists := m.loggers[name]
	if exists {
		return l
	}
	l = m.createLogger(name)
	if l != nil {
		m.loggers[name] = l
	}
	return l
}

func (m *Manager) ReloadConfig() error {
	return m.ApplyConfig(m.cfg)
}

func (m *Manager) GetPackageLogger() *logger {
	return m.getPackageLoggerSkip(3)
}

func (m *Manager) getPackageLoggerSkip(skip int) *logger {
	function := callerFrame(skip).Function
	packageName := function[:strings.LastIndex(function, ".")]
	return m.GetNamedLogger(packageName)
}

func (m *Manager) GetObjectLogger(obj interface{}) *logger {
	v := reflect.ValueOf(obj)
	for v.Kind() == reflect.Ptr {
		v = reflect.Indirect(v)
	}
	return m.GetTypeLogger(v.Type())
}

func (m *Manager) GetTypeLogger(t reflect.Type) *logger {
	pkg := t.PkgPath()
	name := t.Name()
	if len(pkg) == 0 {
		if len(name) == 0 {
			name = t.String()
		}
	} else {
		name = pkg + "." + name
	}
	return m.GetNamedLogger(name)
}

func (m *Manager) createLogger(name string) *logger {
	logger := &logger{name: name}
	m.updateLogger(logger)
	return logger
}

func (m *Manager) updateLogger(logger *logger) {
	var newRouters []*router
	for _, r := range m.routers {
		if r.Matches(logger.name) {
			newRouters = append(newRouters, r)
		}
	}
	logger.routers = newRouters
}

func (m *Manager) watch(filename string) {
	ctx, cancel := context.WithCancel(context.Background())
	m.cancelWatch = cancel

	absPath, err := filepath.Abs(filename)
	if err != nil {
		return
	}
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		return
	}
	defer watcher.Close()
	dir := filepath.Dir(absPath) + string(filepath.Separator)
	fmt.Println("Watching: ", dir)
	watcher.Add(dir)

	changeFound := false
	for !changeFound {
		select {
		case event := <-watcher.Events:
			if event.Name == absPath &&
				(event.Op&fsnotify.Write == fsnotify.Write ||
					event.Op&fsnotify.Create == fsnotify.Create) {
				changeFound = true
			}

		case err := <-watcher.Errors:
			fmt.Println("error:", err)
			return

		case <-ctx.Done():
			return
		}
	}

	time.Sleep(10 * time.Millisecond)

	fmt.Println("Reloading config")
	m.ApplyConfig(LoadConfig(filename))
}

func callerFrame(skip int) runtime.Frame {
	pc := make([]uintptr, 64)
	n := runtime.Callers(skip, pc)
	frames := runtime.CallersFrames(pc[:n])
	frame, _ := frames.Next()
	return frame
}
