# TF2 Weapon scraper

This script scraps all TF2 weapons of a https://wiki.teamfortress.com/wiki/Weapons into a .json file, i did it as a test project to meet myself with the Golang, i have no clue for what reasons you can use it, but, sure.
## Installation

i didn't compiled it, so you'll have to do everything from the whatever IDE you're using, clone the repo and do the thing
```
go get github.com/gocolly/colly/v2
```

## Usage

Just specify a path (or don't, it will dump .json into the working directory as a output.json) with a
```
go run scraper.go /path/to/file.json
```

it will take some time to complete, something around 3 minutes

there is a complete output.json in the repository if you don't feel like launching it yourself, it is current as of 05.05.2026 (considering Valve updating the game, it will remain the same anyways)

## Plans
i will update it whenever i get some free time and will try to impliment routines and some QOL features, oh and i do intend to use it in another project, it will be more user-oriented than this piece of learning example

## License

MIT License, feel free to do whatever

## Credits
[TF2 Wiki](https://wiki.teamfortress.com/wiki/) for providing with info

Valve software for creating this amazing game

[Colly](https://go-colly.org/) for the framework