package engine

import (
	"encoding/xml"
	"io/ioutil"
	"os"
	"sort"
	"sync"
	"time"
)

const (
	optimal       = "optimal"
	duration      = "duration"
	totalAmount   = "TotalAmount"
	passengerType = "SingleAdult"
)

type Data struct {
	XMLName           xml.Name `xml:"AirFareSearchResponse"`
	RequestTime       string   `xml:"RequestTime,attr"`
	ResponseTime      string   `xml:"ResponseTime,attr"`
	RequestId         string   `xml:"RequestId"`
	PricedItineraries struct {
		Flights []*Flights `json:"flights" xml:"Flights"`
	} `json:"priced_itineraries" xml:"PricedItineraries"`
}

type Flights struct {
	//debug information
	TotalPrice            float32 `json:"totalPrice"`
	RouteDuration         float32 `json:"routeDuration"`
	Stops                 int     `json:"stops"`
	OnwardPricedItinerary struct {
		Flights struct {
			Flight []struct {
				Carrier struct {
					Text string `json:"text" xml:",chardata"`
					ID   string `json:"id" xml:"id,attr"`
				} `json:"carrier" xml:"Carrier"`
				FlightNumber       string `json:"flightNumber "xml:"FlightNumber"`
				Source             string `json:"source "xml:"Source"`
				Destination        string `json:"destination "xml:"Destination"`
				DepartureTimeStamp string `json:"departureTimeStamp" xml:"DepartureTimeStamp"`
				ArrivalTimeStamp   string `json:"arrivalTimeStamp"xml:"ArrivalTimeStamp"`
				Class              string `json:"class" xml:"Class"`
				NumberOfStops      string `json:"numberOfStops" xml:"NumberOfStops"`
				FareBasis          string `json:"fareBasis" xml:"FareBasis"`
				WarningText        string `json:"warningText"  xml:"WarningText"`
				TicketType         string `json:"ticketType" xml:"TicketType"`
			} `json:"flight" xml:"Flight"`
		} `json:"flights" xml:"Flights"`
	} `json:"onwardPricedItinerary" xml:"OnwardPricedItinerary"`
	ReturnPricedItinerary *ReturnPricedItinerary `json:"returnPricedItinerary,omitempty" xml:"ReturnPricedItinerary"`
	Pricing               struct {
		Currency       string `json:"currency" xml:"currency,attr"`
		ServiceCharges []struct {
			Text       float32 `json:"text" xml:",chardata" json:"value"`
			Type       string  `json:"type" xml:"type,attr"`
			ChargeType string  `json:"chargeType" xml:"ChargeType,attr"`
		} `json:"serviceCharges" xml:"ServiceCharges"`
	} `json:"pricing" xml:"Pricing"`
}

type ReturnPricedItinerary struct {
	Flights struct {
		Flight []struct {
			Text    string `xml:",chardata" json:"-"`
			Carrier struct {
				ID string `xml:"id,attr"`
			} `xml:"Carrier"`
			FlightNumber       string `json:"flightNumber "xml:"FlightNumber"`
			Source             string `json:"source "xml:"Source"`
			Destination        string `json:"destination "xml:"Destination"`
			DepartureTimeStamp string `json:"departureTimeStamp" xml:"DepartureTimeStamp"`
			ArrivalTimeStamp   string `json:"arrivalTimeStamp"xml:"ArrivalTimeStamp"`
			Class              string `json:"class" xml:"Class"`
			NumberOfStops      string `json:"numberOfStops" xml:"NumberOfStops"`
			FareBasis          string `json:"fareBasis" xml:"FareBasis"`
			WarningText        string `json:"warningText"  xml:"WarningText"`
			TicketType         string `json:"ticketType" xml:"TicketType"`
		} `json:"flight" xml:"Flight"`
	} `json:"flights" xml:"Flights"`
}

//ParseData read xml to struct
func ParseData(fileName string) ([]*Flights, error) {
	// Open our xmlFile
	xmlFile, err := os.Open(fileName)
	if err != nil {
		return nil, err
	}
	defer xmlFile.Close()

	byteValue, _ := ioutil.ReadAll(xmlFile)
	var data Data
	err = xml.Unmarshal(byteValue, &data)
	if err != nil {
		return nil, err
	}
	return data.PricedItineraries.Flights, nil
}

//GetVariants get all variants from xml
func GetVariants(sort, order string) ([]*Flights, error) {
	dataVia3, err := ParseData("RS_Via-3.xml")
	if err != nil {
		return nil, err
	}
	dataViaOW, err := ParseData("RS_ViaOW.xml")
	if err != nil {
		return nil, err
	}
	fullData := append(dataVia3, dataViaOW...)
	wg := &sync.WaitGroup{}
	wg.Add(len(fullData))
	for _, perVariant := range fullData {
		go func(variant *Flights) {
			defer wg.Done()
			price, duration := getPriceAndTime(variant)
			variant.Stops = len(variant.OnwardPricedItinerary.Flights.Flight) - 1
			variant.TotalPrice = price
			variant.RouteDuration = duration
		}(perVariant)
	}
	wg.Wait()
	sortBy(fullData, sort, order)
	if sort == optimal {
		bestPriceTickets := fullData[:len(fullData)/2]
		sortBy(bestPriceTickets, duration, "asc")
		return bestPriceTickets, nil
	}
	return fullData, nil
}

//sorting by filter
func sortBy(flights []*Flights, sortBy, orderBy string) {
	sort.Slice(flights[:], func(i, j int) bool {
		firstValue := flights[i].TotalPrice
		secondValue := flights[j].TotalPrice

		if sortBy == duration {
			firstValue = flights[i].RouteDuration
			secondValue = flights[j].RouteDuration
		}
		if orderBy == "desc" {
			return firstValue > secondValue
		}
		return firstValue < secondValue

	})
}

func getPriceAndTime(flight *Flights) (price, routeDuration float32) {

	for _, charges := range flight.Pricing.ServiceCharges {
		if charges.ChargeType == totalAmount && charges.Type == passengerType {
			price = charges.Text
		}
	}
	layout := "2006-01-02T1504"
	segments := flight.OnwardPricedItinerary.Flights.Flight
	departureTime, _ := time.Parse(layout, segments[0].DepartureTimeStamp)
	arrivalTime, _ := time.Parse(layout, segments[len(segments)-1].ArrivalTimeStamp)
	routeDuration = float32(arrivalTime.Sub(departureTime).Minutes())
	return price, routeDuration
}
