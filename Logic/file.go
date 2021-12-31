package Logic

import (
	"Models"
	"log"
	"os"
	"strings"
)

func ReadLinesFromFile(filePath string, offset int64) string {
	file, err := os.Open(filePath)
	if err != nil {
		log.Panic("Failed to open file: " + err.Error())
	}
	defer func(file *os.File) {
		var err = file.Close()
		if err != nil {
			log.Panic("Failed to close file: " + err.Error())
		}
	}(file)
	buf := make([]byte, offset)
	stat, err := os.Stat(filePath)
	start := stat.Size() - offset
	_, err = file.ReadAt(buf, start)
	if err != nil {
		log.Panic("Failed to read file: " + err.Error())
	}
	return string(buf)
}

func WriteUsersToFile(users []Models.User, path string) int64 {
	var bytesWritten int64
	for _, user := range users {
		bytesWritten += int64(user.ToFile(path))
	}
	return bytesWritten
}

func ReadUsersFromFile(offset int64, path string) []Models.User {
	buf := ReadLinesFromFile(path, offset)
	lines := strings.Split(buf, "\n")
	var fileUsers []Models.User
	if len(lines) > 1 {
		lines = lines[:len(lines)-1]
		for _, line := range lines {
			user := UserFromString(line)
			if user.ID != "" {
				fileUsers = append(fileUsers, user)
			}
		}
	}
	return fileUsers
}

func UserFromString(input string) Models.User {
	fields := strings.Split(input, ",")
	var user Models.User
	if len(fields) == 4 {
		user.ID = strings.TrimSpace(fields[0])
		user.Firstname = strings.TrimSpace(fields[1])
		user.Lastname = strings.TrimSpace(fields[2])
		user.Email = strings.TrimSpace(fields[3])
		return user
	}
	log.Panic("Failed convert string to customer: Invalid field count (?) required 4", len(fields))
	return user
}

func ValidateUsersInFile(fromFile []Models.User, fromDB []Models.User) bool {
	if len(fromFile) != len(fromDB) {
		return false
	}
	for i := range fromFile {
		if !fromDB[i].Compare(fromFile[i]) {
			return false
		}
	}
	return true
}
