import React from 'react';
import { BrowserRouter, Routes, Route, useLocation } from 'react-router-dom';

import {
  Container,
  Flex,
  Textarea,
  Box,

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
function MainLayout({ children }) {
  return (
    <>
      {/* <Header /> */}
      {children}
    </>
  );
}

function App() {
  return (
    <ChakraProvider theme={theme}>
      <Flex direction="column" minH="100vh">
        <Box textAlign="right">
          {/* <ColorModeSwitcher justifySelf="flex-end" /> */}
        </Box>
        <Box textAlign="center" fontSize="xl" flex="1">
          <BrowserRouter>
            <Routes>
              <Route path="/" element={<Landing />} />
              <Route
                path="/register"
                element={
                  <MainLayout>
                    <Register />
                  </MainLayout>
                }
              />
              <Route
                path="/login"
                element={
                  <MainLayout>
                    <Login />
                  </MainLayout>
                }
              />
              <Route
                path="/chat"
                element={
                  <MainLayout>
                    <Chat />
                  </MainLayout>
                }
              />
              <Route
                path="/activity"
                element={
                  <MainLayout>
                    <ActivityPage />
                  </MainLayout>
                }
              />
            </Routes>
          </BrowserRouter>
        </Box>
        <Footer />
      </Flex>
    </ChakraProvider>
  );
}

export default App;