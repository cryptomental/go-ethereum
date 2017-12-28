#include <stdlib.h>
#include <stdio.h>
#include <stdint.h>
#include <string.h>
#include "ecies.h"

int LLVMFuzzerTestOneInput(const uint8_t *data, size_t size)
{
    GoSlice p = {(void*)data, size, size};
    unsigned char mode;

    if ( size < 1 ) {
        return 0;
    }

    mode = *data;
    data++; size--;

    GoResetCoverage();
    run_ecies(p, (int)mode);
    return (int)GoCalcCoverage();
}
