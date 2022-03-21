package cache

import (
	"gorm.io/gorm"
	"time"
)

type Product struct {
	gorm.Model
	Name   string
	CaseId string
}

type Case struct {
	AccountNumber        string
	CaseNumber           string
	ContactName          string
	CreatedByName        string
	CreatedDate          time.Time
	CustomerEscalation   bool
	Id                   string `gorm:"primary_key"`
	LastModifiedByName   string
	LastModifiedDate     time.Time
	LastPublicUpdateBy   string
	LastPublicUpdateDate time.Time
	Owner                string
	Products             []Product
	Severity             string
	Summary              string
	Status               string
	Type                 string
	Uri                  string
	Version              string
}

type Account struct {
	AccountNumber  string `gorm:"primaryKey"`
	GSCSMSegment   string
	Name           string
	CSMUserID      string
	CSMUserName    string
	CSMUserSSOName string
	Strategic      bool
	HasEnhancedSLA bool
	HasSRM         bool
	HasTAM         bool
}
