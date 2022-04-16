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
)

//go:embed content
var content embed.FS

func listFiles(folder string) ([]string, error) {
	var files []string

	err := filepath.Walk(folder, func(path string, info os.FileInfo, err error) error {
		// log.Println(path, info.Name())

		if !info.IsDir() {
			relpath, _ := filepath.Rel(folder, path)
			files = append(files, relpath)
		}
		return nil
	})
	return files, err
}

func handlerIndexJson(dataFolder string) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		files, err := listFiles(dataFolder)
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
func Run(serverAddressPort, jsonDataFolder string) error {

	http.Handle("/", handlerContent())

	// http.Handle("/css/", handlerCSS())
	// http.Handle("/css/", http.StripPrefix("/css", http.FileServer(http.Dir("asset/js"))))

	// http.Handle("/js/", http.StripPrefix("/js", http.FileServer(http.Dir("asset/js"))))

	http.HandleFunc("/data", handlerIndexJson(jsonDataFolder))
	http.Handle("/data/", http.StripPrefix("/data", http.FileServer(http.Dir(jsonDataFolder))))

	fmt.Printf("server listening to %s\n", serverAddressPort)
	err := http.ListenAndServe(serverAddressPort, nil)
	return err
}
