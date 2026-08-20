package utils;

pub type Result<T, E> :: enum {
    Ok(T),
        Err(E),
}

Result < T, E > # {
    pub isOk() -> bool { return @isvariant(self, Result.Ok); }
    pub isErr() -> bool { return @isvariant(self, Result.None); }

    pub unwrap() -> T {
        if self.isNone() { @panic("unwrapped err result"); }
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
