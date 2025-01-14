package utils

import (
	"fmt"
	"os"
)

type Opswfile struct {
	file     *os.File
	filename string
}

func NewOpswfile(filename string) Opswfile {
	return Opswfile{
		filename: filename,
	}
}

func (u Opswfile) Openfile() error {
	file, err := os.Create(u.filename)
	if err != nil {
		return err
	}
	u.file = file
	return nil
}

func (u Opswfile) Closefile() {
	if u.file != nil {
		u.file.Close()
	}

}

func (u Opswfile) Write(content string) error {
	if u.file == nil {
		if err := u.Openfile(); err != nil {
			return err
		}
	}
	if resBytes, err := u.file.WriteString(content); err != nil {
		return err
	} else {
		fmt.Printf("Data writed in file %d.", resBytes)
	}
	return nil

}
