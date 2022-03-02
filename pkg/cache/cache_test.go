package cache

import (
	"fmt"
	"github.com/jwmatthews/case_watcher/pkg/api"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"os"
	"testing"
)

const (
	dbName = "unit_tests.db"
)

var (
	myCache Cache
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
	t.Log("Running TestStoreCases")
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
	err = myCache.DB.Where(&Case{Id: "myid1"}).First(&case2).Error
	require.NoError(t, err)
	assert.Equal(t, len(case2.Products), 3, "Number of products should equal to test data created")

}

func TestMain(m *testing.M) {
	var err error
	// myCache will be used by all tests in this package
	myCache, err = Init(dbName)
	if err != nil {
		fmt.Printf("Failed to initiative database: %s\n", err)
	}
	os.Exit(m.Run())
}