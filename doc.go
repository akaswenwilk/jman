// Package jman provides fast, minimal helpers for building, querying and asserting JSON
// objects (Obj) and arrays (Arr) in tests.
//
// # Core Types
//
//   • Obj — JSON object implemented as `map[string]any`.
//   • Arr — JSON array implemented as `[]any`.
//
// Both types satisfy the JSONEqual interface and can be created:
//
//   // Literal
//   user := jman.Obj{"id": 1, "name": "alice"}
//   tags := jman.Arr{"go", "test", 42}
//
//   // Parse / normalise
//   user, _ := jman.New[jman.Obj](`{"id":1,"name":"alice"}`)
//   tags, _ := jman.New[jman.Arr]([]byte(`["go","test",42]`))
//
// # Path-based Getters & Setters
//
// Paths use JSON-Path–like dot syntax and must start with `$`.
//
//   id, _ := user.Get("$.id").Number()
//   _     = user.Set("$.settings.theme", "dark")
//   _     = tags.Set("$.1", "unit")
//
// Getter returns a Result that can be converted with:
//  
// Setter can set any type as a value, however the value will be normalized into either:
// bool, string, float64, Obj, or Arr.
//
//   String(), Number(), Bool(), Obj(), Arr()
//
// # Deep Equality
//
//   expected := jman.Obj{
//       "id":         "{{uuid}}",
//       "name":       "{{nonEmpty}}",
//       "roles":      jman.Arr{"admin", "editor"},
//       "addresses":  jman.Arr{
//           jman.Obj{"street": "High", "no": 1},
//           jman.Obj{"street": "Low",  "no": 9},
//       },
//   }
//
//   actual := `{
//       "id":"9b74c989-7cdf-41fa-9a49-5290f31e59d3",
//       "name":"alice",
//       "roles":["editor","admin"],
//       "addresses":[
//         {"street":"High","no":1},
//         {"street":"Low","no":9}
//       ]}`
//
//   err := expected.Equal(actual,
//       jman.WithIgnoreArrayOrder("$.roles", "$.addresses"),
//       jman.WithDefaultMatchers(jman.Matchers{
//           jman.IsUUID("{{uuid}}"),
//           jman.NotEmpty("{{nonEmpty}}"),
//       }),
//   )
//   if err != nil {
//       fmt.Println(err) // human-readable diff
//   }
//
// # Inequality Report
//
// For each difference the path and problem is returned, e.g.:
//
//   expected not equal to actual:
//   $.roles expected 2 items - got 3 items
//   $.name expected "alice" - got "bob"
//   $.extra unexpected key
//
// # Matchers
//
// Matchers allow placeholders in the *expected* JSON that are resolved
// at comparison time:
//
//   jman.IsUUID("{{uuid}}")          // any valid UUID
//   jman.NotEmpty("{{nonEmpty}}")    // non-empty string/array/object
//   jman.EqualMatcher("{{id}}", 99)  // equals specific value
//
// Write your own with `jman.Custom`. A placeholder is a string that when found in the expected as a value,
// will find the corresponding value in the actual JSON and compare it using the matcher.
//
// # Options
//
//   • WithIgnoreArrayOrder(paths...) — compare arrays as sets for given paths.
//   • WithDefaultMatchers(ms)       — register Matchers once per comparison.
//
package jman
