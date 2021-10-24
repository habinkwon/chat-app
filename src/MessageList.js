import React, { useState, useEffect } from 'react'
import { gql, useApolloClient } from '@apollo/client'

function MessageList({ chatId }) {
	const [messages, setMessages] = useState([])
	const [after, setAfter] = useState(null)
	const [more, setMore] = useState(true)

	const client = useApolloClient()
	const [loading, setLoading] = useState(true)
	const [error, setError] = useState(null)

	async function fetchMessages() {
		setLoading(true)
		const result = await client.query({
			query: gql`
				query GetMessages($id: ID!, $after: ID) {
					chat(id: $id) {
						messages(first: 5, after: $after, desc: true) {
							id
							content
							sender {
								id
								name
							}
							createdAt
						}
					}
				}
			`,
			variables: { id: chatId, after },
		})
		setLoading(false)
		setError(result.error)
		const messages = Array.from(result.data?.chat?.messages ?? []).reverse()
		if (!messages.length) {
			setMore(false)
			return
		}
		setMessages((prevMessages) => [...messages, ...prevMessages])
		setAfter(messages[0].id)
	}

	useEffect(() => {
		fetchMessages()
	}, [])

	if (loading) return <p>Loading...</p>
	if (error) return <p>{error}</p>
	return (
		<>
			{more && (
				<a onClick={() => fetchMessages()} className="underline text-blue-600 hover:text-blue-800 cursor-pointer">
					Load more
				</a>
			)}
			{messages.map((message) => (
				<div key={message.id} className="mb-2.5">
					<div>
						<span className="font-bold mr-2.5">{message.sender?.name}</span>
						<span className="text-gray-500">{message.createdAt}</span>
					</div>
					<div>{message.content}</div>
				</div>
			))}
		</>
	)
}

export default MessageList
