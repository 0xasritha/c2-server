package main

type ResultManager interface {
	getResults() []Result
	addResult(result Result) error // TODO: rename?
}

type Result struct {
	uuid         uint32
	Contents     map[string]string `json:"contents"`
	IsSuccessful bool              `json:"isSuccessful"`
}

// TODO: rename this to LocalResultsManager? or think of other names
type LocalResultsStore struct {
	results []Result
}

func NewLocalResultsStore() *LocalResultsStore {
	lrs := LocalResultsStore{
		results: make([]Result, 0),
	}
	return &lrs
}

func (lrs LocalResultsStore) getResults() []Result {
	return lrs.results
}

func (lrs LocalResultsStore) addResult(result Result) error {
	lrs.results = append(lrs.results, result)
	return nil // TODO: do we need to return an error in this one + Task one
}

type SQLResultsStore struct{}
