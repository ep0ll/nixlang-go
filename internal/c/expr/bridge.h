#ifndef NIXLANG_GO_EXPR_BRIDGE_H
#define NIXLANG_GO_EXPR_BRIDGE_H

#include <stdint.h>
#include <nix_api_expr.h>

void * nixgo_expr_id_to_userdata(uintptr_t id);
void nixgo_expr_userdata_free(void * userdata);
EvalState * nixgo_expr_clone_callback_state(EvalState * state);
void nixgo_expr_free_callback_state(EvalState * state);
void * nixgo_expr_state_identity(EvalState * state);

#endif
