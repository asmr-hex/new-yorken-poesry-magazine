# XAQT Library
*TODO (cw|4.28.2018) create godocs instead of this file...these are mostly notes for myself during development.*

# Public Types
* **`xaqt.Context`**: entrypoint into all functionality. TODO Propose renaming?
* **`xaqt.Compilers`**: list of compilters. Propose re-typing to `[]Compiler`.
* **`xaqt.Message`**: details on success or failure of execution

# Public Methods

``` go
func NewContext(xaqt.Compilers, ...option) // option should be public!
func (ctx *Context) ReadCompilers(string)
func (ctw *Context) Languages([]string) // rename -> GetSupportedLanguages ?
func (ctw *Context) Evaluate(string, string, []string)
```
