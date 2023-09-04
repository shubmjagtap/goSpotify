import React, { useEffect, useState } from 'react';
import axios from 'axios';

function ChatPage() {

    const [chats, setChats] = useState([]);

    // Define the fetchChats function outside of the JSX
    const fetchChats = async () => {
        try {
        const { data } = await axios.get('http://localhost:8081/api/data');
        setChats(data);
        } catch (error) {
        console.error(error);
        }
    };

    // Use useEffect to call fetchChats when the component mounts
    useEffect(() => {
        fetchChats();
    }, []);

    return (
        <div className="App">
            {chats.map((chat) => (
                <div key={chat._id}>{chat.chatName}</div>
            ))}
        </div>
    );
}

export default ChatPage;
