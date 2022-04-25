package controller

import (
	"errors"
	"io/ioutil"
	"log"

	"github.com/blevesearch/bleve/v2"
	"github.com/blevesearch/bleve/v2/mapping"
)

type controllerImpl struct {
	mapping *mapping.IndexMappingImpl
	index   bleve.Index
	indices map[string]bleve.Index
}

func New() (*controllerImpl, error) {
	mapping := bleve.NewIndexMapping()

	// Variable will keep a map of a pointer to all the open `index` files.
	indices := map[string]bleve.Index{}

	// Iterate through the contents of the 'db' folder (if it exists) and
	// get the names of the directories that are in. If the 'db' directory does
	// not exist then skip this whole step.
	files, err := ioutil.ReadDir("db")
	if err == nil {
		for _, file := range files {
			// Only open up directories.
			if file.IsDir() {
				// Get the name of the directory.
				directoryName := file.Name()

				// The following code will open up the bleve index.
				index, err := bleve.Open("db/" + directoryName)
				if err != nil {
					return nil, err
				}

				indices[directoryName] = index
				log.Println("Opened index:", directoryName)
			}
		}
	}

	return &controllerImpl{
		mapping: mapping,
		indices: indices,
	}, nil
}

func (c *controllerImpl) Register(filenames []string) error {
	for _, filename := range filenames {
		// Save everything into the db folder.
		filepath := "db/" + filename

		// Check if the filename already exists in our `indices` and if it
		// does then that means we have previously registered that filename
		// with an index so we can skip this loop and continue processing the
		// next filename.
		if _, ok := c.indices[filename]; ok {
			log.Println("Already registered index:", filename)
			continue
		}

		// Create an index for the particular filename.
		index, err := bleve.New(filepath, c.mapping)
		if err != nil {
			return err
		}
		// Open the index and attach it to our running application.
		index, err = bleve.Open(filepath)
		if err != nil {
			return err
		}
		c.indices[filename] = index
		log.Println("Registered index:", filename)
	}
	return nil
}

func (c *controllerImpl) Index(filename string, identifier string, data interface{}) error {
	tenantIndex, ok := c.indices[filename]
	if !ok {
		return errors.New("index D.N.E. for filename")
	}
	return tenantIndex.Index(identifier, data)
}

func (c *controllerImpl) Query(filename string, search string) ([]string, error) {
	tenantIndex, ok := c.indices[filename]
	if !ok {
		return nil, errors.New("index D.N.E. for filename")
	}

	query := bleve.NewQueryStringQuery(search)
	searchRequest := bleve.NewSearchRequest(query)
	searchResults, err := tenantIndex.Search(searchRequest)

	var arr []string
	for _, v := range searchResults.Hits {
		arr = append(arr, v.ID)
	}

	return arr[:], err
}

func (c *controllerImpl) Close() {
	for key, tenantIndex := range c.indices {
		tenantIndex.Close()
		log.Println("Closed index:", key)
	}
}
