package main

import (
	"fmt"
	"sort"
	"strings"

	log "gopkg.in/inconshreveable/log15.v2"
)

// CardList represents the Scryfall List API when retrieving multiple cards
type CardList struct {
	Object     string   `json:"object"`
	TotalCards int      `json:"total_cards"`
	Warnings   []string `json:"warnings"`
	HasMore    bool     `json:"has_more"`
	NextPage   string   `json:"next_page"`
	Data       []Card   `json:"data"`
}

// CardRuling contains an individual ruling on a card
type CardRuling struct {
	Object      string `json:"object"`
	OracleID    string `json:"oracle_id"`
	Source      string `json:"source"`
	PublishedAt string `json:"published_at"`
	Comment     string `json:"comment"`
}

func (ruling *CardRuling) formatRuling() string {
	return fmt.Sprintf("%v: %v", ruling.PublishedAt, ruling.Comment)
}

// CardMetadata contains some extraneous extra information we sometimes retrieve
type CardMetadata struct {
	PreviousPrintings     []string
	PreviousFlavourTexts  []string
	PreviousReminderTexts []string
}

// CardRulingResult represents the JSON returned by the /cards/{}/rulings Scryfall API
type CardRulingResult struct {
	Object  string       `json:"object"`
	HasMore bool         `json:"has_more"`
	Data    []CardRuling `json:"data"`
}

// CommonCard stores the things common to both Card and CardFaces
type CommonCard struct {
	ManaCost        string   `json:"mana_cost"`
	TypeLine        string   `json:"type_line"`
	ColorIndicators []string `json:"color_indicator"`
	OracleText      string   `json:"oracle_text"`
	Power           string   `json:"power"`
	Toughness       string   `json:"toughness"`
	Loyalty         string   `json:"loyalty"`
}

// CardFace represents the individual information for each face of a DFC
type CardFace struct {
	CommonCard
	Object         string `json:"object"`
	Name           string `json:"name"`
	Watermark      string `json:"watermark"`
	Artist         string `json:"artist"`
	IllustrationID string `json:"illustration_id,omitempty"`
}

// Card represents the JSON returned by the /cards Scryfall API
type Card struct {
	CommonCard
	Object        string `json:"object"`
	ID            string `json:"id"`
	OracleID      string `json:"oracle_id"`
	MultiverseIds []int  `json:"multiverse_ids"`
	MtgoID        int    `json:"mtgo_id"`
	MtgoFoilID    int    `json:"mtgo_foil_id"`
	TcgplayerID   int    `json:"tcgplayer_id"`
	Name          string `json:"name"`
	Lang          string `json:"lang"`
	ReleasedAt    string `json:"released_at"`
	URI           string `json:"uri"`
	ScryfallURI   string `json:"scryfall_uri"`
	Layout        string `json:"layout"`
	HighresImage  bool   `json:"highres_image"`
	ImageUris     struct {
		Small      string `json:"small"`
		Normal     string `json:"normal"`
		Large      string `json:"large"`
		Png        string `json:"png"`
		ArtCrop    string `json:"art_crop"`
		BorderCrop string `json:"border_crop"`
	} `json:"image_uris"`
	Cmc           float32    `json:"cmc"`
	Colors        []string   `json:"colors"`
	ColorIdentity []string   `json:"color_identity"`
	CardFaces     []CardFace `json:"card_faces"`
	Legalities    struct {
		Standard  string `json:"standard"`
		Future    string `json:"future"`
		Frontier  string `json:"frontier"`
		Modern    string `json:"modern"`
		Legacy    string `json:"legacy"`
		Pauper    string `json:"pauper"`
		Vintage   string `json:"vintage"`
		Penny     string `json:"penny"`
		Commander string `json:"commander"`
		OneV1     string `json:"1v1"`
		Duel      string `json:"duel"`
		Brawl     string `json:"brawl"`
	} `json:"legalities"`
	Games           []string `json:"games"`
	Reserved        bool     `json:"reserved"`
	Foil            bool     `json:"foil"`
	Nonfoil         bool     `json:"nonfoil"`
	Oversized       bool     `json:"oversized"`
	Promo           bool     `json:"promo"`
	Reprint         bool     `json:"reprint"`
	Set             string   `json:"set"`
	SetName         string   `json:"set_name"`
	SetURI          string   `json:"set_uri"`
	SetSearchURI    string   `json:"set_search_uri"`
	ScryfallSetURI  string   `json:"scryfall_set_uri"`
	RulingsURI      string   `json:"rulings_uri"`
	PrintsSearchURI string   `json:"prints_search_uri"`
	CollectorNumber string   `json:"collector_number"`
	Digital         bool     `json:"digital"`
	Rarity          string   `json:"rarity"`
	FlavourText     string   `json:"flavor_text"`
	IllustrationID  string   `json:"illustration_id"`
	Artist          string   `json:"artist"`
	BorderColor     string   `json:"border_color"`
	Frame           string   `json:"frame"`
	FrameEffect     string   `json:"frame_effect"`
	FullArt         bool     `json:"full_art"`
	Timeshifted     bool     `json:"timeshifted"`
	Colorshifted    bool     `json:"colorshifted"`
	Futureshifted   bool     `json:"futureshifted"`
	StorySpotlight  bool     `json:"story_spotlight"`
	EdhrecRank      int      `json:"edhrec_rank"`
	Usd             string   `json:"usd"`
	Eur             string   `json:"eur"`
	Tix             string   `json:"tix"`
	RelatedUris     struct {
		Gatherer       string `json:"gatherer"`
		TcgplayerDecks string `json:"tcgplayer_decks"`
		Edhrec         string `json:"edhrec"`
		Mtgtop8        string `json:"mtgtop8"`
	} `json:"related_uris"`
	PurchaseUris struct {
		Tcgplayer   string `json:"tcgplayer"`
		Cardmarket  string `json:"cardmarket"`
		Cardhoarder string `json:"cardhoarder"`
	} `json:"purchase_uris"`
	Rulings  []CardRuling
	Metadata CardMetadata
}

// CardCatalog stores the result of the catalog/card-names API call
type CardCatalog struct {
	Object      string   `json:"object"`
	URI         string   `json:"uri"`
	TotalValues int      `json:"total_values"`
	Data        []string `json:"data"`
}

// CardSearchResult stores the result of an advanced Card search API call
type CardSearchResult struct {
	Object     string `json:"object"`
	TotalCards int    `json:"total_cards"`
	HasMore    bool   `json:"has_more"`
	NextPage   string `json:"next_page"`
	Data       []Card `json:"data"`
}

func standardiseColorIndicator(ColorIndicators []string) string {
	expandedColors := map[string]string{"W": "White",
		"U": "Blue",
		"B": "Black",
		"R": "Red",
		"G": "Green"}
	mappedColors := map[string]int{"White": 0,
		"Blue":  1,
		"Black": 2,
		"Red":   3,
		"Green": 4}

	var colorWords []string
	for _, color := range ColorIndicators {
		colorWords = append(colorWords, expandedColors[color])
	}

	sort.Slice(colorWords, func(i, j int) bool {
		return mappedColors[colorWords[i]] < mappedColors[colorWords[j]]
	})

	return "[" + strings.Join(colorWords, "/") + "]"
}

func normaliseCardName(input string) string {
	ret := nonAlphaRegex.ReplaceAllString(strings.ToLower(input), "")
	// log.Debug("Normalising", "Input", input, "Output", ret)
	return ret
}

func formatManaCost(input string) string {
	return input
}

func replaceManaCostForSlack(input string) string {
	manaString := strings.Replace(input, "{1000000}", ":mana-1000000-1::mana-1000000-2::mana-1000000-3::mana-1000000-4:", -1)
	for _, match := range emojiRegex.FindAllString(manaString, -1) {
		replacement := strings.Replace(match, "{", ":mana-", -1)
		replacement = strings.Replace(replacement, "}", ":", -1)
		replacement = strings.Replace(replacement, "/", "", -1)
		manaString = strings.Replace(manaString, match, replacement, 1)
	}
	return manaString
}

// TODO: Have a command to see all printing information
func (card *Card) formatExpansions() string {
	var ret []string
	if card.Name != "Plains" && card.Name != "Island" && card.Name != "Swamp" && card.Name != "Mountain" && card.Name != "Forest" {
		if len(card.Metadata.PreviousPrintings) > 0 {
			if len(card.Metadata.PreviousPrintings) < 10 {
				ret = card.Metadata.PreviousPrintings
			} else {
				ret = card.Metadata.PreviousPrintings[:5]
				ret = append(ret, "[...]")
			}
		}
	}
	log.Warn("FE", "card", card)
	ret = append(ret, fmt.Sprintf("%s-%s", strings.ToUpper(card.Set), strings.ToUpper(card.Rarity[0:1])))
	return strings.Join(sliceUniqMap(ret), ",")
}

func (card *Card) formatLegalities() string {
	var ret []string
	switch card.Legalities.Vintage {
	case "legal":
		ret = append(ret, "Vin")
	case "restricted":
		ret = append(ret, "VinRes")
	case "banned":
		ret = append(ret, "VinBan")
	}
	switch card.Legalities.Legacy {
	case "legal":
		ret = append(ret, "Leg")
	case "restricted":
		ret = append(ret, "LegRes")
	case "banned":
		ret = append(ret, "LegBan")
	}
	switch card.Legalities.Modern {
	case "legal":
		ret = append(ret, "Mod")
	case "restricted":
		ret = append(ret, "ModRes")
	case "banned":
		ret = append(ret, "ModBan")
	}
	switch card.Legalities.Standard {
	case "legal":
		ret = append(ret, "Std")
	case "restricted":
		ret = append(ret, "StdRes")
	case "banned":
		ret = append(ret, "StdBan")
	}
	return strings.Join(ret, ",")
}

func lookupUniqueNamePrefix(input string) string {
	ncn := normaliseCardName(input)
	log.Debug("in lookupUniqueNamePrefix", "Input", input, "NCN", ncn, "Length of CN", len(cardNames))
	var err error
	if len(cardNames) == 0 {
		log.Debug("In lookupUniqueNamePrefix -- Importing")
		cardNames, err = importCardNames(false)
		if err != nil {
			log.Warn("Error importing card names", "Error", err)
			return ""
		}
	}
	//c := cardNames[:0]
	var c []string
	for _, x := range cardNames {
		if strings.HasPrefix(normaliseCardName(x), ncn) {
			log.Debug("In lookupUniqueNamePrefix", "Gottem", x)
			c = append(c, x)
		}
	}
	log.Debug("In lookupUniqueNamePrefix", "C", c)
	if len(c) == 1 {
		return c[0]
	}
	// Look for something legendary-ish
	var i int
	var j string
	for _, x := range c {
		if strings.Contains(x, ",") || strings.Contains(x, "the") {
			i++
			j = x
		}
	}
	if i == 1 {
		cardUniquePrefixHits.Add(1)
		return j
	}
	return ""
}
