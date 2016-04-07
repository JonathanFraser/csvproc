package main

import (
	"os"
	"log"
	"path/filepath"
	"errors"
	"github.com/JonathanFraser/csvproc"
)

func main() {
	if len(os.Args) != 2 {
		log.Fatal("usage is: "+os.Args[0]+" outputdir")
		return 
	}

	fullpath, err := filepath.Abs(os.Args[1])
	if err != nil {
		log.Fatal(err)
	}

	info,err := os.Stat(fullpath)
	if err != nil {
		log.Fatal(err)
	}

	if !info.IsDir() {
		log.Fatal(errors.New("Second parameter is not a directory"))
	}

	csvFile,err := csvproc.Load(os.Stdin)
	if err != nil {
		log.Fatal(err)
	}

	waves := csvFile.ExtractWaves()
	for _,v := range waves {
		f,err := os.Create(filepath.Clean(fullpath+string(os.PathSeparator)+v.Name))
		if err != nil {
			log.Fatal(err)
		}
		defer f.Close()
		err = v.Store(f)
		if err != nil {
			log.Fatal(err)
		}
	}
}
