import React, { useState, useCallback } from 'react'
import { gql, useMutation } from '@apollo/client'
import throttle from 'lodash/throttle'

export default function MessageComposer({ chatId }) {
	const [message, setMessage] = useState('')

	const [userTyping] = useMutation(gql`
		mutation UserTyping($chatId: ID!) {
			userTyping(chatId: $chatId) {
				id
			}
		}
	`)
	const throttledUserTyping = useCallback(
		throttle(
			() =>
				userTyping({
					variables: { chatId },
				}),
			3000
		),
		[]
	)

	const handleKeyPress = (e) => {
		if (e.key === 'Enter') {
			e.preventDefault()
			if (e.altKey) {
				setMessage((message) => message + '\r\n')
				return
			} else if (!message.length) {
				return
			}
			postMessage({ variables: { chatId, message } })
			setMessage('')
		} else {
			throttledUserTyping()
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
