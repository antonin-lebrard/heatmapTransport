'use strict'

const mapEl = document.querySelector('#map')
mapEl.style.height = `${window.innerHeight - 40}px`
window.addEventListener('resize', () => mapEl.style.height = `${window.innerHeight - 40}px`)

/**
 * @type {Map}
 */
window.map = L.map('map').setView([48.85, 2.35], 11)
window.map.addLayer(L.tileLayer(
  'https://{s}.tile.openstreetmap.org/{z}/{x}/{y}.png',
  {
    minZoom: 1,
    maxZoom: 30,
    attribution: 'Map data © <a href="https://openstreetmap.org">OpenStreetMap</a> contributors'
  }
))

function drawVisibleStops() {
  if (pathsButton.value === 'R') { return }
  const bounds = window.map.getBounds()
  window.stopsLeafletMarkers.forEach((marker) => {
    if (bounds.contains(marker.getLatLng())) {
      if (!marker._addedToMap) marker.addTo(window.map)
    } else {
      marker.removeFrom(window.map)
      marker._addedToMap = false
    }
  })
}

window.map.on('zoomend', drawVisibleStops)
window.map.on('moveend', drawVisibleStops)

const maxTimeInput = document.querySelector('#maxTimeInput')
const departureTimeInput = document.querySelector('#departureTimeInput')
const heatmapButton = document.querySelector('#heatmapButton')
const pathsButton = document.querySelector('#pathsButton')
// const chosenDepartureDiv = document.querySelector('#chosenDeparture')

heatmapButton.addEventListener('click', () => {
  window.stopsLeafletMarkers.forEach((marker) => {
    marker.setStyle({ color: 'grey' })
  })
  const stopId = window.stopId
  const maxTime = maxTimeInput.value
  const departureTime = departureTimeInput.value
  fetchHeatmap(stopId, maxTime, departureTime, (stopTimes) => {
    stopTimes.forEach((stopTime) => {
      // if (stopTime.total_time === 0) {
      //   chosenDepartureDiv.innerHTML = stopTimes.departure
      // }
      window.stopIdToStop[stopTime.stop_id][1].setStyle({ color: colorForTotalTime(stopTime.total_time, maxTime) })
    })
  })
})

window.lines = []
pathsButton.addEventListener('click', () => {
  if (pathsButton.value === 'P') {
    pathsButton.value = 'R'
    pathsButton.title = 'R'
    const stopId = window.stopId
    const maxTime = maxTimeInput.value
    const departureTime = departureTimeInput.value
    fetchPaths(stopId, maxTime, departureTime, (paths) => {
      paths.forEach((path) => {
        const latlngFrom = [window.stopIdToStop[path.from][0].lat, window.stopIdToStop[path.from][0].lon]
        const latlngTo = [window.stopIdToStop[path.to][0].lat, window.stopIdToStop[path.to][0].lon]
        const line = L.polyline([ latlngFrom, latlngTo ], { color: colorForTotalTime(path.time, maxTime) })
        window.lines.push(line)
        line.addTo(window.map)
      })
    })
    window.stopsLeafletMarkers.forEach(marker => {
      marker.removeFrom(window.map)
      marker._addedToMap = false
    })
  } else {
    pathsButton.value = 'P'
    pathsButton.title = 'P'
    window.lines.forEach(line => line.removeFrom(window.map))
    window.lines = []
    drawVisibleStops()
  }
})

window.stopsLeafletMarkers = []
window.stopIdToStop = {}
fetchStops((stops) => {
  stops.forEach((stop) => {
    const marker = L.circleMarker([stop.lat, stop.lon], { radius: 5, color: 'grey', weight: 1, fillOpacity: 1 })
    const popup = L.popup().setContent(stop.name)
    marker.bindPopup(popup)
    marker.on('click', () => {
      window.stopId = stop.id
      fetchDepTimesForStop(stop.id, departureTimes => {
        departureTimes = departureTimes.map(el => el.substring(0, el.length - 3))
        if (departureTimes.length > 1)
          departureTimes = [departureTimes[0], departureTimes[departureTimes.length - 1]]
        popup.setContent(
          `<span style="margin: auto">${stop.name}</span>` +
          `<div style="display: flex; flex-direction: row; flex-wrap: wrap;">${departureTimes.map(el => `<span>${el}</span>`)}</div>`
        )
      })
    })
    window.stopIdToStop[stop.id] = [stop, marker]
    window.stopsLeafletMarkers.push(marker)
    marker._addedToMap = true
    marker.addTo(window.map)
  })

  window.stopId = stops[0].id
})
