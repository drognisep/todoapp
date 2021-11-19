package data

type Project struct {
	Name  string  `json:"name"`
	Tasks []*Task `json:"tasks"`
}

type Model struct {
	Projects          []*Project `json:"projects"`
	LastProjectOpened string     `json:"lastProjectOpened"`
}

func (m *Model) LastProject() *Project {
	for _, p := range m.Projects {
		if p.Name == m.LastProjectOpened {
			return p
		}
	}
	newProject := &Project{Name: m.LastProjectOpened}
	m.Projects = append(m.Projects, newProject)
	return newProject
}

func (m *Model) ProjectNames() []string {
	var names []string
	for _, p := range m.Projects {
		names = append(names, p.Name)
	}
	return names
}
