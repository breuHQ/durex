package workflows_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"go.breu.io/temporal-tools/workflows"
)

func TestWorkflowMod(t *testing.T) {
	workflow, _ := workflows.NewOptions(
		workflows.WithMod("mod"),
	)

	parent, err := workflow.ParentWorkflowID()
	suffix := workflow.IDSuffix()

	assert.Equal(t, "mod", suffix)
	assert.Equal(t, "", parent)
	assert.ErrorIs(t, err, workflows.ErrParentNil)
}

func TestWorkflowModWithID(t *testing.T) {
	workflow, _ := workflows.NewOptions(
		workflows.WithMod("mod"),
		workflows.WithModID("modid"),
	)

	parent, err := workflow.ParentWorkflowID()
	suffix := workflow.IDSuffix()

	assert.Equal(t, "mod.modid", suffix)
	assert.Equal(t, "", parent)
	assert.ErrorIs(t, err, workflows.ErrParentNil)
}

func TestWorkflowElement(t *testing.T) {
	workflow, _ := workflows.NewOptions(
		workflows.WithElement("elm"),
	)

	parent, err := workflow.ParentWorkflowID()
	suffix := workflow.IDSuffix()

	assert.Equal(t, "elm", suffix)
	assert.Equal(t, "", parent)
	assert.ErrorIs(t, err, workflows.ErrParentNil)
}

func TestWorkflowElementWithID(t *testing.T) {
	workflow, _ := workflows.NewOptions(
		workflows.WithElement("elm"),
		workflows.WithElementID("elmid"),
	)

	parent, err := workflow.ParentWorkflowID()
	suffix := workflow.IDSuffix()

	assert.Equal(t, "elm.elmid", suffix)
	assert.Equal(t, "", parent)
	assert.ErrorIs(t, err, workflows.ErrParentNil)
}

func TestWorkflowBlock(t *testing.T) {
	workflow, _ := workflows.NewOptions(
		workflows.WithBlock("block"),
	)

	parent, err := workflow.ParentWorkflowID()
	suffix := workflow.IDSuffix()

	assert.Equal(t, "block", suffix)
	assert.Equal(t, "", parent)
	assert.ErrorIs(t, err, workflows.ErrParentNil)
}

func TestWorkflowBlockWithID(t *testing.T) {
	workflow, _ := workflows.NewOptions(
		workflows.WithBlock("block"),
		workflows.WithBlockID("blockid"),
	)

	parent, err := workflow.ParentWorkflowID()
	suffix := workflow.IDSuffix()

	assert.Equal(t, "block.blockid", suffix)
	assert.Equal(t, "", parent)
	assert.ErrorIs(t, err, workflows.ErrParentNil)
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

	parent, err := workflow.ParentWorkflowID()
	suffix := workflow.IDSuffix()

	assert.Equal(t, "block.blockid.elm.elmid.mod.modid.prop.value", suffix)
	assert.Equal(t, "", parent)
	assert.ErrorIs(t, err, workflows.ErrParentNil)
}
