import React from 'react';
import { Box, Text, VStack, Container } from '@chakra-ui/react';

const ActivityPage = ({ activities }) => {
  return (
    <Box p={4} bg="gray.900" minHeight="100vh">
      <Container maxW="container.lg">
        <VStack spacing={4} align="stretch">
          {activities.map((activity, index) => (
            <Box
              key={index}
              p={6}
              borderWidth="1px"
              borderRadius="lg"
              bg="gray.800"
              boxShadow="lg"
              transition="transform 0.2s, background-color 0.2s"
              _hover={{ bg: "gray.700", transform: "scale(1.02)" }}
            >
              <Text fontSize="lg" fontWeight="bold" color="purple.300">
                {activity.username}
              </Text>
              <Text fontSize="md" color="gray.300">
                {activity.action}
              </Text>
              <Text fontSize="sm" color="gray.400">
                {new Date(activity.timestamp * 1000).toLocaleString()}
              </Text>
              <Text fontSize="md" color="gray.300">
                {activity.details}
              </Text>
            </Box>
          ))}
        </VStack>
      </Container>
    </Box>
  );
};

export default ActivityPage;