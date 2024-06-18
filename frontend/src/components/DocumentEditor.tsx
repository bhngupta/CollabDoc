import React, { useEffect, useState, useRef } from 'react';
import { useParams, useNavigate } from 'react-router-dom';

import './DocumentEditor.css';

const DocumentEditor: React.FC = () => {
  const { docID } = useParams<{ docID: string }>();
  const navigate = useNavigate();
  const [text, setText] = useState('');
  const ws = useRef<WebSocket | null>(null);
  const textareaRef = useRef<HTMLTextAreaElement>(null);

  useEffect(() => {
    if (!docID) {
      const newDocID = generateRandomDocID();
      navigate(`/${newDocID}`, { replace: true });
      return;
    }

    const socket = new WebSocket(`ws://localhost:8080/ws`);
    ws.current = socket;

    socket.onopen = () => {
      console.log('WebSocket connected');
      socket.send(JSON.stringify({ type: 'get', operation: { docID } }));
    };

    socket.onmessage = (event) => {
      const message = JSON.parse(event.data);
      console.log('Received message:', message);
      if (message.content !== undefined) {
        setText(message.content);
      }
    };

    socket.onerror = (error) => {
      console.error('WebSocket error:', error);
    };

    socket.onclose = (event) => {
      console.log('WebSocket connection closed:', event.code, event.reason);
    };

    return () => {
      if (ws.current) {
        ws.current.close();
      }
    };
  }, [docID, navigate]);

  const handleInputChange = (event: React.ChangeEvent<HTMLTextAreaElement>) => {
    const inputValue = event.target.value;
    const op = {
      docID,
      OpType: 'update',
      Pos: 0,
      Length: text.length,
      Content: inputValue,
    };

    setText(inputValue);

    const operation = {
      type: 'operation',
      operation: op,
    };

    if (ws.current && ws.current.readyState === WebSocket.OPEN) {
      ws.current.send(JSON.stringify(operation));
    }
  };

  return (
    <div className="editor-container">
      <header className="editor-header">
        <div className="logo">CollabDoc</div>
      </header>
      <main className="editor-main">
        <textarea
          ref={textareaRef}
          className="editor-textarea"
          value={text}
          onChange={handleInputChange}
        />
      </main>
      <footer className="editor-footer">
        <div className="footer-top">
          <span><b>Words:</b> {text.split(/\s+/).filter(Boolean).length} | <b>Chars:</b> {text.length}</span>
        </div>
      </footer>
    </div>
  );
};

const generateRandomDocID = () => {
  const characters = 'ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789';
  let result = '';
  for (let i = 0; i < 6; i++) {
    result += characters.charAt(Math.floor(Math.random() * characters.length));
  }
  return result;
};

export default DocumentEditor;
