//Madison Stulir
//Gray-Scott HW
//September 12, 2022

package main

//SimulateGrayScott is a function that takes an initial board,
//a number of generations, a feed and kill rate, diffusion rates for predator
//and prey and returns the gray scott simulation of predator and prey interactions for the number of generations input
//Input: a board, number of generations integer, feed and kill rates of the reaction as well as diffusion rates of the prey and predator as floats, and a kernel for the diffusion
//Output: an array of boards, numGens long
func SimulateGrayScott(initialBoard Board, numGens int, feedRate, killRate, preyDiffusionRate, predatorDiffusionRate float64, kernel [3][3]float64) []Board {
  //make an array of numGens+1 Boards using function initializeBoard

  numRows:=CountRowsMatrix(initialBoard)
  numCols:=CountColsMatrix(initialBoard)
  var boards []Board
  boards = make([]Board, numGens+1)
  boards[0]=initialBoard
  for i:=1;i<numGens+1;i++{
    boards[i]=initializeBoard(numRows,numCols)
  }
  //Multiply the kernel by the predator and prey diffusion rates for use in updating the boards each generation
  preyKernel:=multiplyKernel(kernel,preyDiffusionRate)
  predatorKernel:=multiplyKernel(kernel,predatorDiffusionRate)
  //loop through the number of generations and update the board
  for i:=1;i<numGens+1;i++ {
    boards[i]=UpdateBoard(boards[i-1],feedRate,killRate,preyKernel,predatorKernel)
    //fmt.Println("generation",i)
  }
  //return an array of boards
  return boards
}

//UpdateBoard takes in the current board and updates it according to the grayscott model for 1 generation
//Input: A board, feed and kill rates as float64 and diffusion kernels for both the predator and prey
//Output: A board updated by 1 generation
func UpdateBoard(currentBoard Board,feedRate,killRate float64, preyKernel,predatorKernel [][]float64) Board{
  numRows:=CountRowsMatrix(currentBoard)
  numCols:=CountColsMatrix(currentBoard)
  newBoard:=initializeBoard(numRows,numCols)
  //loop through the rows and columns and update the value of the cell for the next generation at each position
  for i:=0;i<numRows;i++ {
    for j:=0;j<numCols;j++ {
      newBoard[i][j] = UpdateCell(currentBoard, i, j, feedRate, killRate, preyKernel,predatorKernel)
    }
  }
  return newBoard
}

//UpdateCell takes in a current board and coordinates row and col in the board and updates its value for the next generation using feed and kill rates, as well as diffusion kernels for predator and prey
//Input: A board, feed and kill rates as float64 and diffusion kernels for both the predator and prey
//Ouput: A Cell value updated for a generation
func UpdateCell(currentBoard Board,row,col int,feedRate, killRate float64, preyKernel,predatorKernel [][]float64) Cell {
  currentCell:=currentBoard[row][col]
  //determine the change in the value from both diffusion and reactions
  diffusionValues:=ChangeDueToDiffusion(currentBoard, row, col, preyKernel,predatorKernel)
  reactionValues:=ChangeDueToReactions(currentCell, feedRate, killRate)
  //return the sum of the current value of the cell and the calculated changes due to diffusion and reactions
  return SumCells(currentCell, diffusionValues, reactionValues)
}

//SumCells takes in a variable number of cells and sums the values at indexes 0 and 1 to get a single cell as return
//Input: a variable number of cells
//Output: the sum of the cells at indexes 0 and 1 as a single cell
func SumCells(cells ...Cell) Cell {
s :=[2]float64{0,0}
	for _, val := range cells {
		s[0] += val[0]
    s[1]+=val[1]
	}
	return s
}

//ChangeDueToReactions updates the value of the current cell by the provided feed and kill rates according to the gray scott model
//Input: a cell, and feed and kill rates as floats
//Ouput: a new cell updated for a generation
func ChangeDueToReactions(currentCell Cell, feedRate, killRate float64) Cell {
  new:=[2]float64{}
  //update predator and prey values separately in a Cell
  new[0]=feedRate*(1-currentCell[0])-currentCell[0]*currentCell[1]*currentCell[1]
  new[1]=-1*killRate*currentCell[1]+currentCell[0]*currentCell[1]*currentCell[1]
  return new
}

//ChangeDueToDiffusion takes in a current board  and its row and column indexes and updates its value for 1 generation by appling prey and predator kernels to give its updated value
//Input: a current board, row and column indexes, and 2 kernels representing the diffusion for the prey at index 0 and predator at index 1
//Ouput: an updated cell for the input row and col of the current board
func ChangeDueToDiffusion(currentBoard Board, row, col int, preyKernel,predatorKernel [][]float64) Cell {
  //input only the prey current values currentBoard[0] for applying its kernel
  prey:=ApplyKernel(currentBoard,row,col,0,preyKernel)
  //input only the predator current values currentBoard[0] for applying its kernel
  predator:=ApplyKernel(currentBoard,row,col, 1,predatorKernel)
  result:=[2]float64{prey,predator}
  return result
}

//TestInfield checks whether the x and y values input are within the board and returns zero if outside the board and the value at that position if inside the board
//Input: a board, locx and locy coordinates and predvprey indicating which slice of the board we would like to return the value from
//Output: the float64 value at the board position indicated (0 if not InField)
func TestInfield(board Board,locx,locy,predvprey int) float64 {
  if InField(board,locx,locy)==true{
    value:=board[locx][locy][predvprey]
    return float64(value)
  } else { // not on board
    value:=0.0
    return value
  }
}
//InField takes in a board and row col corrdinates and returns true if the coordinates are within the board and false if not
//Input: a board and row and col integers
//Output: a bool true if in board, false if out of board
func InField(board Board,row,col int) bool {
  //check positions that are out of field
  if row<0 || col<0 || row>=CountRowsMatrix(board) || col>=CountColsMatrix(board){
    return false
  }
  //if we survive to here we know we are on the board
  return true
}
//ApplyKernel takes a board, r and c indexes, predvprey indicating which slice of board to update and a kernel and applies the kernel to the indicated r and c index in the board
//Input:a board, r and c indexes, predvprey integers and a kernel [][]float64
//Output: a new float64 value for the position in the board updated by the kernel
func ApplyKernel(board Board, r,c,predvprey int, kernel [][]float64) float64 {
  //define the positions around r and c, checking if they are in field and if not assigning them as zero
  center:=board[r][c][predvprey]
  northwest:=TestInfield(board,r-1,c-1,predvprey)
  north:=TestInfield(board,r-1,c,predvprey)
  northeast:=TestInfield(board,r-1,c+1,predvprey)
  east:=TestInfield(board,r,c+1,predvprey)
  southeast:=TestInfield(board,r+1,c+1,predvprey)
  south:=TestInfield(board,r+1,c,predvprey)
  southwest:=TestInfield(board,r+1,c-1,predvprey)
  west:=TestInfield(board,r,c-1,predvprey)

  //define the coordinates of the kernel for ease of final update
  kcenter:=kernel[1][1]
  knorthwest:=kernel[0][0]
  knorth:=kernel[0][1]
  knortheast:=kernel[0][2]
  keast:=kernel[1][2]
  ksoutheast:=kernel[2][2]
  ksouth:=kernel[2][1]
  ksouthwest:=kernel[2][0]
  kwest:=kernel[1][0]
  //"values multiplied in order", center,kcenter,northwest,knorthwest,north,knorth,northeast,knortheast,east,keast,southeast,ksoutheast,south,ksouth,southwest,ksouthwest,west,kwest
  new:=center*kcenter+northwest*knorthwest+north*knorth+northeast*knortheast+east*keast+southeast*ksoutheast+south*ksouth+southwest*ksouthwest+west*kwest
  return new
}

//multiplyKernel takes an input kernel and multiplies all of its values by a diffusion rate float64
//Input: a kernel [3][3]float64 and a diffusion rate the multiply the kernel by
//Output: a kernel with updated values [][]float64
func multiplyKernel(kernel [3][3]float64, diffusionRate float64) [][]float64 {
  newKernel:=initializekernelBoard(3,3)
  //loop through all positions of kernel and update the value by multiplying it by the diffusion rate
  for i:=0;i<3;i++ {
    for j:=0;j<3;j++ {
      newKernel[i][j]=diffusionRate*kernel[i][j]
    }
  }
  return newKernel
}

//initializekernelBoard takes in a number of rows and columns as integers and returns a matrix of that size of float64
//Input: numRows and numCols integers representing the size of the desired kernel
//Output: a kernel of type [numRows][numCols]float64
func initializekernelBoard(numRows,numCols int) [][]float64 {
  //make a 2-D slice
  //default vals=false
  var board [][]float64
  board=make([]([]float64),numRows)
  //now we need to make the rows the correct length (cols)
  for r:=range board {
    board[r]=make([]float64,numCols)
  }
  return board
}
//initializeBoard takes in a number of rows and columns as integers and returns a Board of that size
//Input: numRows and numCols integers representing the size of the desired Board
//Output: a board of the numRows,numCols size
func initializeBoard(numRows,numCols int) Board {
  var board Board
  board = make(Board, numRows)
  for z := range board { board[z] = make([]Cell, numCols)}
  return board
}

//CountRowsMatrix takes in a Board and returns an integer of the number of rows in the board
//Input: a board
//Output: an integer of the number of rows in the board
func CountRowsMatrix(board Board) int {
  return len(board)
}
//CountRowsMatrix takes in a Board and returns an integer of the number of columns in the board
//Input: a board
//Output: an integer of the number of columns in the board
func CountColsMatrix(board Board) int {
  // assume we have a rectangular board
  if len(board)==0{
    panic("Error: empty board given to CountCols")
  }
  // give number of elements in the 0th row
  return len(board[0])
}
