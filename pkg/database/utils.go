package database

import goqu "gopkg.in/doug-martin/goqu.v4"

// MustDataset panics if a dataset returns an error.
func MustDataset(ds *goqu.Dataset, err error) *goqu.Dataset {
	if err != nil {
		panic(err)
	}

	return ds
}
