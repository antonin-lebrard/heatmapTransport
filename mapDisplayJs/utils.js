'use strict'


Array.prototype.l = function() {
  return this[this.length - 1]
}
Array.prototype.f = function() {
  return this[0]
}

String.prototype.bw = function(char) {
  return this.length > 0 && this[0] === char
}

function get(url, cb) {
  let xhttp = new XMLHttpRequest()
  xhttp.onreadystatechange = function() {
    if (xhttp.readyState === 4 && xhttp.status === 200) {
      cb(xhttp.responseText);
    }
  }
  xhttp.open("GET", url, true)
  xhttp.send()
}

const base = 'http://localhost:4000/'

function colorForTotalTime(totalTime, withMaxTime) {
  let color = 'rgb('
  if (totalTime > (withMaxTime * (1 / 3))) {
    color += '255,'
  } else {
    color += '0,'
  }
  if (totalTime < (withMaxTime * (2 / 3))) {
    color += '255,'
  } else {
    color += '0,'
  }
  color += '0)'
  return color
}

/**
 * @typedef {Object} Stop
 * @property {Number} id
 * @property {String} name
 * @property {Number} lat
 * @property {Number} lon
 */

/**
 * @typedef {Object} StopTime
 * @property {Number} stop_id
 * @property {String} departure
 * @property {Number} total_time
 */

/**
 * @typedef {Object} NodePath
 * @property {Number} from
 * @property {Number} to
 * @property {Number} time
 */

/**
 * @param {Function.<Array.<Stop>>} cb
 */
function fetchStops(cb) {
  get(base + 'stops', (body) => {
    cb(JSON.parse(body))
  })
}

/**
 * @param {String} stopId
 * @param {Function.<Array.<String>>} cb
 */
function fetchDepTimesForStop(stopId, cb) {
  get(base + `departures?stop_id=${stopId}`, (body) => {
    cb(JSON.parse(body))
  })
}

/**
 * @param {String} stopId
 * @param {String} maxTime
 * @param {String} departureTime
 * @param {Function.<Array.<StopTime>>} cb
 */
function fetchHeatmap(stopId, maxTime, departureTime, cb) {
  get(base + `heatmap?stop_id=${stopId}&from_time=${departureTime}&max_time_seconds=${maxTime}`, (body) => {
    cb(JSON.parse(body))
  })
}

/**
 * @param {String} stopId
 * @param {String} maxTime
 * @param {String} departureTime
 * @param {Function.<Array.<NodePath>>} cb
 */
function fetchPaths(stopId, maxTime, departureTime, cb) {
  get(base + `paths?stop_id=${stopId}&from_time=${departureTime}&max_time_seconds=${maxTime}`, (body) => {
    cb(JSON.parse(body))
  })
}
