package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"

	"github.com/signintech/gopdf"
)

// Version 0.1 --> Generates a sample invoice with several defects
// Version 0.2 --> Makes the background image look good
// Version 0.3 --> Find a font with bold and the simbol €
// Version 0.4 --> Use only the fonts FreeSans and FreeSansBold
// Version 0.5 --> Put it all on a better place
// Version 0.6 --> All fixed data should be a variable (this way we prepare for reading from a file)
// Version 0.7 --> Load main data (List Details and other data) from a .tsv (Tab-Separated Values)
// Version 0.8 --> Load static data from a .dat, basically all default variables not included in the .tsv

type tPoint struct {
	x, y float64
}

type tListLine struct {
	concept []string // Multi-line concept
	amount  string
	price   string
	total   int      // Total amount in cents
}

// Default values
var title string           = "Default Invoice"
var author string          = "[Your Name Here]"
var version string         = "0.9.8"
var creator string         = "genFact v" + version
var producer string        = "signintech/gopdf"
var outputFileName string  = "output/default.pdf"
var documentWidth float64  = 595.28
var documentHeight float64 = 841.89
var fontRegular string     = "fonts/FreeSans.ttf"
var fontBold string        = "fonts/FreeSansBold.ttf"
var backgroundImage string = "img/Background_0.png"
var myInfoRect1 tPoint     = tPoint{x: 7.0,   y: 10.0}
var myInfoRect2 tPoint     = tPoint{x: 165.0, y: 79.0}
var myInfoStart tPoint     = tPoint{x: 10.0,  y: 24.0}
var myInfoLine float64     = 17.0
var myInfoFontSize float64 = 14.0
var myInfoData []string    = []string{"YOUR DEFAULT NAME", "DEFAULT ADDRESS", "00000 - CITY", "00.000.000-A"}
var invInfoRect1 tPoint    = tPoint{x: 450.0, y: 10.0}
var invInfoRect2 tPoint    = tPoint{x: 580.0, y: 23.0}
var invInfoStart1 tPoint   = tPoint{x: 472.7, y: 21.0}
var invInfoStart2 tPoint   = tPoint{x: 451.0, y: 32.0}
var invInfoStart3 tPoint   = tPoint{x: 530.0, y: 32.0}
var invInfoLine float64    = 11.0
var invInfoFontSz1 float64 = 13.0
var invInfoFontSz2 float64 = 10.0
var invInfoFontSz3 float64 = 10.0
var invInfoTitle string    = " I N V O I C E "
var invInfoText []string   = []string{"Page", "Invoice nº", "Date", "Client nº"}
var page int               = 1
var invNum string          = "2025-01"
var date string            = "01/01/2025"
var cliNum int             = 100
var cliInfoStart1 tPoint   = tPoint{x: 10.0, y: 100.0}
var cliInfoStart2 tPoint   = tPoint{x: 10.0, y: 114.0}
var cliInfoStart3 tPoint   = tPoint{x: 10.0, y: 142.0}
var cliInfoLine float64    = 14.0
var cliInfoFontSz1 float64 = 14.0
var cliInfoFontSz2 float64 = 13.0
var cliInfoFontSz3 float64 = 13.0
var cliInfoName string     = "DEFAULT CLIENT"
var cliInfoData []string   = []string{"Client's Address", "City (Country)"}
var cliInfoCode string     = "A0000000"
var listDetRect1 tPoint    = tPoint{x:  10.0, y: 150.0}
var listDetRect2 tPoint    = tPoint{x: 580.0, y: 162.0}
var listHFontSz float64    = 10
var listHStartY float64    = 159.0
var listHStartXs []float64 = []float64{11.0, 402.0, 475.0, 545.0}
var listHText []string     = []string{"Description", "Amount.", "Price", "Import"}
var listDFontSz float64    = 10.0
var listDStartY float64    = 165.0
var listDStartXs []float64 = []float64{11.0, 430.0, 503.0, 579.0}
var listDetLine float64    = 12.0
var totalNoTax int         = 0
var coin string            = "€"
var coinFront bool         = false
var listDet []tListLine    = []tListLine{
	tListLine{concept: []string{"Default description/concept"}, amount: "-", price: "9,99 €", total: 999},
	tListLine{concept: []string{"Multi", "Line", "Description", "Example"}, amount: "1", price: "1,25 €", total: 125},
}
var deletedDefaultListDetails = false // Control variable
var listDetSepX1 float64   = 10.0
var listDetSepX2 float64   = 580.0
var taxes float64          = 10.0
var taxMsgStartX float64   = 100.0
var taxText1 string        = "\"NO TAXES APPLY\""
var taxText2 string        = "\"10% TAX\""
var subTotalText string    = "Total Base...............:"
var subTotalStartX float64 = 420.0
var totAmoStartX float64   = 578.0
var totalText string       = "Total Invoice...........:"

func processFTsv(fTsv *os.File) error {

	scn := bufio.NewScanner(fTsv)
	var line int
	for scn.Scan() {
		line++
		l := scn.Text()
		l = strings.TrimSpace(l)
		t := strings.Split(l, "\t")

		switch t[0] {
		case "":
			fmt.Println("Skipping empty line", line)
		case "author":
			if len(t) < 2 {
				return fmt.Errorf("Too few arguments on line %d", line)
			} else if len(t) > 2 {
				return fmt.Errorf("Too many arguments on line %d", line)
			}
			fmt.Println("Author will now be", t[1])
			author = t[1]
		case "staticFile":
			if len(t) < 2 {
				return fmt.Errorf("Too few arguments on line %d", line)
			} else if len(t) > 2 {
				return fmt.Errorf("Too many arguments on line %d", line)
			}
			fmt.Println("Loading static data file", t[1])
			var err error
			var fDat *os.File
			fDat, err = os.Open(t[1])
			if err != nil {
				return fmt.Errorf("Could not open static data file (\"%s\") on line %d", t[1], line)
			}
			defer fDat.Close()
			err = processFDat(fDat)
			if err != nil {
				return err
			}
		case "outputFileName":
			if len(t) < 2 {
				return fmt.Errorf("Too few arguments on line %d", line)
			} else if len(t) > 2 {
				return fmt.Errorf("Too many arguments on line %d", line)
			}
			fmt.Println("Output file will now be", t[1])
			outputFileName = t[1]
		case "invNum", "invoiceNumber":
			if len(t) < 2 {
				return fmt.Errorf("Too few arguments on line %d", line)
			} else if len(t) > 2 {
				return fmt.Errorf("Too many arguments on line %d", line)
			}
			fmt.Println("Invoice number will now be", t[1])
			invNum = t[1]
		case "date", "invoiceDate":
			if len(t) < 2 {
				return fmt.Errorf("Too few arguments on line %d", line)
			} else if len(t) > 2 {
				return fmt.Errorf("Too many arguments on line %d", line)
			}
			fmt.Println("Invoice date will now be", t[1])
			date = t[1]
		case "title":
			if len(t) < 2 {
				return fmt.Errorf("Too few arguments on line %d", line)
			} else if len(t) > 2 {
				return fmt.Errorf("Too many arguments on line %d", line)
			}
			fmt.Println("Base title will now be", t[1])
			title = t[1]
			title = strings.ReplaceAll(title, "[date]", date)
			title = strings.ReplaceAll(title, "[invNum]", invNum)
		case "cliNum", "clientNumber":
			if len(t) < 2 {
				return fmt.Errorf("Too few arguments on line %d", line)
			} else if len(t) > 2 {
				return fmt.Errorf("Too many arguments on line %d", line)
			}
			fmt.Println("Client number will now be", t[1])
			var err error
			cliNum, err = strconv.Atoi(t[1])
			if err != nil {
				return fmt.Errorf("Error converting client number (\"%s\") into an integer at line %d", t[1], line)
			}
		case "coin":
			if len(t) < 2 {
				return fmt.Errorf("Too few arguments on line %d", line)
			} else if len(t) > 2 {
				return fmt.Errorf("Too many arguments on line %d", line)
			}
			fmt.Println("Invoice coin will now be", t[1])
			coin = t[1]
		case "coinFront":
			if len(t) < 2 {
				return fmt.Errorf("Too few arguments on line %d", line)
			} else if len(t) > 2 {
				return fmt.Errorf("Too many arguments on line %d", line)
			}
			if strings.ToUpper(t[1]) == "TRUE" {
				coinFront = true
				fmt.Println("Coin symbol will now be displayed at the front")
			} else {
				coinFront = false
				fmt.Println("Coin symbol will now be displayed at the back")
			}
		case "taxes":
			if len(t) < 2 {
				return fmt.Errorf("Too few arguments on line %d", line)
			} else if len(t) > 2 {
				return fmt.Errorf("Too many arguments on line %d", line)
			}
			fmt.Println("Taxes are now set to", t[1])
			var err error
			taxes, err = strconv.ParseFloat(t[1], 64)
			if err != nil {
				return fmt.Errorf("Error converting taxes number (\"%s\") into a decimal number at line %d", t[1], line)
			}
		case "listDetail":
			if len(t) < 5 {
				return fmt.Errorf("Too few arguments on line %d", line)
			}
			if !deletedDefaultListDetails {
				listDet = []tListLine{}
				deletedDefaultListDetails = true
			}
			fmt.Println("List Detail")
			var tmp tListLine
			var err error
			l := len(t)
			tmp.total, err = strconv.Atoi(t[l-1])
			if err != nil {
				return fmt.Errorf("Error converting list detail total (\"%s\") into an integer at line %d", t[l-1], line)
			}
			tmp.price = t[l-2]
			tmp.amount = t[l-3]
			tmp.concept = t[1:l-3]
			listDet = append(listDet, tmp)
		default:
			if t[0][0] == '#' {
				fmt.Println("Comment at line", line)
			} else {
				return fmt.Errorf("⚠️ Unrecognized command \"%s\" on line %d", t[0], line)
			}
		}
	}

	return nil
}

func processFDat(fDat *os.File) error {

	scn := bufio.NewScanner(fDat)
	var line int
	for scn.Scan() {
		line++
		l := scn.Text()
		t := strings.Split(l, "=")
		t[0] = strings.TrimSpace(t[0])
		if len(t) > 1 {
			t[1] = strings.TrimSpace(t[1])
		}

		switch t[0] {
		case "":
			fmt.Println(" → Skipping empty line", line)
		case "documentWidth":
			if len(t) < 2 {
				return fmt.Errorf(" → Too few arguments on line %d", line)
			} else if len(t) > 2 {
				return fmt.Errorf(" → Too many arguments on line %d", line)
			}
			fmt.Println(" → Document width will now be", t[1])
			var err error
			documentWidth, err = strconv.ParseFloat(t[1], 64)
			if err != nil {
				return fmt.Errorf(" → Error converting document width number (\"%s\") into a decimal number at line %d", t[1], line)
			}
		case "documentHeight":
			if len(t) < 2 {
				return fmt.Errorf(" → Too few arguments on line %d", line)
			} else if len(t) > 2 {
				return fmt.Errorf(" → Too many arguments on line %d", line)
			}
			fmt.Println(" → Document height will now be", t[1])
			var err error
			documentHeight, err = strconv.ParseFloat(t[1], 64)
			if err != nil {
				return fmt.Errorf(" → Error converting document height number (\"%s\") into a decimal number at line %d", t[1], line)
			}
		case "fontRegular":
			if len(t) < 2 {
				return fmt.Errorf(" → Too few arguments on line %d", line)
			} else if len(t) > 2 {
				return fmt.Errorf(" → Too many arguments on line %d", line)
			}
			fmt.Println(" → The regular font file will be", t[1])
			fontRegular = t[1]
		case "fontBold":
			if len(t) < 2 {
				return fmt.Errorf(" → Too few arguments on line %d", line)
			} else if len(t) > 2 {
				return fmt.Errorf(" → Too many arguments on line %d", line)
			}
			fmt.Println(" → The bold font file will be", t[1])
			fontBold = t[1]
		case "backgroundImage":
			if len(t) < 2 {
				return fmt.Errorf(" → Too few arguments on line %d", line)
			} else if len(t) > 2 {
				return fmt.Errorf(" → Too many arguments on line %d", line)
			}
			fmt.Println(" → The background image will be", t[1])
			backgroundImage = t[1]
		case "myInfoRect1":
			if len(t) < 2 {
				return fmt.Errorf(" → Too few arguments on line %d", line)
			} else if len(t) > 2 {
				return fmt.Errorf(" → Too many arguments on line %d", line)
			}
			p := strings.Split(t[1], "; ")
			if len(p) < 2 {
				return fmt.Errorf(" → Too few numbers on line %d", line)
			} else if len(p) > 2 {
				return fmt.Errorf(" → Too many numbers on line %d", line)
			}
			fmt.Println(" → My information rectangle will start at x =", p[0], "and y =", p[1])
			var err error
			var x, y float64
			x, err = strconv.ParseFloat(strings.TrimSpace(p[0]), 64)
			if err != nil {
				return fmt.Errorf(" → Error converting my information rectangle start X number (\"%s\") into a decimal number at line %d", p[0], line)
			}
			y, err = strconv.ParseFloat(strings.TrimSpace(p[1]), 64)
			if err != nil {
				return fmt.Errorf(" → Error converting my information rectangle start Y number (\"%s\") into a decimal number at line %d", p[1], line)
			}
			myInfoRect1 = tPoint{x: x, y: y}
		case "myInfoRect2":
			if len(t) < 2 {
				return fmt.Errorf(" → Too few arguments on line %d", line)
			} else if len(t) > 2 {
				return fmt.Errorf(" → Too many arguments on line %d", line)
			}
			p := strings.Split(t[1], "; ")
			if len(p) < 2 {
				return fmt.Errorf(" → Too few numbers on line %d", line)
			} else if len(p) > 2 {
				return fmt.Errorf(" → Too many numbers on line %d", line)
			}
			fmt.Println(" → My information rectangle will end at x =", p[0], "and y =", p[1])
			var err error
			var x, y float64
			x, err = strconv.ParseFloat(strings.TrimSpace(p[0]), 64)
			if err != nil {
				return fmt.Errorf(" → Error converting my information rectangle end X number (\"%s\") into a decimal number at line %d", p[0], line)
			}
			y, err = strconv.ParseFloat(strings.TrimSpace(p[1]), 64)
			if err != nil {
				return fmt.Errorf(" → Error converting my information rectangle end Y number (\"%s\") into a decimal number at line %d", p[1], line)
			}
			myInfoRect2 = tPoint{x: x, y: y}
		case "myInfoStart":
			if len(t) < 2 {
				return fmt.Errorf(" → Too few arguments on line %d", line)
			} else if len(t) > 2 {
				return fmt.Errorf(" → Too many arguments on line %d", line)
			}
			p := strings.Split(t[1], "; ")
			if len(p) < 2 {
				return fmt.Errorf(" → Too few numbers on line %d", line)
			} else if len(p) > 2 {
				return fmt.Errorf(" → Too many numbers on line %d", line)
			}
			fmt.Println(" → My information text will start at x =", p[0], "and y =", p[1])
			var err error
			var x, y float64
			x, err = strconv.ParseFloat(strings.TrimSpace(p[0]), 64)
			if err != nil {
				return fmt.Errorf("Error converting my information text start X number (\"%s\") into a decimal number at line %d", p[0], line)
			}
			y, err = strconv.ParseFloat(strings.TrimSpace(p[1]), 64)
			if err != nil {
				return fmt.Errorf(" → Error converting my information text start Y number (\"%s\") into a decimal number at line %d", p[1], line)
			}
			myInfoStart = tPoint{x: x, y: y}
		case "myInfoLine":
			if len(t) < 2 {
				return fmt.Errorf(" → Too few arguments on line %d", line)
			} else if len(t) > 2 {
				return fmt.Errorf(" → Too many arguments on line %d", line)
			}
			fmt.Println(" → My information text line height will now be", t[1])
			var err error
			myInfoLine, err = strconv.ParseFloat(t[1], 64)
			if err != nil {
				return fmt.Errorf(" → Error converting my information text line height number (\"%s\") into a decimal number at line %d", t[1], line)
			}
		case "myInfoFontSize":
			if len(t) < 2 {
				return fmt.Errorf(" → Too few arguments on line %d", line)
			} else if len(t) > 2 {
				return fmt.Errorf(" → Too many arguments on line %d", line)
			}
			fmt.Println(" → My information font size will now be", t[1])
			var err error
			myInfoFontSize, err = strconv.ParseFloat(t[1], 64)
			if err != nil {
				return fmt.Errorf(" → Error converting my information font size number (\"%s\") into a decimal number at line %d", t[1], line)
			}
		case "myInfoData":
			if len(t) < 2 {
				return fmt.Errorf(" → Too few arguments on line %d", line)
			} else if len(t) > 2 {
				return fmt.Errorf(" → Too many arguments on line %d", line)
			}
			fmt.Println(" → My information data will now be:", t[1])
			p := strings.Split(t[1], "; ")
			myInfoData = p
		case "invInfoRect1":
			if len(t) < 2 {
				return fmt.Errorf(" → Too few arguments on line %d", line)
			} else if len(t) > 2 {
				return fmt.Errorf(" → Too many arguments on line %d", line)
			}
			p := strings.Split(t[1], "; ")
			if len(p) < 2 {
				return fmt.Errorf(" → Too few numbers on line %d", line)
			} else if len(p) > 2 {
				return fmt.Errorf(" → Too many numbers on line %d", line)
			}
			fmt.Println(" → Invoice information rectangle will start at x =", p[0], "and y =", p[1])
			var err error
			var x, y float64
			x, err = strconv.ParseFloat(strings.TrimSpace(p[0]), 64)
			if err != nil {
				return fmt.Errorf(" → Error converting invoice information rectangle start X number (\"%s\") into a decimal number at line %d", p[0], line)
			}
			y, err = strconv.ParseFloat(strings.TrimSpace(p[1]), 64)
			if err != nil {
				return fmt.Errorf(" → Error converting invoice information rectangle start Y number (\"%s\") into a decimal number at line %d", p[1], line)
			}
			invInfoRect1 = tPoint{x: x, y: y}
		case "invInfoRect2":
			if len(t) < 2 {
				return fmt.Errorf(" → Too few arguments on line %d", line)
			} else if len(t) > 2 {
				return fmt.Errorf(" → Too many arguments on line %d", line)
			}
			p := strings.Split(t[1], "; ")
			if len(p) < 2 {
				return fmt.Errorf(" → Too few numbers on line %d", line)
			} else if len(p) > 2 {
				return fmt.Errorf(" → Too many numbers on line %d", line)
			}
			fmt.Println(" → Invoice information rectangle will end at x =", p[0], "and y =", p[1])
			var err error
			var x, y float64
			x, err = strconv.ParseFloat(strings.TrimSpace(p[0]), 64)
			if err != nil {
				return fmt.Errorf(" → Error converting invoice information rectangle end X number (\"%s\") into a decimal number at line %d", p[0], line)
			}
			y, err = strconv.ParseFloat(strings.TrimSpace(p[1]), 64)
			if err != nil {
				return fmt.Errorf(" → Error converting invoice information rectangle end Y number (\"%s\") into a decimal number at line %d", p[1], line)
			}
			invInfoRect2 = tPoint{x: x, y: y}
		case "invInfoStart1":
			if len(t) < 2 {
				return fmt.Errorf(" → Too few arguments on line %d", line)
			} else if len(t) > 2 {
				return fmt.Errorf(" → Too many arguments on line %d", line)
			}
			p := strings.Split(t[1], "; ")
			if len(p) < 2 {
				return fmt.Errorf(" → Too few numbers on line %d", line)
			} else if len(p) > 2 {
				return fmt.Errorf(" → Too many numbers on line %d", line)
			}
			fmt.Println(" → Invoice information 1 start at x =", p[0], "and y =", p[1])
			var err error
			var x, y float64
			x, err = strconv.ParseFloat(strings.TrimSpace(p[0]), 64)
			if err != nil {
				return fmt.Errorf(" → Error converting invoice information 1 start X number (\"%s\") into a decimal number at line %d", p[0], line)
			}
			y, err = strconv.ParseFloat(strings.TrimSpace(p[1]), 64)
			if err != nil {
				return fmt.Errorf(" → Error converting invoice information 1 start Y number (\"%s\") into a decimal number at line %d", p[1], line)
			}
			invInfoStart1 = tPoint{x: x, y: y}
		case "invInfoStart2":
			if len(t) < 2 {
				return fmt.Errorf(" → Too few arguments on line %d", line)
			} else if len(t) > 2 {
				return fmt.Errorf(" → Too many arguments on line %d", line)
			}
			p := strings.Split(t[1], "; ")
			if len(p) < 2 {
				return fmt.Errorf(" → Too few numbers on line %d", line)
			} else if len(p) > 2 {
				return fmt.Errorf(" → Too many numbers on line %d", line)
			}
			fmt.Println(" → Invoice information 2 start at x =", p[0], "and y =", p[1])
			var err error
			var x, y float64
			x, err = strconv.ParseFloat(strings.TrimSpace(p[0]), 64)
			if err != nil {
				return fmt.Errorf(" → Error converting invoice information 2 start X number (\"%s\") into a decimal number at line %d", p[0], line)
			}
			y, err = strconv.ParseFloat(strings.TrimSpace(p[1]), 64)
			if err != nil {
				return fmt.Errorf(" → Error converting invoice information 2 start Y number (\"%s\") into a decimal number at line %d", p[1], line)
			}
			invInfoStart2 = tPoint{x: x, y: y}
		case "invInfoStart3":
			if len(t) < 2 {
				return fmt.Errorf(" → Too few arguments on line %d", line)
			} else if len(t) > 2 {
				return fmt.Errorf(" → Too many arguments on line %d", line)
			}
			p := strings.Split(t[1], "; ")
			if len(p) < 2 {
				return fmt.Errorf(" → Too few numbers on line %d", line)
			} else if len(p) > 2 {
				return fmt.Errorf(" → Too many numbers on line %d", line)
			}
			fmt.Println(" → Invoice information 3 start at x =", p[0], "and y =", p[1])
			var err error
			var x, y float64
			x, err = strconv.ParseFloat(strings.TrimSpace(p[0]), 64)
			if err != nil {
				return fmt.Errorf(" → Error converting invoice information 3 start X number (\"%s\") into a decimal number at line %d", p[0], line)
			}
			y, err = strconv.ParseFloat(strings.TrimSpace(p[1]), 64)
			if err != nil {
				return fmt.Errorf(" → Error converting invoice information 3 start Y number (\"%s\") into a decimal number at line %d", p[1], line)
			}
			invInfoStart3 = tPoint{x: x, y: y}
		case "invInfoLine":
			if len(t) < 2 {
				return fmt.Errorf(" → Too few arguments on line %d", line)
			} else if len(t) > 2 {
				return fmt.Errorf(" → Too many arguments on line %d", line)
			}
			fmt.Println(" → Invoice text line height will now be", t[1])
			var err error
			invInfoLine, err = strconv.ParseFloat(t[1], 64)
			if err != nil {
				return fmt.Errorf(" → Error converting invoice text line height number (\"%s\") into a decimal number at line %d", t[1], line)
			}
		case "invInfoFontSz1":
			if len(t) < 2 {
				return fmt.Errorf(" → Too few arguments on line %d", line)
			} else if len(t) > 2 {
				return fmt.Errorf(" → Too many arguments on line %d", line)
			}
			fmt.Println(" → Invoice information 1 font size will now be", t[1])
			var err error
			invInfoFontSz1, err = strconv.ParseFloat(t[1], 64)
			if err != nil {
				return fmt.Errorf(" → Error converting invoice information 1 font size number (\"%s\") into a decimal number at line %d", t[1], line)
			}
		case "invInfoFontSz2":
			if len(t) < 2 {
				return fmt.Errorf(" → Too few arguments on line %d", line)
			} else if len(t) > 2 {
				return fmt.Errorf(" → Too many arguments on line %d", line)
			}
			fmt.Println(" → Invoice information 2 font size will now be", t[1])
			var err error
			invInfoFontSz2, err = strconv.ParseFloat(t[1], 64)
			if err != nil {
				return fmt.Errorf(" → Error converting invoice information 2 font size number (\"%s\") into a decimal number at line %d", t[1], line)
			}
		case "invInfoFontSz3":
			if len(t) < 2 {
				return fmt.Errorf(" → Too few arguments on line %d", line)
			} else if len(t) > 2 {
				return fmt.Errorf(" → Too many arguments on line %d", line)
			}
			fmt.Println(" → Invoice information 3 font size will now be", t[1])
			var err error
			invInfoFontSz3, err = strconv.ParseFloat(t[1], 64)
			if err != nil {
				return fmt.Errorf(" → Error converting invoice information 3 font size number (\"%s\") into a decimal number at line %d", t[1], line)
			}
		case "invInfoTitle":
			if len(t) < 2 {
				return fmt.Errorf(" → Too few arguments on line %d", line)
			} else if len(t) > 2 {
				return fmt.Errorf(" → Too many arguments on line %d", line)
			}
			fmt.Println(" → Invoice information title will be", t[1])
			invInfoTitle = t[1]
		case "invInfoText":
			if len(t) < 2 {
				return fmt.Errorf(" → Too few arguments on line %d", line)
			} else if len(t) > 2 {
				return fmt.Errorf(" → Too many arguments on line %d", line)
			}
			fmt.Println(" → Invoice information text will now be:", t[1])
			p := strings.Split(t[1], "; ")
			invInfoText = p
		case "cliInfoStart1":
			if len(t) < 2 {
				return fmt.Errorf(" → Too few arguments on line %d", line)
			} else if len(t) > 2 {
				return fmt.Errorf(" → Too many arguments on line %d", line)
			}
			p := strings.Split(t[1], "; ")
			if len(p) < 2 {
				return fmt.Errorf(" → Too few numbers on line %d", line)
			} else if len(p) > 2 {
				return fmt.Errorf(" → Too many numbers on line %d", line)
			}
			fmt.Println(" → Client information 1 start at x =", p[0], "and y =", p[1])
			var err error
			var x, y float64
			x, err = strconv.ParseFloat(strings.TrimSpace(p[0]), 64)
			if err != nil {
				return fmt.Errorf(" → Error converting client information 1 start X number (\"%s\") into a decimal number at line %d", p[0], line)
			}
			y, err = strconv.ParseFloat(strings.TrimSpace(p[1]), 64)
			if err != nil {
				return fmt.Errorf(" → Error converting client information 1 start Y number (\"%s\") into a decimal number at line %d", p[1], line)
			}
			cliInfoStart1 = tPoint{x: x, y: y}
		case "cliInfoStart2":
			if len(t) < 2 {
				return fmt.Errorf(" → Too few arguments on line %d", line)
			} else if len(t) > 2 {
				return fmt.Errorf(" → Too many arguments on line %d", line)
			}
			p := strings.Split(t[1], "; ")
			if len(p) < 2 {
				return fmt.Errorf(" → Too few numbers on line %d", line)
			} else if len(p) > 2 {
				return fmt.Errorf(" → Too many numbers on line %d", line)
			}
			fmt.Println(" → Client information 2 start at x =", p[0], "and y =", p[1])
			var err error
			var x, y float64
			x, err = strconv.ParseFloat(strings.TrimSpace(p[0]), 64)
			if err != nil {
				return fmt.Errorf(" → Error converting client information 2 start X number (\"%s\") into a decimal number at line %d", p[0], line)
			}
			y, err = strconv.ParseFloat(strings.TrimSpace(p[1]), 64)
			if err != nil {
				return fmt.Errorf(" → Error converting client information 2 start Y number (\"%s\") into a decimal number at line %d", p[1], line)
			}
			cliInfoStart2 = tPoint{x: x, y: y}
		case "cliInfoStart3":
			if len(t) < 2 {
				return fmt.Errorf(" → Too few arguments on line %d", line)
			} else if len(t) > 2 {
				return fmt.Errorf(" → Too many arguments on line %d", line)
			}
			p := strings.Split(t[1], "; ")
			if len(p) < 2 {
				return fmt.Errorf(" → Too few numbers on line %d", line)
			} else if len(p) > 2 {
				return fmt.Errorf(" → Too many numbers on line %d", line)
			}
			fmt.Println(" → Client information 3 start at x =", p[0], "and y =", p[1])
			var err error
			var x, y float64
			x, err = strconv.ParseFloat(strings.TrimSpace(p[0]), 64)
			if err != nil {
				return fmt.Errorf(" → Error converting client information 3 start X number (\"%s\") into a decimal number at line %d", p[0], line)
			}
			y, err = strconv.ParseFloat(strings.TrimSpace(p[1]), 64)
			if err != nil {
				return fmt.Errorf(" → Error converting client information 3 start Y number (\"%s\") into a decimal number at line %d", p[1], line)
			}
			cliInfoStart3 = tPoint{x: x, y: y}
		case "cliInfoLine":
			if len(t) < 2 {
				return fmt.Errorf(" → Too few arguments on line %d", line)
			} else if len(t) > 2 {
				return fmt.Errorf(" → Too many arguments on line %d", line)
			}
			fmt.Println(" → Client information text line height will now be", t[1])
			var err error
			cliInfoLine, err = strconv.ParseFloat(t[1], 64)
			if err != nil {
				return fmt.Errorf(" → Error converting client information text line height number (\"%s\") into a decimal number at line %d", t[1], line)
			}
		case "cliInfoFontSz1":
			if len(t) < 2 {
				return fmt.Errorf(" → Too few arguments on line %d", line)
			} else if len(t) > 2 {
				return fmt.Errorf(" → Too many arguments on line %d", line)
			}
			fmt.Println(" → Client information 1 font size will now be", t[1])
			var err error
			cliInfoFontSz1, err = strconv.ParseFloat(t[1], 64)
			if err != nil {
				return fmt.Errorf(" → Error converting client information 1 font size number (\"%s\") into a decimal number at line %d", t[1], line)
			}
		case "cliInfoFontSz2":
			if len(t) < 2 {
				return fmt.Errorf(" → Too few arguments on line %d", line)
			} else if len(t) > 2 {
				return fmt.Errorf(" → Too many arguments on line %d", line)
			}
			fmt.Println(" → Client information 2 font size will now be", t[1])
			var err error
			cliInfoFontSz2, err = strconv.ParseFloat(t[1], 64)
			if err != nil {
				return fmt.Errorf(" → Error converting client information 2 font size number (\"%s\") into a decimal number at line %d", t[1], line)
			}
		case "cliInfoFontSz3":
			if len(t) < 2 {
				return fmt.Errorf(" → Too few arguments on line %d", line)
			} else if len(t) > 2 {
				return fmt.Errorf(" → Too many arguments on line %d", line)
			}
			fmt.Println(" → Client information 3 font size will now be", t[1])
			var err error
			cliInfoFontSz3, err = strconv.ParseFloat(t[1], 64)
			if err != nil {
				return fmt.Errorf(" → Error converting client information 3 font size number (\"%s\") into a decimal number at line %d", t[1], line)
			}
		case "cliInfoName":
			if len(t) < 2 {
				return fmt.Errorf(" → Too few arguments on line %d", line)
			} else if len(t) > 2 {
				return fmt.Errorf(" → Too many arguments on line %d", line)
			}
			fmt.Println(" → Client name will be", t[1])
			cliInfoName = t[1]
		case "cliInfoData":
			if len(t) < 2 {
				return fmt.Errorf(" → Too few arguments on line %d", line)
			} else if len(t) > 2 {
				return fmt.Errorf(" → Too many arguments on line %d", line)
			}
			fmt.Println(" → Client information data will now be:", t[1])
			p := strings.Split(t[1], "; ")
			cliInfoData = p
		case "cliInfoCode":
			if len(t) < 2 {
				return fmt.Errorf(" → Too few arguments on line %d", line)
			} else if len(t) > 2 {
				return fmt.Errorf(" → Too many arguments on line %d", line)
			}
			fmt.Println(" → Client code will now be", t[1])
			cliInfoCode = t[1]
		case "listDetRect1":
			if len(t) < 2 {
				return fmt.Errorf(" → Too few arguments on line %d", line)
			} else if len(t) > 2 {
				return fmt.Errorf(" → Too many arguments on line %d", line)
			}
			p := strings.Split(t[1], "; ")
			if len(p) < 2 {
				return fmt.Errorf(" → Too few numbers on line %d", line)
			} else if len(p) > 2 {
				return fmt.Errorf(" → Too many numbers on line %d", line)
			}
			fmt.Println(" → List details rectangle will start at x =", p[0], "and y =", p[1])
			var err error
			var x, y float64
			x, err = strconv.ParseFloat(strings.TrimSpace(p[0]), 64)
			if err != nil {
				return fmt.Errorf(" → Error converting list details rectangle start X number (\"%s\") into a decimal number at line %d", p[0], line)
			}
			y, err = strconv.ParseFloat(strings.TrimSpace(p[1]), 64)
			if err != nil {
				return fmt.Errorf(" → Error converting list details rectangle start Y number (\"%s\") into a decimal number at line %d", p[1], line)
			}
			listDetRect1 = tPoint{x: x, y: y}
		case "listDetRect2":
			if len(t) < 2 {
				return fmt.Errorf(" → Too few arguments on line %d", line)
			} else if len(t) > 2 {
				return fmt.Errorf(" → Too many arguments on line %d", line)
			}
			p := strings.Split(t[1], "; ")
			if len(p) < 2 {
				return fmt.Errorf(" → Too few numbers on line %d", line)
			} else if len(p) > 2 {
				return fmt.Errorf(" → Too many numbers on line %d", line)
			}
			fmt.Println(" → List details rectangle will end at x =", p[0], "and y =", p[1])
			var err error
			var x, y float64
			x, err = strconv.ParseFloat(strings.TrimSpace(p[0]), 64)
			if err != nil {
				return fmt.Errorf(" → Error converting list details rectangle end X number (\"%s\") into a decimal number at line %d", p[0], line)
			}
			y, err = strconv.ParseFloat(strings.TrimSpace(p[1]), 64)
			if err != nil {
				return fmt.Errorf(" → Error converting list details rectangle end Y number (\"%s\") into a decimal number at line %d", p[1], line)
			}
			listDetRect2 = tPoint{x: x, y: y}
		case "listHFontSz":
			if len(t) < 2 {
				return fmt.Errorf(" → Too few arguments on line %d", line)
			} else if len(t) > 2 {
				return fmt.Errorf(" → Too many arguments on line %d", line)
			}
			fmt.Println(" → List header font size will now be", t[1])
			var err error
			listHFontSz, err = strconv.ParseFloat(t[1], 64)
			if err != nil {
				return fmt.Errorf(" → Error converting list header font size number (\"%s\") into a decimal number at line %d", t[1], line)
			}
		case "listHStartY":
			if len(t) < 2 {
				return fmt.Errorf(" → Too few arguments on line %d", line)
			} else if len(t) > 2 {
				return fmt.Errorf(" → Too many arguments on line %d", line)
			}
			fmt.Println(" → List header start Y will now be", t[1])
			var err error
			listHStartY, err = strconv.ParseFloat(t[1], 64)
			if err != nil {
				return fmt.Errorf(" → Error converting list header start Y number (\"%s\") into a decimal number at line %d", t[1], line)
			}
		case "listHStartXs":
			if len(t) < 2 {
				return fmt.Errorf(" → Too few arguments on line %d", line)
			} else if len(t) > 2 {
				return fmt.Errorf(" → Too many arguments on line %d", line)
			}
			fmt.Println(" → List header start Xs will now be", t[1])
			p := strings.Split(t[1], "; ")
			if len(p) < 4 {
				return fmt.Errorf(" → Too few values on line %d", line)
			} else if len(t) > 4 {
				return fmt.Errorf(" → Too many values on line %d", line)
			}
			var err error
			listHStartXs = []float64{}
			for _, v := range p {
				var tmp float64
				tmp, err = strconv.ParseFloat(v, 64)
				if err != nil {
					return fmt.Errorf(" → Error converting one of the list header start Xs number (\"%s\") into a decimal number at line %d", v, line)
				}
				listHStartXs = append(listHStartXs, tmp)
			}
		case "listHText":
			if len(t) < 2 {
				return fmt.Errorf(" → Too few arguments on line %d", line)
			} else if len(t) > 2 {
				return fmt.Errorf(" → Too many arguments on line %d", line)
			}
			fmt.Println(" → List header text will now be:", t[1])
			p := strings.Split(t[1], "; ")
			listHText = p
		case "listDFontSz":
			if len(t) < 2 {
				return fmt.Errorf(" → Too few arguments on line %d", line)
			} else if len(t) > 2 {
				return fmt.Errorf(" → Too many arguments on line %d", line)
			}
			fmt.Println(" → List details font size will now be", t[1])
			var err error
			listDFontSz, err = strconv.ParseFloat(t[1], 64)
			if err != nil {
				return fmt.Errorf(" → Error converting list details font size number (\"%s\") into a decimal number at line %d", t[1], line)
			}
		case "listDStartY":
			if len(t) < 2 {
				return fmt.Errorf(" → Too few arguments on line %d", line)
			} else if len(t) > 2 {
				return fmt.Errorf(" → Too many arguments on line %d", line)
			}
			fmt.Println(" → List details start Y will now be", t[1])
			var err error
			listDStartY, err = strconv.ParseFloat(t[1], 64)
			if err != nil {
				return fmt.Errorf(" → Error converting list details start Y number (\"%s\") into a decimal number at line %d", t[1], line)
			}
		case "listDStartXs":
			if len(t) < 2 {
				return fmt.Errorf(" → Too few arguments on line %d", line)
			} else if len(t) > 2 {
				return fmt.Errorf(" → Too many arguments on line %d", line)
			}
			fmt.Println(" → List details start Xs will now be", t[1])
			p := strings.Split(t[1], "; ")
			if len(p) < 4 {
				return fmt.Errorf(" → Too few values on line %d", line)
			} else if len(t) > 4 {
				return fmt.Errorf(" → Too many values on line %d", line)
			}
			var err error
			listDStartXs = []float64{}
			for _, v := range p {
				var tmp float64
				tmp, err = strconv.ParseFloat(v, 64)
				if err != nil {
					return fmt.Errorf(" → Error converting one of the list details start Xs number (\"%s\") into a decimal number at line %d", v, line)
				}
				listDStartXs = append(listDStartXs, tmp)
			}
		case "listDetLine":
			if len(t) < 2 {
				return fmt.Errorf(" → Too few arguments on line %d", line)
			} else if len(t) > 2 {
				return fmt.Errorf(" → Too many arguments on line %d", line)
			}
			fmt.Println(" → List details line height will now be", t[1])
			var err error
			listDetLine, err = strconv.ParseFloat(t[1], 64)
			if err != nil {
				return fmt.Errorf(" → Error converting list details line height number (\"%s\") into a decimal number at line %d", t[1], line)
			}
		case "listDetSepX1":
			if len(t) < 2 {
				return fmt.Errorf(" → Too few arguments on line %d", line)
			} else if len(t) > 2 {
				return fmt.Errorf(" → Too many arguments on line %d", line)
			}
			fmt.Println(" → List details separator start X will now be", t[1])
			var err error
			listDetSepX1, err = strconv.ParseFloat(t[1], 64)
			if err != nil {
				return fmt.Errorf(" → Error converting list details separator start X number (\"%s\") into a decimal number at line %d", t[1], line)
			}
		case "listDetSepX2":
			if len(t) < 2 {
				return fmt.Errorf(" → Too few arguments on line %d", line)
			} else if len(t) > 2 {
				return fmt.Errorf(" → Too many arguments on line %d", line)
			}
			fmt.Println(" → List details separator end X will now be", t[1])
			var err error
			listDetSepX2, err = strconv.ParseFloat(t[1], 64)
			if err != nil {
				return fmt.Errorf(" → Error converting list details separator end X number (\"%s\") into a decimal number at line %d", t[1], line)
			}
		case "taxMsgStartX":
			if len(t) < 2 {
				return fmt.Errorf(" → Too few arguments on line %d", line)
			} else if len(t) > 2 {
				return fmt.Errorf(" → Too many arguments on line %d", line)
			}
			fmt.Println(" → Tax message start X will now be", t[1])
			var err error
			taxMsgStartX, err = strconv.ParseFloat(t[1], 64)
			if err != nil {
				return fmt.Errorf(" → Error converting tax message start X number (\"%s\") into a decimal number at line %d", t[1], line)
			}
		case "taxText1":
			if len(t) < 2 {
				return fmt.Errorf(" → Too few arguments on line %d", line)
			} else if len(t) > 2 {
				return fmt.Errorf(" → Too many arguments on line %d", line)
			}
			fmt.Println(" → Tax text 1 will now be", t[1])
			taxText1 = t[1]
		case "taxText2":
			if len(t) < 2 {
				return fmt.Errorf(" → Too few arguments on line %d", line)
			} else if len(t) > 2 {
				return fmt.Errorf(" → Too many arguments on line %d", line)
			}
			fmt.Println(" → Tax text 2 will now be", t[1])
			taxText2 = t[1]
		case "subTotalStartX":
			if len(t) < 2 {
				return fmt.Errorf(" → Too few arguments on line %d", line)
			} else if len(t) > 2 {
				return fmt.Errorf(" → Too many arguments on line %d", line)
			}
			fmt.Println(" → Subtotal start X will now be", t[1])
			var err error
			subTotalStartX, err = strconv.ParseFloat(t[1], 64)
			if err != nil {
				return fmt.Errorf(" → Error converting Subtotal start X number (\"%s\") into a decimal number at line %d", t[1], line)
			}
		case "subTotalText":
			if len(t) < 2 {
				return fmt.Errorf(" → Too few arguments on line %d", line)
			} else if len(t) > 2 {
				return fmt.Errorf(" → Too many arguments on line %d", line)
			}
			fmt.Println(" → Subtotal text will now be", t[1])
			subTotalText = t[1]
		case "totalText":
			if len(t) < 2 {
				return fmt.Errorf(" → Too few arguments on line %d", line)
			} else if len(t) > 2 {
				return fmt.Errorf(" → Too many arguments on line %d", line)
			}
			fmt.Println(" → Total text will now be", t[1])
			totalText = t[1]
		case "totAmoStartX":
			if len(t) < 2 {
				return fmt.Errorf(" → Too few arguments on line %d", line)
			} else if len(t) > 2 {
				return fmt.Errorf(" → Too many arguments on line %d", line)
			}
			fmt.Println(" → Total amount start X will now be", t[1])
			var err error
			totAmoStartX, err = strconv.ParseFloat(t[1], 64)
			if err != nil {
				return fmt.Errorf(" → Error converting Total amount start X number (\"%s\") into a decimal number at line %d", t[1], line)
			}
		default:
			if t[0][0] == '#' {
				fmt.Println(" → Comment at line", line)
			} else {
				return fmt.Errorf(" → ⚠️ Unrecognized command \"%s\" on line %d", t[0], line)
			}
		}
	}

	return nil
}

func loadStatic(pdf *gopdf.GoPdf) error {

	var err error

	// Create document
	pdf.Start(gopdf.Config{PageSize: gopdf.Rect{W: documentWidth, H: documentHeight}})

	// Adding first page
	pdf.AddPage()

	// Load fonts
	err = pdf.AddTTFFont("Regular", fontRegular)
	if err != nil {
		return err
	}

	err = pdf.AddTTFFont("Bold", fontBold)
	if err != nil {
		return err
	}

	// Load background image
	pdf.Image(backgroundImage, 0.0, 0.0, nil)

	// My info
	err = pdf.Rectangle(myInfoRect1.x, myInfoRect1.y, myInfoRect2.x, myInfoRect2.y, "D", 0, 0)
	if err != nil {
		return err
	}

	err = pdf.SetFont("Regular", "", myInfoFontSize)
	if err != nil {
		return err
	}

	for i, s := range myInfoData {
		pdf.SetXY(myInfoStart.x, myInfoStart.y + (float64(i) * myInfoLine))
		pdf.Text(s)
	}

	// Invoice info
	pdf.SetFillColor(0, 0, 0)
	err = pdf.Rectangle(invInfoRect1.x, invInfoRect1.y, invInfoRect2.x, invInfoRect2.y, "F", 0, 0)
	if err != nil {
		return err
	}

	err = pdf.SetFont("Bold", "", invInfoFontSz1)
	if err != nil {
		return err
	}

	pdf.SetTextColor(255, 255, 255)
	pdf.SetXY(invInfoStart1.x, invInfoStart1.y)
	pdf.Text(invInfoTitle)
	pdf.SetTextColor(0, 0, 0)

	err = pdf.SetFont("Bold", "", invInfoFontSz2)
	if err != nil {
		return err
	}

	for i, s := range invInfoText {
		pdf.SetXY(invInfoStart2.x, invInfoStart2.y + (float64(i) * invInfoLine))
		pdf.Text(s)
	}

	err = pdf.SetFont("Regular", "", invInfoFontSz3)
	if err != nil {
		return err
	}

	pdf.SetXY(invInfoStart3.x, invInfoStart3.y)
	pdf.Text(fmt.Sprint(page))
	page++

	pdf.SetXY(invInfoStart3.x, invInfoStart3.y + (1 * invInfoLine))
	pdf.Text(invNum)

	pdf.SetXY(invInfoStart3.x, invInfoStart3.y + (2 * invInfoLine))
	pdf.Text(date)

	pdf.SetXY(invInfoStart3.x, invInfoStart3.y + (3 * invInfoLine))
	pdf.Text(fmt.Sprint(cliNum))

	// Client info
	err = pdf.SetFont("Bold", "", cliInfoFontSz1)
	if err != nil {
		return err
	}

	pdf.SetXY(cliInfoStart1.x, cliInfoStart1.y)
	pdf.Text(cliInfoName)

	err = pdf.SetFont("Regular", "", cliInfoFontSz2)
	if err != nil {
		return err
	}

	for i, s := range cliInfoData {
		pdf.SetXY(cliInfoStart2.x, cliInfoStart2.y + (float64(i) * cliInfoLine))
		pdf.Text(s)
	}

	pdf.SetTextColor(80, 80, 80)
	err = pdf.SetFont("Bold", "", cliInfoFontSz3)
	if err != nil {
		return err
	}
	pdf.SetXY(cliInfoStart3.x, cliInfoStart3.y)
	pdf.Text(cliInfoCode)

	pdf.SetTextColor(0, 0, 0)

	// List

	err = pdf.Rectangle(listDetRect1.x, listDetRect1.y, listDetRect2.x, listDetRect2.y, "D", 0, 0)

	err = pdf.SetFont("Bold", "", listHFontSz)
	if err != nil {
		return err
	}

	for i, x := range listHStartXs {
		pdf.SetXY(x, listHStartY)
		if i < len(listHText) {
			pdf.Text(listHText[i])
		}
	}

	return nil
}

func loadListDetails(pdf *gopdf.GoPdf) error {

	var err error
	var displacement float64

	err = pdf.SetFont("Regular", "", listDFontSz)
	if err != nil {
		return err
	}

	startY := listDStartY
	for _, l := range listDet {
		for _, c := range l.concept {
			startY += listDetLine
			pdf.SetXY(listDStartXs[0], startY)
			pdf.Text(c)
		}
		
		displacement, err = pdf.MeasureTextWidth(l.amount)
		if err != nil {
			return err
		}
		pdf.SetX(listDStartXs[1] - displacement)
		pdf.Text(l.amount)
		displacement, err = pdf.MeasureTextWidth(l.price)
		if err != nil {
			return err
		}
		pdf.SetX(listDStartXs[2] - displacement)
		pdf.Text(l.price)
		money := convertMoney(l.total, coin, coinFront)
		displacement, err = pdf.MeasureTextWidth(money)
		if err != nil {
			return err
		}
		pdf.SetX(listDStartXs[3] - displacement)
		pdf.Text(money)
		totalNoTax += l.total

		pdf.SetLineType("dotted")
		pdf.SetStrokeColor(150, 150, 150)
		pdf.Line(listDetSepX1, startY + (listDetLine / 2), listDetSepX2, startY + (listDetLine / 2))
		pdf.SetStrokeColor(0, 0, 0)
		startY += listDetLine
	}

	// Total stuf

	err = pdf.SetFont("Bold", "", listDFontSz)
	if err != nil {
		return err
	}

	posY := pdf.GetY()

	pdf.SetXY(taxMsgStartX, posY + (listDetLine * 4))
	if taxes == 0.0 {
		pdf.Text(taxText1)
	} else {
		replace := fmt.Sprintf("%g", taxes)
		text := strings.ReplaceAll(taxText2, "[NUM]", replace)
		pdf.Text(text)
	}

	// Total base (with no taxes)

	pdf.SetXY(subTotalStartX, posY + (listDetLine * 3))
	pdf.Text(subTotalText)
	
	subTotalAmount := convertMoney(totalNoTax, coin, coinFront)

	displacement, err = pdf.MeasureTextWidth(subTotalAmount)
	if err != nil {
		return err
	}

	pdf.SetX(totAmoStartX - displacement)
	pdf.Text(subTotalAmount)

	// Total with taxes

	pdf.SetXY(subTotalStartX, posY + (listDetLine * 5))
	pdf.Text(totalText)

	totalAmount := convertMoney(int(float64(totalNoTax) * (1.0 + taxes/100.0)), coin, coinFront)

	displacement, err = pdf.MeasureTextWidth(totalAmount)
	if err != nil {
		return err
	}

	pdf.SetX(totAmoStartX - displacement)
	pdf.Text(totalAmount)

	// Line under total base

	pdf.SetLineWidth(1)
	pdf.SetLineType("")
	pdf.Line(subTotalStartX, posY + (listDetLine * 3) + 2, listDetRect2.x, posY + (listDetLine * 3) + 2)

	// Line under total with taxes

	pdf.Line(subTotalStartX, posY + (listDetLine * 5) + 2, listDetRect2.x, posY + (listDetLine * 5) + 2)
	pdf.Line(subTotalStartX, posY + (listDetLine * 5) + 4, listDetRect2.x, posY + (listDetLine * 5) + 4)

	return nil
}

func convertMoney(a int, c string, cf bool) string {
	b := fmt.Sprintf("%.2f", float64(a)/100.0)

	d := "," + b[len(b)-2:]
	b = b[:len(b)-3]

	count := 0
	for i := len(b) - 1; i >= 0; i-- {
		if i != len(b)-1 && count%3 == 0 {
			d = "." + d
		}
		d = string(b[i]) + d
		count++
	}

	if cf {
		return c + " " + d
	}
	return d + " " + c
}

func main() {

	var err error
	var fTsv *os.File

	if len(os.Args) == 1 {
		fmt.Println("No input files provided, using default settings.")
	} else { // One or more input files
		if len(os.Args) > 2 {
			fmt.Println("Only the first file will be taken into account.")
		}
		fmt.Println("Opening", os.Args[1])
		fTsv, err = os.Open(os.Args[1])
		if err != nil {
			log.Print(err.Error())
			os.Exit(-1)
		}
		defer fTsv.Close()
	}

	pdf := gopdf.GoPdf{}

	if fTsv != nil {
		err = processFTsv(fTsv)
		if err != nil {
			log.Print(err.Error())
			os.Exit(-1)
		}
	}

	err = loadStatic(&pdf)
	if err != nil {
		log.Print(err.Error())
		os.Exit(-1)
	}

	err = loadListDetails(&pdf)
	if err != nil {
		log.Print(err.Error())
		os.Exit(-1)
	}

	// Save and exit
	fmt.Println()
	fmt.Println("Remember to execute these commands:")
	fmt.Printf("exiftool -Title=\"%s\" -Author=\"%s\" -Creator=\"%s\" -Producer=\"%s\" -overwrite_original %s\n",
	 title, author, creator, producer, outputFileName)
	fmt.Println("pdfinfo " + outputFileName)
	fmt.Println()
	pdf.WritePdf(outputFileName)

}
