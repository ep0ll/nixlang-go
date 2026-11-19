#include <stdlib.h>
#include "bridge.h"

void * nixgo_expr_id_to_userdata(uintptr_t id)
{
    (void) id;
    return malloc(1);
}

void nixgo_expr_userdata_free(void * userdata)
{
    free(userdata);
}
