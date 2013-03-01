
#ifndef __INNER_DDOT_H
#define __INNER_DDOT_H

#include <x86intrin.h>

#ifndef INLINE
#define INLINE static inline
#endif

// SSE3 version with HADDPD (Horizontal Add Packed Double)
// update 4 columns from 2 by 4 block.
INLINE void _inner_ddot4_2_sse3(double *Cr0, double *Cr1, double *Cr2, double *Cr3,
				 const double *a0, const double *a1,
				 const double *b0, const double *b1,
				 const double *b2, const double *b3, double alpha, int nVP)
{
  register int k, i;
  register double cval0, cval1, cval2, cval3;
  register __m128d A0, B0, C0_0, F0, A1, B1, C1_0, F1;
  register __m128d A2, B2, C2_0, F2, A3, B3, C3_0, F3, ALP;
  register __m128d C0_1, C1_1, C2_1, C3_1;

  C0_0 = _mm_set1_pd(0.0);
  C1_0 = C0_0;
  C2_0 = C0_0;
  C3_0 = C0_0;
  C0_1 = C0_0;
  C1_1 = C0_0;
  C2_1 = C0_0;
  C3_1 = C0_0;
  ALP = _mm_set1_pd(alpha);

  // unrolling of loops;
  for (k = nVP-1; k > 0; k -= 2) {
    A0 = _mm_load_pd(a0);
    B0 = _mm_load_pd(b0);
    B1 = _mm_load_pd(b1);
    B2 = _mm_load_pd(b2);
    B3 = _mm_load_pd(b3);
    F0   = A0 * B0;
    F1   = A0 * B1;
    F2   = A0 * B2;
    F3   = A0 * B3;
    C0_0 = C0_0 + F0;
    C1_0 = C1_0 + F1;
    C2_0 = C2_0 + F2;
    C3_0 = C3_0 + F3;

    A0 = _mm_load_pd(a1);
    F0   = A0 * B0;
    F1   = A0 * B1;
    F2   = A0 * B2;
    F3   = A0 * B3;
    C0_1 = C0_1 + F0;
    C1_1 = C1_1 + F1;
    C2_1 = C2_1 + F2;
    C3_1 = C3_1 + F3;

    a0 += 2;
    a1 += 2;
    b0 += 2;
    b1 += 2;
    b2 += 2;
    b3 += 2;
  }
  if (k == 0) {
    cval0 = a0[0] * b0[0];
    cval1 = a0[0] * b1[0]; 
    cval2 = a0[0] * b2[0];
    cval3 = a0[0] * b3[0]; 
    Cr0[0] += cval0 * alpha;
    Cr1[0] += cval1 * alpha;
    Cr2[0] += cval2 * alpha;
    Cr3[0] += cval3 * alpha;
    cval0 = a1[0] * b0[0];
    cval1 = a1[0] * b1[0]; 
    cval2 = a1[0] * b2[0];
    cval3 = a1[0] * b3[0]; 
    Cr0[1] += cval0 * alpha;
    Cr1[1] += cval1 * alpha;
    Cr2[1] += cval2 * alpha;
    Cr3[1] += cval3 * alpha;
  }
  // update result
  C0_0 = C0_0 * ALP;
  C1_0 = C1_0 * ALP;
  C2_0 = C2_0 * ALP;
  C3_0 = C3_0 * ALP;
  F0 = _mm_hadd_pd(C0_0, C1_0);
  F1 = _mm_hadd_pd(C2_0, C3_0);
  Cr0[0] += F0[0];
  Cr1[0] += F0[1];
  Cr2[0] += F1[0];
  Cr3[0] += F1[1];

  C0_1 = C0_1 * ALP;
  C1_1 = C1_1 * ALP;
  C2_1 = C2_1 * ALP;
  C3_1 = C3_1 * ALP;
  F2 = _mm_hadd_pd(C0_1, C1_1);
  F3 = _mm_hadd_pd(C2_1, C3_1);
  Cr0[1] += F2[0];
  Cr1[1] += F2[1];
  Cr2[1] += F3[0];
  Cr3[1] += F3[1];
}


// update 4 columns from 1 by 4 block.
INLINE void _inner_ddot4_sse3(double *Cr0, double *Cr1, double *Cr2, double *Cr3,
			      const double *Ar, const double *b0, const double *b1,
			      const double *b2, const double *b3, double alpha, int nVP)
{
  register int k, i;
  register double cval0, cval1, cval2, cval3;
  register __m128d A0, B0, C0, F0, A1, B1, C1, F1;
  register __m128d A2, B2, C2, F2, A3, B3, C3, F3, ALP;

  C0 = _mm_set1_pd(0.0);
  C1 = C0;
  C2 = C0;
  C3 = C0;
  ALP = _mm_set1_pd(alpha);

  // unrolling of loops;
  for (k = nVP-1; k > 0; k -= 2) {
    A0 = _mm_load_pd(Ar);
    B0 = _mm_load_pd(b0);
    B1 = _mm_load_pd(b1);
    B2 = _mm_load_pd(b2);
    B3 = _mm_load_pd(b3);
    F0 = A0 * B0;
    F1 = A0 * B1;
    F2 = A0 * B2;
    F3 = A0 * B3;
    C0 = C0 + F0;
    C1 = C1 + F1;
    C2 = C2 + F2;
    C3 = C3 + F3;
    Ar += 2;
    b0 += 2;
    b1 += 2;
    b2 += 2;
    b3 += 2;
  }
  if (k == 0) {
    cval0   = Ar[0] * b0[0];
    cval1   = Ar[0] * b1[0]; 
    cval2   = Ar[0] * b2[0];
    cval3   = Ar[0] * b3[0]; 
    Cr0[0] += cval0 * alpha;
    Cr1[0] += cval1 * alpha;
    Cr2[0] += cval2 * alpha;
    Cr3[0] += cval3 * alpha;
  }

  C0 = C0 * ALP;
  C1 = C1 * ALP;
  C2 = C2 * ALP;
  C3 = C3 * ALP;
  F0 = _mm_hadd_pd(C0, C1);
  F1 = _mm_hadd_pd(C2, C3);
  Cr0[0] += F0[0]; //C0[0] + C0[1]; //F0[1];
  Cr1[0] += F0[1]; //C1[0] + C1[1]; //F0[0];
  Cr2[0] += F1[0]; //C2[0] + C2[1]; //F1[1];
  Cr3[0] += F1[1]; //C3[0] + C3[1]; //F1[0];
}

// Update 2 columns with 2 by 2 block.
INLINE void _inner_ddot2_2_sse3(double *Cr0, double *Cr1, 
				const double *a0, const double *a1,
				const double *b0, const double *b1,
				double alpha, int nVP)
{
  register int k, i;
  register double cval0, cval1, cval2, cval3;
  register __m128d A0, B0, C0_0, F0, A1, B1, C1_0, F1;
  register __m128d A2, B2, C2_0, F2, A3, B3, C3_0, F3, ALP;
  register __m128d C0_1, C1_1, C2_1, C3_1;

  C0_0 = _mm_set1_pd(0.0);
  C1_0 = C0_0;
  C0_1 = C0_0;
  C1_1 = C1_0;
  ALP = _mm_set1_pd(alpha);

  // unrolling of loops;
  for (k = nVP-1; k > 0; k -= 2) {
    A0 = _mm_load_pd(a0);
    B0 = _mm_load_pd(b0);
    B1 = _mm_load_pd(b1);
    F0   = A0 * B0;
    F1   = A0 * B1;
    C0_0 = C0_0 + F0;
    C1_0 = C1_0 + F1;

    A0 = _mm_load_pd(a1);
    F0   = A0 * B0;
    F1   = A0 * B1;
    C0_1 = C0_1 + F0;
    C1_1 = C1_1 + F1;

    a0 += 2;
    a1 += 2;
    b0 += 2;
    b1 += 2;
  }
  if (k == 0) {
    cval0   = a0[0] * b0[0];
    cval1   = a0[0] * b1[0]; 
    Cr0[0] += cval0 * alpha;
    Cr1[0] += cval1 * alpha;
    cval0   = a1[0] * b0[0];
    cval1   = a1[0] * b1[0]; 
    Cr0[1] += cval0 * alpha;
    Cr1[1] += cval1 * alpha;
  }

 update:
  C0_0 = C0_0 * ALP;
  C1_0 = C1_0 * ALP;
  F0 = _mm_hadd_pd(C0_0, C1_0);
  Cr0[0] += F0[0];
  Cr1[0] += F0[1];

  C0_1 = C0_1 * ALP;
  C1_1 = C1_1 * ALP;
  F1 = _mm_hadd_pd(C0_1, C1_1);
  Cr0[1] += F1[0];
  Cr1[1] += F1[1];
}


// Update 2 columns with 1 by 2 block.
INLINE void _inner_ddot2_sse3(double *Cr0, double *Cr1, const double *Ar,
			      const double *b0, const double *b1, double alpha, int nVP)
{
  register int k, i;
  register double cval0, cval1;
  register __m128d A0, B0, C0, F0, A1, B1, C1, F1, ALP;

  C0 = _mm_set1_pd(0.0);
  C1 = _mm_set1_pd(0.0);
  ALP = _mm_set1_pd(alpha);

  // unrolling of loops;
  for (k = nVP-1; k > 0; k -= 2) {
    A0 = _mm_load_pd(Ar);
    B0 = _mm_load_pd(b0);
    B1 = _mm_load_pd(b1);
    F0 = A0 * B0;
    F1 = A0 * B1;
    C0 = C0 + F0;
    C1 = C1 + F1;
    Ar += 2;
    b0 += 2;
    b1 += 2;
  }
  if (k == 0) {
    cval0   = Ar[0] * b0[0];
    cval1   = Ar[0] * b1[0];
    Cr0[0] += cval0 * alpha;
    Cr1[0] += cval1 * alpha;
  }
 update:
  C0 = C0 * ALP;
  C1 = C1 * ALP;
  F0 = _mm_hadd_pd(C0, C1);
  Cr0[0] += F0[0];
  Cr1[0] += F0[1];
}

INLINE void _inner_ddot_sse(double *Cr, const double *Ar,
                            const double *Br, double alpha, int nVP)
{
  register int k;
  register double cval = 0.0;
  register __m128d A0, B0, C0, F0, F1, ALP;

  C0 = _mm_set1_pd(0.0);
  ALP = _mm_set1_pd(alpha);

  // unrolling of loops;
  for (k = nVP-3; k > 0; k -= 4) {
    A0 = _mm_load_pd(Ar);
    B0 = _mm_load_pd(Br);
    F0 = A0 * B0;
    A0 = _mm_load_pd(&Ar[2]);
    B0 = _mm_load_pd(&Br[2]);
    C0 = C0 + F0;
    F1 = A0 * B0;
    C0 = C0 + F1;
    Ar += 4;
    Br += 4;
  }
  k += 3;
  switch (k) {
  case 3:
  case 2:
    A0 = _mm_load_pd(Ar);
    B0 = _mm_load_pd(Br);
    F0 = A0 * B0;
    C0 = C0 + F0;
    Ar += 2; Br += 2;
    if (k == 2)
      goto update;
  case 1:
    cval = Ar[0] * Br[0];
    Cr[0] += cval * alpha;
  }
 update:
  C0 = C0 * ALP;
  Cr[0] += (C0[0] + C0[1]);
}



// SSE4.1 versions with DPPD instruction (Dot Product of Packed Double)

// Update 4 columsn with 2 by 4 block. There are propably too many register variables
// here ...
INLINE void _inner_ddot4_2_sse4_1(double *Cr0, double *Cr1, double *Cr2, double *Cr3,
				  const double *a0, const double *a1,
				  const double *b0, const double *b1,
				  const double *b2, const double *b3,
				  double alpha, int nVP)
{
  register int k, i;
  register double cval0, cval1, cval2, cval3, cval4, cval5, cval6, cval7;
  register __m128d A0, B0, F0, B1, F1, B2, F2, B3, F3;

  cval0 = 0.0;
  cval1 = 0.0;
  cval2 = 0.0;
  cval3 = 0.0;
  cval4 = 0.0;
  cval5 = 0.0;
  cval6 = 0.0;
  cval7 = 0.0;

  // unrolling of loops;
  for (k = nVP-1; k > 0; k -= 2) {
    A0 = _mm_load_pd(a0);
    B0 = _mm_load_pd(b0);
    B1 = _mm_load_pd(b1);
    B2 = _mm_load_pd(b2);
    B3 = _mm_load_pd(b3);

    // 0x13: 3 = ddot of low and high elements from source, 1 = broadcast to low
    F0 = _mm_dp_pd(A0, B0, 0x13);
    F1 = _mm_dp_pd(A0, B1, 0x13);
    F2 = _mm_dp_pd(A0, B2, 0x13);
    F3 = _mm_dp_pd(A0, B3, 0x13);
    cval0 += F0[0];
    cval1 += F1[0];
    cval2 += F2[0];
    cval3 += F3[0];

    A0 = _mm_load_pd(a1);
    F0 = _mm_dp_pd(A0, B0, 0x13);
    F1 = _mm_dp_pd(A0, B1, 0x13);
    F2 = _mm_dp_pd(A0, B2, 0x13);
    F3 = _mm_dp_pd(A0, B3, 0x13);
    cval4 += F0[0];
    cval5 += F1[0];
    cval6 += F2[0];
    cval7 += F3[0];

    a0 += 2;
    a1 += 2;
    b0 += 2;
    b1 += 2;
    b2 += 2;
    b3 += 2;
  }
  if (k == 0) {
    cval0 += a0[0] * b0[0];
    cval1 += a0[0] * b1[0]; 
    cval2 += a0[0] * b2[0];
    cval3 += a0[0] * b3[0]; 
    cval4 += a1[0] * b0[0];
    cval5 += a1[0] * b1[0]; 
    cval6 += a1[0] * b2[0];
    cval7 += a1[0] * b3[0]; 
  }
  Cr0[0] += cval0 * alpha;
  Cr1[0] += cval1 * alpha;
  Cr2[0] += cval2 * alpha;
  Cr3[0] += cval3 * alpha;
  Cr0[1] += cval4 * alpha;
  Cr1[1] += cval5 * alpha;
  Cr2[1] += cval6 * alpha;
  Cr3[1] += cval7 * alpha;
}

// Update 4 columns with from 1 by 4 block.
INLINE void _inner_ddot4_sse4_1(double *Cr0, double *Cr1, double *Cr2, double *Cr3,
				const double *Ar, const double *b0, const double *b1,
				const double *b2, const double *b3, double alpha, int nVP)
{
  register int k, i;
  register double cval0, cval1, cval2, cval3;
  register __m128d A0, B0, F0, B1, F1, B2, F2, B3, F3;

  cval0 = 0.0; cval1 = 0.0; cval2 = 0.0; cval3 = 0.0;

  // unrolling of loops;
  for (k = nVP-1; k > 0; k -= 2) {
    A0 = _mm_load_pd(Ar);
    B0 = _mm_load_pd(b0);
    B1 = _mm_load_pd(b1);
    B2 = _mm_load_pd(b2);
    B3 = _mm_load_pd(b3);

    // 0x13: 3 = ddot of low and high elements from source, 1 = broadcast to low
    F0 = _mm_dp_pd(A0, B0, 0x13);
    F1 = _mm_dp_pd(A0, B1, 0x13);
    F2 = _mm_dp_pd(A0, B2, 0x13);
    F3 = _mm_dp_pd(A0, B3, 0x13);
    cval0 += F0[0];
    cval1 += F1[0];
    cval2 += F2[0];
    cval3 += F3[0];

    Ar += 2;
    b0 += 2;
    b1 += 2;
    b2 += 2;
    b3 += 2;
  }
  if (k != 0) {
    cval0 += Ar[0] * b0[0];
    cval1 += Ar[0] * b1[0]; 
    cval2 += Ar[0] * b2[0];
    cval3 += Ar[0] * b3[0]; 
  }
  Cr0[0] += cval0 * alpha;
  Cr1[0] += cval1 * alpha;
  Cr2[0] += cval2 * alpha;
  Cr3[0] += cval3 * alpha;
}

// Update 2 columsn with 2 by 2 block.
INLINE void _inner_ddot2_2_sse4_1(double *Cr0, double *Cr1, 
				  const double *a0, const double *a1,
				  const double *b0, const double *b1,
				  double alpha, int nVP)
{
  register int k, i;
  register double cval0, cval1, cval2, cval3;
  register __m128d A0, B0, F0, B1, F1, F2;

  cval0 = 0.0; cval1 = 0.0;

  // unrolling of loops;
  for (k = nVP-1; k > 0; k -= 2) {
    A0 = _mm_load_pd(a0);
    B0 = _mm_load_pd(b0);
    B1 = _mm_load_pd(b1);

    F0 = _mm_dp_pd(A0, B0, 0x13);
    F1 = _mm_dp_pd(A0, B1, 0x13);
    cval0 += F0[0];
    cval1 += F1[0];

    // 0x13: 3 = ddot of low and high elements from source, 1 = broadcast to low
    A0 = _mm_load_pd(a1);
    F2 = _mm_dp_pd(A0, B0, 0x13);
    F0 = _mm_dp_pd(A0, B1, 0x13);
    cval2 += F2[0];
    cval3 += F0[0];

    a0 += 2;
    a1 += 2;
    b0 += 2;
    b1 += 2;
  }
  if (k != 0) {
    cval0 += a0[0] * b0[0];
    cval1 += a0[0] * b1[0]; 
    cval2 += a1[0] * b0[0];
    cval3 += a1[0] * b1[0]; 
  }
  Cr0[0] += cval0 * alpha;
  Cr1[0] += cval1 * alpha;
  Cr0[1] += cval2 * alpha;
  Cr1[1] += cval3 * alpha;
}

// Update 2 columns with 1 by 2 block of length nVP.
INLINE void _inner_ddot2_sse4_1(double *Cr0, double *Cr1, 
				const double *Ar, const double *b0, const double *b1,
				double alpha, int nVP)
{
  register int k, i;
  register double cval0, cval1;
  register __m128d A0, B0, F0, B1, F1;

  cval0 = 0.0; cval1 = 0.0;

  // unrolling of loops;
  for (k = nVP-1; k > 0; k -= 2) {
    A0 = _mm_load_pd(Ar);
    B0 = _mm_load_pd(b0);
    B1 = _mm_load_pd(b1);

    // 0x13: 3 = ddot of low and high elements from source, 1 = broadcast to low
    F0 = _mm_dp_pd(A0, B0, 0x13);
    F1 = _mm_dp_pd(A0, B1, 0x13);
    cval0 += F0[0];
    cval1 += F1[0];

    Ar += 2;
    b0 += 2;
    b1 += 2;
  }
  if (k == 0) {
    cval0 += Ar[0] * b0[0];
    cval1 += Ar[0] * b1[0]; 
  }
  Cr0[0] += cval0 * alpha;
  Cr1[0] += cval1 * alpha;
}


// Non-vectorized versions.
INLINE void _inner_ddot4(double *c0, double *c1, double *c2, double *c3,
                         const double *Ar, const double *b0,
                         const double *b1, const double *b2, const double *b3,
                         double alpha, int nVP)
{
  register int k;
  register double cval0, cval1, cval2, cval3;
  register double cval4, cval5, cval6, cval7;

  cval0 = 0.0; cval1 = 0.0; cval2 = 0.0; cval3 = 0.0;
  cval4 = 0.0; cval5 = 0.0; cval6 = 0.0; cval7 = 0.0;

  // unrolling of loops;
  for (k = nVP-1; k > 0; k -= 2) {
    cval0 += Ar[0] * b0[0];
    cval1 += Ar[1] * b0[1];

    cval2 += Ar[0] * b1[0];
    cval3 += Ar[1] * b1[1];

    cval4 += Ar[0] * b2[0];
    cval5 += Ar[1] * b2[1];

    cval6 += Ar[0] * b3[0];
    cval7 += Ar[1] * b3[1];

    Ar += 2;
    b0 += 2;
    b1 += 2;
    b2 += 2;
    b3 += 2;
  }
  if (k == 0) {
    cval0 += Ar[0] * b0[0];
    cval2 += Ar[0] * b1[0];
    cval4 += Ar[0] * b2[1];
    cval6 += Ar[0] * b3[1];
  }
 update:
  c0[0] += (cval0 + cval1) * alpha;
  c1[0] += (cval2 + cval3) * alpha;
  c2[0] += (cval4 + cval5) * alpha;
  c3[0] += (cval6 + cval7) * alpha;
}

INLINE void _inner_ddot2(double *c0, double *c1, const double *Ar,
                         const double *b0, const double *b1, double alpha, int nVP)
{
  register int k;
  register double f0, f1, f2, f3, cval0, cval1, cval2, cval3;

  cval0 = 0.0; cval1 = 0.0; cval2 = 0.0; cval3 = 0.0;
  // unrolling of loops;
  for (k = 0; k < nVP-3; k += 4) {
    f0 = Ar[0] * b0[0];
    f1 = Ar[1] * b0[1];
    f2 = Ar[2] * b0[2];
    f3 = Ar[3] * b0[3];
    cval0 += f0;
    cval1 += f1;
    cval0 += f2;
    cval1 += f3;
    f0 = Ar[0] * b1[0];
    f1 = Ar[1] * b1[1];
    f2 = Ar[2] * b1[2];
    f3 = Ar[3] * b1[3];
    cval2 += f0;
    cval3 += f1;
    cval2 += f2;
    cval3 += f3;
    b0 += 4;
    b1 += 4;
    Ar += 4;
  }
  if (k == nVP)
    goto update;

  if (k < nVP-1) {
    f0 = Ar[0] * b0[0];
    f1 = Ar[1] * b0[1];
    cval0 += f0;
    cval1 += f1;
    f0 = Ar[0] * b1[0];
    f1 = Ar[1] * b1[1];
    cval2 += f0;
    cval3 += f1;
    b0 += 2;
    b1 += 2;
    Ar += 2;
    k += 2;
  }
  if (k < nVP) {
    cval0 += Ar[0] * b0[0];
    cval2 += Ar[0] * b1[0];
  }
 update:
  c0[0] += (cval0 + cval1) * alpha;
  c1[0] += (cval2 + cval3) * alpha;
}

INLINE void _inner_ddot(double *Cr, const double *Ar, const double *Br, double alpha, int nVP)
{
  register int k;
  register double f0, f1, f2, f3, cval0, cval1;

  cval0 = 0.0; cval1 = 0.0;
  // unrolling of loops;
  for (k = 0; k < nVP-3; k += 4) {
    f0 = Ar[0] * Br[0];
    f1 = Ar[1] * Br[1];
    f2 = Ar[2] * Br[2];
    f3 = Ar[3] * Br[3];
    cval0 += f0;
    cval1 += f1;
    cval0 += f2;
    cval1 += f3;
    Br += 4;
    Ar += 4;
  }
  if (k == nVP)
    goto update;

  if (k < nVP-1) {
    f0 = Ar[0] * Br[0];
    f1 = Ar[1] * Br[1];
    cval0 += f0;
    cval1 += f1;
    Br += 2;
    Ar += 2;
    k += 2;
  }
  if (k < nVP) {
    cval0 += Ar[0] * Br[0];
    Br++;
    Ar++;
    k++;
  }
 update:
  Cr[0] += (cval0 + cval1) * alpha;
}

#endif

// Local Variables:
// indent-tabs-mode: nil
// End:
