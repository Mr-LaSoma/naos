# Compiler actions

## List

| Action | Description |
| ------ | ----------- |
| . | REFLECTION |
| . | . |
| @sizeof(type) | returns the size in bytes of the type |
| @alignof(type) | returns the alignment in the memory of the type |
| @typeof(variable) | returns the type of the variable |
| @valueof(enum) | returns the value of an enum (if done one no value enum this will crash the program) |
| @isvariant(enum, enum.variant) | returns whether a certain enum is a certain variant |
| . | . |
| . | MEMORY OPERATION |
| . | . |
| @alloc(byte_size) | equivalent of malloc in c |
| @free(pointer) | frees the memory of the pointer |
| @alloca(byte_size) | returns the pointer to memory allocated in the stack |
| @memcopy(dest, src, size) | equivalent of memcopy in c |
| @memmove(dest, src, size) | equivalent of memmove in c |
| @memset(dest, value_byte, size) | equivalent of memset in c |
| . | . |
| . | TYPES CONTROL |
| . | . |
| expr @as(type) | casts the expression in the type, if not possible the program won't compile |
| expr @bitcast(type) | casts the previous expression in the type, this is unsafe |
| @null | nullptr in c++ or NULL in c |
| . | . |
| . | LIFE CYCLE |
| . | . |
| @overload @func | overloads a compiler function for a certain type (only usable in # {}, linking block ) |
| @func @invalid | makes a compiler function of a type invalid (won't compile if is invalidated) |
| @delete, @move, @copy | compiler functions overloadable in the # {}, called on specific times |
| . | . |
| . | OPTIMIZATIONS |
| . | . |
| @syscall(number, arguments) | calls the #number syscall |
| @panic(message) | crashes the program and logs a message |
| @unreachable | indication for optimizations |
| @inline func , @noinline func | indication to inline or not inline a function (can't use on compiler functions) |
| . | . |
| . | MISC |
| . | . |
| @log(message) | logs the message as a compilation message |
