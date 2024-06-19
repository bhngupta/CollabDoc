import React, { useEffect, useState, useRef } from 'react';
import { useParams, useNavigate } from 'react-router-dom';
import './DocumentEditor.css';

const DocumentEditor: React.FC = () => {
  const { docID } = useParams<{ docID: string }>();
  const navigate = useNavigate();
  const [text, setText] = useState('');
  const [prevText, setPrevText] = useState('');
  const ws = useRef<WebSocket | null>(null);
  const textareaRef = useRef<HTMLTextAreaElement>(null);

  useEffect(() => {
    if (!docID) {
      const newDocID = generateRandomDocID();
      navigate(`/${newDocID}`, { replace: true });
      return;
    }

    const socket = new WebSocket(`ws://localhost:8080/ws?docID=${docID}`);
    ws.current = socket;

    socket.onopen = () => {
      console.log('WebSocket connected');
      // Only Sending GET Type of operation
      socket.send(JSON.stringify({ type: 'get', operation: { docID } }));
    };

    socket.onmessage = (event) => {
      const message = JSON.parse(event.data);
      console.log('Received message:', message);
      if (message.type === 'heartbeat_ack') {
        // Handle heartbeat acknowledgement if needed
      } else if (message.content !== undefined) {
        setText(message.content);
        setPrevText(message.content); // Ensure previous text is set to the current text
      }
    };    

    socket.onerror = (error) => {
      console.error('WebSocket error:', error);
    };

    socket.onclose = (event) => {
      console.log('WebSocket connection closed:', event.code, event.reason);
    };
    
    // Heartbeat interval setup
    const heartbeatInterval = setInterval(() => {
      if (ws.current && ws.current.readyState === WebSocket.OPEN) {
        ws.current.send(JSON.stringify({ type: 'heartbeat' }));
      }
    }, 500); // Send heartbeat every second (1000ms)

    return () => {
      clearInterval(heartbeatInterval);

      if (ws.current && ws.current.readyState === WebSocket.OPEN) {
        ws.current.close();
      }
    };
  }, [docID, navigate]);

  const handleInputChange = (event: React.ChangeEvent<HTMLTextAreaElement>) => {
    const inputValue = event.target.value;
    const selectionStart = event.target.selectionStart;
    const prevLength = prevText.length;
    const currLength = inputValue.length;

    let opType = '';
    let pos = 0;
    let length = 0;
    let content = '';

    if (currLength > prevLength) {
      // Insertion
      opType = 'insert';
      pos = selectionStart - (currLength - prevLength);
      length = 0;
      content = inputValue.slice(pos, selectionStart);
    } else if (currLength < prevLength) {
    
        opType = 'delete';
        pos = selectionStart;
        length = prevLength - currLength;
      
    }

    if (opType) {
      const op = {
        docID,
        OpType: opType,
        Pos: pos,
        Length: length,
        Content: content,
      };

      console.log('Generated operation:', op); // Log the operation

      if (ws.current && ws.current.readyState === WebSocket.OPEN) {
        ws.current.send(JSON.stringify({ type: 'operation', operation: op }));
      }
    }

    setText(inputValue);
    setPrevText(inputValue);
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
