package Models

import (
	"bufio"
	"fmt"
	"log"
	"os"
)

type User struct {
	ID        string `gorm:"column:id"`
	Firstname string `gorm:"column:firstName"`
	Lastname  string `gorm:"column:lastName"`
	Email     string `gorm:"column:email"`
}

func (user User) ToFile(path string) int {
	f, err := os.OpenFile(path,
		os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	w := bufio.NewWriter(f)
	if err != nil {
		log.Println(err)
		return 0
	}
	bytesWritten, err := fmt.Fprintf(w, "%s,%s,%s,%s\n", user.ID,
		user.Firstname,
		user.Lastname,
		user.Email)
	if err != nil {
		log.Println(err)
		return 0
	}
	if err := w.Flush(); err != nil {
		log.Println(err)
		return 0
	}
	if err := f.Close(); err != nil {
		log.Println(err)
		return 0
	}
	return bytesWritten
}

func (a User) Compare(b User) bool {
	if a.ID != b.ID {
		return false
	}
	if a.Firstname != b.Firstname {
		return false
	}
	if a.Lastname != b.Lastname {
		return false
	}
	if a.Email != b.Email {
		return false
	}
	return true
}
