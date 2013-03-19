
// Copyright (c) Harri Rautila, 2012,2013

// This file is part of github.com/hrautila/matops package. It is free software,
// distributed under the terms of GNU Lesser General Public License Version 3, or
// any later version. See the COPYING tile included in this archive.

package matops

import (
    "github.com/hrautila/matrix"
    "github.com/hrautila/matops/calgo"
    "errors"
    "math"
    //"fmt"
)

type Flags int
const (
    TRANSA = calgo.TRANSA
    TRANSB = calgo.TRANSB
    LOWER  = calgo.LOWER
    UPPER  = calgo.UPPER
    LEFT   = calgo.LEFT
    RIGHT  = calgo.RIGHT
    UNIT   = calgo.UNIT
    TRANS  = calgo.TRANS
    NOTRANS = calgo.NOTRANS
    NONE   = calgo.NOTRANS
)
    
// blocking parameter size for DOT based algorithms
var vpLen int = 196
var nB int = 68
var mB int = 68

// Number of parallel workers to use.
var nWorker int = 1

// problems small than this do not benefit from parallelism
var limitOne int64 = 200*200*200

// Set blocking parameters for blocked algorithms. Hard limits are set
// by actual C-implementation in matops.calgo package.
// Parameter nb defines column direction block size, mb defines row direction
// block size and vplen defines viewport length for matrix-matrix multiplication.
func BlockingParams(vplen, nb, mb int) {
    vpLen = vplen
    nB = nb
    mB = mb
}

func NumWorkers(newWorkers int) int {
    oldWorkers := nWorker
    nWorker = newWorkers
    return oldWorkers
}

func blockIndex4(i, r, sz int) int {
    if (i == r) {
        return sz;
    }
    return i*sz/r - ((i*sz/r) & 0x3);
}

func blockIndex2(i, r, sz int) int {
    if (i == r) {
        return sz;
    }
    return i*sz/r - ((i*sz/r) & 0x1);
}

func isSquared(num int) (int, bool) {
    nsqrt := int(math.Sqrt(float64(num)))
    issquared := nsqrt*nsqrt == num
    return nsqrt, issquared
}


func divideWork(rows, cols, workers int) (colWorkers int, rowWorkers int) {
    colWorkers = 0
    rowWorkers = 0
    nwsqrt := int(math.Sqrt(float64(workers)))
    issquare := nwsqrt*nwsqrt == workers
    if workers == 2 || (workers & 0x1) != 0 {
        // odd number of workers
        if cols > rows {
            colWorkers = workers
            rowWorkers = 1
        } else {
            rowWorkers = workers
            colWorkers = 1
        }
    } else if issquare {
        // square number 
        colWorkers = nwsqrt
        rowWorkers = nwsqrt
    } else {
        // even number of workers
        if cols > rows {
            rowWorkers = 2
            colWorkers = workers/2
        } else {
            colWorkers = 2
            rowWorkers = workers/2
        }
    }    
    //fmt.Printf("divideWork: c=%d, r=%d\n", colWorkers, rowWorkers)
    return
}

type task func(int, int, int, int, chan int)

func scheduleWork(colworks, rowworks, cols, rows int, worker task) {
    ntask := colworks*rowworks
    ch := make(chan int, ntask)
    for k := 0; k < colworks; k++ {
        colstart := blockIndex4(k, colworks, cols)
        colend   := blockIndex4(k+1, colworks, cols)
        for l := 0; l < rowworks; l++ {
            rowstart := blockIndex4(l, rowworks, rows)
            rowend   := blockIndex4(l+1, rowworks, rows)
            //fmt.Printf("schedule: S=%d, L=%d, R=%d, E=%d\n", colstart, colend, rowstart, rowend)
            go worker(colstart, colend, rowstart, rowend, ch)
        }
    }
    nready := 0
    for nready < ntask {
        nready += <- ch
    }
}

// Generic matrix-matrix multpily. (blas.GEMM). Calculates
//   C = beta*C + alpha*A*B     (default)
//   C = beta*C + alpha*A.T*B   flags&TRANSA
//   C = beta*C + alpha*A*B.T   flags&TRANSB
//   C = beta*C + alpha*A.T*B.T flags&(TRANSA|TRANSB)
//
// C is M*N, A is M*P or P*M if flags&TRANSA. B is P*N or N*P if flags&TRANSB.
//
func Mult(C, A, B *matrix.FloatMatrix, alpha, beta float64, flags Flags) error {
    if A.Cols() != B.Rows() {
        return errors.New("A.cols != B.rows: size mismatch")
    }
    psize := int64(C.NumElements()*A.Cols())
    Ar := A.FloatArray()
    ldA := A.LeadingIndex()
    Br := B.FloatArray()
    ldB := B.LeadingIndex()
    Cr := C.FloatArray()
    ldC := C.LeadingIndex()

    if nWorker <= 1 || psize <= limitOne {
        calgo.DMult(Cr, Ar, Br, alpha, beta, calgo.Flags(flags), ldC, ldA, ldB, B.Rows(),
            0, C.Cols(), 0, C.Rows(),
            vpLen, nB, mB)
        return nil
    } 
    // here we have more than one worker available
    worker := func(cstart, cend, rstart, rend int, ready chan int) {
        calgo.DMult(Cr, Ar, Br, alpha, beta, calgo.Flags(flags), ldC, ldA, ldB, B.Rows(),
            cstart, cend, rstart, rend, vpLen, nB, mB)
        ready <- 1
    }
    colworks, rowworks := divideWork(C.Rows(), C.Cols(), nWorker)
    scheduleWork(colworks, rowworks, C.Cols(), C.Rows(), worker)
    return nil
}

// Symmetric matrix multiply. (blas.SYMM)
//   C = beta*C + alpha*A*B     (default)
//   C = beta*C + alpha*A.T*B   flags&TRANSA
//   C = beta*C + alpha*A*B.T   flags&TRANSB
//   C = beta*C + alpha*A.T*B.T flags&(TRANSA|TRANSB)
//
// C is N*P, A is N*N symmetric matrix. B is N*P or P*N if flags&TRANSB.
//
func MultSym(C, A, B *matrix.FloatMatrix, alpha, beta float64, flags Flags) error {

    if A.Rows() != A.Cols() {
        return errors.New("A matrix not square matrix.");
    }
    if A.Cols() != B.Rows() {
        return errors.New("A.cols != B.rows: size mismatch")
    }
    psize := int64(C.NumElements())*int64(A.Cols())
    Ar := A.FloatArray()
    ldA := A.LeadingIndex()
    Br := B.FloatArray()
    ldB := B.LeadingIndex()
    Cr := C.FloatArray()
    ldC := C.LeadingIndex()

    if nWorker <= 1 || psize <= limitOne {
        calgo.DMultSymm(Cr, Ar, Br, alpha, beta, calgo.Flags(flags), ldC, ldA, ldB,
            A.Cols(),  0, C.Cols(), 0, C.Rows(), vpLen, nB, mB)
        return nil
    } 
    // here we have more than one worker available
    worker := func(cstart, cend, rstart, rend int, ready chan int) {
        calgo.DMultSymm(Cr, Ar, Br, alpha, beta, calgo.Flags(flags), ldC, ldA, ldB,
            A.Cols(), cstart, cend, rstart, rend, vpLen, nB, mB)
        ready <- 1
    }
    colworks, rowworks := divideWork(C.Rows(), C.Cols(), nWorker)
    scheduleWork(colworks, rowworks, C.Cols(), C.Rows(), worker)
    return nil
}

// Triangular matrix multiply. (blas.TRMM) 
//    B = alpha*A*B    if flags&LEFT 
//    B = alpha*A.T*B  if flags&(LEFT|TRANSA) 
//    B = alpha*B*A    if flags&RIGHT
//    B = alpha*B*A.T  if flags&(RIGHT|TRANSA)
//
// Matrix A is N*N triangular defined with flags bits as follow
//  LOWER       non-unit lower triangular
//  LOWER|UNIT  unit lower triangular, A diagonal not used
//  UPPER       non-unit upper triangular
//  UPPER|UNIT  unit upper triangular, A diagonal not used
// 
// Matrix B is N*P if flags&LEFT or P*N if flags&RIGHT.
//
func MultTrm(B, A *matrix.FloatMatrix, alpha float64, flags Flags) error {

    if flags&LEFT != 0 && B.Rows() != A.Cols() {
        return errors.New("A.cols != B.rows: size mismatch")
    } else if flags&RIGHT != 0 && B.Cols() != A.Rows() {
        return errors.New("A.cols != B.rows: size mismatch")
    }

    Ar := A.FloatArray()
    ldA := A.LeadingIndex()
    Br := B.FloatArray()
    ldB := B.LeadingIndex()
    
    E := B.Cols()
    if flags & RIGHT != 0 {
        E = B.Rows()
    }
    // if more workers available can divide to tasks by B columns if flags&LEFT or by
    // B rows if flags&RIGHT.
    calgo.DTrmmBlk(Br, Ar, alpha, calgo.Flags(flags), ldB, ldA, A.Cols(), 0, E, nB)
    return nil
}

// Solve multiple right sides. If flags&UNIT then A diagonal is assumed to
// to unit and is not referenced. (blas.TRSM)
//      alpha*B = A.-1*B if flags&LEFT
//      alpha*B = A.-T*B if flags&(LEFT|TRANS)
//      alpha*B = B*A.-1 if flags&RIGHT
//      alpha*B = B*A.-T if flags&(RIGHT|TRANS)
//
// Matrix A is N*N triangular matrix defined with flags bits as follow
//  LOWER       non-unit lower triangular
//  LOWER|UNIT  unit lower triangular
//  UPPER       non-unit upper triangular
//  UPPER|UNIT  unit upper triangular
//
// Matrix B is N*P if flags&LEFT or P*N if flags&RIGHT.
//
func Solve(B, A *matrix.FloatMatrix, alpha float64, flags Flags) error {

    if flags&LEFT != 0 && B.Rows() != A.Cols() {
        return errors.New("A.cols != B.rows: size mismatch")
    } else if flags&RIGHT != 0 && B.Cols() != A.Rows() {
        return errors.New("A.cols != B.rows: size mismatch")
    }

    Ar := A.FloatArray()
    ldA := A.LeadingIndex()
    Br := B.FloatArray()
    ldB := B.LeadingIndex()
    
    E := B.Cols()
    if flags & RIGHT != 0 {
        E = B.Rows()
    }
    // if more workers available can divide to tasks by B columns if flags&LEFT or by
    // B rows if flags&RIGHT.
    calgo.DSolveBlk(Br, Ar, alpha, calgo.Flags(flags), ldB, ldA, A.Cols(), 0, E, nB)
    return nil
}

// Rank update for symmetric lower or upper matrix (blas.SYRK)
//      C = beta*C + alpha*A*A.T + alpha*A.T*A
func RankUpdateSym(C, A *matrix.FloatMatrix, alpha, beta float64, flags Flags) error {
    if C.Rows() != C.Cols() {
        return errors.New("C not a square matrix")
    }
    Ar := A.FloatArray()
    ldA := A.LeadingIndex()
    Cr := C.FloatArray()
    ldC := C.LeadingIndex()
    S := 0
    E := C.Rows()
    // if more workers available C can be divided to blocks [S:E, S:E] along diagonal
    // and updated in separate tasks. 
    calgo.DSymmRankBlk(Cr, Ar, alpha, beta, calgo.Flags(flags), ldC, ldA, A.Cols(), S, E,
        vpLen, nB)
    return nil
}

// Rank 2 update for symmetric lower or upper matrix. (blas.SYR2K)
//      C = beta*C + alpha*A*B.T + alpha*B*A.T 
//      C = beta*C + alpha*A.T*B + alpha*B.T*A  if flags&TRANS
// matrix C
//   lower triangular if flags&LOWER
//   upper triangular if flags&UPPER
func RankUpdate2Sym(C, A, B *matrix.FloatMatrix, alpha, beta float64, flags Flags) error {
    Ar := A.FloatArray()
    ldA := A.LeadingIndex()
    Br := B.FloatArray()
    ldB := B.LeadingIndex()
    Cr := C.FloatArray()
    ldC := C.LeadingIndex()
    S := 0
    E := C.Rows()
    // if more workers available C can be divided to blocks [S:E, S:E] along diagonal
    // and updated in separate tasks. 
    calgo.DSymmRank2Blk(Cr, Ar, Br, alpha, beta, calgo.Flags(flags), ldC, ldA, ldB,
        A.Cols(), S, E, vpLen, nB)
    return nil
}


// Local Variables:
// tab-width: 4
// indent-tabs-mode: nil
// End:
