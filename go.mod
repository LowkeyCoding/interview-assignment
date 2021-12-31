module Test

go 1.17

replace (
	Logic => ./Logic
	Models => ./Models
)

require Logic v1.0.0

require (
	Models v1.0.0 // indirect
	github.com/jinzhu/inflection v1.0.0 // indirect
	github.com/jinzhu/now v1.1.3 // indirect
	github.com/mattn/go-sqlite3 v1.14.9 // indirect
	gorm.io/driver/sqlite v1.2.6 // indirect
	gorm.io/gorm v1.22.4 // indirect
)
