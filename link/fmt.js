
fmt = {
    Println(...args) {
        console.log(...args)
    },
    Errorf(format, ...args) {
        return args.reduce((p,c) => p.replace(/%[svd]/,c), format)
    }
}

function GetFloat64FromEvent(e) {
    return JSON.parse(e.data)
}

function panic(err) {
    document.getElementsByTagName('canvas')[0].style.display = 'none';
    document.writeln(err)
}

