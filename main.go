package main

import (
	"log"
	"manga-cli/cmd"
)

func main(){
	if err := cmd.Execute(); err != nil{
		log.Fatalf("Error %v", err)
	}
}