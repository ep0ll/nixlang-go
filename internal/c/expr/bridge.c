#include "bridge.h"
void * nixgo_expr_id_to_userdata(uintptr_t id) { return (void *) id; }
uintptr_t nixgo_expr_userdata_to_id(void * userdata) { return (uintptr_t) userdata; }
