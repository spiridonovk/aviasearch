package main

import (
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"os"
)


type Flights struct {
	XMLName           xml.Name `xml:"AirFareSearchResponse"`
	Text              string   `xml:",chardata"`
	RequestTime       string   `xml:"RequestTime,attr"`
	ResponseTime      string   `xml:"ResponseTime,attr"`
	RequestId         string   `xml:"RequestId"`
	PricedItineraries struct {
		Text    string `xml:",chardata"`
		Flights []struct {
			Text                  string `xml:",chardata"`
			OnwardPricedItinerary struct {
				Text    string `xml:",chardata"`
				Flights struct {
					Text   string `xml:",chardata"`
					Flight []struct {
						Text    string `xml:",chardata"`
						Carrier struct {
							Text string `xml:",chardata"`
							ID   string `xml:"id,attr"`
						} `xml:"Carrier"`
						FlightNumber       string `xml:"FlightNumber"`
						Source             string `xml:"Source"`
						Destination        string `xml:"Destination"`
						DepartureTimeStamp string `xml:"DepartureTimeStamp"`
						ArrivalTimeStamp   string `xml:"ArrivalTimeStamp"`
						Class              string `xml:"Class"`
						NumberOfStops      string `xml:"NumberOfStops"`
						FareBasis          string `xml:"FareBasis"`
						WarningText        string `xml:"WarningText"`
						TicketType         string `xml:"TicketType"`
					} `xml:"Flight"`
				} `xml:"Flights"`
			} `xml:"OnwardPricedItinerary"`
			ReturnPricedItinerary struct {
				Text    string `xml:",chardata"`
				Flights struct {
					Text   string `xml:",chardata"`
					Flight []struct {
						Text    string `xml:",chardata"`
						Carrier struct {
							Text string `xml:",chardata"`
							ID   string `xml:"id,attr"`
						} `xml:"Carrier"`
						FlightNumber       string `xml:"FlightNumber"`
						Source             string `xml:"Source"`
						Destination        string `xml:"Destination"`
						DepartureTimeStamp string `xml:"DepartureTimeStamp"`
						ArrivalTimeStamp   string `xml:"ArrivalTimeStamp"`
						Class              string `xml:"Class"`
						NumberOfStops      string `xml:"NumberOfStops"`
						FareBasis          string `xml:"FareBasis"`
						WarningText        string `xml:"WarningText"`
						TicketType         string `xml:"TicketType"`
					} `xml:"Flight"`
				} `xml:"Flights"`
			} `xml:"ReturnPricedItinerary"`
			Pricing struct {
				Text           string `xml:",chardata"`
				Currency       string `xml:"currency,attr"`
				ServiceCharges []struct {
					Text       string `xml:",chardata"`
					Type       string `xml:"type,attr"`
					ChargeType string `xml:"ChargeType,attr"`
				} `xml:"ServiceCharges"`
			} `xml:"Pricing"`
		} `xml:"Flights"`
	} `xml:"PricedItineraries"`
}

func main() {

	// Open our xmlFile
	xmlFile, err := os.Open("RS_Via-3.xml")
	// if we os.Open returns an error then handle it
	if err != nil {
		fmt.Println(err)
	}

	fmt.Println("Successfully Opened flighs.xml")
	// defer the closing of our xmlFile so that we can parse it later on
	defer xmlFile.Close()

	// read our opened xmlFile as a byte array.
	byteValue, _ := ioutil.ReadAll(xmlFile)

	// we initialize our Users array
	var flights Flights
	// we unmarshal our byteArray which contains our
	// xmlFiles content into 'flighs' which we defined above
	xml.Unmarshal(byteValue, &flights)

	// we iterate through every user within our flighs array and
	// print out the user Type, their name, and their facebook url
	// as just an example
	for _,allFlight := range flights.PricedItineraries.Flights {
		for _, foo := range allFlight.OnwardPricedItinerary.Flights.Flight {
			fmt.Println("SOurce: " + foo.Source)
			fmt.Println("Destination: " + foo.Destination)

		}
	}
}

