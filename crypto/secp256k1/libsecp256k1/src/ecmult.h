/**********************************************************************
 * Copyright (c) 2013, 2014 Pieter Wuille                             *
 * Distributed under the MIT software license, see the accompanying   *
 * file COPYING or http://www.opensource.org/licenses/mit-license.php.*
 **********************************************************************/

#ifndef _SECP256K1_ECMULT_
#define _SECP256K1_ECMULT_

#include "num.h"
#include "group.h"

typedef struct {
    /* For accelerating the computation of a*P + b*G: */
    xsecp256k1_ge_storage (*pre_g)[];    /* odd multiples of the generator */
#ifdef USE_ENDOMORPHISM
    xsecp256k1_ge_storage (*pre_g_128)[]; /* odd multiples of 2^128*generator */
#endif
} xsecp256k1_ecmult_context;

static void xsecp256k1_ecmult_context_init(xsecp256k1_ecmult_context *ctx);
static void xsecp256k1_ecmult_context_build(xsecp256k1_ecmult_context *ctx, const xsecp256k1_callback *cb);
static void xsecp256k1_ecmult_context_clone(xsecp256k1_ecmult_context *dst,
                                           const xsecp256k1_ecmult_context *src, const xsecp256k1_callback *cb);
static void xsecp256k1_ecmult_context_clear(xsecp256k1_ecmult_context *ctx);
static int xsecp256k1_ecmult_context_is_built(const xsecp256k1_ecmult_context *ctx);

/** Double multiply: R = na*A + ng*G */
static void xsecp256k1_ecmult(const xsecp256k1_ecmult_context *ctx, xsecp256k1_gej *r, const xsecp256k1_gej *a, const xsecp256k1_scalar *na, const xsecp256k1_scalar *ng);

#endif
