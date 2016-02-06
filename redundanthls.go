package redundanthls

import (
	"errors"
	"github.com/franela/goreq"
	"strings"
)

func getRawManifest(path string) (string, error) {
	res, err := goreq.Request{Uri: path}.Do()
	if err != nil || res.StatusCode != 200 {
		return "", err
	}
	body, err := res.Body.ToString()
	if err != nil {
		return "", err
	}
	return body, nil
}
func getRedundantManifest(rawManifest string, hosts []string) string {
	var redundantManifest string
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
					for _, host := range hosts {
						redundantManifest = redundantManifest + "\n" + levelInfo + "\nhttp://" + host + "/" + line
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

func RedundantManifestFromString(rawManifest string, hosts []string) (string, error) {

	if rawManifest == "" {
		return "", errors.New("Manifest is empty")
	}
	redundantManifest := getRedundantManifest(rawManifest, hosts)
	return redundantManifest, nil

}
func RedundantManifestFromUrl(url string, hosts []string) (string, error) {

	res, err := goreq.Request{Uri: url}.Do()
	if err != nil || res.StatusCode != 200 {
		return "", errors.New("Error retrieving manifest.")
	}
	rawManifest, err := res.Body.ToString()
	if err != nil {
		return "", errors.New("Error getting the body of the manifest.")
	}
	redundantManifest := getRedundantManifest(rawManifest, hosts)
	return redundantManifest, nil

}
