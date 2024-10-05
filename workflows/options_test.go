package workflows_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"go.breu.io/durex/workflows"
)

func TestWorkflowMod(t *testing.T) {
	workflow, _ := workflows.NewOptions(
		workflows.WithMod("mod"),
	)

	parent := workflow.ParentWorkflowID()
	suffix := workflow.IDSuffix()

	assert.Equal(t, "mod", suffix)
	assert.Equal(t, "", parent)
}

func TestWorkflowModWithID(t *testing.T) {
	workflow, _ := workflows.NewOptions(
		workflows.WithMod("mod"),
		workflows.WithModID("modid"),
	)

	parent := workflow.ParentWorkflowID()
	suffix := workflow.IDSuffix()

	assert.Equal(t, "mod.modid", suffix)
	assert.Equal(t, "", parent)
}

func TestWorkflowElement(t *testing.T) {
	workflow, _ := workflows.NewOptions(
		workflows.WithElement("elm"),
	)

	parent := workflow.ParentWorkflowID()
	suffix := workflow.IDSuffix()

	assert.Equal(t, "elm", suffix)
	assert.Equal(t, "", parent)
}

func TestWorkflowElementWithID(t *testing.T) {
	workflow, _ := workflows.NewOptions(
		workflows.WithElement("elm"),
		workflows.WithElementID("elmid"),
	)

	parent := workflow.ParentWorkflowID()
	suffix := workflow.IDSuffix()

	assert.Equal(t, "elm.elmid", suffix)
	assert.Equal(t, "", parent)
}

func TestWorkflowBlock(t *testing.T) {
	workflow, _ := workflows.NewOptions(
		workflows.WithBlock("block"),
	)

	parent := workflow.ParentWorkflowID()
	suffix := workflow.IDSuffix()

	assert.Equal(t, "block", suffix)
	assert.Equal(t, "", parent)
}

func TestWorkflowBlockWithID(t *testing.T) {
	workflow, _ := workflows.NewOptions(
		workflows.WithBlock("block"),
		workflows.WithBlockID("blockid"),
	)

	parent := workflow.ParentWorkflowID()
	suffix := workflow.IDSuffix()

	assert.Equal(t, "block.blockid", suffix)
	assert.Equal(t, "", parent)
}

func TestWithProp(t *testing.T) {
	workflow, _ := workflows.NewOptions(
		workflows.WithBlock("block"),
		workflows.WithBlockID("blockid"),
		workflows.WithElement("elm"),
		workflows.WithElementID("elmid"),
		workflows.WithMod("mod"),
		workflows.WithModID("modid"),
		workflows.WithProp("prop", "value"),
	)

	parent := workflow.ParentWorkflowID()
	suffix := workflow.IDSuffix()

	assert.Equal(t, "block.blockid.elm.elmid.mod.modid.prop.value", suffix)
	assert.Equal(t, "", parent)
}
func TestWithPropMultiple(t *testing.T) {
	workflow, _ := workflows.NewOptions(
		workflows.WithProp("prop1", "value1"),
		workflows.WithProp("prop2", "value2"),
		workflows.WithProp("prop3", "value3"),
	)

	suffix := workflow.IDSuffix()

	assert.Equal(t, "prop1.value1.prop2.value2.prop3.value3", suffix)
}

func TestWithPropEmptyValue(t *testing.T) {
	workflow, _ := workflows.NewOptions(
		workflows.WithProp("prop", ""),
	)

	suffix := workflow.IDSuffix()

	assert.Equal(t, "prop", suffix)
}
