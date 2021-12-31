package Main

import (
	. "Logic"
	"Models"
	"fmt"
	"io"
	"os"
	"testing"
)

func setup() {
	deleteFile("temp.db")
	deleteFile("data.out")
	_, err := copyFile("sqldump.db", "temp.db")
	if err != nil {
		return
	}
	createFile("data.out")
}

func copyFile(src, dst string) (int64, error) {
	sourceFileStat, err := os.Stat(src)
	if err != nil {
		return 0, err
	}

	if !sourceFileStat.Mode().IsRegular() {
		return 0, fmt.Errorf("%s is not a regular file", src)
	}

	source, err := os.Open(src)
	if err != nil {
		return 0, err
	}
	defer source.Close()

	destination, err := os.Create(dst)
	if err != nil {
		return 0, err
	}
	defer destination.Close()
	nBytes, err := io.Copy(destination, source)
	return nBytes, err
}

func deleteFile(path string) bool {
	var _, err = os.Stat(path)

	// create file if not exists
	if !os.IsNotExist(err) {
		var err = os.Remove(path)
		if err != nil {
			fmt.Printf("Failed to delete file %s", path)
			return false
		}
	}
	return true
}

func createFile(path string) bool {
	// check if file exists
	var _, err = os.Stat(path)

	// create file if not exists
	if os.IsNotExist(err) {
		var file, err = os.Create(path)
		if err != nil {
			fmt.Printf("Failed to create file %s", path)
			return false
		}
		defer file.Close()
	}

	fmt.Println("File Created Successfully", path)
	return true
}

func TestMain(m *testing.M) {
	setup()
	code := m.Run()
	os.Exit(code)
}

func TestQueryGenerator(t *testing.T) {
	queryA := "ID = ?"
	argsA := "42"
	queryAExpected := "ID = 42"
	resultA := GenerateQuery(queryA, argsA)
	if resultA != queryAExpected {
		t.Errorf("Got query (%s) expected (%s)", resultA, queryAExpected)
	}

	argsB := "E42"
	queryBExpected := "ID = \"E42\""
	resultB := GenerateQuery(queryA, argsB)
	if resultB != queryBExpected {
		t.Errorf("Got query (%s) expected (%s)", resultB, queryBExpected)
	}

	queryB := "ID = ? OR ?"
	argsC := "42;55"
	queryCExpected := "ID = 42 OR 55"
	resultC := GenerateQuery(queryB, argsC)
	if resultC != queryCExpected {
		t.Errorf("Got query (%s) expected (%s)", resultC, queryCExpected)
	}

	argsD := "E42;D55"
	queryDExpected := "ID = \"E42\" OR \"D55\""
	resultD := GenerateQuery(queryB, argsD)
	if resultD != queryDExpected {
		t.Errorf("Got query (%s) expected (%s)", resultD, queryDExpected)
	}

	queryC := "ID = ? OR ?"
	argsE := "42,55;64,128"
	queryEExpected := "ID = (42,55) OR (64,128)"
	resultE := GenerateQuery(queryC, argsE)
	if resultE != queryEExpected {
		t.Errorf("Got query (%s) expected (%s)", resultE, queryEExpected)
	}
	argsF := "E42,D55;F64,G128"
	queryFExpected := "ID = (\"E42\",\"D55\") OR (\"F64\",\"G128\")"
	resultF := GenerateQuery(queryC, argsF)
	if resultF != queryFExpected {
		t.Errorf("Got query (%s) expected (%s)", resultF, queryFExpected)
	}

}

func TestCompareUsers(t *testing.T) {
	var a Models.User
	var b Models.User
	var c Models.User

	a.ID = "0D00A443-93D3-8573-135D-8946397866A1"
	a.Firstname = "Hayes"
	a.Lastname = "Cruz"
	a.Email = "dui@aliquetsem.org"

	b.ID = "352AFADF-0A34-7933-7196-294A3AEEA6CE"
	b.Firstname = "Addison"
	b.Lastname = "Mcdonald"
	b.Email = "consectetuer.adipiscing.elit@Duissit.co.uk"

	c.ID = "352AFADF-0A34-7933-7196-294A3AEEA6CE"
	c.Firstname = "Addison"
	c.Lastname = "Mcdonald"
	c.Email = "consectetuer.adipiscing.elit@Duissit.co.uk"

	if a.Compare(a) == false {
		t.Errorf("User A should be equivilant to user A")
	}
	if a.Compare(b) == true {
		t.Errorf("User A should not be equivilant to user B")
	}
	if b.Compare(c) == false {
		t.Errorf("User B should not be equivilant to user C")
	}
}

func TestReadLinesFromFile(t *testing.T) {
	var a Models.User
	var b Models.User

	a.ID = "0D00A443-93D3-8573-135D-8946397866A1"
	a.Firstname = "Hayes"
	a.Lastname = "Cruz"
	a.Email = "dui@aliquetsem.org"

	b.ID = "352AFADF-0A34-7933-7196-294A3AEEA6CE"
	b.Firstname = "Addison"
	b.Lastname = "Mcdonald"
	b.Email = "consectetuer.adipiscing.elit@Duissit.co.uk"
	bytesWritten := WriteUsersToFile([]Models.User{a, b}, "data.out")
	users := ReadUsersFromFile(bytesWritten, "data.out")

	if a.Compare(users[0]) == false {
		t.Errorf("User A should be equivilant to first user from file")
	}
	if b.Compare(users[1]) == false {
		t.Errorf("User A should be equivilant to second user from file")
	}
}

func TestGetUsersFromDB(t *testing.T) {
	connectionString := "temp.db"
	db := ConnectToDatabase(&connectionString)

	var users []Models.User

	db.Limit(2).Find(&users)

	var a Models.User
	var b Models.User

	a.ID = "0D00A443-93D3-8573-135D-8946397866A1"
	a.Firstname = "Hayes"
	a.Lastname = "Cruz"
	a.Email = "dui@aliquetsem.org"

	b.ID = "352AFADF-0A34-7933-7196-294A3AEEA6CE"
	b.Firstname = "Addison"
	b.Lastname = "Mcdonald"
	b.Email = "consectetuer.adipiscing.elit@Duissit.co.uk"

	if a.Compare(users[0]) == false {
		t.Errorf("User A should be equivilant to first user in database")
	}
	if b.Compare(users[1]) == false {
		t.Errorf("User A should be equivilant to second user in database")
	}

	CloseConnection(db)
}

func TestMatchDBWithFile(t *testing.T) {
	connectionString := "temp.db"
	db := ConnectToDatabase(&connectionString)

	var fromDB []Models.User

	db.Limit(2).Find(&fromDB)

	bytesWritten := WriteUsersToFile(fromDB, "data.out")

	fromFile := ReadUsersFromFile(bytesWritten, "data.out")

	if ValidateUsersInFile(fromDB, fromFile) == false {
		t.Errorf("Users from database should match users from file!")
	}
	CloseConnection(db)
}

func TestDeleteUsersFromDB(t *testing.T) {
	connectionString := "temp.db"
	db := ConnectToDatabase(&connectionString)

	var fromDB []Models.User
	query := "ID IN (\"0D00A443-93D3-8573-135D-8946397866A1\",\"352AFADF-0A34-7933-7196-294A3AEEA6CE\")"
	db.Where(query).Delete(Models.User{})
	db.Find(&fromDB, query)
	if len(fromDB) != 0 {
		t.Errorf("Users retrieved (%d) expected (%d)", len(fromDB), 0)
	}
	CloseConnection(db)
}
