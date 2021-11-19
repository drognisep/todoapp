package data

import "sync"

type Project struct {
	Name  string  `json:"name"`
	Tasks []*Task `json:"tasks"`
}

type Model struct {
	mux               sync.RWMutex
	Projects          []*Project `json:"projects"`
	LastProjectOpened string     `json:"lastProjectOpened"`
}

func (m *Model) ProjectNames() []string {
	m.mux.RLock()
	defer m.mux.RUnlock()
	var names []string
	for _, p := range m.Projects {
		names = append(names, p.Name)
	}
	return names
}

func (m *Model) HasProjectName(name string) bool {
	m.mux.RLock()
	defer m.mux.RUnlock()
	for _, p := range m.Projects {
		if p.Name == name {
			return true
		}
	}
	return false
}

func (m *Model) RemoveProject(project string) {
	m.mux.Lock()
	defer m.mux.Unlock()
	for i, p := range m.Projects {
		if p.Name == project {
			m.Projects = append(m.Projects[:i], m.Projects[i+1:]...)
			break
		}
	}
}

func (m *Model) AddProject(project *Project) {
	m.mux.Lock()
	defer m.mux.Unlock()
	m.Projects = append(m.Projects, project)
}

func (m *Model) RenameProject(project string, newName string) {
	m.mux.Lock()
	defer m.mux.Unlock()
	for _, p := range m.Projects {
		if p.Name == project {
			p.Name = newName
		}
	}
}
