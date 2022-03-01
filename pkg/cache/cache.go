package cache

import (
	"github.com/jwmatthews/case_watcher/pkg/api"
	"gorm.io/driver/sqlite" // Sqlite driver based on GGO
	"gorm.io/gorm"
	"log"
)

type Cache struct {
	db *gorm.DB
}

func Init(dbName string) (Cache, error) {
	var err error
	c := Cache{}
	c.db, err = gorm.Open(sqlite.Open(dbName), &gorm.Config{})
	if err != nil {
		log.Printf("Error opening db:  %s\n", err)
		return c, err
	}
	err = c.db.AutoMigrate(&Case{}, &Product{})
	if err != nil {
		log.Printf("Error migrating schema: %s\n", err)
		return c, err
	}
	return c, nil
}

func (c Cache) StoreCases(cases []api.Case) error {
	// Convert from the API representation of cases to our DB model
	myCases := c.ConvertToDBCases(cases)

	for _, tmpCase := range myCases {
		log.Printf("Attempting to save: %v\n", tmpCase)
		if c.db.Model(&tmpCase).Where("id = ?", tmpCase.Id).Updates(&tmpCase).RowsAffected == 0 {
			c.db.Create(&tmpCase)
			log.Printf("Saved new case: %s\n", tmpCase.Id)
		}
	}
	return nil
}

// ConvertToDBCases converts from the format returned from remote API
// to the format we will save to the database
func (c Cache) ConvertToDBCases(cases []api.Case) []Case {
	myCases := make([]Case, 0, len(cases))
	for _, apiCase := range cases {
		myCases = append(myCases, c.ConvertToDBCase(apiCase))
	}
	return myCases
}

func (c Cache) ConvertToDBCase(ac api.Case) Case {
	myCase := Case{}
	myCase.Id = ac.Id
	myCase.Uri = ac.Uri
	myCase.CreatedByName = ac.CreatedByName
	myCase.ContactName = ac.ContactName
	myCase.Version = ac.Version
	for _, p := range ac.Products {
		prod := Product{}
		// Look up an existing product of this name or create a new one
		c.db.Where("name = ?", p).FirstOrCreate(&prod, Product{Name: p})
		log.Printf("ConvertToDBCase:: Product = %v", prod)
	}
	return myCase
}
