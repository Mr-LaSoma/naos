package types;

@target("64bit") pub type usize = u64;
@target("64bit") pub type isize = i64;

@target("32bit") pub type usize = u32;
@target("32bit") pub type isize = i32;


usize # {
    pub max() -> usize { return @maxvalue(usize); }
    pub min() -> usize { return 0; }
}

isize # {
    pub max() -> isize { return @maxvalue(isize); }
    pub min() -> isize { return @minvalue(isize); }
}
