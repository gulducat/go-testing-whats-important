package main

// Testable Examples are great because they are a clear signal
// directly to the user about how your Thing should be used.

func ExampleThing() {
	t := NewThing("hello")
	err := t.WriteWriter()
	if err != nil {
		panic(err)
	}
	// Output: hello
}
