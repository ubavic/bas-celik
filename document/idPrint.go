package document

import (
	"fmt"
	"math"
	"time"

	"github.com/signintech/gopdf"
	"github.com/ubavic/bas-celik/v2/localization"
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

func (idw *IdPdfWriter) cell(s string, transliterate bool) {
	if transliterate && !idw.doc.pdfCyrillicLabels {
		s = localization.CyrillicToLatin(s)
	}

	err := idw.pdf.Cell(nil, s)
	if err != nil {
		panic(fmt.Errorf("putting text: %w", err))
	}
}

func (idw *IdPdfWriter) putData(label, data string) {
	y := idw.pdf.GetY()

	if !idw.doc.pdfCyrillicLabels {
		label = localization.CyrillicToLatin(label)
	}

	idw.pdf.SetX(idw.textLeftMargin)
	texts, err := idw.pdf.SplitTextWithWordWrap(label, 120)
	if err != nil && err != gopdf.ErrEmptyString {
		panic(err)
	}

	for i, text := range texts {
		idw.cell(text, false)
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
		idw.cell(text, false)
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
	ipw.cell("ЧИТАЧ ЕЛЕКТРОНСКЕ ЛИЧНЕ КАРТЕ: ШТАМПА ПОДАТАКА", true)

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

	ipw.cell("Подаци о грађанину", true)

	ipw.moveY(16)
	ipw.line(0)
	ipw.moveY(9)

	ipw.putData("Презиме:", ipw.doc.Surname)
	ipw.putData("Име:", ipw.doc.GivenName)
	ipw.putData("Име једног родитеља:", ipw.doc.ParentGivenName)
	ipw.putData("Датум рођења:", ipw.doc.DateOfBirth)
	ipw.putData("Место рођења\nопштина и држава:", ipw.doc.GetFullPlaceOfBirth())
	addressLabel := "Пребивалиште\nи адреса стана:"
	if ipw.doc.AddressLabel == "prebivalište" {
		addressLabel = "Пребивалиште:"
	}
	ipw.putData(addressLabel, ipw.doc.GetFullAddress(true))
	ipw.putData("Датум промене адресе:", ipw.doc.AddressDate)
	ipw.putData("ЈМБГ:", ipw.doc.PersonalNumber)
	ipw.putData("Пол:", ipw.doc.Sex)

	ipw.moveY(-8.67)
	ipw.line(0)
	ipw.moveY(9)
	ipw.cell("Подаци о документу", true)
	ipw.moveY(16)

	ipw.line(0)
	ipw.moveY(9)
	ipw.putData("Документ издаје:", ipw.doc.IssuingAuthority)
	ipw.putData("Број документа:", ipw.doc.DocRegNo)
	ipw.putData("Датум издавања:", ipw.doc.IssuingDate)
	ipw.putData("Важи до:", ipw.doc.ExpiryDate)

	ipw.moveY(-8.67)
	ipw.line(0)
	ipw.moveY(3)
	ipw.line(0)
	ipw.moveY(9)

	ipw.cell("Датум штампе: "+time.Now().Format("02.01.2006."), true)

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

	if ipw.doc.pdfCyrillicLabels {
		ipw.cell("1. У чипу личне карте, подаци о имену и презимену имаоца личне карте исписани су на националном писму онако", false)
		ipw.pdf.SetX(ipw.leftMargin)
		ipw.moveY(9.7)
		ipw.cell("како су исписани на самом обрасцу личне карте, док су остали подаци исписани латиничким писмом.", false)
		ipw.pdf.SetX(ipw.leftMargin)
		ipw.moveY(9.7)
		ipw.cell("2. Ако се име лица састоји од две речи чија је укупна дужина између 20 и 30 карактера или презимена од две речи", false)
		ipw.pdf.SetX(ipw.leftMargin)
		ipw.moveY(9.7)
		ipw.cell("чија је укупна дужина између 30 и 36 карактера, у чипу личне карте издате пре 18.08.2014. године, друга реч у", false)
		ipw.pdf.SetX(ipw.leftMargin)
		ipw.moveY(9.7)
		ipw.cell("имену или презимену скраћује се на прва два карактера", false)
	} else {
		ipw.cell("1. U čipu lične karte, podaci o imenu i prezimenu imaoca lične karte ispisani su na nacionalnom pismu onako kako su", false)
		ipw.pdf.SetX(ipw.leftMargin)
		ipw.moveY(9.7)
		ipw.cell("ispisani na samom obrascu lične karte, dok su ostali podaci ispisani latiničkim pismom.", false)
		ipw.pdf.SetX(ipw.leftMargin)
		ipw.moveY(9.7)
		ipw.cell("2. Ako se ime lica sastoji od dve reči čija je ukupna dužina između 20 i 30 karaktera ili prezimena od dve reči čija je", false)
		ipw.pdf.SetX(ipw.leftMargin)
		ipw.moveY(9.7)
		ipw.cell("ukupna dužina između 30 i 36 karaktera, u čipu lične karte izdate pre 18.08.2014. godine, druga reč u imenu ili prezimenu", false)
		ipw.pdf.SetX(ipw.leftMargin)
		ipw.moveY(9.7)
		ipw.cell("skraćuje se na prva dva karaktera", false)
	}

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
	ipw.cell("ЧИТАЧ ЕЛЕКТРОНСКЕ ЛИЧНЕ КАРТЕ: ШТАМПА ПОДАТАКА", true)

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

	ipw.cell("Подаци о странцу", true)

	ipw.moveY(16)
	ipw.line(0)
	ipw.moveY(9)

	ipw.putData("Презиме:", ipw.doc.Surname)
	ipw.putData("Име:", ipw.doc.GivenName)
	ipw.putData("Држављанство:", ipw.doc.NationalityFull)
	ipw.putData("Датум рођења:", ipw.doc.DateOfBirth)
	ipw.putData("Основ боравка:", ipw.doc.PurposeOfStay)
	addressLabel := "Пребивалиште\nи адреса стана:"
	if ipw.doc.AddressLabel == "prebivalište" {
		addressLabel = "Пребивалиште:"
	}
	ipw.putData(addressLabel, localization.JoinWithComma(ipw.doc.State, ipw.doc.GetFullAddress(true)))
	ipw.putData("Датум промене адресе:", ipw.doc.AddressDate)
	ipw.putData("Евиденцијски број\nстранца:", ipw.doc.PersonalNumber)
	ipw.putData("Пол:", ipw.doc.Sex)

	ipw.moveY(-8.67)
	ipw.line(0)
	ipw.moveY(9)
	ipw.cell("Подаци о документу", true)
	ipw.moveY(16)

	ipw.line(0)
	ipw.moveY(9)
	ipw.putData("Документ издаје:", ipw.doc.IssuingAuthority)
	ipw.putData("Број документа:", ipw.doc.DocRegNo)
	ipw.putData("Датум издавања:", ipw.doc.IssuingDate)
	ipw.putData("Важи до:", ipw.doc.ExpiryDate)

	ipw.moveY(-8.67)
	ipw.line(0)
	ipw.moveY(3)
	ipw.line(0)
	ipw.moveY(9)

	ipw.cell("Датум штампе: "+time.Now().Format("02.01.2006."), true)

	ipw.moveY(19)

	ipw.line(0.83)

	err = ipw.pdf.SetFontSize(9)
	if err != nil {
		panic(err)
	}

	ipw.moveY(4)

	ipw.pdf.SetX(ipw.leftMargin)

	if ipw.doc.pdfCyrillicLabels {
		ipw.cell("1. У чипу личне карте за странце, подаци о имену и презимену странца исписани су онако како су исписани на", false)
		ipw.pdf.SetX(ipw.leftMargin)
		ipw.moveY(9.7)
		ipw.cell("самом обрасцу личне карте за странца латиничним писмом.", false)
		ipw.pdf.SetX(ipw.leftMargin)
		ipw.moveY(9.7)
		ipw.cell("2. Ако се име или презиме странца састоји од две или више речи чија дужина прелази 30 карактера за име, односно", false)
		ipw.pdf.SetX(ipw.leftMargin)
		ipw.moveY(9.7)
		ipw.cell("односно 36 карактера за презиме, у чип се уписује пуно име и презиме странца, а на обрасцу личне карте за", false)
		ipw.pdf.SetX(ipw.leftMargin)
		ipw.moveY(9.7)
		ipw.cell("странца се уписује до 30 карактера за име, односно 36 карактера за презиме.", false)
	} else {
		ipw.cell("1. U čipu lične karte za strance, podaci o imenu i prezimenu stranca ispisani su onako kako su ispisani na samom", false)
		ipw.pdf.SetX(ipw.leftMargin)
		ipw.moveY(9.7)
		ipw.cell("obrascu lične karte za stranca latiničnim pismom.", false)
		ipw.pdf.SetX(ipw.leftMargin)
		ipw.moveY(9.7)
		ipw.cell("2. Ako se ime ili prezime stranca sastoji od dve ili više reči čija dužina prelazi 30 karaktera za ime, odnosno 36", false)
		ipw.pdf.SetX(ipw.leftMargin)
		ipw.moveY(9.7)
		ipw.cell("karaktera za prezime, u čip se upisuje puno ime i prezime stranca, a na obrascu lične karte za stranca se upisuje do", false)
		ipw.pdf.SetX(ipw.leftMargin)
		ipw.moveY(9.7)
		ipw.cell("30 karaktera za ime, odnosno 36 karaktera za prezime.", false)
	}

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
	ipw.cell("ЧИТАЧ ЕЛЕКТРОНСКЕ ЛИЧНЕ КАРТЕ: ШТАМПА ПОДАТАКА", true)

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

	ipw.cell("Подаци о странцу", true)

	ipw.moveY(16)
	ipw.line(0)
	ipw.moveY(9)

	ipw.putData("Презиме:", ipw.doc.Surname)
	ipw.putData("Име:", ipw.doc.GivenName)
	ipw.putData("Држављанство:", ipw.doc.NationalityFull)
	ipw.putData("Датум рођења:", ipw.doc.DateOfBirth)
	ipw.putData("Место рођења,\nопштина и држава:", ipw.doc.GetFullPlaceOfBirth())
	addressLabel := "Пребивалиште\nи адреса стана:"
	if ipw.doc.AddressLabel == "prebivalište" {
		addressLabel = "Пребивалиште:"
	}
	ipw.putData(addressLabel, ipw.doc.GetFullAddress(true))
	ipw.putData("Датум промене адресе:", ipw.doc.AddressDate)
	ipw.putData("Евиденцијски број\nстранца:", ipw.doc.PersonalNumber)
	ipw.putData("Пол:", ipw.doc.Sex)
	ipw.putData("Основ боравка:", ipw.doc.PurposeOfStay)
	ipw.putData("Напомена:", ipw.doc.ENote)

	ipw.moveY(-8.67)
	ipw.line(0)
	ipw.moveY(9)
	ipw.cell("Подаци о документу", true)
	ipw.moveY(16)

	ipw.line(0)
	ipw.moveY(9)
	ipw.putData("Назив документа:", ipw.doc.DocumentName)
	ipw.putData("Документ издаје:", ipw.doc.IssuingAuthority)
	ipw.putData("Број документа:", ipw.doc.DocRegNo)
	ipw.putData("Датум издавања:", ipw.doc.IssuingDate)
	ipw.putData("Важи до:", ipw.doc.ExpiryDate)

	ipw.moveY(-8.67)
	ipw.line(0)
	ipw.moveY(3)
	ipw.line(0)
	ipw.moveY(9)

	ipw.cell("Датум штампе: "+time.Now().Format("02.01.2006."), true)

	ipw.moveY(19)

	ipw.line(0.83)

	err = ipw.pdf.SetFontSize(9)
	if err != nil {
		panic(err)
	}

	ipw.moveY(4)

	ipw.pdf.SetX(ipw.leftMargin)

	if ipw.doc.pdfCyrillicLabels {
		ipw.cell("1. У чипу дозволе за привремени боравак и рад, подаци о имену и презимену имаоца дозволе исписани су онако", false)
		ipw.pdf.SetX(ipw.leftMargin)
		ipw.moveY(9.7)
		ipw.cell("како су исписани на самом обрасцу дозволе за привремени боравак латиничним писмом.", false)
		ipw.pdf.SetX(ipw.leftMargin)
		ipw.moveY(9.7)
		ipw.cell("2. Ако се име или презиме странца састоји од две или више речи чија дужина прелази 30 карактера за име,", false)
		ipw.pdf.SetX(ipw.leftMargin)
		ipw.moveY(9.7)
		ipw.cell("односно 36 карактера за презиме у чип се уписује пуно име странца, а на обрасцу дозволе за привремени боравак", false)
		ipw.pdf.SetX(ipw.leftMargin)
		ipw.moveY(9.7)
		ipw.cell("се уписује до 30 карактера за име, односно 36 карактера за презиме.", false)

	} else {
		ipw.cell("1. U čipu dozvole za privremeni boravak i rad, podaci o imenu i prezimenu imaoca dozvole ispisani su onako kako su", false)
		ipw.pdf.SetX(ipw.leftMargin)
		ipw.moveY(9.7)
		ipw.cell("ispisani na samom obrascu dozvole za privremeni boravak latiničnim pismom.", false)
		ipw.pdf.SetX(ipw.leftMargin)
		ipw.moveY(9.7)
		ipw.cell("2. Ako se ime ili prezime stranca sastoji od dve ili više reči čija dužina prelazi 30 karaktera za ime, odnosno 36 karaktera", false)
		ipw.pdf.SetX(ipw.leftMargin)
		ipw.moveY(9.7)
		ipw.cell("za prezime u čip se upisuje puno ime stranca, a na obrascu dozvole za privremeni boravak se upisuje do 30 karaktera za", false)
		ipw.pdf.SetX(ipw.leftMargin)
		ipw.moveY(9.7)
		ipw.cell("ime, odnosno 36 karaktera za prezime.", false)
	}

	ipw.moveY(9.7)

	ipw.line(0)
}
