package main

import "fmt"

func main() {
  word := "Hello!"
  ptr := &word
  fmt.Println(ptr) //O endereço de memória para o qual está apontando
  fmt.Println(*ptr) //O valor armazenado no endereço de memória

  word = "World!"
  fmt.Println(ptr) //O endereço de memória para o qual está apontando
  fmt.Println(*ptr) //O valor armazenado no endereço de memória
}