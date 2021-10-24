import React, { useEffect, useState } from 'react'
import { gql, useMutation } from '@apollo/client'

export default function MessageComposer({ chatId }) {
	const [message, setMessage] = useState('')

	const [postMessage] = useMutation(gql`
		mutation PostMessage($chatId: ID!, $message: String!) {
			postMessage(chatId: $chatId, text: $message) {
				id
			}
		}
	`)

	const handleKeyPress = (e) => {
		if (e.key === 'Enter') {
			e.preventDefault()
			if (e.altKey) {
				setMessage((message) => message + '\r\n')
				return
			}
			postMessage({ variables: { chatId, message } })
			setMessage('')
		}
	}

	return (
		<textarea
			rows="1"
			className="border border-gray-800"
			placeholder="Type a message"
			value={message}
			onChange={(e) => setMessage(e.target.value)}
			onKeyPress={handleKeyPress}
		/>
	)
}
