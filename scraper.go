package main

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/gocolly/colly/v2"
)

func main() {
	//out of curiosity, i wonder how much this process takes now
	start := time.Now()

	if len(os.Args) > 2 {
		fmt.Println("you like arguments, dont you?")
	}
	//i assume its something you would call initialization (theres a chance i didnt spelled that correctly)
	c := colly.NewCollector(
		colly.Async(true),
	)

	c.Limit(&colly.LimitRule{
		DomainGlob:  "*",
		Parallelism: 25,
	})

	gridItteration := 0
	type Weapon struct {
		ID       int      `json:"id"`
		Name     string   `json:"name"`
		WikiURL  string   `json:"wikiurl"`
		ImageURL string   `json:"imageurl"`
		Slot     string   `json:"slot"`
		IsSkin   bool     `json:"isskin"`
		Classes  []string `json:"classes"`
	}
	var weapons []Weapon

	c.OnHTML("table.wikitable.grid th", func(e *colly.HTMLElement) {

		weaponName := e.ChildText("a b")
		if weaponName == "" {
			return //skipping empty headers
		}

		gridItteration += 1
		//for some reason, link on the wiki arent including domain, i guess
		wikiLink := "https://wiki.teamfortress.com" + e.ChildAttr("a", "href")
		imgSrc := "https://wiki.teamfortress.com" + e.ChildAttr("a img", "src")
		values := scrapeWeaponItself(wikiLink)

		var humbleJson Weapon

		humbleJson.ID = gridItteration
		humbleJson.Name = weaponName
		humbleJson.WikiURL = wikiLink
		humbleJson.ImageURL = imgSrc

		//actually, nevermind, i like it, its so cursed im enjoying every second of it
		if val, ok := values["isSkin"].(bool); ok {
			humbleJson.IsSkin = val
		} else {
			humbleJson.IsSkin = false
		}

		//holy syntax
		val, ok := values["classes"].([]string)
		if ok {
			humbleJson.Classes = val
		}

		if val, ok := values["slot"].(string); ok {
			humbleJson.Slot = val
		}

		weapons = append(weapons, humbleJson)
		fmt.Printf("\rLoaded: %d", gridItteration)
	})

	//took it right off the wiki and im NOT going to change it, it stays
	c.OnRequest(func(r *colly.Request) {
		fmt.Println("Visiting", r.URL)
	})

	c.OnError(func(r *colly.Response, err error) {
		fmt.Println("Request URL:", r.Request.URL, "failed with response:", r, "\nError:", err)
	})

	c.Visit("https://wiki.teamfortress.com/wiki/Weapons")

	c.Wait()

	fmt.Println()

	data, err := json.MarshalIndent(weapons, "", " ")
	if err != nil {
		//i love the panic function, it is hilarious
		panic(err)
	}

	if len(os.Args) >= 2 {
		err = os.WriteFile(os.Args[1], data, 0644)
		if err != nil {
			fmt.Println(err)
			fmt.Println("saving as output.json in working folder")
			err = os.WriteFile("output.json", data, 0644)
			if err != nil {
				fmt.Println("tried to save in working directory, failed miserably")
				panic(err)
			}
		}
	} else {
		err = os.WriteFile("output.json", data, 0644)
		if err != nil {
			fmt.Println("tried to save in working directory, failed miserably")
			panic(err)
		}
	}
	fmt.Println("done! have a good one")
	elapsed := time.Since(start)
	fmt.Println("this took", elapsed)
}

func scrapeWeaponItself(linkToWeapon string) map[string]any {
	c := colly.NewCollector(
		colly.Async(true),
	)

	c.Limit(&colly.LimitRule{
		DomainGlob:  "*",
		Parallelism: 25,
	})

	//i dont like it tbh
	var values = make(map[string]any)

	c.OnHTML("div.mw-parser-output p", func(e *colly.HTMLElement) {
		if strings.Contains(e.Text, "This weapon functions identically") {
			values["isSkin"] = true
		}
	})

	c.OnHTML("table.infobox.item-infobox.weapon-infobox tr", func(e *colly.HTMLElement) {
		if e.ChildText("td.infobox-label") == "Used by:" {
			var classes []string

			//apparently, weapons could be used by multiple classes, suprising for a person who spent
			//2k hours of his life in tf2
			e.ForEach("td a", func(_ int, el *colly.HTMLElement) {
				//also, i didnt know _ stands for unused values, good to know!
				classes = append(classes, el.Text)
			})
			values["classes"] = classes
		} else if e.ChildText("td.infobox-label") == "Slot:" {
			values["slot"] = e.ChildText("td a")
		}
	})

	c.OnError(func(r *colly.Response, err error) {
		fmt.Println("Request URL:", r.Request.URL, "failed with response:", r, "\nError:", err)
	})

	c.Visit(linkToWeapon)
	c.Wait()
	return values
}
