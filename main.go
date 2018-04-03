package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"

	"github.com/gorilla/mux"
)

func main() {
	router := mux.NewRouter()

	fs := http.FileServer(http.Dir("static"))

	router.HandleFunc("/createlist", CreateList).Methods("POST")
	router.Handle("/", fs)

	http.ListenAndServe("localhost:8080", router)
}

func CreateList(w http.ResponseWriter, r *http.Request) {
	words := strings.Split(r.FormValue("list"), "\n")

	list := []string{"Kanji", "\t", "Kana", "\t", "English", "\t", "Part of Speech", "\n"}

	for _, w := range words {
		resp, err := http.Get("http://jisho.org/api/v1/search/words?keyword=" + url.QueryEscape(w))
		if err != nil {
			fmt.Println(err)
			continue
		}
		defer resp.Body.Close()

		decoder := json.NewDecoder(resp.Body)
		var jisho Jisho
		decoder.Decode(&jisho)

		for i, p := range jisho.Data[0].Senses[0].PartsOfSpeech {
			p = strings.ToLower(p)
			switch p {
			case "godan verb with u ending":
				jisho.Data[0].Senses[0].PartsOfSpeech[i] = "5v-う"
			case "godan verb with ku ending":
				jisho.Data[0].Senses[0].PartsOfSpeech[i] = "5v-く"
			case "godan verb with su ending":
				jisho.Data[0].Senses[0].PartsOfSpeech[i] = "5v-す"
			case "godan verb with tsu ending":
				jisho.Data[0].Senses[0].PartsOfSpeech[i] = "5v-つ"
			case "godan verb with nu ending":
				jisho.Data[0].Senses[0].PartsOfSpeech[i] = "5v-ぬ"
			case "godan verb with bu ending":
				jisho.Data[0].Senses[0].PartsOfSpeech[i] = "5v-ぶ"
			case "godan verb with mu ending":
				jisho.Data[0].Senses[0].PartsOfSpeech[i] = "5v-む"
			case "godan verb with ru ending":
				jisho.Data[0].Senses[0].PartsOfSpeech[i] = "5v-る"
			case "ichidan verb":
				jisho.Data[0].Senses[0].PartsOfSpeech[i] = "1v"
			case "suru verb":
				jisho.Data[0].Senses[0].PartsOfSpeech[i] = "v-する"
			case "transitive verb":
				jisho.Data[0].Senses[0].PartsOfSpeech[i] = "vt"
			case "intransitive verb":
				jisho.Data[0].Senses[0].PartsOfSpeech[i] = "vi"
			case "noun":
				jisho.Data[0].Senses[0].PartsOfSpeech[i] = "n"
			case "i-adjective":
				jisho.Data[0].Senses[0].PartsOfSpeech[i] = "い-adj"
			case "na-adjective":
				jisho.Data[0].Senses[0].PartsOfSpeech[i] = "な-adj"
			default:
				jisho.Data[0].Senses[0].PartsOfSpeech[i] = p
			}
		}

		list = append(list, []string{jisho.Data[0].Japanese[0].Word, "\t", jisho.Data[0].Japanese[0].Reading, "\t", strings.Join(jisho.Data[0].Senses[0].EnglishDefinitions, ", "), "\t", strings.Join(jisho.Data[0].Senses[0].PartsOfSpeech, ", "), "\n"}...)
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte(strings.Join(list, "")))
}

type Jisho struct {
	Meta struct {
		Status int `json:"status"`
	} `json:"meta"`
	Data []struct {
		IsCommon bool     `json:"is_common"`
		Tags     []string `json:"tags"`
		Japanese []struct {
			Word    string `json:"word"`
			Reading string `json:"reading"`
		} `json:"japanese"`
		Senses []struct {
			EnglishDefinitions []string `json:"english_definitions"`
			PartsOfSpeech      []string `json:"parts_of_speech"`
			Links              []string `json:"links"`
			Tags               []string `json:"tags"`
			Restrictions       []string `json:"restrictions"`
			SeeAlso            []string `json:"see_also"`
			Antonyms           []string `json:"antonyms"`
			Source             []string `json:"source"`
			Info               []string `json:"info"`
		} `json:"senses"`
		Attribution struct {
			Jmdict   bool `json:"jmdict"`
			Jmnedict bool `json:"jmnedict"`
			Dbpedia  bool `json:"dbpedia"`
		} `json:"attribution"`
	} `json:"data"`
}
