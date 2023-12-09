package main

import (
	"errors"
	"fmt"
	"os"
)

//screen interface
type screen interface {
  initialize(maxX, maxY int)
  getMaxXY() (maxX, maxY int)
  drawPixel(x, y int, c Color) (err error)
  getPixel(x, y int) (c Color, err error)
  clearScreen()
  screenShot(f string) (err error)
}

//geometry interface
type geometry interface {
  draw(scn *screen) (err error)
  shape() (s string)
}

//Point
type Point struct {
  x, y int
}

//rectangle
type Rectangle struct {
  ll, ur Point
  c int
}

//Triangle
type Triangle struct {
  pt0, pt1, pt2 Point
  c int
}

//CIRCLE
type Circle struct {
  cp Point
  r int
  c int
}

//display
type Display struct {
  maxX , maxY int
  matrix [][]Color
}

//basis for colors
type Color struct {
  r, g, b int
}

// display 
//sets up what we use for the program
var display Display
var cmap = [9]Color{{255, 0, 0}, {0, 255, 0}, {0, 0, 255}, {255, 255, 0}, {255, 164, 0}, {128, 0, 128}, {165, 42, 42}, {0, 0, 0}, {255, 255, 255}}

//-------------------------------------------SCREEN-------------------------------------------
//inititalize the screen
func (d Display) initialize(maxX, maxY int) {
  display.maxX = maxX
  display.maxY = maxY
  display.matrix = make([][]Color, maxX)
  for i := range display.matrix {
    display.matrix[i] = make([]Color, maxY)
  }
  display.clearScreen()
}

//get the x,y dimensions of the screen
func (display Display) getMaxXY() (maxX, maxY int) {
  return display.maxX, display.maxY
}

//helper for the color checker
func colorUnknown(c int) (result bool) {
  if c >= 9 {
    return false
  } else if c < 0 {
    return false
  }
  return true
}

//helper for the color checker (drawPixel as it takes in Color Struct)
func colorUnknown2(c Color) (result bool) {
  var checker = false
  for _, color := range cmap {
    if color == c {
      checker = true
    }
  }
  if checker {
    return true
  }
  return false
}

//out of bounds helper
//https://blog.logrocket.com/exploring-structs-interfaces-go/
func outOfBounds(p Point, scn screen) (result bool){
  x, y := scn.getMaxXY()
  if p.x >= x || p.x < 0 {
    return false
  } else if p.y >= y || p.y < 0 {
    return false
  }
  
  return true
}

//draw the pixel with color c at the location x,y
func (d Display) drawPixel(x, y int, c Color) (err error) {
  
  if !colorUnknown2(c) {
    return colorUnknownErr
  }
  
  newPoint := Point{x, y}
  if !outOfBounds(newPoint, &display) {
    return outOfBoundsErr
  }
  
  d.matrix[x][y] = c
  
  return nil
}

//Get the pixel color at location x,y
func (d Display) getPixel(x, y int) (c Color, err error) {
  newPoint := Point{x, y}
  if outOfBounds(newPoint, &display) == false {
    return Color{}, outOfBoundsErr
  }
  return display.matrix[x][y], nil
}

//clear the screen by setting each pixel to white
func (d Display) clearScreen() {
  for a := 0; a < display.maxX; a++ {
    for b := 0; b < display.maxY; b++ {
      display.matrix[a][b] = cmap[white];
    }
  }
}

//needs work
//dump the screen as a f.pmm
func (d Display) screenShot(f string) (err error) {
  fileMaker, err := os.Create(f + ".ppm")
  if (err != nil) {
    return errors.New("Does not work")
  }
  
  //header
  var height = fmt.Sprintf("%d", d.maxY)
  var width = fmt.Sprintf("%d", d.maxX)
  
  fileMaker.WriteString("P3\n")
  var depth = 0
  for a := 0; a < display.maxX; a++ {
    for b := 0; b < display.maxY; b++ {
      if depth < display.matrix[a][b].r {
        depth = display.matrix[a][b].r
      } else if depth < display.matrix[a][b].g {
        depth = display.matrix[a][b].g
      } else if depth < display.matrix[a][b].b {
        depth = display.matrix[a][b].b
      }
    }
    fileMaker.WriteString("\n")
  }
  var depthh = fmt.Sprintf("%d", depth)
  fileMaker.WriteString(width + " " + height + "\n" + depthh + "\n")
  
  
  //data
  //reference: PROJECT 2 UTILITY.cs
  for a := 0; a < display.maxX; a++ {
    for b := 0; b < display.maxY; b++ {
      var a1 = fmt.Sprintf("%d", display.matrix[a][b].r)
      var a2 = fmt.Sprintf("%d", display.matrix[a][b].g)
      var a3 = fmt.Sprintf("%d", display.matrix[a][b].b)
      fileMaker.WriteString(a1 + " " + a2 + " " + a3 + " ")
    }
    fileMaker.WriteString("\n")
  }
  return nil
}
//-------------------------------------------GEOMETRY-------------------------------------------
//  https://gabrielgambetta.com/computer-graphics-from-scratch/07-filled-triangles.html
func interpolate (l0, d0, l1, d1 int) (values []int) {
  a := float64(d1 - d0) / float64(l1 - l0)
  d  := float64(d0)

  count := l1-l0+1
  for ; count>0; count-- {
    values = append(values, int(d))
    d = d+a
  }
  return
}

//circle draw
//https://stackoverflow.com/questions/51626905/drawing-circles-with-two-radius-in-golang
func (circ Circle) draw(scn screen) (err error) {
  upY := Point{circ.cp.x, circ.cp.y + circ.r}
  downY := Point{circ.cp.x, circ.cp.y - circ.r}
  upX := Point{circ.cp.x + circ.r, circ.cp.y}
  downX := Point{circ.cp.x - circ.r, circ.cp.y}
  if !outOfBounds(upY, scn) || !outOfBounds(downY, scn) || !outOfBounds(downX, scn) || !outOfBounds(upX, scn){
    return outOfBoundsErr
  }

  if !colorUnknown(circ.c) {
    return colorUnknownErr
  }
  
  var baseX = circ.cp.x
  var baseY = circ.cp.y
  var testR = circ.r
  var theC = cmap[circ.c]
  for testR > -1 {
    var r = testR
    x, y, dx, dy := r - 1, 0, 1, 1
    err1 := dx - (r * 2)
    for x > y {
      scn.drawPixel(baseX + x, baseY + y, theC)
      scn.drawPixel(baseX + y, baseY + x, theC)
      scn.drawPixel(baseX - y, baseY + x, theC)
      scn.drawPixel(baseX - x, baseY + y, theC)
      scn.drawPixel(baseX - x, baseY - y, theC)
      scn.drawPixel(baseX - y, baseY - x, theC)
      scn.drawPixel(baseX + y, baseY - x, theC)
      scn.drawPixel(baseX + x, baseY - y, theC)
      if err1 <= 0 {
        y++
        err1 += dy
        dy += 2
      }
      if err1 > 0{
        x--
        dx += 2
        err1 += dx - (r*2)
      }
    }
    testR--
  }
  for a := baseX - baseX; a <= baseX + baseX; a++ {
    check := false
    for b := baseY - baseY; b <= baseY + baseY; b++ {
      result, _ := scn.getPixel(a, b)
      if result == theC && !check{
        check = true
        scn.drawPixel(a, b, theC)
      }
      if b != baseY + baseY {
        result2, _ := scn.getPixel(a, b + 1)
        if check && result2 == theC{
          scn.drawPixel(a, b, theC)
        }
      }
    }
  }
  for a := baseX - baseX; a <= baseX + baseX; a++ {
    check := false
    for b := baseY - baseY; b <= baseY + baseY; b++ {
      result, _ := scn.getPixel(a, b)
      if result == theC && !check{
        check = true
        scn.drawPixel(a, b, theC)
      }
      if b != baseY + baseY {
        result2, _ := scn.getPixel(a, b + 1)
        if check && result2 == theC{
          scn.drawPixel(a, b, theC)
        }
      }
    }
  }
  return nil
}

//rectangle draw
func (rect Rectangle) draw(scn screen) (err error) {
  
  if !outOfBounds(rect.ll, scn) || !outOfBounds(rect.ur, scn) {
    return outOfBoundsErr
  }
  if !colorUnknown(rect.c) {
    return colorUnknownErr
  }
  
  for a := rect.ll.x; a <= rect.ur.x; a++ {
    for b := rect.ll.y; b <= rect.ur.y; b++ {
      scn.drawPixel(a, b, cmap[rect.c])
    }
  }
  
  return nil
}

//triangle method
//  https://gabrielgambetta.com/computer-graphics-from-scratch/07-filled-triangles.html
func (tri Triangle) draw(scn screen) (err error) {
  if !outOfBounds(tri.pt0, scn) || !outOfBounds(tri.pt1,scn)  || !outOfBounds(tri.pt2,scn){
    return outOfBoundsErr
  }
  if !colorUnknown(tri.c) {
    return colorUnknownErr
  }

  y0 := tri.pt0.y
  y1 := tri.pt1.y
  y2 := tri.pt2.y

  // Sort the points so that y0 <= y1 <= y2
  if y1 < y0 { tri.pt1, tri.pt0 = tri.pt0, tri.pt1 }
  if y2 < y0 { tri.pt2, tri.pt0 = tri.pt0, tri.pt2 }
  if y2 < y1 { tri.pt2, tri.pt1 = tri.pt1, tri.pt2 }

  x0,y0,x1,y1,x2,y2 := tri.pt0.x, tri.pt0.y, tri.pt1.x, tri.pt1.y, tri.pt2.x, tri.pt2.y

  x01 := interpolate(y0, x0, y1, x1)
  x12 := interpolate(y1, x1, y2, x2)
  x02 := interpolate(y0, x0, y2, x2)

  // Concatenate the short sides

  x012 := append(x01[:len(x01)-1],  x12...)

  // Determine which is left and which is right
  var x_left, x_right []int
  m := len(x012) / 2
  if x02[m] < x012[m] {
    x_left = x02
    x_right = x012
  } else {
    x_left = x012
    x_right = x02
  }

  // Draw the horizontal segments
  for y := y0; y<= y2; y++  {
    for x := x_left[y - y0]; x <=x_right[y - y0]; x++ {
      scn.drawPixel(x, y, cmap[tri.c])
    }
  }
  return
}

//return the type of the object per shape type

func (rect Rectangle) shape() (s string) {
  return "Rectangle"
}

func (tri Triangle) shape() (s string) {
  return "Triangle"
}

func (circ Circle) shape() (s string) {
  return "Circle"
}
//--------------------------------------------------------------------------------------
//colors pre-defined
var red = 0
var green = 1
var blue = 2
var yellow = 3
var orange = 4
var purple = 5
var brown = 6
var black = 7
var white = 8

//predefined errors
var outOfBoundsErr = errors.New("geometry out of bounds")
var colorUnknownErr = errors.New("color unknown")

//-------------------------------------------MAIN-------------------------------------------
func main() {
  fmt.Println("starting ...")
  display.initialize(1024, 1024)
  
  rect :=  Rectangle{Point{100,300}, Point{600,900}, red}
  err := rect.draw(&display)
  if err != nil {
    fmt.Println("rect: ", err)
  }

  fmt.Println("rect is a ", rect.shape())

  rect2 := Rectangle{Point{0,0}, Point{100, 1024}, green}
  err = rect2.draw(&display)
  if err != nil {
    fmt.Println("rect2: ", err)
  }

  rect3 := Rectangle{Point{100,300}, Point{100, 1022}, 102}
  err = rect3.draw(&display)
  if err != nil {
    fmt.Println("rect3: ", err)
  }

  circ := Circle{Point{500,500}, 200, green}
  err = circ.draw(&display)
  if err != nil {
    fmt.Println("circ: ", err)
  }

  circ2 := Circle{Point{0,0}, 200, purple}
  err = circ2.draw(&display)
  if err != nil {
    fmt.Println("circ2: ", err)
  }

  circ3 := Circle{Point{500,500}, 200, 18}
  err = circ3.draw(&display)
  if err != nil {
    fmt.Println("circ3: ", err)
  }

  fmt.Println("circ is a ", circ.shape())

  tri := Triangle{Point{100, 100}, Point{600, 300},  Point{859,850}, yellow}
  err = tri.draw(&display)
  if err != nil {
    fmt.Println("tri: ", err)
  }

  tri2 := Triangle{Point{1, -1}, Point{7, 300},  Point{859,10000}, yellow}
  err = tri2.draw(&display)
  if err != nil {
    fmt.Println("tri2: ", err)
  }

  tri3 := Triangle{Point{100, 100}, Point{600, 300},  Point{859,850}, -1}
  err = tri3.draw(&display)
  if err != nil {
    fmt.Println("tri3: ", err)
  }

  fmt.Println("tri is a ", tri.shape())

  fmt.Println(display.maxX)
  fmt.Println(display.maxY)
  fmt.Println(display.matrix[0][0].b)
  display.drawPixel(0,0, cmap[7])
  fmt.Println(display.getPixel(0,0))

  display.screenShot("output")
  /*
  fmt.Println("Pick your palett size")
  xDim := 1
  yDim := 1
  fmt.Println("Pick your x Dimension")
  fmt.Scanln(&xDim)
  fmt.Println("Pick your y dimension")
  fmt.Scanln(&yDim)
  display.initialize(xDim, yDim)

  //the loop of the program, basically gets the programming running until exit clause
  var exit = false
  for (!exit) {
    var option = -1
    fmt.Println("\nHow will you like to run this program?")
    fmt.Println("1. Run Test Cases")
    fmt.Println("2. Make Shapes Manually")
    fmt.Println("3. Exit")
    fmt.Scanln(&option)
    
    //per each option
    if option == 2 {
      var shape = -1
      var color = -1
      fmt.Println("What shape do you want to make?")
      fmt.Println("1. Triangle\n2. Rectange\n3. Circle")
      fmt.Scanln(&shape)
      
      fmt.Println("Color do you want your shape?")
      fmt.Println("1. Red\n2. Green\n3. Blue\n4. Yellow\n5. Orange\n6. Purple\n7. Brown\n8. Black\n9. White")
      fmt.Scanln(&color)

      
    } else if option == 3 {
      exit = true
      fmt.Println("Exitting Program")
      display.screenShot("output")
      
    } else {
      rect :=  Rectangle{Point{100,300}, Point{600,900}, red}
      err := rect.draw(&display)
      if err != nil {
        fmt.Println("rect: ", err)
      }

      fmt.Println("rect is a ", rect.shape())

      rect2 := Rectangle{Point{0,0}, Point{100, 1024}, green}
      err = rect2.draw(&display)
      if err != nil {
        fmt.Println("rect2: ", err)
      }

      rect3 := Rectangle{Point{100,300}, Point{100, 1022}, 102}
      err = rect3.draw(&display)
      if err != nil {
        fmt.Println("rect3: ", err)
      }

      circ := Circle{Point{500,500}, 200, green}
      err = circ.draw(&display)
      if err != nil {
        fmt.Println("circ: ", err)
      }
      
      circ2 := Circle{Point{0,0}, 200, purple}
      err = circ2.draw(&display)
      if err != nil {
        fmt.Println("circ2: ", err)
      }

      circ3 := Circle{Point{500,500}, 200, 18}
      err = circ3.draw(&display)
      if err != nil {
        fmt.Println("circ3: ", err)
      }
      
      fmt.Println("circ is a ", circ.shape())

      tri := Triangle{Point{100, 100}, Point{600, 300},  Point{859,850}, yellow}
      err = tri.draw(&display)
      if err != nil {
        fmt.Println("tri: ", err)
      }

      tri2 := Triangle{Point{1, -1}, Point{7, 300},  Point{859,10000}, yellow}
      err = tri2.draw(&display)
      if err != nil {
        fmt.Println("tri2: ", err)
      }

      tri3 := Triangle{Point{100, 100}, Point{600, 300},  Point{859,850}, -1}
      err = tri3.draw(&display)
      if err != nil {
        fmt.Println("tri3: ", err)
      }
      
      fmt.Println("tri is a ", tri.shape())

      fmt.Println(display.maxX)
      fmt.Println(display.maxY)
      fmt.Println(display.matrix[0][0].b)
      display.drawPixel(0,0, cmap[7])
      fmt.Println(display.getPixel(0,0))

      display.screenShot("output")
    }
  }
  */
}
