package document

import (
	"fmt"
	"math"
	"time"

	"github.com/signintech/gopdf"
	"github.com/ubavic/bas-celik/localization"
)

type IdPdfWriter struct {
	pdf            *gopdf.GoPdf
	leftMargin     float64
	rightMargin    float64
	textLeftMargin float64
	doc            *IdDocument
}

func (idw *IdPdfWriter) line(width float64) {
	if width > 0 {
		idw.pdf.SetLineWidth(width)
	}

	y := idw.pdf.GetY()
	idw.pdf.Line(idw.leftMargin, y, idw.rightMargin, y)
}

func (idw *IdPdfWriter) moveY(y float64) {
	idw.pdf.SetXY(idw.pdf.GetX(), idw.pdf.GetY()+y)
}

func (idw *IdPdfWriter) cell(s string) {
	err := idw.pdf.Cell(nil, s)
	if err != nil {
		panic(fmt.Errorf("putting text: %w", err))
	}
}

func (idw *IdPdfWriter) putData(label, data string) {
	y := idw.pdf.GetY()

	idw.pdf.SetX(idw.textLeftMargin)
	texts, err := idw.pdf.SplitTextWithWordWrap(label, 120)
	if err != nil && err != gopdf.ErrEmptyString {
		panic(err)
	}

	for i, text := range texts {
		idw.cell(text)
		if i < len(texts)-1 {
			idw.pdf.SetXY(idw.textLeftMargin, idw.pdf.GetY()+12)
		}
	}

	y1 := idw.pdf.GetY()

	idw.pdf.SetXY(idw.textLeftMargin+128, y)
	texts, err = idw.pdf.SplitTextWithWordWrap(data, 350)
	if err != nil && err != gopdf.ErrEmptyString {
		panic(err)
	}

	for i, text := range texts {
		idw.cell(text)
		if i < len(texts)-1 {
			idw.pdf.SetXY(idw.textLeftMargin+128, idw.pdf.GetY()+12)
		}
	}

	y2 := idw.pdf.GetY()

	idw.pdf.SetXY(idw.textLeftMargin, math.Max(y1, y2)+24.67)
}

func (ipw *IdPdfWriter) printRegularId() {
	ipw.pdf.SetLineType("solid")
	ipw.pdf.SetY(59.041)
	ipw.line(0.83)

	ipw.pdf.SetXY(ipw.textLeftMargin+1.0, 68.5)

	err := ipw.pdf.SetCharSpacing(-0.2)
	if err != nil {
		panic(err)
	}
	ipw.cell("ČITAČ ELEKTRONSKE LIČNE KARTE: ŠTAMPA PODATAKA")

	err = ipw.pdf.SetCharSpacing(-0.1)
	if err != nil {
		panic(err)
	}

	ipw.pdf.SetY(88)

	ipw.line(0)

	imageY := 102.8
	imageHeight := 159.0

	err = ipw.pdf.ImageFrom(ipw.doc.Portrait, ipw.leftMargin, imageY, &gopdf.Rect{W: 119.9, H: imageHeight})
	if err != nil {
		panic(err)
	}

	ipw.pdf.SetLineWidth(0.48)
	ipw.pdf.SetFillColor(255, 255, 255)
	err = ipw.pdf.Rectangle(ipw.leftMargin, imageY, 179, imageY+imageHeight, "D", 0, 0)
	if err != nil {
		panic(err)
	}

	ipw.pdf.SetFillColor(0, 0, 0)

	ipw.pdf.SetY(276)

	ipw.line(1.08)
	ipw.moveY(8)
	ipw.pdf.SetX(ipw.textLeftMargin)
	err = ipw.pdf.SetFontSize(11.1)
	if err != nil {
		panic(err)
	}

	ipw.cell("Podaci o građaninu")

	ipw.moveY(16)
	ipw.line(0)
	ipw.moveY(9)

	ipw.putData("Prezime:", ipw.doc.Surname)
	ipw.putData("Ime:", ipw.doc.GivenName)
	ipw.putData("Ime jednog roditelja:", ipw.doc.ParentGivenName)
	ipw.putData("Datum rođenja:", ipw.doc.DateOfBirth)
	ipw.putData("Mesto rođenja,\nopština i država:", ipw.doc.GetFullPlaceOfBirth())
	ipw.putData("Prebivalište:", ipw.doc.GetFullAddress(true))
	ipw.putData("Datum promene adrese:", ipw.doc.AddressDate)
	ipw.putData("JMBG:", ipw.doc.PersonalNumber)
	ipw.putData("Pol:", ipw.doc.Sex)

	ipw.moveY(-8.67)
	ipw.line(0)
	ipw.moveY(9)
	ipw.cell("Podaci o dokumentu")
	ipw.moveY(16)

	ipw.line(0)
	ipw.moveY(9)
	ipw.putData("Dokument izdaje:", ipw.doc.IssuingAuthority)
	ipw.putData("Broj dokumenta:", ipw.doc.DocRegNo)
	ipw.putData("Datum izdavanja:", ipw.doc.IssuingDate)
	ipw.putData("Važi do:", ipw.doc.ExpiryDate)

	ipw.moveY(-8.67)
	ipw.line(0)
	ipw.moveY(3)
	ipw.line(0)
	ipw.moveY(9)

	ipw.cell("Datum štampe: " + time.Now().Format("02.01.2006."))

	ipw.moveY(19)

	if ipw.pdf.GetY() < 700 {
		ipw.pdf.SetY(730.6)
	}

	ipw.line(0.83)

	err = ipw.pdf.SetFontSize(9)
	if err != nil {
		panic(err)
	}

	ipw.moveY(10)
	ipw.pdf.SetX(ipw.leftMargin)

	ipw.cell("1. U čipu lične karte, podaci o imenu i prezimenu imaoca lične karte ispisani su na nacionalnom pismu onako kako su")
	ipw.pdf.SetX(ipw.leftMargin)
	ipw.moveY(9.7)
	ipw.cell("ispisani na samom obrascu lične karte, dok su ostali podaci ispisani latiničkim pismom.")
	ipw.pdf.SetX(ipw.leftMargin)
	ipw.moveY(9.7)
	ipw.cell("2. Ako se ime lica sastoji od dve reči čija je ukupna dužina između 20 i 30 karaktera ili prezimena od dve reči čija je")
	ipw.pdf.SetX(ipw.leftMargin)
	ipw.moveY(9.7)
	ipw.cell("ukupna dužina između 30 i 36 karaktera, u čipu lične karte izdate pre 18.08.2014. godine, druga reč u imenu ili prezimenu")
	ipw.pdf.SetX(ipw.leftMargin)
	ipw.moveY(9.7)
	ipw.cell("skraćuje se na prva dva karaktera")

	ipw.moveY(15.7)
	ipw.line(0)
}

func (ipw *IdPdfWriter) printForeignerId() {
	ipw.pdf.SetLineType("solid")
	ipw.pdf.SetY(59.041)
	ipw.line(0.83)

	ipw.pdf.SetXY(ipw.textLeftMargin+1.0, 64.95)

	err := ipw.pdf.SetCharSpacing(-0.2)
	if err != nil {
		panic(err)
	}
	ipw.cell("ČITAČ ELEKTRONSKE LIČNE KARTE: ŠTAMPA PODATAKA")

	err = ipw.pdf.SetCharSpacing(-0.1)
	if err != nil {
		panic(err)
	}

	ipw.pdf.SetY(79.8)

	ipw.line(0)

	imageY := 86.0
	imageHeight := 159.0

	err = ipw.pdf.ImageFrom(ipw.doc.Portrait, ipw.leftMargin, imageY, &gopdf.Rect{W: 119.9, H: imageHeight})
	if err != nil {
		panic(err)
	}

	ipw.pdf.SetLineWidth(0.48)
	ipw.pdf.SetFillColor(255, 255, 255)
	err = ipw.pdf.Rectangle(ipw.leftMargin, imageY, 179, imageY+imageHeight, "D", 0, 0)
	if err != nil {
		panic(err)
	}

	ipw.pdf.SetFillColor(0, 0, 0)

	ipw.pdf.SetY(250)

	ipw.line(1.08)
	ipw.moveY(8)
	ipw.pdf.SetX(ipw.textLeftMargin)
	err = ipw.pdf.SetFontSize(11.1)
	if err != nil {
		panic(err)
	}

	ipw.cell("Podaci o strancu")

	ipw.moveY(16)
	ipw.line(0)
	ipw.moveY(9)

	ipw.putData("Prezime:", ipw.doc.Surname)
	ipw.putData("Ime:", ipw.doc.GivenName)
	ipw.putData("Državljanstvo:", ipw.doc.NationalityFull)
	ipw.putData("Datum rođenja:", ipw.doc.DateOfBirth)
	ipw.putData("Osnov boravka:", ipw.doc.PurposeOfStay)
	ipw.putData("Prebivalište:", localization.JoinWithComma(ipw.doc.State, ipw.doc.GetFullAddress(true)))
	ipw.putData("Datum promene adrese:", ipw.doc.AddressDate)
	ipw.putData("Evidencijski broj\nstranca:", ipw.doc.PersonalNumber)
	ipw.putData("Pol:", ipw.doc.Sex)

	ipw.moveY(-8.67)
	ipw.line(0)
	ipw.moveY(9)
	ipw.cell("Podaci o dokumentu")
	ipw.moveY(16)

	ipw.line(0)
	ipw.moveY(9)
	ipw.putData("Dokument izdaje:", ipw.doc.IssuingAuthority)
	ipw.putData("Broj dokumenta:", ipw.doc.DocRegNo)
	ipw.putData("Datum izdavanja:", ipw.doc.IssuingDate)
	ipw.putData("Važi do:", ipw.doc.ExpiryDate)

	ipw.moveY(-8.67)
	ipw.line(0)
	ipw.moveY(3)
	ipw.line(0)
	ipw.moveY(9)

	ipw.cell("Datum štampe: " + time.Now().Format("02.01.2006."))

	ipw.moveY(19)

	ipw.line(0.83)

	err = ipw.pdf.SetFontSize(9)
	if err != nil {
		panic(err)
	}

	ipw.moveY(4)

	ipw.pdf.SetX(ipw.leftMargin)

	ipw.cell("1. U čipu lične karte za strance, podaci o imenu i prezimenu stranca ispisani su onako kako su ispisani na samom")
	ipw.pdf.SetX(ipw.leftMargin)
	ipw.moveY(9.7)
	ipw.cell("obrascu lične karte za stranca latiničnim pismom.")
	ipw.pdf.SetX(ipw.leftMargin)
	ipw.moveY(9.7)
	ipw.cell("2. Ako se ime ili prezime stranca sastoji od dve ili više reči čija dužina prelazi 30 karaktera za ime, odnosno 36")
	ipw.pdf.SetX(ipw.leftMargin)
	ipw.moveY(9.7)
	ipw.cell("karaktera za prezime, u čip se upisuje puno ime i prezime stranca, a na obrascu lične karte za stranca se upisuje do")
	ipw.pdf.SetX(ipw.leftMargin)
	ipw.moveY(9.7)
	ipw.cell("30 karaktera za ime, odnosno 36 karaktera za prezime.")

	ipw.moveY(9.7)

	ipw.line(0)
}

func (ipw *IdPdfWriter) printResidencePermit() {
	ipw.pdf.SetLineType("solid")
	ipw.pdf.SetY(59.041)
	ipw.line(0.83)

	ipw.pdf.SetXY(ipw.textLeftMargin+1.0, 64.95)

	err := ipw.pdf.SetCharSpacing(-0.2)
	if err != nil {
		panic(err)
	}
	ipw.cell("ČITAČ ELEKTRONSKE LIČNE KARTE: ŠTAMPA PODATAKA")

	err = ipw.pdf.SetCharSpacing(-0.1)
	if err != nil {
		panic(err)
	}

	ipw.pdf.SetY(79.8)

	ipw.line(0)

	imageY := 86.0
	imageHeight := 159.0

	err = ipw.pdf.ImageFrom(ipw.doc.Portrait, ipw.leftMargin, imageY, &gopdf.Rect{W: 119.9, H: imageHeight})
	if err != nil {
		panic(err)
	}

	ipw.pdf.SetLineWidth(0.48)
	ipw.pdf.SetFillColor(255, 255, 255)
	err = ipw.pdf.Rectangle(ipw.leftMargin, imageY, 179, imageY+imageHeight, "D", 0, 0)
	if err != nil {
		panic(err)
	}

	ipw.pdf.SetFillColor(0, 0, 0)

	ipw.pdf.SetY(250)

	ipw.line(1.08)
	ipw.moveY(8)
	ipw.pdf.SetX(ipw.textLeftMargin)
	err = ipw.pdf.SetFontSize(11.1)
	if err != nil {
		panic(err)
	}

	ipw.cell("Podaci o strancu")

	ipw.moveY(16)
	ipw.line(0)
	ipw.moveY(9)

	ipw.putData("Prezime:", ipw.doc.Surname)
	ipw.putData("Ime:", ipw.doc.GivenName)
	ipw.putData("Državljanstvo:", ipw.doc.NationalityFull)
	ipw.putData("Datum rođenja:", ipw.doc.DateOfBirth)
	ipw.putData("Mesto rođenja,\nopština i država:", ipw.doc.GetFullPlaceOfBirth())
	ipw.putData("Prebivalište:", ipw.doc.GetFullAddress(true))
	ipw.putData("Datum promene adrese:", ipw.doc.AddressDate)
	ipw.putData("Evidencijski broj\nstranca:", ipw.doc.PersonalNumber)
	ipw.putData("Pol:", ipw.doc.Sex)
	ipw.putData("Osnov boravka:", ipw.doc.PurposeOfStay)
	ipw.putData("Napomena:", ipw.doc.ENote)

	ipw.moveY(-8.67)
	ipw.line(0)
	ipw.moveY(9)
	ipw.cell("Podaci o dokumentu")
	ipw.moveY(16)

	ipw.line(0)
	ipw.moveY(9)
	ipw.putData("Naziv dokumenta:", ipw.doc.DocumentName)
	ipw.putData("Dokument izdaje:", ipw.doc.IssuingAuthority)
	ipw.putData("Broj dokumenta:", ipw.doc.DocRegNo)
	ipw.putData("Datum izdavanja:", ipw.doc.IssuingDate)
	ipw.putData("Važi do:", ipw.doc.ExpiryDate)

	ipw.moveY(-8.67)
	ipw.line(0)
	ipw.moveY(3)
	ipw.line(0)
	ipw.moveY(9)

	ipw.cell("Datum štampe: " + time.Now().Format("02.01.2006."))

	ipw.moveY(19)

	ipw.line(0.83)

	err = ipw.pdf.SetFontSize(9)
	if err != nil {
		panic(err)
	}

	ipw.moveY(4)

	ipw.pdf.SetX(ipw.leftMargin)

	ipw.cell("1. U čipu dozvole za privremeni boravak i rad, podaci o imenu i prezimenu imaoca dozvole ispisani su onako")
	ipw.pdf.SetX(ipw.leftMargin)
	ipw.moveY(9.7)
	ipw.cell("kako su ispisani na samom obrascu dozvole za privremeni boravak latiničnim pismom.")
	ipw.pdf.SetX(ipw.leftMargin)
	ipw.moveY(9.7)
	ipw.cell("2. Ako se ime ili prezime stranca sastoji od dve ili više reči čija dužina prelazi 30 karaktera za ime,")
	ipw.pdf.SetX(ipw.leftMargin)
	ipw.moveY(9.7)
	ipw.cell("odnosno 36 karaktera za prezime u čip se upisuje puno ime stranca, a na obrascu dozvole za privremeni boravak")
	ipw.pdf.SetX(ipw.leftMargin)
	ipw.moveY(9.7)
	ipw.cell("se upisuje do 30 karaktera za ime, odnosno 36 karaktera za prezime.")

	ipw.moveY(9.7)

	ipw.line(0)
}
