if (typeof runtime === 'undefined') {
    runtime = {}
}

class Type {
    constructor(name, kind) {
        this.name = name;
        this.kind = kind;
    }

    Name() {
        return this.name;
    }

    String() {
        return this.name;
    }
}

class Value {
    constructor(val, type) {
        this.val = val;
        this.type = type;
    }

    Type() {
        return this.type;
    }

    String() {
        return String(this.val);
    }

    Interface() {
        return this.val;
    }
}

const __types = new Map();

function registerType(name, methods) {
    __types.set(name, new Type(name, "struct"));
}

function TypeOf(i) {
    if (i === null || i === undefined) {
        return null; // reflect.TypeOf(nil) is nil
    }

    let name = "";
    if (typeof i === "number") name = "float64"; // Default to float64 for now
    else if (typeof i === "string") name = "string";
    else if (typeof i === "boolean") name = "bool";
    else if (i.__gs_goType) name = i.__gs_goType;
    else name = "unknown";

    if (!__types.has(name)) {
        __types.set(name, new Type(name, "unknown"));
    }
    return __types.get(name);
}

function ValueOf(i) {
    return new Value(i, TypeOf(i));
}

runtime.reflect = {
    TypeOf: TypeOf,
    ValueOf: ValueOf,
    registerType: registerType,
};
