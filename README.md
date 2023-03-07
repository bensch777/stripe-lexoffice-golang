# Rechnungen und Rechnungsdaten aus Stripe nach Lexoffice exportieren
 
Export der PDF Rechnungen aus Stripe und Import als Belege mit Rechnungsdaten in Lexoffice.

In Lexoffice erscheint die PDF-Rechnung als Einnahmebeleg unter dem Sammelkunden. Als Belegnummer, Belegdatum, Betrag, Steuern werden die Informationen aus Stripe übernommen und der Beleg unter der Einnahmeart verbucht.

# Einrichtung

- Stripe: [Eingeschränkten API-Schlüssel](https://dashboard.stripe.com/apikeys/create) erstellen.
- lexoffice: [API-Schlüssel](https://app.lexoffice.de/addons/public-api) erstellen.
- Anschließend die API-Schlüssel in der main.go hinterlegen und den Zeitraum für die Exporte anpassen.
```golang
var stripeKey = "rk_";
var lexofficeKey = "";
var startdate = "2023-01-01T00:00:00"
var enddate = "2023-01-02T23:59:59"```