#include "bridge.h"
void * nixgo_store_id_to_userdata(uintptr_t id) { return (void *) id; }
uintptr_t nixgo_store_userdata_to_id(void * userdata) { return (uintptr_t) userdata; }
