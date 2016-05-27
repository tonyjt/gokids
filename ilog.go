package gokids

import(
  "fmt"
)

var log ILog

type ILog interface{
  Error(format string, a ...interface{})
  Info(format string, a ...interface{})
}


type LogDefault struct{

}

func NewLogDefault()*LogDefault{
  return &LogDefault{}
}

func (l *LogDefault) Error(format string, a ...interface{}){
  fmt.Print("error:")
  fmt.Printf(format, a...)
  fmt.Print("\n")
}

func (l *LogDefault) Info(format string, a ...interface{}){
  fmt.Print("info:")
  fmt.Printf(format, a...)
  fmt.Print("\n")
}
