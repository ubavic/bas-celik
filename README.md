# Baš Čelik

**Baš Čelik** je čitač elektronskih ličnih karata i zdravstvenih knjižica. Program je osmišljen kao zamena za zvanične aplikacije poput *Čelika*. Nažalost, zvanične aplikacije mogu se pokrenuti samo na Windows operativnom sistemu, dok Baš Čelik funkcioniše na tri operativna sistema (Windows/Linux/OSX).

![Interfejs](assets/ui.png)

Aplikacija bi trebalo da podržava očitavanje svih ličnih karata, ali za sada nije testirana na starim (izdate pre avgusta 2014. godine), kao ni na najnovijim (izdate nakon februara 2023. godine). Unapred sam zahvalan za bilo kakvu povratnu informaciju.

## Upotreba

Povežite čitač za računar i pokrenite aplikaciju. Ubacite karticu u čitač. Aplikacija će pročitati informacije sa kartice i prikazati ih. Tada možete sačuvati PDF pritiskom na donje desno dugme.

Kreirani PDF dokument izgleda maksimalno približno dokumentu koji se dobija sa zvaničnim aplikacijama.

### Pokretanje na Linuksu

Aplikacija zahteva instalirane `ccid` i `opensc`/`pcscd` pakete. Nakon instalacije ovih paketa, neophodno je i pokrenuti `pcscd` servis:

```
sudo systemctl start pcscd
sudo systemctl enable pcscd
```

### Pokretanje u komandnoj liniji

Aplikacija prihvata sledeće opcije:

 + `-verbose`: tokom rada aplikacije detalji o greškama će biti ispisani u komandnu liniju
 + `-pdf PATH`: grafički interfejs neće biti pokrenut, a sadržaj dokumenta će biti direktno sačuvan u PDF na `PATH` lokaciji.
 + `-help`: informacija o opcijama će biti prikazana

## Preuzimanje 

Izvršne datoteke poslednje verzije programa možete preuzeti sa [Releases](https://github.com/ubavic/bas-celik/releases) stranice.

Dostupan je i [AUR paket](https://aur.archlinux.org/packages/bas-celik) za Arch korisnike. 

## Kompilacija

Potrebno je posedovati `go` kompajler. Na Linuksu je potrebno instalirati i `libpcsclite-dev` i [pakete za Fyne](https://developer.fyne.io/started/#prerequisites) (možda i `pkg-config`).

Nakon preuzimanja repozitorijuma, dovoljno je pokrenuti

```
go build main.go
```

### Kroskompilacija

Uz pomoć [fyne-cross](https://github.com/fyne-io/fyne-cross) programa moguće je na jednom operativnom sistemu iskompajlirati program za sva tri operativna sistema. Ovaj program zahteva Docker na vašem operativnom sistemu.

## Planirane nadogradnje

 + Omogućavanje potpisivanja dokumenata sa LK

## Poznati problemi

 + Na Windowsu, aplikacija u nekim slučajevima neće pročitati karticu ako je kartica ubačena u čitač nakon pokretanja programa. U tom slučaju, dovoljno je restartovati program.
 + Podaci na zdravstvenoj kartici su kodirani sa (meni) nepoznatim kodranjem. Program dekodira uspešno većinu karaktera, ali ne sve. Zbog toga se mogu desiti greške prilikom ispisa podataka.

Ni jedan od problema ne utiče na "sigurnost" vašeg dokumenta. Baš Čelik isključivo čita podatke sa kartice.

## Licenca 

Program i izvorni kôd su objavljeni pod [*MIT* licencom](LICENSE).

Font `free-sans-regular` je objavljen pod [*SIL Open Font* licencom](assets/LICENSE).
