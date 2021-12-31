package main

import (
	. "Logic"
	"fmt"
	"log"
)

func main() {
	var limit int
	var query string
	var path string
	var connectionString string

	GetArguments(&path, &connectionString, &limit, &query)

	db := ConnectToDatabase(&connectionString)

	usersFromDB := ExecuteQuery(db, &limit, &query)

	bytesOffset := WriteUsersToFile(usersFromDB, path)

	usersFromFile := ReadUsersFromFile(bytesOffset, path)

	if ValidateUsersInFile(usersFromFile, usersFromDB) {
		err := DeleteUsers(usersFromDB, db)
		if err != nil {
			log.Panic(err)
		}
	} else {
		log.Panic("Failed to validate users in file")
	}
	fmt.Printf("A total of %d users has been written to %s", len(usersFromDB), path)

	CloseConnection(db)
}
