package lazyloader

import (
	"context"
	"testing"

	lsp "github.com/senforsce/tndr/lsp/protocol"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"golang.org/x/tools/go/packages"
)

func TestTndrPkgTraverserOpenTopologically(t *testing.T) {
	tests := []struct {
		name             string
		traverser        *goPkgTraverser
		pkg              *packages.Package
		wantPkgsRefCount map[string]int
		wantOpenedDocs   map[string]int
		wantClosedDocs   map[string]int
		wantErrContains  string
	}{
		{
			name: "file read failed",
			traverser: &goPkgTraverser{
				fileReader: mockFileReader{error: assert.AnError},
			},
			pkg: &packages.Package{
				OtherFiles: []string{"/foo.go", "/foo.t1"},
			},
			wantErrContains: "read file \"/foo.t1\"",
		},
		{
			name: "did open failed",
			traverser: &goPkgTraverser{
				fileReader: mockFileReader{error: nil},
				tndrDocHandler: &mockTndrDocHandler{
					openedDocs: map[string]int{},
					error:      assert.AnError,
				},
			},
			pkg: &packages.Package{
				OtherFiles: []string{"/foo.go", "/foo.t1"},
			},
			wantErrContains: "did open file \"/foo.t1\"",
		},
		{
			name: "open a->b->c",
			traverser: &goPkgTraverser{
				fileReader: mockFileReader{error: nil},
				tndrDocHandler: &mockTndrDocHandler{
					openedDocs: map[string]int{},
				},
				pkgsRefCount: map[string]int{},
			},
			pkg: &packages.Package{
				PkgPath:    "a",
				OtherFiles: []string{"/a.t1", "/a_other.t1", "/a.go"},
				Imports: map[string]*packages.Package{
					"b": {
						PkgPath:    "b",
						OtherFiles: []string{"/b.t1"},
						Imports: map[string]*packages.Package{
							"c": {
								PkgPath:    "c",
								OtherFiles: []string{"/c.t1", "/c_other.t1"},
							},
						},
					},
				},
			},
			wantPkgsRefCount: map[string]int{
				"a": 1,
				"b": 1,
				"c": 1,
			},
			wantOpenedDocs: map[string]int{
				"/c.t1":       1,
				"/c_other.t1": 1,
				"/b.t1":       1,
				"/a.t1":       1,
				"/a_other.t1": 1,
			},
		},
		{
			name: "open a->b->c with c already open",
			traverser: &goPkgTraverser{
				fileReader: mockFileReader{error: nil},
				tndrDocHandler: &mockTndrDocHandler{
					openedDocs: map[string]int{},
				},
				pkgsRefCount: map[string]int{
					"c": 1,
				},
			},
			pkg: &packages.Package{
				PkgPath:    "a",
				OtherFiles: []string{"/a.t1", "/a_other.t1"},
				Imports: map[string]*packages.Package{
					"b": {
						PkgPath:    "b",
						OtherFiles: []string{"/b.t1"},
						Imports: map[string]*packages.Package{
							"c": {
								PkgPath:    "c",
								OtherFiles: []string{"/c.t1", "/c_other.t1"},
							},
						},
					},
				},
			},
			wantPkgsRefCount: map[string]int{
				"a": 1,
				"b": 1,
				"c": 2,
			},
			wantOpenedDocs: map[string]int{
				"/b.t1":       1,
				"/a.t1":       1,
				"/a_other.t1": 1,
			},
		},
		{
			name: "open a->b->c a->d->c",
			traverser: &goPkgTraverser{
				fileReader: mockFileReader{error: nil},
				tndrDocHandler: &mockTndrDocHandler{
					openedDocs: map[string]int{},
				},
				pkgsRefCount: map[string]int{},
			},
			pkg: &packages.Package{
				PkgPath:    "a",
				OtherFiles: []string{"/a.t1", "/a_other.t1"},
				Imports: map[string]*packages.Package{
					"b": {
						PkgPath:    "b",
						OtherFiles: []string{"/b.t1"},
						Imports: map[string]*packages.Package{
							"c": {
								PkgPath:    "c",
								OtherFiles: []string{"/c.t1", "/c_other.t1"},
							},
						},
					},
					"d": {
						PkgPath:    "d",
						OtherFiles: []string{"/d.t1"},
						Imports: map[string]*packages.Package{
							"c": {
								PkgPath:    "c",
								OtherFiles: []string{"/c.t1", "/c_other.t1"},
							},
						},
					},
				},
			},
			wantPkgsRefCount: map[string]int{
				"a": 1,
				"b": 1,
				"c": 2,
				"d": 1,
			},
			wantOpenedDocs: map[string]int{
				"/d.t1":       1,
				"/c.t1":       1,
				"/c_other.t1": 1,
				"/b.t1":       1,
				"/a.t1":       1,
				"/a_other.t1": 1,
			},
		},
		{
			name: "open a->b->c a->d->c with d open",
			traverser: &goPkgTraverser{
				fileReader: mockFileReader{error: nil},
				tndrDocHandler: &mockTndrDocHandler{
					openedDocs: map[string]int{},
				},
				pkgsRefCount: map[string]int{
					"d": 1,
					"c": 1,
				},
			},
			pkg: &packages.Package{
				PkgPath:    "a",
				OtherFiles: []string{"/a.t1", "/a_other.t1"},
				Imports: map[string]*packages.Package{
					"b": {
						PkgPath:    "b",
						OtherFiles: []string{"/b.t1"},
						Imports: map[string]*packages.Package{
							"c": {
								PkgPath:    "c",
								OtherFiles: []string{"/c.t1", "/c_other.t1"},
							},
						},
					},
					"d": {
						PkgPath:    "d",
						OtherFiles: []string{"/d.t1"},
						Imports: map[string]*packages.Package{
							"c": {
								PkgPath:    "c",
								OtherFiles: []string{"/c.t1", "/c_other.t1"},
							},
						},
					},
				},
			},
			wantPkgsRefCount: map[string]int{
				"a": 1,
				"b": 1,
				"c": 1,
				"d": 2,
			},
			wantOpenedDocs: map[string]int{
				"/b.t1":       1,
				"/a.t1":       1,
				"/a_other.t1": 1,
			},
		},
		{
			name: "open a->b->c a->d->c with a open",
			traverser: &goPkgTraverser{
				fileReader: mockFileReader{error: nil},
				tndrDocHandler: &mockTndrDocHandler{
					openedDocs: map[string]int{},
				},
				pkgsRefCount: map[string]int{
					"a": 1,
					"b": 1,
					"d": 1,
					"c": 1,
				},
			},
			pkg: &packages.Package{
				PkgPath:    "a",
				OtherFiles: []string{"/a.t1", "/a_other.t1"},
				Imports: map[string]*packages.Package{
					"b": {
						PkgPath:    "b",
						OtherFiles: []string{"/b.t1"},
						Imports: map[string]*packages.Package{
							"c": {
								PkgPath:    "c",
								OtherFiles: []string{"/c.t1", "/c_other.t1"},
							},
						},
					},
					"d": {
						PkgPath:    "d",
						OtherFiles: []string{"/d.t1"},
						Imports: map[string]*packages.Package{
							"c": {
								PkgPath:    "c",
								OtherFiles: []string{"/c.t1", "/c_other.t1"},
							},
						},
					},
				},
			},
			wantPkgsRefCount: map[string]int{
				"a": 2,
				"b": 1,
				"c": 1,
				"d": 2,
			},
			wantOpenedDocs: map[string]int{},
		},
		{
			name: "open a->b->c a->d->c c->e with c open",
			traverser: &goPkgTraverser{
				fileReader: mockFileReader{error: nil},
				tndrDocHandler: &mockTndrDocHandler{
					openedDocs: map[string]int{},
				},
				pkgsRefCount: map[string]int{
					"c": 1,
					"e": 1,
				},
			},
			pkg: &packages.Package{
				PkgPath:    "a",
				OtherFiles: []string{"/a.t1", "/a_other.t1"},
				Imports: map[string]*packages.Package{
					"b": {
						PkgPath:    "b",
						OtherFiles: []string{"/b.t1"},
						Imports: map[string]*packages.Package{
							"c": {
								PkgPath:    "c",
								OtherFiles: []string{"/c.t1", "/c_other.t1"},
								Imports: map[string]*packages.Package{
									"e": {
										PkgPath:    "e",
										OtherFiles: []string{"/e.t1"},
									},
								},
							},
						},
					},
					"d": {
						PkgPath:    "d",
						OtherFiles: []string{"/d.t1"},
						Imports: map[string]*packages.Package{
							"c": {
								PkgPath:    "c",
								OtherFiles: []string{"/c.t1", "/c_other.t1"},
								Imports: map[string]*packages.Package{
									"e": {
										PkgPath:    "e",
										OtherFiles: []string{"/e.t1"},
									},
								},
							},
						},
					},
				},
			},
			wantPkgsRefCount: map[string]int{
				"a": 1,
				"b": 1,
				"c": 3,
				"d": 1,
				"e": 1,
			},
			wantOpenedDocs: map[string]int{
				"/d.t1":       1,
				"/b.t1":       1,
				"/a.t1":       1,
				"/a_other.t1": 1,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.traverser.openTopologically(context.Background(), tt.pkg)

			if tt.wantErrContains != "" {
				assert.Error(t, err)
				assert.ErrorContains(t, err, tt.wantErrContains)
			} else {
				assert.NoError(t, err)

				handler, ok := tt.traverser.tndrDocHandler.(*mockTndrDocHandler)
				require.True(t, ok)
				assert.Equal(t, tt.wantOpenedDocs, handler.openedDocs)
			}
		})
	}
}

func TestTndrPkgTraverserCloseTopologically(t *testing.T) {
	tests := []struct {
		name             string
		traverser        *goPkgTraverser
		pkg              *packages.Package
		wantPkgsRefCount map[string]int
		wantClosedDocs   map[string]int
		wantErrContains  string
	}{
		{
			name: "did close failed",
			traverser: &goPkgTraverser{
				fileReader: mockFileReader{error: nil},
				tndrDocHandler: &mockTndrDocHandler{
					closedDocs: map[string]int{},
					error:      assert.AnError,
				},
				pkgsRefCount: map[string]int{"foo": 1},
			},
			pkg: &packages.Package{
				PkgPath:    "foo",
				OtherFiles: []string{"/foo.go", "/foo.t1"},
			},
			wantErrContains: "did close file \"/foo.t1\"",
		},
		{
			name: "close a->b->c",
			traverser: &goPkgTraverser{
				fileReader: mockFileReader{error: nil},
				tndrDocHandler: &mockTndrDocHandler{
					closedDocs: map[string]int{},
				},
				pkgsRefCount: map[string]int{
					"a": 1, "b": 1, "c": 1,
				},
			},
			pkg: &packages.Package{
				PkgPath:    "a",
				OtherFiles: []string{"/a.t1", "/a_other.t1", "/a.go"},
				Imports: map[string]*packages.Package{
					"b": {
						PkgPath:    "b",
						OtherFiles: []string{"/b.t1"},
						Imports: map[string]*packages.Package{
							"c": {
								PkgPath:    "c",
								OtherFiles: []string{"/c.t1", "/c_other.t1"},
							},
						},
					},
				},
			},
			wantPkgsRefCount: map[string]int{},
			wantClosedDocs: map[string]int{
				"/c.t1":       1,
				"/c_other.t1": 1,
				"/b.t1":       1,
				"/a.t1":       1,
				"/a_other.t1": 1,
			},
		},
		{
			name: "close a->b->c with c open",
			traverser: &goPkgTraverser{
				fileReader: mockFileReader{error: nil},
				tndrDocHandler: &mockTndrDocHandler{
					closedDocs: map[string]int{},
				},
				pkgsRefCount: map[string]int{
					"a": 1, "b": 1, "c": 2,
				},
			},
			pkg: &packages.Package{
				PkgPath:    "a",
				OtherFiles: []string{"/a.t1", "/a_other.t1"},
				Imports: map[string]*packages.Package{
					"b": {
						PkgPath:    "b",
						OtherFiles: []string{"/b.t1"},
						Imports: map[string]*packages.Package{
							"c": {
								PkgPath:    "c",
								OtherFiles: []string{"/c.t1", "/c_other.t1"},
							},
						},
					},
				},
			},
			wantPkgsRefCount: map[string]int{"c": 1},
			wantClosedDocs: map[string]int{
				"/b.t1":       1,
				"/a.t1":       1,
				"/a_other.t1": 1,
			},
		},
		{
			name: "close a->b->c a->d->c",
			traverser: &goPkgTraverser{
				fileReader: mockFileReader{error: nil},
				tndrDocHandler: &mockTndrDocHandler{
					closedDocs: map[string]int{},
				},
				pkgsRefCount: map[string]int{
					"a": 1, "b": 1, "d": 1, "c": 2,
				},
			},
			pkg: &packages.Package{
				PkgPath:    "a",
				OtherFiles: []string{"/a.t1", "/a_other.t1"},
				Imports: map[string]*packages.Package{
					"b": {
						PkgPath:    "b",
						OtherFiles: []string{"/b.t1"},
						Imports: map[string]*packages.Package{
							"c": {
								PkgPath:    "c",
								OtherFiles: []string{"/c.t1", "/c_other.t1"},
							},
						},
					},
					"d": {
						PkgPath:    "d",
						OtherFiles: []string{"/d.t1"},
						Imports: map[string]*packages.Package{
							"c": {
								PkgPath:    "c",
								OtherFiles: []string{"/c.t1", "/c_other.t1"},
							},
						},
					},
				},
			},
			wantPkgsRefCount: map[string]int{},
			wantClosedDocs: map[string]int{
				"/d.t1":       1,
				"/b.t1":       1,
				"/a.t1":       1,
				"/a_other.t1": 1,
				"/c.t1":       1,
				"/c_other.t1": 1,
			},
		},
		{
			name: "close a->b->c a->d->c with d open",
			traverser: &goPkgTraverser{
				fileReader: mockFileReader{error: nil},
				tndrDocHandler: &mockTndrDocHandler{
					closedDocs: map[string]int{},
				},
				pkgsRefCount: map[string]int{
					"a": 1,
					"b": 1,
					"c": 2,
					"d": 2,
				},
			},
			pkg: &packages.Package{
				PkgPath:    "a",
				OtherFiles: []string{"/a.t1", "/a_other.t1"},
				Imports: map[string]*packages.Package{
					"b": {
						PkgPath:    "b",
						OtherFiles: []string{"/b.t1"},
						Imports: map[string]*packages.Package{
							"c": {
								PkgPath:    "c",
								OtherFiles: []string{"/c.t1", "/c_other.t1"},
							},
						},
					},
					"d": {
						PkgPath:    "d",
						OtherFiles: []string{"/d.t1"},
						Imports: map[string]*packages.Package{
							"c": {
								PkgPath:    "c",
								OtherFiles: []string{"/c.t1", "/c_other.t1"},
							},
						},
					},
				},
			},
			wantPkgsRefCount: map[string]int{
				"c": 1,
				"d": 1,
			},
			wantClosedDocs: map[string]int{
				"/a.t1":       1,
				"/a_other.t1": 1,
				"/b.t1":       1,
			},
		},
		{
			name: "close a->b->c a->d->c with a open",
			traverser: &goPkgTraverser{
				fileReader: mockFileReader{error: nil},
				tndrDocHandler: &mockTndrDocHandler{
					closedDocs: map[string]int{},
				},
				pkgsRefCount: map[string]int{
					"a": 2,
					"b": 1,
					"c": 1,
					"d": 2,
				},
			},
			pkg: &packages.Package{
				PkgPath:    "a",
				OtherFiles: []string{"/a.t1", "/a_other.t1"},
				Imports: map[string]*packages.Package{
					"b": {
						PkgPath:    "b",
						OtherFiles: []string{"/b.t1"},
						Imports: map[string]*packages.Package{
							"c": {
								PkgPath:    "c",
								OtherFiles: []string{"/c.t1", "/c_other.t1"},
							},
						},
					},
					"d": {
						PkgPath:    "d",
						OtherFiles: []string{"/d.t1"},
						Imports: map[string]*packages.Package{
							"c": {
								PkgPath:    "c",
								OtherFiles: []string{"/c.t1", "/c_other.t1"},
							},
						},
					},
				},
			},
			wantPkgsRefCount: map[string]int{
				"a": 1,
				"b": 1,
				"c": 1,
				"d": 2,
			},
			wantClosedDocs: map[string]int{},
		},
		{
			name: "close a->b->c a->d->c c->e with c open twice",
			traverser: &goPkgTraverser{
				fileReader: mockFileReader{error: nil},
				tndrDocHandler: &mockTndrDocHandler{
					closedDocs: map[string]int{},
				},
				pkgsRefCount: map[string]int{
					"a": 1,
					"b": 1,
					"c": 3,
					"d": 1,
					"e": 1,
				},
			},
			pkg: &packages.Package{
				PkgPath:    "a",
				OtherFiles: []string{"/a.t1", "/a_other.t1"},
				Imports: map[string]*packages.Package{
					"b": {
						PkgPath:    "b",
						OtherFiles: []string{"/b.t1"},
						Imports: map[string]*packages.Package{
							"c": {
								PkgPath:    "c",
								OtherFiles: []string{"/c.t1", "/c_other.t1"},
								Imports: map[string]*packages.Package{
									"e": {
										PkgPath:    "e",
										OtherFiles: []string{"/e.t1"},
									},
								},
							},
						},
					},
					"d": {
						PkgPath:    "d",
						OtherFiles: []string{"/d.t1"},
						Imports: map[string]*packages.Package{
							"c": {
								PkgPath:    "c",
								OtherFiles: []string{"/c.t1", "/c_other.t1"},
								Imports: map[string]*packages.Package{
									"e": {
										PkgPath:    "e",
										OtherFiles: []string{"/e.t1"},
									},
								},
							},
						},
					},
				},
			},
			wantPkgsRefCount: map[string]int{
				"c": 1,
				"e": 1,
			},
			wantClosedDocs: map[string]int{
				"/a.t1":       1,
				"/a_other.t1": 1,
				"/b.t1":       1,
				"/d.t1":       1,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			err := tt.traverser.closeTopologically(context.Background(), tt.pkg)

			if tt.wantErrContains != "" {
				assert.Error(t, err)
				assert.ErrorContains(t, err, tt.wantErrContains)
			} else {
				assert.NoError(t, err)

				handler, ok := tt.traverser.tndrDocHandler.(*mockTndrDocHandler)
				require.True(t, ok)
				assert.Equal(t, tt.wantClosedDocs, handler.closedDocs)
				assert.Equal(t, tt.wantPkgsRefCount, tt.traverser.pkgsRefCount)
			}
		})
	}
}

type mockFileReader struct {
	bytes []byte
	error error
}

func (r mockFileReader) read(_ string) ([]byte, error) {
	return r.bytes, r.error
}

type mockTndrDocHandler struct {
	openedDocs map[string]int
	closedDocs map[string]int
	error      error
}

func (h *mockTndrDocHandler) HandleDidOpen(_ context.Context, params *lsp.DidOpenTextDocumentParams) error {
	h.openedDocs[params.TextDocument.URI.Filename()]++
	return h.error
}

func (h *mockTndrDocHandler) HandleDidClose(_ context.Context, params *lsp.DidCloseTextDocumentParams) error {
	h.closedDocs[params.TextDocument.URI.Filename()]++
	return h.error
}
