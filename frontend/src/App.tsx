import React from 'react';
import { BrowserRouter as Router, Routes, Route, Navigate } from 'react-router-dom';
import DocumentEditor from './components/DocumentEditor';

const generateRandomDocID = () => {
  const characters = 'ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789';
  let result = '';
  for (let i = 0; i < 6; i++) {
    result += characters.charAt(Math.floor(Math.random() * characters.length));
  }
  return result;
};

const App: React.FC = () => {
  const docID = generateRandomDocID();

  return (
    <Router>
      <Routes>
        <Route path="/:docID" element={<DocumentEditor />} />
        <Route path="/" element={<Navigate to={`/${docID}`} replace />} />
      </Routes>
    </Router>
  );
};

export default App;