#include <stdlib.h>
#include <nix_api_external.h>
#include <nix_api_expr.h>
#include "bridge.h"

extern void cgoExternalPrint(void *, nix_printer *);
extern void cgoExternalShowType(void *, nix_string_return *);
extern void cgoExternalTypeOf(void *, nix_string_return *);
extern void cgoExternalCoerceToString(void *, nix_string_context *, int, int, nix_string_return *);
extern int cgoExternalEqual(void *, void *);
extern void cgoExternalPrintValueAsJSON(void *, EvalState *, bool, nix_string_context *, bool, nix_string_return *);
extern void cgoExternalPrintValueAsXML(void *, EvalState *, int, int, void *, nix_string_context *, void *, int);
extern void cgoExternalFinalizer(void *, void *);

NixCExternalValueDesc * nixgo_external_desc_new(uint32_t flags)
{
    NixCExternalValueDesc * desc = calloc(1, sizeof(NixCExternalValueDesc));
    if (desc == NULL)
        return NULL;
    if (flags & NIXGO_EXTERNAL_PRINT) desc->print = cgoExternalPrint;
    if (flags & NIXGO_EXTERNAL_SHOW_TYPE) desc->showType = cgoExternalShowType;
    if (flags & NIXGO_EXTERNAL_TYPE_OF) desc->typeOf = cgoExternalTypeOf;
    if (flags & NIXGO_EXTERNAL_COERCE) desc->coerceToString = cgoExternalCoerceToString;
    if (flags & NIXGO_EXTERNAL_EQUAL) desc->equal = cgoExternalEqual;
    if (flags & NIXGO_EXTERNAL_JSON) desc->printValueAsJSON = cgoExternalPrintValueAsJSON;
    if (flags & NIXGO_EXTERNAL_XML) desc->printValueAsXML = cgoExternalPrintValueAsXML;
    return desc;
}

void nixgo_external_desc_free(NixCExternalValueDesc * desc)
{
    free(desc);
}

void * nixgo_external_handle_to_ptr(uintptr_t handle)
{
    uintptr_t * token = malloc(sizeof(uintptr_t));
    if (token != NULL)
        *token = handle;
    return token;
}

uintptr_t nixgo_external_handle_from_ptr(void * ptr)
{
    return ptr == NULL ? 0 : *(uintptr_t *) ptr;
}

void nixgo_external_handle_free(void * ptr)
{
    free(ptr);
}

void nixgo_external_register_finalizer(void * object, void * userdata)
{
    nix_gc_register_finalizer(object, userdata, cgoExternalFinalizer);
}
