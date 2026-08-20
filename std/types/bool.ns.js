package types;

pub type bool :: u1;
pub const true: bool = 1;
pub const false: bool = 0;

bool # {
    @overload @as(source: u1) -> bool {
        return if (source != 0) { => true; } else { => false; };
    }

    @overload @as(source: bool) -> u1 {
        return if (source) { => 1; } else { => 0 };
    }
}
