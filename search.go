package main

import (
	"log"

	"github.com/blevesearch/bleve"
)

type SearchIndex interface {
	Search(string) []string
	Put(string, File)
}

type BleveSeachIndex struct {
	Index bleve.Index
}

func GetIndex() BleveSeachIndex {
	indexpath := "jotting.searchindex"
	mapping := bleve.NewIndexMapping()
	index, err := bleve.New(indexpath, mapping)

	if err == bleve.ErrorIndexPathExists {
		index, err = bleve.Open(indexpath)
		if err != nil {
			log.Print(err)
		}
	}

	return BleveSeachIndex{
		Index: index,
	}
}

func (si *BleveSeachIndex) Put(idx string, f File) {
	si.Index.Index(idx, f)
}

func (si *BleveSeachIndex) Search(q string) []string {
	query := bleve.NewQueryStringQuery(q)
	search := bleve.NewSearchRequest(query)
	searchResults, err := si.Index.Search(search)

	if err != nil {
		log.Print(err)
	}

	var ids []string

	for _, val := range searchResults.Hits {
		ids = append(ids, val.ID)
	}

	return ids
}
