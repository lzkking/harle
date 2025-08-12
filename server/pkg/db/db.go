package db

import "gorm.io/gorm"

var Client = newDBClient()

func Session() *gorm.DB {
	return Client.Session(&gorm.Session{
		FullSaveAssociations: true,
	})
}
