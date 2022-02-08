# alang

An interepreted (maybe?) turing complete language in Go.

## Examples

`go run main.go`

```
>> let add = fn(x,y) { return x+y; };
>> add(5, 10); // 15
...
>> 10 / 2 + 3; // 8
...
>> let SOME_VAR = 5
>> fn(x) { SOME_VAR * x }(10) // 50
```

## Run tests

```
go test ./...
```

## Next steps 

- [ ] Support to complex structures, arrays, slices, hashmaps
- [ ] Be compiled

