#ifndef NIXLANG_GO_EXTERNAL_BRIDGE_H
#define NIXLANG_GO_EXTERNAL_BRIDGE_H

#include <stdint.h>
#include <nix_api_external.h>

#define NIXGO_EXTERNAL_PRINT       (1u << 0)
#define NIXGO_EXTERNAL_SHOW_TYPE   (1u << 1)
#define NIXGO_EXTERNAL_TYPE_OF     (1u << 2)
#define NIXGO_EXTERNAL_COERCE      (1u << 3)
#define NIXGO_EXTERNAL_EQUAL       (1u << 4)
#define NIXGO_EXTERNAL_JSON        (1u << 5)
#define NIXGO_EXTERNAL_XML         (1u << 6)

NixCExternalValueDesc * nixgo_external_desc_new(uint32_t flags);
void nixgo_external_desc_free(NixCExternalValueDesc * desc);
void * nixgo_external_handle_to_ptr(uintptr_t handle);
uintptr_t nixgo_external_handle_from_ptr(void * ptr);
void nixgo_external_register_finalizer(void * object, uintptr_t handle);

#endif
