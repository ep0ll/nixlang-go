module github.com/ep0ll/nixlang-go

go 1.22

// Requires Nix C API from master (libutil-c, libstore-c, libexpr-c, libfetchers-c, libflake-c, libmain-c).
// Build with: pkg-config available for nix-*-c packages, e.g. via nix develop / flake.
