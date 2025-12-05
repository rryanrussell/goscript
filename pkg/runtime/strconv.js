if (typeof runtime === 'undefined') {
    runtime = {}
}

runtime.strconv = {
    Atoi: (s) => {
        const n = parseInt(s, 10);
        if (isNaN(n)) {
            return [0, new Error("strconv.Atoi: parsing " + JSON.stringify(s) + ": invalid syntax")];
        }
        return [n, null];
    },
    Itoa: (n) => String(n),
    ParseFloat: (s, bitSize) => {
        const f = parseFloat(s);
        if (isNaN(f)) {
            return [0, new Error("strconv.ParseFloat: parsing " + JSON.stringify(s) + ": invalid syntax")];
        }
        return [f, null];
    },
    FormatInt: (i, base) => {
        return i.toString(base);
    },
    FormatFloat: (f, fmt, prec, bitSize) => {
        // Simplified implementation
        return f.toString();
    },
    ParseBool: (str) => {
        if (str === "true") return [true, null];
        if (str === "false") return [false, null];
        return [false, new Error("strconv.ParseBool: parsing " + JSON.stringify(str) + ": invalid syntax")];
    },
    FormatBool: (b) => {
        return b ? "true" : "false";
    }
};
