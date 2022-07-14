package conn

import (
	"fmt"
)

type DBConnInfo struct {
	DBHost   string
	DBPort   int32
	DBName   string
	Username string
	Password string
}

func (x *DBConnInfo) PasswordAsEnv() *string {
	if x.Password != "" {
		s := fmt.Sprintf(`PGPASSWORD=%v`, x.Password)
		return &s
	}

	return nil
}

func (x *DBConnInfo) DBHostAsCmdArg() *string {

	if x.DBHost != "" {
		s := fmt.Sprintf("--host=%v", x.DBHost)
		return &s
	}

	return nil
}

func (x *DBConnInfo) DBPortAsCmdArg() *string {

	if x.DBPort != 0 {
		s := fmt.Sprintf(`--port=%v`, x.DBPort)
		return &s
	}

	return nil
}

func (x *DBConnInfo) UsernameAsCmdArg() *string {

	if x.Username != "" {
		s := fmt.Sprintf(`--username=%v`, x.Username)
		return &s
	}

	return nil
}

func (x *DBConnInfo) DBNameAsCmdArg(argKey *string) *string {

	if x.DBName != "" {
		var s string
		if argKey != nil {
			s = fmt.Sprintf(`%v=%v`, *argKey, x.DBName)
		} else {
			s = x.DBName
		}

		return &s
	}

	return nil
}
