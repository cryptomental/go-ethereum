/**********************************************************************
 * Copyright (c) 2014 Pieter Wuille                                   *
 * Distributed under the MIT software license, see the accompanying   *
 * file COPYING or http://www.opensource.org/licenses/mit-license.php.*
 **********************************************************************/

#ifndef _SECP256K1_HASH_
#define _SECP256K1_HASH_

#include <stdlib.h>
#include <stdint.h>

typedef struct {
    uint32_t s[8];
    uint32_t buf[16]; /* In big endian */
    size_t bytes;
} xsecp256k1_sha256_t;

static void xsecp256k1_sha256_initialize(xsecp256k1_sha256_t *hash);
static void xsecp256k1_sha256_write(xsecp256k1_sha256_t *hash, const unsigned char *data, size_t size);
static void xsecp256k1_sha256_finalize(xsecp256k1_sha256_t *hash, unsigned char *out32);

typedef struct {
    xsecp256k1_sha256_t inner, outer;
} xsecp256k1_hmac_sha256_t;

static void xsecp256k1_hmac_sha256_initialize(xsecp256k1_hmac_sha256_t *hash, const unsigned char *key, size_t size);
static void xsecp256k1_hmac_sha256_write(xsecp256k1_hmac_sha256_t *hash, const unsigned char *data, size_t size);
static void xsecp256k1_hmac_sha256_finalize(xsecp256k1_hmac_sha256_t *hash, unsigned char *out32);

typedef struct {
    unsigned char v[32];
    unsigned char k[32];
    int retry;
} xsecp256k1_rfc6979_hmac_sha256_t;

static void xsecp256k1_rfc6979_hmac_sha256_initialize(xsecp256k1_rfc6979_hmac_sha256_t *rng, const unsigned char *key, size_t keylen);
static void xsecp256k1_rfc6979_hmac_sha256_generate(xsecp256k1_rfc6979_hmac_sha256_t *rng, unsigned char *out, size_t outlen);
static void xsecp256k1_rfc6979_hmac_sha256_finalize(xsecp256k1_rfc6979_hmac_sha256_t *rng);

#endif
