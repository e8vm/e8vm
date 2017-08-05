package dagvis

// Repo is the overview dependency structure of a repository.
type Repo struct {
	Name     string
	RepoTopo *M
	PkgTopos map[string]*M
}

// NewRepo creates an empty overview for a repo.
func NewRepo(name string) *Repo {
	return &Repo{
		Name:     name,
		PkgTopos: make(map[string]*M),
	}
}
