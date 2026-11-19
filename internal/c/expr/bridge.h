#ifndef NIXLANG_GO_EXPR_BRIDGE_H
#define NIXLANG_GO_EXPR_BRIDGE_H
#include <stdint.h>
void * nixgo_expr_id_to_userdata(uintptr_t id);
uintptr_t nixgo_expr_userdata_to_id(void * userdata);
#endif
