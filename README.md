# Lernpfad

Verwaltung von Lern-Kursen.
Die Lernkurse können strukturiert aufgebaut werden. Der Benutzer muss einen Account im "data" haben und sich dort authentifizieren.
Im Lernpfad müssen die Berechtigungen für die Kurse vergeben werden.

## Bermerkungen

Es wird go in der Version 1.18 genutzt, weil die azure Bibiothek nur in dieser Version arbeitet.
Das Deployment benötigt ca. 5-10 min - es wird das go-programm compiliert.

## Information zur Installation

Es wird ein Go-Server für die Anwendung genutzt, der über eine `.env` Datei auf eine cosmos-db zugreift. Diese enthält eine Datenbank `data` mit einem Container `user`. Hier werden Name und Password authentifiziert und die Rolle `admin` übernommen/geprüft.

Die Client-Anwendung nutzt ausschließlich den go-server. Mit der Authentifizierung wird ein jwt-token mit der Gültigkeit für eine Stunde als Coockie gesetzt.

Dokumente und Filme aus den Kursen können auf Fremd-Verzeichnisse zeigen (z.B. youtube) oder auf den azure storage Container. Hierzu muss mit dem Namen des Kurs ein Unterverzeichnis angelegt werden, das dann die Dokuemnte enthält.

## Starten im Dev Mode

server.go

```
	LOCAL_TEMPDIR = "lernpfad"
	PORT          = "8081"
```

Zugriff auf lokales Verzeichnis "lernpfad" und Port 8081.

Main.Layout.Vue

```
server: "http://localhost:8081",
```

Der Client greift auf den Port des Go-Server zu.

Anschließend kann der go Debugger gestartet werden - der Server läuft dann auf Port 8081.
Außerdem wird mit
`quasar dev`
der Client-Server via qausar zum debuggen auf Port 8080 gestratet.

## Deploymend azure

server.go

```
	LOCAL_TEMPDIR = "/lernpfad"
	PORT          = "8080"
```

Zugriff auf azure Verzeichnis "/lernpfad" und Port 8080.
Das Verzeichnis muss im azure-Container bereitgestellt/eingestellt werden.

Main.Layout.Vue

```
server: "",
```

Der Client nutzt die Standard URL, Zugriff erfolgt über den gleichen Server.

Anschließend kann mit
`npm run build`
das Go-Programm compiliert werden und mit
`quasar build`
auch die Client-Software. Die Software steht dann unter `"./dist/spa"`. Im npm-Build ist das Generieren des quasar-Clients auskommentiert, damit nicht in dem azure-Container die Qusar-Anwendung generiert werden muss - die Anwendung wird mit deployed.

Anschlißend kann auf den azure-Container deplayed werden. Dort wird der Go-Server noch einmal compiliert.

## Voraussetzungen

-   go Version 1.18
-   Verzeichnis für Daten "/lernpfad"
-   der go-server sendet auf Port 8080

## Kurse

Die Kurse werden im "data" Verzeichnis für den Benutzer "all" als JSON abgelegt. Hier liegen die Vorlagen.

## Benutzer und Berechtigungen

Die Benutzer benötigen einen Account im 'data' und müssen sich dort auhtentifizieren. Als Benutzer sihet man aber nur die Kurse, zu denen man Berechtigung hat.

Für die Berechtigungen werden Gruppen (group) definiert, z.B. "gast" - bekommt ein Benutzer diese Gruppe, dann sieht er alle Kurse dieser Gruppe.
An diesen Kursen kann er "Teilnehmen" - dazu wird der Kurs in das Verzeichnis des Benutzers kopiert und steht dann zur Bearbeitung zur Verfügung.

Über Regeln (policy) werden den Gruppen die Berechtigungen zu den Kursen zugeordnet.

Regeln: Gruppe -> Kurse
Gruppen: Benutzer -> Gruppe

Beispiel Regeln:
p, klasse1, data, all, Kurse, get
p, klasse1, data, all, Kurs_G000, get
p, klasse1, data, all, Kurs_M001, get

Die Gruppe "klasse1" beinhaltet die Kurse Kurs_G000 und Kurs_M001. Es kann die Funktion "Kurse" - also eine Liste aller Kurse - aufgerufen werden.
Im Verzeichnis "all" darf nur lesend (get) zugegriffen werden.

Beispiel Gruppen:
g, jank, klasse1, data
g, superjoda, admin, data
g, superjoda, klasse1, data

Der Benuter "jank" darf auf die Gruppe "klasse1" zugreifen.
Der Benutzer "superjoda" darf auf Gruppe "klasse1" und auf die "admin" Berechtigungen zugreifen.
