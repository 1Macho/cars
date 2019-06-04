package main

import (
  "github.com/1Macho/geometry"
  "github.com/ojrac/opensimplex-go"
  "math/rand"
  "time"
)

func BuildSimplexTrack (radius float64, width float64, amplification float64, zoom float64, segments int) Loop {
  degreesPerSegment := 360 / float64(segments)
  innerPoints := []geometry.Point{}
  outerPoints := []geometry.Point{}
  checkPointsPoints := []geometry.Point{}
  currentAngle := geometry.Angle{0}
  noiseGen := opensimplex.New(time.Now().UnixNano())
  shouldInvert := false
  if (rand.Float64() > 0.5) {
    shouldInvert = true
  }
  var finishLine geometry.Segment
  var startPosition geometry.Point
  for i := 0; i < segments; i++ {
    basePoint := geometry.Point{0,0}.Translate(radius, currentAngle)
    noiseValue :=radius + noiseGen.Eval2(basePoint.X / zoom, basePoint.Y / zoom) * amplification
    innerPoint := geometry.Point{0,0}.Translate(noiseValue, currentAngle)
    outerPoint := geometry.Point{0,0}.Translate(noiseValue + width, currentAngle)
    if (shouldInvert) {
      outerPoint.X = outerPoint.X * -1
      innerPoint.X = innerPoint.X * -1
    }
    innerPoints = append(innerPoints, innerPoint)
    outerPoints = append(outerPoints, outerPoint)
    if (i == 0) {
      midPoint := noiseValue + width / 2
      startPosition = geometry.Point{0,0}.Translate(midPoint, currentAngle)
      if (shouldInvert) {
        startPosition.X = startPosition.X * -1
      }
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
  return Loop{walls, startPosition, finishLine, checkPoints, geometry.Angle{90}}
}

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
  return Loop{walls, startPosition, finishLine, checkPoints, geometry.Angle{90}}
}

func BuildRandomizedCircularTrack (variation float64, innerRadius float64, outerRadius float64, segments int) Loop {
  degreesPerSegment := 360 / float64(segments)
  innerPoints := []geometry.Point{}
  outerPoints := []geometry.Point{}
  checkPointsPoints := []geometry.Point{}
  shouldInvert := false
  if (rand.Float64() > 0.5) {
    shouldInvert = true
  }
  currentAngle := geometry.Angle{0}
  var finishLine geometry.Segment
  var startPosition geometry.Point
  for i := 0; i < segments; i++ {
    displacement := (rand.Float64() * variation * 2) - variation
    newInnerRadius := displacement + innerRadius
    newOuterRadius := displacement + outerRadius
    innerPoint := geometry.Point{0,0}.Translate(newInnerRadius, currentAngle)
    outerPoint := geometry.Point{0,0}.Translate(newOuterRadius, currentAngle)
    if (shouldInvert) {
      outerPoint.X = outerPoint.X * -1
      innerPoint.X = innerPoint.X * -1
    }
    innerPoint.Y = innerPoint.Y * 2
    outerPoint.Y = outerPoint.Y * 2
    innerPoints = append(innerPoints, innerPoint)
    outerPoints = append(outerPoints, outerPoint)
    if (i == 0) {
      midPoint := (newInnerRadius + newOuterRadius) / 2
      startPosition = geometry.Point{0,0}.Translate(midPoint, currentAngle)
      if (shouldInvert) {
        startPosition.X = startPosition.X * -1
      }
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
  return Loop{walls, startPosition, finishLine, checkPoints, geometry.Angle{90}}
}
