import { ApolloClient, createHttpLink, split, InMemoryCache, gql } from '@apollo/client/core'
import { setContext } from '@apollo/client/link/context'
import { WebSocketLink } from '@apollo/client/link/ws'
import { getMainDefinition } from '@apollo/client/utilities'

import './style.css'

export class ChatClient {
  constructor({ httpUrl, wsUrl, token }) {
    const httpLink = createHttpLink({
      uri: httpUrl,
    })
    const wsLink = new WebSocketLink({
      uri: wsUrl,
      options: {
        reconnect: true,
        connectionParams: () => {
          return {
            headers: {
              Authorization: token ? `Bearer ${token}` : '',
            },
          }
        },
      },
    })
    const splitLink = split(
      ({ query }) => {
        const def = getMainDefinition(query)
        return def.kind === 'OperationDefinition' && def.operation === 'subscription'
      },
      wsLink,
      httpLink
    )
    const authLink = setContext((_, { headers }) => {
      return {
        headers: {
          ...headers,
          Authorization: token ? `Bearer ${token}` : '',
        },
      }
    })
    this.client = new ApolloClient({
      link: authLink.concat(splitLink),
      cache: new InMemoryCache(),
    })
    this.token = token
  }

  async createChat(userId) {
    const result = await this.client.mutate({
      mutation: gql`
        mutation CreateChat($userId: ID!) {
          createChat(userIds: [$userId]) {
            id
          }
        }
      `,
      variables: {
        userId,
      },
    })
    return result.data?.createChat?.id
  }

  async deleteChat(chatId) {
    const result = await this.client.mutate({
      mutation: gql`
        mutation DeleteChat($chatId: ID!) {
          deleteChat(id: $chatId) {
            id
          }
        }
      `,
      variables: {
        chatId,
      },
    })
    return result.data?.deleteChat?.id
  }

  async getChats() {
    const result = await this.client.query({
      query: gql`
        query GetChats {
          chats {
            id
            name
            members {
              id
              name
              nickname
            }
            messages(first: 1, desc: true) {
              content
              createdAt
            }
          }
        }
      `,
    })
    const chats = []
    result.data.chats.forEach((chat) => {
      const members = chat.members?.map((member) => ({
        id: member.id,
        name: member.name,
        nickname: member.nickname,
      }))
      chats.push({
        id: chat.id,
        name: chat.name,
        members,
        lastMessage: chat.messages?.[0]?.content ?? '',
        lastPostedAt: chat.messages?.[0]?.createdAt ?? '',
      })
    })
    return chats
  }

  async postMessage({ chatId, message }) {
    const result = await this.client.mutate({
      mutation: gql`
        mutation PostMessage($chatId: ID!, $message: String!) {
          postMessage(chatId: $chatId, text: $message) {
            id
          }
        }
      `,
      variables: {
        chatId,
        message,
      },
    })
    return result.data?.postMessage?.id
  }

  async listMessages({ chatId, limit = 10, after }) {
    const result = await this.client.query({
      query: gql`
        query GetMessages($chatId: ID!, $limit: Int, $after: ID) {
          chat(id: $chatId) {
            messages(first: $limit, after: $after, desc: true) {
              id
              content
              sender {
                id
                name
                nickname
              }
              createdAt
            }
          }
        }
      `,
      variables: {
        chatId,
        limit,
        after,
      },
    })
    let messages = Array.from(result.data?.chat?.messages ?? []).reverse()
    messages = messages.map((message) => ({
      id: message.id,
      message: message.content,
      senderId: message.sender?.id,
      sender: message.sender?.nickname || message.sender?.name,
      senderName: message.sender?.name,
      senderNickname: message.sender?.nickname,
      createdAt: message.createdAt,
    }))
    return messages
  }

  async streamMessages({ onMessagePosted, onUserTyping }) {
    this.client
      .subscribe({
        query: gql`
          subscription OnChatEvent {
            chatEvent {
              type
              chatId
              message {
                id
                content
                sender {
                  id
                  name
                  nickname
                }
                createdAt
              }
              user {
                name
              }
            }
          }
        `,
      })
      .subscribe({
        next({ data }) {
          const event = data.chatEvent
          switch (event.type) {
            case 'MESSAGE_POSTED':
              const { message } = event
              onMessagePosted &&
                onMessagePosted({
                  chatId: event.chatId,
                  id: message.id,
                  message: message.content,
                  senderId: message.sender?.id,
                  sender: message.sender?.nickname || message.sender?.name,
                  senderName: message.sender?.name,
                  senderNickname: message.sender?.nickname,
                  createdAt: message.createdAt,
                })
              break
            case 'USER_TYPING':
              onUserTyping && onUserTyping(event.user?.name)
              break
          }
        },
      })
  }
}

export class ChatView {}
