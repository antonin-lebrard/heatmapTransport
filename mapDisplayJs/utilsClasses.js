'use strict'

function rdrgb() {
  let str = 'rgb('
  str += Math.floor(Math.random() * 255).toString() + ','
  str += Math.floor(Math.random() * 255).toString() + ','
  str += Math.floor(Math.random() * 255).toString() + ')'
  return str
}

class MapPolygons {

  constructor(map, polygons) {
    this.map = map
    this.polygons = polygons.map(poly => {
      poly.pop()
      return Array.from(poly, (el => [el.x, el.y]))
    })
    this.drawPolygons()
  }

  drawPolygons() {
    for (let i = 0; i < this.polygons.length; i++) {
      let polygonL = L.polygon(this.polygons[i], {
        color: rd(),
        weight: 10,
        fill: false,
        opacity: 1,
        fillOpacity: 1
      })
      polygonL.addTo(this.map)
    }
  }

}
