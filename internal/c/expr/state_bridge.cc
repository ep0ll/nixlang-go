#include "nix_api_expr_internal.h"
#include "bridge.h"

extern "C" EvalState * nixgo_expr_clone_callback_state(EvalState * state)
{
    if (state == nullptr)
        return nullptr;
    try {
        return new EvalState{state->state};
    } catch (...) {
        return nullptr;
    }
}

extern "C" void nixgo_expr_free_callback_state(EvalState * state)
{
    delete state;
}
