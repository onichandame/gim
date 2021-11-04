package graphql

type GraphqlModule struct{}

func (m *GraphqlModule) Providers() []interface{} {
	return []interface{}{newGraphqlService}
}
func (m *GraphqlModule) Exports() []interface{} {
	return []interface{}{newGraphqlService}
}
func (m *GraphqlModule) Controllers() []interface{} {
	return []interface{}{newGraphqlBuilderController}
}
