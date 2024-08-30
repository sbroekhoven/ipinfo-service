package main

import (
	"log"
	"net"

	"github.com/oschwald/geoip2-golang"
)

func getASN(database string, ip net.IP) (*geoip2.ASN, error) {
	// Get GeoIP ASN information
	asndb, err := geoip2.Open(database)
	if err != nil {
		log.Fatal(err)
	}
	defer asndb.Close()
	asnr, err := asndb.ASN(ip)
	if err != nil {
		log.Fatal(err)
	}
	return asnr, err
}

func getCity(database string, ip net.IP) (*geoip2.City, error) {
	// Get GeoIP City information
	citydb, err := geoip2.Open(database)
	if err != nil {
		log.Fatal(err)
	}
	defer citydb.Close()
	cr, err := citydb.City(ip)
	if err != nil {
		log.Fatal(err)
	}
	return cr, err
}
