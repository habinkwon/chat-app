import React, { useState, useEffect } from 'react'
import { gql, useApolloClient, useSubscription } from '@apollo/client'

export default function MessageList({ chatId, userId }) {
	const [messages, setMessages] = useState([])

	async function fetchMessages() {
		setMessages((prevMessages) => [...messages, ...prevMessages])
		setAfter(messages[0].id)
	}
	useEffect(() => fetchMessages(), [])

	const [userTyping, setUserTyping] = useState(null)
	useEffect(() => {
		if (!userTyping) return
		const timer = setTimeout(() => setUserTyping(null), 2000)
		return () => clearTimeout(timer)
	}, [userTyping])

	useSubscription({
		onSubscriptionData: ({ subscriptionData: { data } }) => {
			const event = data.chatEvent
			if (event.chatId !== chatId) return
			switch (event.type) {
				case 'MESSAGE_POSTED':
					setMessages((prevMessages) => [...prevMessages, event.message])
					break
				case 'USER_TYPING':
					setUserTyping(event.user?.name)
					break
			}
		},
	})

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
			{userTyping && <div>{userTyping} is typing...</div>}
		</>
	)
}
