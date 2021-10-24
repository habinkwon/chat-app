import React from 'react'
import { gql, useQuery } from '@apollo/client'

function ChatList({ setChatId }) {
	const { loading, error, data } = useQuery(gql`
		query GetChats {
			chats {
				id
				name
				messages(first: 1, desc: true) {
					content
				}
			}
		}
	`)
	if (loading) return <p>Loading...</p>
	if (error) return <p>{error.message}</p>
	return data.chats.map((chat) => (
		<div key={chat.id} onClick={() => setChatId(chat.id)} className="cursor-pointer hover:bg-gray-300">
			<div className="font-bold">{chat.name}</div>
			<div>{chat.messages?.[0]?.content}</div>
		</div>
	))
}

export default ChatList
