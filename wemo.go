package wemo

import (
	"encoding/xml"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/huin/goupnp"
)

type Switch struct {
	Host string
	Name string
}

type Setup struct {
	XMLName      xml.Name `xml:"root"`
	FriendlyName string   `xml:"device>friendlyName"`
}

type BinaryState struct {
	XMLName     xml.Name `xml:"Envelope"`
	BinaryState int      `xml:"Body>GetBinaryStateResponse>BinaryState"`
}

func Get(url string) []byte {
	resp, err := http.Get(url)
	if err != nil {
		// handle error
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	return body
}

func Switches() ([]Switch, error) {
	var switches []Switch
	devices, err := goupnp.DiscoverDevices("urn:Belkin:device:controllee:1")
	if err != nil {
		return switches, err
	}

	for i := range devices {
		var wemoSwitch Switch
		var setup Setup
		setupXml := Get(devices[i].Root.URLBaseStr)
		xml.Unmarshal(setupXml, &setup)
		wemoSwitch.Host = devices[i].Root.URLBase.Host
		wemoSwitch.Name = setup.FriendlyName
		switches = append(switches, wemoSwitch)
	}
	return switches, nil
}

func (s *Switch) On() {
	s.setBinaryState("1")
}

func (s *Switch) Off() {
	s.setBinaryState("0")
}

func (s *Switch) Status() int {
	var binaryState BinaryState
	reqXml := `<?xml version="1.0" encoding="utf-8"?><s:Envelope xmlns:s="http://schemas.xmlsoap.org/soap/envelope/" s:encodingStyle="http://schemas.xmlsoap.org/soap/encoding/"><s:Body><u:GetBinaryState xmlns:u="urn:Belkin:service:basicevent:1"></u:GetBinaryState></s:Body></s:Envelope>`
	url := "http://" + s.Host + "/upnp/control/basicevent1"
	req, _ := http.NewRequest("POST", url, strings.NewReader(reqXml))
	req.Header.Add("SOAPACTION", `"urn:Belkin:service:basicevent:1#GetBinaryState"`)
	req.Header.Add("Content-type", `text/xml; charset="utf-8"`)
	client := http.Client{}
	resp, _ := client.Do(req)
	defer resp.Body.Close()
	body, _ := ioutil.ReadAll(resp.Body)
	xml.Unmarshal(body, &binaryState)
	return binaryState.BinaryState
}
func (s *Switch) setBinaryState(signal string) {
	binaryState := `<?xml version="1.0" encoding="utf-8"?><s:Envelope xmlns:s="http://schemas.xmlsoap.org/soap/envelope/" s:encodingStyle="http://schemas.xmlsoap.org/soap/encoding/"><s:Body><u:SetBinaryState xmlns:u="urn:Belkin:service:basicevent:1"><BinaryState>` + signal + `</BinaryState></u:SetBinaryState></s:Body></s:Envelope>`
	url := "http://" + s.Host + "/upnp/control/basicevent1"
	req, _ := http.NewRequest("POST", url, strings.NewReader(binaryState))
	req.Header.Add("SOAPACTION", `"urn:Belkin:service:basicevent:1#SetBinaryState"`)
	req.Header.Add("Content-type", `text/xml; charset="utf-8"`)
	client := http.Client{}
	client.Do(req)
}
