
#include "cmops.h"
#include "colcpy.h"
#include "inner_ddot.h"

// this will compute sub-block matrix product: Cij += Aik * Bkj using
// successive inner vector product (DOT) function.
void _dblock_ddot_sse(double *Cc, const double *Aroot, const double *Bc, double alpha,
                      int ldC, int ldA, int ldB, int nSL, int nRE, int nVP)
{
  register int i, j;
  register double *c0, *c1, *c2, *c3;
  register const double *Br0, *Br1, *Br2, *Br3;
  const double *Ac;
  double *Cx = Cc;
  const double *Bx = Bc;

  for (j = 0; j < nSL-3; j += 4) {
    Ac = Aroot;
    Br0 = Bc;
    Br1 = Br0 + ldB;
    Br2 = Br1 + ldB;
    Br3 = Br2 + ldB;
    c0 = Cc;
    c1 = c0 + ldC;
    c2 = c1 + ldC;
    c3 = c2 + ldC;

    for (i = 0; i < nRE-1; i += 2) {
      _inner_ddot4_2_sse3(c0, c1, c2, c3, Ac, Ac+ldA, Br0, Br1, Br2, Br3, alpha, nVP);
      //_inner_ddot4_sse3(c0, c1, c2, c3, Ac, Br0, Br1, Br2, Br3, alpha, nVP);
      //_inner_ddot4_sse3(c0+1, c1+1, c2+1, c3+1, Ac+ldA, Br0, Br1, Br2, Br3, alpha, nVP);
      Ac += (ldA << 1);
      c0 +=2;
      c1 +=2;
      c2 +=2;
      c3 +=2;
    }
    if (i < nRE) {
      _inner_ddot4_sse3(c0, c1, c2, c3, Ac, Br0, Br1, Br2, Br3, alpha, nVP);
      Ac += ldA;
      c0++;
      c1++;
      c2++;
      c3++;
    }
    // forward 4 columns in C, B
    Cc += (ldC << 2);
    Bc += (ldB << 2);
  }
  // Here if j == nSL --> nSL mod 4 == 0 and we are done
  // If work is divided right this should happen most of the time.
  if (j == nSL) {
    return;
  }

  // do the not-multiples of 4 cases....
  if (j < nSL-1) {
    Ac = Aroot;
    Br0 = Bc;
    Br1 = Br0 + ldB;
    c0 = Cc;
    c1 = c0 + ldC;
    for (i = 0; i < nRE; i++) {
      _inner_ddot2_sse3(c0, c1, Ac, Br0, Br1, alpha, nVP);
      Ac += ldA;
      c0++;
      c1++;
    }
    // forward 2 columns in C, B
    Cc += (ldC << 1);
    Bc += (ldB << 1);
    j += 2;
  }

  if (j < nSL) {
    // not multiple of 2
    Ac = Aroot;
    Br0 = Bc;
    c0 = Cc;
    for (i = 0; i < nRE; i++) {
      _inner_ddot_sse(c0, Ac, Br0, alpha, nVP);
      Ac += ldA;
      c0++;
    }
    j++;
  }
}


void _dblock_mult_cpy(mdata_t *C, const mdata_t *A, const mdata_t *B,
                      double alpha, double beta, int flags, 
                      int nP, int S, int L, int R, int E, int vlen)
{
  int j, k, vpS, vpL, nC, nB, nA, L1, L2;
  const double *Bc, *Ac, *AvpS, *Bp;
  double *Cc, *Cp; 
  double Acpy[MAX_VP_ROWS*MAX_VP_COLS] __attribute__((aligned(64)));
  double Bcpy[MAX_VP_ROWS*MAX_VP_COLS] __attribute__((aligned(64)));


  if (vlen > nP) {
    vlen = nP;
  }
  vpS = 0;
  vpL = vlen < nP ? vlen : nP;

  // start of block, element [0, 0]
  Cc = &C->md[S*C->step+R];

  // scaling the block with beta
  dscale_tile(Cc, C->step, beta, E-R, L-S);

  while (vpS < nP) {
    nB = vpL-vpS;
    //nB += (nB & 0x1);
    //nA = nB;

    // column viewport start in panel B[:,S]
    // row viewport start A[R,:]

    // transpose A on copy to be able to DOT operations.
    nA = nB = MAX_VP_ROWS;
    if (flags & MTX_TRANSB) {
      Bc = &B->md[vpS*B->step + S];
      colcpy4_trans(Bcpy, nB, Bc, B->step, L-S, vpL-vpS);
    } else {
      Bc = &B->md[S*B->step + vpS];
      colcpy(Bcpy, nB, Bc, B->step, vpL-vpS, L-S);
    }

    if (flags & MTX_TRANSA) {
      AvpS = &A->md[R*A->step + vpS];
      colcpy(Acpy, nA, AvpS, A->step, vpL-vpS, E-R);
    } else {
      AvpS = &A->md[vpS*A->step + R];
      colcpy4_trans(Acpy, nA, AvpS, A->step, E-R, vpL-vpS);
    }

    //printf(".. R=%d, E=%d, S=%d, L=%d, vpS=%d, vpL=%d, F=0x%x\n", R, E, S, L, vpS, vpL, flags);
    //printf(".. A=\n"); print_tile(Acpy, nA, vpL-vpS, E-R);
    //printf(".. B=\n"); print_tile(Bc, B->step, vpL-vpS, L-S);
    _dblock_ddot_sse(Cc, Acpy, Bcpy, alpha, C->step, nA, nB, L-S, E-R, vpL-vpS);
    //printf(".. C=\n"); print_tile(Cc, C->step, E-R, L-S);

    vpS = vpL;
    vpL += vlen;
    if (vpL > nP) {
      vpL = nP;
    }
  }
}

void dmult_mm_blocked(mdata_t *C, const mdata_t *A, const mdata_t *B,
                      double alpha, double beta, int flags,
                      int P, int S, int L, int R, int E, 
                      int vlen, int NB, int MB)
{
  int i, j, nI, nJ;

  // restrict block sizes as data is copied to aligned buffers of predefined max sizes.
  if (NB > MAX_VP_ROWS || NB <= 0) {
    NB = MAX_VP_ROWS;
  }
  if (MB > MAX_VP_COLS || MB <= 0) {
    MB = MAX_VP_COLS;
  }
  if (vlen > MAX_VP_ROWS || vlen <= 0) {
    vlen = MAX_VP_ROWS;
  }

  for (j = S; j < L; j += NB) {
    nJ = L - j < NB ? L - j : NB;
    for (i = R; i < E; i += MB) {
      nI = E - i < MB ? E - i : MB;
      //printf("pre : C=\n"); print_tile(C->md, C->step, E-R, L-S);
      _dblock_mult_cpy(C, A, B, alpha, beta, flags, P, j, j+nJ, i, i+nI, vlen);
      //printf("post: C=\n"); print_tile(C->md, C->step, E-R, L-S);
    }
  }
}

// Local Variables:
// indent-tabs-mode: nil
// End: