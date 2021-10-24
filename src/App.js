import React, { useState } from 'react'
import ChatList from './ChatList'
import MessageList from './MessageList'

function App() {
	const [chatId, setChatId] = useState(0)

	return (
		<div className="container">
			<div className="row">
				<div className="col-2">
					<ChatList setChatId={setChatId} />
				</div>
				<div className="col">
					<MessageList chatId={chatId} />
				</div>
			</div>
		</div>
	)
}

export default App
