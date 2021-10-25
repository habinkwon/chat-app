import React from 'react'
import { render } from 'react-dom'
import { ApolloClient, createHttpLink, split, InMemoryCache, ApolloProvider } from '@apollo/client'
import { setContext } from '@apollo/client/link/context'
import { WebSocketLink } from '@apollo/client/link/ws'
import { getMainDefinition } from '@apollo/client/utilities'

import App from './App'
import './styles.css'

const httpLink = createHttpLink({
	uri: process.env.HTTP_URL,
})
const wsLink = new WebSocketLink({
	uri: process.env.WS_URL,
	options: {
		reconnect: true,
		connectionParams: () => {
			return {
				Authorization: 'Bearer 1',
				headers: {
					Authorization: 'Bearer 2',
					'user-id': 1,
				},
			}
		},
	},
})
const splitLink = split(
	({ query }) => {
		const def = getMainDefinition(query)
		return def.kind === 'OperationDefinition' && def.operation === 'subscription'
	},
	wsLink,
	httpLink
)
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
	link: authLink.concat(splitLink),
	cache: new InMemoryCache(),
})

render(
	<ApolloProvider client={client}>
		<App />
	</ApolloProvider>,
	document.getElementById('root')
)
