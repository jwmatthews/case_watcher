package cache

import (
	"fmt"
	"github.com/jwmatthews/case_watcher/pkg/api"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"os"
	"testing"
	"time"
)

// Desired test
// - Create a Case with several Products and ensure we can query and get correct Products
// - Update a Case 3+ times, ensure no extra Products are created
// - Update a Case and remove a Product
// - Update a Case and add a Product
// - Delete a Case and see the associated Product mappings are deletedgi
const (
	dbName = "unit_tests.db"
)

func getSampleCase() api.Case {
	products := make([]string, 0)
	for i := 0; i < 3; i++ {
		products = append(products, fmt.Sprintf("TestName%d", i))
	}
	myCase := api.Case{Id: "myid1", Summary: "My summary", ContactName: "bob smith", Products: products}
	return myCase
}

func TestStoreCases(t *testing.T) {
	myCache := InitCache(t, dbName)
	defer CleanUpDB(dbName)

	myCases := make([]api.Case, 0)
	myCase := getSampleCase()
	myCases = append(myCases, myCase)
	err := myCache.StoreCases(myCases)
	require.NoError(t, err)

	product := Product{}
	err = myCache.DB.Where(&Product{Name: "TestName2", CaseId: "myid1"}).First(&product).Error
	require.NoError(t, err)
	assert.Equal(t, product.Name, "TestName2", "Product name should match what was queried")
	assert.Equal(t, product.CaseId, "myid1", "CaseId should match what was queried")

	case2 := Case{}
	err = myCache.DB.Preload("Products").Where(&Case{Id: "myid1"}).First(&case2).Error
	require.NoError(t, err)
	assert.Equal(t, 3, len(case2.Products), "Number of products should equal to test data created")

}

// TestStoreCasesCalledMultipleTimes will ensure when we update a Case
// with Products we do not create duplicate Product entries
func TestStoreCasesCalledMultipleTimes(t *testing.T) {
	myCache := InitCache(t, dbName)
	defer CleanUpDB(dbName)

	myCases := make([]api.Case, 0)
	myCase := getSampleCase()
	myCases = append(myCases, myCase)

	// Call #1
	err := myCache.StoreCases(myCases)
	require.NoError(t, err)
	// Call #2
	err = myCache.StoreCases(myCases)
	require.NoError(t, err)
	// Call #3
	err = myCache.StoreCases(myCases)
	require.NoError(t, err)

	cases := make([]Case, 0)
	result := myCache.DB.Where(&Case{Id: "myid1"}).Find(&cases)
	assert.Equal(t, int64(1), result.RowsAffected)
	assert.Equal(t, 1, len(cases))

	case2 := Case{}
	err = myCache.DB.Preload("Products").Where(&Case{Id: "myid1"}).First(&case2).Error
	require.NoError(t, err)
	assert.Equal(t, 3, len(case2.Products), "Number of products should equal to test data created")

	products := make([]Product, 0)
	result = myCache.DB.Where(&Product{CaseId: "myid1"}).Find(&products)
	assert.Equal(t, int64(3), result.RowsAffected)
}

func TestGetMissingAccountIDs(t *testing.T) {
	myCache := InitCache(t, dbName)
	defer CleanUpDB(dbName)
	// Store 2 new Cases with Accounts
	myCase1 := Case{Id: "case1", AccountNumber: "1"}
	myCase2 := Case{Id: "case2", AccountNumber: "2"}
	myCase3 := Case{Id: "case3", AccountNumber: ""} // intentional empty account number
	err := myCache.StoreCase(myCase1)
	require.NoError(t, err)
	err = myCache.StoreCase(myCase2)
	require.NoError(t, err)
	err = myCache.StoreCase(myCase3)
	require.NoError(t, err)

	foundAccountIDs := myCache.GetMissingAccountIDs()
	t.Log("foundAccountIDs = ", foundAccountIDs)
	assert.Equal(t, 2, len(foundAccountIDs))
	assert.Contains(t, foundAccountIDs, "1")
	assert.Contains(t, foundAccountIDs, "2")

	// Next is to save an Account, and ensure that the
	// saved Account is not included with missing account IDs
}

func TestCache_GetAllCases(t *testing.T) {
	myCache := InitCache(t, dbName)
	defer CleanUpDB(dbName)
	myCase1 := Case{Id: "case1", AccountNumber: "1"}
	myCase2 := Case{Id: "case2", AccountNumber: "2"}
	myCase3 := Case{Id: "case3", AccountNumber: "3"}
	err := myCache.StoreCase(myCase1)
	require.NoError(t, err)
	err = myCache.StoreCase(myCase2)
	require.NoError(t, err)
	err = myCache.StoreCase(myCase3)
	require.NoError(t, err)

	foundCases, err := myCache.GetAllCases()
	require.NoError(t, err)
	t.Log("foundCases = ", foundCases)
	assert.Equal(t, 3, len(foundCases))
}

func TestCache_GetOpenCases(t *testing.T) {
	myCache := InitCache(t, dbName)
	defer CleanUpDB(dbName)
	myCase1 := Case{Id: "case1", AccountNumber: "1", Status: "Waiting on Customer"}
	myCase2 := Case{Id: "case2", AccountNumber: "2", Status: "Unknown"}
	myCase3 := Case{Id: "case3", AccountNumber: "3", Status: "Closed"}
	myCase4 := Case{Id: "case4", AccountNumber: "4", Status: "Closed"}
	myCase5 := Case{Id: "case5", AccountNumber: "5", Status: "Closed"}
	err := myCache.StoreCase(myCase1)
	require.NoError(t, err)
	err = myCache.StoreCase(myCase2)
	require.NoError(t, err)
	err = myCache.StoreCase(myCase3)
	require.NoError(t, err)
	err = myCache.StoreCase(myCase4)
	require.NoError(t, err)
	err = myCache.StoreCase(myCase5)
	require.NoError(t, err)

	foundCases, err := myCache.GetOpenCases()
	require.NoError(t, err)
	t.Log("foundCases = ", foundCases)
	assert.Equal(t, 2, len(foundCases))
}

func TestCache_GetCasesActiveFrom(t *testing.T) {
	myCache := InitCache(t, dbName)
	defer CleanUpDB(dbName)
	myCase1 := Case{Id: "case1", AccountNumber: "1", Status: "Waiting on Customer", LastModifiedDate: time.Now().AddDate(0, 0, -1)}
	myCase2 := Case{Id: "case2", AccountNumber: "2", Status: "Unknown", LastModifiedDate: time.Now().AddDate(0, 0, -25)}
	myCase3 := Case{Id: "case3", AccountNumber: "3", Status: "Closed", LastModifiedDate: time.Now().AddDate(0, 0, -3)}
	myCase4 := Case{Id: "case4", AccountNumber: "4", Status: "Closed", LastModifiedDate: time.Now().AddDate(0, -1, 0)}
	myCase5 := Case{Id: "case5", AccountNumber: "5", Status: "Closed", LastModifiedDate: time.Now().AddDate(-1, 0, 0)}
	err := myCache.StoreCase(myCase1)
	require.NoError(t, err)
	err = myCache.StoreCase(myCase2)
	require.NoError(t, err)
	err = myCache.StoreCase(myCase3)
	require.NoError(t, err)
	err = myCache.StoreCase(myCase4)
	require.NoError(t, err)
	err = myCache.StoreCase(myCase5)
	require.NoError(t, err)

	foundCases, err := myCache.GetCasesActiveFrom(time.Now().AddDate(0, 0, -7))
	require.NoError(t, err)
	t.Log("foundCases = ", foundCases)
	assert.Equal(t, 2, len(foundCases))
	for _, c := range foundCases {
		// We expect only 'case1' and 'case3' have been active in the past week
		assert.Contains(t, []string{"case1", "case3"}, c.Id)
	}
}

func TestCache_GetClosedCases(t *testing.T) {
	myCache := InitCache(t, dbName)
	defer CleanUpDB(dbName)
	myCase1 := Case{Id: "case1", AccountNumber: "1", Status: "Waiting on Customer"}
	myCase2 := Case{Id: "case2", AccountNumber: "2", Status: "Unknown"}
	myCase3 := Case{Id: "case3", AccountNumber: "3", Status: "Closed"}
	myCase4 := Case{Id: "case4", AccountNumber: "4", Status: "Closed"}
	myCase5 := Case{Id: "case5", AccountNumber: "5", Status: "Closed"}
	err := myCache.StoreCase(myCase1)
	require.NoError(t, err)
	err = myCache.StoreCase(myCase2)
	require.NoError(t, err)
	err = myCache.StoreCase(myCase3)
	require.NoError(t, err)
	err = myCache.StoreCase(myCase4)
	require.NoError(t, err)
	err = myCache.StoreCase(myCase5)
	require.NoError(t, err)

	foundCases, err := myCache.GetClosedCases()
	require.NoError(t, err)
	t.Log("foundCases = ", foundCases)
	assert.Equal(t, 3, len(foundCases))
	for _, c := range foundCases {
		// We expect only 'case3', 'case4', 'case5' to be Closed
		assert.Contains(t, []string{"case3", "case4", "case5"}, c.Id)
	}
}

func TestCache_GetUniqueCaseStatusValues(t *testing.T) {
	myCache := InitCache(t, dbName)
	defer CleanUpDB(dbName)
	myCase1 := Case{Id: "case1", AccountNumber: "1", Status: "Waiting on Customer"}
	myCase2 := Case{Id: "case2", AccountNumber: "2", Status: "Unknown"}
	myCase3 := Case{Id: "case3", AccountNumber: "3", Status: "Foo"}
	myCase4 := Case{Id: "case4", AccountNumber: "4", Status: "Foo"}
	myCase5 := Case{Id: "case5", AccountNumber: "5", Status: "Closed"}
	myCase6 := Case{Id: "case6", AccountNumber: "6", Status: "Closed"}
	err := myCache.StoreCase(myCase1)
	require.NoError(t, err)
	err = myCache.StoreCase(myCase2)
	require.NoError(t, err)
	err = myCache.StoreCase(myCase3)
	require.NoError(t, err)
	err = myCache.StoreCase(myCase4)
	require.NoError(t, err)
	err = myCache.StoreCase(myCase5)
	require.NoError(t, err)
	err = myCache.StoreCase(myCase6)
	require.NoError(t, err)

	foundValues, err := myCache.GetUniqueCaseStatusValues()
	require.NoError(t, err)
	t.Log("foundValues = ", foundValues)
	assert.Equal(t, 4, len(foundValues))
	for _, v := range foundValues {
		assert.Contains(t, []string{"Waiting on Customer", "Foo", "Unknown", "Closed"}, v)
	}
}

func CleanUpDB(dbName string) error {
	_, err := os.Stat(dbName)
	if err != nil {
		fmt.Printf("Error looking up test db file, %s: %s\n", dbName, err)
		return err
	}
	err = os.Remove(dbName)
	if err != nil {
		fmt.Printf("Error removing test db file, %s: %s\n", dbName, err)
		return err
	}
	return nil
}

func InitCache(t *testing.T, dbName string) *Cache {
	myCache, err := Init(dbName)
	if err != nil {
		t.Fatalf("Failed to initiative database: %s\n", err)
	}
	return &myCache
}
