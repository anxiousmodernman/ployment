package webhook

import (
	"archive/zip"
	"bytes"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/anxiousmodernman/ployment/config"
)

// get rid of these
//const repoUrl string = "https://github.com/anxiousmodernman/coleman-dot-codes/zipball/master"
//const tempDir string = "/Users/coleman/temp"

type AppContext struct {
	Config config.PloymentConfig
}

type Hook struct {
	*AppContext
	Handler func(http.ResponseWriter, *http.Request, *AppContext) error
}

//type Hook func(res http.ResponseWriter, req *http.Request) error

// make Hook satisfy http.Handler
func (fn Hook) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if err := fn.Handler(w, r, fn.AppContext); err != nil {
		log.Fatal("Internal server error", err)
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
	}
}

func WebhookHandler(w http.ResponseWriter, r *http.Request, ctx *AppContext) error {
	urlReader, err := getReaderFromUrl(ctx.Config.RepositoryUrl)
	if err != nil {
		return fmt.Errorf("Could not get GitHub url: %s", err)
	}

	zr, err := zip.NewReader(urlReader, int64(urlReader.Len()))
	if err != nil {
		return fmt.Errorf("Unable to read zip: %s", err)
	}

	for _, zf := range zr.File {
		if err := writeFile(zf, ctx.Config.TargetDirectory); err != nil {
			return fmt.Errorf("Unable to write file %s. Error: %s", zf.Name, err)
		}

	}
	return nil
}

func getReaderFromUrl(url string) (*bytes.Reader, error) {
	res, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	buf := &bytes.Buffer{}

	_, err = io.Copy(buf, res.Body)
	if err != nil {
		return nil, err
	}

	return bytes.NewReader(buf.Bytes()), nil
}

func writeFile(zf *zip.File, currentDir string) error {
	fr, err := zf.Open()
	if err != nil {
		log.Printf("Error opening file for write: %s", zf.Name)
		return err
	}
	defer fr.Close()

	// File is a struct that containts FileHeader
	fileHeader := zf.FileHeader
	fileInfo := fileHeader.FileInfo()
	isDir := fileInfo.IsDir()

	// zip only uses forward slash in golang. Replace with correct os separater, if necessary.
	path := strings.Replace(filepath.Join(currentDir, zf.Name), `/`, string(filepath.Separator), -1)
	dir, _ := filepath.Split(path)

	if isDir {

		err = os.MkdirAll(path, 0777)
		if err != nil {
			log.Printf("Unable to create directory: %s", dir)
			return err
		}
		return nil
	}

	f, err := os.Create(path)
	if err != nil {
		log.Printf("Unable to create file: %s", path)
		return err
	}
	defer f.Close()

	_, err = io.Copy(f, fr)
	if err != nil {
		log.Printf("Issue writing file <%s>: %s", path, err)
		return err
	}
	return nil
}
