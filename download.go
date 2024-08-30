package main

import (
	"archive/tar"
	"compress/gzip"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"

	"github.com/labstack/gommon/log"
)

// Get GeoLite2 database files
func GetGeoLite2(databaseName string, licenseKey string) (string, error) {

	// Database names can be:
	// - GeoLite2-ASN
	// - GeoLite2-City

	// Downloaded files can be:
	// - GeoLite2-ASN.tar.gz
	// - GeoLite2-City.tar.gz

	// Database files are:
	// - ./folder/GeoLite2-ASN.mmdb
	// - ./folder/GeoLite2-City.mmdb

	// First, check if directories and files are existing.
	pattern := fmt.Sprintf("./%s_*", databaseName)
	log.Infof("Check patern: %s", pattern)
	matches, err := filepath.Glob(pattern)
	if err != nil {
		return "", err
	}
	if len(matches) > 0 {
		last := matches[len(matches)-1]
		log.Infof("Last directory found: %s", last)
		// Check if the database file exists and if it exits, return the filename and path.
		databaseFile := fmt.Sprintf("./%s/%s.mmdb", last, databaseName)
		if _, err := os.Stat(databaseFile); err == nil {
			return databaseFile, nil
		}
	}

	// If exist, remove the old database archive
	if _, err := os.Stat("./" + databaseName + ".tar.gz"); err == nil {
		e := os.Remove("./" + databaseName + ".tar.gz")
		if e != nil {
			return "", err
		}
	}

	// Get the databases from maxmind
	getUL := fmt.Sprintf("https://download.maxmind.com/app/geoip_download?edition_id=%s&license_key=%s&suffix=tar.gz", databaseName, licenseKey)
	resp, err := http.Get(getUL)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	out, err := os.Create(databaseName + ".tar.gz")
	if err != nil {
		return "", err
	}
	defer out.Close()
	_, err = io.Copy(out, resp.Body)
	if err != nil {
		return "", err
	}

	// Extract the database archive file
	dbArchive, err := os.Open(databaseName + ".tar.gz")
	if err != nil {
		return "", err
	}
	err = ExtractTarGz(dbArchive)
	if err != nil {
		return "", err
	}

	// Delete database archive
	e := os.Remove("./" + databaseName + ".tar.gz")
	if e != nil {
		return "", err
	}

	// Check if directories and files are existing.
	pattern = fmt.Sprintf("./%s_*", databaseName)
	matches, err = filepath.Glob(pattern)
	if err != nil {
		return "", err
	}
	if len(matches) > 0 {
		last := matches[len(matches)-1]
		log.Infof("Last directory found: %s", last)
		// Check if the database file exists and if it exits, return the filename and path.
		databaseFile := fmt.Sprintf("./%s/%s.mmdb", last, databaseName)
		if _, err := os.Stat(databaseFile); err == nil {
			return databaseFile, nil
		}
	}

	return "", nil
}

func ExtractTarGz(gzipStream io.Reader) error {
	uncompressedStream, err := gzip.NewReader(gzipStream)
	if err != nil {
		return err
	}

	tarReader := tar.NewReader(uncompressedStream)

	for {
		header, err := tarReader.Next()

		if err == io.EOF {
			break
		}

		if err != nil {
			return err
		}

		switch header.Typeflag {
		case tar.TypeDir:
			if err := os.Mkdir(header.Name, 0755); err != nil {
				return err
			}
		case tar.TypeReg:
			outFile, err := os.Create(header.Name)
			if err != nil {
				return err
			}
			if _, err := io.Copy(outFile, tarReader); err != nil {
				return err
			}
			outFile.Close()

		default:
			return err
		}

	}
	return nil
}
