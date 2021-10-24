import React from 'react'
import { render } from 'react-dom'
import { ApolloClient, createHttpLink, InMemoryCache, ApolloProvider } from '@apollo/client'
import { setContext } from '@apollo/client/link/context'
import { WebSocketLink } from '@apollo/client/link/ws'
import App from './App'
import './styles.css'

const httpLink = createHttpLink({
	uri: 'http://localhost:8080/query',
})
const authLink = setContext((_, { headers }) => {
	const token = localStorage.getItem('token')
	return {
		headers: {
			...headers,
			authorization: token ? `Bearer ${token}` : '',
			'user-id': 1,
		},
	}
})
const client = new ApolloClient({
	link: authLink.concat(httpLink),
	cache: new InMemoryCache(),
})

render(
	<ApolloProvider client={client}>
		<App />
	</ApolloProvider>,
	document.getElementById('root')
)
