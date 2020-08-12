# Search your Google location history

Simple program I wrote to explore my location history. In particular I wanted
to see how frequently I visited specific places.

Your location history can be exported using https://takeout.google.com/settings/takeout.

Pass one or more of the exported JSON files as arguments to the script. Use the
`--name` flag to search by the name of a place. Results are output as JSON. You
can use [jq](https://stedolan.github.io/jq/) to parse and format the output as
needed.

Example:

```
$ go run main.go --name "boulder+" ~/Downloads/Takeout/Location\ History/Semantic\ Location\ History/**/*.json | jq -c '.[] | {Name,Start,End,PlaceConfidence}'
{"Name":"boulder+ Gym","Start":"2019-11-22T11:35:08.313+08:00","End":"2019-11-22T13:56:16.208+08:00","PlaceConfidence":"HIGH_CONFIDENCE"}
{"Name":"boulder+ Gym","Start":"2019-11-25T20:37:20.528+08:00","End":"2019-11-25T22:08:55.303+08:00","PlaceConfidence":"HIGH_CONFIDENCE"}
{"Name":"boulder+ Gym","Start":"2019-11-27T20:43:14.507+08:00","End":"2019-11-27T22:22:39.478+08:00","PlaceConfidence":"HIGH_CONFIDENCE"}
```
