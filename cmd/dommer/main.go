package main

import (
	"errors"
	"flag"
	"os"

	"github.com/rryanrussell/goscript/pkg/dommer"
)

var (
	input    = flag.String("input", "lib/lib.dom.d.ts", "Input file path")
	output   = flag.String("output", "app/dom.extern.go", "Output file path")
	debugMax = flag.Int("debug-max", 1000, "Max debug depth/items")
)

func data() []byte {
	data := `
	//start
	
	/** This IndexedDB API interface provides a connection to a database; you can use an IDBDatabase object to open a transaction on your database then create, manipulate, and delete objects (data) in that database. The interface provides the only way to get and manage versions of the database. */
	interface IDBDatabase extends EventTarget {
		/** Returns the name of the database. */
		readonly name: string;
		/** Returns a list of the names of object stores in the database. */
		readonly objectStoreNames: DOMStringList;
		onabort: ((this: IDBDatabase, ev: Event) => any) | null;
		onclose: ((this: IDBDatabase, ev: Event) => any) | null;
		onerror: ((this: IDBDatabase, ev: Event) => any) | null;
		onversionchange: ((this: IDBDatabase, ev: IDBVersionChangeEvent) => any) | null;
		/** Returns the version of the database. */
		readonly version: number;
		/** Closes the connection once all running transactions have finished. */
		close(): void;
		/**
		 * Creates a new object store with the given name and options and returns a new IDBObjectStore.
		 *
		 * Throws a "InvalidStateError" DOMException if not called within an upgrade transaction.
		 */
		createObjectStore(name: string, options?: IDBObjectStoreParameters): IDBObjectStore;
		/**
		 * Deletes the object store with the given name.
		 *
		 * Throws a "InvalidStateError" DOMException if not called within an upgrade transaction.
		 */
		deleteObjectStore(name: string): void;
		/** Returns a new transaction with the given mode ("readonly" or "readwrite") and scope which can be a single object store name or an array of names. */
		transaction(storeNames: string | string[], mode?: IDBTransactionMode): IDBTransaction;
		addEventListener<K extends keyof IDBDatabaseEventMap>(type: K, listener: (this: IDBDatabase, ev: IDBDatabaseEventMap[K]) => any, options?: boolean | AddEventListenerOptions): void;
		addEventListener(type: string, listener: EventListenerOrEventListenerObject, options?: boolean | AddEventListenerOptions): void;
		removeEventListener<K extends keyof IDBDatabaseEventMap>(type: K, listener: (this: IDBDatabase, ev: IDBDatabaseEventMap[K]) => any, options?: boolean | EventListenerOptions): void;
		removeEventListener(type: string, listener: EventListenerOrEventListenerObject, options?: boolean | EventListenerOptions): void;
	}
		`

	dataBytes := []byte(data)

	if *input != "" {
		var err error
		dataBytes, err = os.ReadFile(*input)
		if err != nil {
			panic(err)
		}
	}

	return dataBytes
}

func main() {
	flag.Parse()

	data := data()

	results, err := dommer.LoadCache(data)

	if errors.Is(err, dommer.ErrNoCache) {
		results = dommer.Parse(data, *debugMax)
		dommer.PutCache(data, results)
	} else if err != nil {
		panic(err)
	}

	f, err := os.OpenFile(*output, os.O_CREATE|os.O_RDWR|os.O_TRUNC, 0600)

	if err != nil {
		panic(err)
	}

	defer f.Close()

	dommer.Generate(f, results)
}
