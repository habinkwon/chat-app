render(
	<ApolloProvider client={client}>
		<App userId={userId} />
	</ApolloProvider>,
	document.getElementById('root')
)
