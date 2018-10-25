package main

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"os"
    "net/http"
    "log"
    "strings"		
)

// Sites struct - array of features
type Sites struct {
	Sites []Site `json:"features"`
}

// Site struct - name, type geometry links list
type Site struct {
	Type   string `json:"type"`
	Geometry Geometry `json:"geometry"`
	Param Param `json:"properties"`
}

// Geometry struct - coordinates
type Geometry struct {
	Type string `json:"type"`
	Coordinates []float32 `json:"coordinates"` 
}

// Param struct selected site information
type Param struct {	
	Cupname string `json:"cleanup_site_name"`
}

func main() {

     // Create output file
     newFile, err := os.Create("response.json")
     if err != nil {
          log.Fatal(err)
     }
     defer newFile.Close()

     var lat string = "47.59"
     var long string = "-122.33"
	 
     url := "https://data.wa.gov/resource/2tkm-ssw6.geojson?%24where=within_circle(location,%20" + lat + ",%20" + long + ",%20250)"

     response, err := http.Get(url)
     defer response.Body.Close()

     // Write bytes from HTTP response to file.
     // response.Body satisfies the reader interface.
     // newFile satisfies the writer interface.
     // That allows us to use io.Copy which accepts
     // any type that implements reader and writer interface
	 
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
	
	road := `{\"type": \"Feature\",\"geometry\":{\"type\":\"Point\",\"coordinates\":[`
	mark := `\"properties\":{\"marker-size\":\"small\",\"marker-color\":\"#fc0602\",\"marker-symbol\":\"star-stroked\",\"site_name\":\"`
	t := []string{}		
	
	for i := 0; i < len(features.Sites); i++ {
	  	fmt.Println("\n")	
		fmt.Println("Site Name: " + features.Sites[i].Param.Cupname)
	  	fmt.Printf("lon: %f \n",features.Sites[i].Geometry.Coordinates[0])
	  	fmt.Printf("lat: %f \n",features.Sites[i].Geometry.Coordinates[1])		
		t = append(t, road + fmt.Sprintf("%f",features.Sites[i].Geometry.Coordinates[0]) + "," + fmt.Sprintf("%f",features.Sites[i].Geometry.Coordinates[1]) + "]}," + mark + features.Sites[i].Param.Cupname + "\"}},")			
	}
	
	fmt.Println("\n")
	s := `{\"type": \"Feature\",\"geometry\":{\"type\":\"Point\",\"coordinates\":[` + long + `,` + lat + `]},\"properties\":{\"marker-size\":\"small\",\"marker-color\":\"#fc0602\",\"marker-symbol\":\"star-stroked\",\"site_name\":\"center"}}`
	t = append(t,s) 

	fmt.Println( `{"type": "FeatureCollection", "features": [` +  strings.Join(t, "") + `]}`)		
}

