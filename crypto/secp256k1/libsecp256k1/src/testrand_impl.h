/**********************************************************************
 * Copyright (c) 2013-2015 Pieter Wuille                              *
 * Distributed under the MIT software license, see the accompanying   *
 * file COPYING or http://www.opensource.org/licenses/mit-license.php.*
 **********************************************************************/

#ifndef _SECP256K1_TESTRAND_IMPL_H_
#define _SECP256K1_TESTRAND_IMPL_H_

#include <stdint.h>
#include <string.h>

#include "testrand.h"
#include "hash.h"

static xsecp256k1_rfc6979_hmac_sha256_t xsecp256k1_test_rng;
static uint32_t xsecp256k1_test_rng_precomputed[8];
static int xsecp256k1_test_rng_precomputed_used = 8;
static uint64_t xsecp256k1_test_rng_integer;
static int xsecp256k1_test_rng_integer_bits_left = 0;

SECP256K1_INLINE static void xsecp256k1_rand_seed(const unsigned char *seed16) {
    xsecp256k1_rfc6979_hmac_sha256_initialize(&xsecp256k1_test_rng, seed16, 16);
}

SECP256K1_INLINE static uint32_t xsecp256k1_rand32(void) {
    if (xsecp256k1_test_rng_precomputed_used == 8) {
        xsecp256k1_rfc6979_hmac_sha256_generate(&xsecp256k1_test_rng, (unsigned char*)(&xsecp256k1_test_rng_precomputed[0]), sizeof(xsecp256k1_test_rng_precomputed));
        xsecp256k1_test_rng_precomputed_used = 0;
    }
    return xsecp256k1_test_rng_precomputed[xsecp256k1_test_rng_precomputed_used++];
}

static uint32_t xsecp256k1_rand_bits(int bits) {
    uint32_t ret;
    if (xsecp256k1_test_rng_integer_bits_left < bits) {
        xsecp256k1_test_rng_integer |= (((uint64_t)xsecp256k1_rand32()) << xsecp256k1_test_rng_integer_bits_left);
        xsecp256k1_test_rng_integer_bits_left += 32;
    }
    ret = xsecp256k1_test_rng_integer;
    xsecp256k1_test_rng_integer >>= bits;
    xsecp256k1_test_rng_integer_bits_left -= bits;
    ret &= ((~((uint32_t)0)) >> (32 - bits));
    return ret;
}

static uint32_t xsecp256k1_rand_int(uint32_t range) {
    /* We want a uniform integer between 0 and range-1, inclusive.
     * B is the smallest number such that range <= 2**B.
     * two mechanisms implemented here:
     * - generate B bits numbers until one below range is found, and return it
     * - find the largest multiple M of range that is <= 2**(B+A), generate B+A
     *   bits numbers until one below M is found, and return it modulo range
     * The second mechanism consumes A more bits of entropy in every iteration,
     * but may need fewer iterations due to M being closer to 2**(B+A) then
     * range is to 2**B. The array below (indexed by B) contains a 0 when the
     * first mechanism is to be used, and the number A otherwise.
     */
    static const int addbits[] = {0, 0, 0, 0, 0, 0, 0, 0, 0, 1, 2, 2, 2, 2, 2, 2, 2, 2, 2, 3, 3, 3, 3, 3, 3, 3, 3, 3, 3, 3, 2, 1, 0};
    uint32_t trange, mult;
    int bits = 0;
    if (range <= 1) {
        return 0;
    }
    trange = range - 1;
    while (trange > 0) {
        trange >>= 1;
        bits++;
    }
    if (addbits[bits]) {
        bits = bits + addbits[bits];
        mult = ((~((uint32_t)0)) >> (32 - bits)) / range;
        trange = range * mult;
    } else {
        trange = range;
        mult = 1;
    }
    while(1) {
        uint32_t x = xsecp256k1_rand_bits(bits);
        if (x < trange) {
            return (mult == 1) ? x : (x % range);
        }
    }
}

static void xsecp256k1_rand256(unsigned char *b32) {
    xsecp256k1_rfc6979_hmac_sha256_generate(&xsecp256k1_test_rng, b32, 32);
}

static void xsecp256k1_rand_bytes_test(unsigned char *bytes, size_t len) {
    size_t bits = 0;
    memset(bytes, 0, len);
    while (bits < len * 8) {
        int now;
        uint32_t val;
        now = 1 + (xsecp256k1_rand_bits(6) * xsecp256k1_rand_bits(5) + 16) / 31;
        val = xsecp256k1_rand_bits(1);
        while (now > 0 && bits < len * 8) {
            bytes[bits / 8] |= val << (bits % 8);
            now--;
            bits++;
        }
    }
}

static void xsecp256k1_rand256_test(unsigned char *b32) {
    xsecp256k1_rand_bytes_test(b32, 32);
}

#endif
