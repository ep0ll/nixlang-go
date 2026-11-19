#include <stdlib.h>
#include "primop_bridge.h"

void * nixgo_primop_handle_to_ptr(uintptr_t handle)
{
    uintptr_t *token = malloc(sizeof(uintptr_t));
    if (token != NULL)
        *token = handle;
    return token;
}

uintptr_t nixgo_primop_handle_from_ptr(void * ptr)
{
    return ptr == NULL ? 0 : *(uintptr_t *) ptr;
}

void nixgo_primop_handle_free(void * ptr)
{
    free(ptr);
}
