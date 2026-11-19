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

NixCExternalValueDesc * nixgo_external_desc_new(void)
{
    NixCExternalValueDesc * desc = calloc(1, sizeof(NixCExternalValueDesc));
    if (desc == NULL)
        return NULL;
    desc->print = cgoExternalPrint;
    desc->showType = cgoExternalShowType;
    desc->typeOf = cgoExternalTypeOf;
    desc->coerceToString = cgoExternalCoerceToString;
    desc->equal = cgoExternalEqual;
    desc->printValueAsJSON = cgoExternalPrintValueAsJSON;
    desc->printValueAsXML = cgoExternalPrintValueAsXML;
    return desc;
}

void nixgo_external_desc_free(NixCExternalValueDesc * desc)
{
    free(desc);
}

void * nixgo_external_handle_to_ptr(uintptr_t handle)
{
    return (void *) handle;
}

uintptr_t nixgo_external_handle_from_ptr(void * ptr)
{
    return (uintptr_t) ptr;
}

void nixgo_external_register_finalizer(void * object, uintptr_t handle)
{
    nix_gc_register_finalizer(object, (void *) handle, cgoExternalFinalizer);
}
