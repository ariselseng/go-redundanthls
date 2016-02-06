package redundanthls

import (
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

func RedundantManifestFromUrl(url string, hosts []string) (string, error) {

	res, err := goreq.Request{Uri: url}.Do()
	if err != nil || res.StatusCode != 200 {
		return "", err
	}
	rawManifest, err := res.Body.ToString()
	if err != nil {
		return "", err
	}
	redundantManifest := getRedundantManifest(rawManifest, hosts)
	return redundantManifest, nil

}
