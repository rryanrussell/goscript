package dommer

import (
	"bytes"
	"strings"
	"testing"

	"github.com/rryanrussell/goscript/pkg/dommer/ast"
)

func TestGenerate_KeywordEscaping(t *testing.T) {
	program := &ast.DtsModule{
		Decls: []ast.DtsDecl{
			&ast.InterfaceDecl{
				Ident: &ast.Ident{Name: &ast.TokenExpr{Value: "TestInterface"}},
				Body: &ast.TypeBodyExpr{
					Members: []ast.TypeBodyDecl{
						&ast.MethodDecl{
							Ident: &ast.Ident{Name: &ast.TokenExpr{Value: "method"}},
							Params: []*ast.ParamExpr{
								{
									Ident: &ast.KeywordIdent{Name: &ast.TokenExpr{Value: "range"}},
									Type:  &ast.TypeAnnotation{Type: &ast.TypeIdentExpr{Name: &ast.Ident{Name: &ast.TokenExpr{Value: "number"}}}},
								},
							},
							Return: &ast.TypeAnnotation{Type: &ast.TypeIdentExpr{Name: &ast.Ident{Name: &ast.TokenExpr{Value: "void"}}}},
						},
					},
				},
			},
		},
	}

	var buf bytes.Buffer
	Generate(&buf, []Result{{Value: program}})
	output := buf.String()

	if strings.Contains(output, "range float64") {
		t.Errorf("Output contains unescaped keyword 'range':\n%s", output)
	}
}

func TestGenerate_Conflict(t *testing.T) {
	program := &ast.DtsModule{
		Decls: []ast.DtsDecl{
			&ast.InterfaceDecl{
				Ident: &ast.Ident{Name: &ast.TokenExpr{Value: "Event"}},
				Body:  &ast.TypeBodyExpr{},
			},
			&ast.GlobalVarDecl{
				Variable: &ast.VariableDecl{
					Ident: &ast.StringIdent{Name: &ast.TokenExpr{Value: "Event"}},
					Type:  &ast.TypeAnnotation{Type: &ast.TypeIdentExpr{Name: &ast.Ident{Name: &ast.TokenExpr{Value: "Event"}}}},
				},
			},
		},
	}

	var buf bytes.Buffer
	Generate(&buf, []Result{{Value: program}})
	output := buf.String()
	t.Log("Forward order output:\n", output)

	lines := strings.Split(output, "\n")
	hasTypeEvent := false
	hasVarEvent := false

	for _, line := range lines {
		if strings.HasPrefix(strings.TrimSpace(line), "type Event struct") {
			hasTypeEvent = true
		}
		if strings.HasPrefix(strings.TrimSpace(line), "var Event ") {
			hasVarEvent = true
		}
	}

	if hasTypeEvent && hasVarEvent {
		t.Errorf("Output contains conflicting definitions for Event:\n%s", output)
	}
}

func TestGenerate_Conflict_Reverse(t *testing.T) {
	program := &ast.DtsModule{
		Decls: []ast.DtsDecl{
			&ast.GlobalVarDecl{
				Variable: &ast.VariableDecl{
					Ident: &ast.StringIdent{Name: &ast.TokenExpr{Value: "Event"}},
					Type:  &ast.TypeAnnotation{Type: &ast.TypeIdentExpr{Name: &ast.Ident{Name: &ast.TokenExpr{Value: "Event"}}}},
				},
			},
			&ast.InterfaceDecl{
				Ident: &ast.Ident{Name: &ast.TokenExpr{Value: "Event"}},
				Body:  &ast.TypeBodyExpr{},
			},
		},
	}

	var buf bytes.Buffer
	Generate(&buf, []Result{{Value: program}})
	output := buf.String()
	t.Log("Reverse order output:\n", output)

	lines := strings.Split(output, "\n")
	hasTypeEvent := false
	hasVarEvent := false

	for _, line := range lines {
		if strings.HasPrefix(strings.TrimSpace(line), "type Event struct") {
			hasTypeEvent = true
		}
		if strings.HasPrefix(strings.TrimSpace(line), "var Event ") {
			hasVarEvent = true
		}
	}

	if hasTypeEvent && hasVarEvent {
		t.Errorf("Output contains conflicting definitions for Event:\n%s", output)
	}
}
