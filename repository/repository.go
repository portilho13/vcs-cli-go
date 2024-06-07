package repository

type Repository struct {
	Name string
	LocalPath string
	RemotePath string
}

func (r *Repository) Init(Name string, LocalPath string, RemotePath string) (Repository) {
	return Repository{
		Name: Name,
		LocalPath: LocalPath,
		RemotePath: RemotePath,
	}
}