# Heatmap Public Transport

<a href="https://github.com/antonin-lebrard/heatmapTransport/blob/master/showcase.mp4?raw=true" rel="Presentation of the interface, not usable for blindness affected persons, so I don't know why I write this">![](https://raw.githubusercontent.com/antonin-lebrard/heatmapTransport/master/showcase.gif)</a>

This is a tool to interactively generate heatmaps of the time to reach other stations from one chosen station.

### Usage

This project has been developped with the Paris public transports network in mind, precisely the [RATP](https://www.ratp.fr/) part of its network.<br>
There has been some choices made in the structuration of the data while writing the code and reading the RATP files, and it should pose some problems with other data sources.<br>
But it should be adaptable if jumping in the code do not scares you too much. The file that handles every data types from the input data is [here](https://github.com/antonin-lebrard/heatmapTransport/blob/master/internal/pkg/csvLoading.go)

The initial data for the stations, timetables, and transfers is taken from here: https://data.ratp.fr/explore/dataset/offre-transport-de-la-ratp-format-gtfs/information/

There two input option when interacting with the map:
- Max time (in seconds): to consider as the maximum limit for transport time between each station and chosen initial station.
   <br>The color code is derived from this value, as:
   - green: <= 1/3 of Max time
   - yellow: > 1/3 and <= 2/3 of Max time
   - red: > 2/3 of Max time<br>
   (The colors are configurable in [this function](https://github.com/antonin-lebrard/heatmapTransport/blob/master/mapDisplayJs/utils.js#L28))
- Departure Time: to consider as the starting moment from which to traverse the graph of stations. The time to reach one station at 17h00 should be different at 4h00.

There is one button next to the `Get Heatmap` one, to display the paths taken to reach all considered station, instead of the stations. The color code is taken from the one for the normal heatmap.

### Installation

- Install golang: https://github.com/golang/go/wiki#working-with-go
    - There should be a correct installation ready with $GOPATH set and your go binary in your $PATH. It did takes some time for me to do a complete and correct install.
- `git clone https://github.com/antonin-lebrard/heatmapTransport.git` in your `$GOPATH/src`
- `cd heatmapTransport`
- `go build -o heatmap .`
- Download the full GTFS archive from RATP [here](https://data.ratp.fr/explore/dataset/offre-transport-de-la-ratp-format-gtfs/information/), its the `RATP_GTFS_FULL` link.
- Extract the archive and take the `stops.txt`, `stop_times.txt`, `transfers.txt` files and put them in `heatmapTransport/ratp`
- Now execute `./heatmap`. It will takes a long time (and RAM, like ~4Go of it) to gather all the necessary data and construct an adapted graph format.
   - This step has taken me ~1h30 to do, but at the end of it, it will save the graph into a gigantic text file (1.6 Go for me) which will only takes 20 seconds to reload at the next launch

When `./heatmap` has finished loading the GTFS data and saved its graph to disk (or read it back), the application will open a server which the [`mapDisplayJs`](https://github.com/antonin-lebrard/heatmapTransport/tree/master/mapDisplayJs) part of this project will use.

You can directly launch it with the `file:///` url in your browser, for me it would be `file:////home/antonin/go/src/heatmapTransport/mapDisplayJs/index.html`

### Credits

Thanks to the [OpenTripPlanner (OTP) project](http://www.opentripplanner.org/) and its [technical documentation](http://docs.opentripplanner.org/en/latest), particularly its [Bibliography page](http://docs.opentripplanner.org/en/latest/Bibliography)
which helped me to realise this might be possible.<br>
Thanks to [leaflet](https://leafletjs.com/) as usual the simplest tool to present simple data on a map.


