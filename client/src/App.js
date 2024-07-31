import React from 'react';
import { BrowserRouter, Routes, Route } from 'react-router-dom';

import {
  Container,
  Flex,
  Textarea,
  Box,
  FormControl,
  FormErrorMessage,
  ModalFooter,
  InputGroup,
  InputRightElement,
  Button,
  Input,
  Tabs,
  TabList,
  Tab,
  TabPanels,
  TabPanel,
  VStack,
  Text,
  Modal,
  ModalOverlay,
  ModalContent,
  ModalHeader,
  ModalBody,
  ModalCloseButton,
  ChakraProvider,
} from '@chakra-ui/react';


import theme from './theme';
import { ColorModeSwitcher } from './ColorModeSwitcher';

import Header from './Components/Header';
import Landing from './Components/Landing';
import Register from './Components/Register';
import Login from './Components/Login';
import Chat from './Components/Chat/Chat';
import Footer from './Components/Footer';
import ActivityPage from './Components/Activitypage';

// theme.styles.global['font-family'] = 'roboto';
function App() {
  return (
    <ChakraProvider theme={theme}>
      <Flex direction="column" minH="100vh">
        {/* <Box textAlign="right">
          <ColorModeSwitcher justifySelf="flex-end" />
        </Box> */}
        <Box textAlign="center" fontSize="xl" flex="1">
          <BrowserRouter>
            <Header />
            <Routes>
              <Route path="/" element={<Landing />} />
              <Route path="/register" element={<Register />} />
              <Route path="/login" element={<Login />} />
              <Route path="/chat" element={<Chat />} />
              <Route path="/activity" element={<ActivityPage />} />
            </Routes>
          </BrowserRouter>
        </Box>
        <Footer />
      </Flex>
    </ChakraProvider>
  );
}

export default App;
