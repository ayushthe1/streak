import React from 'react';

import { Text, Box, Container,VStack } from '@chakra-ui/react';
import './Chat.css';

const ChatHistory = (currentUser, chats) => {
  const history = chats.map(m => {
    // incoming message on left side
    let margin = '0%';
    let bgcolor = 'darkgray';
    let textAlign = 'left';
    let isCurrentUser = false

    if (m.from === currentUser) {
      // outgoing message to right
      margin = '20%';
      bgcolor = 'purple.400';
      textAlign = 'right';
      isCurrentUser = true
    }

    const ts = new Date(m.timestamp * 1000);

    return (
      <Box
        key={m.id}
        textAlign={textAlign}
        width={'80%'}
        p={2}
        m={2}
        marginTop={2}
        marginBottom={2}
        marginLeft={margin}
        paddingRight={2}
        boxShadow="md"
        bg={bgcolor}
        borderRadius={20}
        // transform={isCurrentUser ? 'translateX(20%)' : 'translateX(-20%)'}
        _hover={{ bg: isCurrentUser ? 'purple.500' : 'gray.600' }}
        transition="background-color 0.2s"
      >
        <Text> {m.message} </Text>
        <Text as="sub" fontSize="xs" color="gray.300" mt={2}>
          {ts.toUTCString()}
        </Text>
      </Box>
    );
  });

  return <Container>{history}</Container>;
};

export default ChatHistory;


// import React from 'react';
// import { Text, Box, Container, VStack } from '@chakra-ui/react';

// const ChatHistory = ({ currentUser, chats }) => {
//   const history = chats.map((m) => {
//     const isCurrentUser = m.from === currentUser;
//     const ts = new Date(m.timestamp * 1000);

//     return (
//       <Box
//         key={m.id}
//         alignSelf={isCurrentUser ? 'flex-end' : 'flex-start'}
//         bg={isCurrentUser ? 'purple.400' : 'gray.700'}
//         color="white"
//         borderRadius="lg"
//         p={4}
//         m={2}
//         maxWidth="80%"
//         boxShadow="md"
//         transform={isCurrentUser ? 'translateX(20%)' : 'translateX(-20%)'}
//         _hover={{ bg: isCurrentUser ? 'purple.500' : 'gray.600' }}
//         transition="background-color 0.2s"
//       >
//         <Text>{m.message}</Text>
//         <Text as="sub" fontSize="xs" color="gray.300" mt={2}>
//           {ts.toUTCString()}
//         </Text>
//       </Box>
//     );
//   });

//   return (
//     <Container maxW="container.md" p={4}>
//       <VStack align="stretch" spacing={3}>
//         {history}
//       </VStack>
//     </Container>
//   );
// };

// export default ChatHistory;
