package report

import (
	"fmt"
	"github.com/jwmatthews/case_watcher/pkg/cache"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"os"
	"testing"
	"time"
)

const (
	dbName = "unit_tests.db"
)

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

func InitCache(t *testing.T, dbName string) *cache.Cache {
	myCache, err := cache.Init(dbName)
	if err != nil {
		t.Fatalf("Failed to initiative database: %s\n", err)
	}
	return &myCache
}

func TestReport_GetSubjectLine(t *testing.T) {
	myCache := InitCache(t, dbName)
	defer CleanUpDB(dbName)

	r := GetReport(myCache, "myspreadsheetID")
	subjLine := r.GetSubjectLine()
	now := time.Now().Format("2006-01-02")
	assert.Contains(t, subjLine, now)
}

func TestReport_GetSpreadsheetURL(t *testing.T) {
	myCache := InitCache(t, dbName)
	defer CleanUpDB(dbName)

	sID := "spreadsheetIDX1243434"
	r := GetReport(myCache, sID)
	foundURL := r.GetSpreadsheetURL()
	assert.Contains(t, foundURL, sID)
}

func TestReport_ToHTMLWithEmptyCache(t *testing.T) {
	myCache := InitCache(t, dbName)
	defer CleanUpDB(dbName)

	sID := "spreadsheetIDX1243434"
	r := GetReport(myCache, sID)
	html := r.ToHTML()
	assert.NotContains(t, html, "Error")
}

func TestReport_GetActiveCasesFromWithEmptyCache(t *testing.T) {
	myCache := InitCache(t, dbName)
	defer CleanUpDB(dbName)

	sID := "spreadsheetIDX1243434"
	r := GetReport(myCache, sID)
	cases, err := r.GetActiveCasesFrom(time.Now())
	require.NoError(t, err)
	assert.Len(t, cases, 0)
}
