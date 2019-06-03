package main

import (
  "github.com/1Macho/geometry"
  "github.com/veandco/go-sdl2/sdl"
)

type Loop struct {
  Walls []geometry.Segment
  Start geometry.Point
  FinishLine geometry.Segment
  CheckPoints []geometry.Segment
  StartAngle geometry.Angle
}

func DrawSegment (renderer *sdl.Renderer, offset geometry.Point, toDraw geometry.Segment) {
  a := toDraw.Line.Origin
  b := toDraw.EndPoint()
  renderer.DrawLine(int32(a.X + offset.X), int32(a.Y + offset.Y), int32(b.X + offset.X), int32(b.Y + offset.Y))
}

func DrawSegments (renderer *sdl.Renderer, offset geometry.Point, toDraw []geometry.Segment) {
  for _, target := range toDraw {
    DrawSegment(renderer, offset, target)
  }
}

func (l Loop) Draw (renderer *sdl.Renderer, offset geometry.Point) {
  renderer.SetDrawColor(255,255,255,255)
  DrawSegments(renderer, offset, l.Walls)
  renderer.SetDrawColor(0,0,120,255)
  DrawSegments(renderer, offset, l.CheckPoints)
  renderer.SetDrawColor(255,0,0,255)
  DrawSegment(renderer, offset, l.FinishLine)
}
