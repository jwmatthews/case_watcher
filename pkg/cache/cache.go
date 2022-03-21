package cache

import (
	"github.com/jwmatthews/case_watcher/pkg/api"
	"gorm.io/driver/sqlite" // Sqlite driver based on GGO
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"log"
	"os"
	"time"
)

type Cache struct {
	DB *gorm.DB
}

func Init(dbName string) (Cache, error) {
	var err error
	c := Cache{}
	newLogger := logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags), // io writer
		logger.Config{
			SlowThreshold: time.Second, // Slow SQL threshold
			//LogLevel:                  logger.Info, // Log level
			LogLevel:                  logger.Warn, // Log level
			IgnoreRecordNotFoundError: true,        // Ignore ErrRecordNotFound error for logger
			Colorful:                  true,        // Disable color
		},
	)

	c.DB, err = gorm.Open(sqlite.Open(dbName),
		&gorm.Config{
			Logger: newLogger,
		})
	if err != nil {
		log.Printf("Error opening db:  %s\n", err)
		return c, err
	}
	err = c.DB.AutoMigrate(&Case{}, &Product{}, &Account{})
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
		err := c.StoreCase(tmpCase)
		if err != nil {
			log.Printf("Error storing case '%s':  %s", tmpCase.Id, err)
			return err
		}
	}
	return nil
}

func (c Cache) StoreCase(myCase Case) error {
	log.Printf("Attempting to save: %v\n", myCase)
	if c.DB.Model(&myCase).Where("id = ?", myCase.Id).Updates(&myCase).RowsAffected == 0 {
		c.DB.Create(&myCase)
		log.Printf("Saved new case: %s\n", myCase.Id)
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

// ConvertToDBCase converts from api format to DB
// Ideally I would have 1 structure type and reuse for both
// I *think* I need the 2 different formats to handle some subtle
// differences with JSON vs Gorm modeling of Products, but I am not 100% certain
// TODO: See if there is a cleaner way to refactor to a single reused structure definition between incoming API data and DB
//
func (c Cache) ConvertToDBCase(ac api.Case) Case {
	myCase := Case{}
	myCase.AccountNumber = ac.AccountNumber
	myCase.CaseNumber = ac.CaseNumber
	myCase.ContactName = ac.ContactName
	myCase.CreatedByName = ac.CreatedByName
	myCase.CreatedDate = ac.CreatedDate
	myCase.CustomerEscalation = ac.CustomerEscalation
	myCase.Id = ac.Id
	myCase.LastModifiedByName = ac.LastModifiedByName
	myCase.LastModifiedDate = ac.LastModifiedDate
	myCase.LastPublicUpdateDate = ac.LastPublicUpdateDate
	myCase.LastPublicUpdateBy = ac.LastPublicUpdateBy
	myCase.Owner = ac.Owner
	// Products
	for _, p := range ac.Products {
		prod := Product{}
		if c.DB.Where(&Product{Name: p, CaseId: ac.Id}).First(&prod).RowsAffected == 0 {
			prod.Name = p
		}
		myCase.Products = append(myCase.Products, prod)
	}
	myCase.Severity = ac.Severity
	myCase.Summary = ac.Summary
	myCase.Status = ac.Status
	myCase.Type = ac.Type
	myCase.Uri = ac.Uri
	myCase.Version = ac.Version
	return myCase
}

// GetMissingAccountIDs will return a slice of account ids we lack details on
func (c Cache) GetMissingAccountIDs() []string {
	// Find all Account IDs from all Cases saved.
	accountNums := make([]string, 0)
	result := make([]string, 0)
	c.DB.Raw("SELECT account_number FROM cases WHERE 'account_number' is NOT NULL").Scan(&accountNums)
	for i := range accountNums {
		if accountNums[i] != "" {
			result = append(result, accountNums[i])
		}
	}
	return result
	//c.DB.Table("cases").Select("account_number").Scan(&accountNums)
	// Find all Account IDs which do not have an entry in the Database
}

func (c Cache) GetAllCases() ([]Case, error) {
	cases := make([]Case, 0)
	err := c.DB.Where(&Case{}).Find(&cases).Error
	if err != nil {
		return []Case{}, err
	}
	return cases, nil
}

func (c Cache) GetOpenCases() ([]Case, error) {
	cases := make([]Case, 0)
	err := c.DB.Where("status != 'Closed'").Find(&cases).Error
	if err != nil {
		return []Case{}, err
	}
	return cases, nil
}

func (c Cache) GetClosedCases() ([]Case, error) {
	cases := make([]Case, 0)
	err := c.DB.Where("status == 'Closed'").Find(&cases).Error
	if err != nil {
		return []Case{}, err
	}
	return cases, nil
}

func (c Cache) GetCasesActiveFrom(since time.Time) ([]Case, error) {
	cases := make([]Case, 0)
	err := c.DB.Where("last_modified_date >= ?", since).Find(&cases).Error
	if err != nil {
		return []Case{}, err
	}
	return cases, nil
}

func (c Cache) GetUniqueCaseStatusValues() ([]string, error) {
	results := make([]string, 0)
	err := c.DB.Model(&Case{}).Distinct("status").Order("status asc").Find(&results).Error
	if err != nil {
		return []string{}, err
	}
	return results, nil

}
