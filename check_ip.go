package main

import (
	"bufio"
	"encoding/csv"
	"encoding/json"
	"log"
	"net"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

// Handler
func (conf Config) getIPInfo(c echo.Context) error {
	data := new(Data)
	data.Title = "IP information"

	return c.Render(http.StatusOK, "form.html", map[string]interface{}{
		"Data": data,
	})
}

type Data struct {
	Title        string
	DownloadJSON string
	DownloadCSV  string
	Results      []*ipinformation
}

type ipinformation struct {
	IP          string
	ASN         uint
	ASNorg      string
	City        string
	SubDivision string
	Country     string
	CoutryCode  string
}

func (conf Config) postIPInfo(c echo.Context) error {
	data := new(Data)
	data.Title = "IP information"

	iplist := c.FormValue("iplist")

	// Generate something for creating a file
	id := uuid.New()

	// Create CSV file stuff
	csvfilename := "public/csv/" + id.String() + ".csv"
	csvfile, err := os.Create(csvfilename)
	if err != nil {
		log.Fatal(err)
	}
	data.DownloadCSV = "csv/" + id.String() + ".csv"

	// Create JSON file stuff
	jsonfile := "public/json/" + id.String() + ".json"
	data.DownloadJSON = "json/" + id.String() + ".json"

	csvw := csv.NewWriter(csvfile)
	csvl := []string{"ip", "asn", "asnorg", "city", "subdivision", "country", "coutrycode"}
	csvw.Write(csvl)

	scanner := bufio.NewScanner(strings.NewReader(iplist))
	for scanner.Scan() {
		ips := new(ipinformation)

		// Get this line's IP and remove any spaces
		lip := strings.ReplaceAll(scanner.Text(), " ", "")

		// Get this line's IP and remove any tabs
		lip = strings.ReplaceAll(lip, "\t", "")

		// Get this line's IP and remove IPv6 starting bracket
		lip = strings.ReplaceAll(lip, "[", "")

		// Get this line's IP and remove IPv6 ending bracket
		lip = strings.ReplaceAll(lip, "]", "")

		// Check if the IP address will parse and is valid
		pip := net.ParseIP(lip)
		if pip == nil {
			// Of not. log not an IP and go further
			ips.IP = lip
			ips.ASN = 0
			ips.ASNorg = "invalid ip"
			ips.City = "invalid ip"
			ips.SubDivision = "invalid ip"
			ips.Country = "invalid ip"
			ips.CoutryCode = "invalid ip"
			data.Results = append(data.Results, ips)
			csvl := []string{ips.IP, "0", ips.ASNorg, ips.City, ips.SubDivision, ips.Country, ips.CoutryCode}
			csvw.Write(csvl)
			continue
		}
		ips.IP = pip.String()

		masn, err := getASN(conf.Maxmind.ASN, pip)
		if err != nil {
			log.Fatal(err)
		}
		ips.ASN = masn.AutonomousSystemNumber
		ips.ASNorg = masn.AutonomousSystemOrganization

		mcity, err := getCity(conf.Maxmind.City, pip)
		if err != nil {
			log.Fatal(err)
		}
		if len(mcity.Subdivisions) > 0 {
			ips.SubDivision = mcity.Subdivisions[0].Names["en"]
		}
		ips.City = mcity.City.Names["en"]
		ips.Country = mcity.Country.Names["en"]
		ips.CoutryCode = mcity.Country.IsoCode

		data.Results = append(data.Results, ips)

		var asnStr string = strconv.FormatUint(uint64(ips.ASN), 10)
		csvl := []string{ips.IP, asnStr, ips.ASNorg, ips.City, ips.SubDivision, ips.Country, ips.CoutryCode}
		csvw.Write(csvl)

	}
	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}

	csvw.Flush()

	resultsJSON, _ := json.Marshal(data.Results)
	err = os.WriteFile(jsonfile, resultsJSON, 0644)
	if err != nil {
		log.Fatal(err)
	}

	return c.Render(http.StatusOK, "form.html", map[string]interface{}{
		"Data": data,
	})
}
