/**********************************************************************
 * Copyright (c) 2013-2015 Pieter Wuille                              *
 * Distributed under the MIT software license, see the accompanying   *
 * file COPYING or http://www.opensource.org/licenses/mit-license.php.*
 **********************************************************************/

#include "include/secp256k1.h"

#include "util.h"
#include "num_impl.h"
#include "field_impl.h"
#include "scalar_impl.h"
#include "group_impl.h"
#include "ecmult_impl.h"
#include "ecmult_const_impl.h"
#include "ecmult_gen_impl.h"
#include "ecdsa_impl.h"
#include "eckey_impl.h"
#include "hash_impl.h"

#define ARG_CHECK(cond) do { \
    if (EXPECT(!(cond), 0)) { \
        xsecp256k1_callback_call(&ctx->illegal_callback, #cond); \
        return 0; \
    } \
} while(0)

static void default_illegal_callback_fn(const char* str, void* data) {
    fprintf(stderr, "[libxsecp256k1] illegal argument: %s\n", str);
    abort();
}

static const xsecp256k1_callback default_illegal_callback = {
    default_illegal_callback_fn,
    NULL
};

static void default_error_callback_fn(const char* str, void* data) {
    fprintf(stderr, "[libxsecp256k1] internal consistency check failed: %s\n", str);
    abort();
}

static const xsecp256k1_callback default_error_callback = {
    default_error_callback_fn,
    NULL
};


struct xsecp256k1_context_struct {
    xsecp256k1_ecmult_context ecmult_ctx;
    xsecp256k1_ecmult_gen_context ecmult_gen_ctx;
    xsecp256k1_callback illegal_callback;
    xsecp256k1_callback error_callback;
};

xsecp256k1_context* xsecp256k1_context_create(unsigned int flags) {
    xsecp256k1_context* ret = (xsecp256k1_context*)checked_malloc(&default_error_callback, sizeof(xsecp256k1_context));
    ret->illegal_callback = default_illegal_callback;
    ret->error_callback = default_error_callback;

    if (EXPECT((flags & SECP256K1_FLAGS_TYPE_MASK) != SECP256K1_FLAGS_TYPE_CONTEXT, 0)) {
            xsecp256k1_callback_call(&ret->illegal_callback,
                                    "Invalid flags");
            free(ret);
            return NULL;
    }

    xsecp256k1_ecmult_context_init(&ret->ecmult_ctx);
    xsecp256k1_ecmult_gen_context_init(&ret->ecmult_gen_ctx);

    if (flags & SECP256K1_FLAGS_BIT_CONTEXT_SIGN) {
        xsecp256k1_ecmult_gen_context_build(&ret->ecmult_gen_ctx, &ret->error_callback);
    }
    if (flags & SECP256K1_FLAGS_BIT_CONTEXT_VERIFY) {
        xsecp256k1_ecmult_context_build(&ret->ecmult_ctx, &ret->error_callback);
    }

    return ret;
}

xsecp256k1_context* xsecp256k1_context_clone(const xsecp256k1_context* ctx) {
    xsecp256k1_context* ret = (xsecp256k1_context*)checked_malloc(&ctx->error_callback, sizeof(xsecp256k1_context));
    ret->illegal_callback = ctx->illegal_callback;
    ret->error_callback = ctx->error_callback;
    xsecp256k1_ecmult_context_clone(&ret->ecmult_ctx, &ctx->ecmult_ctx, &ctx->error_callback);
    xsecp256k1_ecmult_gen_context_clone(&ret->ecmult_gen_ctx, &ctx->ecmult_gen_ctx, &ctx->error_callback);
    return ret;
}

void xsecp256k1_context_destroy(xsecp256k1_context* ctx) {
    if (ctx != NULL) {
        xsecp256k1_ecmult_context_clear(&ctx->ecmult_ctx);
        xsecp256k1_ecmult_gen_context_clear(&ctx->ecmult_gen_ctx);

        free(ctx);
    }
}

void xsecp256k1_context_set_illegal_callback(xsecp256k1_context* ctx, void (*fun)(const char* message, void* data), const void* data) {
    if (fun == NULL) {
        fun = default_illegal_callback_fn;
    }
    ctx->illegal_callback.fn = fun;
    ctx->illegal_callback.data = data;
}

void xsecp256k1_context_set_error_callback(xsecp256k1_context* ctx, void (*fun)(const char* message, void* data), const void* data) {
    if (fun == NULL) {
        fun = default_error_callback_fn;
    }
    ctx->error_callback.fn = fun;
    ctx->error_callback.data = data;
}

static int xsecp256k1_pubkey_load(const xsecp256k1_context* ctx, xsecp256k1_ge* ge, const xsecp256k1_pubkey* pubkey) {
    if (sizeof(xsecp256k1_ge_storage) == 64) {
        /* When the xsecp256k1_ge_storage type is exactly 64 byte, use its
         * representation inside xsecp256k1_pubkey, as conversion is very fast.
         * Note that xsecp256k1_pubkey_save must use the same representation. */
        xsecp256k1_ge_storage s;
        memcpy(&s, &pubkey->data[0], 64);
        xsecp256k1_ge_from_storage(ge, &s);
    } else {
        /* Otherwise, fall back to 32-byte big endian for X and Y. */
        xsecp256k1_fe x, y;
        xsecp256k1_fe_set_b32(&x, pubkey->data);
        xsecp256k1_fe_set_b32(&y, pubkey->data + 32);
        xsecp256k1_ge_set_xy(ge, &x, &y);
    }
    ARG_CHECK(!xsecp256k1_fe_is_zero(&ge->x));
    return 1;
}

static void xsecp256k1_pubkey_save(xsecp256k1_pubkey* pubkey, xsecp256k1_ge* ge) {
    if (sizeof(xsecp256k1_ge_storage) == 64) {
        xsecp256k1_ge_storage s;
        xsecp256k1_ge_to_storage(&s, ge);
        memcpy(&pubkey->data[0], &s, 64);
    } else {
        VERIFY_CHECK(!xsecp256k1_ge_is_infinity(ge));
        xsecp256k1_fe_normalize_var(&ge->x);
        xsecp256k1_fe_normalize_var(&ge->y);
        xsecp256k1_fe_get_b32(pubkey->data, &ge->x);
        xsecp256k1_fe_get_b32(pubkey->data + 32, &ge->y);
    }
}

int xsecp256k1_ec_pubkey_parse(const xsecp256k1_context* ctx, xsecp256k1_pubkey* pubkey, const unsigned char *input, size_t inputlen) {
    xsecp256k1_ge Q;

    VERIFY_CHECK(ctx != NULL);
    ARG_CHECK(pubkey != NULL);
    memset(pubkey, 0, sizeof(*pubkey));
    ARG_CHECK(input != NULL);
    if (!xsecp256k1_eckey_pubkey_parse(&Q, input, inputlen)) {
        return 0;
    }
    xsecp256k1_pubkey_save(pubkey, &Q);
    xsecp256k1_ge_clear(&Q);
    return 1;
}

int xsecp256k1_ec_pubkey_serialize(const xsecp256k1_context* ctx, unsigned char *output, size_t *outputlen, const xsecp256k1_pubkey* pubkey, unsigned int flags) {
    xsecp256k1_ge Q;
    size_t len;
    int ret = 0;

    VERIFY_CHECK(ctx != NULL);
    ARG_CHECK(outputlen != NULL);
    ARG_CHECK(*outputlen >= ((flags & SECP256K1_FLAGS_BIT_COMPRESSION) ? 33 : 65));
    len = *outputlen;
    *outputlen = 0;
    ARG_CHECK(output != NULL);
    memset(output, 0, len);
    ARG_CHECK(pubkey != NULL);
    ARG_CHECK((flags & SECP256K1_FLAGS_TYPE_MASK) == SECP256K1_FLAGS_TYPE_COMPRESSION);
    if (xsecp256k1_pubkey_load(ctx, &Q, pubkey)) {
        ret = xsecp256k1_eckey_pubkey_serialize(&Q, output, &len, flags & SECP256K1_FLAGS_BIT_COMPRESSION);
        if (ret) {
            *outputlen = len;
        }
    }
    return ret;
}

static void xsecp256k1_ecdsa_signature_load(const xsecp256k1_context* ctx, xsecp256k1_scalar* r, xsecp256k1_scalar* s, const xsecp256k1_ecdsa_signature* sig) {
    (void)ctx;
    if (sizeof(xsecp256k1_scalar) == 32) {
        /* When the xsecp256k1_scalar type is exactly 32 byte, use its
         * representation inside xsecp256k1_ecdsa_signature, as conversion is very fast.
         * Note that xsecp256k1_ecdsa_signature_save must use the same representation. */
        memcpy(r, &sig->data[0], 32);
        memcpy(s, &sig->data[32], 32);
    } else {
        xsecp256k1_scalar_set_b32(r, &sig->data[0], NULL);
        xsecp256k1_scalar_set_b32(s, &sig->data[32], NULL);
    }
}

static void xsecp256k1_ecdsa_signature_save(xsecp256k1_ecdsa_signature* sig, const xsecp256k1_scalar* r, const xsecp256k1_scalar* s) {
    if (sizeof(xsecp256k1_scalar) == 32) {
        memcpy(&sig->data[0], r, 32);
        memcpy(&sig->data[32], s, 32);
    } else {
        xsecp256k1_scalar_get_b32(&sig->data[0], r);
        xsecp256k1_scalar_get_b32(&sig->data[32], s);
    }
}

int xsecp256k1_ecdsa_signature_parse_der(const xsecp256k1_context* ctx, xsecp256k1_ecdsa_signature* sig, const unsigned char *input, size_t inputlen) {
    xsecp256k1_scalar r, s;

    VERIFY_CHECK(ctx != NULL);
    ARG_CHECK(sig != NULL);
    ARG_CHECK(input != NULL);

    if (xsecp256k1_ecdsa_sig_parse(&r, &s, input, inputlen)) {
        xsecp256k1_ecdsa_signature_save(sig, &r, &s);
        return 1;
    } else {
        memset(sig, 0, sizeof(*sig));
        return 0;
    }
}

int xsecp256k1_ecdsa_signature_parse_compact(const xsecp256k1_context* ctx, xsecp256k1_ecdsa_signature* sig, const unsigned char *input64) {
    xsecp256k1_scalar r, s;
    int ret = 1;
    int overflow = 0;

    VERIFY_CHECK(ctx != NULL);
    ARG_CHECK(sig != NULL);
    ARG_CHECK(input64 != NULL);

    xsecp256k1_scalar_set_b32(&r, &input64[0], &overflow);
    ret &= !overflow;
    xsecp256k1_scalar_set_b32(&s, &input64[32], &overflow);
    ret &= !overflow;
    if (ret) {
        xsecp256k1_ecdsa_signature_save(sig, &r, &s);
    } else {
        memset(sig, 0, sizeof(*sig));
    }
    return ret;
}

int xsecp256k1_ecdsa_signature_serialize_der(const xsecp256k1_context* ctx, unsigned char *output, size_t *outputlen, const xsecp256k1_ecdsa_signature* sig) {
    xsecp256k1_scalar r, s;

    VERIFY_CHECK(ctx != NULL);
    ARG_CHECK(output != NULL);
    ARG_CHECK(outputlen != NULL);
    ARG_CHECK(sig != NULL);

    xsecp256k1_ecdsa_signature_load(ctx, &r, &s, sig);
    return xsecp256k1_ecdsa_sig_serialize(output, outputlen, &r, &s);
}

int xsecp256k1_ecdsa_signature_serialize_compact(const xsecp256k1_context* ctx, unsigned char *output64, const xsecp256k1_ecdsa_signature* sig) {
    xsecp256k1_scalar r, s;

    VERIFY_CHECK(ctx != NULL);
    ARG_CHECK(output64 != NULL);
    ARG_CHECK(sig != NULL);

    xsecp256k1_ecdsa_signature_load(ctx, &r, &s, sig);
    xsecp256k1_scalar_get_b32(&output64[0], &r);
    xsecp256k1_scalar_get_b32(&output64[32], &s);
    return 1;
}

int xsecp256k1_ecdsa_signature_normalize(const xsecp256k1_context* ctx, xsecp256k1_ecdsa_signature *sigout, const xsecp256k1_ecdsa_signature *sigin) {
    xsecp256k1_scalar r, s;
    int ret = 0;

    VERIFY_CHECK(ctx != NULL);
    ARG_CHECK(sigin != NULL);

    xsecp256k1_ecdsa_signature_load(ctx, &r, &s, sigin);
    ret = xsecp256k1_scalar_is_high(&s);
    if (sigout != NULL) {
        if (ret) {
            xsecp256k1_scalar_negate(&s, &s);
        }
        xsecp256k1_ecdsa_signature_save(sigout, &r, &s);
    }

    return ret;
}

int xsecp256k1_ecdsa_verify(const xsecp256k1_context* ctx, const xsecp256k1_ecdsa_signature *sig, const unsigned char *msg32, const xsecp256k1_pubkey *pubkey) {
    xsecp256k1_ge q;
    xsecp256k1_scalar r, s;
    xsecp256k1_scalar m;
    VERIFY_CHECK(ctx != NULL);
    ARG_CHECK(xsecp256k1_ecmult_context_is_built(&ctx->ecmult_ctx));
    ARG_CHECK(msg32 != NULL);
    ARG_CHECK(sig != NULL);
    ARG_CHECK(pubkey != NULL);

    xsecp256k1_scalar_set_b32(&m, msg32, NULL);
    xsecp256k1_ecdsa_signature_load(ctx, &r, &s, sig);
    return (!xsecp256k1_scalar_is_high(&s) &&
            xsecp256k1_pubkey_load(ctx, &q, pubkey) &&
            xsecp256k1_ecdsa_sig_verify(&ctx->ecmult_ctx, &r, &s, &q, &m));
}

static int nonce_function_rfc6979(unsigned char *nonce32, const unsigned char *msg32, const unsigned char *key32, const unsigned char *algo16, void *data, unsigned int counter) {
   unsigned char keydata[112];
   int keylen = 64;
   xsecp256k1_rfc6979_hmac_sha256_t rng;
   unsigned int i;
   /* We feed a byte array to the PRNG as input, consisting of:
    * - the private key (32 bytes) and message (32 bytes), see RFC 6979 3.2d.
    * - optionally 32 extra bytes of data, see RFC 6979 3.6 Additional Data.
    * - optionally 16 extra bytes with the algorithm name.
    * Because the arguments have distinct fixed lengths it is not possible for
    *  different argument mixtures to emulate each other and result in the same
    *  nonces.
    */
   memcpy(keydata, key32, 32);
   memcpy(keydata + 32, msg32, 32);
   if (data != NULL) {
       memcpy(keydata + 64, data, 32);
       keylen = 96;
   }
   if (algo16 != NULL) {
       memcpy(keydata + keylen, algo16, 16);
       keylen += 16;
   }
   xsecp256k1_rfc6979_hmac_sha256_initialize(&rng, keydata, keylen);
   memset(keydata, 0, sizeof(keydata));
   for (i = 0; i <= counter; i++) {
       xsecp256k1_rfc6979_hmac_sha256_generate(&rng, nonce32, 32);
   }
   xsecp256k1_rfc6979_hmac_sha256_finalize(&rng);
   return 1;
}

const xsecp256k1_nonce_function xsecp256k1_nonce_function_rfc6979 = nonce_function_rfc6979;
const xsecp256k1_nonce_function xsecp256k1_nonce_function_default = nonce_function_rfc6979;

int xsecp256k1_ecdsa_sign(const xsecp256k1_context* ctx, xsecp256k1_ecdsa_signature *signature, const unsigned char *msg32, const unsigned char *seckey, xsecp256k1_nonce_function noncefp, const void* noncedata) {
    xsecp256k1_scalar r, s;
    xsecp256k1_scalar sec, non, msg;
    int ret = 0;
    int overflow = 0;
    VERIFY_CHECK(ctx != NULL);
    ARG_CHECK(xsecp256k1_ecmult_gen_context_is_built(&ctx->ecmult_gen_ctx));
    ARG_CHECK(msg32 != NULL);
    ARG_CHECK(signature != NULL);
    ARG_CHECK(seckey != NULL);
    if (noncefp == NULL) {
        noncefp = xsecp256k1_nonce_function_default;
    }

    xsecp256k1_scalar_set_b32(&sec, seckey, &overflow);
    /* Fail if the secret key is invalid. */
    if (!overflow && !xsecp256k1_scalar_is_zero(&sec)) {
        unsigned char nonce32[32];
        unsigned int count = 0;
        xsecp256k1_scalar_set_b32(&msg, msg32, NULL);
        while (1) {
            ret = noncefp(nonce32, msg32, seckey, NULL, (void*)noncedata, count);
            if (!ret) {
                break;
            }
            xsecp256k1_scalar_set_b32(&non, nonce32, &overflow);
            if (!overflow && !xsecp256k1_scalar_is_zero(&non)) {
                if (xsecp256k1_ecdsa_sig_sign(&ctx->ecmult_gen_ctx, &r, &s, &sec, &msg, &non, NULL)) {
                    break;
                }
            }
            count++;
        }
        memset(nonce32, 0, 32);
        xsecp256k1_scalar_clear(&msg);
        xsecp256k1_scalar_clear(&non);
        xsecp256k1_scalar_clear(&sec);
    }
    if (ret) {
        xsecp256k1_ecdsa_signature_save(signature, &r, &s);
    } else {
        memset(signature, 0, sizeof(*signature));
    }
    return ret;
}

int xsecp256k1_ec_seckey_verify(const xsecp256k1_context* ctx, const unsigned char *seckey) {
    xsecp256k1_scalar sec;
    int ret;
    int overflow;
    VERIFY_CHECK(ctx != NULL);
    ARG_CHECK(seckey != NULL);

    xsecp256k1_scalar_set_b32(&sec, seckey, &overflow);
    ret = !overflow && !xsecp256k1_scalar_is_zero(&sec);
    xsecp256k1_scalar_clear(&sec);
    return ret;
}

int xsecp256k1_ec_pubkey_create(const xsecp256k1_context* ctx, xsecp256k1_pubkey *pubkey, const unsigned char *seckey) {
    xsecp256k1_gej pj;
    xsecp256k1_ge p;
    xsecp256k1_scalar sec;
    int overflow;
    int ret = 0;
    VERIFY_CHECK(ctx != NULL);
    ARG_CHECK(pubkey != NULL);
    memset(pubkey, 0, sizeof(*pubkey));
    ARG_CHECK(xsecp256k1_ecmult_gen_context_is_built(&ctx->ecmult_gen_ctx));
    ARG_CHECK(seckey != NULL);

    xsecp256k1_scalar_set_b32(&sec, seckey, &overflow);
    ret = (!overflow) & (!xsecp256k1_scalar_is_zero(&sec));
    if (ret) {
        xsecp256k1_ecmult_gen(&ctx->ecmult_gen_ctx, &pj, &sec);
        xsecp256k1_ge_set_gej(&p, &pj);
        xsecp256k1_pubkey_save(pubkey, &p);
    }
    xsecp256k1_scalar_clear(&sec);
    return ret;
}

int xsecp256k1_ec_privkey_tweak_add(const xsecp256k1_context* ctx, unsigned char *seckey, const unsigned char *tweak) {
    xsecp256k1_scalar term;
    xsecp256k1_scalar sec;
    int ret = 0;
    int overflow = 0;
    VERIFY_CHECK(ctx != NULL);
    ARG_CHECK(seckey != NULL);
    ARG_CHECK(tweak != NULL);

    xsecp256k1_scalar_set_b32(&term, tweak, &overflow);
    xsecp256k1_scalar_set_b32(&sec, seckey, NULL);

    ret = !overflow && xsecp256k1_eckey_privkey_tweak_add(&sec, &term);
    memset(seckey, 0, 32);
    if (ret) {
        xsecp256k1_scalar_get_b32(seckey, &sec);
    }

    xsecp256k1_scalar_clear(&sec);
    xsecp256k1_scalar_clear(&term);
    return ret;
}

int xsecp256k1_ec_pubkey_tweak_add(const xsecp256k1_context* ctx, xsecp256k1_pubkey *pubkey, const unsigned char *tweak) {
    xsecp256k1_ge p;
    xsecp256k1_scalar term;
    int ret = 0;
    int overflow = 0;
    VERIFY_CHECK(ctx != NULL);
    ARG_CHECK(xsecp256k1_ecmult_context_is_built(&ctx->ecmult_ctx));
    ARG_CHECK(pubkey != NULL);
    ARG_CHECK(tweak != NULL);

    xsecp256k1_scalar_set_b32(&term, tweak, &overflow);
    ret = !overflow && xsecp256k1_pubkey_load(ctx, &p, pubkey);
    memset(pubkey, 0, sizeof(*pubkey));
    if (ret) {
        if (xsecp256k1_eckey_pubkey_tweak_add(&ctx->ecmult_ctx, &p, &term)) {
            xsecp256k1_pubkey_save(pubkey, &p);
        } else {
            ret = 0;
        }
    }

    return ret;
}

int xsecp256k1_ec_privkey_tweak_mul(const xsecp256k1_context* ctx, unsigned char *seckey, const unsigned char *tweak) {
    xsecp256k1_scalar factor;
    xsecp256k1_scalar sec;
    int ret = 0;
    int overflow = 0;
    VERIFY_CHECK(ctx != NULL);
    ARG_CHECK(seckey != NULL);
    ARG_CHECK(tweak != NULL);

    xsecp256k1_scalar_set_b32(&factor, tweak, &overflow);
    xsecp256k1_scalar_set_b32(&sec, seckey, NULL);
    ret = !overflow && xsecp256k1_eckey_privkey_tweak_mul(&sec, &factor);
    memset(seckey, 0, 32);
    if (ret) {
        xsecp256k1_scalar_get_b32(seckey, &sec);
    }

    xsecp256k1_scalar_clear(&sec);
    xsecp256k1_scalar_clear(&factor);
    return ret;
}

int xsecp256k1_ec_pubkey_tweak_mul(const xsecp256k1_context* ctx, xsecp256k1_pubkey *pubkey, const unsigned char *tweak) {
    xsecp256k1_ge p;
    xsecp256k1_scalar factor;
    int ret = 0;
    int overflow = 0;
    VERIFY_CHECK(ctx != NULL);
    ARG_CHECK(xsecp256k1_ecmult_context_is_built(&ctx->ecmult_ctx));
    ARG_CHECK(pubkey != NULL);
    ARG_CHECK(tweak != NULL);

    xsecp256k1_scalar_set_b32(&factor, tweak, &overflow);
    ret = !overflow && xsecp256k1_pubkey_load(ctx, &p, pubkey);
    memset(pubkey, 0, sizeof(*pubkey));
    if (ret) {
        if (xsecp256k1_eckey_pubkey_tweak_mul(&ctx->ecmult_ctx, &p, &factor)) {
            xsecp256k1_pubkey_save(pubkey, &p);
        } else {
            ret = 0;
        }
    }

    return ret;
}

int xsecp256k1_context_randomize(xsecp256k1_context* ctx, const unsigned char *seed32) {
    VERIFY_CHECK(ctx != NULL);
    ARG_CHECK(xsecp256k1_ecmult_gen_context_is_built(&ctx->ecmult_gen_ctx));
    xsecp256k1_ecmult_gen_blind(&ctx->ecmult_gen_ctx, seed32);
    return 1;
}

int xsecp256k1_ec_pubkey_combine(const xsecp256k1_context* ctx, xsecp256k1_pubkey *pubnonce, const xsecp256k1_pubkey * const *pubnonces, size_t n) {
    size_t i;
    xsecp256k1_gej Qj;
    xsecp256k1_ge Q;

    ARG_CHECK(pubnonce != NULL);
    memset(pubnonce, 0, sizeof(*pubnonce));
    ARG_CHECK(n >= 1);
    ARG_CHECK(pubnonces != NULL);

    xsecp256k1_gej_set_infinity(&Qj);

    for (i = 0; i < n; i++) {
        xsecp256k1_pubkey_load(ctx, &Q, pubnonces[i]);
        xsecp256k1_gej_add_ge(&Qj, &Qj, &Q);
    }
    if (xsecp256k1_gej_is_infinity(&Qj)) {
        return 0;
    }
    xsecp256k1_ge_set_gej(&Q, &Qj);
    xsecp256k1_pubkey_save(pubnonce, &Q);
    return 1;
}

#ifdef ENABLE_MODULE_ECDH
# include "modules/ecdh/main_impl.h"
#endif

#ifdef ENABLE_MODULE_SCHNORR
# include "modules/schnorr/main_impl.h"
#endif

#ifdef ENABLE_MODULE_RECOVERY
# include "modules/recovery/main_impl.h"
#endif
