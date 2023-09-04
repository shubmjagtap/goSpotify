import logo from './logo.svg';
import { ChakraProvider, Button, ButtonGroup } from '@chakra-ui/react'
import { BrowserRouter as Router, Route, Switch } from 'react-router-dom';
import { BrowserRouter } from 'react-router-dom/cjs/react-router-dom.min';
import HomePage from "./Pages/HomePage"
import ChatPage from "./Pages/ChatPage"

function App() {
  return (
    <BrowserRouter>
      <ChakraProvider>
        <div className="App">
          <Route path="/" component={HomePage} exact />
          <Route path="/chats" component={ChatPage} />
        </div>
      </ChakraProvider>
    </BrowserRouter>
  );
}

export default App;
