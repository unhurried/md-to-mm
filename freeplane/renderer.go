package freeplane

import (
	"fmt"
	"io"

	bf "github.com/russross/blackfriday/v2"
)

type Renderer struct{}

/* Flags to store the current rendering status. */

// true while rendering the current paragraph node to fit in a node element.
var renderingAsNodeElement = false

// Store headding hierarchy so that child nodes can be placed under its headding.
var headdingStack = []int{0}

func (r *Renderer) RenderNode(w io.Writer, node *bf.Node, entering bool) bf.WalkStatus {
	switch node.Type {
	case bf.Document:
		if entering {
			r.out(w, []byte("<map version=\"freeplane 1.8.0\">"))
			r.out(w, []byte("<node TEXT=\"Result\">"))
		} else {
			for i := 0; i < len(headdingStack)-1; i++ {
				r.out(w, []byte("</node>"))
			}
			r.out(w, []byte("</node></map>"))
		}
	case bf.BlockQuote:
		if entering {
			r.out(w, []byte("<node><richcontent TYPE=\"NODE\"><html><body>"))
		} else {
			r.out(w, []byte("</body></html></richcontent></node>"))
		}
	case bf.List:
	case bf.Item:
		if !entering {
			r.out(w, []byte("</node>"))
		}
	case bf.Paragraph:
		parentType := node.Parent.Type
		switch parentType {
		case bf.Document:
			if entering {
				if renderableAsNodeElement(node) {
					renderingAsNodeElement = true
					r.renderParagraphAsNodeElement(w, node, entering)
				} else {
					r.out(w, []byte("<node><richcontent TYPE=\"NODE\"><html><body>"))
					renderAsHtml(w, node, entering)
				}
			} else {
				if !renderingAsNodeElement {
					renderAsHtml(w, node, entering)
					r.out(w, []byte("</body></html></richcontent>"))
				}
				r.out(w, []byte("</node>"))
				renderingAsNodeElement = false
			}
		case bf.BlockQuote:
			renderAsHtml(w, node, entering)
		case bf.Item:
			if entering {
				if renderableAsNodeElement(node) {
					renderingAsNodeElement = true
					r.renderParagraphAsNodeElement(w, node, entering)
				} else {
					r.out(w, []byte("<node><richcontent TYPE=\"NODE\"><html><body>"))
					renderAsHtml(w, node, entering)
				}
			} else {
				if !renderingAsNodeElement {
					renderAsHtml(w, node, entering)
					r.out(w, []byte("</body></html></richcontent>"))
				}
				renderingAsNodeElement = false
			}
		}
	case bf.Heading:
		if entering {
			// Search for the headding hierarchy whose level is equal or less than one of current node.
			for headdingStack[len(headdingStack)-1] >= node.HeadingData.Level {
				r.out(w, []byte("</node>"))
				headdingStack = headdingStack[:len(headdingStack)-1]
			}
			headdingStack = append(headdingStack, node.HeadingData.Level)

			r.out(w, []byte("<node><richcontent TYPE=\"NODE\"><html><body>"))
			renderAsHtml(w, node, entering)
		} else {
			renderAsHtml(w, node, entering)
			r.out(w, []byte("</body></html></richcontent>"))
		}
	case bf.HorizontalRule:
		renderAsHtml(w, node, entering)
	case bf.Emph:
		renderAsHtml(w, node, entering)
	case bf.Strong:
		renderAsHtml(w, node, entering)
	case bf.Del:
		renderAsHtml(w, node, entering)
	case bf.Link:
		renderAsHtml(w, node, entering)
	case bf.Image:
		renderAsHtml(w, node, entering)
	case bf.Text:
		renderAsHtml(w, node, entering)
	case bf.HTMLBlock:
		r.out(w, []byte("<node><richcontent TYPE=\"NODE\"><html><body>"))
		r.out(w, node.Literal)
		r.out(w, []byte("</body></html></richcontent></node>"))
	case bf.CodeBlock:
		r.out(w, []byte("<node><richcontent TYPE=\"NODE\"><html><body><pre><code>"))
		r.out(w, node.Literal)
		r.out(w, []byte("</code></pre></body></html></richcontent></node>"))
	case bf.Softbreak:
		r.out(w, []byte("<br />"))
	case bf.Hardbreak:
		r.out(w, []byte("<br />"))
	case bf.Code:
		renderAsHtml(w, node, entering)
	case bf.HTMLSpan:
		renderAsHtml(w, node, entering)
	case bf.Table:
		if entering {
			r.out(w, []byte("<node><richcontent TYPE=\"NODE\"><html><body><table>"))
		} else {
			r.out(w, []byte("</table></body></html></richcontent></node>"))
		}
	case bf.TableCell:
		renderAsHtml(w, node, entering)
	case bf.TableHead:
		renderAsHtml(w, node, entering)
	case bf.TableBody:
		renderAsHtml(w, node, entering)
	case bf.TableRow:
		renderAsHtml(w, node, entering)
	default:
		panic("Unknown node type " + node.Type.String())
	}
	return bf.GoToNext
}

func (r *Renderer) RenderHeader(w io.Writer, ast *bf.Node) {
}

func (r *Renderer) RenderFooter(w io.Writer, ast *bf.Node) {
}

func (r *Renderer) out(w io.Writer, text []byte) {
	w.Write(text)
}

var htmlRenderer = bf.NewHTMLRenderer(bf.HTMLRendererParameters{})

// Render a node in a html format with NewHTMLRenderer.
func renderAsHtml(w io.Writer, node *bf.Node, entering bool) {
	if !renderingAsNodeElement {
		htmlRenderer.RenderNode(w, node, entering)
	}
}

// Retrieve children of a given node excluding empty text nodes.
func getChildren(node *bf.Node) []*bf.Node {
	children := []*bf.Node{}
	for child := node.FirstChild; child != nil; child = child.Next {
		// Skip empty text nodes, that don't have literals.
		if child.Type != bf.Text || len(child.Literal) > 0 {
			children = append(children, child)
		}
	}
	return children
}

// True if a node can be rendered into an node element in Freeplane file.
// (can be expressed in the following format: <node TEXT="..." LINK="..."></node>)
func renderableAsNodeElement(node *bf.Node) bool {
	children := getChildren(node)
	return len(children) == 1 && (isSingleTextNode(children[0]) || isSingleLinkNode(children[0]))
}

// Return true if the node is a single text node.
func isSingleTextNode(node *bf.Node) bool {
	if node == nil || node.Type != bf.Text {
		return false
	}

	children := getChildren(node)
	return len(children) == 0
}

// Return true if the node is a link node that only has a child text node.
func isSingleLinkNode(node *bf.Node) bool {
	if node == nil || node.Type != bf.Link {
		return false
	}

	children := getChildren(node)
	return len(children) == 1 && children[0].Type == bf.Text
}

// Render a paragraph node to fit in a node element.
func (r *Renderer) renderParagraphAsNodeElement(w io.Writer, node *bf.Node, entering bool) {
	if entering {
		text := ""
		link := ""

		// Search ancestors of the node and concatenate text nodes and extract a link destination.
		for c1 := node.FirstChild; c1 != nil; c1 = c1.Next {
			for c2 := c1; c2 != nil; c2 = c2.FirstChild {
				text += string(c2.Literal)
				if c2.Destination != nil {
					link = string(c2.Destination)
				}
			}
		}

		if link == "" {
			r.out(w, []byte(fmt.Sprintf("<node TEXT=\"%s\">", text)))
		} else {
			r.out(w, []byte(fmt.Sprintf("<node TEXT=\"%s\" LINK=\"%s\">", text, link)))
		}
	}
}
