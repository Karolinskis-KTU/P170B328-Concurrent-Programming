# Inžinerinis projektas - duomenų lygiagretumo priemonių taikymas

## Vertinimo tvarka
- Ataskaita - maks. 5 balai
- Programos veikimas ir darbų tarp gijų paskirstymas - maks. 2,5 balo
- Atsakymai į klausimus - maks. 2,5 balo

## Ataskaitos strukūra
1. Titulinis lapas ir užduotis.
2. **Užduoties analizė ir sprendimo metodas** - trumpai aprašyti sprendžiamą problemą ir koks pasirinktas sprendimo būdas bei kaip pritaikomos duomenų lygiagretumo priemonės.
3. **Testavimas ir programos vykdymo istrukcija** - parodyti, kad programa gauna **teisingą** rezultatą ir instrukciją, kaip pasileisti programą.
4. **Vykdymo laiko kitimo tyrimas:**
    - Paruošti bent 8 skirtingos apimties duomenų rinkinius (rekomenduojama, kad mažiausias rinkinys būtų apdorojamas vienos gijos greičiau nei per 10ms, didžiausias - ilgiau nei per 2 min).
    - Kiekvienam duomenų rinkiniui rasti optimalų gijų ar procesų kiekį, parodyti grafikus, kaip kinta vykdymo laikas priklausomai nuo gijų ar procesų kiekio kiekvienam rinkiniui. Vykdymo laiką matuoti keletą kartų ir pateikti tų vykdymo laikų vidurkius.
    - Gautus rezultatus aprašyti: pakomentuoti diagramas, aprašyti, kodėl, jūsų nuomone, reultatai yra būtent tokie, kokius gavote.
5. **Išvados ir literatūra** - ką pavyko naujo sužinoti, kas tikėto ar netikėto buvo pastebėta darbo metu, pakomentuoti pasirinktos priemonės tinkamumą pasirinktai problemai spręsti.

## Priemonės
Tikslas - išbandyti duomenų lygiagretumo priemones, kurios gijas valdo automatiškai, bet paleidžiamų gijų kiekį galima nurodyti kokį norima. Galima rinktis savo priemones (turi būti suderinta su lab. darbų dėstytoju) arba iš siūlomų:
* C++ su OpenMP `parallel for`
* C++ su MPI `Scatter / Gather`
* C++ su `CUDA`
* Haskell `rpar` ir `parListChunk` ar kt. strategija
* Python `multiprocessing.Pool`
* Ruby su `parallel gem`
* Java / Kotlin `parallel stream`
* C# `parallel LINQ`
* Scala `parallel collections`
* F# `Pseq`
* Rust `rayon`

## Pateikimo ir gynimo tvarka
* Projekto pristatymo laikas pateiktas Moodle, gintis reikia paskirtą savaitę.
* Jei naudojamos ne siūlomos priemonės, jos turi būti iš anksto suderintos su lab. darbų dėstytoju.
* Pateikiamos ataskaitos formatas - PDF.
* Ataskaitą, programą ir duomenų failus (1 variantą) įkelti į Moodle prieš projekto gynimą.
* Neatsiskaičius projekto laiku pakartotinai gintis galima 16 savaitę.

## Užduotis
* Tiems, kas turi skaitinių metodų modulį, **rekomenduojama** išlygiagretinti antrojo namų darbo optimiavimo uždavinį.
* Galima susigalvoti savo užduotį, bet reikia iš anksto sudetinti su laboratorinių darbų dėstytoju. Užduotis turi būti duomenų lygiagretumo užduotis, t.y., visiems duomenų rinkinio elementams taikoma ta pati operacija ir gaunamas rezultatas.

## Reikalavimai programai
* Turi būti galimybė keisti gijų ar procesų kiekį nepriklausomai nuo pasirinkto duomenų rinkinio;
* Turi būti galimybė gynimo metu lengvai keisti, kuris duomenų rinkinys apdorojamas ir kiek gijų ar procesų naudojama;
* Darbo tikslas yra gauti kuo geresnė pagreitėjimą kuo mažiau komplikuojant programos kodą;
* Darant vykdymo laiko tyrim1 program1 vykdykite ne derinimo re=imu, pv,. jei dirbate su Visual Studio ir C++, pakeiskite Debug konfigūraciją į Release - įsijungs visos kompiliatoriaus optimizacijos ir programa veiks greičiau.