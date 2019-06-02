package main

import (
  "github.com/1Macho/geometry"
)

func BuildCircularTrack (innerRadius float64, outerRadius float64, segments int) Loop {
  degreesPerSegment := 360 / float64(segments)
  innerPoints := []geometry.Point{}
  outerPoints := []geometry.Point{}
  checkPointsPoints := []geometry.Point{}
  currentAngle := geometry.Angle{0}
  var finishLine geometry.Segment
  var startPosition geometry.Point
  for i := 0; i < segments; i++ {
    innerPoint := geometry.Point{0,0}.Translate(innerRadius, currentAngle)
    outerPoint := geometry.Point{0,0}.Translate(outerRadius, currentAngle)
    innerPoints = append(innerPoints, innerPoint)
    outerPoints = append(outerPoints, outerPoint)
    if (i == 0) {
      midPoint := (innerRadius + outerRadius) / 2
      startPosition = geometry.Point{0,0}.Translate(midPoint, currentAngle)
      finishLine = geometry.SegmentFromPoints(innerPoint, outerPoint)
    } else {
      checkPointsPoints = append(checkPointsPoints, innerPoint, outerPoint)
    }
    currentAngle = currentAngle.Add(degreesPerSegment)
  }
  innerShape := geometry.ShapeFromPoints(innerPoints)
  outerShape := geometry.ShapeFromPoints(outerPoints)
  walls := append(innerShape, outerShape...)
  checkPoints := geometry.IntermitentShapeFromPoints(checkPointsPoints)
  return Loop{walls, startPosition, finishLine, checkPoints}
}
