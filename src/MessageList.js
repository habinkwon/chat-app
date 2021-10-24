import React from 'react'
import { gql, useQuery } from '@apollo/client'

function MessageList({ chatId }) {
	if (!chatId) return <p>Select a chat</p>
	const { loading, error, data } = useQuery(
		gql`
			query GetMessages($id: ID!) {
				chat(id: $id) {
					messages(first: 10, desc: true) {
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
		{ variables: { id: chatId } }
	)
	if (loading) return <p>Loading...</p>
	if (error) return <p>{error.message}</p>
	if (!data.chat) return <p>Chat not found</p>
	return data.chat.messages?.map((message) => (
		<div key={message.id} className="mb-2.5">
			<div>
				<span className="font-bold mr-2.5">{message.sender?.name}</span>
				<span className="text-gray-500">{message.createdAt}</span>
			</div>
			<div>{message.content}</div>
		</div>
	))
}

export default MessageList
