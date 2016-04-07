package main

import (
	"github.com/JonathanFraser/csvproc"
	"os"
	"math/rand"
	"time"
)


func main() {
	rand.Seed(time.Now().UnixNano())
	f := csvproc.Generate(120000,10)
	f.Store(os.Stdout)
}
