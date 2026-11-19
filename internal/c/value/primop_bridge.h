#ifndef NIXLANG_GO_PRIMOP_BRIDGE_H
#define NIXLANG_GO_PRIMOP_BRIDGE_H

#include <stdint.h>

void * nixgo_primop_handle_to_ptr(uintptr_t handle);
uintptr_t nixgo_primop_handle_from_ptr(void * ptr);

#endif
