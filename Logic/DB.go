package Logic

import (
	"Models"
	"flag"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"log"
	"os"
	"regexp"
	"strings"
	"time"
	"unicode"
)

func SetupLogger() logger.Interface {
	return logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags), // io writer
		logger.Config{
			SlowThreshold:             time.Second,   // Slow SQL threshold
			LogLevel:                  logger.Silent, // Log level
			IgnoreRecordNotFoundError: true,          // Ignore ErrRecordNotFound error for logger
			Colorful:                  false,         // Disable color
		},
	)
}

func ConnectToDatabase(connectionString *string) *gorm.DB {
	db, err := gorm.Open(sqlite.Open(*connectionString), &gorm.Config{
		Logger: SetupLogger(),
	})

	if err != nil {
		panic("failed to connect database")
	}

	return db
}

func CloseConnection(db *gorm.DB) {
	rawConnection, err := db.DB()
	if err != nil {
		log.Fatalln(err)
	}
	defer rawConnection.Close()
}

func GetArguments(path *string, connectionString *string, limit *int, query *string) {
	flag.IntVar(limit, "limit", 1, "Limits the amount of entries pulled from the database. Use 0 for no limit (DEFAULT: 1)")
	flag.StringVar(query, "query", "", "The query used to select data from the database (DEFAULT \"SELECT * From Users\")")
	flag.StringVar(path, "path", "data.out", "The path to file containing deleted users from the database (DEFAULT: \"data.out\")")
	flag.StringVar(connectionString, "db", "sqldump.db", "The database connection string (DEFAULT: \"sqldump.db\")")
	ArgsPtr := flag.String("args", "", "Arguments replace question marks in queries. Lists are separated by commas and arguments are seperated by semi colons (NO DEFAULT)")
	flag.Parse()
	*query = GenerateQuery(*query, *ArgsPtr)
}

func GenerateQuery(query string, rawArgs string) string {
	args := strings.Split(rawArgs, ";")
	reStr := regexp.MustCompile("^(.*?)\\?(.*)$")
	for _, arg := range args {
		var repStr string
		if IsArgInt(arg) {
			repStr = arg
		} else {
			argValues := strings.Split(arg, ",")
			for i, value := range argValues {

				repStr += "\"" + value + "\""
				if i != len(argValues)-1 {
					repStr += ","
				}
			}
		}
		if strings.ContainsAny(repStr, ",") {
			repStr = "${1}(" + repStr + ")$2"
		} else {
			repStr = "${1}" + repStr + "$2"
		}
		query = reStr.ReplaceAllString(query, repStr)
	}
	return query
}

func IsArgInt(s string) bool {
	for _, c := range s {
		if !unicode.IsDigit(c) && c != ',' {
			return false
		}
	}
	return true
}

func ExecuteQuery(db *gorm.DB, limit *int, query *string) []Models.User {
	var users []Models.User
	if *limit > 0 {
		db.Limit(*limit).Find(&users, *query)
	} else {
		db.Find(&users, *query)
	}
	return users
}

func DeleteUsers(users []Models.User, db *gorm.DB) error {
	tx := db.Begin()

	// In case something goes wrong rollback
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	if err := tx.Delete(&users).Error; err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit().Error
}
