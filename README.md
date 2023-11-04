# Baš Čelik

**Baš Čelik** je čitač elektronskih ličnih karata i zdravstvenih knjižica. Program je osmišljen kao zamena za zvanične aplikacije poput *Čelika*. Nažalost, zvanične aplikacije mogu se pokrenuti samo na Windows operativnom sistemu, dok Baš Čelik funkcioniše na tri operativna sistema (Windows/Linux/OSX).

Baš Čelik je besplatan program, sa potpuno otvorenim kodom dostupnim na adresi [github.com/ubavic/bas-celik](https://github.com/ubavic/bas-celik). 

![Interfejs](assets/ui.png)

## Upotreba

Povežite čitač za računar i pokrenite Baš Čelik. Ubacite karticu u čitač. Program će pročitati informacije sa kartice i prikazati ih. Tada možete sačuvati PDF pritiskom na donje desno dugme.

Kreirani PDF dokument izgleda maksimalno približno dokumentu koji se dobija sa zvaničnim aplikacijama.

### Pokretanje na Linuksu

Baš Čelik zahteva instalirane `ccid` i `opensc`/`pcscd` pakete. Nakon instalacije ovih paketa, neophodno je i pokrenuti `pcscd` servis:

```
sudo systemctl start pcscd
sudo systemctl enable pcscd
```

### Pokretanje u komandnoj liniji

Baš Čelik prihvata sledeće opcije:
 
 + `-atr`: ATR kôd kartice biće prikazan u konzoli. 
 + `-help`: informacija o opcijama biće prikazana u konzoli.
 + `-json PATH`: grafički interfejs neće biti pokrenut, a sadržaj dokumenta biće direktno sačuvan u JSON datoteku na `PATH` lokaciji.
 + `-list`: lista raspoloživih čitača biće prikazana u konzoli.
 + `-pdf PATH`: grafički interfejs neće biti pokrenut, a sadržaj dokumenta biće direktno sačuvan u PDF datoteku na `PATH` lokaciji.
 + `-verbose`: tokom rada aplikacije detalji o greškama biće prikazani u konzoli.
 + `-version`: informacija o verziji programa biće prikazana u konzoli.
 + `-reader INDEX`: postavlja odabrani čitač za čitanje podataka. Parametar `INDEX` označava prirodan broj koji je naveden u ispisu `list` komande. Izbor utiče samo na čitanje sa `pdf` i `json` opcijama.

U slučaju `json` i `pdf` opcija, program ne dodaje ekstenziju na kraj lokacije koju je korisnik naveo.

Pri pokretanju sa `atr`, `json` ili `pdf` opcijom, program očekuje da je kartica smeštena u čitač i neće čekati na ubacivanje kartice kao što je to slučaj sa grafičkim okruženjem.

Pri pokretanju sa `atr`, `help`, `list` ili `version` opcijama podaci sa kartice neće biti očitani (osim eventualno ATR koda). Program će prestati izvršavanje nakon ispisa odgovarajuće informacije.  

### Čitači i drajveri

Baš Čelik bi trebalo da funkcioniše sa svim čitačima pametnih kartica koji su trenutno dostupni u prodaji (Gemalto, Hama, Samtec...). Korisnici Windows (7, 8, 10, 11) i macOS operativnih sistema ne moraju da instaliraju ni jedan dodatni program (drajver).

## Preuzimanje 

Izvršne datoteke poslednje verzije programa možete preuzeti sa [Releases](https://github.com/ubavic/bas-celik/releases) stranice.

Dostupan je i [AUR paket](https://aur.archlinux.org/packages/bas-celik) za Arch korisnike. 

## Kompilacija

Potrebno je posedovati `go` kompajler. Na Linuksu je potrebno instalirati i `libpcsclite-dev` i [pakete za Fyne](https://developer.fyne.io/started/#prerequisites) (možda i `pkg-config`).

Nakon preuzimanja repozitorijuma, dovoljno je pokrenuti

```
go build main.go
```

Prva kompilacija može potrajati nekoliko minuta (i do deset), jer je neophodno da se preuzmu i kompajliraju sve Golang biblioteke. Sve naredne kompilacije se izvršavaju u nekoliko sekundi.

### Kroskompilacija

Uz pomoć [fyne-cross](https://github.com/fyne-io/fyne-cross) programa moguće je na jednom operativnom sistemu iskompajlirati program za sva tri operativna sistema. Ovaj program zahteva Docker na vašem operativnom sistemu.

## Planirane nadogradnje

 + Potpisivanje dokumenata sa sertifikatom smeštenim na LK
 + Verifikacija podataka na karticama
 + Opcija promene PIN koda na LK
 + Podrška za više čitača
 + Podrška za proveru dostupnosti novih verzija

## Poznati problemi

Baš Čelik bi trebalo da podržava očitavanje svih ličnih karata i zdravstvenih knjižica. Unapred sam zahvalan za bilo kakvu povratnu informaciju.

Za sada su registrovani naredni problemi:

 + Na Windowsu, aplikacija u nekim slučajevima neće pročitati karticu ako je kartica ubačena u čitač nakon pokretanja programa. U tom slučaju, dovoljno je restartovati program.

Ni jedan od problema ne utiče na "sigurnost" vašeg dokumenta. Baš Čelik isključivo čita podatke sa kartice.

### Ćirilica i latinica

Program prikazuje i eksportuje podatke onako kako su zapisani na kartici. Ako na nekom dokumentu uočite podatke na oba pisma, u pitanju nije *bug* već stanje na kartici.

## Slični projekti

Postoje takođe i drugi projekti otvorenog koda koji imaju izvesne sličnosti sa *Baš Čelikom*:

 + [JFreesteel](https://github.com/grakic/jfreesteel) i [jevrc](https://github.com/grakic/jevrc) Java programi čiji kôd mi je pomogao pri implementaciji nekih delova Baš Čelika.
 + [SerbianIdReader](https://github.com/lazarbankovic/serbianIdReader) Rust program za očitavanje ličnih karata.
 + [mup-rs-api-delphi](https://github.com/obucina/mup-rs-api-delphi), [BashChelik](https://github.com/neman/BashChelik) i [Saobracajna.NET](https://github.com/clearpath/Saobracajna.NET) wraperi u različitim jezicima za svanične MUP-ove biblioteke (sličost u nazivu sa jednom od biblioteka je slučajna).

## Licenca 

Program i izvorni kôd su objavljeni pod [*MIT* licencom](LICENSE).

Font [*Liberation*](https://github.com/liberationfonts/liberation-fonts) je objavljen pod [*SIL Open Font* licencom](assets/LICENSE).
