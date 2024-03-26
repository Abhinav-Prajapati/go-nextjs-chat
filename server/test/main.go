package main

import "fmt"

type msg interface {
	printMsg()
}

type warning struct {
	msg string
}

type err struct {
	msg string
}

func (w *warning) printMsg() {
	fmt.Println("Warning ", w.msg)
}

func (e *err) printMsg() {
	fmt.Println("err ", e.msg)
}

func log(msg msg) {
	msg.printMsg()
}

func main() {
	// war := &warning{msg: "this is a warning"}
	err := &err{msg: "this is an error"}
	log(err)

}
