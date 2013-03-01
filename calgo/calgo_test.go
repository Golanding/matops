
// Copyright (c) Harri Rautila, 2012,2013

// This file is part of github.com/hrautila/matops package. It is free software,
// distributed under the terms of GNU Lesser General Public License Version 3, or
// any later version. See the COPYING tile included in this archive.


package calgo

import (
    "github.com/hrautila/matrix"
    "github.com/hrautila/linalg/blas"
    "github.com/hrautila/linalg"
    "testing"
    "math/rand"
    "time"
)

const M = 8
const N = 8
const P = 8


func TestMakeData(t *testing.T) {
    rand.Seed(time.Now().UnixNano())
}

func _TestMultSmall(t *testing.T) {
    bM := 7
    bN := 7
    bP := 7
    /*
    Ddata := [][]float64{
        []float64{1.0, 1.0, 1.0, 1.0, 1.0, 1.0, 1.0},
        []float64{2.0, 2.0, 2.0, 2.0, 2.0, 2.0, 2.0},
        []float64{3.0, 3.0, 3.0, 3.0, 3.0, 3.0, 3.0},
        []float64{4.0, 4.0, 4.0, 4.0, 4.0, 4.0, 4.0},
        []float64{5.0, 5.0, 5.0, 5.0, 5.0, 5.0, 5.0},
        []float64{6.0, 6.0, 6.0, 6.0, 6.0, 6.0, 6.0},
        []float64{7.0, 7.0, 7.0, 7.0, 7.0, 7.0, 7.0}}
    D := matrix.FloatMatrixFromTable(Ddata, matrix.RowOrder)
     */
    //E := matrix.FloatMatrixFromTable(Ddata, matrix.RowOrder)
    D := matrix.FloatNormal(bM, bP)
    E := matrix.FloatNormal(bP, bN)
    //D := matrix.FloatWithValue(bM, bP, 1.0)
    //E := matrix.FloatWithValue(bP, bN, 1.0)
    //C0 := matrix.FloatZeros(bM, bN)
    //C1 := matrix.FloatZeros(bM, bN)
    C0 := matrix.FloatWithValue(bM, bN, 1.0)
    C1 := C0.Copy()

    Dr := D.FloatArray()
    Er := E.FloatArray()
    C1r := C1.FloatArray()

    blas.GemmFloat(D, E, C0, 1.0, 2.0)
    t.Logf("blas: C=D*E\n%v\n", C0)

    DMult2(C1r, Dr, Er, 1.0, 2.0, NOTRANS, bM, bM, bP, bP, 0,  bN, 0,  bM, 4, 4, 4)
    t.Logf("C0 == C1: %v\n", C0.AllClose(C1))
    t.Logf("C1: C1=D*E\n%v\n", C1)
}


func _TestMultBig(t *testing.T) {
    bM := 100*M + 3
    bN := 100*N + 3
    bP := 100*P + 3
    D := matrix.FloatNormal(bM, bP)
    E := matrix.FloatNormal(bP, bN)
    C0 := matrix.FloatZeros(bM, bN)
    C1 := matrix.FloatZeros(bM, bN)

    Dr := D.FloatArray()
    Er := E.FloatArray()
    C1r := C1.FloatArray()

    blas.GemmFloat(D, E, C0, 1.0, 1.0)
    //t.Logf("blas: C=D*E\n%v\n", C0)

    DMult(C1r, Dr, Er, 1.0, 1.0, NOTRANS, bM, bM, bP, bP, 0,  bN, 0,  bM, 32, 32, 32)
    t.Logf("C0 == C1: %v\n", C0.AllClose(C1))
}


func _TestMultTransASmall(t *testing.T) {
    bM := 7
    bN := 7
    bP := 7
    /*
    Ddata := [][]float64{
        []float64{1.0, 1.0, 1.0, 1.0, 1.0},
        []float64{2.0, 2.0, 2.0, 2.0, 2.0},
        []float64{3.0, 3.0, 3.0, 3.0, 3.0},
        []float64{4.0, 4.0, 4.0, 4.0, 4.0},
        []float64{5.0, 5.0, 5.0, 5.0, 5.0}}
    D := matrix.FloatMatrixFromTable(Ddata, matrix.RowOrder)
     */
    D := matrix.FloatNormal(bM, bP)
    E := matrix.FloatNormal(bP, bN)
    //D := matrix.FloatWithValue(bM, bP, 2.0)
    //E := matrix.FloatWithValue(bP, bN, 2.0)
    //C0 := matrix.FloatZeros(bM, bN)
    //C1 := matrix.FloatZeros(bM, bN)
    C0 := matrix.FloatWithValue(bM, bN, 0.0)
    C1 := C0.Copy()
    Dt := D.Transpose()

    Dr := Dt.FloatArray()
    Er := E.FloatArray()
    C1r := C1.FloatArray()
    //t.Logf("Dt:\n%v\n", Dt)
    //t.Logf("E:\n%v\n", E)
    blas.GemmFloat(Dt, E, C0, 1.0, 1.0, linalg.OptTransA)
    t.Logf("blas: C=D*E\n%v\n", C0)

    DMult2(C1r, Dr, Er, 1.0, 1.0, TRANSA, bM, bM, bP, bP, 0,  bN, 0,  bM, 4, 4, 4)
    t.Logf("C0 == C1: %v\n", C0.AllClose(C1))
    t.Logf("C1: C1=D*E\n%v\n", C1)
}


func _TestMultTransABig(t *testing.T) {
    bM := 100*M + 3
    bN := 100*N + 3
    bP := 100*P + 3
    D := matrix.FloatNormal(bM, bP)
    E := matrix.FloatNormal(bP, bN)
    C0 := matrix.FloatZeros(bM, bN)
    C1 := matrix.FloatZeros(bM, bN)
    Dt := D.Transpose()

    Dr := Dt.FloatArray()
    Er := E.FloatArray()
    C1r := C1.FloatArray()

    blas.GemmFloat(Dt, E, C0, 1.0, 1.0, linalg.OptTransA)

    DMult(C1r, Dr, Er, 1.0, 1.0, TRANSA, bM, bM, bP, bP, 0,  bN, 0,  bM, 32, 32, 32)
    t.Logf("C0 == C1: %v\n", C0.AllClose(C1))
}

func _TestMultTransBSmall(t *testing.T) {
    bM := 7
    bN := 7
    bP := 7
    D := matrix.FloatNormal(bM, bP)
    E := matrix.FloatNormal(bP, bN)
    //D := matrix.FloatWithValue(bM, bP, 2.0)
    //E := matrix.FloatWithValue(bP, bN, 1.0)
    //C0 := matrix.FloatZeros(bM, bN)
    //C1 := matrix.FloatZeros(bM, bN)
    C0 := matrix.FloatWithValue(bP, bN, 1.0)
    C1 := C0.Copy()
    Et := E.Transpose()

    Dr := D.FloatArray()
    Er := Et.FloatArray()
    C1r := C1.FloatArray()

    blas.GemmFloat(D, Et, C0, 1.0, 1.0, linalg.OptTransB)
    t.Logf("blas: C=D*E.T\n%v\n", C0)

    DMult(C1r, Dr, Er, 1.0, 1.0, TRANSB, bM, bM, bP, bP, 0,  bN, 0,  bM, 4, 4, 4)
    t.Logf("C0 == C1: %v\n", C0.AllClose(C1))
    t.Logf("C1: C1=D*E.T\n%v\n", C1)
}


func _TestMultTransBBig(t *testing.T) {
    bM := 100*M + 3
    bN := 100*N + 3
    bP := 100*P + 3
    D := matrix.FloatNormal(bM, bP)
    E := matrix.FloatNormal(bP, bN)
    C0 := matrix.FloatZeros(bM, bN)
    C1 := matrix.FloatZeros(bM, bN)
    Et := E.Transpose()

    Dr := D.FloatArray()
    Er := Et.FloatArray()
    C1r := C1.FloatArray()

    blas.GemmFloat(D, Et, C0, 1.0, 1.0, linalg.OptTransB)
    //t.Logf("blas: C=D*E\n%v\n", C0)

    DMult(C1r, Dr, Er, 1.0, 1.0, TRANSB, bM, bM, bP, bP, 0,  bN, 0,  bM, 32, 32, 32)
    t.Logf("C0 == C1: %v\n", C0.AllClose(C1))
}

func _TestMultTransABSmall(t *testing.T) {
    bM := 7
    bN := 7
    bP := 7
    D := matrix.FloatNormal(bM, bP)
    E := matrix.FloatNormal(bP, bN)
    //D := matrix.FloatWithValue(bM, bP, 2.0)
    //E := matrix.FloatWithValue(bP, bN, 1.0)
    C0 := matrix.FloatZeros(bM, bN)
    C1 := matrix.FloatZeros(bM, bN)
    Dt := D.Transpose()
    Et := E.Transpose()

    Dr := Dt.FloatArray()
    Er := Et.FloatArray()
    C1r := C1.FloatArray()

    blas.GemmFloat(Dt, Et, C0, 1.0, 1.0, linalg.OptTransA, linalg.OptTransB)
    t.Logf("blas: C=D.T*E.T\n%v\n", C0)

    DMult(C1r, Dr, Er, 1.0, 1.0, TRANSA|TRANSB, bM, bM, bP, bP, 0,  bN, 0,  bM, 4, 4, 4)
    t.Logf("C0 == C1: %v\n", C0.AllClose(C1))
    t.Logf("C1: C1=D.T*E.T\n%v\n", C1)
}


func _TestMultTransABBig(t *testing.T) {
    bM := 100*M
    bN := 100*N
    bP := 100*P
    D := matrix.FloatNormal(bM, bP)
    E := matrix.FloatNormal(bP, bN)
    C0 := matrix.FloatZeros(bM, bN)
    C1 := matrix.FloatZeros(bM, bN)
    Dt := D.Transpose()
    Et := E.Transpose()

    Dr := Dt.FloatArray()
    Er := Et.FloatArray()
    C1r := C1.FloatArray()

    blas.GemmFloat(Dt, Et, C0, 1.0, 1.0, linalg.OptTransA, linalg.OptTransB)

    DMult(C1r, Dr, Er, 1.0, 1.0, TRANSA|TRANSB, bM, bM, bP, bP, 0,  bN, 0,  bM, 32, 32, 32)
    t.Logf("C0 == C1: %v\n", C0.AllClose(C1))
}

func _TestMultMVSmall(t *testing.T) {
    bM := 5
    bN := 5
    A := matrix.FloatNormal(bM, bN)
    //X := matrix.FloatNormal(bN, 1)
    //A := matrix.FloatWithValue(bM, bN, 2.0)
    X := matrix.FloatVector([]float64{1.0, 2.0, 3.0, 4.0, 5.0})
    Y1 := matrix.FloatZeros(bM, 1)
    Y0 := matrix.FloatZeros(bM, 1)

    Ar := A.FloatArray()
    Xr := X.FloatArray()
    Y1r := Y1.FloatArray()

    blas.GemvFloat(A, X, Y0, 1.0, 1.0)
    t.Logf("blas: Y=A*X\n%v\n", Y0)

    DMultMV(Y1r, Ar, Xr, 1.0, 1.0, NOTRANS, 1, A.LeadingIndex(), 1, 0,  bN, 0,  bM, 4, 4)
    t.Logf("Y0 == Y1: %v\n", Y0.AllClose(Y1))
    t.Logf("Y1: Y1 = A*X\n%v\n", Y1)
}

func _TestMultMV(t *testing.T) {
    bM := 100*M
    bN := 100*N
    A := matrix.FloatNormal(bM, bN)
    X := matrix.FloatNormal(bN, 1)
    Y1 := matrix.FloatZeros(bM, 1)
    Y0 := matrix.FloatZeros(bM, 1)

    Ar := A.FloatArray()
    Xr := X.FloatArray()
    Y1r := Y1.FloatArray()

    blas.GemvFloat(A, X, Y0, 1.0, 1.0)

    DMultMV(Y1r, Ar, Xr, 1.0, 1.0, NOTRANS, 1, A.LeadingIndex(), 1, 0,  bN, 0,  bM, 32, 32)
    t.Logf("Y0 == Y1: %v\n", Y0.AllClose(Y1))
    if ! Y0.AllClose(Y1) {
        y0 := Y0.SubMatrix(0, 0, 2, 1)
        y1 := Y1.SubMatrix(0, 0, 2, 1)
        t.Logf("y0=\n%v\n", y0)
        t.Logf("y1=\n%v\n", y1)
    }
}

func _TestMultMVTransASmall(t *testing.T) {
    bM := 10
    bN := 10
    /*
    Adata := [][]float64{
     []float64{1.0, 1.0, 1.0, 1.0, 1.0},
     []float64{2.0, 2.0, 2.0, 2.0, 2.0},
     []float64{3.0, 3.0, 3.0, 3.0, 3.0},
     []float64{4.0, 4.0, 4.0, 4.0, 4.0},
     []float64{5.0, 5.0, 5.0, 5.0, 5.0}}
    A := matrix.FloatMatrixFromTable(Adata)
     */
    A := matrix.FloatNormal(bN, bM)
    //X := matrix.FloatNormal(bN, 1)
    X := matrix.FloatWithValue(bN, 1, 1.0)
    //A := matrix.FloatWithValue(bM, bN, 2.0)
    //X := matrix.FloatVector([]float64{1.0, 2.0, 3.0, 4.0, 5.0})
    //At := A.Transpose()
    Y1 := matrix.FloatZeros(bM, 1)
    Y0 := matrix.FloatZeros(bM, 1)

    Ar := A.FloatArray()
    Xr := X.FloatArray()
    Y1r := Y1.FloatArray()

    blas.GemvFloat(A, X, Y0, 1.0, 1.0, linalg.OptTrans)
    t.Logf("blas: Y=A.T*X\n%v\n", Y0)

    DMultMV(Y1r, Ar, Xr, 1.0, 1.0, TRANSA, 1, A.LeadingIndex(), 1, 0,  bN, 0,  bM, 4, 4)
    t.Logf("Y0 == Y1: %v\n", Y0.AllClose(Y1))
    t.Logf("Y1: Y1 = A*X\n%v\n", Y1)
}

func _TestMultMVTransA(t *testing.T) {
    bM := 1000*M
    bN := 1000*N
    A := matrix.FloatNormal(bN, bM)
    //X := matrix.FloatNormal(bN, 1)
    X := matrix.FloatWithValue(bN, 1, 1.0)
    Y1 := matrix.FloatZeros(bM, 1)
    Y0 := matrix.FloatZeros(bM, 1)

    Ar := A.FloatArray()
    Xr := X.FloatArray()
    Y1r := Y1.FloatArray()

    blas.GemvFloat(A, X, Y0, 1.0, 1.0, linalg.OptTrans)
    //t.Logf("blas: Y=A.T*X\n%v\n", Y0)

    DMultMV(Y1r, Ar, Xr, 1.0, 1.0, TRANSA, 1, A.LeadingIndex(), 1, 0,  bN, 0,  bM, 4, 4)
    ok := Y0.AllClose(Y1)
    t.Logf("Y0 == Y1: %v\n", ok)
    if ! ok {
        y1 := Y1.SubMatrix(0, 0, 5, 1)
        t.Logf("Y1[0:5]:\n%v\n", y1)
        y0 := Y0.SubMatrix(0, 0, 5, 1)
        t.Logf("Y0[0:5]:\n%v\n", y0)
    }
}


func TestMultSymmSmall(t *testing.T) {
    //bM := 5
    bN := 7
    bP := 7
    Adata := [][]float64{
        []float64{1.0, 1.0, 1.0, 1.0, 1.0, 1.0, 1.0},
        []float64{0.0, 2.0, 2.0, 2.0, 2.0, 2.0, 2.0},
        []float64{0.0, 0.0, 3.0, 3.0, 3.0, 3.0, 3.0},
        []float64{0.0, 0.0, 0.0, 4.0, 4.0, 4.0, 4.0},
        []float64{0.0, 0.0, 0.0, 0.0, 5.0, 5.0, 5.0},
        []float64{0.0, 0.0, 0.0, 0.0, 0.0, 6.0, 6.0},
        []float64{0.0, 0.0, 0.0, 0.0, 0.0, 0.0, 7.0}}

    //A := matrix.FloatNormal(bN, bN)
    A := matrix.FloatMatrixFromTable(Adata, matrix.RowOrder)
    //B := matrix.FloatNormal(bN, bP)
    //A := matrix.FloatWithValue(bM, bP, 2.0)
    B := matrix.FloatWithValue(bN, bP, 2.0)
    C0 := matrix.FloatZeros(bN, bP)
    C1 := matrix.FloatZeros(bN, bP)

    Ar := A.FloatArray()
    Br := B.FloatArray()
    C1r := C1.FloatArray()

    t.Logf("A=\n%v\n", A)
    blas.SymmFloat(A, B, C0, 1.0, 1.0, linalg.OptUpper)
    t.Logf("blas: C=A*B\n%v\n", C0)

    DMultSymm2(C1r, Ar, Br, 1.0, 1.0, UPPER|LEFT, bN, A.LeadingIndex(), bN, bN, 0,  bP, 0,  bN, 2, 2, 2)
    t.Logf("C0 == C1: %v\n", C0.AllClose(C1))
    t.Logf("C1: C1 = A*X\n%v\n", C1)
}

func _TestMultSymmLowerSmall(t *testing.T) {
    //bM := 5
    bN := 5
    bP := 5
    Adata := [][]float64{
     []float64{1.0, 0.0, 0.0, 0.0, 0.0},
     []float64{1.0, 2.0, 0.0, 0.0, 0.0},
     []float64{1.0, 2.0, 3.0, 0.0, 0.0},
     []float64{1.0, 2.0, 3.0, 4.0, 0.0},
     []float64{1.0, 2.0, 3.0, 4.0, 5.0}}

    //A := matrix.FloatNormal(bN, bN)
    A := matrix.FloatMatrixFromTable(Adata, matrix.RowOrder)
    //B := matrix.FloatNormal(bN, bP)
    //A := matrix.FloatWithValue(bM, bP, 2.0)
    B := matrix.FloatWithValue(bN, bP, 1.0)
    C0 := matrix.FloatZeros(bN, bP)
    C1 := matrix.FloatZeros(bN, bP)

    Ar := A.FloatArray()
    Br := B.FloatArray()
    C1r := C1.FloatArray()

    t.Logf("A=\n%v\n", A)
    blas.SymmFloat(A, B, C0, 1.0, 1.0, linalg.OptLower)
    t.Logf("blas: C=A*B\n%v\n", C0)

    DMultSymm(C1r, Ar, Br, 1.0, 1.0, LOWER|LEFT, bN, A.LeadingIndex(), bN,
        bN, 0,  bP, 0,  bN, 4, 4, 4)
    t.Logf("C0 == C1: %v\n", C0.AllClose(C1))
    t.Logf("C1: C1 = A*X\n%v\n", C1)
}

func TestMultSymmUpper(t *testing.T) {
    //bM := 5
    bN := 100*N
    bP := 100*P
    A := matrix.FloatNormalSymmetric(bN, matrix.Upper)
    B := matrix.FloatNormal(bN, bP)
    C0 := matrix.FloatZeros(bN, bP)
    C1 := matrix.FloatZeros(bN, bP)

    Ar := A.FloatArray()
    Br := B.FloatArray()
    C1r := C1.FloatArray()

    blas.SymmFloat(A, B, C0, 1.0, 1.0, linalg.OptUpper)

    DMultSymm2(C1r, Ar, Br, 1.0, 1.0, UPPER|LEFT, bN, A.LeadingIndex(), bN,
        bN, 0,  bP, 0,  bN, 32, 32, 32)
    t.Logf("C0 == C1: %v\n", C0.AllClose(C1))
}

func _TestMultSymmLower(t *testing.T) {
    //bM := 5
    bN := 100*N
    bP := 100*P
    A := matrix.FloatNormalSymmetric(bN, matrix.Lower)
    B := matrix.FloatNormal(bN, bP)
    C0 := matrix.FloatZeros(bN, bP)
    C1 := matrix.FloatZeros(bN, bP)

    Ar := A.FloatArray()
    Br := B.FloatArray()
    C1r := C1.FloatArray()

    blas.SymmFloat(A, B, C0, 1.0, 1.0, linalg.OptLower)

    DMultSymm(C1r, Ar, Br, 1.0, 1.0, LOWER|LEFT, bN, A.LeadingIndex(), bN,
        bN, 0,  bP, 0,  bN, 32, 32, 32)
    t.Logf("C0 == C1: %v\n", C0.AllClose(C1))
}

func _TestRankSmall(t *testing.T) {
    bM := 5
    bN := 5
    //bP := 5
    Adata := [][]float64{
     []float64{1.0, 1.0, 1.0, 1.0, 1.0},
     []float64{2.0, 2.0, 2.0, 2.0, 2.0},
     []float64{3.0, 3.0, 3.0, 3.0, 3.0},
     []float64{4.0, 4.0, 4.0, 4.0, 4.0},
     []float64{5.0, 5.0, 5.0, 5.0, 5.0}}

    A := matrix.FloatMatrixFromTable(Adata, matrix.RowOrder)
    A0 := matrix.FloatMatrixFromTable(Adata, matrix.RowOrder)
    X := matrix.FloatVector([]float64{1.0, 2.0, 3.0, 4.0, 5.0})
    Y := matrix.FloatWithValue(bN, 1, 2.0)

    Ar := A.FloatArray()
    Xr := X.FloatArray()
    Yr := Y.FloatArray()

    t.Logf("A=\n%v\n", A)
    blas.GerFloat(X, Y, A0, 1.0)
    t.Logf("blas ger:\n%v\n", A0)

    DRankMV(Ar, Xr, Yr, 1.0, A.LeadingIndex(), 1, 1, 0,  bN, 0,  bM, 4, 4)
    t.Logf("A0 == A1: %v\n", A0.AllClose(A))
    t.Logf("A1: \n%v\n", A)
}

func _TestRank(t *testing.T) {
    bM := M*100
    bN := N*100
    //bP := 5

    A := matrix.FloatWithValue(bM, bN, 1.0);
    A0 := matrix.FloatWithValue(bM, bN, 1.0);
    X := matrix.FloatNormal(bM, 1);
    Y := matrix.FloatNormal(bN, 1);

    Ar := A.FloatArray()
    Xr := X.FloatArray()
    Yr := Y.FloatArray()

    blas.GerFloat(X, Y, A0, 1.0)

    DRankMV(Ar, Xr, Yr, 1.0, A.LeadingIndex(), 1, 1, 0,  bN, 0,  bM, 4, 4)
    t.Logf("A0 == A1: %v\n", A0.AllClose(A))
}

func TestMultSyrSmall(t *testing.T) {
    bN := 7
    //A := matrix.FloatNormal(bN, bN)
    //B := matrix.FloatNormal(bN, bP)
    //A := matrix.FloatWithValue(bM, bP, 1.0)
    X := matrix.FloatWithValue(bN, 1, 1.0)
    C0 := matrix.FloatZeros(bN, bN)
    C1 := matrix.FloatZeros(bN, bN)
    for i := 0; i < bN; i++ {
        X.Add(1.0+float64(i), i)
    }
    t.Logf("X=\n%v\n", X)

    Xr := X.FloatArray()
    C1r := C1.FloatArray()

    blas.SyrFloat(X, C0, 1.0, linalg.OptUpper)
    t.Logf("blas: C0\n%v\n", C0)

    DSymmRankMV(C1r, Xr, 1.0, UPPER, C1.LeadingIndex(), 1, 0,  bN, 4)
    t.Logf("C0 == C1: %v\n", C0.AllClose(C1))
    t.Logf("C1: C1 = A*X\n%v\n", C1)

    blas.SyrFloat(X, C0, 1.0, linalg.OptLower)
    t.Logf("blas: C0\n%v\n", C0)

    DSymmRankMV(C1r, Xr, 1.0, LOWER, C1.LeadingIndex(), 1, 0,  bN, 4)
    t.Logf("C0 == C1: %v\n", C0.AllClose(C1))
    t.Logf("C1: C1 = A*X\n%v\n", C1)
}

func TestMultSyr2Small(t *testing.T) {
    bN := 7
    //A := matrix.FloatNormal(bN, bN)
    //B := matrix.FloatNormal(bN, bP)
    //A := matrix.FloatWithValue(bM, bP, 1.0)
    X := matrix.FloatWithValue(bN, 1, 1.0)
    Y := matrix.FloatWithValue(bN, 1, 1.0)
    C0 := matrix.FloatZeros(bN, bN)
    C1 := matrix.FloatZeros(bN, bN)
    for i := 0; i < bN; i++ {
        X.Add(1.0+float64(i), i)
        Y.Add(2.0+float64(i), i)
    }
    t.Logf("X=\n%v\nY=\n%v\n", X, Y)

    Xr := X.FloatArray()
    Yr := Y.FloatArray()
    C1r := C1.FloatArray()

    blas.Syr2Float(X, Y, C0, 1.0, linalg.OptUpper)
    t.Logf("blas: C0\n%v\n", C0)

    DSymmRank2MV(C1r, Xr, Yr, 1.0, UPPER, C1.LeadingIndex(), 1, 1, 0,  bN, 4)
    t.Logf("C0 == C1: %v\n", C0.AllClose(C1))
    t.Logf("C1: C1 = A*X\n%v\n", C1)

    blas.Syr2Float(X, Y, C0, 1.0, linalg.OptLower)
    t.Logf("blas: C0\n%v\n", C0)

    DSymmRank2MV(C1r, Xr, Yr, 1.0, LOWER, C1.LeadingIndex(), 1, 1, 0,  bN, 4)
    t.Logf("C0 == C1: %v\n", C0.AllClose(C1))
    
    t.Logf("C1: C1 = A*X\n%v\n", C1)
}

func solveForwardTest(t *testing.T, A, X0 *matrix.FloatMatrix, unit bool, bN, bNB int) {
    X1 := X0.Copy()
    Ar := A.FloatArray()
    Xr := X1.FloatArray()

    if bN < 8 {
        t.Logf("A=\n%v\n", A)
        t.Logf("X0=\n%v\n", X0)
    }
    blas.TrsvFloat(A, X0, linalg.OptLower)
    if bN < 8 {
        t.Logf("blas: X0\n%v\n", X0)
    }

    if bN == bNB {
        DSolveLower(Xr, Ar, unit, 1, A.LeadingIndex(), bN, bN)
    } else {
        DSolveLowerBlocked(Xr, Ar, unit, 1, A.LeadingIndex(), bN, bNB)
    }
    t.Logf("X0 == X1: %v\n", X0.AllClose(X1))
    if bN < 8 {
        t.Logf("X1:\n%v\n", X1)
    }

}

func solveBackwardTest(t *testing.T, A, X0 *matrix.FloatMatrix, unit bool, bN, bNB int) {
    X1 := X0.Copy()

    if bN < 8 {
        t.Logf("A=\n%v\n", A)
        t.Logf("X0=\n%v\n", X0)
    }
    blas.TrsvFloat(A, X0, linalg.OptUpper)
    if bN < 8 {
        t.Logf("blas: X0\n%v\n", X0)
    }

    Ar := A.FloatArray()
    Xr := X1.FloatArray()
    if bN == bNB {
        DSolveUpper(Xr, Ar, unit,  1, A.LeadingIndex(), bN, bN)
    } else {
        DSolveUpperBlocked(Xr, Ar, unit,  1, A.LeadingIndex(), bN, bNB)
    }
    t.Logf("X1 == X0: %v\n", X1.AllClose(X0))
    if bN < 8 {
        t.Logf("X1:\n%v\n", X1)
    }
}


func TestSolveSmall(t *testing.T) {
    Adata := [][]float64{
        []float64{1.0, 0.0, 0.0, 0.0, 0.0},
        []float64{1.0, 2.0, 0.0, 0.0, 0.0},
        []float64{1.0, 2.0, 3.0, 0.0, 0.0},
        []float64{1.0, 2.0, 3.0, 4.0, 0.0},
        []float64{1.0, 2.0, 3.0, 4.0, 5.0}}

    A := matrix.FloatMatrixFromTable(Adata, matrix.RowOrder)
    bN := A.Rows()
    At := A.Transpose()
    X0 := matrix.FloatWithValue(A.Rows(), 1, 1.0)
    X1 := X0.Copy()
    xsum := 0.0
    for i := 0; i < bN; i++ {
        xsum += float64(i)
        X0.Add(xsum, i)
        X1.Add(xsum, -(i+1))
    }
    X2 := X0.Copy()
    X3 := X0.Copy()

    t.Logf("-- SOLVE NON-UNIT ---\n")
    solveForwardTest(t, A, X0, false, A.Rows(), A.Rows())
    t.Logf("-- SOLVE UNIT ---\n")
    A.Diag().SetIndexes(1.0)
    solveForwardTest(t, A, X1, true, A.Rows(), A.Rows())

    t.Logf("-- SOLVE NON-UNIT BACKWARD ---\n")
    solveBackwardTest(t, At, X2, false, At.Rows(), At.Rows())
    t.Logf("-- SOLVE UNIT BACKWARD ---\n")
    At.Diag().SetIndexes(1.0)
    solveBackwardTest(t, At, X3, true, At.Rows(), At.Rows())
}

func TestSolveBlockedSmall(t *testing.T) {
    Adata := [][]float64{
        []float64{1.0, 0.0, 0.0, 0.0, 0.0},
        []float64{1.0, 2.0, 0.0, 0.0, 0.0},
        []float64{1.0, 2.0, 3.0, 0.0, 0.0},
        []float64{1.0, 2.0, 3.0, 4.0, 0.0},
        []float64{1.0, 2.0, 3.0, 4.0, 5.0}}

    A := matrix.FloatMatrixFromTable(Adata, matrix.RowOrder)
    X0 := matrix.FloatWithValue(A.Rows(), 1, 1.0)
    X1 := X0.Copy()
    X2 := X0.Copy()
    xsum := 0.0
    for i := 0; i < A.Rows(); i++ {
        xsum += float64(i)
        X0.Add(xsum, i)
        X2.Add(xsum, -(i+1))
    }
    t.Logf("-- SOLVE NON-UNIT ---\n")
    solveForwardTest(t, A, X0, false, A.Rows(), 4)
    solveBackwardTest(t, A.Transpose(), X2, false, A.Rows(), 4)
    
    t.Logf("-- SOLVE UNIT ---\n")
    A.Diag().SetIndexes(1.0)
    solveForwardTest(t, A, X1, true, A.Rows(), 4)
}

func TestSolveRandom(t *testing.T) {
    bN := 22
    A := matrix.FloatNormalSymmetric(bN, matrix.Lower)
    At := A.Transpose()
    X0 := matrix.FloatWithValue(A.Rows(), 1, 1.0)
    X1 := X0.Copy()
    X2 := X0.Copy()
    t.Logf("-- BLOCKED SOLVE FORWARD NON-UNIT ---\n")
    solveForwardTest(t, A, X0, false, bN, 4)
    t.Logf("-- BLOCKED SOLVE FORWARD UNIT ---\n")
    A.Diag().SetIndexes(1.0)
    solveForwardTest(t, A, X1, true, bN, 4)
    t.Logf("-- BLOCKED SOLVE BACKWARD NON-UNIT ---\n")
    solveBackwardTest(t, At, X2, false, bN, 4)
}


func tridiagSmall(t *testing.T, unit bool) {
    //bM := 5
    bN := 5
    Adata := [][]float64{
        []float64{1.0, 0.0, 0.0, 0.0, 0.0},
        []float64{1.0, 2.0, 0.0, 0.0, 0.0},
        []float64{1.0, 2.0, 3.0, 0.0, 0.0},
        []float64{1.0, 2.0, 3.0, 4.0, 0.0},
        []float64{1.0, 2.0, 3.0, 4.0, 5.0}}

    diag := linalg.OptNonUnit
    if unit {
        diag = linalg.OptUnit
    }
    Al := matrix.FloatMatrixFromTable(Adata, matrix.RowOrder)
    //A := matrix.FloatNormal(bN, bN)
    //A := matrix.FloatWithValue(bM, bP, 2.0)
    //Z := matrix.FloatNormal(bN, 1);
    X0 := matrix.FloatWithValue(bN, 1, 2.0)
    X2 := matrix.FloatWithValue(bN, 1, 2.0)
    xsum := 0.0
    for i := 0; i < bN; i++ {
        xsum += float64(i) + 1.0
        //X0.Add(xsum, i)
        X2.Add(xsum, -(i+1))
    }
    //X0.Mul(Z)
    X1 := X0.Copy()
    //X2.Mul(Z)
    //X3 := X2.Copy()
    Au := Al.Transpose()

    t.Logf("X0=\n%v\n", X0)
    t.Logf("A(upper)=\n%v\n", Au)
    t.Logf("A(lower)=\n%v\n", Al)
    //t.Logf("Z=\n%v\n", Z)
    blas.TrmvFloat(Al, X0, linalg.OptUpper, diag)
    t.Logf("1. A(lower), blas(upper): X0 = Al*X0\n%v\n", X0)

    Ar := Al.FloatArray()
    Xr := X1.FloatArray()
    DTrimvUpper(Xr, Ar, unit, 1, Al.LeadingIndex(), bN, bN)
    t.Logf("   X0 == X1: %v\n", X0.AllClose(X1))
    t.Logf("   A(lower), X1(fwd) = Al*X1:\n%v\n", X1)
    
    X0.SetIndexes(2.0)
    X1.SetIndexes(2.0)

    blas.TrmvFloat(Au, X0, linalg.OptUpper, diag)
    t.Logf("2. A(upper), blas(upper): X0 = Au*X0\n%v\n", X0)

    Ar = Au.FloatArray()
    Xr = X1.FloatArray()
    DTrimvUpper(Xr, Ar, unit, 1, Au.LeadingIndex(), bN, bN)
    t.Logf("   X0 == X1: %v\n", X0.AllClose(X1))
    t.Logf("   A(upper), X1(fwd) = Au*X1:\n%v\n", X1)

    X0.SetIndexes(2.0)
    X1.SetIndexes(2.0)
    blas.TrmvFloat(Al, X0, linalg.OptLower, diag)
    t.Logf("3. A(lower), blas(lower): X0 = Al*X0\n%v\n", X0)

    Ar = Al.FloatArray()
    Xr = X1.FloatArray()
    DTrimvLower(Xr, Ar, unit, 1, Al.LeadingIndex(), bN, bN)
    t.Logf("   X0 == X1: %v\n", X0.AllClose(X1))
    t.Logf("   A(lower), X1(backwd) = Al*X1:\n%v\n", X1)

    X0.SetIndexes(2.0)
    X1.SetIndexes(2.0)
    blas.TrmvFloat(Au, X0, linalg.OptLower, diag)
    t.Logf("4. A(upper), blas(lower): X0 = Au*X0\n%v\n", X0)

    Ar = Au.FloatArray()
    Xr = X1.FloatArray()
    DTrimvLower(Xr, Ar, unit, 1, Au.LeadingIndex(), bN, bN)
    t.Logf("   X0 == X1: %v\n", X0.AllClose(X1))
    t.Logf("   A(upper), X1(backwd) = Au*X1:\n%v\n", X1)

    t.Logf("-- TRANSPOSED --\n")
    Au_T := Al
    Al_T := Au
    t.Logf("A(upper).T=\n%v\n", Au_T)
    t.Logf("A(lower).T=\n%v\n\n", Al_T)

    X0.SetIndexes(2.0)
    X1.SetIndexes(2.0)

    blas.TrmvFloat(Au_T, X0, linalg.OptUpper, linalg.OptTrans, diag)
    t.Logf("5. A(upper).T, blas(upper,trans): X0 = Au.T*X0\n%v\n", X0)

    Ar = Au_T.FloatArray()
    Xr = X1.FloatArray()
    DTrimvUpperTransA(Xr, Ar, unit, 1, Au_T.LeadingIndex(), bN, bN)
    t.Logf("   X0 == X1: %v\n", X0.AllClose(X1))
    t.Logf("   A(upper).T, X1(fwd,trans) = Au.T*X1:\n%v\n", X1)

    X0.SetIndexes(2.0)
    X1.SetIndexes(2.0)

    blas.TrmvFloat(Au_T, X0, linalg.OptLower, linalg.OptTrans, diag)
    t.Logf("6. A(upper).T, blas(lower,trans): X0 = Au.T*X0\n%v\n", X0)

    Ar = Au_T.FloatArray()
    Xr = X1.FloatArray()
    DTrimvLowerTransA(Xr, Ar, unit, 1, Au_T.LeadingIndex(), bN, bN)
    t.Logf("   X0 == X1: %v\n", X0.AllClose(X1))
    t.Logf("   A(upper).T, X1(backwd,transA) = Au.T*X1:\n%v\n", X1)

    X0.SetIndexes(2.0)
    X1.SetIndexes(2.0)

    blas.TrmvFloat(Al_T, X0, linalg.OptLower, linalg.OptTrans, diag)
    t.Logf("7. A(lower).T, blas(lower,trans): X0 = Al.T*X0\n%v\n", X0)

    Ar = Al_T.FloatArray()
    Xr = X1.FloatArray()
    DTrimvLowerTransA(Xr, Ar, unit, 1, Al.LeadingIndex(), bN, bN)
    t.Logf("   X0 == X1: %v\n", X0.AllClose(X1))
    t.Logf("   A(lower).T, X1(backwd,trans) = Al.T*X1:\n%v\n", X1)

    X0.SetIndexes(2.0)
    X1.SetIndexes(2.0)

    blas.TrmvFloat(Al_T, X0, linalg.OptUpper, linalg.OptTrans, diag)
    t.Logf("8. A(lower).T, blas(upper,trans): X0 = Al.T*X0\n%v\n", X0)

    Ar = Al_T.FloatArray()
    Xr = X1.FloatArray()
    DTrimvUpperTransA(Xr, Ar, unit, 1, Al_T.LeadingIndex(), bN, bN)
    t.Logf("   X0 == X1: %v\n", X0.AllClose(X1))
    t.Logf("   A(lower).T, X1(fwd,trans) = Al.T*X1:\n%v\n", X1)

}

func TestTridiagNonUnitSmall(t *testing.T) {
    tridiagSmall(t, false) 
}

func TestTridiagUnitSmall(t *testing.T) {
    tridiagSmall(t, true) 
}

func trmmUpper(t *testing.T, A *matrix.FloatMatrix, unit bool) {
    B0 := matrix.FloatWithValue(A.Rows(), 2, 2.0)
    B1 := B0.Copy()
    diag := linalg.OptNonUnit
    if unit {
        diag = linalg.OptUnit
    }

    blas.TrmmFloat(A, B0, 1.0, linalg.OptUpper, diag)
    if A.Rows() < 8 {
        t.Logf("  BLAS: B0 = A*B0\n%v\n", B0)
    }

    Ar := A.FloatArray()
    Br := B1.FloatArray()
    DTrmmUpper(Br, Ar, 1.0, unit, B1.LeadingIndex(), A.LeadingIndex(), A.Cols(), 0, B1.Cols())
    t.Logf("   B0 == B1: %v\n", B0.AllClose(B1))
    if A.Rows() < 8 {
        t.Logf("  DTrmmUpper: B1 = A*B1\n%v\n", B1)
    }
}

func trmmUpperTransA(t *testing.T, A *matrix.FloatMatrix, unit bool) {
    B0 := matrix.FloatWithValue(A.Rows(), 2, 2.0)
    B1 := B0.Copy()
    diag := linalg.OptNonUnit
    if unit {
        diag = linalg.OptUnit
    }

    blas.TrmmFloat(A, B0, 1.0, linalg.OptUpper, linalg.OptTransA, diag)
    if A.Rows() < 8 {
        t.Logf("  BLAS: B0 = A*B0\n%v\n", B0)
    }

    Ar := A.FloatArray()
    Br := B1.FloatArray()
    DTrmmUpperTransA(Br, Ar, 1.0, unit, B1.LeadingIndex(), A.LeadingIndex(), A.Cols(), 0, B1.Cols())
    t.Logf("   B0 == B1: %v\n", B0.AllClose(B1))
    if A.Rows() < 8 {
        t.Logf("  DTrmmUpperTransA: B1 = A*B1\n%v\n", B1)
    }
}

func trmmLower(t *testing.T, A *matrix.FloatMatrix, unit bool) {
    B0 := matrix.FloatWithValue(A.Rows(), 2, 2.0)
    B1 := B0.Copy()
    diag := linalg.OptNonUnit
    if unit {
        diag = linalg.OptUnit
    }

    blas.TrmmFloat(A, B0, 1.0, linalg.OptLower, diag)
    if A.Rows() < 8 {
        t.Logf("  BLAS: B0 = A*B0\n%v\n", B0)
    }

    Ar := A.FloatArray()
    Br := B1.FloatArray()
    DTrmmLower(Br, Ar, 1.0, unit, B1.LeadingIndex(), A.LeadingIndex(), A.Cols(), 0, B1.Cols())
    t.Logf("   B0 == B1: %v\n", B0.AllClose(B1))
    if A.Rows() < 8 {
        t.Logf("  DTrmmLower: B1 = A*B1\n%v\n", B1)
    }
}

func trmmLowerTransA(t *testing.T, A *matrix.FloatMatrix, unit bool) {
    B0 := matrix.FloatWithValue(A.Rows(), 2, 2.0)
    B1 := B0.Copy()
    diag := linalg.OptNonUnit
    if unit {
        diag = linalg.OptUnit
    }

    blas.TrmmFloat(A, B0, 1.0, linalg.OptLower, linalg.OptTransA, diag)
    if A.Rows() < 8 {
        t.Logf("  BLAS: B0 = A*B0\n%v\n", B0)
    }

    Ar := A.FloatArray()
    Br := B1.FloatArray()
    DTrmmLowerTransA(Br, Ar, 1.0, unit, B1.LeadingIndex(), A.LeadingIndex(), A.Cols(), 0, B1.Cols())
    t.Logf("   B0 == B1: %v\n", B0.AllClose(B1))
    if A.Rows() < 8 {
        t.Logf("  DTrmmLowerTransA: B1 = A*B1\n%v\n", B1)
    }
}

func TestTrmm(t *testing.T) {
    //bN := 5
    Adata := [][]float64{
        []float64{1.0, 0.0, 0.0, 0.0, 0.0},
        []float64{1.0, 2.0, 0.0, 0.0, 0.0},
        []float64{1.0, 2.0, 3.0, 0.0, 0.0},
        []float64{1.0, 2.0, 3.0, 4.0, 0.0},
        []float64{1.0, 2.0, 3.0, 4.0, 5.0}}

    A := matrix.FloatMatrixFromTable(Adata, matrix.RowOrder)
    t.Logf("-- TRMM-UPPER, NON-UNIT ---")
    trmmUpper(t, A.Transpose(), false)
    t.Logf("-- TRMM-UPPER, NON-UNIT, MATRIX-LOWER ---")
    trmmUpper(t, A, false)
    t.Logf("-- TRMM-UPPER, NON-UNIT, TRANSA ---")
    trmmUpperTransA(t, A, false)
    t.Logf("-- TRMM-LOWER, NON-UNIT ---")
    trmmLower(t, A, false)
    t.Logf("-- TRMM-LOWER, NON-UNIT, MATRIX-UPPER ---")
    trmmLower(t, A.Transpose(), false)
    t.Logf("-- TRMM-LOWER, NON-UNIT, TRANSA ---")
    trmmLowerTransA(t, A, false)
    A.Diag().SetIndexes(1.0)
    t.Logf("-- TRMM-UPPER, UNIT, TRANSA ---")
    trmmUpperTransA(t, A, true)
    t.Logf("-- TRMM-LOWER, UNIT, TRANSA ---")
    trmmLowerTransA(t, A, true)
}

// Local Variables:
// tab-width: 4
// indent-tabs-mode: nil
// End:
