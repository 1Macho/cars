package main

import (
  "github.com/veandco/go-sdl2/sdl"
  "math/rand"
  "time"
)

func main() {
  if err := sdl.Init(sdl.INIT_EVERYTHING); err != nil {
    panic(err)
  }
  defer sdl.Quit()

  window, err := sdl.CreateWindow("Cars", sdl.WINDOWPOS_UNDEFINED, sdl.WINDOWPOS_UNDEFINED,
		1200,1000, sdl.WINDOW_SHOWN)
  if err != nil {
    panic(err)
  }
  defer window.Destroy()

  renderer, err := sdl.CreateRenderer(window, -1, sdl.RENDERER_ACCELERATED)
	if err != nil {
		panic(err)
	}
	defer renderer.Destroy()

	renderer.Clear()

  //testAngle := geometry.Angle{90}
  //println(testAngle.Slope())
  //testLine := geometry.Line{geometry.Angle{90+45},geometry.Point{0,1}}
  //println(testLine.TrueOrigin().X)
  //println(testLine.TrueOrigin().Y)
  //testRay := geometry.Ray{geometry.Line{geometry.Angle{47},geometry.Point{0,0}}}
  //println(geometry.InterceptRayLine(testRay,testLine).X)
  //println(geometry.InterceptRayLine(testRay,testLine).Y)

  rand.Seed(time.Now().UnixNano())

  testLoop := BuildCircularTrack(1000, 1400, 16)

  testSimulation := CreateSimulation(testLoop, 20)

  running := true
  for running {
    testSimulation.Tick()
    testSimulation.Draw(renderer)
//    println(testCar.Drivable.Direction.Value)
    //testDrivable.Draw(renderer, geometry.Point{0,0})
    renderer.Present()
    //testParticle.Tick()
    for event := sdl.PollEvent(); event != nil; event = sdl.PollEvent() {
      switch event.(type) {
      case *sdl.QuitEvent:
        println("Quit")
        running = false
        break
      }
    }
    time.Sleep(1000000000/60)
  }
}
