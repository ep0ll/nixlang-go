#include <stdlib.h>
#include "bridge.h"

void * nixgo_store_id_to_userdata(uintptr_t id)
{
    uintptr_t *token = malloc(sizeof(uintptr_t));
    if (token != NULL)
        *token = id;
    return token;
}

uintptr_t nixgo_store_userdata_to_id(void * userdata)
{
    return userdata == NULL ? 0 : *(uintptr_t *) userdata;
}

void nixgo_store_userdata_free(void * userdata)
{
    free(userdata);
}
