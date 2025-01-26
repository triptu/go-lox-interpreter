package lox

import (
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"strings"
)

/*
Generated with Claude to visualise the AST. It creates a DOT file which is
used with Graphviz to generate a PNG and then open it.
*/

type visualiseTreeVisitor struct {
	nodeCount int
	builder   strings.Builder
}

func (v *visualiseTreeVisitor) Visualize(e expr, outputPath string) error {
	// Generate DOT content
	e.accept(v)
	v.builder.WriteString("}\n")

	// Write DOT file
	dotFile := outputPath + ".dot"
	err := os.WriteFile(dotFile, []byte(v.builder.String()), 0644)
	if err != nil {
		return fmt.Errorf("failed to write DOT file: %v", err)
	}

	// Generate PNG using Graphviz with higher DPI for better quality
	pngFile := outputPath + ".png"
	cmd := exec.Command("dot", "-Tpng", "-Gdpi=300", dotFile, "-o", pngFile)
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to generate PNG: %v", err)
	}

	// Open the image
	if err := openImage(pngFile); err != nil {
		return fmt.Errorf("failed to open image: %v", err)
	}

	return nil
}

func NewVisualiseTreeVisitor() *visualiseTreeVisitor {
	v := &visualiseTreeVisitor{}
	// Enhanced styling for the graph
	v.builder.WriteString("digraph AST {\n")
	v.builder.WriteString("  bgcolor=\"#FFFFFF\";\n")
	v.builder.WriteString("  node [fontname=\"Arial\", fontsize=12, shape=box, style=\"rounded,filled\", margin=0.2];\n")
	v.builder.WriteString("  edge [color=\"#2B2B2B\", penwidth=1.2];\n")
	return v
}

func (v *visualiseTreeVisitor) getNextNodeID() string {
	v.nodeCount++
	return fmt.Sprintf("node%d", v.nodeCount)
}

// Get node color based on type
func (v *visualiseTreeVisitor) getNodeStyle(nodeType string) (string, string) {
	switch nodeType {
	case "Binary", "Unary", "Logical":
		return "#E3F2FD", "#1565C0" // Light blue background, dark blue text
	case "Literal":
		return "#F3E5F5", "#6A1B9A" // Light purple background, dark purple text
	case "Variable", "This":
		return "#E8F5E9", "#2E7D32" // Light green background, dark green text
	case "Assign", "Set":
		return "#FFF3E0", "#E65100" // Light orange background, dark orange text
	case "Call", "Get":
		return "#F3E5F5", "#4A148C" // Light purple background, darker purple text
	case "Group":
		return "#FAFAFA", "#424242" // Light gray background, dark gray text
	case "Super":
		return "#FCE4EC", "#880E4F" // Light pink background, dark pink text
	default:
		return "#FFFFFF", "#000000" // White background, black text
	}
}

func (v *visualiseTreeVisitor) addNode(nodeID, nodeType, label string) {
	bgColor, textColor := v.getNodeStyle(nodeType)
	v.builder.WriteString(fmt.Sprintf("  %s [label=%q, fillcolor=%q, fontcolor=%q];\n",
		nodeID, label, bgColor, textColor))
}

func (v *visualiseTreeVisitor) addEdge(fromID, toID string) {
	v.builder.WriteString(fmt.Sprintf("  %s -> %s;\n", fromID, toID))
}

// openImage opens the generated image file using the default system viewer
func openImage(imagePath string) error {
	var cmd *exec.Cmd

	switch runtime.GOOS {
	case "darwin": // macOS
		cmd = exec.Command("open", imagePath)
	case "windows":
		cmd = exec.Command("cmd", "/c", "start", imagePath)
	default: // Linux and other Unix-like
		cmd = exec.Command("xdg-open", imagePath)
	}

	return cmd.Run()
}

func (v *visualiseTreeVisitor) visitAssignExpr(e eAssign) (any, error) {
	nodeID := v.getNextNodeID()
	v.addNode(nodeID, "Assign", fmt.Sprintf("Assign\n%s", e.name.lexeme))

	valueID := getVal(e.value.accept(v)).(string)
	v.addEdge(nodeID, valueID)

	return nodeID, nil
}

func (v *visualiseTreeVisitor) visitBinaryExpr(e eBinary) (any, error) {
	nodeID := v.getNextNodeID()
	v.addNode(nodeID, "Binary", fmt.Sprintf("Binary\n%s", e.operator.lexeme))

	leftID := getVal(e.left.accept(v)).(string)
	rightID := getVal(e.right.accept(v)).(string)

	v.addEdge(nodeID, leftID)
	v.addEdge(nodeID, rightID)

	return nodeID, nil
}

func (v *visualiseTreeVisitor) visitCallExpr(e eCall) (any, error) {
	nodeID := v.getNextNodeID()
	v.addNode(nodeID, "Call", "Call")

	calleeID := getVal(e.callee.accept(v)).(string)
	v.addEdge(nodeID, calleeID)

	for _, arg := range e.arguments {
		argID := getVal(arg.accept(v)).(string)
		v.addEdge(nodeID, argID)
	}

	return nodeID, nil
}

func (v *visualiseTreeVisitor) visitGetExpr(e eGet) (any, error) {
	nodeID := v.getNextNodeID()
	v.addNode(nodeID, "Get", fmt.Sprintf("Get\n%s", e.name.lexeme))

	objectID := getVal(e.object.accept(v)).(string)
	v.addEdge(nodeID, objectID)

	return nodeID, nil
}

func (v *visualiseTreeVisitor) visitGroupingExpr(e eGrouping) (any, error) {
	nodeID := v.getNextNodeID()
	v.addNode(nodeID, "Group", "Group")

	exprID := getVal(e.expression.accept(v)).(string)
	v.addEdge(nodeID, exprID)

	return nodeID, nil
}

func (v *visualiseTreeVisitor) visitLiteralExpr(e eLiteral) (any, error) {
	nodeID := v.getNextNodeID()
	v.addNode(nodeID, "Literal", fmt.Sprintf("Literal\n%v", getTokenLiteralStr(e.value)))
	return nodeID, nil
}

func (v *visualiseTreeVisitor) visitLogicalExpr(e eLogical) (any, error) {
	nodeID := v.getNextNodeID()
	v.addNode(nodeID, "Logical", fmt.Sprintf("Logical\n%s", e.operator.lexeme))

	leftID := getVal(e.left.accept(v)).(string)
	rightID := getVal(e.right.accept(v)).(string)

	v.addEdge(nodeID, leftID)
	v.addEdge(nodeID, rightID)

	return nodeID, nil
}

func (v *visualiseTreeVisitor) visitSetExpr(e eSet) (any, error) {
	nodeID := v.getNextNodeID()
	v.addNode(nodeID, "Set", fmt.Sprintf("Set\n%s", e.name.lexeme))

	objectID := getVal(e.object.accept(v)).(string)
	valueID := getVal(e.value.accept(v)).(string)

	v.addEdge(nodeID, objectID)
	v.addEdge(nodeID, valueID)

	return nodeID, nil
}

func (v *visualiseTreeVisitor) visitSuperExpr(e eSuper) (any, error) {
	nodeID := v.getNextNodeID()
	v.addNode(nodeID, "Super", fmt.Sprintf("Super\n%s.%s", e.keyword.lexeme, e.method.lexeme))
	return nodeID, nil
}

func (v *visualiseTreeVisitor) visitThisExpr(e eThis) (any, error) {
	nodeID := v.getNextNodeID()
	v.addNode(nodeID, "This", "This")
	return nodeID, nil
}

func (v *visualiseTreeVisitor) visitUnaryExpr(e eUnary) (any, error) {
	nodeID := v.getNextNodeID()
	v.addNode(nodeID, "Unary", fmt.Sprintf("Unary\n%s", e.operator.lexeme))

	rightID := getVal(e.right.accept(v)).(string)
	v.addEdge(nodeID, rightID)

	return nodeID, nil
}

func (v *visualiseTreeVisitor) visitVariableExpr(e eVariable) (any, error) {
	nodeID := v.getNextNodeID()
	v.addNode(nodeID, "Variable", fmt.Sprintf("Variable\n%s", e.name.lexeme))
	return nodeID, nil
}

func (v *visualiseTreeVisitor) visitListExpr(e eList) (any, error) {
	panic("not implemented")
}

func (v *visualiseTreeVisitor) visitGetIndexExpr(e eGetIndex) (any, error) {
	panic("not implemented")
}

func (v *visualiseTreeVisitor) visitSetIndexExpr(e eSetIndex) (any, error) {
	panic("not implemented")
}

func getVal(val any, _ error) any {
	return val
}
