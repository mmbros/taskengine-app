package server

import (
	"embed"
	_ "embed"
	"encoding/json"
	"fmt"
	"io/fs"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
)

//go:embed content
var content embed.FS

func listFiles(folder string, recursive bool) ([]string, error) {
	var files []string

	var isSubDir bool

	re := regexp.MustCompile(`\.json$`)

	err := filepath.Walk(folder, func(path string, info os.FileInfo, err error) error {
		// log.Println(path, info.Name())

		if info.IsDir() {
			if !recursive && isSubDir {
				return filepath.SkipDir
			}
			isSubDir = true
		} else {
			if re.MatchString(path) == false {
				return nil
			}
			relpath, _ := filepath.Rel(folder, path)
			files = append(files, relpath)
		}
		return nil
	})
	return files, err
}

func handlerIndexJson(dataFolder string, recursive bool) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		files, err := listFiles(dataFolder, recursive)
		if err != nil {
			http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(files)
	}
}

/*
https://bhupesh.me/embedding-static-files-in-golang/


import _ "embed"

//go:embed sample.txt
var s string

*/
func HelloHandler(w http.ResponseWriter, r *http.Request) {

	// fmt.Fprintf(w, "Hello, there\n")

	buf, err := ioutil.ReadFile("asset/html/index.html")

	if err != nil {
		log.Fatal(err)
	}

	w.Header().Set("Content-Type", "text/html")
	w.Write(buf)
}

// func HelloHandlerEnbed(w http.ResponseWriter, r *http.Request) {
// 	w.Header().Set("Content-Type", "text/html")
// 	w.Write([]byte(html))
// }

func handlerCSS() http.Handler {

	fsys := fs.FS(content)
	css, err := fs.Sub(fsys, "content/css")
	if err != nil {
		panic(err)
	}

	return http.StripPrefix("/css/", http.FileServer(http.FS(css)))
}
func handlerContent() http.Handler {

	fsys := fs.FS(content)
	root, err := fs.Sub(fsys, "content")
	if err != nil {
		panic(err)
	}

	return http.FileServer(http.FS(root))
}
func Run(serverAddressPort, jsonDataFolder string, recursive bool) error {

	http.Handle("/", handlerContent())

	// http.Handle("/css/", handlerCSS())
	// http.Handle("/css/", http.StripPrefix("/css", http.FileServer(http.Dir("asset/js"))))

	// http.Handle("/js/", http.StripPrefix("/js", http.FileServer(http.Dir("asset/js"))))

	http.HandleFunc("/data", handlerIndexJson(jsonDataFolder, recursive))
	http.Handle("/data/", http.StripPrefix("/data", http.FileServer(http.Dir(jsonDataFolder))))

	fmt.Printf("server listening to %s\n", serverAddressPort)
	err := http.ListenAndServe(serverAddressPort, nil)
	return err
}
