#ifndef NIXLANG_GO_EXTERNAL_BRIDGE_H
#define NIXLANG_GO_EXTERNAL_BRIDGE_H

#include <stdint.h>
#include <nix_api_external.h>

NixCExternalValueDesc * nixgo_external_desc_new(void);
void nixgo_external_desc_free(NixCExternalValueDesc * desc);
void * nixgo_external_handle_to_ptr(uintptr_t handle);
uintptr_t nixgo_external_handle_from_ptr(void * ptr);
void nixgo_external_register_finalizer(void * object, uintptr_t handle);

#endif
