# Rechnungen und Rechnungsdaten aus Stripe nach Lexoffice exportieren
 
Export der PDF Rechnungen aus Stripe und Import als Belege mit Rechnungsdaten in Lexoffice.

In Lexoffice erscheint die PDF-Rechnung als Einnahmebeleg unter dem Sammelkunden. Als Belegnummer, Belegdatum, Betrag, Steuern werden die Informationen aus Stripe übernommen und der Beleg unter der Einnahmeart verbucht.

# Einrichtung

1. Stripe: [Eingeschränkten API-Schlüssel](https://dashboard.stripe.com/apikeys/create) erstellen. Es wird nur die Berechtigung "Rechnung lesen" benötigt.
2. lexoffice: [API-Schlüssel](https://app.lexoffice.de/addons/public-api) erstellen.
3. API-Schlüssel in der main.go hinterlegen und den Zeitraum für die Exporte anpassen.
```golang
var stripeKey = "rk_";
var lexofficeKey = "";
var startdate = "2023-01-01T00:00:00"
var enddate = "2023-01-02T23:59:59"
```
4. Terminal öffnen und die main.go ausführen.
```bash
go run .\main.go
```

**Hinweis:** Die Rechnungen werden gleichzeitig in das Verzeichnis **/invoices/2023/** heruntergeladen. Rechnungen, die sich in diesem Verzeichnis befinden, werden nicht erneut heruntergeladen und somit auch nicht erneut an lexoffice übergeben.