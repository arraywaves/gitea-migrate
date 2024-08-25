package models

type Repository struct {
	Name     string `json:"name"`
	CloneURL string `json:"clone_url"`
	Private  bool   `json:"private"`
	Owner    string `json:"owner"`
}

func NewRepository(name, cloneURL, owner string, private bool) *Repository {
	return &Repository{
		Name:     name,
		CloneURL: cloneURL,
		Private:  private,
		Owner:    owner,
	}
}
