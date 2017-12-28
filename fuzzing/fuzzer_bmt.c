#include <stdlib.h>
#include <stdio.h>
#include <stdint.h>
#include <string.h>
#include "bmt.h"

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
    run_bmt(p, (int)mode);
    return (int)GoCalcCoverage();
}
