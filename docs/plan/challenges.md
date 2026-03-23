# Plan: The Great Go Massive Challenge Library

## 1. Objective
Construct an exhaustive library of **105 "fixme" challenges** (5 scenarios across 21 topics). These are "Detective Scenarios" designed to be bite-sized, high-stakes, and educational, forcing participants to break their habits from dynamic languages and adopt the "Gopher mindset."

## 2. Directory Structure
Challenges will be grouped by fine-grained topics in the `internal/challenges/basics/fixme/` directory:
```text
internal/challenges/basics/fixme/
├── 01_slices_basics_test.go
├── 02_slices_memory_test.go
├── 03_arrays_test.go
├── 04_maps_init_test.go
├── 05_maps_usage_test.go
├── 06_pointers_basics_test.go
├── 07_pointers_traps_test.go
├── 08_structs_test.go
├── 09_embedding_test.go
├── 10_receivers_test.go
├── 11_interfaces_nil_test.go
├── 12_interfaces_assertions_test.go
├── 13_errors_basics_test.go
├── 14_errors_advanced_test.go
├── 15_goroutines_test.go
├── 16_channels_basics_test.go
├── 17_channels_select_test.go
├── 18_sync_mutex_test.go
├── 19_sync_waitgroup_test.go
├── 20_context_test.go
└── 21_testing_bench_test.go
```

---

## 3. The Library (105 Scenarios)

### Topic 1: Slices (Basics & Headers)
1.  **"The Phantom Append"**: Function appends to a slice, but the caller doesn't see the original length. (Lesson: Slice headers are passed by value).
2.  **"The Slice Reset"**: Setting a slice to `nil` inside a function doesn't clear it for the caller.
3.  **"The Partial Update"**: Modifying `s[0]` works, but `append(s, "new")` doesn't reflect in the caller.
4.  **"The Reassignment Ghost"**: Reassigning `s = s[1:]` inside a function leaves the caller with the original first element.
5.  **"The In-place Sort Surprise"**: Participants are shocked to find that passing a slice to a sort function modifies their original data.

### Topic 2: Slices (Memory & Capacity)
1.  **"The Zero-Length Trap"**: `make([]int, 5)` then `append` leads to `[0 0 0 0 0 1]`. (Lesson: Length vs Capacity).
2.  **"The Hidden Giant"**: Slicing a tiny window `massive[0:1]` keeps a 1GB underlying array in memory forever.
3.  **"The Allocation Breakup"**: Appending to a slice near capacity suddenly stops reflecting in other slices sharing the same array.
4.  **"The Copy Cat"**: Using `copy(dst, src)` where `dst` is uninitialized, resulting in an empty destination.
5.  **"The Capacity Explosion"**: A loop that appends millions of items without pre-allocation, causing a performance death spiral.

### Topic 3: Arrays (Fixed Constraints)
1.  **"The 8KB Copy"**: Passing a `[1000]int` by value and wondering why the benchmark is so slow.
2.  **"The Size Wall"**: Trying to pass a `[5]int` into a function expecting `[10]int`. (Lesson: Size is part of the type).
3.  **"The Array Pointer Mutation"**: Using `*[5]int` to finally successfully mutate an array from another function.
4.  **"The Equality Paradox"**: Arrays can be compared with `==`, but slices cannot. Why?
5.  **"The Range Value Snapshot"**: Modifying the loop variable `v` in `for _, v := range array` doesn't change the array.

### Topic 4: Maps (Lifecycle & Initialization)
1.  **"The Dead Map"**: Writing to `var m map[string]int` (nil) causes an immediate panic.
2.  **"The Read Silence"**: Reading from a nil map returns `0` instead of panicking, leading to silent logical errors.
3.  **"The Literal Shortcut"**: Identifying when to use `map[string]int{}` vs `make()`.
4.  **"The Clear Conundrum"**: Trying to "clear" a map by setting it to `nil` (incorrect) vs a new `make`.
5.  **"The Pointer Map Uselessness"**: Debugging code that uses `*map[string]int`—a redundant and confusing pattern.

### Topic 5: Maps (Read/Write Semantics)
1.  **"The Zero vs Missing"**: A test that fails because it can't tell if a score is `0` or if the player isn't in the map. (Lesson: `comma-ok`).
2.  **"The Unaddressable Struct"**: `m["key"].Field = 1` fails to compile. Why can't we mutate struct fields in maps directly?
3.  **"The Mystery Order"**: A test that depends on map iteration order and fails 50% of the time.
4.  **"The Concurrent Crash"**: A visceral "fatal error: concurrent map writes" crash during a web request simulation.
5.  **"The Slice-as-Key Crime"**: Trying to use a slice as a map key and getting a compiler slap on the wrist.

### Topic 6: Pointers (Addressing & Mutation)
1.  **"The Literal Address"**: Trying to do `&42` or `&"hello"` and failing to understand why only variables have addresses.
2.  **"The Nil Pointer Roulette"**: A function that returns a pointer—dereferencing it without a nil check.
3.  **"The Swap Failure"**: The classic `Swap(a, b)` that doesn't swap because it didn't use pointers.
4.  **"The Double Pointer Maze"**: Navigating an API response that uses `**User` for optional updates.
5.  **"The Value Isolation"**: Modifying a struct passed by value and seeing the change "vanish" outside the function.

### Topic 7: Pointers (The Loop Variable Trap)
1.  **"The Last Element Echo"**: Storing `&v` from a `range` loop into a slice—every entry is the last item!
2.  **"The Goroutine Closure Capture"**: A goroutine inside a loop captures `i`, printing `10` ten times.
3.  **"The Recycled Pointer"**: Using a pointer receiver method inside a loop on a value-type range variable.
4.  **"The Slice of Pointers Update"**: Updating a `[]*User` correctly vs incorrectly updating the local variable.
5.  **"The Escaping Local"**: Returning `&x` from a function—proving Go isn't C++ and the stack isn't the limit.

### Topic 8: Structs (Encapsulation & Factory Functions)
1.  **"The Secret Field"**: Trying to access `u.password` from another package. (Lesson: Exporting rules).
2.  **"The Shallow Copy Trap"**: Copying a struct that contains a slice—changing the slice in one copies it to the other.
3.  **"The JSON Invisibility"**: A field `ID` not appearing in JSON because the tag is `id` but the field is unexported `id`.
4.  **"The Anonymous Comparison"**: Comparing two `struct{Name string}`—when are they "identical"?
5.  **"The Factory Fail"**: A `NewUser()` function that returns a value but should return a pointer for consistency.

### Topic 9: Embedding (Composition & Shadowing)
1.  **"The Shadowed ID"**: Both `Base` and `User` have `ID`. `u.ID` gives the wrong one!
2.  **"The Ambiguous Method"**: `Admin` embeds `Manager` and `Employee`, both have `GetName()`. Calling it crashes the compiler.
3.  **"The promotion shortcut"**: Calling an inner method directly on the outer instance and wondering how it works.
4.  **"The Hidden Interface"**: An outer struct satisfies an interface because its *embedded* inner struct does.
5.  **"The Embedding != Inheritance"**: Trying to pass `Admin` into a function that takes `User` (it fails!).

### Topic 10: Receivers (Methods & Mutability)
1.  **"The Snapshot Method"**: `(u User) SetEmail(...)` doesn't change the email. (Lesson: Value vs Pointer receivers).
2.  **"The Mutex Copy Disaster"**: A value receiver copies a struct with a `sync.Mutex`, creating two independent, broken locks.
3.  **"The Nil Receiver Logic"**: Calling a method on a nil struct pointer and having it actually work (but it shouldn't).
4.  **"The Implicit Pointer"**: Why `v.Method()` works even if `Method` takes a pointer, but only if `v` is addressable.
5.  **"The Receiver Consistency"**: Mixing value and pointer receivers on the same type and seeing the interface implementation break.

### Topic 11: Interfaces (The Nil-Not-Nil Mystery)
1.  **"The Interface Lie"**: `var err *MyErr = nil; var i error = err; i == nil` is FALSE. The most famous Go gotcha.
2.  **"The Empty Interface Guess"**: Storing an `int` in `any` and trying to assert it as an `int64`.
3.  **"The Name Doesn't Matter"**: A struct satisfies an interface it doesn't even know exists.
4.  **"The Pointer Requirement"**: An interface requires a method `M()`, implemented by `*T`. Passing `T` fails.
5.  **"The Internal Leak"**: An interface that exposes methods returning unexported types.

### Topic 12: Interfaces (Type Assertions & Switches)
1.  **"The Panic Assertion"**: `i.(string)` when `i` is an `int`. Crash!
2.  **"The Comma-OK Safety"**: Correctly using `s, ok := i.(string)` to prevent the crash.
3.  **"The Type Switch Fallthrough"**: Discovering that `type` switches are special and don't allow fallthrough.
4.  **"The Interface Wrapping Chain"**: Using `errors.As` to find a specific error type 3 layers deep.
5.  **"The Pointer to Interface Anti-pattern"**: Debugging why `func(err *error)` is almost always a logical error.

### Topic 13: Error Handling (Wrapping & Scoping)
1.  **"The Ignored Error"**: A function returns `(Result, error)` and the participant uses `res, _ := ...`.
2.  **"The Lost Cause"**: Using `fmt.Errorf("failed: %v", err)`—now `errors.Is` can't find the root cause.
3.  **"The Shadowed Err"**: `if res, err := ...` inside a loop hiding the outer `err` that needs to be returned.
4.  **"The Multi-Failure"**: Using `errors.Join` to return 3 separate validation errors at once.
5.  **"The Fragile String"**: `if err.Error() == "not found"`—explaining why string comparison is the path to pain.

### Topic 14: Error Handling (Custom Types & Sentinel Errors)
1.  **"The Sentinel Check"**: Refactoring string checks to use `errors.Is(err, ErrNotFound)`.
2.  **"The Custom Field Access"**: Using `errors.As` to extract a `RetryAfter` field from a custom error.
3.  **"The Error Implementation"**: Forgetting the `*` in `func (e *MyErr) Error()` and breaking the interface.
4.  **"The Typed Nil Return"**: Returning a nil pointer of type `*MyErr` from a function returning `error`.
5.  **"The Error Wrapping Loop"**: An error that wraps itself, causing a stack overflow on `Error()`.

### Topic 15: Goroutines (Closures & Lifetime)
1.  **"The Vanishing Main"**: Main exits before the goroutine even starts its work.
2.  **"The Index Capture"**: 10 goroutines in a loop all thinking they are index `10`.
3.  **"The Shared Counter Race"**: `count++` from two goroutines resulting in `1` instead of `2`.
4.  **"The Leaking Worker"**: A goroutine started in a loop that never finishes and eats memory.
5.  **"The Panic Exit"**: A goroutine panics and kills the entire server.

### Topic 16: Channels (Buffering & Deadlocks)
1.  **"The Eternal Wait"**: Sending to an unbuffered channel with no receiver.
2.  **"The 11th Item"**: Sending to a buffered channel of size 10 that is already full.
3.  **"The Nil Channel Hang"**: Sending to a channel that was declared but never `make()`'d.
4.  **"The Closed Read Zero"**: Reading from a closed channel and getting `0` forever, leading to an infinite loop.
5.  **"The Unclosed Range"**: `for v := range ch` hanging because the sender forgot to `close(ch)`.

### Topic 17: Channels (Closing & Select)
1.  **"The Dead Write"**: Sending a "Goodbye" message to a channel that was already closed. Crash!
2.  **"The Panic Close"**: Closing a channel twice.
3.  **"The Select Default Spin"**: A `select` with `default` inside a loop that hits 100% CPU.
4.  **"The Random Choice"**: A test that fails because `select` picked `ch1` instead of `ch2` randomly.
5.  **"The Timeout Pattern"**: Fixing a hanging request using `select` and `time.After`.

### Topic 18: Sync Primitives (Mutex & RWMutex)
1.  **"The Return Lock"**: Returning early from a function but forgetting to call `mu.Unlock()`. (Lesson: `defer`).
2.  **"The Re-entrant Deadlock"**: A method `A` locks, calls `B`, which also tries to lock.
3.  **"The Mutex Copy"**: Passing a struct containing a mutex by value to another function.
4.  **"The RWMutex Writer Starvation"**: 1000 readers keeping a writer locked out forever.
5.  **"The Unlocked Unlock"**: Calling `Unlock` on a mutex that isn't locked.

### Topic 19: Sync Primitives (WaitGroup & Once)
1.  **"The Done Overdose"**: Calling `wg.Done()` 11 times for 10 goroutines.
2.  **"The Late Add"**: Calling `wg.Add(1)` inside the `go func()`, causing `wg.Wait()` to return immediately.
3.  **"The Shared WaitGroup"**: Passing `wg` by value instead of pointer to workers.
4.  **"The Recursive Once"**: `once.Do` calling another function that also calls `once.Do`.
5.  **"The Forgotten Wait"**: Starting 100 goroutines but forgetting to call `wg.Wait()` at the end.

### Topic 20: Context (Propagation & Leaks)
1.  **"The Context Root"**: Using `context.Background()` inside a database call instead of the request context.
2.  **"The Timer Leak"**: `WithTimeout` without `defer cancel()` leading to memory growth.
3.  **"The Ignored Cancel"**: A 30-second processing loop that never checks `ctx.Done()`.
4.  **"The Value Misuse"**: Putting a `db.Conn` or `logger` inside context values (Antipattern).
5.  **"The Context vs Channel"**: Explaining why you shouldn't send a context through a channel to "cancel" work.

### Topic 21: Testing & Benchmarking (Helpers & Timers)
1.  **"The Wrong Line"**: A test helper fails, but the error points to the helper code, not the test. (Lesson: `t.Helper()`).
2.  **"The Setup Penalty"**: A benchmark that is 100x slower because it includes 5 seconds of setup time. (Lesson: `b.ResetTimer()`).
3.  **"The Parallel Race"**: `t.Parallel()` tests that accidentally share a global variable.
4.  **"The Compiler Deletion"**: A benchmark that runs in 0.00ns because the compiler optimized the code away.
5.  **"The Assert Crash"**: Using `assert.Nil(t, err)` instead of `require.Nil(t, err)`, causing a nil-pointer panic on the next line.

## 4. Engagement & Educational Value
Every challenge is a **Detective Brief**. The goal is to solve the "Mystery of the Nil Map" or "The Case of the Missing Item." 
By solving these, participants don't just fix code—they build an intuition for Go's memory layout and runtime behavior.

## 5. Next Steps
Once this plan is approved, I will begin implementing the challenges in batches of 3 topics.
