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
	Id                   string `gorm:"primary_key"`
	Uri                  string
	CreatedByName        string
	ContactName          string
	Version              string
	Products             []Product
	Number               string
	LastPublicUpdateBy   string
	Severity             string
	Owner                string
	LastPublicUpdateDate time.Time
	CreatedDate          time.Time
	Summary              string
	LastModifiedDate     time.Time
	AccountNumber        string
	Type                 string
	LastModifiedByName   string
	CustomerEscalation   bool
	Status               string
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
