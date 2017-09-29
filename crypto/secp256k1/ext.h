// Copyright 2015 Jeffrey Wilcke, Felix Lange, Gustav Simonsson. All rights reserved.
// Use of this source code is governed by a BSD-style license that can be found in
// the LICENSE file.

// xsecp256k1_context_create_sign_verify creates a context for signing and signature verification.
static xsecp256k1_context* xsecp256k1_context_create_sign_verify() {
	return xsecp256k1_context_create(SECP256K1_CONTEXT_SIGN | SECP256K1_CONTEXT_VERIFY);
}

// xsecp256k1_ecdsa_recover_pubkey recovers the public key of an encoded compact signature.
//
// Returns: 1: recovery was successful
//          0: recovery was not successful
// Args:    ctx:        pointer to a context object (cannot be NULL)
//  Out:    pubkey_out: the serialized 65-byte public key of the signer (cannot be NULL)
//  In:     sigdata:    pointer to a 65-byte signature with the recovery id at the end (cannot be NULL)
//          msgdata:    pointer to a 32-byte message (cannot be NULL)
static int xsecp256k1_ecdsa_recover_pubkey(
	const xsecp256k1_context* ctx,
	unsigned char *pubkey_out,
	const unsigned char *sigdata,
	const unsigned char *msgdata
) {
	xsecp256k1_ecdsa_recoverable_signature sig;
	xsecp256k1_pubkey pubkey;

	if (!xsecp256k1_ecdsa_recoverable_signature_parse_compact(ctx, &sig, sigdata, (int)sigdata[64])) {
		return 0;
	}
	if (!xsecp256k1_ecdsa_recover(ctx, &pubkey, &sig, msgdata)) {
		return 0;
	}
	size_t outputlen = 65;
	return xsecp256k1_ec_pubkey_serialize(ctx, pubkey_out, &outputlen, &pubkey, SECP256K1_EC_UNCOMPRESSED);
}

// xsecp256k1_pubkey_scalar_mul multiplies a point by a scalar in constant time.
//
// Returns: 1: multiplication was successful
//          0: scalar was invalid (zero or overflow)
// Args:    ctx:      pointer to a context object (cannot be NULL)
//  Out:    point:    the multiplied point (usually secret)
//  In:     point:    pointer to a 64-byte public point,
//                    encoded as two 256bit big-endian numbers.
//          scalar:   a 32-byte scalar with which to multiply the point
int xsecp256k1_pubkey_scalar_mul(const xsecp256k1_context* ctx, unsigned char *point, const unsigned char *scalar) {
	int ret = 0;
	int overflow = 0;
	xsecp256k1_fe feX, feY;
	xsecp256k1_gej res;
	xsecp256k1_ge ge;
	xsecp256k1_scalar s;
	ARG_CHECK(point != NULL);
	ARG_CHECK(scalar != NULL);
	(void)ctx;

	xsecp256k1_fe_set_b32(&feX, point);
	xsecp256k1_fe_set_b32(&feY, point+32);
	xsecp256k1_ge_set_xy(&ge, &feX, &feY);
	xsecp256k1_scalar_set_b32(&s, scalar, &overflow);
	if (overflow || xsecp256k1_scalar_is_zero(&s)) {
		ret = 0;
	} else {
		xsecp256k1_ecmult_const(&res, &ge, &s);
		xsecp256k1_ge_set_gej(&ge, &res);
		/* Note: can't use xsecp256k1_pubkey_save here because it is not constant time. */
		xsecp256k1_fe_normalize(&ge.x);
		xsecp256k1_fe_normalize(&ge.y);
		xsecp256k1_fe_get_b32(point, &ge.x);
		xsecp256k1_fe_get_b32(point+32, &ge.y);
		ret = 1;
	}
	xsecp256k1_scalar_clear(&s);
	return ret;
}
