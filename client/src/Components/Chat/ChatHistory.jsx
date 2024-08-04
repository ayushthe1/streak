import React from 'react';
import { Text, Box, Flex, Icon, VStack } from '@chakra-ui/react';
import { FaFileDownload } from 'react-icons/fa';

const ChatHistory = (currentUser, chats) => {
  const handleFileDownload = (url, fileName) => {
    fetch(url)
      .then(response => response.blob())
      .then(blob => {
        const url = window.URL.createObjectURL(blob);
        const a = document.createElement('a');
        a.style.display = 'none';
        a.href = url;
        a.download = fileName;
        document.body.appendChild(a);
        a.click();
        window.URL.revokeObjectURL(url);
      })
      .catch(() => alert('An error occurred while downloading the file.'));
  };

  const history = chats.map(m => {
    const isCurrentUser = m.from === currentUser;
    const alignSelf = isCurrentUser ? 'flex-end' : 'flex-start';
    const bg = isCurrentUser ? 'purple.500' : 'gray.600';
    const textAlign = isCurrentUser ? 'right' : 'left';

    const ts = new Date(m.timestamp * 1000);

    return (
      <Box
        key={m.id}
        alignSelf={alignSelf}
        maxWidth="70%"
        mb={4}
      >
        <Flex direction="column" bg={bg} p={3} borderRadius="lg" boxShadow="md">
          {/* <Text color="white" fontSize="sm" mb={1}>
            {m.from}
          </Text> */}
          <Text color="white" textAlign={textAlign} wordBreak="break-word">
            {m.message}
          </Text>
          {m.file_url && (
            <Flex align="center" mt={2} justify={isCurrentUser ? 'flex-end' : 'flex-start'}>
              <Icon
                as={FaFileDownload}
                color="white"
                cursor="pointer"
                onClick={() => handleFileDownload(m.file_url, m.file_name)}
                mr={2}
              />
              <Text fontSize="xs" color="gray.300">
                {m.file_name}
              </Text>
            </Flex>
          )}
          <Text as="sub" fontSize="xs" color="gray.300" mt={2} alignSelf={isCurrentUser ? 'flex-end' : 'flex-start'}>
            {ts.toLocaleTimeString()}
          </Text>
        </Flex>
      </Box>
    );
  });

  return <VStack spacing={4} align="stretch" w="100%">{history}</VStack>;
};

export default ChatHistory;