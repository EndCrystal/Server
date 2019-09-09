package block

func init() {
	Register(Block{
		Name:       "core:air",
		Attributes: AttributeNone,
	})
	Register(Block{
		Name:       "core:bedrock",
		Attributes: AttributeSolid,
	})
}
