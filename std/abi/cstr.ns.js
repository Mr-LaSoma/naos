package abi;

using @import("std.types.byte");
using @import("std.types.size");

pub type cstr :: [*]byte;

cstr # {
    pub fromBytes(source[*]byte, len: usize) -> cstr {
        ptr: [*]byte = @alloc(len + 1);
        @memcopy(ptr, source, len);
        ptr[len] = '\0' as byte;

        return ptr @bitcast(cstr);
    }

    @overload @as(source: [*]byte) -> cstr @invalid
    @overload @as(source: cstr) -> [*]byte {
        return source @bitcast([*]byte);
    }

    @overload @delete (self: * cstr) -> void {
        @free(cstr)
}
}
