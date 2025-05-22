/*
Copyright © 2025 Kodo Robotics

*/
package virtualbox

import (
	"encoding/xml"
	"fmt"
	"os"
)

func ExtractVMNameFromOVF(ovfPath string) (string, error) {
	xmlFile, err := os.Open(ovfPath)
	if err != nil {
		return "", err
	}
	defer xmlFile.Close()

	decoder := xml.NewDecoder(xmlFile)

	for {
		tok, err := decoder.Token()
		if err != nil {
			break
		}

		switch se := tok.(type) {
		case xml.StartElement:
			if se.Name.Local == "Machine" && se.Name.Space == "http://www.virtualbox.org/ovf/machine" {
				for _, attr := range se.Attr {
					if attr.Name.Local == "name" {
						return attr.Value, nil
					}
				}
			}
		}
	}

	return "", fmt.Errorf("could not find <vbox:Machine name=...> in OVF")
}