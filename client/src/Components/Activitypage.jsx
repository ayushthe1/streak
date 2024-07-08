import React, { useEffect, useState } from 'react';
import { Box, Text, VStack } from '@chakra-ui/react';

const ActivityPage = () => {
  const [activities, setActivities] = useState([]);

  useEffect(() => {
    const ws = new WebSocket('ws://localhost:3000/ws'); // Update the WebSocket URL accordingly

    ws.onmessage = (event) => {
      const message = JSON.parse(event.data);
      if (message.type === 'activity') {
        setActivities((prevActivities) => [message, ...prevActivities]);
      }
    };

    return () => {
      ws.close();
    };
  }, []);

  return (
    <Box p={4}>
      <VStack spacing={4} align="stretch">
        {activities.map((activity, index) => (
          <Box key={index} p={4} borderWidth="1px" borderRadius="lg">
            <Text>{activity.username}</Text>
            <Text>{activity.action}</Text>
            <Text>{new Date(activity.timestamp * 1000).toLocaleString()}</Text>
            <Text>{activity.details}</Text>
          </Box>
        ))}
      </VStack>
    </Box>
  );
};

export default ActivityPage;
