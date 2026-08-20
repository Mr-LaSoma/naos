package utils;

@import("std.abi.cstr");

pub type Option<T> :: enum { // love you rust btw
    None,
    Some(T),
}

Option < T > # {
    pub isSome() -> bool { return @isvariant(self, Option.Some); }
    pub isNone() -> bool { return @isvariant(self, Option.None); }

    pub unwrap() -> T {
        if self.isNone() { @panic("unwrapped none option"); }
        return @valueof(self);
    }

    pub expect(errorMsg: cstr) -> T {
        if self.isNone() { @panic(errorMsg); }
        return @valueof(self);
    }

    pub unwrapOr(defaultValue: T) -> T {
        if self.isNone() { return defaultValue; }
        return @valueof(self);
    }
}
