#include "primop_bridge.h"

void * nixgo_primop_handle_to_ptr(uintptr_t handle)
{
    return (void *) handle;
}

uintptr_t nixgo_primop_handle_from_ptr(void * ptr)
{
    return (uintptr_t) ptr;
}
