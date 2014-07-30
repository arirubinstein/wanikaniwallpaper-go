package main

import "encoding/json"
import "net/http"
import "fmt"
import "log"
import "io/ioutil"
import "container/list"

func GetKanjiForApiKey(key string) *list.List {
	url := "https://www.wanikani.com/api/v1.2/user/" + key + "/kanji"
	fmt.Println("Getting kanji for url " + url)

	response, err := http.Get(url)
	content, err := ioutil.ReadAll(response.Body)
	defer response.Body.Close()

	if err != nil {
		log.Fatalln("Error: http.Get.", err)
	}

	result := list.New()

	var jsonData interface{}
	json.Unmarshal(content, &jsonData)

	requestedInformation := jsonData.(map[string]interface{})

	for _, v := range requestedInformation {
		switch v.(type) {
		case []interface{}:
			for i := range v.([]interface{}) {
				data := v.([]interface{})[i].(map[string]interface{})
				character := data["character"].(string)
				if data["user_specific"] == nil {
					continue
				}
				srs := data["user_specific"].(map[string]interface{})["srs"].(string)
				kanjiStats := KanjiStats{srs}
				kanji := Kanji{character, kanjiStats.Status()}

				result.PushBack(kanji)
			}
		default:
			continue
		}
	}

	return result
}
