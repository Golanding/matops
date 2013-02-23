
// Copyright (c) Harri Rautila, 2012,2013

// This file is part of github.com/hrautila/matops package. It is free software,
// distributed under the terms of GNU Lesser General Public License Version 3, or
// any later version. See the COPYING tile included in this archive.

#include <stdio.h>
#include <stdlib.h>
#include <stdint.h>

#include "cmops.h"
#include "inner_vec_axpy.h"
#include "inner_vec_dot.h"

// Functions here implement various versions of TRMV operation.

// Calculates backward a diagonal block and updates Xc values from last to first.
// Updates are calculated in breadth first manner by with successive AXPY operations.
static void
_dmvec_trid_axpy_backward(double *Xc, const double *Ac, int unit,
                          int incX, int ldA, int nRE)
{
  // Y is 
  register int i;
  register double *xr, xtmp;
  register const double *Ar, *Acl;

  // diagonal matrix of nRE rows/cols and vector X of length nRE
  // move to point to last column and last entry of X.
  Acl = Ac + (nRE-1)*ldA;
  xr = Xc + (nRE-1)*incX;

  // xr is the current X element, Ar is row in current A column.
  for (i = nRE; i > 0; i--) {
    Ar = Acl + i - 1; // move on diagonal

    // update all x-values below with the current A column and current X
    _inner_vec_daxpy(xr+incX, incX, Ar+1, xr, incX, 1.0, nRE-i);
    xr[0] *= unit ? 1.0 : Ar[0];

    // previous X, previous column in A 
    xr  -= incX;
    Acl -= ldA;
  }
}

// Calculates backward a diagonal block and updates Xc values from last to first.
// Updates are calculated in depth first manner with DOT operations.
static void
_dmvec_trid_dot_backward(double *Xc, const double *Ac, int unit,
                         int incX, int ldA, int nRE)
{
  // Y is 
  register int i;
  register double *xr, xtmp;
  register const double *Ar, *Acl;

  // lower diagonal matrix (transposed) of nRE rows/cols and vector X of length nRE
  // we do it really forward!! unlike the _axpy method above.
  Acl = Ac + (nRE-1)*ldA;
  xr = Xc + (nRE-1)*incX;

  // xr is the current X element, Ar is row in current A column.
  for (i = 0; i < nRE; i++) {
    Ar = Ac + i; // move on diagonal

    // update all x-values below with the current A column and current X
    xtmp = unit ? 1.0 : 0.0;
    _inner_vec_ddot(&xtmp, 1, Acl, Xc, incX, 1.0, nRE-unit-i);
    xr[0] = xtmp;

    // previous X, previous column in A 
    xr -= incX;
    Acl -= ldA;
  }
}

// Calculate forward a diagonal block and updates Xc values from first to last.
static void
_dmvec_trid_axpy_forward(double *Xc, const double *Ac, double unit,
                         int incX, int ldA, int nRE)
{
  // Y is 
  register int i;
  register double *xr, xtmp;
  register const double *Ar;

  // upper diagonal matrix of nRE rows/cols and vector X, Y of length nRE
  xr = Xc;

  // xr is the current X element, Ar is row in current A column.
  for (i = 0; i < nRE; i++) {
    // update all previous x-values with current A column and current X
    _inner_vec_daxpy(Xc, incX, Ac, xr, incX, 1.0, i);
    Ar = Ac + i;
    xr[0] *= unit ? 1.0 : Ar[0];
    // next X, next column in A 
    xr += incX;
    Ac += ldA;
  }
}

// Calculate forward a diagonal block and updates Xc values from first to last.
static void
_dmvec_trid_dot_forward(double *Xc, const double *Ac, int unit, int incX, int ldA, int nRE)
{
  // Y is 
  register int i;
  register double *xr, xtmp;
  register const double *Ar;

  xr = Xc;
  // xr is the current X element, Ar is row in current A column.
  for (i = 0; i < nRE; i++) {
    Ar = Ac + i + unit;
    // update all previous x-values with current A column and current X
    xtmp = unit ? 1.0 : 0.0;
    _inner_vec_ddot(&xtmp, 1, Ar, xr, incX, 1.0, nRE-unit-i);
    xr[0] = xtmp;
    // next X, next column in A 
    xr += incX;
    Ac += ldA;
  }
}

//extern void memset(void *, int, size_t);

// X = A*X; unblocked version
void dmvec_trid_unb(mvec_t *X, const mdata_t *A, int flags, int N)
{
  // indicates if diagonal entry is unit (=1.0) or non-unit.
  int unit = flags & MTX_UNIT ? 1 : 0;

  if (flags & MTX_UPPER) {
    if (flags & MTX_TRANSA) {
      _dmvec_trid_dot_backward(X->md, A->md, unit, X->inc, A->step, N);
    } else {
      _dmvec_trid_axpy_forward(X->md, A->md, unit, X->inc, A->step, N);
    }
  } else {
    if (flags & MTX_TRANSA) {
      _dmvec_trid_dot_forward(X->md, A->md, unit, X->inc, A->step, N);
    } else {
      _dmvec_trid_axpy_backward(X->md, A->md, unit, X->inc, A->step, N);
    }
  }
}


void dmvec_trid_blocked(mvec_t *X, const mdata_t *A, double alpha, int flags, int N, int NB)
{
  int i, nI;
  mvec_t Y;
  // indicates if diagonal entry is unit (=1.0) or non-unit.
  int unit = flags & MTX_UNIT ? 1 : 0;

  if (NB <= 0) {
    NB = 68;
  }
  //memset(cB, 0, sizeof(cB));

  if (flags & MTX_UPPER) {
    for (i = 0; i < N; i += NB) {
      nI = N - i < NB ? N - i : NB;
      // solve forward using Y values 
      _dmvec_trid_axpy_forward(&X->md[i], &A->md[i*A->step+i], unit, X->inc, A->step, nI);
    }
  }
}

// Local Variables:
// indent-tabs-mode: nil
// End:
