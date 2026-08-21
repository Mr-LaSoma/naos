package utils;

using @import("std.abi.cstr");

pub type Option<T> :: enum {
    Some(T);
None;
};

Option < T > # {
    pub isSome(const self: *Option) -> bool { return @isvariant(self, Option.Some); }
    pub isNone(const self: *Option) -> bool { return @isvariant(self, Option.None); }

    pub unwrap(const self: *Option) -> T {
        if (self.isNone()) { @panic("unwrapped none option"); }
        return @valueof(self);
    }

    pub expect(const errorMsg: cstr) -> T {
        if (self.isNone()) { @panic(errorMsg); }
        return @valueof(self);
    }

    pub unwrapOr(const value: T) -> T {
        if (self.isNone()) { return value; }
        return @valueof(self);
    }
}
