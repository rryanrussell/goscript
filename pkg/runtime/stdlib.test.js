import { test, expect } from "bun:test";
import { readFileSync } from "fs";
import { join } from "path";

// Setup global runtime
global.runtime = {};

// Load runtime files
// Note: These are raw scripts intended for concatenation.
const load = (file) => {
    const content = readFileSync(join(import.meta.dir, file), "utf8");
    eval(content);
};

load("fmt.js");
load("strings.js");
load("errors.js");
load("strconv.js");
load("math.js");
load("time.js");
load("sort.js");

test("fmt.sprintf", () => {
    expect(runtime.fmt.sprintf("Hello %s", "World")).toBe("Hello World");
    expect(runtime.fmt.sprintf("Value: %d", 42)).toBe("Value: 42");
    expect(runtime.fmt.sprintf("Float: %f", 3.14)).toBe("Float: 3.14");
    expect(runtime.fmt.sprintf("Bool: %t", true)).toBe("Bool: true");
    expect(runtime.fmt.sprintf("Type: %T", 123)).toBe("Type: number");
    expect(runtime.fmt.sprintf("Type: %T", { __gs_goType: "main.Foo" })).toBe("Type: main.Foo");
    expect(runtime.fmt.sprintf("Percent: %%")).toBe("Percent: %");
});

test("strings.Contains", () => {
    expect(runtime.strings.Contains("seafood", "foo")).toBe(true);
    expect(runtime.strings.Contains("seafood", "bar")).toBe(false);
});

test("strings.HasPrefix", () => {
    expect(runtime.strings.HasPrefix("Gopher", "Go")).toBe(true);
    expect(runtime.strings.HasPrefix("Gopher", "C")).toBe(false);
});

test("strings.HasSuffix", () => {
    expect(runtime.strings.HasSuffix("Amigo", "go")).toBe(true);
    expect(runtime.strings.HasSuffix("Amigo", "Ami")).toBe(false);
});

test("strings.Index", () => {
    expect(runtime.strings.Index("chicken", "ken")).toBe(4);
    expect(runtime.strings.Index("chicken", "dmr")).toBe(-1);
});

test("strings.Split", () => {
    expect(runtime.strings.Split("a,b,c", ",")).toEqual(["a", "b", "c"]);
});

test("strings.Join", () => {
    expect(runtime.strings.Join(["a", "b", "c"], ",")).toBe("a,b,c");
});

test("strings.Replace", () => {
    expect(runtime.strings.Replace("oink oink oink", "k", "ky", 2)).toBe("oinky oinky oink");
    expect(runtime.strings.Replace("oink oink oink", "oink", "moo", -1)).toBe("moo moo moo");
});

test("strings.TrimSpace", () => {
    expect(runtime.strings.TrimSpace(" \t\n Hello, Gophers \n\t\r\n")).toBe("Hello, Gophers");
});

test("errors", () => {
    const err = runtime.errors.New("failed");
    expect(err).toBeInstanceOf(Error);
    expect(err.message).toBe("failed");
    expect(runtime.errors.Is(err, err)).toBe(true);
    expect(runtime.errors.Is(err, new Error("failed"))).toBe(true);
    expect(runtime.errors.Is(err, new Error("other"))).toBe(false);
});

test("strconv", () => {
    expect(runtime.strconv.Atoi("123")).toEqual([123, null]);
    expect(runtime.strconv.Atoi("abc")[1]).toBeInstanceOf(Error);
    expect(runtime.strconv.Itoa(123)).toBe("123");
    expect(runtime.strconv.ParseFloat("3.14")).toEqual([3.14, null]);
    expect(runtime.strconv.ParseBool("true")).toEqual([true, null]);
    expect(runtime.strconv.FormatBool(false)).toBe("false");
});

test("math", () => {
    expect(runtime.math.Abs(-5)).toBe(5);
    expect(runtime.math.Floor(3.9)).toBe(3);
    expect(runtime.math.Ceil(3.1)).toBe(4);
    expect(runtime.math.Max(1, 2)).toBe(2);
    expect(runtime.math.Min(1, 2)).toBe(1);
    expect(runtime.math.Sqrt(16)).toBe(4);
    expect(runtime.math.IsNaN(NaN)).toBe(true);
});

test("time", () => {
    const now = runtime.time.Now();
    expect(now).toBeInstanceOf(Date);
    expect(runtime.time.Second).toBe(1000000000);
    // Since is hard to test precisely, but we can check it returns a number
    expect(typeof runtime.time.Since(now)).toBe("number");
});

test("sort", () => {
    const ints = [3, 1, 2];
    runtime.sort.Ints(ints);
    expect(ints).toEqual([1, 2, 3]);

    const strings = ["c", "a", "b"];
    runtime.sort.Strings(strings);
    expect(strings).toEqual(["a", "b", "c"]);

    const slice = [3, 1, 2];
    runtime.sort.Slice(slice, (i, j) => slice[i] < slice[j]);
    expect(slice).toEqual([1, 2, 3]);
});
