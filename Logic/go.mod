module Logic

go 1.17

require (
	Models v1.0.0
	gorm.io/driver/sqlite v1.2.6
	gorm.io/gorm v1.22.4
)

require (
	github.com/jinzhu/inflection v1.0.0 // indirect
	github.com/jinzhu/now v1.1.3 // indirect
	github.com/mattn/go-sqlite3 v1.14.9 // indirect
)

replace Models => ../Models
