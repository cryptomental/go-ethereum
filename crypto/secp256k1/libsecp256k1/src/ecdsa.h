/**********************************************************************
 * Copyright (c) 2013, 2014 Pieter Wuille                             *
 * Distributed under the MIT software license, see the accompanying   *
 * file COPYING or http://www.opensource.org/licenses/mit-license.php.*
 **********************************************************************/

#ifndef _SECP256K1_ECDSA_
#define _SECP256K1_ECDSA_

#include <stddef.h>

#include "scalar.h"
#include "group.h"
#include "ecmult.h"

static int xsecp256k1_ecdsa_sig_parse(xsecp256k1_scalar *r, xsecp256k1_scalar *s, const unsigned char *sig, size_t size);
static int xsecp256k1_ecdsa_sig_serialize(unsigned char *sig, size_t *size, const xsecp256k1_scalar *r, const xsecp256k1_scalar *s);
static int xsecp256k1_ecdsa_sig_verify(const xsecp256k1_ecmult_context *ctx, const xsecp256k1_scalar* r, const xsecp256k1_scalar* s, const xsecp256k1_ge *pubkey, const xsecp256k1_scalar *message);
static int xsecp256k1_ecdsa_sig_sign(const xsecp256k1_ecmult_gen_context *ctx, xsecp256k1_scalar* r, xsecp256k1_scalar* s, const xsecp256k1_scalar *seckey, const xsecp256k1_scalar *message, const xsecp256k1_scalar *nonce, int *recid);

#endif
