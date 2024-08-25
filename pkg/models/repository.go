package models

type Repository struct {
	Name     string `json:"name"`
	CloneURL string `json:"clone_url"`
	Private  bool   `json:"private"`
	Owner    struct {
		Login string `json:"login"`
	} `json:"owner"`
}

func NewRepository(name, cloneURL string, private bool, ownerLogin string) *Repository {
	return &Repository{
		Name:     name,
		CloneURL: cloneURL,
		Private:  private,
		Owner: struct {
			Login string `json:"login"`
		}{
			Login: ownerLogin,
		},
	}
}
