#ifndef NIXLANG_GO_STORE_BRIDGE_H
#define NIXLANG_GO_STORE_BRIDGE_H
#include <stdint.h>
void * nixgo_store_id_to_userdata(uintptr_t id);
uintptr_t nixgo_store_userdata_to_id(void * userdata);
void nixgo_store_userdata_free(void * userdata);
#endif
