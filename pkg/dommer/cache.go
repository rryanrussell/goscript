package dommer

import (
	"bytes"
	"encoding/base64"
	"encoding/gob"
	"errors"
	"fmt"
	"hash/fnv"
	"os"

	"github.com/rryanrussell/goscript/pkg/dommer/ast"
)

func getHash(data []byte) string {
	hasher := fnv.New32a()
	_, err := hasher.Write(data)
	if err != nil {
		panic(err)
	}

	return base64.RawStdEncoding.EncodeToString(hasher.Sum(nil))
}

func cacheFile(hash string) string {
	return fmt.Sprintf("%s.parsecache", hash)
}

var ErrNoCache = errors.New("parse cache for this data not found")

func LoadCache(data []byte) ([]Result, error) {
	file := cacheFile(getHash(data))

	f, err := os.Open(file)

	if errors.Is(err, os.ErrNotExist) {
		return nil, ErrNoCache
	} else if err != nil {
		return nil, err
	}

	defer f.Close()

	dec := gob.NewDecoder(f)

	var result []Result
	if err := dec.Decode(&result); err != nil {
		return nil, err
	}

	return result, nil
}

func PutCache(data []byte, parsed []Result) error {
	var buf bytes.Buffer

	enc := gob.NewEncoder(&buf)

	if err := enc.Encode(parsed); err != nil {
		panic(err)
	}

	cacheLength := buf.Len()

	file := cacheFile(getHash(data))

	if err := os.WriteFile(file, buf.Bytes(), os.ModePerm&0600); err != nil {
		panic(err)
	}

	fmt.Printf("wrote %d bytes cache to %s\n\n", cacheLength, file)

	return nil
}

func init() {
	gob.Register(ast.DtsModule{})

	gob.Register(ast.GenericTypeExpr{})
	gob.Register(ast.TypeIdentExpr{})

	// bindIface(reg, typeOf[ast.InterfaceExpr](),
	gob.Register(ast.TypeIdentExpr{})
	gob.Register(ast.TypeBodyExpr{})
	gob.Register(ast.TupleExpr{})
	gob.Register(ast.LambdaTypeExpr{})
	gob.Register(ast.GenericTypeExpr{})
	gob.Register(ast.TypeIntersectionExpr{})
	gob.Register(ast.KeyofTypeExpr{})
	gob.Register(ast.TypeArrayExpr{})
	gob.Register(ast.TypeArrayExpr{})
	gob.Register(ast.ParenthesesTypeExpr{})
	gob.Register(ast.StringLiteralTypeExpr{})
	gob.Register(ast.TypeIdentExpr{})

	// bindIface(reg, typeOf[ast.TypeExpr](),
	gob.Register(ast.VariableDecl{})
	gob.Register(ast.SetterDecl{})
	gob.Register(ast.GetterDecl{})
	gob.Register(ast.MethodDecl{})
	gob.Register(ast.CallableDecl{})
	gob.Register(ast.IndexerDecl{})

	// bindIface(reg, typeOf[ast.TypeBodyDecl](),
	gob.Register(ast.InterfaceDecl{})
	gob.Register(ast.GlobalConstDecl{})
	gob.Register(ast.GlobalVarDecl{})
	gob.Register(ast.TypeDecl{})
	gob.Register(ast.NamespaceDecl{})
	gob.Register(ast.FunctionDecl{})

	// bindIface(reg, typeOf[ast.DtsDecl](),
	gob.Register(ast.StringIdent{})
	gob.Register(ast.KeywordIdent{})

	gob.Register(ast.TypeUnionExpr{})
	gob.Register(ast.TypeIndex{})
	gob.Register(ast.TypeofTypeExpr{})
}
