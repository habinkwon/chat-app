import React, { useState } from 'react'
import { render } from 'react-dom'
import { ApolloClient, createHttpLink, InMemoryCache, ApolloProvider, useQuery, gql } from '@apollo/client'
import { setContext } from '@apollo/client/link/context'

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

client
	.query({
		query: gql`
			query {
				chats {
					id
				}
			}
		`,
	})
	.then((result) => console.log(result))

function App() {
	const [state, setState] = useState('CLICK ME')

	return <button onClick={() => setState('CLICKED')}>{state}</button>
}

render(<App />, document.getElementById('root'))
