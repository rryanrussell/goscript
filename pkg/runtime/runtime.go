package runtime

import _ "embed"

//go:embed runtime.js
var RuntimeJS []byte

//go:embed strings.js
var StringsJS []byte

//go:embed reflect.js
var ReflectJS []byte

//go:embed fmt.js
var FmtJS []byte

//go:embed errors.js
var ErrorsJS []byte

//go:embed strconv.js
var StrconvJS []byte

//go:embed math.js
var MathJS []byte

//go:embed time.js
var TimeJS []byte

//go:embed sort.js
var SortJS []byte

//go:embed sync.js
var SyncJS []byte

//go:embed chan.js
var ChanJS []byte
