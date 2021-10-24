import React from 'react'
import { gql, useQuery } from '@apollo/client'

export default function ChatList({ setChatId }) {
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

	const handleClick = (chat) => {
		setChatId(chat.id)
	}

	if (loading) return <p>Loading...</p>
	if (error) return <p>{error}</p>
	return data.chats.map((chat) => (
		<div key={chat.id} onClick={() => handleClick(chat)} className="cursor-pointer hover:bg-gray-300">
			<div className="font-bold">{chat.name}</div>
			<div>{chat.messages?.[0]?.content}</div>
		</div>
	))
}
