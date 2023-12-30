# Doprinos razvoju Baš Čelik aplikacije

Zahvalan sam svima koji žele da doprinesu razvoju Baš Čelik aplikacije. Bez mnogobrojnih komentara i ispravki [saradnika](https://github.com/ubavic/bas-celik/graphs/contributors), ova aplikacija bi bila mnogo lošija.

Da bi komunikacija i rad bili što efikasniji, navodim neke smernice za nove saradnike koji otvaraju *issue*-e ili *pull request*-ove.  

## Otvaranje novih *issue*-a

1. [Issues](https://github.com/ubavic/bas-celik/issues) je mesto gde se iznose uočeni problemi sa aplikacijom, ali i mesto gde se mogu predlagati nove funkcionalnosti/poboljšanja, iznositi kritike, ili prosto postavljati pitanja.
2. Pre otvaranja *issue*-a proverite da li već postoji *issue* na istu temu. Ako postoji, napišite vaš komentar u postojećem *issue*-u, čak i ako je isti trenutno zatvoren.
3. Ako ste uočili više različitih problema sa aplikacijom koji ne deluju povezano, za svaki problem otvorite novi *issue*.
4. Prilikom prijave problema sa čitanjem nekog dokumenta, bilo bi dobro da navedete tip dokumenta, okvirni datum izdavanja dokumenta, *ATR* kôd dokumenta, verziju operativnog sistema na kom koristite Baš Čelik. Idealno bi bilo i da navedete da li je problem uočen i prilikom čitanja drugih kartica istog tipa.
5. *ATR* kod možete slobodno deliti jer ne sadrži identifikacione informacije, ali *PDF* i *JSON* datoteke nemojte deliti jer sadrže privatne informacije! Ako se problem javi prilikom generisanja ovih datoteka, opišite detaljno problem. 
6. Nema potrebe postavljati *screenshot*-ove aplikacije, osim u slučaju kada je sam korisnički interfejs problem (što se nikad nije desilo).
7. Ako otvarate *issue* radi postavljanja pitanja, prvo dobro pročitajte postojeću dokumentaciju koja se nalazi u [README.md](README.md) datoteci. 

## Otvaranje novih *pull request*-ova

1. Ako se *pull request* tiče male ispravke (npr, ispravljanje očiglednog bug-a, update neke biblioteke ili predlog korišćenja neke druge funkcije, itd), nije neophodno otvarati poseban *issue*. U komentaru *PR*-a navedite razlog izmene pogotovu ako se radi o ispravci nekog bug-a. Ako je ispravka trivijalna (npr. ispravka greške u pisanju ili formatiranju koda, dodavanje korisnog komentara), ne morate navoditi komentar u *PR*-u.
2. Ako *PR* podrazumeva uvođenje dodatnih funkcionalnosti, ili značajne izmene koda, tada je neophodno da prvo otvorite poseban *issue* u kom će te opisati razloge otvaranja *PR*-a. Da ne biste došli u situaciju da vam odbijem veliki *PR* na koji ste utrošili vreme, pre početka proverite da li sam uopšte zainteresovan za takav *PR*.
3. Trudim se da git istoriju projekta držim čistom. Ako vaš *PR* sadrži više komitova, potrudite se da svaki komit ima smislenu poruku napisanu na engleskom. Kod svakog komita mora se uspešno kompajlirati i formatiran je sa `gofmt`. Svaki komit mora predstavljati smislenu izmenu koda.

