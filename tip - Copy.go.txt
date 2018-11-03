package main

import (
	"encoding/json"
	"fmt"     
	"os"
    "io"
	"io/ioutil"	
    "log"
    "net/http"
	"strings"
	"bytes"
)

func geocoder(adrs string) (string, string) {

type Match struct {
	Ads string `json:"matchedAddress"`
	Coord struct {
		Lon float64 `json:"x"`
		Lat float64 `json:"y"`
	} `json:"coordinates"`	
}

type Ress struct {
	Res struct {
	 Matchs []Match `json:"addressMatches"` 
	} `json:"result"`
}	
 
//    adrs := os.Args[1]

    fmt.Println("Address: " + adrs)		
	
    // Create file to receive response from census 
    newFile1, err := os.Create("geocoded.json")
    if err != nil {
        log.Fatal(err)
    }
    defer newFile1.Close()

    url1 := "https://geocoding.geo.census.gov/geocoder/locations/onelineaddress?benchmark=9&format=json&address=" + adrs
	 
    geocoded, err := http.Get(url1)
    defer geocoded.Body.Close()

    // Write bytes from HTTP response to file. geocoded.Body satisfies the reader interface.
    // newFile1 satisfies the writer interface. That enables use of io.Copy which accepts
    // any type that implements reader and writer interface.
	 
    numBytesWritten1, err := io.Copy(newFile1, geocoded.Body)
    if err != nil {
        log.Fatal(err)
    }
    log.Printf("Downloaded %d byte file.\n", numBytesWritten1)
	 
	// Open jsonFile1
	jsonFile1, err := os.Open("geocoded.json")
	if err != nil {
		fmt.Println(err)
	}
	defer jsonFile1.Close()

	// read our opened file as a byte array.
	byteValue1, _ := ioutil.ReadAll(jsonFile1)
	
	// initialize result, which is a Ress struct
	var result Ress
	
	json.Unmarshal(byteValue1, &result)
	
	fmt.Println("Address: " + result.Res.Matchs[0].Ads)		
	fmt.Printf("lon: %f \n",result.Res.Matchs[0].Coord.Lon)
	fmt.Printf("lat: %f \n",result.Res.Matchs[0].Coord.Lat)	

    var lat string = fmt.Sprintf("%f", result.Res.Matchs[0].Coord.Lat)
    var long string = fmt.Sprintf("%f", result.Res.Matchs[0].Coord.Lon)
	
	return lat, long
	
}

func wast(lat, long string) string {

type Site struct {
	Geometry struct {
		Coordinates []float32 `json:"coordinates"`
	} `json:"geometry"`
	Param struct {	
		Cupname string `json:"cleanup_site_name"`
	} `json:"properties"`
}

type Sites struct {
	Sites []Site `json:"features"`
}


//    var lat string = fmt.Sprintf("%f", result.Res.Matchs[0].Coord.Lat)
//    var long string = fmt.Sprintf("%f", result.Res.Matchs[0].Coord.Lon)
	
	 
    url := "https://data.wa.gov/resource/2tkm-ssw6.geojson?%24where=within_circle(location,%20" + lat + ",%20" + long + ",%20250)"

    response, err := http.Get(url)
    defer response.Body.Close()

    // Write bytes from HTTP response to file. response.Body satisfies the reader interface.
    // newFile satisfies the writer interface. That enables use of io.Copy which accepts
    // any type that implements reader and writer interface.

    // Create output file
    newFile, err := os.Create("response.json")
    if err != nil {
         log.Fatal(err)
    }
    defer newFile.Close()
	
    numBytesWritten, err := io.Copy(newFile, response.Body)
    if err != nil {
         log.Fatal(err)
    }
    log.Printf("Downloaded %d byte file.\n", numBytesWritten)

	// Open our jsonFile
	jsonFile, err := os.Open("response.json")
	// if we os.Open returns an error then handle it
	if err != nil {
		fmt.Println(err)
	}

	fmt.Println("Successfully Opened response.json")
	// defer the closing of our jsonFile so that we can parse it later on
	defer jsonFile.Close()

	// read our opened file as a byte array.
	byteValue, _ := ioutil.ReadAll(jsonFile)

	// we initialize our Sites array
	var features Sites

	// we unmarshal our byteArray which contains our
	// jsonFile's content into 'features' which we defined above
	json.Unmarshal(byteValue, &features)
	
	// we iterate through every user within our features array 
	
	road := `{\"type\": \"Feature\",\"geometry\":{\"type\":\"Point\",\"coordinates\":[`
	mark := `\"properties\":{\"marker-size\":\"small\",\"marker-color\":\"#fc0602\",\"marker-symbol\":\"star-stroked\",\"site_name\":\"`
	t := []string{}		
	
	for i := 0; i < len(features.Sites); i++ {
	  	fmt.Println("\n")	
		fmt.Println("Site Name: " + features.Sites[i].Param.Cupname)
	  	fmt.Printf("lon: %f \n",features.Sites[i].Geometry.Coordinates[0])
	  	fmt.Printf("lat: %f \n",features.Sites[i].Geometry.Coordinates[1])		
		t = append(t, road + fmt.Sprintf("%f",features.Sites[i].Geometry.Coordinates[0]) + "," + fmt.Sprintf("%f",features.Sites[i].Geometry.Coordinates[1]) + "]}," + mark + features.Sites[i].Param.Cupname + `\"}},`)			
	}
	
	fmt.Println("\n")
	s := `{\"type\": \"Feature\",\"geometry\":{\"type\":\"Point\",\"coordinates\":[` + long + `,` + lat + `]},\"properties\":{\"marker-size\":\"small\",\"marker-color\":\"#3dd33b\",\"marker-symbol\":\"star\",\"site_name\":\"center\"}}`
	t = append(t,s) 	
	strb := `"{\"type\": \"FeatureCollection\", \"features\": [` +  strings.Join(t, "") + `]}"`
	
	return strb		
}

func gst(strb string) string {
	
    var jsonStr = []byte(`{"description": "A (secret) gist","public": false,"files": {"file1.geojson": { "content":` + strb + `}}}`)

    xurl := "https://api.github.com/gists"	
    req, err := http.NewRequest("POST", xurl, bytes.NewBuffer(jsonStr))
    req.Header.Set("Authorization", "token xxxxxxxxxxx") // The token
    req.Header.Set("Content-Type", "application/json")

    client := &http.Client{}
    resp, err := client.Do(req)
    if err != nil {
        panic(err)
    }
    defer resp.Body.Close()

    if resp.Status == "201 Created" {
        fmt.Println("Success")
        fmt.Println("Go to the following address to access the secret gist")
        fmt.Println(resp.Header.Get("Location"))
    } else {
        fmt.Println("Failed creating secret gist")
    }	
	return strb
}   

func main() {
    adrs := os.Args[1]			
	lat, long := geocoder(adrs)
	fmt.Println(lat, long)

	strb := wast(lat, long) 
	fmt.Println(strb)
	
//	strc := strb
	gst(strb)	
}

