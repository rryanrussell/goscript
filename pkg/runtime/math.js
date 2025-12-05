if (typeof runtime === 'undefined') {
    runtime = {}
}

runtime.math = {
    PI: Math.PI,
    E: Math.E,
    Abs: Math.abs,
    Floor: Math.floor,
    Ceil: Math.ceil,
    Max: Math.max,
    Min: Math.min,
    Sqrt: Math.sqrt,
    Pow: Math.pow,
    Sin: Math.sin,
    Cos: Math.cos,
    Tan: Math.tan,
    Log: Math.log,
    Log2: Math.log2,
    Log10: Math.log10,
    Exp: Math.exp,
    Round: Math.round,
    Trunc: Math.trunc,
    Inf: (sign) => (sign >= 0 ? Infinity : -Infinity),
    IsInf: (f, sign) => {
        if (sign === 0) return !isFinite(f);
        if (sign > 0) return f === Infinity;
        return f === -Infinity;
    },
    IsNaN: Number.isNaN,
};
