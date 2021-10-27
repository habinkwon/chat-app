import { ApolloClient, createHttpLink, split, InMemoryCache, gql } from '@apollo/client/core'
import { setContext } from '@apollo/client/link/context'
import { WebSocketLink } from '@apollo/client/link/ws'
import { getMainDefinition } from '@apollo/client/utilities'

import './style.css'

export class ChatClient {
	constructor({ httpUrl, wsUrl, token }) {
		const userId = token
		const httpLink = createHttpLink({
			uri: httpUrl,
		})
		const wsLink = new WebSocketLink({
			uri: wsUrl,
			options: {
				reconnect: true,
				connectionParams: () => {
					return {
						headers: {
							'user-id': userId,
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
			return {
				headers: {
					...headers,
					authorization: token ? `Bearer ${token}` : '',
					'user-id': userId,
				},
			}
		})
		this.client = new ApolloClient({
			link: authLink.concat(splitLink),
			cache: new InMemoryCache(),
		})
	}
}

export class ChatView {}
