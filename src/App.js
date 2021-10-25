import React, { useState, useEffect } from 'react'
import { gql, useMutation } from '@apollo/client'
import ChatList from './ChatList'
import MessageList from './MessageList'
import MessageComposer from './MessageComposer'

export default function App({ userId }) {
	const [chatId, setChatId] = useState(0)

	const [setAsOnline] = useMutation(gql`
		mutation SetAsOnline {
			setAsOnline {
				status
			}
		}
	`)
	useEffect(() => {
		const interval = setInterval(setAsOnline, 30 * 1000)
		setAsOnline()
		return () => clearInterval(interval)
	}, [])

	return (
		<div className="container">
			<div className="row">
				<div className="col-2">
					<ChatList setChatId={setChatId} />
				</div>
				<div className="col">
					{chatId ? (
						<>
							<MessageList chatId={chatId} userId={userId} />
							<MessageComposer chatId={chatId} />
						</>
					) : (
						<p>Select a chat</p>
					)}
				</div>
			</div>
		</div>
	)
}
