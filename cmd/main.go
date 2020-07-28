package main

/*
import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/Alvarios/guards/guards"
	"net/http"
	"os"
	"time"
)*/

import (
	"github.com/Alvarios/guards/server"
	"net/http"
	"time"
)

type Log struct {
	Id           int       `json:"id"`
	ErrorMessage string    `json:"error_message"`
	Time         time.Time `json:"-"`
	Message      string    `json:"message"`
	Error        string    `json:"error"`
}

func (l Log) Compare(log *Log) bool {
	return (l.Id == log.Id && l.ErrorMessage == log.ErrorMessage &&
		l.Message == log.Message && l.Error == log.Error)
}

var expected Log = Log{
	Id:           http.StatusBadRequest,
	ErrorMessage: http.StatusText(http.StatusBadRequest),
	Time:         time.Now(),
	Message:      "test_init",
	Error:        "Error test",
}

func main() {

	s := server.InitializeEvent()
	s.Run()
	/*fileName := "test_init_file.json"
	file, err := os.Create(fileName)
	if err != nil {
		//		t.Errorf("Failed to create file : %s", err.Error())
		return
	}
	g := guards.NewLogger(file, true)
	g.InvalidRequest(errors.New("Error test"), "test_init")
	// read from the file
	log := &Log{}
	file.Close()
	file, _ = os.Open(fileName)
	json.NewDecoder(file).Decode(log)

	if expected.Compare(log) == false {
		// error
		fmt.Printf("Not equal log. Should be equal. have : %v | wanted %v", log, expected)
	}

	fmt.Print("Everything went good")*/
}
