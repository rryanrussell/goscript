if (typeof runtime === 'undefined') {
    runtime = {}
}

runtime.fmt = {
    sprintf(format, ...args) {
        let argIndex = 0;
        let result = "";

        for (let i = 0; i < format.length; i++) {
            if (format[i] === '%' && i + 1 < format.length) {
                const verb = format[i + 1];
                let val = args[argIndex++];

                switch (verb) {
                    case 'v':
                    case 's':
                        if (val === undefined || val === null) {
                            result += "<nil>";
                        } else if (typeof val === 'object' && val.__gs_goType) {
                            // TODO: Call String() method if exists
                            result += JSON.stringify(val); // Simple fallback
                        } else {
                            result += String(val);
                        }
                        break;
                    case 'd':
                        result += parseInt(val);
                        break;
                    case 'f':
                        result += parseFloat(val);
                        break;
                    case 't':
                        result += String(!!val);
                        break;
                    case 'T':
                        if (val === undefined || val === null) {
                            result += "nil";
                        } else if (val.__gs_goType) {
                            result += val.__gs_goType;
                        } else {
                            result += typeof val;
                        }
                        break;
                    case '%':
                        result += '%';
                        argIndex--; // Consume nothing
                        break;
                    default:
                        result += '%' + verb;
                        argIndex--;
                }
                i++; // Skip verb
            } else {
                result += format[i];
            }
        }
        return result;
    },

    printf(format, ...args) {
        const s = this.sprintf(format, ...args);
        console.log(s);
    },

    println(...args) {
        const s = args.map(arg => {
            if (arg === undefined || arg === null) {
                return "<nil>";
            }
            if (typeof arg === 'object' && arg.__gs_goType) {
                return JSON.stringify(arg);
            }
            return String(arg);
        }).join(" ");
        console.log(s);
    }
};
