package redundanthls

import (
	"errors"
	"github.com/franela/goreq"
	"strings"
)

// GetRawManifest returns body of url as string
func GetRawManifest(url string) (string, error) {
	res, err := goreq.Request{Uri: url}.Do()
	if err != nil || res.StatusCode != 200 {
		return "", errors.New("Error retrieving raw manifest.")
	}
	body, err := res.Body.ToString()
	if err != nil {
		return "", err
	}
	return body, nil
}
func getRedundantManifest(rawManifest string, hosts []string, maxLevel string, protocol string) string {
	var redundantManifest string
	var hasFoundMaxLevel bool
	var looksForMaxLevel bool
	if maxLevel != "" {
		looksForMaxLevel = true
	}
	rawLines := strings.Split(rawManifest, "\n")
	// fmt.Println(rawLines[0])
	var levelInfo string
	for key, line := range rawLines {
		if key == 0 {
			redundantManifest = line
		} else {

			if strings.HasPrefix(line, "#EXT-X-STREAM-INF:PROGRAM-ID=1") {
				levelInfo = line
				// redundantManifest = redundantManifest + "\n"
			} else {

				if strings.HasSuffix(line, ".m3u8") {

					if looksForMaxLevel && !hasFoundMaxLevel || !looksForMaxLevel {
						if strings.Contains(line, maxLevel) {
							hasFoundMaxLevel = true
						}
						for _, host := range hosts {
							redundantManifest = redundantManifest + "\n" + levelInfo + "\n" + protocol + host + "/" + line
						}
					}

				} else {
					redundantManifest = redundantManifest + "\n" + line
				}
			}

		}
	}

	return redundantManifest
}

type error interface {
	Error() string
}

// RedundantManifestFromString returns redundant manifest from string
func RedundantManifestFromString(rawManifest string, hosts []string, maxLevel string, protocol string) (string, error) {

	if rawManifest == "" {
		return "", errors.New("Manifest is empty")
	}
	redundantManifest := getRedundantManifest(rawManifest, hosts, maxLevel, protocol)
	return redundantManifest, nil

}

// RedundantManifestFromURL returns redundant manifest from url
func RedundantManifestFromURL(url string, hosts []string, protocol string) (string, error) {

	res, err := goreq.Request{Uri: url}.Do()
	if err != nil || res.StatusCode != 200 {
		return "", errors.New("Error retrieving manifest.")
	}
	rawManifest, err := res.Body.ToString()
	if err != nil {
		return "", errors.New("Error getting the body of the manifest.")
	}
	redundantManifest := getRedundantManifest(rawManifest, hosts, "", protocol)
	return redundantManifest, nil

}
