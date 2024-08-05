// ContactList component
import React from 'react';
import { Text, Box, Divider, VStack } from '@chakra-ui/react';

const ContactList = (contacts, sendMessage, selectedContact) => {
  const contactList = contacts.map(c => {
    const ts = new Date(c.last_activity * 1000);
    const isSelected = c.username === selectedContact;

    return (
      <Box key={c.username}>
        <Box
          as="button"
          textAlign="left"
          p={4}
          m={2}
          borderRadius="md"
          bg={isSelected ? "purple.600" : "gray.700"}
          color="white"
          boxShadow="md"
          _hover={{ bg: isSelected ? "purple.500" : "gray.600" }}
          onClick={() => sendMessage(c.username)}
          transition="background-color 0.2s"
          w={{ base: "100%", md: "80%" }}
        >
          <Text fontSize="lg" fontWeight="bold" color={isSelected ? "white" : "purple.300"}>
            {c.username}
          </Text>
          {/* <Text as="sub" fontSize="xs" color={isSelected ? "gray.200" : "gray.400"}>
            Last active: {ts.toDateString()}
          </Text> */}
        </Box>
        <Divider borderColor="gray.600" />
      </Box>
    );
  });

  return <VStack align="stretch" spacing={2}>{contactList}</VStack>;
};

export default ContactList;