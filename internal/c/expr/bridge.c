#include <stdlib.h>
#include "bridge.h"

void * nixgo_expr_id_to_userdata(uintptr_t id)
{
    uintptr_t *token = malloc(sizeof(uintptr_t));
    if (token != NULL)
        *token = id;
    return token;
}

uintptr_t nixgo_expr_userdata_take_id(void * userdata)
{
    if (userdata == NULL)
        return 0;
    uintptr_t id = *(uintptr_t *) userdata;
    free(userdata);
    return id;
}
