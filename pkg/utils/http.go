package utils

import (
	"github.com/mattfenwick/collections/pkg/file"
	"github.com/pkg/errors"
	"io/ioutil"
	"net/http"
)

func GetFileFromURL(url string, path string) error {
	bytes, err := GetURL(url)
	if err != nil {
		return err
	}
	return file.Write(path, bytes, 0777)
}

func GetURL(url string) ([]byte, error) {
	response, err := http.Get(url)
	if err != nil {
		return nil, errors.Wrapf(err, "unable to GET %s", url)
	}
	defer response.Body.Close()
	if response.StatusCode < 200 || response.StatusCode > 299 {
		return nil, errors.Errorf("GET request to %s failed with status code %d", url, response.StatusCode)
	}
	bytes, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return nil, errors.Wrapf(err, "unable to read body from GET to %s", url)
	}

	return bytes, nil
}
