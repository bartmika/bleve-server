package controller

import (
	"errors"
	"io/ioutil"

	"github.com/rs/zerolog"
	"github.com/blevesearch/bleve/v2"
	"github.com/blevesearch/bleve/v2/mapping"
)

type controllerImpl struct {
	Logger        zerolog.Logger
	directoryPath string
	mapping       *mapping.IndexMappingImpl
	index         bleve.Index
	indices       map[string]bleve.Index
}

func New(directoryPath string, logger zerolog.Logger) (*controllerImpl, error) {
	mapping := bleve.NewIndexMapping()

	// Variable will keep a map of a pointer to all the open `index` files.
	indices := map[string]bleve.Index{}

	// Iterate through the contents of the 'db' folder (if it exists) and
	// get the names of the directories that are in. If the 'db' directory does
	// not exist then skip this whole step.
	files, err := ioutil.ReadDir(directoryPath)
	if err == nil {
		for _, file := range files {
			// Only open up directories.
			if file.IsDir() {
				// Get the name of the directory.
				directoryName := file.Name()

				// The following code will open up the bleve index.
				index, err := bleve.Open(directoryPath + "/" + directoryName)
				if err != nil {
					// Skip the file since it's not a "bleve" file we can
					// use for our application.
					continue
				}

				indices[directoryName] = index
				logger.Info().Msgf("opened index: %s", directoryName)
			}
		}
	}

	return &controllerImpl{
		Logger:        logger,
		directoryPath: directoryPath,
		mapping:       mapping,
		indices:       indices,
	}, nil
}

func (c *controllerImpl) Register(filenames []string) error {
	for _, filename := range filenames {
		// Save everything into the db folder.
		filepath := c.directoryPath + "/" + filename

		// Check if the filename already exists in our `indices` and if it
		// does then that means we have previously registered that filename
		// with an index so we can skip this loop and continue processing the
		// next filename.
		if _, ok := c.indices[filename]; ok {
			continue
		}

		// Create an index for the particular filename.
		index, err := bleve.New(filepath, c.mapping)
		if err != nil {
			return err
		}

		c.indices[filename] = index
		c.Logger.Info().Msgf("new index: %s", filename)
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

	// DEVELOPERS NOTE: Increase the limit from `10` to `1 million` so in essence
	// we return unrestricted number of entries.
	searchRequest := bleve.NewSearchRequestOptions(query, 1_000_000_000, 0, false)
	searchResults, err := tenantIndex.Search(searchRequest)

	var arr []string
	for _, v := range searchResults.Hits {
		arr = append(arr, v.ID)
	}

	return arr[:], err
}

func (c *controllerImpl) Delete(filename string, identifier string) error {
	tenantIndex, ok := c.indices[filename]
	if !ok {
		return errors.New("index D.N.E. for filename")
	}
	return tenantIndex.Delete(identifier)
}

func (c *controllerImpl) Close() {
	for key, tenantIndex := range c.indices {
		tenantIndex.Close()
		c.Logger.Info().Msgf("closed index: %s", key)
	}
}
