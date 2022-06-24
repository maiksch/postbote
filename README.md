# Postbote

API Explorer in the terminal

## Workflow

Ich befinde mich in einem API Projekt in Neovim und möchte einen gerade entwickelten Endpunkt testen.
Ich verwende einen konfiguriertes Keymapping (tmux, neovim,.. hier ist vieles vorstellbar,
der User muss es einfach konfigurieren können) und öffne Postbote. Es werden mir zuvor
verwendete Endpunkte vorgeschlagen von denen ich per fuzzy search einen auswählen könnte,
ich möchte aber einen neuen verwenden. 

Ich konfiguriere den Endpunkt, in dem ich einen Namen, eine URL, Query Parameter und einen Body
angebe. Der Request wird abgeschickt und ich sehe die Antwort

## TODOs

### MVP
* [ ] Angeben von Name, HTTP-Method, URL, Query Parametern und eines Bodies
* [ ] Speichere einen Request in einer Datei (.postbote), um diese wieder zu verwenden
* [ ] Fuzzy Finder für alte Requests

