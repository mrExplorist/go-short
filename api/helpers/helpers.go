package helpers

import (
	"os"
	"strings"
)

func EnforceHTTP(url string) string { 
			if url[:4] != "http" {  // if the url does not start with http, add it
				url = "http://" + url
			}
			return url

}

func IsDomainValid(url string) bool {


	// if the url is the same as the domain, return false 
	if url == os.Getenv("DOMAIN"){
		return false
	}


	newURL := strings.Replace(url, "http://", "", 1) // remove http:// from the url
	newURL = strings.Replace(newURL, "https://", "", 1) // remove https:// from the url 
	newURL = strings.Replace(newURL, "www.", "", 1) // remove www. from the url
	newURL = strings.Split(newURL, "/")[0] // remove everything after the first / in the url

	 if newURL == os.Getenv("DOMAIN") { // if the url is the same as the domain, return false
		return false
	}

	return true
}